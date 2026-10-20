package blueprint

import (
	"testing"

	"github.com/vlumi/cobolith/layout"
)

// Each rule of fieldKind, and the order they are tried in.
func TestFieldKind(t *testing.T) {
	numeric := func(name string, digits, scale int) *layout.Field {
		return &layout.Field{Name: name, Category: layout.Numeric, Digits: digits, Scale: scale}
	}
	coded := numeric("XX-STATUS", 1, 0)
	coded.Conditions = []layout.Condition{{Name: "XX-OPEN", Values: []string{"1"}}}
	keyedGroup := &layout.Field{Name: "XX-KEY", Category: layout.Group, Children: []*layout.Field{numeric("XX-NO", 10, 0)}}
	for _, c := range []struct {
		field *layout.Field
		want  string
	}{
		{keyedGroup, KindGroup},
		{&layout.Field{Name: "FILLER", Category: layout.Alphanumeric}, KindFiller},
		{numeric("XX-ACCT-NO", 10, 0), KindKey},
		{coded, KindCode},
		{&layout.Field{Name: "XX-PRINTED", Category: layout.NumericEdited}, KindEdited},
		{numeric("XX-BALANCE", 13, 2), KindAmount},
		{numeric("XX-RATE-DATE", 8, 2), KindAmount},
		{numeric("XX-OPEN-DATE", 8, 0), KindDate},
		{numeric("XX-YYMM", 6, 0), KindDate},
		{numeric("XX-STAMP-TIMESTAMP", 14, 0), KindDate},
		{numeric("XX-MONTH", 2, 0), KindNumber},
		{numeric("XX-COUNT", 8, 0), KindNumber},
		{&layout.Field{Name: "XX-NAME", Category: layout.Alphanumeric}, KindText},
		{&layout.Field{Name: "XX-DATE", Category: layout.Alphanumeric}, KindText},
	} {
		keys := map[string]bool{"XX-ACCT-NO": true, "XX-KEY": true}
		if got := fieldKind(c.field, keys); got != c.want {
			t.Errorf("%s: %s, want %s", c.field.Name, got, c.want)
		}
	}
}
