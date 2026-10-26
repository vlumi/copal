package explore

import (
	"fmt"
	"html/template"
	"slices"

	"github.com/vlumi/copal/bedrock/runlog"
	"github.com/vlumi/copal/bedrock/runner"
	"github.com/vlumi/copal/survey/blueprint"
)

// runIndex is the system's run record as its machine room kept it, by
// date and job, with the jobs that run each job as a step of theirs.
type runIndex struct {
	jobs    map[string]*blueprint.Job
	byDate  map[string][]runlog.Run
	dates   []string
	callers map[string][]string
}

func newRunIndex(jobs map[string]*blueprint.Job, runs []runlog.Run) *runIndex {
	x := &runIndex{jobs: jobs, byDate: map[string][]runlog.Run{}, callers: map[string][]string{}}
	for _, r := range runs {
		if _, ok := x.byDate[r.Date]; !ok {
			x.dates = append(x.dates, r.Date)
		}
		x.byDate[r.Date] = append(x.byDate[r.Date], r)
	}
	slices.Sort(x.dates)
	for _, name := range slices.Sorted(mapKeys(jobs)) {
		for _, step := range jobs[name].Steps {
			if step.Runs != "" && !slices.Contains(x.callers[step.Runs], name) {
				x.callers[step.Runs] = append(x.callers[step.Runs], name)
			}
		}
	}
	return x
}

func mapKeys[V any](m map[string]V) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// run is a job's own record on a date, the last if it ran twice.
func (x *runIndex) run(job, date string) (runlog.Run, bool) {
	runs := x.byDate[date]
	for _, r := range slices.Backward(runs) {
		if r.Job == job {
			return r, true
		}
	}
	return runlog.Run{}, false
}

// outcomes are how a job's steps ended on a date, by step name: from the
// job's own record, or else from the record of a job that runs it, whose
// log has its steps; a step that runs another job ends as the worst of
// that job's steps did.
func (x *runIndex) outcomes(job *blueprint.Job, date string) map[string]runner.Result {
	return x.outcomesSeen(job, date, map[string]bool{})
}

func (x *runIndex) outcomesSeen(job *blueprint.Job, date string, seen map[string]bool) map[string]runner.Result {
	seen[job.Name] = true
	found := map[string]runner.Result{}
	names := map[string]bool{}
	for _, s := range job.Steps {
		names[s.Name] = true
	}
	if run, ok := x.run(job.Name, date); ok {
		for _, s := range run.Steps {
			found[s.Step] = s
		}
	} else {
		for _, caller := range x.callers[job.Name] {
			if run, ok := x.run(caller, date); ok {
				for _, s := range run.Steps {
					if names[s.Step] {
						found[s.Step] = s
					}
				}
			}
		}
	}
	for _, step := range job.Steps {
		called, ok := x.jobs[step.Runs]
		if _, done := found[step.Name]; done || !ok || seen[called.Name] {
			continue
		}
		if w, ok := worst(x.outcomesSeen(called, date, seen)); ok {
			w.Step, w.Program = step.Name, ""
			found[step.Name] = w
		}
	}
	return found
}

// worst is the step that ended worst: an abend, a step bypassed, or the
// highest return code.
func worst(results map[string]runner.Result) (runner.Result, bool) {
	var w runner.Result
	rank := -1
	for _, name := range slices.Sorted(mapKeys(results)) {
		r := results[name]
		n := r.RC
		switch {
		case r.Abended:
			n = 1 << 30
		case r.Bypassed:
			n = 1 << 29
		}
		if n > rank {
			w, rank = r, n
		}
	}
	return w, rank >= 0
}

// RunSummary is a job's run on a date, in a line: the outcome of the step
// that ended worst.
type RunSummary struct {
	Date  string
	Worst string
}

// runsOf lists the dates a job ran, newest first.
func (x *runIndex) runsOf(job *blueprint.Job) []RunSummary {
	var runs []RunSummary
	for _, date := range slices.Backward(x.dates) {
		if w, ok := worst(x.outcomes(job, date)); ok {
			runs = append(runs, RunSummary{Date: date, Worst: w.Outcome()})
		}
	}
	return runs
}

// RunView is a job's run on a date: its network with each step's outcome,
// and the steps in a table.
type RunView struct {
	Job     *blueprint.Job
	Date    string
	Network template.HTML
	Steps   []StepOutcome
}

// StepOutcome is a step of a job and how it ended, empty when the run did
// not take it.
type StepOutcome struct {
	Step    blueprint.Step
	Outcome string
}

func (s *Server) jobRun(job *blueprint.Job, date string) (RunView, error) {
	outcomes := s.runs.outcomes(job, date)
	if len(outcomes) == 0 {
		return RunView{}, fmt.Errorf("a run of %s on %s: %w", job.Name, date, errNotFound)
	}
	v := RunView{Job: job, Date: date}
	for _, step := range job.Steps {
		o := StepOutcome{Step: step}
		if r, ok := outcomes[step.Name]; ok {
			o.Outcome = r.Outcome()
		}
		v.Steps = append(v.Steps, o)
	}
	network, err := svg(blueprint.JobGraph(*job, blueprint.JobNotes{
		Step: func(step blueprint.Step) string {
			if r, ok := outcomes[step.Name]; ok {
				return r.Outcome()
			}
			return "not run"
		},
		StepLink:    func(step blueprint.Step) string { return stepLink(blueprint.StepRef{Job: job.Name, Step: step.Name}) },
		DatasetLink: func(ds string) string { return datasetLink(datasetBase(ds)) },
	}))
	v.Network = network
	return v, err
}
