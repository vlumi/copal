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
	"slices"
	"strings"

	"github.com/vlumi/copal/jcl"
	"github.com/vlumi/copal/pangaea/dataset"
)

// JobLog is the file in the sysout directory that lists each step's outcome.
const JobLog = "JOBLOG"

// Messages is the file in the sysout directory that lists, for each step
// that ran, the dataset each DD statement was bound to and the records it
// held when the step ended, as a mainframe's job entry subsystem does in
// its JESYSMSG.
const Messages = "JESYSMSG"

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
			for _, d := range append([]jcl.DD{dd}, dd.Concatenated...) {
				if !d.Temporary {
					run.newest[d.DSN] = r.Datasets.LatestGeneration(d.DSN)
				}
			}
		}
	}
	results, err := run.steps(job)
	if cleanupErr := run.deleteTemporary(); err == nil {
		err = cleanupErr
	}
	if run.scratch != "" {
		if cleanupErr := os.RemoveAll(run.scratch); err == nil {
			err = cleanupErr
		}
	}
	if err != nil {
		return results, err
	}
	if err := os.WriteFile(filepath.Join(r.Sysout, Messages), []byte(run.messages.String()), 0o644); err != nil {
		return results, err
	}
	return results, r.writeJobLog(job, results)
}

// recordAllocations writes a step's outcome to the messages, and each of
// its DDs with the datasets it was bound to and their records.
func (run *jobRun) recordAllocations(step jcl.Step, result Result) {
	fmt.Fprintln(&run.messages, result)
	for _, dd := range step.DDs {
		name := dd.Name
		switch {
		case dd.Sysout:
			fmt.Fprintf(&run.messages, "  %-8s SYSOUT\n", name)
			continue
		case dd.Dummy:
			fmt.Fprintf(&run.messages, "  %-8s DUMMY\n", name)
			continue
		}
		for _, d := range append([]jcl.DD{dd}, dd.Concatenated...) {
			dsns := []string{run.dsn(d)}
			if d.Disp == jcl.Shr || d.Disp == jcl.Old {
				if inputs, err := run.inputs(d); err == nil {
					dsns = inputs
				}
			}
			for _, dsn := range dsns {
				fmt.Fprintf(&run.messages, "  %-8s %-28s %s\n", name, dsn, run.records(dsn))
				name = ""
			}
		}
	}
}

// records says how many records a dataset holds.
func (run *jobRun) records(dsn string) string {
	lrecl, ok := run.Datasets.LRECL(dsn)
	info, err := os.Stat(run.Datasets.Path(dsn))
	if !ok || err != nil {
		return "NOT CATALOGED"
	}
	return fmt.Sprintf("%d RECORDS", info.Size()/int64(lrecl))
}

// jobRun is the state of one run of a job: the newest generation of each
// generation data group when the job started, which relative generations
// count from, the temporary datasets it created, and a scratch directory
// for the files that concatenated datasets are read from.
type jobRun struct {
	Runner
	newest    map[string]int
	temporary []string
	scratch   string
	messages  strings.Builder
}

func (run *jobRun) steps(job *jcl.Job) ([]Result, error) {
	var results []Result
	for _, step := range job.Steps {
		if abended(results) || bypassed(step, results) {
			result := Result{Step: step.Name, Program: step.Program, Bypassed: true}
			results = append(results, result)
			fmt.Fprintln(&run.messages, result)
			continue
		}
		result, err := run.runStep(step)
		if err != nil {
			return results, fmt.Errorf("step %s: %w", step.Name, err)
		}
		results = append(results, result)
		run.recordAllocations(step, result)
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
	if dd.Disp == jcl.Shr || dd.Disp == jcl.Old {
		return run.allocateInput(step, dd)
	}
	dsn := run.dsn(dd)
	_, cataloged := run.Datasets.LRECL(dsn)
	switch {
	case dd.Disp == jcl.New && cataloged:
		return "", fmt.Errorf("%s already exists", dsn)
	case dd.Disp == jcl.Mod && !cataloged && dd.LRECL == 0:
		return "", fmt.Errorf("%s does not exist yet and has no DCB to create it with", dsn)
	}
	return run.Datasets.Path(dsn), nil
}

// allocateInput returns the file a program reads a DD from: the dataset's
// own, or, for a concatenation or every generation of a group, a scratch
// file of all their records in order.
func (run *jobRun) allocateInput(step jcl.Step, dd jcl.DD) (string, error) {
	var members []string
	for _, d := range append([]jcl.DD{dd}, dd.Concatenated...) {
		dsns, err := run.inputs(d)
		if err != nil {
			return "", err
		}
		members = append(members, dsns...)
	}
	if len(members) == 1 {
		return run.Datasets.Path(members[0]), nil
	}
	return run.concatenate(step.Name+"."+dd.Name, members)
}

// inputs lists the datasets a DD reads: the one it names, or every
// generation of the group it names, newest first, as of the job's start.
func (run *jobRun) inputs(dd jcl.DD) ([]string, error) {
	dsn := run.dsn(dd)
	if _, cataloged := run.Datasets.LRECL(dsn); cataloged {
		return []string{dsn}, nil
	}
	var dsns []string
	if !dd.Relative && !dd.Temporary {
		for _, n := range slices.Backward(run.Datasets.Generations(dd.DSN)) {
			if n <= run.newest[dd.DSN] {
				dsns = append(dsns, dataset.GenerationName(dd.DSN, n))
			}
		}
	}
	if len(dsns) == 0 {
		return nil, fmt.Errorf("%s not cataloged", dsn)
	}
	return dsns, nil
}

// concatenate writes the records of the datasets, which must have one
// record length, to one scratch file.
func (run *jobRun) concatenate(name string, dsns []string) (string, error) {
	lrecl, _ := run.Datasets.LRECL(dsns[0])
	var data []byte
	for _, dsn := range dsns {
		if n, _ := run.Datasets.LRECL(dsn); n != lrecl {
			return "", fmt.Errorf("%s has records of %d bytes, %s of %d", dsn, n, dsns[0], lrecl)
		}
		member, err := os.ReadFile(run.Datasets.Path(dsn))
		if err != nil {
			return "", err
		}
		data = append(data, member...)
	}
	if run.scratch == "" {
		var err error
		if run.scratch, err = os.MkdirTemp("", "pangaea-concatenation-"); err != nil {
			return "", err
		}
	}
	path := filepath.Join(run.scratch, name)
	return path, os.WriteFile(path, data, 0o644)
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
