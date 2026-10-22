package batchfile

import (
	"strings"

	"github.com/vlumi/copal/bedrock/script"
)

// Script reads the file as the job it runs, in the model a shell script
// reads into: its steps, the programs, batch files and moves, with the
// DD_ bindings in force at each, in the paths they are written in, with
// %NAME% as $NAME and backslashes as slashes; a variable set to a
// parameter or a variable of the environment alone, as DATE from %1 or
// DATA from %GLS_DATA%, stays a name of its own, as in a shell script.
func (f *File) Script() *script.Script {
	s := &script.Script{Name: f.Name, Comment: f.Comment, Parameters: map[string]string{}}
	vars, bindings := map[string]string{}, map[string]string{}
	value := func(name string) (string, bool) {
		if v, ok := vars[name]; ok {
			return v, true
		}
		return "$" + name, true
	}
	var last *script.Step
	for _, line := range f.Lines {
		var step script.Step
		switch line.Kind {
		case Set:
			if binding, ok := strings.CutPrefix(line.Var, "DD_"); ok {
				bindings[binding] = expand(line.Value, value)
				continue
			}
			if name, ok := parameter(line.Value); ok {
				s.Parameters[line.Var] = name
				vars[line.Var] = "$" + line.Var
				continue
			}
			vars[line.Var] = expand(line.Value, value)
			continue
		case Run:
			words := expandAll(line.Words, value)
			step = script.Step{Kind: script.Program, Command: words[0], Program: programName(words[0]), Args: words[1:],
				Bindings: map[string]string{}, Stdout: expand(line.Output, value)}
			for name, path := range bindings {
				step.Bindings[name] = path
			}
		case Call:
			words := expandAll(line.Words, value)
			step = script.Step{Kind: script.Run, Command: words[0], Program: batchName(words[0]), Args: words[1:]}
		case Move:
			words := expandAll(line.Words, value)
			step = script.Step{Kind: script.Move, Command: "move", Program: "move", Input: words[:1], Output: words[1]}
		case Test:
			if last != nil {
				last.StopsIf = &script.Test{Op: "-ge", Value: line.Level}
			}
			continue
		default:
			continue
		}
		step.Name, step.Comment, step.Line = line.Step, line.Comment, line.N
		if step.Name == "" {
			step.Name = strings.ToUpper(step.Program)
		}
		s.Steps = append(s.Steps, step)
		last = &s.Steps[len(s.Steps)-1]
	}
	return s
}

// parameter is the name a value that is one variable alone, %1 or
// %GLS_DATA%, refers to.
func parameter(value string) (string, bool) {
	switch {
	case len(value) == 2 && value[0] == '%' && value[1] >= '1' && value[1] <= '9':
		return value[1:], true
	case len(value) > 2 && value[0] == '%' && value[len(value)-1] == '%' && !strings.Contains(value[1:len(value)-1], "%"):
		return strings.ToUpper(value[1 : len(value)-1]), true
	}
	return "", false
}

func expandAll(words []string, value func(string) (string, bool)) []string {
	expanded := make([]string, len(words))
	for i, w := range words {
		expanded[i] = expand(w, value)
	}
	return expanded
}
