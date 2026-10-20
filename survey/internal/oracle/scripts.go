package oracle

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	jobs "github.com/vlumi/copal/bedrock/runner"
	"github.com/vlumi/copal/bedrock/script"
	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/coverage"
)

// scriptTimeout is how long a script may run before it is killed.
const scriptTimeout = 2 * time.Minute

// RunScripts runs scripts as a batch, in the order given, with sh, as they
// are run: each on the directories of its environment the scripts before
// it wrote to, and on inputs taken from the cases of its steps' programs
// for the rest, as RunJobs runs jobs. A variable of the environment that
// names the directory of the programs or of the scripts a step runs names
// the oracle's own, which run each and record how it ended; any other
// names a directory of its own. A parameter of a script takes its value
// from the command line of the cleanest case of the first program it is
// passed to, which that case is consistent with, and parameters the
// scripts stand for by the same variable, as DATE, take the value of the
// first script's across the batch.
func RunScripts(system *blueprint.System, scripts []*blueprint.Job, results map[string]*Result, scratch string) ([]*Batch, error) {
	root := filepath.Join(scratch, "scripts")
	if err := os.RemoveAll(root); err != nil {
		return nil, err
	}
	sr := &scriptRun{system: system, results: results, root: root, programs: filepath.Join(scratch, "programs")}
	var plans []*scriptPlan
	for _, job := range scripts {
		plan, err := sr.plan(job)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", job.Name, err)
		}
		plans = append(plans, plan)
	}
	named := map[string]string{}
	var batches []*Batch
	for _, plan := range plans {
		for n, name := range plan.named {
			if value, ok := named[name]; ok {
				plan.args[n] = value
			} else if name != "" {
				named[name] = plan.args[n]
			}
		}
		plan.fill()
		b, err := sr.run(plan)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", plan.job.Name, err)
		}
		batches = append(batches, b)
	}
	return batches, nil
}

// scriptRun is a batch of scripts and the directories their environment
// names.
type scriptRun struct {
	system   *blueprint.System
	results  map[string]*Result
	root     string
	programs string
}

// scriptPlan is a script to run: its steps, and those of the scripts it
// runs, the directory it runs in, its environment and its arguments, and
// the variable that stands for each of its parameters, if one does.
type scriptPlan struct {
	job   *blueprint.Job
	dir   string
	steps []scriptStep
	env   map[string]string
	args  []string
	named []string
}

// unset stands for a parameter of a script until it has its value.
const unset = "\x00"

// scriptStep is a step of a script, or of a script it runs, with the
// script it is a step of and the values of its parameters.
type scriptStep struct {
	script.Step
	name   string // the step's name, after those of the steps that run its script
	script *script.Script
	values map[string]string
}

func (s scriptStep) expand(text string) string {
	return s.script.Expand(text, func(parameter string) string { return s.values[parameter] })
}

// path is a file's path as the step names it, expanded and clean, so
// that a parameter left blank, as in $DATA/work/$DATE/file, names the
// same file wherever it is written.
func (s scriptStep) path(text string) string {
	return filepath.Clean(s.expand(text))
}

// command is the name of what the step runs, as the oracle's own stands in
// for it.
func (s scriptStep) command() string {
	return filepath.Base(s.Command)
}

func (sr *scriptRun) plan(job *blueprint.Job) (*scriptPlan, error) {
	top, err := script.ParseFile(job.Source)
	if err != nil {
		return nil, err
	}
	plan := &scriptPlan{job: job, dir: filepath.Join(sr.root, job.Name)}
	if plan.env, err = sr.environment(top, plan.dir); err != nil {
		return nil, err
	}
	values := maps.Clone(plan.env)
	plan.args = make([]string, positionals(top))
	plan.named = make([]string, len(plan.args))
	for n := range plan.args {
		values[strconv.Itoa(n+1)] = unset + strconv.Itoa(n+1)
		for _, variable := range slices.Sorted(maps.Keys(top.Parameters)) {
			if top.Parameters[variable] == strconv.Itoa(n+1) && plan.named[n] == "" {
				plan.named[n] = variable
			}
		}
	}
	if plan.steps, err = sr.flatten(top, values, ""); err != nil {
		return nil, err
	}
	for n := range plan.args {
		plan.args[n] = sr.commandLine(plan.steps, unset+strconv.Itoa(n+1))
	}
	return plan, nil
}

