package blueprint

import (
	"reflect"
	"testing"

	"github.com/vlumi/cobolith/cobol"
)

// A field written from an argument a program called writes has, behind
// the call, the origins of what the called program gives back in it: its
// own literals, marked as its, and what the CALL passes in, traced in the
// caller.
func TestOriginsThroughCalls(t *testing.T) {
	s, err := Analyze([]string{"testdata/WRITER.cbl", "testdata/CONVERT.cbl"}, "testdata/no-copybooks", cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	writer := s.Programs[1]
	if writer.Name != "WRITER" || len(writer.Outputs) != 1 {
		t.Fatalf("outputs of %s: %+v", writer.Name, writer.Outputs)
	}
	got := map[string][]string{}
	for _, f := range writer.Outputs[0].Fields {
		got[f.Name] = f.Origins
	}
	want := map[string][]string{
		"OUT-DATE": {"call CONVERT argument 3 WS-ERA-DATE", "in CONVERT, literal 'R'", "file IN-FILE IN-DATE"},
		"OUT-NAME": {"file IN-FILE IN-NAME"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("origins %q, want %q", got, want)
	}
}
