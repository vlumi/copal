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

// TestLaurasia replays the cases recorded in oracle/laurasia, and the
// batch of Laurasia's scripts in the order `make oracle` runs them, and
// compares what each reaches and writes, and the coverage, with what was
// recorded, as TestPangaea does for Pangaea. Run with -update to record what the
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
	byName := map[string]*Result{}
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
		byName[p.Name] = result
	}
	var scriptJobs []*blueprint.Job
	for _, name := range []string{"open", "nightly", "cycle"} {
		scriptJobs = append(scriptJobs, &system.Jobs[slices.IndexFunc(system.Jobs, func(j blueprint.Job) bool { return j.Name == name })])
	}
	batches, err := RunBatch(system, scriptJobs, byName, scratch)
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range batches {
		golden.CompareFiles(t, b.Files(), filepath.Join("oracle/laurasia/jobs", b.Job))
	}
	want, err := os.ReadFile("oracle/laurasia/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if got := Index("laurasia", results, batches); string(got) != string(want) {
		t.Errorf("oracle/laurasia/README.md differs:\n%s", got)
	}
}
