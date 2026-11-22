package oracle

import (
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/survey/blueprint"
)

// The search reaches every rule of the fixture's tables that a run can
// reach: a missing header and trailer, a header of another date, totals
// that differ and totals that match, an account that exists, one that does
// not and one that is closed, a large debit, and files that cannot be
// opened or read, from a base case that passes and the levers the
// program's own conditions and files give. The one rule it cannot reach,
// a count the program keeps twice found to differ, is the one it marks as
// totals the program adds up itself.
func TestSearch(t *testing.T) {
	system, err := blueprint.Analyze([]string{"testdata/BATCH.cbl"}, "testdata/copybooks", cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Search(system, &system.Programs[0], "testdata/copybooks", t.TempDir(), nil)
	if err != nil {
		t.Fatal(err)
	}
	own := 0
	for ti, table := range res.Units[0].Tables {
		for ri, r := range table.Rules {
			_, reached := res.First[rule{0, ti, ri}]
			marked := res.Units[0].checksOwnTotals(table, r)
			if marked {
				own++
			}
			if reached == marked {
				t.Errorf("%s %s rule %d reached %v, marked as own totals %v: %v %v", table.Paragraph, table.Loop, ri+1, reached, marked, table.Conditions, r.Entries)
			}
		}
	}
	if own != 1 {
		t.Errorf("%d rules marked as own totals, want 1", own)
	}
}

// checkUnreachable fails if a case reaches a rule marked as one no run of
// the program can reach: taking a parameter to be what no call passes, or
// in a paragraph performed only from such rules; or reaches a rule marked
// as finding a file to fail once open without failing a file itself.
func checkUnreachable(t *testing.T, r *Result) {
	t.Helper()
	marks := r.program.marks()
	for reached, first := range r.First {
		why := marks[reached]
		_, _, fails := r.Cases[first].failedFile()
		if why == unpassed || why == unperformed || why == failedLater && !fails {
			table := r.Units[reached.unit].Tables[reached.table]
			t.Errorf("%s: %s %s rule %d is reached, and marked as %s", r.Program, r.Units[reached.unit].Program, table.Paragraph, reached.rule+1, why)
		}
	}
}
