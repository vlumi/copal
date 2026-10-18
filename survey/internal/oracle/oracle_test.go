package oracle

import (
	"testing"

	"github.com/vlumi/copal/survey/blueprint"
)

// The search reaches every rule of the fixture's tables that a run can
// reach: a missing header and trailer, a header of another date, totals
// that differ and totals that match, an account that exists, one that does
// not and one that is closed, and a large debit, from a base case that
// passes and the levers the program's own conditions give.
func TestSearch(t *testing.T) {
	system, err := blueprint.Analyze([]string{"testdata/BATCH.cbl"}, "testdata/copybooks")
	if err != nil {
		t.Fatal(err)
	}
	res, err := Search(system, &system.Programs[0], "testdata/copybooks", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for ti, table := range res.Units[0].Tables {
		for ri, r := range table.Rules {
			if _, ok := res.First[rule{0, ti, ri}]; !ok {
				t.Errorf("%s %s rule %d not reached: %v %v", table.Paragraph, table.Loop, ri+1, table.Conditions, r.Entries)
			}
		}
	}
}
