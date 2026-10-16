package console

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/vlumi/copal/blueprint"
	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/pangaea/compile"
)

// TestConsole resets a bank and runs two days from the console, as from
// the browser, then reads every page.
func TestConsole(t *testing.T) {
	programs := t.TempDir()
	if _, err := compile.Programs("../programs", "../copybooks", programs); err != nil {
		t.Fatal(err)
	}
	b := bank.Bank{Dir: t.TempDir(), Programs: programs, Copybooks: "../copybooks", DailyJob: "../jobs/DAILY.jcl", MonthlyJob: "../jobs/MONTHLY.jcl"}
	sources, err := filepath.Glob("../programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(append(sources, b.DailyJob, b.MonthlyJob), "../copybooks")
	if err != nil {
		t.Fatal(err)
	}
	c, err := New(b, system, "../../blueprints/pangaea")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(c.Handler())
	defer server.Close()
	client := server.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

	get(t, server.URL+"/", http.StatusOK, "No bank yet")
	post(t, client, server.URL+"/reset", url.Values{"seed": {"3"}, "accounts": {"200"}}, "/")
	get(t, server.URL+"/", http.StatusOK, "The bank on 2026-09-01", "No days run yet")
	post(t, client, server.URL+"/run", url.Values{"days": {"2"}}, "/compare?from=2026-09-01&to=2026-09-03")
	get(t, server.URL+"/compare?from=2026-09-01&to=2026-09-03", http.StatusOK,
		"ACCT.MASTER.G0001V00", "ACCT.MASTER.G0003V00", "Reconciled")
	get(t, server.URL+"/", http.StatusOK, `href="/days/2026-09-02"`, "RC=0004")
	get(t, server.URL+"/days/2026-09-02", http.StatusOK, "TRAN.AUDIT.G0002V00", "VALID.VALRPT", "POST.POSTRPT",
		"Job DAILY", "<svg", ">program ACCTMNT<", "/blueprints/programs/ACCTMNT.md", "/blueprints/copybooks/ACCTMAST.md", " records<")
	get(t, server.URL+"/days/2026-09-01", http.StatusOK, "Job MONTHLY", "STMT.STMTRPT")
	get(t, server.URL+"/days/2026-09-02/sysout/JOBLOG", http.StatusOK, "JOB DAILY")
	get(t, server.URL+"/days/2026-09-01/sysout/MONTHLY/JOBLOG", http.StatusOK, "JOB MONTHLY")
	get(t, server.URL+"/days/2026-09-02/sysout/../SEED", http.StatusNotFound)
	get(t, server.URL+"/blueprints/programs/ACCTMNT.md", http.StatusOK, "# ACCTMNT")
	accounts := get(t, server.URL+"/accounts", http.StatusOK)
	if !regexp.MustCompile(`(\d+) of (\d+) accounts`).MatchString(accounts) {
		t.Error("the account list does not say how many accounts it shows")
	}
	first := regexp.MustCompile(`/accounts/(\d+)`).FindStringSubmatch(accounts)
	if first == nil {
		t.Fatal("no account in the list")
	}
	get(t, server.URL+"/accounts?q="+first[1], http.StatusOK, "1 of ")
	get(t, server.URL+first[0], http.StatusOK, "AM-BALANCE", "S9(11)V99")
	get(t, server.URL+"/accounts/42", http.StatusNotFound, "not found")
	get(t, server.URL+"/days/2026-09-03", http.StatusNotFound, "not found")
	post(t, client, server.URL+"/run", url.Values{"days": {"none"}}, "")

	foreign, err := http.NewRequest("POST", server.URL+"/reset", strings.NewReader("seed=1&accounts=100"))
	if err != nil {
		t.Fatal(err)
	}
	foreign.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	foreign.Header.Set("Sec-Fetch-Site", "cross-site")
	resp, err := client.Do(foreign)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("a reset from another site: %s, want it forbidden", resp.Status)
	}
}

func get(t *testing.T, url string, status int, want ...string) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body := new(strings.Builder)
	if _, err := io.Copy(body, resp.Body); err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != status {
		t.Errorf("GET %s: %s, want %d", url, resp.Status, status)
	}
	for _, w := range want {
		if !strings.Contains(body.String(), w) {
			t.Errorf("GET %s: no %q in the page", url, w)
		}
	}
	return body.String()
}

// post submits a form and checks where it redirects to, or that it is
// refused when there is nowhere.
func post(t *testing.T, client *http.Client, url string, form url.Values, redirect string) {
	t.Helper()
	resp, err := client.PostForm(url, form)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if redirect == "" {
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("POST %s %v: %s, want it refused", url, form, resp.Status)
		}
		return
	}
	if got := resp.Header.Get("Location"); resp.StatusCode != http.StatusSeeOther || got != redirect {
		t.Errorf("POST %s %v: %s to %q, want %q", url, form, resp.Status, got, redirect)
	}
}
