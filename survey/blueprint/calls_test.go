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
