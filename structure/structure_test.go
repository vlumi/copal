package structure

import (
	"reflect"
	"testing"

	"github.com/vlumi/copal/cobol"
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
			s.performs = append(s.performs, target)
		}
		got = append(got, s)
	}
	want := []summary{
		{name: "MAIN-PARA", performs: []string{"INIT-PARA", "WORK-PARA THRU WORK-EXIT", "INIT-PARA"}},
		{name: "INIT-PARA", performedBy: []string{"MAIN-PARA"}},
		{name: "WORK-PARA", performedBy: []string{"MAIN-PARA"}},
		{name: "WORK-EXIT"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
	if pos := paragraphs[0].Performs[1].Pos; pos.Line != 9 {
		t.Errorf("PERFORM THRU at line %d, want 9", pos.Line)
	}
}
