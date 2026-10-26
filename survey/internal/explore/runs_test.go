package explore

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/bedrock/runlog"
	"github.com/vlumi/copal/bedrock/runner"
	"github.com/vlumi/copal/survey/blueprint"
)

// TestRuns serves the card processor with a run record of two days and
// checks the runs on the job pages and a run's page: the night's record
// holds the steps of the scripts it calls, as its log does, so the
// payments script's run comes from it, and the night's step that runs the
// payments script ends as that script's steps did.
func TestRuns(t *testing.T) {
	t.Chdir("../../..")
	paths, err := filepath.Glob("laurasia/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	scripts, err := filepath.Glob("laurasia/scripts/*.sh")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(append(paths, scripts...), "laurasia/copybooks", cobol.MicroFocus)
	if err != nil {
		t.Fatal(err)
	}
	runs := []runlog.Run{
		{Date: "2026-10-01", Job: "nightly", Steps: []runner.Result{{Step: "CLRIN", Program: "clrin", RC: 4}, {Step: "PAYIN", Program: "payin", RC: 8}, {Step: "ACCTUL", Program: "acctul"}}},
		{Date: "2026-10-01", Job: "cycle", Steps: []runner.Result{{Step: "CYCLE", Program: "cycle"}}},
		{Date: "2026-10-03", Job: "cycle", Steps: []runner.Result{{Step: "CYCLE", Program: "cycle"}, {Step: "ACCTUL", Program: "acctul"}}},
	}
	x, err := New(system, "Laurasia", "", runs)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(x.Handler())
	t.Cleanup(server.Close)
	for path, want := range map[string][]string{
		"/jobs/cycle":                   {`href="/jobs/cycle/runs/2026-10-03"`, `href="/jobs/cycle/runs/2026-10-01"`},
		"/jobs/nightly":                 {`href="/jobs/nightly/runs/2026-10-01">2026-10-01</a></td><td>RC=0008`},
		"/jobs/payments":                {`href="/jobs/payments/runs/2026-10-01">2026-10-01</a></td><td>RC=0008`},
		"/jobs/nightly/runs/2026-10-01": {"<svg", "RC=0004", ">PAYMENTS</a></td><td><a href=\"/jobs/payments\">payments</a></td><td>RC=0008", "not run"},
	} {
		body := get(t, server.URL+path, http.StatusOK)
		for _, w := range want {
			if !strings.Contains(body, w) {
				t.Errorf("%s: no %q", path, w)
			}
		}
	}
	get(t, server.URL+"/jobs/nightly/runs/2026-10-03", http.StatusNotFound)
	if body := get(t, server.URL+"/jobs/disputes", http.StatusOK); strings.Contains(body, "<h2>Runs</h2>") {
		t.Error("the disputes script has runs, though no step of it is in the record")
	}
}

func get(t *testing.T, url string, status int) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != status {
		t.Errorf("GET %s: %d, want %d", url, resp.StatusCode, status)
	}
	return string(body)
}
