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
		"Step THIRD runs NOPE, which is not among the programs analyzed.",
		"Step SECOND reads &&U before any step writes it.",
		"&&T is written and never read.",
		"&&V is written and never read.",
	}
	if !slices.Equal(job.Findings, want) {
		t.Errorf("findings\n%q\nwant\n%q", job.Findings, want)
	}
	if dd := job.Steps[0].DDs[0]; dd.File != "IN-FILE" || dd.Use != "read" || !slices.Equal(dd.Datasets, []string{"A"}) {
		t.Errorf("first DD %+v", dd)
	}
}
