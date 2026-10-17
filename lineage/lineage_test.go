package lineage

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/vlumi/copal/blueprint"
)

func pangaea(t *testing.T) *blueprint.System {
	t.Helper()
	programs, err := filepath.Glob("../pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := filepath.Glob("../pangaea/jobs/*.jcl")
	if err != nil {
		t.Fatal(err)
	}
	s, err := blueprint.Analyze(append(programs, jobs...), "../pangaea/copybooks")
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
		"[&&SORTED &&VALID &&CARDVAL] [{DAILY VALID TXVALID TD-AMOUNT} {DAILY VALIDC TXVALID TD-AMOUNT}]" {
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
