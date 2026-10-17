package fileio

import (
	"fmt"
	"slices"
	"testing"

	"github.com/vlumi/copal/cobol"
)

func TestUses(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/FILES.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, u := range Uses(p) {
		got = append(got, fmt.Sprintf("%s %s", u.Mode, u.File))
	}
	want := []string{"INPUT IN-FILE", "OUTPUT OUT-FILE", "EXTEND LOG-FILE", "OUTPUT SORTED"}
	if !slices.Equal(got, want) {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestSorts(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/FILES.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	sorts := Sorts(p)
	if len(sorts) != 1 || sorts[0].Work != "WORK-SORT" || !slices.Equal(sorts[0].Using, []string{"IN-FILE"}) ||
		!slices.Equal(sorts[0].Giving, []string{"SORTED"}) {
		t.Errorf("sorts %+v", sorts)
	}
}
