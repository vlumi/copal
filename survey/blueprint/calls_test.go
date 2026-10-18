package blueprint

import (
	"reflect"
	"strings"
	"testing"
)

func TestCalls(t *testing.T) {
	s, err := Analyze([]string{"testdata/CALLER.cbl", "testdata/CALLEE.cbl"}, "testdata/no-copybooks")
	if err != nil {
		t.Fatal(err)
	}
	callee, caller := s.Programs[0], s.Programs[1]
	wantCalls := []Call{
		{Program: "CALLEE", Using: []string{"WS-REQUEST", "WS-DATE"}, Paragraph: "MAIN-PARA", Source: "testdata/CALLER.cbl:9"},
		{Program: "MISSING", External: true, Paragraph: "MAIN-PARA", Source: "testdata/CALLER.cbl:10"},
	}
	if !reflect.DeepEqual(caller.Calls, wantCalls) {
		t.Errorf("calls %+v, want %+v", caller.Calls, wantCalls)
	}
	if want := []string{"CALLER"}; !reflect.DeepEqual(callee.CalledBy, want) {
		t.Errorf("CALLEE called by %q, want %q", callee.CalledBy, want)
	}
	if callee.Kind != KindCalled || caller.Kind != KindMain {
		t.Errorf("CALLEE is %s, CALLER %s", callee.Kind, caller.Kind)
	}
	if callee.Records[0].Section != "linkage" || !reflect.DeepEqual(callee.Using, []string{"LK-REQUEST", "LK-DATE"}) {
		t.Errorf("CALLEE records %+v, using %q", callee.Records, callee.Using)
	}
	md := programMarkdown("blueprints", callee)
	if want := "Called with LK-REQUEST, LK-DATE by [CALLER](CALLER.md)."; !strings.Contains(md, want) {
		t.Errorf("CALLEE blueprint has no %q:\n%s", want, md)
	}
	if want := "| MISSING (not analyzed) | | MAIN-PARA |"; !strings.Contains(programMarkdown("blueprints", caller), want) {
		t.Errorf("CALLER blueprint has no %q", want)
	}
}

// A program's DISPLAYs to the job's output are listed, those upon another
// device not.
func TestDisplays(t *testing.T) {
	s, err := Analyze([]string{"testdata/REPLY.cbl"}, "testdata/no-copybooks")
	if err != nil {
		t.Fatal(err)
	}
	want := []Display{{Items: []string{"'REPLY: '", "WS-REPLY"}, Paragraph: "MAIN-PARA", Source: "testdata/REPLY.cbl:22"}}
	if p := s.Programs[0]; !p.ReadsSysin || !reflect.DeepEqual(p.Displays, want) {
		t.Errorf("reads SYSIN %v, displays %+v, want %+v", p.ReadsSysin, p.Displays, want)
	}
}

// A program's inputs are the records it reads from files, a record read
// INTO another described by that one, with where each field's values go.
func TestInputs(t *testing.T) {
	s, err := Analyze([]string{"testdata/REPLY.cbl"}, "testdata/no-copybooks")
	if err != nil {
		t.Fatal(err)
	}
	want := []Input{{File: "IN-FILE", Record: "WS-IN", Fields: []FieldEffect{
		{Name: "WS-IN-KEY", Start: 1, Length: 4, Reaches: []string{"WS-REPLY"}},
		{Name: "WS-IN-REST", Start: 5, Length: 6},
	}}}
	if got := s.Programs[0].Inputs; !reflect.DeepEqual(got, want) {
		t.Errorf("inputs %+v, want %+v", got, want)
	}
}

// A call writes the caller's fields over the bytes the program called
// writes of what it is passed; a program outside the system may write all
// it is passed by reference.
func TestCallWrites(t *testing.T) {
	s, err := Analyze([]string{"testdata/ASKER.cbl", "testdata/ANSWERER.cbl"}, "testdata/no-copybooks")
	if err != nil {
		t.Fatal(err)
	}
	asker := s.Programs[1]
	writes := callWrites(asker, s.Programs)
	calls := asker.source.Paragraphs[0].Statements
	if items, known := writes(calls[0]); !known || !reflect.DeepEqual(items, []string{"WS-ANSWER"}) {
		t.Errorf("ANSWERER writes %q, %v, want WS-ANSWER", items, known)
	}
	if _, known := writes(calls[1]); known {
		t.Error("what ELSEWHERE writes is known")
	}
}