// fill puts the values of the script's parameters in its steps' paths.
func (plan *scriptPlan) fill() {
	for n, value := range plan.args {
		for i := range plan.steps {
			for parameter, v := range plan.steps[i].values {
				plan.steps[i].values[parameter] = strings.ReplaceAll(v, unset+strconv.Itoa(n+1), value)
			}
		}
	}
}

func (sr *scriptRun) run(plan *scriptPlan) (*Batch, error) {
	b := &Batch{Job: plan.job.Name, results: sr.results, reached: map[string]map[rule]bool{}, about: scriptAbout}
	if err := b.provideFiles(plan.steps, sr, plan.dir); err != nil {
		return nil, err
	}
	before, err := readTree(filepath.Join(sr.root, "data"))
	if err != nil {
		return nil, err
	}
	traces, standIns := filepath.Join(plan.dir, "traces"), filepath.Join(plan.dir, "stand-ins")
	if err := sr.standIns(plan.steps, traces, standIns); err != nil {
		return nil, err
	}
	out, status, err := sr.execute(plan.job.Source, plan.args, plan.env, standIns)
	if err != nil {
		return nil, err
	}
	b.ended = fmt.Sprintf("The script ended with exit status %d.", status)
	if len(out) > 0 {
		b.sysout = map[string][]byte{"output": out}
	}
	if err := b.record(plan.steps, traces); err != nil {
		return nil, err
	}
	b.outputs, err = b.writtenFiles(plan.steps, sr, before, plan.dir)
	return b, err
}

