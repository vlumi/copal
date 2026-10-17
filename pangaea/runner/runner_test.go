package runner

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/vlumi/copal/bedrock/jcl"
	"github.com/vlumi/copal/pangaea/dataset"
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

func TestRunGenerationsTemporaryAndAppend(t *testing.T) {
	r := newRunner(t, map[string]string{
		// WORK and NEXT write the next master from the current one, through a
		// temporary dataset.
		"WORK": `cat "$DD_IN" > "$DD_WORK"`,
		"NEXT": `cat "$DD_WORK" > "$DD_OUT"; printf X >> "$DD_OUT"`,
		"LOG":  `printf L >> "$DD_LOG"`,
		"NONE": `cat "$DD_IN" > "$DD_GONE"`,
	})
	if err := r.Datasets.Load("MASTER.G0001V00", strings.NewReader("A\n")); err != nil {
		t.Fatal(err)
	}
	job := parse(t, `//DAY JOB
//S1 EXEC PGM=WORK
//IN DD DSN=MASTER(0),DISP=SHR
//WORK DD DSN=&&WORK,DISP=NEW,DCB=(RECFM=FB,LRECL=1)
//S2 EXEC PGM=NEXT
//WORK DD DSN=&&WORK,DISP=OLD
//OUT DD DSN=MASTER(+1),DISP=NEW,DCB=(RECFM=FB,LRECL=1)
//S3 EXEC PGM=LOG
//LOG DD DSN=LOG(+1),DISP=MOD,DCB=(RECFM=FB,LRECL=1)
//S4 EXEC PGM=LOG
//LOG DD DSN=LOG(+1),DISP=MOD
//S5 EXEC PGM=NONE
//IN DD DSN=MASTER(+1),DISP=SHR
//GONE DD DUMMY
`)
	for range 2 {
		if _, err := r.Run(job); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{"LOG.G0001V00", "LOG.G0002V00", "MASTER.G0001V00", "MASTER.G0002V00", "MASTER.G0003V00"}
	if got := r.Datasets.Names(); !slices.Equal(got, want) {
		t.Errorf("cataloged %q, want %q: no temporary datasets left", got, want)
	}
	assertFile(t, r.Datasets.Path("MASTER.G0003V00"), "AXX")
	assertFile(t, r.Datasets.Path("LOG.G0002V00"), "LL")
}

func TestRunRejectsAppendWithoutRecordLength(t *testing.T) {
	r := newRunner(t, map[string]string{"OK": "exit 0"})
	if _, err := r.Run(parse(t, "//J JOB\n//S EXEC PGM=OK\n//LOG DD DSN=LOG,DISP=MOD\n")); err == nil {
		t.Error("appended to a dataset that does not exist without a DCB")
	}
}

// A group named by its base reads every generation it had when the job
// started, newest first; a concatenation reads its datasets one after the
// other.
func TestRunConcatenatesGenerations(t *testing.T) {
	r := newRunner(t, map[string]string{
		"WRITE": `printf C > "$DD_OUT"`,
		"READ":  `cat "$DD_IN" > "$DD_COPY"`,
	})
	for dsn, text := range map[string]string{"LOG.G0001V00": "A\n", "LOG.G0002V00": "B\n", "EXTRA": "D\n"} {
		if err := r.Datasets.Load(dsn, strings.NewReader(text)); err != nil {
			t.Fatal(err)
		}
	}
	job := parse(t, `//J JOB
//S1 EXEC PGM=WRITE
//OUT DD DSN=LOG(+1),DISP=NEW,DCB=(RECFM=FB,LRECL=1)
//S2 EXEC PGM=READ
//IN DD DSN=LOG,DISP=SHR
// DD DSN=LOG(+1),DISP=SHR
// DD DSN=EXTRA,DISP=SHR
//COPY DD DSN=COPY,DISP=NEW,DCB=(RECFM=FB,LRECL=1)
`)
	if _, err := r.Run(job); err != nil {
		t.Fatal(err)
	}
	assertFile(t, r.Datasets.Path("COPY"), "BACD")

	if err := r.Datasets.Load("WIDE", strings.NewReader("WW\n")); err != nil {
		t.Fatal(err)
	}
	mixed := parse(t, "//J JOB\n//S EXEC PGM=READ\n//IN DD DSN=EXTRA,DISP=SHR\n// DD DSN=WIDE,DISP=SHR\n//COPY DD DUMMY\n")
	if _, err := r.Run(mixed); err == nil {
		t.Error("concatenated datasets of different record lengths")
	}
}

// An indexed dataset is cataloged as such and is read only on its own.
func TestRunIndexed(t *testing.T) {
	r := newRunner(t, map[string]string{
		"LOAD": `printf AB > "$DD_OUT"`,
		"READ": `exit 0`,
	})
	if _, err := r.Run(parse(t, "//J JOB\n//S EXEC PGM=LOAD\n//OUT DD DSN=K(+1),DISP=NEW,RECORG=KS,LRECL=2\n")); err != nil {
		t.Fatal(err)
	}
	if !r.Datasets.Indexed("K.G0001V00") {
		t.Error("K.G0001V00 not cataloged as indexed")
	}
	if err := r.Datasets.Load("PLAIN", strings.NewReader("CD\n")); err != nil {
		t.Fatal(err)
	}
	concatenated := parse(t, "//J JOB\n//S EXEC PGM=READ\n//IN DD DSN=PLAIN,DISP=SHR\n// DD DSN=K(0),DISP=SHR\n")
	if _, err := r.Run(concatenated); err == nil {
		t.Error("concatenated an indexed dataset")
	}
}

// A step's SYSIN DD * reads the job's in-stream data as standard input.
func TestRunInStream(t *testing.T) {
	r := newRunner(t, map[string]string{"ECHO": `cat`})
	r.InStream = strings.NewReader("REQUEST\n")
	if _, err := r.Run(parse(t, "//J JOB\n//S EXEC PGM=ECHO\n//SYSIN DD *\n")); err != nil {
		t.Fatal(err)
	}
	assertFile(t, filepath.Join(r.Sysout, "S.STDOUT"), "REQUEST\n")
	var output strings.Builder
	r.InStream, r.Output = strings.NewReader("AGAIN\n"), &output
	if _, err := r.Run(parse(t, "//J JOB\n//S EXEC PGM=ECHO\n//SYSIN DD *\n")); err != nil {
		t.Fatal(err)
	}
	if output.String() != "AGAIN\n" {
		t.Errorf("output %q, want the program's", output.String())
	}
	if _, err := r.Run(parse(t, "//J JOB\n//S EXEC PGM=ECHO\n//CARDS DD *\n")); err == nil {
		t.Error("ran a step with in-stream data on CARDS")
	}
}
