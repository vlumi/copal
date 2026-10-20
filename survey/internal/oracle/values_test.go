package oracle

import (
	"slices"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
)

// A literal becomes a value for its field: a number in units of the
// field's scale, text without its quotes, or a figurative constant.
func TestLiteral(t *testing.T) {
	amount := &layout.Field{Name: "AMOUNT", Category: layout.Numeric, Digits: 7, Scale: 2, Length: 7}
	code := &layout.Field{Name: "CODE", Category: "alphanumeric", Length: 2}
	for _, c := range []struct {
		field *layout.Field
		lit   string
		want  value
	}{
		{amount, "500", number(50000)},
		{amount, "12.5", number(1250)},
		{amount, "ZERO", number(0)},
		{code, "'DP'", text("DP")},
		{code, "SPACES", text("")},
		{code, "HIGH-VALUES", filled(0xff)},
		{code, "ZEROS", text("00")},
	} {
		if got, ok := literal(c.field, c.lit); !ok || got != c.want {
			t.Errorf("literal %s for %s: %v, %v, want %v", c.lit, c.field.Name, got, ok, c.want)
		}
	}
}

// What a FUNCTION MOD or REM divides an item by is a value worth trying
// for it: a literal as written, or the VALUE of the field named.
func TestDivisor(t *testing.T) {
	records := &layout.ProgramRecords{WorkingStorage: []*layout.Field{{Name: "WS-NOTE", Values: []string{"20"}}}}
	for _, c := range []struct {
		text string
		want string
		ok   bool
	}{
		{"FUNCTION MOD(RQ-AMOUNT WS-NOTE)", "20", true},
		{"FUNCTION REM(RQ-AMOUNT 50)", "50", true},
		{"FUNCTION MOD(OTHER WS-NOTE)", "", false},
		{"RQ-AMOUNT", "", false},
	} {
		got, ok := divisor(records, &cobol.Operand{Text: c.text}, "RQ-AMOUNT")
		if got != c.want || ok != c.ok {
			t.Errorf("divisor of %s: %q, %v, want %q, %v", c.text, got, ok, c.want, c.ok)
		}
	}
}

// A condition name on a table's entry gives its field its values, and the
// condition names of another item give it none of theirs.
func TestConditionNameValues(t *testing.T) {
	kind := &layout.Field{Name: "WS-RT-KIND", Category: "alphanumeric", Length: 1, Conditions: []layout.Condition{
		{Name: "WS-RT-CREDIT", Values: []string{"'C'"}}, {Name: "WS-RT-DEBIT", Values: []string{"'D'"}}}}
	found := &layout.Field{Name: "WS-FOUND", Category: "alphanumeric", Length: 1, Conditions: []layout.Condition{
		{Name: "RATE-FOUND", Values: []string{"'Y'"}}}}
	records := &layout.ProgramRecords{WorkingStorage: []*layout.Field{{Name: "WS", Children: []*layout.Field{kind, found}}}}
	input := &layout.Field{Name: "RT-KIND", Category: "alphanumeric", Length: 1}
	debit := &cobol.Condition{Op: "NAME", Left: &cobol.Operand{Text: "WS-RT-DEBIT(WS-RX)", Name: "WS-RT-DEBIT"}}
	both := &cobol.Condition{Op: "AND", Parts: []*cobol.Condition{debit,
		{Op: "NAME", Left: &cobol.Operand{Text: "RATE-FOUND", Name: "RATE-FOUND"}}}}
	got := valuesFor(records, nil, both, "WS-RT-DEBIT", input)
	if want := []value{text("C"), text("D"), text("?"), text("")}; !slices.Equal(got, want) {
		t.Errorf("values %v, want %v", got, want)
	}
}

// A part of a field of digits is a number of the part's digits, set alone.
func TestPart(t *testing.T) {
	root := &layout.Field{Name: "R", Start: 1, Length: 10}
	date := &layout.Field{Name: "DATE", Start: 3, Length: 8, Digits: 8, Category: layout.Numeric}
	root.Children = []*layout.Field{{Name: "CODE", Start: 1, Length: 2, Category: "alphanumeric"}, date}
	if got := digitsOver(root, 7, 2); got != date {
		t.Fatalf("the digits over 7-8 are of %v, want DATE", got)
	}
	month := part{offset: 4, length: 2}
	r := []byte("AB20261016")
	if !month.set(r, root, date, number(12)) || string(r) != "AB20261216" {
		t.Errorf("month 12 set: %q", r)
	}
	if month.set(r, root, date, number(100)) {
		t.Errorf("month 100 set: %q", r)
	}
}

// A comparison of an item with a number, within a condition, is a count
// worth repeating records to; one with another item is not.
func TestComparisons(t *testing.T) {
	count := &cobol.Condition{Op: "=", Left: &cobol.Operand{Text: "WS-COUNT", Name: "WS-COUNT", Refs: []string{"WS-COUNT"}}, Right: &cobol.Operand{Text: "20"}}
	other := &cobol.Condition{Op: ">", Left: &cobol.Operand{Text: "WS-A", Name: "WS-A", Refs: []string{"WS-A"}}, Right: &cobol.Operand{Text: "WS-B", Name: "WS-B", Refs: []string{"WS-B"}}}
	got := comparisons(&cobol.Condition{Op: "OR", Parts: []*cobol.Condition{count, other}})
	if len(got) != 1 || got[0] != count {
		t.Errorf("comparisons %v, want the count's alone", got)
	}
}
