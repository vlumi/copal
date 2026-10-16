package runner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Outcome is a step's outcome as the job log writes it: RC=0004, BYPASSED
// or ABENDED.
func (r Result) Outcome() string {
	switch {
	case r.Bypassed:
		return "BYPASSED"
	case r.Abended:
		return "ABENDED"
	}
	return fmt.Sprintf("RC=%04d", r.RC)
}

func parseResult(line string) (Result, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return Result{}, fmt.Errorf("%q is not a step's outcome", line)
	}
	r := Result{Step: fields[0], Program: fields[1]}
	switch outcome := fields[2]; outcome {
	case "BYPASSED":
		r.Bypassed = true
	case "ABENDED":
		r.Abended = true
	default:
		rc, ok := strings.CutPrefix(outcome, "RC=")
		n, err := strconv.Atoi(rc)
		if !ok || err != nil {
			return Result{}, fmt.Errorf("%q is not a step's outcome", line)
		}
		r.RC = n
	}
	return r, nil
}

// ReadJobLog reads the job log in a sysout directory: the job's name and
// each step's outcome.
func ReadJobLog(dir string) (string, []Result, error) {
	lines, err := readLines(filepath.Join(dir, JobLog))
	if err != nil {
		return "", nil, err
	}
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "JOB ") {
		return "", nil, fmt.Errorf("%s: no JOB line", filepath.Join(dir, JobLog))
	}
	var results []Result
	for _, line := range lines[1:] {
		r, err := parseResult(line)
		if err != nil {
			return "", nil, fmt.Errorf("%s: %w", filepath.Join(dir, JobLog), err)
		}
		results = append(results, r)
	}
	return strings.TrimPrefix(lines[0], "JOB "), results, nil
}

// Allocation is a dataset a step's DD statement was bound to, as the
// allocation messages give it, and the records it held when the step
// ended; a concatenation has one for each of its datasets. A DD bound to
// the job's output or to nothing has no dataset.
type Allocation struct {
	Step, DD  string
	Dataset   string
	Records   int
	Cataloged bool
	Sysout    bool
	Dummy     bool
}

// ReadMessages reads the allocation messages in a sysout directory.
func ReadMessages(dir string) ([]Allocation, error) {
	path := filepath.Join(dir, Messages)
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}
	var allocations []Allocation
	var step, dd string
	for _, line := range lines {
		if !strings.HasPrefix(line, "  ") {
			r, err := parseResult(line)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			step, dd = r.Step, ""
			continue
		}
		a, err := parseAllocation(line, step, dd)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		dd = a.DD
		allocations = append(allocations, a)
	}
	return allocations, nil
}

// parseAllocation reads a DD's line of the messages: its name, blank for
// the next dataset of a concatenation, then SYSOUT, DUMMY, or a dataset
// with its records or NOT CATALOGED.
func parseAllocation(line, step, previous string) (Allocation, error) {
	if step == "" || len(line) < 11 {
		return Allocation{}, fmt.Errorf("%q is not a DD of a step", line)
	}
	a := Allocation{Step: step, DD: strings.TrimSpace(line[2:10])}
	if a.DD == "" {
		a.DD = previous
	}
	switch fields := strings.Fields(line[10:]); {
	case a.DD == "":
		return Allocation{}, fmt.Errorf("%q continues no DD", line)
	case len(fields) == 1 && fields[0] == "SYSOUT":
		a.Sysout = true
	case len(fields) == 1 && fields[0] == "DUMMY":
		a.Dummy = true
	case len(fields) == 3 && fields[1] == "NOT" && fields[2] == "CATALOGED":
		a.Dataset = fields[0]
	case len(fields) == 3 && fields[2] == "RECORDS":
		n, err := strconv.Atoi(fields[1])
		if err != nil {
			return Allocation{}, fmt.Errorf("%q: %w", line, err)
		}
		a.Dataset, a.Records, a.Cataloged = fields[0], n, true
	default:
		return Allocation{}, fmt.Errorf("%q is not a DD of a step", line)
	}
	return a, nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}
