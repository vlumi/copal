// Package batch builds Gondwana's programs and runs its batch files: a
// fixture's, its input copied to a data directory and the batch files its
// RUN file names run in order, or one at a time; each takes the programs,
// the other batch files and the data from the environment the files
// expect, GLS_BIN, GLS_SCRIPTS and GLS_DATA.
package batch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/bedrock/batchfile"
	"github.com/vlumi/copal/bedrock/run"
)

// Build compiles Gondwana's programs into a directory, in Fujitsu
// NetCOBOL dialect. Root is the gondwana directory.
func Build(root, bin string) error {
	_, err := run.Build{Dialect: cobol.Fujitsu, Sources: filepath.Join(root, "programs"),
		Copybooks: filepath.Join(root, "copybooks"), Out: bin}.Programs()
	return err
}

// Runner runs the batch files with the programs, the batch files and the
// data they expect from the environment.
type Runner struct {
	Programs string
	Scripts  string
	Data     string
}

// Run runs a batch file by name, with its arguments, through bedrock's
// interpreter. A file that exits with an errorlevel other than 0 is an
// error.
func (r Runner) Run(file string, args ...string) error {
	dirs := []string{r.Programs, r.Scripts, r.Data}
	for i, dir := range dirs {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		dirs[i] = abs
	}
	programs, scripts, data := dirs[0], dirs[1], dirs[2]
	interpreter := batchfile.Runner{Env: map[string]string{"GLS_BIN": programs, "GLS_SCRIPTS": scripts, "GLS_DATA": data}}
	rc, err := interpreter.Run(filepath.Join(scripts, file), args...)
	if err != nil {
		return err
	}
	if rc != 0 {
		return fmt.Errorf("%s: errorlevel %d", strings.Join(append([]string{file}, args...), " "), rc)
	}
	return nil
}

// Run copies a fixture's input into the data directory and runs the steps
// of its RUN file on it, a batch file and its arguments a line, with the
// programs from bin and the batch files from root/scripts. A step that
// fails stops the run.
func Run(root, fixture, bin, data string) error {
	if err := os.CopyFS(data, os.DirFS(filepath.Join(fixture, "input"))); err != nil {
		return err
	}
	steps, err := os.ReadFile(filepath.Join(fixture, "RUN"))
	if err != nil {
		return err
	}
	r := Runner{Programs: bin, Scripts: filepath.Join(root, "scripts"), Data: data}
	for _, step := range strings.Split(strings.TrimSpace(string(steps)), "\n") {
		args := strings.Fields(step)
		if err := r.Run(args[0], args[1:]...); err != nil {
			return err
		}
	}
	return nil
}
