package structure

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
)

func TestParagraphs(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/FLOW.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	paragraphs, err := Paragraphs(p)
	if err != nil {
		t.Fatal(err)
	}
	type summary struct {
		name        string
		performs    []string
		performedBy []string
	}
	var got []summary
	for _, para := range paragraphs {
		s := summary{name: para.Name, performedBy: para.PerformedBy}
		for _, perform := range para.Performs {
			target := perform.Target
			if perform.Thru != "" {
				target += " THRU " + perform.Thru
			}
			if perform.Sort {
				target = "SORT " + target
			}
			s.performs = append(s.performs, target)
		}
		got = append(got, s)
	}
	want := []summary{
		{name: "MAIN-PARA", performs: []string{"INIT-PARA", "WORK-PARA THRU WORK-EXIT", "INIT-PARA", "SORT INIT-PARA", "SORT WORK-PARA THRU WORK-EXIT"}},
		{name: "INIT-PARA", performedBy: []string{"MAIN-PARA"}},
		{name: "WORK-PARA", performedBy: []string{"MAIN-PARA"}},
		{name: "WORK-EXIT", performedBy: []string{"MAIN-PARA"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
	if pos := paragraphs[0].Performs[1].Pos; pos.Line != 16 {
		t.Errorf("PERFORM THRU at line %d, want 16", pos.Line)
	}
}

func TestGoTosAndFallThrough(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/GOTO.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	paragraphs, err := Paragraphs(p)
	if err != nil {
		t.Fatal(err)
	}
	type summary struct {
		name        string
		goTos       []string
		fallsInto   string
		performedBy []string
		reachedBy   []string
	}
	var got []summary
	for _, para := range paragraphs {
		s := summary{name: para.Name, fallsInto: para.FallsInto, performedBy: para.PerformedBy, reachedBy: para.ReachedBy}
		for _, goTo := range para.GoTos {
			s.goTos = append(s.goTos, goTo.Target)
		}
		got = append(got, s)
	}
	main := []string{"MAIN-PARA"}
	want := []summary{
		{name: "MAIN-PARA"},
		{name: "CHECK-PARA", goTos: []string{"CHECK-EXIT", "CHECK-ONE", "CHECK-TWO"}, fallsInto: "CHECK-ONE", performedBy: main},
		{name: "CHECK-ONE", goTos: []string{"CHECK-EXIT"}, performedBy: main, reachedBy: []string{"CHECK-PARA", "CHECK-TWO"}},
		{name: "CHECK-TWO", goTos: []string{"CHECK-EXIT", "CHECK-ONE"}, performedBy: main, reachedBy: []string{"CHECK-PARA"}},
		{name: "CHECK-EXIT", performedBy: main, reachedBy: []string{"CHECK-PARA", "CHECK-ONE", "CHECK-TWO"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

func TestSections(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/SECTIONS.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	paragraphs, err := Paragraphs(p)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, para := range paragraphs {
		got = append(got, fmt.Sprintf("%s/%s %v", para.Section, para.Name, para.PerformedBy))
	}
	want := []string{"MAIN/MAIN []", "CHECK/CHECK [MAIN]", "CHECK/CHECK-MORE [MAIN]", "TAIL-ONLY/TAIL-STEP [MAIN]"}
	if !slices.Equal(got, want) {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}
