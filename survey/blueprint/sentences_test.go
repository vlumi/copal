package blueprint

import (
	"reflect"
	"testing"

	"github.com/vlumi/cobolith/cobol"
)

// TestSentences checks the rules of a table written as sentences: a
// conditional phrase raised and not, a condition name with its field and
// values, a compound found false turned into its parts, and the actions
// in words, a PERFORM with its UNTIL.
func TestSentences(t *testing.T) {
	s, err := Analyze([]string{"testdata/SENTENCE.cbl"}, "testdata/no-copybooks", cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range s.Programs[0].Rules.Tables[0].Rules {
		got = append(got, r.Sentence)
	}
	want := []string{
		"When reading `IN-FILE` finds no more records, it opens `IN-FILE` for input, reads `IN-FILE`, displays `'EMPTY'`, and stops the run.",
		"When reading `IN-FILE` finds a record; and `IN-KIND` is `'O'` or `'R'` (`IN-OPEN`) and either `IN-AMOUNT` is more than `WS-LIMIT` or `WS-TOTAL` is 0, " +
			"it opens `IN-FILE` for input, reads `IN-FILE`, adds `IN-AMOUNT` to `WS-TOTAL`, and stops the run.",
		"When reading `IN-FILE` finds a record; and `IN-KIND` is neither `'O'` nor `'R'` (not `IN-OPEN`) or both `IN-AMOUNT` is at most `WS-LIMIT` and `WS-TOTAL` is not 0, " +
			"it opens `IN-FILE` for input; reads `IN-FILE`; sets `WS-TOTAL` and `WS-LIMIT` to 0; performs `CHECK-PARA` until `WS-TOTAL` is at least 10; and stops the run.",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("sentences:\n%q\nwant\n%q", got, want)
	}
}
