package oracle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/internal/golden"
	"github.com/vlumi/copal/survey/blueprint"
)

// TestLaurasia replays the cases recorded in oracle/laurasia and compares
// what each reaches and writes, and the coverage, with what was recorded,
// as TestPangaea does for Pangaea. Run with -update to record what the
// cases now do; `make oracle` searches for them again.
func TestLaurasia(t *testing.T) {
	t.Chdir("../../..")
	programs, err := filepath.Glob("laurasia/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	scripts, err := filepath.Glob("laurasia/scripts/*.sh")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(append(programs, scripts...), "laurasia/copybooks", cobol.MicroFocus)
	if err != nil {
		t.Fatal(err)
	}
	scratch := t.TempDir()
	var results []*Result
	for i := range system.Programs {
		p := &system.Programs[i]
		if p.Kind != blueprint.KindMain {
			continue
		}
		dir := filepath.Join("oracle/laurasia", p.Name)
		result, err := Replay(system, p, "laurasia/copybooks", scratch, dir)
		if err != nil {
			t.Fatalf("%s: %v", p.Name, err)
		}
		golden.CompareFiles(t, result.Files("../../../blueprints/laurasia/rules"), dir)
		results = append(results, result)
	}
	want, err := os.ReadFile("oracle/laurasia/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if got := Index("laurasia", results, nil); string(got) != string(want) {
		t.Errorf("oracle/laurasia/README.md differs:\n%s", got)
	}
}
