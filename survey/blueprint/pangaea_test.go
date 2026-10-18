package blueprint

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/vlumi/copal/internal/golden"
)

// TestPangaea generates the blueprints for Pangaea and compares them with
// the golden ones in blueprints/pangaea. Run with -update to rewrite them.
func TestPangaea(t *testing.T) {
	t.Chdir("../..")
	programs, err := filepath.Glob("pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := filepath.Glob("pangaea/jobs/*.jcl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := Analyze(append(programs, jobs...), "pangaea/copybooks")
	if err != nil {
		t.Fatal(err)
	}
	files, err := system.Render("blueprints/pangaea", "pangaea")
	if err != nil {
		t.Fatal(err)
	}
	golden.CompareFiles(t, files, "blueprints/pangaea")
}

// One business date in three layouts: eight digits in the bank control
// record, six on the maintenance requests, which open accounts dated by
// them, and text on the rejects, strung from the business date.
func TestDateLayouts(t *testing.T) {
	t.Chdir("../..")
	programs, err := filepath.Glob("pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := Analyze(programs, "pangaea/copybooks")
	if err != nil {
		t.Fatal(err)
	}
	held := map[string]string{}
	for _, l := range DateLayouts(system.Entities) {
		for _, f := range l.Fields {
			held[f.Field] = l.Layout
		}
	}
	for field, layout := range map[string]string{"BC-BUSINESS-DATE": "9(8)", "MR-EFFECTIVE-DATE": "9(6)", "RJ-REJECT-DATE": "X(10)"} {
		if got := held[field]; got != layout {
			t.Errorf("%s is held as %q, want %q", field, got, layout)
		}
	}
	for record, want := range map[string]FieldRef{
		"REJECT-RECORD":      {Record: "BANK-CONTROL-RECORD", Field: "BC-BUSINESS-DATE"},
		"ACCT-MASTER-RECORD": {Record: "MAINT-REQUEST", Field: "MR-EFFECTIVE-DATE"},
	} {
		found := false
		for _, e := range system.Entities {
			for _, f := range e.Fields {
				found = found || e.Record == record && slices.Contains(f.DatesFrom, want)
			}
		}
		if !found {
			t.Errorf("no field of %s takes its dates from %v", record, want)
		}
	}
}
