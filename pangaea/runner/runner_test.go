package runner

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/jcl"
)

// newRunner sets up a runner whose programs are shell scripts.
func newRunner(t *testing.T, programs map[string]string) Runner {
	t.Helper()
	dir := t.TempDir()
	programDir := filepath.Join(dir, "programs")
	if err := os.Mkdir(programDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, script := range programs {
		path := filepath.Join(programDir, name)
		if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script+"\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	store, err := dataset.Open(filepath.Join(dir, "datasets"))
	if err != nil {
		t.Fatal(err)
	}
	return Runner{Programs: programDir, Datasets: store, Sysout: filepath.Join(dir, "sysout")}
}

func parse(t *testing.T, text string) *jcl.Job {
	t.Helper()
	job, err := jcl.Parse(strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	return job
}

func outcomes(results []Result) []string {
	var lines []string
	for _, result := range results {
		lines = append(lines, result.String())
	}
	return lines
}

func TestRunBindsDatasetsAndCatalogsNewOnes(t *testing.T) {
	r := newRunner(t, map[string]string{
		"COPY": `cat "$DD_IN" > "$DD_OUT"; echo copied > "$DD_REPORT"; echo done`,
	})
	if err := r.Datasets.Load("A", strings.NewReader("1234\n5678\n")); err != nil {
		t.Fatal(err)
	}
	job := parse(t, `//J JOB
//S EXEC PGM=COPY
//IN DD DSN=A,DISP=SHR
//OUT DD DSN=B,DISP=NEW,DCB=(RECFM=FB,LRECL=4)
//REPORT DD SYSOUT=*
`)
	results, err := r.Run(job)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := outcomes(results), []string{"S        COPY     RC=0000"}; !slices.Equal(got, want) {
		t.Errorf("results %q, want %q", got, want)
	}
	if lrecl, ok := r.Datasets.LRECL("B"); !ok || lrecl != 4 {
		t.Errorf("B cataloged with LRECL %d, %v, want 4", lrecl, ok)
	}
	assertFile(t, filepath.Join(r.Sysout, "S.REPORT"), "copied\n")
	assertFile(t, filepath.Join(r.Sysout, "S.STDOUT"), "done\n")
	assertFile(t, filepath.Join(r.Sysout, JobLog), "JOB J\nS        COPY     RC=0000\n")
}

func TestRunBypassesByCond(t *testing.T) {
	r := newRunner(t, map[string]string{
		"WARN": "exit 4",
		"FAIL": "exit 12",
		"OK":   "exit 0",
	})
	job := parse(t, `//J JOB
//S1 EXEC PGM=WARN
//S2 EXEC PGM=OK,COND=(4,LT)
//S3 EXEC PGM=FAIL
//S4 EXEC PGM=OK,COND=(4,LT)
//S5 EXEC PGM=OK
`)
	results, err := r.Run(job)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"S1       WARN     RC=0004",
		"S2       OK       RC=0000",
		"S3       FAIL     RC=0012",
		"S4       OK       BYPASSED",
		"S5       OK       RC=0000",
	}
	if got := outcomes(results); !slices.Equal(got, want) {
		t.Errorf("results\n%q\nwant\n%q", got, want)
	}
}

func TestRunBypassesEverythingAfterAnAbend(t *testing.T) {
	r := newRunner(t, map[string]string{
		"CRASH": `echo partial > "$DD_OUT"; kill -SEGV $$`,
		"OK":    "exit 0",
	})
	job := parse(t, `//J JOB
//S1 EXEC PGM=CRASH
//OUT DD DSN=B,DISP=NEW,DCB=(RECFM=FB,LRECL=8)
//S2 EXEC PGM=OK
`)
	results, err := r.Run(job)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"S1       CRASH    ABENDED", "S2       OK       BYPASSED"}
	if got := outcomes(results); !slices.Equal(got, want) {
		t.Errorf("results %q, want %q", got, want)
	}
	if _, err := os.Stat(r.Datasets.Path("B")); !os.IsNotExist(err) {
		t.Error("the abended step's new dataset was kept")
	}
}

func TestRunChecksDispositions(t *testing.T) {
	r := newRunner(t, map[string]string{"OK": "exit 0"})
	if err := r.Datasets.Load("A", strings.NewReader("1234\n")); err != nil {
		t.Fatal(err)
	}
	tests := map[string]string{
		"missing input": "//IN DD DSN=MISSING,DISP=SHR\n",
		"existing new":  "//OUT DD DSN=A,DISP=NEW,DCB=(RECFM=FB,LRECL=4)\n",
	}
	for name, dd := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := r.Run(parse(t, "//J JOB\n//S EXEC PGM=OK\n"+dd)); err == nil {
				t.Error("ran without error")
			}
		})
	}
}

func assertFile(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Errorf("%s: %q, want %q", filepath.Base(path), data, want)
	}
}
