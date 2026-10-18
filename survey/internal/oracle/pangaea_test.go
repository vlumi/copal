package oracle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vlumi/copal/internal/golden"
	"github.com/vlumi/copal/survey/blueprint"
)

// TestPangaea replays the cases recorded in oracle/pangaea and compares
// what each reaches and writes, and the coverage, with what was recorded:
// a change in a program's behavior, or in the survey's tables, shows as a
// difference. Run with -update to record what the cases now do; `make
// oracle` searches for the cases again.
func TestPangaea(t *testing.T) {
	t.Chdir("../../..")
	programs, err := filepath.Glob("pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(programs, "pangaea/copybooks")
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
		dir := filepath.Join("oracle/pangaea", p.Name)
		result, err := Replay(system, p, "pangaea/copybooks", scratch, dir)
		if err != nil {
			t.Fatalf("%s: %v", p.Name, err)
		}
		golden.CompareFiles(t, result.Files("../../../blueprints/pangaea/rules"), dir)
		results = append(results, result)
	}
	want, err := os.ReadFile("oracle/pangaea/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if got := Index("pangaea", results); string(got) != string(want) {
		t.Errorf("oracle/pangaea/README.md differs:\n%s", got)
	}
}
