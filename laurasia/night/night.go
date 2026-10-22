// Package night builds Laurasia's programs and runs its scripts: a
// fixture's, its input copied to a data directory and the scripts its RUN
// file names run in order, or one script at a time; each takes the
// programs, the other scripts and the data from the environment the
// scripts expect.
package night

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/bedrock/run"
)

// Build compiles Laurasia's programs into a directory, in Micro Focus
// dialect. Root is the laurasia directory.
func Build(root, bin string) error {
	_, err := run.Build{Dialect: cobol.MicroFocus, Sources: filepath.Join(root, "programs"),
		Copybooks: filepath.Join(root, "copybooks"), Out: bin}.Programs()
	return err
}

// Runner runs the scripts with the programs, the scripts and the data they
// expect from the environment.
type Runner struct {
	Programs string
	Scripts  string
	Data     string
}

// Run runs a script by name, with its arguments, through sh. A script
// that fails returns what it printed with the error.
func (r Runner) Run(script string, args ...string) error {
	dirs := []string{r.Programs, r.Scripts, r.Data}
	for i, dir := range dirs {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		dirs[i] = abs
	}
	programs, scripts, data := dirs[0], dirs[1], dirs[2]
	cmd := exec.Command("sh", append([]string{filepath.Join(scripts, script)}, args...)...)
	cmd.Env = append(os.Environ(), "LCS_BIN="+programs, "LCS_SCRIPTS="+scripts, "LCS_DATA="+data, "COB_LIBRARY_PATH="+programs)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w\n%s", strings.Join(append([]string{script}, args...), " "), err, out.String())
	}
	return nil
}

// Run copies a fixture's input into the data directory and runs the steps
// of its RUN file on it, a script and its arguments a line, with the
// programs from bin and the scripts from root/scripts. A step that fails
// stops the run, with what it printed.
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
