package oracle

import (
	"maps"
	"slices"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
)

// A table that redefines an area of literals holds, in each occurrence of
// each field, the bytes the literals put there, and an item compared with
// one of its fields takes those values.
func TestPresets(t *testing.T) {
	area := &layout.Field{Name: "ACCOUNTS", Start: 1, Length: 8, Children: []*layout.Field{
		{Name: "FILLER", Start: 1, Length: 4, Values: []string{`"1110"`}},
		{Name: "FILLER", Start: 5, Length: 4, Values: []string{`"12"`}},
	}}
	table := &layout.Field{Name: "ACCOUNT-TABLE", Start: 1, Length: 8, Redefines: "ACCOUNTS", Children: []*layout.Field{
		{Name: "ACCOUNT", Start: 1, Length: 4, Occurs: 2, Children: []*layout.Field{
			{Name: "ACCOUNT-CLASS", Start: 1, Length: 2},
			{Name: "ACCOUNT-NUMBER", Start: 3, Length: 2},
		}},
	}}
	records := &layout.ProgramRecords{WorkingStorage: []*layout.Field{area, table}}
	got := presets(records, "")
	want := map[string][]string{"ACCOUNT-CLASS": {"11", "12"}, "ACCOUNT-NUMBER": {"10", ""}}
	if !maps.EqualFunc(got, want, slices.Equal) {
		t.Errorf("presets %v, want %v", got, want)
	}
	input := &layout.Field{Name: "JR-CLASS", Category: layout.Numeric, Length: 2, Digits: 2}
	equal := &cobol.Condition{Op: "=", Left: &cobol.Operand{Text: "ACCOUNT-CLASS(IX)", Name: "ACCOUNT-CLASS", Refs: []string{"IX"}},
		Right: &cobol.Operand{Text: "WS-CLASS", Name: "WS-CLASS"}}
	values := valuesFor(records, nil, got, equal, "WS-CLASS", input)
	if want := []value{number(11), number(12), number(0), number(1), number(-1), number(99)}; !slices.Equal(values, want) {
		t.Errorf("values %v, want %v", values, want)
	}
}

// An area with a field of no literal, or one that repeats, is no table of
// values.
func TestPresetsOfNoLiterals(t *testing.T) {
	area := &layout.Field{Name: "TOTALS", Start: 1, Length: 4, Children: []*layout.Field{{Name: "FILLER", Start: 1, Length: 4}}}
	table := &layout.Field{Name: "TOTAL-TABLE", Start: 1, Length: 4, Redefines: "TOTALS", Children: []*layout.Field{
		{Name: "TOTAL", Start: 1, Length: 2, Occurs: 2}}}
	if got := presets(&layout.ProgramRecords{WorkingStorage: []*layout.Field{area, table}}, ""); len(got) != 0 {
		t.Errorf("presets %v, want none", got)
	}
}
