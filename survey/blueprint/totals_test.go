package blueprint

import (
	"maps"
	"testing"

	"github.com/vlumi/cobolith/cobol"
)

// TestTotalsWhere checks that a control total adds up only the records
// where the conditions on their own fields hold that the ADD runs under:
// the count of details, and of them, those paid. A condition that is not
// on the record's fields is not among them.
func TestTotalsWhere(t *testing.T) {
	s, err := Analyze([]string{"testdata/TOTALS.cbl"}, "testdata/no-copybooks", cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	when := map[string]string{}
	for _, total := range s.Programs[0].Totals {
		when[total.Field] = total.When
	}
	if want := map[string]string{"RT-COUNT": "RS-DETAIL", "RT-PAID-COUNT": "RS-DETAIL AND RS-WAS-PAID"}; !maps.Equal(when, want) {
		t.Errorf("totals where %q, want %q", when, want)
	}
}
