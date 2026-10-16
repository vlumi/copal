// Package golden compares generated files with checked-in golden files.
package golden

import (
	"bytes"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the golden files")

// CompareDir checks that gotDir and goldenDir hold the same files, at any
// depth, with the same contents. With -update it makes goldenDir a copy of
// gotDir instead.
func CompareDir(t *testing.T, gotDir, goldenDir string) {
	t.Helper()
	got := readTree(t, gotDir)
	if *update {
		write(t, goldenDir, got)
		return
	}
	want := readTree(t, goldenDir)
	for name, data := range got {
		golden, ok := want[name]
		switch {
		case !ok:
			t.Errorf("%s: not in %s", name, goldenDir)
		case !bytes.Equal(data, golden):
			t.Errorf("%s differs from %s:\n%s", name, goldenDir, data)
		}
	}
	for name := range want {
		if _, ok := got[name]; !ok {
			t.Errorf("%s: in %s but not generated", name, goldenDir)
		}
	}
}

func readTree(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	files := map[string][]byte{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		name, err := filepath.Rel(dir, path)
		files[name] = data
		return err
	})
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return files
}

func write(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	for name, data := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
