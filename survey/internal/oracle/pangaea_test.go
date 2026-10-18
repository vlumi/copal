package oracle

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/vlumi/copal/internal/golden"
	"github.com/vlumi/copal/survey/blueprint"
)

// TestPangaea replays the cases recorded in oracle/pangaea, and the batch
// of Pangaea's jobs in the order `make oracle` runs them, and compares
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
	jobFiles, err := filepath.Glob("pangaea/jobs/*.jcl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(append(programs, jobFiles...), "pangaea/copybooks")
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
		dir := filepath.Join("oracle/pangaea", p.Name)
		result, err := Replay(system, p, "pangaea/copybooks", scratch, dir)
		if err != nil {
			t.Fatalf("%s: %v", p.Name, err)
		}
		golden.CompareFiles(t, result.Files("../../../blueprints/pangaea/rules"), dir)
		results = append(results, result)
		byName[p.Name] = result
	}
	var jobs []*blueprint.Job
	for _, name := range []string{"START", "ONLINE", "DAILY", "MONTHLY"} {
		jobs = append(jobs, &system.Jobs[slices.IndexFunc(system.Jobs, func(j blueprint.Job) bool { return j.Name == name })])
	}
	batches, err := RunJobs(jobs, byName, scratch)
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range batches {
		golden.CompareFiles(t, b.Files(), filepath.Join("oracle/pangaea/jobs", b.Job))
	}
	want, err := os.ReadFile("oracle/pangaea/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if got := Index("pangaea", results, batches); string(got) != string(want) {
		t.Errorf("oracle/pangaea/README.md differs:\n%s", got)
	}
}
