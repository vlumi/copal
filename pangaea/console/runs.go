package console

import (
	"html/template"
	"os"
	"path/filepath"
	"slices"

	"github.com/vlumi/copal/blueprint"
	"github.com/vlumi/copal/pangaea/runner"
)

// JobRun is a job a business day ran: its steps' outcomes, its network as
// it ran, and its reports. Dir is its sysout directory within the day's.
type JobRun struct {
	Name    string
	Dir     string
	Steps   []runner.Result
	Network template.HTML
	Reports []Report
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
		run, err := l.jobRun(filepath.Join(sysout, dir))
		if err != nil {
			return nil, err
		}
		run.Dir = dir
		runs = append(runs, run)
	}
	return runs, nil
}

func (l *Ledger) jobRun(dir string) (JobRun, error) {
	var run JobRun
	var err error
	if run.Name, run.Steps, err = runner.ReadJobLog(dir); err != nil {
		return JobRun{}, err
	}
	if run.Reports, err = reports(dir, run.Steps); err != nil {
		return JobRun{}, err
	}
	i := slices.IndexFunc(l.Jobs, func(j blueprint.Job) bool { return j.Name == run.Name })
	if i < 0 {
		return run, nil
	}
	allocations, err := runner.ReadMessages(dir)
	if err != nil {
		return JobRun{}, err
	}
	records, indexed := map[[2]string]int{}, map[[2]string]bool{}
	for _, a := range allocations {
		records[[2]string{a.Step, a.DD}] += a.Records
		indexed[[2]string{a.Step, a.DD}] = a.Indexed
	}
	outcomes := map[string]string{}
	for _, s := range run.Steps {
		outcomes[s.Step] = s.Outcome()
	}
	svg, err := blueprint.JobGraph(l.Jobs[i], blueprint.JobNotes{
		Step: func(s blueprint.Step) string { return outcomes[s.Name] },
		DD: func(s blueprint.Step, dd blueprint.DD) string {
			n, ok := records[[2]string{s.Name, dd.Name}]
			switch {
			case !ok:
				return ""
			case indexed[[2]string{s.Name, dd.Name}]:
				return "indexed"
			case n == 1:
				return "1 record"
			}
			return count(n) + " records"
		},
		ProgramLink:  func(name string) string { return "/blueprints/programs/" + name + ".md" },
		CopybookLink: func(name string) string { return "/blueprints/copybooks/" + name + ".md" },
	}).SVG()
	run.Network = template.HTML(svg) //nolint:gosec // drawn by Copal, its labels escaped
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
