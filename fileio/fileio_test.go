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
