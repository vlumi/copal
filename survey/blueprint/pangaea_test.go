package blueprint

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
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
	system, err := Analyze(append(programs, jobs...), "pangaea/copybooks", cobol.IBM)
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
	system, err := Analyze(programs, "pangaea/copybooks", cobol.IBM)
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
	for _, e := range system.Entities {
		if f, ok := e.Field("RJ-REJECT-DATE"); ok && f.Kind != KindDate {
			t.Errorf("RJ-REJECT-DATE, text filled from a date, is of kind %s", f.Kind)
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

// HOLDEXT writes a header, details and a trailer from one record laid out
// in redefinitions; each type's fields get the origins of its own writes.
func TestRecordTypes(t *testing.T) {
	t.Chdir("../..")
	system, err := Analyze([]string{"pangaea/programs/HOLDEXT.cbl"}, "pangaea/copybooks", cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	var types []string
	for _, out := range system.Programs[0].Outputs {
		if out.File != "TRAN-OUT" {
			continue
		}
		types = append(types, out.Type)
		if out.Type != "TR-TRAILER" {
			continue
		}
		for _, f := range out.Fields {
			if f.Name == "TT-AMOUNT-TOTAL" && !slices.Equal(f.Origins, []string{"file HOLD-FILE HD-AMOUNT", "file HOLD-FILE HD-FEE", "initial value of WS-AMOUNT-TOTAL"}) {
				t.Errorf("the trailer's amount total comes from %q", f.Origins)
			}
		}
	}
	if !slices.Equal(types, []string{"TR-HEADER", "TR-DETAIL", "TR-TRAILER"}) {
		t.Errorf("TRAN-OUT is written as %q", types)
	}
}

// The facts between records: TXPOST matches transactions to accounts by
// account number, AUTHSRV reads a card by the request's card number,
// TXVALID checks the trailer's totals over the details, and ATMBAL writes
// the result it is passed and reads nothing of the request.
func TestRelations(t *testing.T) {
	t.Chdir("../..")
	programs, err := filepath.Glob("pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := Analyze(programs, "pangaea/copybooks", cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]Program{}
	for _, p := range system.Programs {
		byName[p.Name] = p
	}
	hasLink := func(program, kind, from, to string) bool {
		return slices.ContainsFunc(byName[program].Links, func(l Link) bool { return l.Kind == kind && l.From.Field == from && l.To.Field == to })
	}
	if !hasLink("TXPOST", "merge", "TD-ACCT-NO", "AM-ACCT-NO") || !hasLink("AUTHSRV", "key", "RQ-CARD-NO", "CR-CARD-NO") {
		t.Errorf("links: TXPOST %+v, AUTHSRV %+v", byName["TXPOST"].Links, byName["AUTHSRV"].Links)
	}
	if want := (Total{File: "TRAN-IN", Field: "TT-AMOUNT-TOTAL", Sum: "TD-AMOUNT", Type: "TR-DETAIL"}); !slices.ContainsFunc(byName["TXVALID"].Totals, func(u Total) bool {
		u.Source, u.added = "", ""
		return u == want
	}) {
		t.Errorf("TXVALID's totals %+v", byName["TXVALID"].Totals)
	}
	for _, param := range byName["ATMBAL"].Parameters {
		for _, f := range param.Fields {
			if param.Item == "AUTH-REQUEST" && (f.Reads || f.Writes) || f.Name == "RS-RESULT" && !f.Writes {
				t.Errorf("ATMBAL's parameter %s field %+v", param.Item, f)
			}
		}
	}
}
