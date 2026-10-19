package oracle

import (
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
