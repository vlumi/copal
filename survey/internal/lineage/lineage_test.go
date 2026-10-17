package lineage

import (
	"fmt"
	"path/filepath"
	"slices"
	"testing"

	"github.com/vlumi/copal/survey/blueprint"
)

func pangaea(t *testing.T) *blueprint.System {
	t.Helper()
	programs, err := filepath.Glob("../../../pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := filepath.Glob("../../../pangaea/jobs/*.jcl")
	if err != nil {
		t.Fatal(err)
	}
	s, err := blueprint.Analyze(append(programs, jobs...), "../../../pangaea/copybooks")
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// A statement's total in the month-end report goes back through the
// postings of the daily batch, its sort and its validation, to the files
// that come from outside: the branches' transactions and the card file of
// the online day.
func TestBackFromReport(t *testing.T) {
	trail := Back(pangaea(t), Node{Program: "STMTPRT", Field: "WS-RA-VALUE"})
	hops := map[Node]Hop{}
	for _, h := range trail.Hops {
		hops[h.Node] = h
	}
	if first := trail.Hops[0].Node; first != (Node{Job: "MONTHLY", Step: "STMT", Program: "STMTPRT", Field: "WS-RA-VALUE"}) {
		t.Errorf("trail starts at %+v", first)
	}
	post, ok := hops[Node{Job: "DAILY", Step: "POST", Program: "TXPOST", Field: "AU-AMOUNT"}]
	if !ok {
		t.Fatal("no hop at the posting of the day's transactions")
	}
	if got := fmt.Sprint(post.Sources[0].Datasets, post.Sources[0].Writers); got !=
		"[&&SORTED &&VALID &&CARDVAL] [{DAILY VALID TXVALID TD-AMOUNT 0 0} {DAILY VALIDC TXVALID TD-AMOUNT 0 0}]" {
		t.Errorf("the posted amount comes, through the sort, from %s", got)
	}
	valid := hops[Node{Job: "DAILY", Step: "VALID", Program: "TXVALID", Field: "TD-AMOUNT"}]
	if len(valid.Sources) != 1 || fmt.Sprint(valid.Sources[0].External) != "[TRAN.DAILY(0)]" {
		t.Errorf("the branches' amounts come from %+v, want TRAN.DAILY(0) from outside", valid.Sources)
	}
	if _, ok := hops[Node{Job: "ONLINE", Step: "EXTRACT", Program: "HOLDEXT", Field: "TD-AMOUNT"}]; !ok {
		t.Error("the card file's amounts do not go back to the online day")
	}
	if _, ok := hops[Node{Job: "DAILY", Step: "CHARGE", Program: "MONCHG", Field: "TD-AMOUNT"}]; !ok {
		t.Error("the month-end charges do not go back to MONCHG")
	}
}

// A field read from a dataset no step writes, by a program no job runs,
// ends where it is read.
func TestBackWithoutJob(t *testing.T) {
	s := pangaea(t)
	s.Jobs = nil
	trail := Back(s, Node{Program: "TXVALID", Field: "TD-AMOUNT"})
	if len(trail.Hops) != 1 || len(trail.Hops[0].Sources) != 1 || trail.Hops[0].Sources[0].Writers != nil {
		t.Errorf("trail %+v, want the field alone", trail)
	}
}

// An amount of the branches' transactions goes, through the sort and the
// postings, into the balances of the master, the audit trail, and from
// there into the month's statements and the ledger.
func TestForwardFromInput(t *testing.T) {
	trail := Forward(pangaea(t), Node{Program: "TXVALID", Field: "TD-AMOUNT"})
	hops := map[Node]Hop{}
	for _, h := range trail.Hops {
		hops[h.Node] = h
	}
	post, ok := hops[Node{Job: "DAILY", Step: "POST", Program: "TXPOST", Field: "TD-AMOUNT"}]
	if !ok || !slices.Contains(post.Effect.Items, "AM-BALANCE") || !slices.Contains(post.Effect.Items, "AU-AMOUNT") {
		t.Fatalf("the posting does not move the amount into the balance and the audit trail: %+v", post.Effect.Items)
	}
	valid := hops[Node{Job: "DAILY", Step: "VALID", Program: "TXVALID", Field: "TD-AMOUNT"}]
	sorted := false
	for _, tgt := range valid.Targets {
		sorted = sorted || fmt.Sprint(tgt.Datasets, tgt.Sorts, tgt.Readers) ==
			"[&&VALID &&SORTED] [{DAILY SORT TXSORT  0 0}] [{DAILY POST TXPOST TD-AMOUNT 0 0}]"
	}
	if !sorted {
		t.Errorf("the valid transactions do not go through the sort: %+v", valid.Targets)
	}
	statement := hops[Node{Job: "MONTHLY", Step: "STMT", Program: "STMTPRT", Field: "AU-AMOUNT"}]
	if !slices.ContainsFunc(statement.Targets, func(tgt Target) bool { return tgt.Report == "STMTOUT" }) {
		t.Errorf("the audited amounts are not printed on the statements: %+v", statement.Targets)
	}
	ledger := hops[Node{Job: "MONTHLY", Step: "LEDGER", Program: "GLSUM", Field: "AU-AMOUNT"}]
	if !slices.ContainsFunc(ledger.Targets, func(tgt Target) bool { return fmt.Sprint(tgt.Outputs) == "[GL.SUMMARY(+1)]" }) {
		t.Errorf("the audited amounts do not reach the ledger: %+v", ledger.Targets)
	}
	maint := hops[Node{Job: "DAILY", Step: "MAINT", Program: "ACCTMNT", Field: "ACCT-IN-RECORD", Start: 57, Length: 13}]
	if !slices.Contains(maint.Effect.Items, "AM-BALANCE") {
		t.Errorf("the balance read in a FILLER is not followed by its bytes: %+v", maint.Effect)
	}
}

// AUTHSRV DISPLAYs its response, the end of the online day's trail: the
// result is displayed with the record it is in.
func TestForwardToDisplay(t *testing.T) {
	trail := Forward(pangaea(t), Node{Program: "AUTHSRV", Field: "RS-RESULT"})
	if len(trail.Hops) == 0 || len(trail.Hops[0].Displays) == 0 || trail.Hops[0].Displays[0].Items[0] != "AUTH-RESPONSE" {
		t.Errorf("trail %+v, want the response displayed", trail.Hops)
	}
}

// A trail stops at its cap, and says so.
func TestForwardCapped(t *testing.T) {
	tr := newTracer(pangaea(t))
	tr.maxHops = 5
	trail := tr.forwardFrom(Node{Program: "BANKDATE", Field: "BC-BUSINESS-DATE"})
	if len(trail.Hops) != 5 || !trail.Capped {
		t.Errorf("%d hops, capped %v", len(trail.Hops), trail.Capped)
	}
}
