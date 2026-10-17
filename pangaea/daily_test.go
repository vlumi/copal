package pangaea_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/copal/bedrock/jcl"
	"github.com/vlumi/copal/internal/golden"
	"github.com/vlumi/copal/pangaea/compile"
	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/fixture"
	"github.com/vlumi/copal/pangaea/records"
	"github.com/vlumi/copal/pangaea/runner"
)

// TestDaily runs the daily batch on each fixture in testdata, whose input
// directory holds the datasets to start from, and compares every dataset the
// batch creates, and its job log and reports, with the fixture's golden
// directory. A fixture whose JOBS file lists more jobs runs them after the
// daily batch, as the scheduler would, with their job logs and reports in
// a directory of their own; a job's in-stream data is the fixture's
// JOB.SYSIN. A step that abends or ends with RC 16, which Pangaea's
// programs keep for what is broken, fails the test. Run with -update to
// rewrite the golden files.
func TestDaily(t *testing.T) {
	programs := t.TempDir()
	if _, err := compile.Programs("programs", "copybooks", programs); err != nil {
		t.Fatal(err)
	}
	codec := newCodec(t)
	fixtures, err := filepath.Glob("testdata/*/input")
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range fixtures {
		dir := filepath.Dir(input)
		t.Run(filepath.Base(dir), func(t *testing.T) {
			runFixture(t, codec, programs, dir)
		})
	}
}

func runFixture(t *testing.T, codec *fixture.Codec, programs, fixtureDir string) {
	dir := t.TempDir()
	store, err := dataset.Open(filepath.Join(dir, "datasets"))
	if err != nil {
		t.Fatal(err)
	}
	inputs, err := filepath.Glob(filepath.Join(fixtureDir, "input", "*"))
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range inputs {
		if err := codec.LoadFile(store, input); err != nil {
			t.Fatal(err)
		}
	}
	loaded := store.Names()

	sysout := filepath.Join(dir, "sysout")
	for i, job := range fixtureJobs(t, fixtureDir) {
		r := runner.Runner{Programs: programs, Datasets: store, Sysout: sysout}
		if i > 0 {
			r.Sysout = filepath.Join(sysout, job)
		}
		if data, err := os.ReadFile(filepath.Join(fixtureDir, job+".SYSIN")); err == nil {
			r.InStream = bytes.NewReader(data)
		} else if !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
		parsed, err := jcl.ParseFile(filepath.Join("jobs", job+".jcl"))
		if err != nil {
			t.Fatal(err)
		}
		results, err := r.Run(parsed)
		if err != nil {
			t.Fatal(err)
		}
		for _, result := range results {
			if result.Abended || result.RC >= 16 {
				t.Errorf("%s: %s, which is no outcome to keep as golden", job, result)
			}
		}
	}

	created := filepath.Join(dir, "created")
	if err := codec.UnloadAll(store, created); err != nil {
		t.Fatal(err)
	}
	for _, dsn := range loaded {
		for _, ext := range []string{".txt", fixture.FieldsExtension} {
			if err := os.Remove(filepath.Join(created, dsn+ext)); err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
		}
	}
	golden.CompareDir(t, created, filepath.Join(fixtureDir, "golden", "datasets"))
	golden.CompareDir(t, sysout, filepath.Join(fixtureDir, "golden", "sysout"))
}

// newCodec writes the datasets with packed decimal fields field by field.
func newCodec(t *testing.T) *fixture.Codec {
	l, err := records.Load("copybooks")
	if err != nil {
		t.Fatal(err)
	}
	return fixture.NewCodec(l)
}

// fixtureJobs lists the jobs a fixture runs: those in its JOBS file, or
// the daily batch alone.
func fixtureJobs(t *testing.T, dir string) []string {
	data, err := os.ReadFile(filepath.Join(dir, "JOBS"))
	if errors.Is(err, os.ErrNotExist) {
		return []string{"DAILY"}
	} else if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(data))
}
