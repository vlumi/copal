package blueprint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/script"
)

// Kinds of job, by what it is written in, and of step, by what it runs.
const (
	KindScript = "script"

	StepSort   = "sort"
	StepMove   = "move"
	StepScript = "script"
)

// analyzeScript reads a batch script as a job: each step that runs a
// program binds the program's files to the DD_ variables exported before
// it, by the names the program assigns them; a sort or a move reads its
// input and writes its output, and passes the records on unchanged; and a
// step may run another script, a job of its own.
func analyzeScript(path string, programs []Program) (Job, error) {
	parsed, err := script.ParseFile(path)
	if err != nil {
		return Job{}, err
	}
	job := Job{Name: parsed.Name, Kind: KindScript, Source: path, Description: cobol.CommentText(parsed.Comment)}
	for _, s := range parsed.Steps {
		step := Step{Name: s.Name, Program: s.Program, Arguments: s.Args, Description: cobol.CommentText(s.Comment),
			Source: fmt.Sprintf("%s:%d", path, s.Line)}
		if s.StopsIf != nil {
			step.StopsIf = fmt.Sprintf("RC %s %d", testOperators[s.StopsIf.Op], s.StopsIf.Value)
		}
		switch s.Kind {
		case script.Program:
			step.Program = strings.ToUpper(s.Program)
			job.Findings = append(job.Findings, bindScriptStep(&step, s, programs)...)
		case script.Sort, script.Move:
			step.Kind = map[script.Kind]string{script.Sort: StepSort, script.Move: StepMove}[s.Kind]
			step.DDs = []DD{
				{Name: "input", Datasets: s.Input, Uses: []string{useRead}, Source: step.Source},
				{Name: "output", Datasets: []string{s.Output}, Uses: []string{useWrite}, Source: step.Source},
			}
			for _, k := range s.Keys {
				step.Keys = append(step.Keys, fmt.Sprintf("%d-%d", k.Start, k.End))
			}
		case script.Run:
			step.Kind, step.Program, step.Runs = StepScript, "", s.Program
		}
		job.Steps = append(job.Steps, step)
	}
	job.Datasets = jobDatasets(job)
	return job, nil
}

// testOperators write the test command's operators as relations.
var testOperators = map[string]string{"-eq": "=", "-ne": "!=", "-lt": "<", "-le": "<=", "-gt": ">", "-ge": ">="}

// bindScriptStep binds the files of a step's program to the DD_ variables
// exported before it, and its standard output to the file it is
// redirected to, which the step writes whether the program displays
// anything or not, and reports a program not analyzed and a file no DD_
// variable binds.
func bindScriptStep(step *Step, s script.Step, programs []Program) []string {
	i := slices.IndexFunc(programs, func(p Program) bool { return p.Name == step.Program })
	if i < 0 {
		return []string{fmt.Sprintf("Step %s runs %s, which is not among the programs analyzed.", step.Name, step.Program)}
	}
	p := programs[i]
	var findings []string
	for _, f := range p.Files {
		path, ok := s.Bindings[f.Assign]
		if !ok {
			findings = append(findings, fmt.Sprintf("Step %s exports no DD_%s for %s's file %s.", step.Name, f.Assign, p.Name, f.Name))
			continue
		}
		dd := DD{Name: f.Assign, Datasets: []string{path}, File: f.Name, Copybook: fileCopybook(p, f), Source: step.Source}
		for _, mode := range f.Opened {
			dd.Uses = append(dd.Uses, uses[mode])
		}
		dd.Indexed = f.Organization == "INDEXED"
		step.DDs = append(step.DDs, dd)
	}
	if s.Stdout != "" {
		step.DDs = append(step.DDs, DD{Name: "stdout", Datasets: []string{s.Stdout}, Uses: []string{useWrite}, Source: step.Source})
	}
	return findings
}
