package pangaea_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/copal/internal/golden"
	"github.com/vlumi/copal/jcl"
	"github.com/vlumi/copal/pangaea/compile"
	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/runner"
)

// TestDaily runs the daily batch on each fixture in testdata, whose input
// directory holds the datasets to start from, and compares every dataset the
// batch creates, and its job log and reports, with the fixture's golden
// directory. A fixture whose JOBS file lists more jobs runs them after the
// daily batch, as the scheduler would, with their job logs and reports in
// a directory of their own. Run with -update to rewrite the golden files.
func TestDaily(t *testing.T) {
	programs := t.TempDir()
	if _, err := compile.Programs("programs", "copybooks", programs); err != nil {
		t.Fatal(err)
	}
	fixtures, err := filepath.Glob("testdata/*/input")
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range fixtures {
		fixture := filepath.Dir(input)
		t.Run(filepath.Base(fixture), func(t *testing.T) {
			runFixture(t, programs, fixture)
		})
	}
}

func runFixture(t *testing.T, programs, fixture string) {
	dir := t.TempDir()
	store, err := dataset.Open(filepath.Join(dir, "datasets"))
	if err != nil {
		t.Fatal(err)
	}
	inputs, err := filepath.Glob(filepath.Join(fixture, "input", "*.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range inputs {
		if err := store.LoadFile(input); err != nil {
			t.Fatal(err)
		}
	}
	loaded := store.Names()

	sysout := filepath.Join(dir, "sysout")
	for i, job := range fixtureJobs(t, fixture) {
		r := runner.Runner{Programs: programs, Datasets: store, Sysout: sysout}
		if i > 0 {
			r.Sysout = filepath.Join(sysout, job)
		}
		if _, err := r.Run(parseJob(t, filepath.Join("jobs", job+".jcl"))); err != nil {
			t.Fatal(err)
		}
	}

	created := filepath.Join(dir, "created")
	if err := store.UnloadAll(created); err != nil {
		t.Fatal(err)
	}
	for _, dsn := range loaded {
		if err := os.Remove(filepath.Join(created, dsn+".txt")); err != nil {
			t.Fatal(err)
		}
	}
	golden.CompareDir(t, created, filepath.Join(fixture, "golden", "datasets"))
	golden.CompareDir(t, sysout, filepath.Join(fixture, "golden", "sysout"))
}

// fixtureJobs lists the jobs a fixture runs: those in its JOBS file, or
// the daily batch alone.
func fixtureJobs(t *testing.T, fixture string) []string {
	data, err := os.ReadFile(filepath.Join(fixture, "JOBS"))
	if errors.Is(err, os.ErrNotExist) {
		return []string{"DAILY"}
	} else if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(data))
}

func parseJob(t *testing.T, path string) *jcl.Job {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	job, err := jcl.Parse(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	return job
}
