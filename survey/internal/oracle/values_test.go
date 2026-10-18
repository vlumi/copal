package oracle

import (
	"testing"

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
