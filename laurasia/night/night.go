// Package night runs Laurasia's scripts on a fixture: its input copied to a
// data directory, and the scripts its RUN file names run in order, each
// taking the programs, the other scripts and the data from the environment
// the scripts expect.
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

// Run copies a fixture's input into the data directory and runs the steps
// of its RUN file on it, a script and its arguments a line, with the
// programs from bin and the scripts from root/scripts. A step that fails
// stops the run, with what it printed.
func Run(root, fixture, bin, data string) error {
	for _, dir := range []*string{&root, &bin, &data} {
		abs, err := filepath.Abs(*dir)
		if err != nil {
			return err
		}
		*dir = abs
	}
	if err := os.CopyFS(data, os.DirFS(filepath.Join(fixture, "input"))); err != nil {
		return err
	}
	steps, err := os.ReadFile(filepath.Join(fixture, "RUN"))
	if err != nil {
		return err
	}
	scripts := filepath.Join(root, "scripts")
	for _, step := range strings.Split(strings.TrimSpace(string(steps)), "\n") {
		args := strings.Fields(step)
		cmd := exec.Command("sh", append([]string{filepath.Join(scripts, args[0])}, args[1:]...)...)
		cmd.Env = append(os.Environ(), "LCS_BIN="+bin, "LCS_SCRIPTS="+scripts, "LCS_DATA="+data, "COB_LIBRARY_PATH="+bin)
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w\n%s", step, err, out.String())
		}
	}
	return nil
}
