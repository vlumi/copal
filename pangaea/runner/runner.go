// Package runner executes Pangaea's jobs, standing in for the mainframe's
// job entry subsystem: it runs each step's program with its DD statements
// bound to datasets, bypasses steps by their COND tests, catalogs the new
// datasets and generations, deletes the temporary ones, and writes a job
// log.
package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/jcl"
)

// JobLog is the file in the sysout directory that lists each step's outcome.
const JobLog = "JOBLOG"

// Runner runs jobs against a dataset store.
type Runner struct {
	Programs string // directory of the compiled programs
	Datasets *dataset.Store
	Sysout   string // directory for the job log, reports and program output
}

// Result is the outcome of one step.
type Result struct {
	Step     string
	Program  string
	RC       int
	Bypassed bool
	Abended  bool
}

func (r Result) String() string {
	outcome := fmt.Sprintf("RC=%04d", r.RC)
	switch {
	case r.Bypassed:
		outcome = "BYPASSED"
	case r.Abended:
		outcome = "ABENDED"
	}
	return fmt.Sprintf("%-8s %-8s %s", r.Step, r.Program, outcome)
}

// Run executes the job's steps in order. A step is bypassed when its COND
// test holds for an earlier step's return code; after an abend, every
// remaining step is bypassed. Temporary datasets are deleted when the job
// ends. The error is for a job that could not run as written, such as a
// missing dataset or program, not for a step's return code.
func (r Runner) Run(job *jcl.Job) ([]Result, error) {
	if err := os.MkdirAll(r.Sysout, 0o755); err != nil {
		return nil, err
	}
	run := &jobRun{Runner: r, newest: map[string]int{}}
	for _, step := range job.Steps {
		for _, dd := range step.DDs {
			if dd.Relative {
				run.newest[dd.DSN] = r.Datasets.LatestGeneration(dd.DSN)
			}
		}
	}
	results, err := run.steps(job)
	if cleanupErr := run.deleteTemporary(); err == nil {
		err = cleanupErr
	}
	if err != nil {
		return results, err
	}
	return results, r.writeJobLog(job, results)
}

// jobRun is the state of one run of a job: the newest generation of each
// generation data group when the job started, which relative generations
// count from, and the temporary datasets it created.
type jobRun struct {
	Runner
	newest    map[string]int
	temporary []string
}

func (run *jobRun) steps(job *jcl.Job) ([]Result, error) {
	var results []Result
	for _, step := range job.Steps {
		if abended(results) || bypassed(step, results) {
			results = append(results, Result{Step: step.Name, Program: step.Program, Bypassed: true})
			continue
		}
		result, err := run.runStep(step)
		if err != nil {
			return results, fmt.Errorf("step %s: %w", step.Name, err)
		}
		results = append(results, result)
	}
	return results, nil
}

func abended(results []Result) bool {
	for _, result := range results {
		if result.Abended {
			return true
		}
	}
	return false
}

func bypassed(step jcl.Step, results []Result) bool {
	if step.Cond == nil {
		return false
	}
	for _, result := range results {
		if !result.Bypassed && step.Cond.Bypasses(result.RC) {
			return true
		}
	}
	return false
}

func (run *jobRun) runStep(step jcl.Step) (Result, error) {
	result := Result{Step: step.Name, Program: step.Program}
	env, err := run.allocate(step)
	if err != nil {
		return result, err
	}
	output, err := os.Create(filepath.Join(run.Sysout, step.Name+".STDOUT"))
	if err != nil {
		return result, err
	}
	cmd := exec.Command(filepath.Join(run.Programs, step.Program))
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = output
	cmd.Stderr = output
	runErr := cmd.Run()
	if err := closeOutput(output); err != nil {
		return result, err
	}

	var exit *exec.ExitError
	switch {
	case errors.As(runErr, &exit) && exit.Exited():
		result.RC = exit.ExitCode()
	case errors.As(runErr, &exit):
		result.Abended = true
		return result, run.discardNew(step)
	case runErr != nil:
		return result, runErr
	}
	return result, run.catalogWritten(step)
}

// dsn is the cataloged name of a DD's dataset.
func (run *jobRun) dsn(dd jcl.DD) string {
	switch {
	case dd.Temporary:
		return "TEMP." + dd.DSN
	case dd.Relative:
		return dataset.GenerationName(dd.DSN, run.newest[dd.DSN]+dd.Generation)
	}
	return dd.DSN
}

// allocate checks the step's datasets and returns the environment that binds
// each DD name to its file, the way GnuCOBOL resolves ASSIGN names.
func (run *jobRun) allocate(step jcl.Step) ([]string, error) {
	var env []string
	for _, dd := range step.DDs {
		path, err := run.allocateDD(step, dd)
		if err != nil {
			return nil, fmt.Errorf("DD %s: %w", dd.Name, err)
		}
		env = append(env, "DD_"+dd.Name+"="+path)
	}
	return env, nil
}

func (run *jobRun) allocateDD(step jcl.Step, dd jcl.DD) (string, error) {
	switch {
	case dd.Sysout:
		return filepath.Join(run.Sysout, step.Name+"."+dd.Name), nil
	case dd.Dummy:
		return os.DevNull, nil
	}
	dsn := run.dsn(dd)
	_, cataloged := run.Datasets.LRECL(dsn)
	switch {
	case dd.Disp == jcl.New && cataloged:
		return "", fmt.Errorf("%s already exists", dsn)
	case (dd.Disp == jcl.Shr || dd.Disp == jcl.Old) && !cataloged:
		return "", fmt.Errorf("%s not cataloged", dsn)
	case dd.Disp == jcl.Mod && !cataloged && dd.LRECL == 0:
		return "", fmt.Errorf("%s does not exist yet and has no DCB to create it with", dsn)
	}
	return run.Datasets.Path(dsn), nil
}

// catalogWritten catalogs the datasets the step created or appended to.
func (run *jobRun) catalogWritten(step jcl.Step) error {
	for _, dd := range step.DDs {
		if dd.Disp != jcl.New && dd.Disp != jcl.Mod {
			continue
		}
		dsn := run.dsn(dd)
		lrecl, cataloged := run.Datasets.LRECL(dsn)
		if !cataloged {
			lrecl = dd.LRECL
			if dd.Temporary {
				run.temporary = append(run.temporary, dsn)
			}
		}
		if err := run.Datasets.Catalog(dsn, lrecl); err != nil {
			return err
		}
	}
	return nil
}

func (run *jobRun) discardNew(step jcl.Step) error {
	for _, dd := range step.DDs {
		if dd.Disp == jcl.New {
			if err := os.Remove(run.Datasets.Path(run.dsn(dd))); err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
	}
	return nil
}

func (run *jobRun) deleteTemporary() error {
	for _, dsn := range run.temporary {
		if err := run.Datasets.Delete(dsn); err != nil {
			return err
		}
	}
	return nil
}

// closeOutput closes a step's captured output and removes it if the program
// wrote nothing.
func closeOutput(f *os.File) error {
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if info.Size() == 0 {
		return os.Remove(f.Name())
	}
	return nil
}

func (r Runner) writeJobLog(job *jcl.Job, results []Result) error {
	var b strings.Builder
	fmt.Fprintf(&b, "JOB %s\n", job.Name)
	for _, result := range results {
		fmt.Fprintln(&b, result)
	}
	return os.WriteFile(filepath.Join(r.Sysout, JobLog), []byte(b.String()), 0o644)
}