// environment gives each variable of the environment a script's steps
// name a directory: the oracle's programs or scripts, for a variable that
// names the directory of what a step runs, and for any other a directory
// of its own, shared by the batch.
func (sr *scriptRun) environment(s *script.Script, dir string) (map[string]string, error) {
	env := map[string]string{}
	var walk func(s *script.Script) error
	walk = func(s *script.Script) error {
		for _, step := range s.Steps {
			if step.Kind == script.Program || step.Kind == script.Run {
				s.Expand(filepath.Dir(step.Command), func(parameter string) string {
					env[parameter] = filepath.Join(dir, "stand-ins", string(step.Kind))
					return ""
				})
			}
			for _, text := range append(append(slices.Collect(maps.Values(step.Bindings)), step.Input...), step.Output, step.Stdout) {
				s.Expand(text, func(parameter string) string {
					if _, ok := env[parameter]; !ok && !isPositional(parameter) {
						env[parameter] = filepath.Join(sr.root, "data", parameter)
					}
					return ""
				})
			}
			if step.Kind == script.Run {
				child, err := sr.parse(step.Program)
				if err != nil {
					return err
				}
				if err := walk(child); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err := walk(s); err != nil {
		return nil, err
	}
	for _, path := range env {
		if !strings.HasPrefix(path, filepath.Join(dir, "stand-ins")) {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return nil, err
			}
		}
	}
	return env, nil
}

// positionals is how many parameters a script's steps are passed or
// written in.
func positionals(s *script.Script) int {
	most := 0
	for _, step := range s.Steps {
		texts := append(append(append(slices.Collect(maps.Values(step.Bindings)), step.Input...), step.Args...), step.Output, step.Stdout)
		for _, text := range texts {
			s.Expand(text, func(parameter string) string {
				if n, err := strconv.Atoi(parameter); err == nil {
					most = max(most, n)
				}
				return ""
			})
		}
	}
	return most
}

func isPositional(parameter string) bool {
	_, err := strconv.Atoi(parameter)
	return err == nil
}

// parse parses the script of the system's that a step runs.
func (sr *scriptRun) parse(name string) (*script.Script, error) {
	i := slices.IndexFunc(sr.system.Jobs, func(j blueprint.Job) bool { return j.Kind == "script" && j.Name == name })
	if i < 0 {
		return nil, fmt.Errorf("no script %s", name)
	}
	return script.ParseFile(sr.system.Jobs[i].Source)
}

// flatten lists a script's steps, and in place of a step that runs another
// script, that step and the other script's steps after it.
func (sr *scriptRun) flatten(s *script.Script, values map[string]string, prefix string) ([]scriptStep, error) {
	var steps []scriptStep
	for _, step := range s.Steps {
		st := scriptStep{Step: step, name: prefix + step.Name, script: s, values: values}
		steps = append(steps, st)
		if step.Kind != script.Run {
			continue
		}
		child, err := sr.parse(step.Program)
		if err != nil {
			return nil, err
		}
		childValues := map[string]string{}
		for parameter, value := range values {
			if !isPositional(parameter) {
				childValues[parameter] = value
			}
		}
		for n, arg := range step.Args {
			childValues[strconv.Itoa(n+1)] = st.expand(arg)
		}
		nested, err := sr.flatten(child, childValues, st.name+".")
		if err != nil {
			return nil, err
		}
		steps = append(steps, nested...)
	}
	return steps, nil
}

// commandLine is the command line of the cleanest case of the first
// program of the steps that is passed a value.
func (sr *scriptRun) commandLine(steps []scriptStep, value string) string {
	for _, step := range steps {
		r := resultOf(sr.results, step)
		if r == nil || !slices.ContainsFunc(step.Args, func(arg string) bool { return step.expand(arg) == value }) {
			continue
		}
		if c, ok := r.cleanest(); ok && len(c[commandLine]) > 0 {
			return strings.TrimRight(string(c[commandLine][0]), " ")
		}
		return ""
	}
	return ""
}

// resultOf is the result of the program a step runs, if it runs one.
func resultOf(results map[string]*Result, step scriptStep) *Result {
	if step.Kind != script.Program {
		return nil
	}
	for _, r := range results {
		if executable(r.program.blueprint) == step.command() {
			return r
		}
	}
	return nil
}

// provideFiles writes the files the script reads that no step of it before
// writes, and no script before it wrote: each with the records of the
// cleanest case of the first program that reads it, or that reads what a
// sort or a move makes of it.
func (b *Batch) provideFiles(steps []scriptStep, sr *scriptRun, dir string) error {
	written := map[string]bool{}
	for i, step := range steps {
		switch step.Kind {
		case script.Program:
			r := resultOf(sr.results, step)
			if r == nil {
				continue
			}
			for _, dd := range slices.Sorted(maps.Keys(step.Bindings)) {
				path := step.path(step.Bindings[dd])
				f, ok := fileAssigned(r, dd)
				if !ok {
					continue
				}
				if reads(f) && !written[path] && !exists(path) {
					if err := b.provideFile(r, dd, path, sr, dir); err != nil {
						return err
					}
				}
				if writes(f) || slices.Contains(f.Opened, "I-O") {
					written[path] = true
				}
			}
		case script.Sort, script.Move:
			for _, in := range step.Input {
				path := step.path(in)
				if written[path] || exists(path) {
					continue
				}
				if r, dd, ok := sr.reader(steps[i+1:], step.path(step.Output)); ok {
					if err := b.provideFile(r, dd, path, sr, dir); err != nil {
						return err
					}
				}
			}
			written[step.path(step.Output)] = true
		}
	}
	return nil
}

// reader is the first program of the steps that reads a file, or what a
// sort or a move of it makes, and the DD name it reads it by.
func (sr *scriptRun) reader(steps []scriptStep, path string) (*Result, string, bool) {
	for i, step := range steps {
		switch step.Kind {
		case script.Program:
			r := resultOf(sr.results, step)
			if r == nil {
				continue
			}
			for dd, bound := range step.Bindings {
				if f, ok := fileAssigned(r, dd); ok && reads(f) && step.path(bound) == path {
					return r, dd, true
				}
			}
		case script.Sort, script.Move:
			if slices.ContainsFunc(step.Input, func(in string) bool { return step.path(in) == path }) {
				return sr.reader(steps[i+1:], step.path(step.Output))
			}
		}
	}
	return nil, "", false
}

func fileAssigned(r *Result, dd string) (blueprint.File, bool) {
	i := slices.IndexFunc(r.program.blueprint.Files, func(f blueprint.File) bool { return f.Assign == dd })
	if i < 0 {
		return blueprint.File{}, false
	}
	return r.program.blueprint.Files[i], true
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (b *Batch) provideFile(r *Result, dd, path string, sr *scriptRun, dir string) error {
	in, f, ok := r.program.inputFor(dd)
	if !ok {
		return nil
	}
	var records [][]byte
	if c, ok := r.cleanest(); ok {
		records = c[in.name()]
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if keyed(f) {
		loaded := filepath.Join(dir, "loaded")
		if err := os.MkdirAll(loaded, 0o755); err != nil {
			return err
		}
		staged := filepath.Join(loaded, filepath.Base(path))
		if err := r.runner.load(f, records, staged); err != nil {
			return err
		}
		if err := os.Rename(staged, path); err != nil {
			return err
		}
	} else if err := os.WriteFile(path, sequential(f, records), 0o644); err != nil {
		return err
	}
	b.inputs = append(b.inputs, batchInput{dsn: sr.fileName(path), reader: r, input: in, records: records})
	return nil
}

// fileName names a file by the variable of the environment whose
// directory holds it and its path there, without a .txt it ends in, which
// the batch's files add to every name.
func (sr *scriptRun) fileName(path string) string {
	rel, err := filepath.Rel(filepath.Join(sr.root, "data"), path)
	if err != nil {
		rel = filepath.Base(path)
	}
	return strings.TrimSuffix(filepath.ToSlash(rel), ".txt")
}

// standIns writes the oracle's own programs and scripts for those the
// steps run, and for sort and mv: each runs what it stands for, a program
// with its trace written to a file named after the process, and records
// that process, what it ran and its exit status, in the order they end.
func (sr *scriptRun) standIns(steps []scriptStep, traces, dir string) error {
	if err := os.MkdirAll(traces, 0o755); err != nil {
		return err
	}
	for _, step := range steps {
		var stood, command string
		switch step.Kind {
		case script.Program:
			stood, command = string(script.Program), strconv.Quote(filepath.Join(sr.programs, step.command()))
		case script.Run:
			i := slices.IndexFunc(sr.system.Jobs, func(j blueprint.Job) bool { return j.Kind == "script" && j.Name == step.Program })
			source, err := filepath.Abs(sr.system.Jobs[i].Source)
			if err != nil {
				return err
			}
			stood, command = string(script.Run), "sh "+strconv.Quote(source)
		default:
			path, err := exec.LookPath(step.Command)
			if err != nil {
				return err
			}
			stood, command = "path", strconv.Quote(path)
		}
		standIn := fmt.Sprintf("#!/bin/sh\nCOB_SET_TRACE=Y COB_TRACE_FILE=%q %s \"$@\"\nrc=$?\necho \"$$ %s $rc\" >> %q\nexit $rc\n",
			traces+"/$$.trace", command, step.command(), filepath.Join(traces, "steps"))
		path := filepath.Join(dir, stood, step.command())
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(standIn), 0o755); err != nil {
			return err
		}
	}
	return nil
}

// execute runs a script with sh, with its arguments and its environment,
// and sort and mv the oracle's own, and returns its output and its exit
// status.
func (sr *scriptRun) execute(source string, args []string, env map[string]string, standIns string) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), scriptTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", append([]string{source}, args...)...)
	cmd.Env = append(os.Environ(), "COB_LIBRARY_PATH="+sr.programs, "PATH="+filepath.Join(standIns, "path")+":"+os.Getenv("PATH"))
	for _, parameter := range slices.Sorted(maps.Keys(env)) {
		cmd.Env = append(cmd.Env, parameter+"="+env[parameter])
	}
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	err := cmd.Run()
	if exit, ok := err.(*exec.ExitError); ok && exit.Exited() {
		return out.Bytes(), exit.ExitCode(), nil
	}
	return out.Bytes(), 0, err
}

// record reads how each step ended from what the oracle's own recorded,
// in the order the steps ran, and replays each program's trace against
// its decision tables. A step with no record did not run.
func (b *Batch) record(steps []scriptStep, traces string) error {
	results := make([]jobs.Result, len(steps))
	ran := make([]bool, len(steps))
	for i, step := range steps {
		results[i] = jobs.Result{Step: step.name, Program: step.Program, Bypassed: true}
		if r := resultOf(b.results, step); r != nil {
			results[i].Program = r.Program
		}
	}
	data, err := os.ReadFile(filepath.Join(traces, "steps"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 3 {
			continue
		}
		i := 0
		for i < len(steps) && (ran[i] || steps[i].command() != fields[1]) {
			i++
		}
		if i == len(steps) {
			continue
		}
		ran[i] = true
		rc, _ := strconv.Atoi(fields[2])
		results[i].RC, results[i].Bypassed = rc, false
		if r := resultOf(b.results, steps[i]); r != nil {
			if err := b.replayTrace(r, filepath.Join(traces, fields[0]+".trace")); err != nil {
				return err
			}
		}
	}
	b.steps = results
	return nil
}

func (b *Batch) replayTrace(r *Result, path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	events, err := coverage.Read(bytes.NewReader(data))
	if err != nil {
		return err
	}
	if b.reached[r.Program] == nil {
		b.reached[r.Program] = map[rule]bool{}
	}
	for u, unit := range r.Units {
		for i, t := range coverage.Replay(unit.source, unit.Tables, events) {
			for j, n := range t.Reached {
				if n > 0 {
					b.reached[r.Program][rule{u, i, j}] = true
				}
			}
		}
	}
	return nil
}

// writtenFiles reads the files the script created or changed, a record a
// line: one a program's step binds as that program reads it, an indexed
// or relative one through its unloader, and any other as it is.
func (b *Batch) writtenFiles(steps []scriptStep, sr *scriptRun, before map[string][]byte, dir string) (map[string][]byte, error) {
	data := filepath.Join(sr.root, "data")
	after, err := readTree(data)
	if err != nil {
		return nil, err
	}
	outputs := map[string][]byte{}
	for rel, content := range after {
		if old, ok := before[rel]; ok && bytes.Equal(old, content) {
			continue
		}
		path := filepath.Join(data, filepath.FromSlash(rel))
		read := content
		if r, f, ok := b.binding(steps, path); ok {
			if read, err = r.runner.read(f, path, filepath.Join(dir, "unloaded")); err != nil {
				return nil, err
			}
		}
		outputs[sr.fileName(path)] = read
	}
	return outputs, nil
}

// binding is the program and file of the last step that binds a path.
func (b *Batch) binding(steps []scriptStep, path string) (*Result, blueprint.File, bool) {
	for _, step := range slices.Backward(steps) {
		r := resultOf(b.results, step)
		if r == nil {
			continue
		}
		for dd, bound := range step.Bindings {
			if step.path(bound) == path {
				if f, ok := fileAssigned(r, dd); ok {
					return r, f, true
				}
			}
		}
	}
	return nil, blueprint.File{}, false
}

const scriptAbout = "The script run with sh after the scripts before it in the batch, on what they wrote, and on inputs taken from the cases of " +
	"its steps' programs for the rest: each file it reads that no script before it, and no step of it before, wrote, with the records of the cleanest case of the first program that reads it, the first that ended with " +
	"the lowest return code, and each of its parameters with the command line of the cleanest case of the first program it is passed to. The files it wrote and its output are recorded beside."
