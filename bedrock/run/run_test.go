package run

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A run binds the program's files by DD name, gives it its standard
// output, ends with its return code, and with trace code compiled in,
// writes the statements it executes to the trace file.
func TestStep(t *testing.T) {
	programs, dir := t.TempDir(), t.TempDir()
	if _, err := (Build{Sources: "testdata/files", Copybooks: "testdata/copybooks", Out: programs, Trace: true}).Programs(); err != nil {
		t.Fatal(err)
	}
	in, out, trace := filepath.Join(dir, "in"), filepath.Join(dir, "out"), filepath.Join(dir, "trace")
	if err := os.WriteFile(in, []byte("ALPHABRAVO"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	outcome, err := Step{Programs: programs, Program: "COPYIT", DDs: map[string]string{"INFILE": in, "OUTFILE": out},
		Stdout: &stdout, TraceFile: trace}.Run()
	if err != nil {
		t.Fatal(err)
	}
	if outcome != (Outcome{RC: 2}) {
		t.Errorf("outcome %+v, want return code 2", outcome)
	}
	if got, want := stdout.String(), "COPIED 2\n"; got != want {
		t.Errorf("output %q, want %q", got, want)
	}
	if copied, err := os.ReadFile(out); err != nil || string(copied) != "ALPHABRAVO" {
		t.Errorf("copied %q, %v, want both records", copied, err)
	}
	traced, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(traced), "Paragraph: MAIN-PARA") || !strings.Contains(string(traced), "WRITE") {
		t.Errorf("trace has no paragraph or WRITE:\n%s", traced)
	}
}
