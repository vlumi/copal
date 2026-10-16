package console

import (
	"bufio"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/blueprint"
	"github.com/vlumi/copal/pangaea/runner"
)

// JobRun is a job a business day ran: its steps' outcomes, its network as
// it ran, and its reports. Dir is its sysout directory within the day's.
type JobRun struct {
	Name    string
	Dir     string
	Steps   []Step
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
	data, err := os.ReadFile(filepath.Join(dir, runner.JobLog))
	if err != nil {
		return JobRun{}, err
	}
	run := JobRun{Name: strings.TrimPrefix(strings.SplitN(string(data), "\n", 2)[0], "JOB ")}
	if run.Steps, err = jobSteps(dir); err != nil {
		return JobRun{}, err
	}
	if run.Reports, err = reports(dir, run.Steps); err != nil {
		return JobRun{}, err
	}
	i := slices.IndexFunc(l.Jobs, func(j blueprint.Job) bool { return j.Name == run.Name })
	if i < 0 {
		return run, nil
	}
	records, err := allocatedRecords(dir)
	if err != nil {
		return JobRun{}, err
	}
	outcomes := map[string]string{}
	for _, s := range run.Steps {
		outcomes[s.Name] = s.Outcome
	}
	svg, err := blueprint.JobGraph(l.Jobs[i], blueprint.JobNotes{
		Step: func(s blueprint.Step) string { return outcomes[s.Name] },
		DD: func(s blueprint.Step, dd blueprint.DD) string {
			if n, ok := records[s.Name][dd.Name]; ok {
				return count(n) + " records"
			}
			return ""
		},
		ProgramLink:  func(name string) string { return "/blueprints/programs/" + name + ".md" },
		CopybookLink: func(name string) string { return "/blueprints/copybooks/" + name + ".md" },
	}).SVG()
	run.Network = template.HTML(svg) //nolint:gosec // drawn by Copal, its labels escaped
	return run, err
}

// allocatedRecords reads the records each step's DDs held, from the
// runner's allocation messages: by step, by DD, all its datasets together.
func allocatedRecords(dir string) (map[string]map[string]int, error) {
	f, err := os.Open(filepath.Join(dir, runner.Messages))
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	records := map[string]map[string]int{}
	var step, dd string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, " ") {
			step = strings.Fields(line)[0]
			records[step] = map[string]int{}
			continue
		}
		fields := strings.Fields(line)
		if !strings.HasPrefix(line, "   ") {
			dd = fields[0]
		}
		if len(fields) >= 2 && fields[len(fields)-1] == "RECORDS" {
			n, err := strconv.Atoi(fields[len(fields)-2])
			if err != nil {
				return nil, err
			}
			records[step][dd] += n
		}
	}
	return records, scanner.Err()
}

// reports reads the reports the steps wrote, in step order.
func reports(sysout string, steps []Step) ([]Report, error) {
	var out []Report
	for _, step := range steps {
		paths, err := filepath.Glob(filepath.Join(sysout, step.Name+".*"))
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
