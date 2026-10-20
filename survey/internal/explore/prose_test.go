package explore

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/survey/blueprint"
)

// A paragraph's page shows what a model wrote about it, marked as a
// model's, and says when it is stale.
func TestProse(t *testing.T) {
	t.Chdir("../../..")
	paths, err := filepath.Glob("laurasia/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(paths, "laurasia/copybooks", cobol.MicroFocus)
	if err != nil {
		t.Fatal(err)
	}
	facts, err := system.ParagraphFacts("MINPAY", "MAIN")
	if err != nil {
		t.Fatal(err)
	}
	system.Annotations = []blueprint.Annotation{
		{Kind: blueprint.KindParagraph, Program: "MINPAY", Subject: "MAIN", Hash: facts.Hash,
			Name: "Work out the minimum payment", Text: "At least the floor.", Writer: "test", Prompt: "1"},
	}
	x, err := New(system, "Laurasia", "")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(x.Handler())
	t.Cleanup(server.Close)
	page := func() string {
		resp, err := http.Get(server.URL + "/programs/MINPAY/paragraphs/MAIN")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = resp.Body.Close() }()
		body, _ := io.ReadAll(resp.Body)
		return string(body)
	}
	if body := page(); !strings.Contains(body, "By a model, test, prompt 1") || !strings.Contains(body, "Work out the minimum payment.</strong> At least the floor.") {
		t.Errorf("the page does not show the prose, marked:\n%s", body)
	}
	system.Annotations[0].Hash = "written from other facts"
	if body := page(); !strings.Contains(body, "Stale: the facts have changed") || strings.Contains(body, "At least the floor.") {
		t.Error("a stale annotation is shown as if it were not")
	}
}
