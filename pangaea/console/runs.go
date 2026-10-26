package console

import (
	"os"
	"path/filepath"

	"github.com/vlumi/copal/bedrock/runner"
)

// JobRun is a job a business day ran: its steps, each with its outcome and
// the datasets it had, and its reports. Dir is its sysout directory within
// the day's.
type JobRun struct {
	Name    string
	Dir     string
	Steps   []StepRun
	Reports []Report
}

// StepRun is a step as it ran: its outcome, and each dataset its DD
// statements allocated, with the records it held after the step, an
// indexed one's not counted.
type StepRun struct {
	runner.Result
	Datasets []runner.Allocation
}

// Report is a report a step wrote. One too large to show is linked to.
type Report struct {
	Name, Text string
	Large      bool
}

// largeReport is the size from which a report is linked to, not shown.
const largeReport = 64 << 10

// jobRuns reads the jobs a day ran: the daily batch in the day's sysout
// directory, and any other in a directory of its own within it.
func (l *Ledger) jobRuns(sysout string) ([]JobRun, error) {
	dirs := []string{""}
	entries, err := os.ReadDir(sysout)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if _, err := os.Stat(filepath.Join(sysout, e.Name(), runner.JobLog)); e.IsDir() && err == nil {
			dirs = append(dirs, e.Name())
		}
	}
	var runs []JobRun
	for _, dir := range dirs {
		run, err := jobRun(filepath.Join(sysout, dir))
		if err != nil {
			return nil, err
		}
		run.Dir = dir
		runs = append(runs, run)
	}
	return runs, nil
}

func jobRun(dir string) (JobRun, error) {
	name, results, err := runner.ReadJobLog(dir)
	if err != nil {
		return JobRun{}, err
	}
	allocations, err := runner.ReadMessages(dir)
	if err != nil {
		return JobRun{}, err
	}
	run := JobRun{Name: name}
	for _, r := range results {
		step := StepRun{Result: r}
		for _, a := range allocations {
			if a.Step == r.Step && a.Dataset != "" && !a.Sysout && !a.Dummy && !a.InStream {
				step.Datasets = append(step.Datasets, a)
			}
		}
		run.Steps = append(run.Steps, step)
	}
	run.Reports, err = reports(dir, results)
	return run, err
}

// reports reads the reports the steps wrote, in step order.
func reports(sysout string, steps []runner.Result) ([]Report, error) {
	var out []Report
	for _, step := range steps {
		paths, err := filepath.Glob(filepath.Join(sysout, step.Step+".*"))
		if err != nil {
			return nil, err
		}
		for _, path := range paths {
			r := Report{Name: filepath.Base(path)}
			info, err := os.Stat(path)
			if err != nil {
				return nil, err
			}
			if info.Size() >= largeReport {
				r.Large = true
			} else {
				data, err := os.ReadFile(path)
				if err != nil {
					return nil, err
				}
				r.Text = string(data)
			}
			out = append(out, r)
		}
	}
	return out, nil
}
