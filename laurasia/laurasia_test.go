// Package laurasia is Laurasia Card Services: a card processor's COBOL, in
// Micro Focus free format, run by shell scripts. Its Go is only this test,
// which runs its fixtures and compares what they leave with their golden
// files.
package laurasia

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/copal/bedrock/run"
	"github.com/vlumi/copal/internal/golden"
)

// TestNights runs each fixture in testdata: its input copied to a data
// directory, the scripts its RUN file names run in order, and what the
// data directory then holds, but what came in and the indexed and relative
// files, which only programs read, compared with its golden files. Run
// with -update to rewrite them.
func TestNights(t *testing.T) {
	bin, err := filepath.Abs(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := (run.Build{Dialect: run.MicroFocus, Sources: "programs", Copybooks: "copybooks", Out: bin}).Programs(); err != nil {
		t.Fatal(err)
	}
	fixtures, err := filepath.Glob("testdata/*/RUN")
	if err != nil {
		t.Fatal(err)
	}
	for _, runFile := range fixtures {
		dir := filepath.Dir(runFile)
		t.Run(filepath.Base(dir), func(t *testing.T) { runFixture(t, bin, dir) })
	}
}

func runFixture(t *testing.T, bin, dir string) {
	data, err := filepath.Abs(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.CopyFS(data, os.DirFS(filepath.Join(dir, "input"))); err != nil {
		t.Fatal(err)
	}
	inputs := files(t, data)
	scripts, err := filepath.Abs("scripts")
	if err != nil {
		t.Fatal(err)
	}
	steps, err := os.ReadFile(filepath.Join(dir, "RUN"))
	if err != nil {
		t.Fatal(err)
	}
	for _, step := range strings.Split(strings.TrimSpace(string(steps)), "\n") {
		args := strings.Fields(step)
		cmd := exec.Command("sh", append([]string{filepath.Join(scripts, args[0])}, args[1:]...)...)
		cmd.Env = append(os.Environ(), "LCS_BIN="+bin, "LCS_SCRIPTS="+scripts, "LCS_DATA="+data, "COB_LIBRARY_PATH="+bin)
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		if err := cmd.Run(); err != nil {
			t.Fatalf("%s: %v\n%s", step, err, out.String())
		}
	}
	got := map[string][]byte{}
	for name, content := range files(t, data) {
		if _, in := inputs[name]; in && bytes.Equal(inputs[name], content) || isBinary(name) {
			continue
		}
		got[name] = content
	}
	golden.CompareFiles(t, got, filepath.Join(dir, "golden"))
}

// files reads every file under a directory, by its path in it.
func files(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	all := map[string][]byte{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		content, err := os.ReadFile(path)
		rel, _ := filepath.Rel(dir, path)
		all[filepath.ToSlash(rel)] = content
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return all
}

// isBinary reports whether a file is an indexed or relative file, which
// only the programs read, and the unloads show.
func isBinary(name string) bool {
	return strings.HasSuffix(name, ".idx") || strings.HasSuffix(name, ".rel")
}
