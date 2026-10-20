package lineage

import (
	"fmt"
	"path/filepath"
	"slices"
	"testing"

	"github.com/vlumi/cobolith/cobol"
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
	s, err := blueprint.Analyze(append(programs, jobs...), "../../../pangaea/copybooks", cobol.IBM)
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

// A field over the bytes of one record type of a file finds the writer's
// field over the same bytes of the same record type.
func TestBackFromHeader(t *testing.T) {
	hops := hopsOf(Back(pangaea(t), Node{Job: "DAILY", Step: "VALIDC", Program: "TXVALID", Field: "TH-BUSINESS-DATE"}))
	if _, ok := hops[Node{Job: "ONLINE", Step: "EXTRACT", Program: "HOLDEXT", Field: "TH-BUSINESS-DATE"}]; !ok {
		t.Error("the card file's business date does not go back to the online day's header")
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

// hopsOf maps a trail's hops by node, without the bytes of a node.
func hopsOf(trail Trail) map[Node]Hop {
	hops := map[Node]Hop{}
	for _, h := range trail.Hops {
		hops[h.Node] = h
	}
	return hops
}

func online(program, field string) Node {
	return Node{Job: "ONLINE", Step: "AUTH", Program: program, Field: field}
}

// A hold's amount in the card file goes back into the online day: the
// hold record AUTHSRV writes, the programs it calls that set the amount,
// the request they are passed, and the job's SYSIN it is accepted from.
// A program called that does not change what it is passed gives nothing.
func TestBackAcrossCalls(t *testing.T) {
	hops := hopsOf(Back(pangaea(t), Node{Job: "ONLINE", Step: "EXTRACT", Program: "HOLDEXT", Field: "TD-AMOUNT"}))
	held, ok := hops[online("AUTHSRV", "HD-AMOUNT")]
	if !ok {
		t.Fatal("no hop at AUTHSRV's hold amount")
	}
	var called []string
	for _, src := range held.Sources {
		if src.Call != "" {
			called = append(called, src.Call)
		}
	}
	if fmt.Sprint(called) != "[ATMWDL AUTHPUR]" {
		t.Errorf("the hold amount comes from %v, want ATMWDL and AUTHPUR", called)
	}
	purchase, ok := hops[online("AUTHPUR", "HD-AMOUNT")]
	if !ok || !slices.ContainsFunc(purchase.Sources, func(s Source) bool {
		return s.Parameter == 1 && s.Field == "RQ-AMOUNT" && slices.Contains(s.Writers, online("AUTHSRV", "RQ-AMOUNT"))
	}) {
		t.Errorf("AUTHPUR's hold amount does not come from the request: %+v", purchase.Sources)
	}
	request := hops[online("AUTHSRV", "RQ-AMOUNT")]
	if len(request.Sources) != 1 || request.Sources[0].Stream != "SYSIN" {
		t.Errorf("the request's amount comes from %+v, want SYSIN", request.Sources)
	}
}

// The available balance AUTHSRV answers with goes back through the programs
// called to the online account, and to the steps that loaded it.
func TestBackFromResponse(t *testing.T) {
	hops := hopsOf(Back(pangaea(t), online("AUTHSRV", "RS-AVAILABLE")))
	if _, ok := hops[online("AUTHPUR", "RS-AVAILABLE")]; !ok {
		t.Fatal("no hop at AUTHPUR's available balance")
	}
	for _, field := range []string{"AO-BALANCE", "AO-OVERDRAFT-LIMIT", "AO-HOLDS"} {
		h, ok := hops[online("AUTHSRV", field)]
		if !ok || !slices.ContainsFunc(h.Sources, func(s Source) bool {
			return s.File == "ACCT-FILE" && slices.ContainsFunc(s.Writers, func(n Node) bool { return n.Program == "ACCTLOAD" })
		}) {
			t.Errorf("%s does not go back to ACCTLOAD through the online account file", field)
		}
	}
}

// The request's amount goes on into the programs called, back out through
// the hold record, into the hold file and the card file.
func TestForwardAcrossCalls(t *testing.T) {
	hops := hopsOf(Forward(pangaea(t), online("AUTHSRV", "RQ-AMOUNT")))
	for _, n := range []Node{online("AUTHPUR", "RQ-AMOUNT"), online("AUTHSRV", "HD-AMOUNT"),
		{Job: "ONLINE", Step: "EXTRACT", Program: "HOLDEXT", Field: "HD-AMOUNT"}} {
		if _, ok := hops[n]; !ok {
			t.Errorf("the request's amount does not reach %+v", n)
		}
	}
	if len(hops[online("AUTHSRV", "RS-AMOUNT")].Displays) == 0 {
		t.Error("the response that shows the amount is not displayed")
	}
}

// A writer's fields over the bytes are found in every record type it
// writes to the file, a header's as well as a detail's.
func TestOverEveryRecordType(t *testing.T) {
	p := &blueprint.Program{
		Records: []blueprint.Record{
			{Field: blueprint.Field{Name: "WS-HEADER", Start: 1, Length: 20}},
			{Field: blueprint.Field{Name: "WS-DETAIL", Start: 1, Length: 20}},
		},
		Outputs: []blueprint.Output{
			{File: "OUT", Record: "WS-HEADER", Fields: []blueprint.FieldOrigin{{Name: "WS-H-DATE", Start: 1, Length: 8}}},
			{File: "OUT", Record: "WS-DETAIL", Fields: []blueprint.FieldOrigin{{Name: "WS-D-ACCT", Start: 1, Length: 8}}},
			{File: "OTHER", Record: "WS-OTHER", Fields: []blueprint.FieldOrigin{{Name: "WS-O-FIELD", Start: 1, Length: 8}}},
		},
	}
	at := Node{Program: "WRITER"}
	got := overEach(at, outputs(p, "OUT"), 1, 8)
	want := []Node{{Program: "WRITER", Field: "WS-H-DATE"}, {Program: "WRITER", Field: "WS-D-ACCT"}}
	if !slices.Equal(got, want) {
		t.Errorf("fields over bytes 1-8: %v, want %v", got, want)
	}
}

// A step that updates a dataset in place, as AUTHSRV does the holds, is
// among the writers of what it reads.
func TestBackInPlace(t *testing.T) {
	hops := hopsOf(Back(pangaea(t), online("AUTHSRV", "HD-AMOUNT")))
	held := hops[online("AUTHSRV", "HD-AMOUNT")]
	if !slices.ContainsFunc(held.Sources, func(s Source) bool {
		return s.File == "HOLD-FILE" && slices.Contains(s.Writers, online("AUTHSRV", "HD-AMOUNT"))
	}) {
		t.Errorf("the holds AUTHSRV reads are not written by AUTHSRV: %+v", held.Sources)
	}
}
