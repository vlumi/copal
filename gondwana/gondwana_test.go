// Package gondwana is Gondwana Life Insurance: a life insurer's COBOL, in
// Fujitsu NetCOBOL dialect and Shift_JIS, run by Windows batch files.
// This test runs its fixtures and compares what they leave with their
// golden files.
package gondwana

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/copal/gondwana/batch"
	"github.com/vlumi/copal/internal/golden"
)

// TestMonths runs each fixture in testdata: its input copied to a data
// directory, the batch files its RUN file names run in order, and what the
// data directory then holds, but what came in and the indexed master and
// the binary control record, which only programs read and whose contents
// the unload's report shows, compared with its golden files. Run with
// -update to rewrite them.
func TestMonths(t *testing.T) {
	bin := t.TempDir()
	if err := batch.Build(".", bin); err != nil {
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
	data := t.TempDir()
	if err := batch.Run(".", dir, bin, data); err != nil {
		t.Fatal(err)
	}
	inputs := files(t, filepath.Join(dir, "input"))
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

// isBinary reports whether a file is the indexed master or the binary
// control record, which only the programs read.
func isBinary(name string) bool {
	return strings.HasSuffix(name, ".IDX") || strings.HasSuffix(name, ".BIN")
}
