package blueprint

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
)

func laurasiaSystem(t *testing.T) *System {
	t.Helper()
	programs, err := filepath.Glob("../../laurasia/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	s, err := Analyze(programs, "../../laurasia/copybooks", cobol.MicroFocus)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestParagraphFacts(t *testing.T) {
	s := laurasiaSystem(t)
	f, err := s.ParagraphFacts("DISPUTE", "APPLY-EVENT")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Program DISPUTE, paragraph APPLY-EVENT.",
		"apply-event section.",
		"when holding-case also dc-charged-back and de-representment",
		"Performed by: APPLY-EVENTS.",
		"Writes: ",
		"HOLDING-CASE AND DC-CHARGED-BACK AND DE-REPRESENTMENT",
	} {
		if !strings.Contains(f.Text, want) {
			t.Errorf("facts have no %q:\n%s", want, f.Text)
		}
	}
	if strings.Contains(f.Text, "write-adjustment section.") {
		t.Error("the facts run on into the next section's source")
	}
	again, _ := s.ParagraphFacts("DISPUTE", "APPLY-EVENT")
	if len(f.Hash) != 64 || again.Hash != f.Hash {
		t.Errorf("hash %q, then %q, want the same of 64 hex digits", f.Hash, again.Hash)
	}
	if _, err := s.ParagraphFacts("DISPUTE", "NO-SUCH-PARA"); err == nil {
		t.Error("facts of a paragraph the program does not have")
	}
}

func TestCheck(t *testing.T) {
	f := Facts{Kind: KindParagraph, Program: "P", Subject: "S", Text: "move dc-amount to ad-amount\nWrites: AD-AMOUNT."}
	for text, ok := range map[string]bool{
		"Moves the case's amount into `ad-amount`.":          true,
		"Copies DC-AMOUNT for the over-limit and month-end.": true,
		"Sets `AD-AMOUNT-TOTAL` from DC-AMOUNT.":             false,
		"Clears WS-NOTHING.":                                 false,
		"Writes `AMOUNT`, part of a name only.":              false,
	} {
		if err := f.Check(Annotation{Text: text}); (err == nil) != ok {
			t.Errorf("check of %q: %v, want ok %v", text, err, ok)
		}
	}
}

func TestAnnotationsOnPage(t *testing.T) {
	s := laurasiaSystem(t)
	f, err := s.ParagraphFacts("MINPAY", "MAIN")
	if err != nil {
		t.Fatal(err)
	}
	s.Annotations = []Annotation{
		{Kind: KindParagraph, Program: "MINPAY", Subject: "MAIN", Hash: f.Hash, Name: "Work out the minimum payment",
			Text: "The payment is the floor at least.", Writer: "test", Prompt: "1"},
	}
	i := 0
	for s.Programs[i].Name != "MINPAY" {
		i++
	}
	page := s.programMarkdown("blueprints", s.Programs[i])
	if want := "| MAIN | Work out the minimum payment | The payment is the floor at least. | test, prompt 1 |"; !strings.Contains(page, want) {
		t.Errorf("page has no %q:\n%s", want, page)
	}
	s.Annotations[0].Hash = "written from other facts"
	if page := s.programMarkdown("blueprints", s.Programs[i]); !strings.Contains(page, "| MAIN | | stale: ") {
		t.Errorf("a stale annotation is shown as if it were not:\n%s", page)
	}
}

func TestAnnotationFiles(t *testing.T) {
	dir := t.TempDir()
	written := []Annotation{
		{Kind: KindParagraph, Program: "A", Subject: "P1", Hash: "h1", Text: "one", Writer: "w", Prompt: "1"},
		{Kind: KindParagraph, Program: "B", Subject: "P2", Hash: "h2", Text: "two", Writer: "w", Prompt: "1"},
	}
	if err := WriteAnnotations(dir, written); err != nil {
		t.Fatal(err)
	}
	read, err := ReadAnnotations(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(read) != 2 || read[0] != written[0] || read[1] != written[1] {
		t.Errorf("read back %+v", read)
	}
	if err := WriteAnnotations(dir, written[1:]); err != nil {
		t.Fatal(err)
	}
	if read, _ := ReadAnnotations(dir); len(read) != 1 {
		t.Errorf("%d annotations after writing one, want the old files gone", len(read))
	}
	if read, err := ReadAnnotations(filepath.Join(dir, "none")); err != nil || len(read) != 0 {
		t.Errorf("a directory that is not there: %v, %v", read, err)
	}
}
