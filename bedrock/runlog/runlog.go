// Package runlog keeps a system's run record: every step its machine room
// ran, a line each, the business date, the job, the step, the program and
// the step's outcome, as a job log writes it:
//
//	2026-09-02 DAILY VALID TXVALID RC=0004
//
// The specimens log their runs each in their own way, a JCL job log, a
// script's log, a batch file's in Shift_JIS; the run record says the same
// of all of them, in one form, for whatever reads any of them. A job is
// named as its source names it: the JCL's job, the script's or the batch
// file's name. The steps of a script another one calls may be recorded
// under the caller, as its log has them.
package runlog

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/vlumi/copal/bedrock/runner"
)

// File is the run record's name in a system's directory.
const File = "RUNS"

// Run is a job a business day ran, and its steps in the order they ran.
type Run struct {
	Date  string // yyyy-mm-dd
	Job   string
	Steps []runner.Result
}

// Append adds runs to the record at path, creating it if it is not there.
func Append(path string, runs ...Run) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	for _, run := range runs {
		for _, s := range run.Steps {
			program := s.Program
			if program == "" {
				program = "-"
			}
			_, _ = fmt.Fprintf(w, "%s %s %s %s %s\n", run.Date, run.Job, s.Step, program, s.Outcome())
		}
	}
	return errors.Join(w.Flush(), f.Close())
}

// Read reads the record at path, a run for each stretch of lines of one
// date and job; a record that is not there has no runs.
func Read(path string) ([]Run, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var runs []Run
	for i, line := range strings.Split(strings.TrimRight(string(content), "\n"), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 5 {
			return nil, fmt.Errorf("%s:%d: %q is not a step's run", path, i+1, line)
		}
		step, err := runner.ParseResult(strings.Join(fields[2:], " "))
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, i+1, err)
		}
		if step.Program == "-" {
			step.Program = ""
		}
		if n := len(runs); n == 0 || runs[n-1].Date != fields[0] || runs[n-1].Job != fields[1] {
			runs = append(runs, Run{Date: fields[0], Job: fields[1]})
		}
		runs[len(runs)-1].Steps = append(runs[len(runs)-1].Steps, step)
	}
	return runs, nil
}
