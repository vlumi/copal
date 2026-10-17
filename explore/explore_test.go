package explore

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/vlumi/copal/blueprint"
)

// newPangaea serves the explorer of the First Bank of Pangaea, from the
// repository's root, where its paths start.
func newPangaea(t *testing.T) *httptest.Server {
	t.Helper()
	t.Chdir("..")
	paths, err := filepath.Glob("pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := filepath.Glob("pangaea/jobs/*.jcl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(append(paths, jobs...), "pangaea/copybooks")
	if err != nil {
		t.Fatal(err)
	}
	x, err := New(system, "Pangaea")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(x.Handler())
	t.Cleanup(server.Close)
	return server
}

var href = regexp.MustCompile(`href="(/[^"#]*)`)

// Every page the explorer links to, from the overview on, is there.
func TestEveryLinkResolves(t *testing.T) {
	server := newPangaea(t)
	seen := map[string]bool{"/": true}
	queue := []string{"/"}
	for len(queue) > 0 {
		link := queue[0]
		queue = queue[1:]
		resp, err := http.Get(server.URL + link)
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s: %s", link, resp.Status)
			continue
		}
		for _, m := range href.FindAllStringSubmatch(string(body), -1) {
			next := strings.ReplaceAll(m[1], "&amp;", "&")
			if !seen[next] && !strings.HasPrefix(next, "/static/") {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	for _, want := range []string{"/jobs/DAILY", "/jobs/DAILY/steps/POST", "/datasets/ACCT.MASTER", "/programs/TXPOST",
		"/programs/TXPOST/paragraphs/3200-POST-DEBIT", "/programs/TXPOST/fields/AM-BALANCE", "/copybooks/ACCTMAST",
		"/source/pangaea/programs/TXPOST.cbl"} {
		if !seen[want] {
			t.Errorf("no link to %s among the %d pages", want, len(seen))
		}
	}
}

func TestPages(t *testing.T) {
	server := newPangaea(t)
	for path, want := range map[string][]string{
		"/search?q=balance":                      {"AM-BALANCE", `href="/programs/TXPOST/fields/AM-BALANCE"`},
		"/jobs/DAILY":                            {"<svg", `href="/jobs/DAILY/steps/POST"`, `href="/datasets/ACCT.MASTER"`},
		"/datasets/ACCT.MASTER":                  {"ACCTMAST", "MAINT", "ACCT.MASTER(0)"},
		"/programs/AUTHSRV":                      {"<svg", `href="/programs/AUTHPUR"`, "run by"},
		"/programs/TXPOST/fields/AM-BALANCE":     {"S9(11)V99", "Values moved into it", "3200-POST-DEBIT", "Follow its values back"},
		"/programs/TXPOST/paragraphs/0000-MAIN":  {"Performs", "1000-INITIALIZE"},
		"/source/pangaea/copybooks/ACCTMAST.cpy": {`id="L1"`, "ACCT-MASTER-RECORD"},
		"/trace?program=STMTPRT&field=WS-RA-VALUE": {"TRAN.DAILY(0)", "which no job writes", `class="box highlight"`,
			`href="#hop-DAILY-VALID-TXVALID-TD-AMOUNT"`, `id="hop-DAILY-VALID-TXVALID-TD-AMOUNT"`, "HOLDEXT", "through the sort of SORT of DAILY"},
	} {
		resp, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s: %s", path, resp.Status)
		}
		for _, w := range want {
			if !strings.Contains(string(body), w) {
				t.Errorf("%s: no %q", path, w)
			}
		}
	}
	for _, path := range []string{"/programs/NOPE", "/jobs/DAILY/steps/NOPE", "/source/go.mod", "/programs/TXPOST/fields/NOPE",
		"/trace?program=TXPOST&field=NOPE"} {
		resp, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("%s: %s, want 404", path, resp.Status)
		}
	}
}
