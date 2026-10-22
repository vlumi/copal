package batchfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vlumi/copal/bedrock/run"
)

// Runner runs batch files as cmd.exe would, with the variables of its
// environment, such as GLS_BIN, GLS_SCRIPTS and GLS_DATA, which the files
// take their programs, each other and their data from.
type Runner struct {
	Env map[string]string
}

// Run runs a batch file with its arguments and returns its errorlevel:
// what it exits with, or the last step's return code when it runs off
// its end. The error is for a file that cannot be read or a step that
// cannot be run at all, not for a return code.
func (r Runner) Run(path string, args ...string) (int, error) {
	vars := map[string]string{}
	for name, v := range r.Env {
		vars[strings.ToUpper(name)] = v
	}
	return r.run(path, args, vars)
}

// run runs a file with the variables in force where it is called, which
// it sets for itself alone, as a file that starts with setlocal does.
func (r Runner) run(path string, args []string, outer map[string]string) (int, error) {
	file, err := ParseFile(path)
	if err != nil {
		return 0, err
	}
	vars := map[string]string{}
	for name, v := range outer {
		vars[name] = v
	}
	rc := 0
	value := func(name string) (string, bool) {
		switch {
		case name == "ERRORLEVEL":
			return strconv.Itoa(rc), true
		case len(name) == 1 && name[0] >= '1' && name[0] <= '9':
			i := int(name[0] - '1')
			if i < len(args) {
				return args[i], true
			}
			return "", true
		}
		v, ok := vars[name]
		return v, ok
	}
	fail := func(line Line, err error) (int, error) {
		return 0, fmt.Errorf("%s line %d: %w", path, line.N, err)
	}
	for _, line := range file.Lines {
		switch line.Kind {
		case Set:
			vars[line.Var] = expand(line.Value, value)
		case Mkdir:
			if err := os.MkdirAll(expand(line.Words[0], value), 0o755); err != nil {
				return fail(line, err)
			}
		case Echo:
			if err := appendLine(expand(line.Output, value), expand(line.Text, value)); err != nil {
				return fail(line, err)
			}
		case Run:
			if rc, err = r.program(expandAll(line.Words, value), expand(line.Output, value), vars); err != nil {
				return fail(line, err)
			}
		case Call:
			words := expandAll(line.Words, value)
			if rc, err = r.run(words[0], words[1:], vars); err != nil {
				return 0, err
			}
		case Move:
			words := expandAll(line.Words, value)
			if err := os.Rename(words[0], words[1]); err != nil {
				return fail(line, err)
			}
		case Test:
			if rc >= line.Level {
				return code(line.Code, value), nil
			}
		case Exit:
			return code(line.Code, value), nil
		}
	}
	return rc, nil
}

// program runs a program with its arguments, its standard output and
// error in the output file, and its DD names bound to what the DD_
// variables name, and returns its return code; a program that abends
// returns 255, as a process that is killed does.
func (r Runner) program(words []string, output string, vars map[string]string) (int, error) {
	step := run.Step{Programs: filepath.Dir(words[0]), Program: programName(words[0]), Args: words[1:], DDs: map[string]string{}}
	for name, v := range vars {
		if dd, ok := strings.CutPrefix(name, "DD_"); ok {
			step.DDs[dd] = v
		}
	}
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return 0, err
		}
		defer func() { _ = f.Close() }()
		step.Stdout = f
	}
	outcome, err := step.Run()
	switch {
	case err != nil:
		return 0, err
	case outcome.Abended:
		return 255, nil
	}
	return outcome.RC, nil
}

// code is the errorlevel an exit gives: the number written, or the
// variable's value.
func code(text string, value func(string) (string, bool)) int {
	n, _ := strconv.Atoi(expand(text, value))
	return n
}

// appendLine appends a line to a file, creating it if need be.
func appendLine(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(f, text); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
