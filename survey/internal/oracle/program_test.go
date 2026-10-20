package oracle

import (
	"fmt"
	"slices"
	"testing"

	"github.com/vlumi/cobolith/layout"
)

// A record that leaves bytes FILLER, moved whole into a record as long,
// takes that record's fields over those bytes, at its own positions.
func TestOverFiller(t *testing.T) {
	r := &layout.Field{Name: "IN-RECORD", Start: 1, Length: 18, Children: []*layout.Field{
		{Name: "IN-KEY", Start: 1, Length: 4},
		{Name: "FILLER", Start: 5, Length: 14},
	}}
	other := &layout.Field{Name: "WS-RECORD", Start: 101, Length: 18, Children: []*layout.Field{
		{Name: "WS-KEY", Start: 101, Length: 4},
		{Name: "WS-DATE", Start: 105, Length: 8},
		{Name: "FILLER", Start: 113, Length: 2},
		{Name: "WS-CODE", Start: 115, Length: 2},
	}}
	filled := overFiller(r, other)
	if filled == nil {
		t.Fatal("no fields over the FILLER")
	}
	var got []string
	for _, f := range filled.Children {
		got = append(got, fmt.Sprintf("%s %d-%d", f.Name, f.Start, f.Start+f.Size()-1))
	}
	want := []string{"IN-KEY 1-4", "WS-DATE 5-12", "FILLER 13-14", "WS-CODE 15-16", "FILLER 17-18"}
	if !slices.Equal(got, want) {
		t.Errorf("fields %q, want %q", got, want)
	}
	if overFiller(other, r) != nil {
		t.Error("fields over a FILLER of none")
	}
}
