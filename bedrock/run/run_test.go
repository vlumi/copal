package run

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A run binds the program's files by DD name, gives it its standard
// output, ends with its return code, and with trace code compiled in,
// writes the statements it executes to the trace file.
func TestStep(t *testing.T) {
	programs, dir := t.TempDir(), t.TempDir()
	if _, err := (Build{Dialect: IBM, Sources: "testdata/files", Copybooks: "testdata/copybooks", Out: programs, Trace: true}).Programs(); err != nil {
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

// A program that runs past its timeout is killed, and the outcome says so.
func TestStepTimeout(t *testing.T) {
	programs := t.TempDir()
	script := filepath.Join(programs, "LOOP")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nsleep 5\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	outcome, err := Step{Programs: programs, Program: "LOOP", Stdout: io.Discard, Timeout: 100 * time.Millisecond}.Run()
	if err != nil {
		t.Fatal(err)
	}
	if outcome != (Outcome{Abended: true, TimedOut: true}) {
		t.Errorf("outcome %+v, want killed for running out of time", outcome)
	}
}
