package blueprint

import (
	"strings"
	"testing"
)

// The rows of a layout draw the record's tree: condition names under their
// field, the last branch closed, the branches of open groups continued.
func TestLayoutRows(t *testing.T) {
	record := Field{Level: 1, Name: "REC", Fields: []Field{
		{Level: 5, Name: "A", Conditions: []Condition{{Name: "A-ON", Values: []string{"'Y'"}}}},
		{Level: 5, Name: "B", Fields: []Field{{Level: 10, Name: "B1"}, {Level: 10, Name: "B2"}}},
		{Level: 5, Name: "C"},
		{Level: 5, Name: "D", Redefines: "C"},
	}}
	var got []string
	for _, r := range LayoutRows(record) {
		got = append(got, strings.ReplaceAll(r.Prefix, " ", " ")+r.Name)
	}
	want := []string{"REC", "├ A", "│ └ A-ON", "├ B", "│ ├ B1", "│ └ B2", "├ C", "╘ D"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("rows\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
