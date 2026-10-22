package oracle

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/internal/golden"
	"github.com/vlumi/copal/survey/blueprint"
)

// TestGondwana replays the cases recorded in oracle/gondwana, and the
// batch of Gondwana's batch files in the order `make oracle` runs them,
// and compares what each reaches and writes, and the coverage, with what
// was recorded, as TestPangaea does for Pangaea. Run with -update to
// record what the cases now do; `make oracle` searches for them again.
func TestGondwana(t *testing.T) {
	t.Chdir("../../..")
	programs, err := filepath.Glob("gondwana/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	batchFiles, err := filepath.Glob("gondwana/scripts/*.BAT")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(append(programs, batchFiles...), "gondwana/copybooks", cobol.Fujitsu)
	if err != nil {
		t.Fatal(err)
	}
	scratch := t.TempDir()
	var results []*Result
	byName := map[string]*Result{}
	for i := range system.Programs {
		p := &system.Programs[i]
		if p.Kind != blueprint.KindMain {
			continue
		}
		dir := filepath.Join("oracle/gondwana", p.Name)
		result, err := Replay(system, p, "gondwana/copybooks", scratch, dir)
		if err != nil {
			t.Fatalf("%s: %v", p.Name, err)
		}
		golden.CompareFiles(t, result.Files("../../../blueprints/gondwana/rules"), dir)
		results = append(results, result)
		byName[p.Name] = result
	}
	var jobs []*blueprint.Job
	for _, name := range []string{"OPEN", "DAILY", "BILLING", "RESULTS", "CLOSE", "ANNUAL"} {
		jobs = append(jobs, &system.Jobs[slices.IndexFunc(system.Jobs, func(j blueprint.Job) bool { return j.Name == name })])
	}
	batches, err := RunBatch(system, jobs, byName, scratch)
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range batches {
		golden.CompareFiles(t, b.Files(), filepath.Join("oracle/gondwana/jobs", b.Job))
	}
	want, err := os.ReadFile("oracle/gondwana/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if got := Index("gondwana", results, batches); string(got) != string(want) {
		t.Errorf("oracle/gondwana/README.md differs:\n%s", got)
	}
}
