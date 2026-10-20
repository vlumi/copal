package lineage

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/survey/blueprint"
)

func laurasia(t *testing.T) *blueprint.System {
	t.Helper()
	programs, err := filepath.Glob("../../../laurasia/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	scripts, err := filepath.Glob("../../../laurasia/scripts/*.sh")
	if err != nil {
		t.Fatal(err)
	}
	s, err := blueprint.Analyze(append(programs, scripts...), "../../../laurasia/copybooks", cobol.MicroFocus)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// The authorization code MATCH reads comes, through the script's sort into
// authorization code order, from CLRIN, which took it from the clearing
// file's presentment; and the code CLRIN writes goes on, through the sort,
// to MATCH.
func TestThroughScriptSort(t *testing.T) {
	s := laurasia(t)
	hops := map[Node]Hop{}
	for _, h := range Back(s, Node{Program: "MATCH", Field: "PS-AUTH-CODE"}).Hops {
		hops[h.Node] = h
	}
	match, ok := hops[Node{Job: "nightly", Step: "MATCH", Program: "MATCH", Field: "PS-AUTH-CODE"}]
	if !ok || len(match.Sources) == 0 {
		t.Fatalf("no sources of MATCH's PS-AUTH-CODE: %+v", match)
	}
	src := match.Sources[0]
	if got := fmt.Sprint(src.Sorts, src.Writers); got != "[{nightly PRESSORT sort  0 0}] [{nightly CLRIN CLRIN PS-AUTH-CODE 0 0}]" {
		t.Errorf("MATCH's authorization code comes from %s", got)
	}
	clrin := hops[Node{Job: "nightly", Step: "CLRIN", Program: "CLRIN", Field: "PS-AUTH-CODE"}]
	if got := fmt.Sprint(clrin.Trace.Origins); got != "[file CLEARING-IN PR-AUTH-CODE]" || len(clrin.Sources) != 1 || fmt.Sprint(clrin.Sources[0].External) != "[$DATA/inbound/CLR.$DATE.txt]" {
		t.Errorf("CLRIN's authorization code comes from %s, %+v; want the clearing file's presentment, from outside", got, clrin.Sources)
	}
	forward := map[Node]bool{}
	for _, h := range Forward(s, Node{Job: "nightly", Step: "CLRIN", Program: "CLRIN", Field: "PS-AUTH-CODE"}).Hops {
		forward[h.Node] = true
	}
	if !forward[Node{Job: "nightly", Step: "MATCH", Program: "MATCH", Field: "PS-AUTH-CODE"}] {
		t.Error("the authorization code CLRIN writes does not reach MATCH through the sort")
	}
}
