package annotate

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/survey/blueprint"
)

// factsWriter writes from the facts alone, with no model: the paragraph's
// name, and what it performs and writes, in the code's names. For one
// paragraph it names an item the facts do not, as a model might.
type factsWriter struct{ calls int }

func (w *factsWriter) Identity() (string, string) { return "facts", "1" }

func (w *factsWriter) Write(_ context.Context, f blueprint.Facts) (string, string, error) {
	w.calls++
	if f.Subject == "MATCH-NEXT" {
		return "Match the next", "Moves `WS-NO-SUCH-ITEM` on.", nil
	}
	var writes string
	for _, line := range strings.Split(f.Text, "\n") {
		if rest, ok := strings.CutPrefix(line, "Writes: "); ok {
			writes = strings.TrimSuffix(rest, ".")
		}
	}
	return "Do " + strings.ToLower(f.Subject), "Writes `" + strings.ReplaceAll(writes, ", ", "`, `") + "`.", nil
}

func laurasia(t *testing.T) *blueprint.System {
	t.Helper()
	programs, err := filepath.Glob("../../../laurasia/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	s, err := blueprint.Analyze(programs, "../../../laurasia/copybooks", cobol.MicroFocus)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestParagraphs(t *testing.T) {
	s := laurasia(t)
	w := &factsWriter{}
	kept, rejected, err := Paragraphs(context.Background(), s, nil, w)
	if err != nil {
		t.Fatal(err)
	}
	paragraphs := 0
	for _, p := range s.Programs {
		paragraphs += len(p.Paragraphs)
	}
	if w.calls != paragraphs || len(kept) != paragraphs-1 {
		t.Fatalf("%d calls and %d kept for %d paragraphs, want a call each and all but one kept", w.calls, len(kept), paragraphs)
	}
	if len(rejected) != 1 || rejected[0].Annotation.Subject != "MATCH-NEXT" || !strings.Contains(rejected[0].Reason, "WS-NO-SUCH-ITEM") {
		t.Errorf("rejected %+v, want MATCH-NEXT for naming WS-NO-SUCH-ITEM", rejected)
	}

	again, _, err := Paragraphs(context.Background(), s, kept, w)
	if err != nil {
		t.Fatal(err)
	}
	if w.calls != paragraphs+1 || len(again) != len(kept) {
		t.Errorf("%d more calls annotating again, want 1, for the one rejected", w.calls-paragraphs)
	}

	stale := kept[0]
	stale.Hash = "changed"
	if _, _, err := Paragraphs(context.Background(), s, []blueprint.Annotation{stale}, w); err != nil {
		t.Fatal(err)
	}
	if w.calls != 2*paragraphs+1 {
		t.Errorf("an annotation whose facts changed was kept")
	}
}
