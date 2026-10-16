package blueprint

import (
	"slices"
	"testing"
)

// TestJobFindings analyzes a job that does not fit its program in every
// way Copal looks for.
func TestJobFindings(t *testing.T) {
	program := Program{Name: "COPY", Files: []File{
		{Name: "IN-FILE", Kind: "FD", Assign: "IN", Opened: []string{"INPUT"}},
		{Name: "OUT-FILE", Kind: "FD", Assign: "OUT", Opened: []string{"OUTPUT"}},
		{Name: "LOG-FILE", Kind: "FD", Assign: "LOG", Opened: []string{"EXTEND"}},
		{Name: "WORK", Kind: "SD", Assign: "SORTWK"},
		{Name: "KEYED-FILE", Kind: "FD", Assign: "KEYED", Organization: "INDEXED", Opened: []string{"OUTPUT"}},
	}}
	job, err := analyzeJob("testdata/BROKEN.jcl", []Program{program})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"Step FIRST: COPY only reads DD IN, which it creates with DISP=NEW.",
		"Step FIRST: COPY writes DD OUT, a dataset it shares with DISP=SHR.",
		"Step FIRST binds DD EXTRA, which COPY assigns no file to.",
		"Step FIRST binds no DD LOG for COPY's file LOG-FILE.",
		"Step SECOND: COPY appends to DD LOG, which is not DISP=MOD.",
		"Step SECOND: COPY's file KEYED-FILE is INDEXED, but DD KEYED creates a sequential dataset.",
		"Step THIRD runs NOPE, which is not among the programs analyzed.",
		"Step SECOND reads &&U before any step writes it.",
		"&&T is written and never read.",
		"&&V is written and never read.",
	}
	if !slices.Equal(job.Findings, want) {
		t.Errorf("findings\n%q\nwant\n%q", job.Findings, want)
	}
	if dd := job.Steps[0].DDs[0]; dd.File != "IN-FILE" || !slices.Equal(dd.Uses, []string{"read"}) || !slices.Equal(dd.Datasets, []string{"A"}) {
		t.Errorf("first DD %+v", dd)
	}
}

// A file opened both ways reads and writes its dataset; one opened I-O
// updates it, which is both.
func TestDDReadsAndWrites(t *testing.T) {
	for _, c := range []struct {
		uses          []string
		reads, writes bool
	}{
		{[]string{"read"}, true, false},
		{[]string{"write"}, false, true},
		{[]string{"read", "write"}, true, true},
		{[]string{"update"}, true, true},
		{nil, false, false},
	} {
		dd := DD{Uses: c.uses}
		if dd.reads() != c.reads || dd.writes() != c.writes {
			t.Errorf("%q: reads %v, writes %v, want %v, %v", c.uses, dd.reads(), dd.writes(), c.reads, c.writes)
		}
	}
}
