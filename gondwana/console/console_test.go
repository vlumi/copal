package console

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/vlumi/copal/gondwana/batch"
	"github.com/vlumi/copal/gondwana/insurer"
)

// TestConsole opens an insurer and runs its first eight business days
// from the console, as from the browser, the billing on the last, then
// reads every page.
func TestConsole(t *testing.T) {
	bin := t.TempDir()
	if err := batch.Build("..", bin); err != nil {
		t.Fatal(err)
	}
	c, err := New(insurer.Insurer{Dir: t.TempDir(), Programs: bin, Scripts: "../scripts", Copybooks: "../copybooks"})
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(c.Handler())
	defer server.Close()
	client := server.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

	get(t, server.URL+"/", http.StatusOK, "No insurer yet")
	post(t, client, server.URL+"/open", url.Values{"seed": {"7"}, "policies": {"300"}}, "/")
	get(t, server.URL+"/", http.StatusOK, "The insurer before 2026-04-01", "No days run yet", "in force")
	post(t, client, server.URL+"/run", url.Values{"days": {"8"}}, "/")
	get(t, server.URL+"/", http.StatusOK, `href="/days/2026-04-10"`, "DAILY, UNLOAD, BILLING", "RC=")
	get(t, server.URL+"/days/2026-04-10", http.StatusOK, "DAILY.BAT", "BILLING.BAT", "GLM010", "OUT/PRMREQ.20260410.DAT", "RPT/GLD010.20260410.RPT", "Debits billed")
	get(t, server.URL+"/days/2026-04-04", http.StatusNotFound)
	get(t, server.URL+"/files/OUT/POLMST.20260410.TXT", http.StatusOK)
	get(t, server.URL+"/files/../SEED", http.StatusNotFound)
	policies := get(t, server.URL+"/policies", http.StatusOK, "of 3")
	first := regexp.MustCompile(`/policies/([0-9A-Z]+)`).FindStringSubmatch(policies)
	if first == nil {
		t.Fatal("no policy in the list")
	}
	get(t, server.URL+"/policies?q="+first[1], http.StatusOK, first[1])
	page := get(t, server.URL+"/policies/"+first[1], http.StatusOK, "PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ", "PM-ｼﾒｲ")
	if strings.ContainsRune(page, '�') {
		t.Error("the policy's page has bytes not read from Shift_JIS")
	}
	get(t, server.URL+"/policies/0", http.StatusNotFound)
	post(t, client, server.URL+"/open", url.Values{"seed": {"7"}, "policies": {"1"}}, "")
}

func get(t *testing.T, url string, status int, want ...string) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != status {
		t.Errorf("GET %s: %d, want %d", url, resp.StatusCode, status)
	}
	for _, w := range want {
		if !strings.Contains(string(body), w) {
			t.Errorf("GET %s: no %q in the page", url, w)
		}
	}
	return string(body)
}

// post submits a form and checks it redirects to the page given, or with
// none given, that it is refused.
func post(t *testing.T, client *http.Client, url string, form url.Values, redirect string) {
	t.Helper()
	resp, err := client.PostForm(url, form)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	switch {
	case redirect == "" && resp.StatusCode != http.StatusBadRequest:
		t.Errorf("POST %s %v: %d, want it refused", url, form, resp.StatusCode)
	case redirect != "" && (resp.StatusCode != http.StatusSeeOther || resp.Header.Get("Location") != redirect):
		t.Errorf("POST %s %v: %d to %q, want a redirect to %s", url, form, resp.StatusCode, resp.Header.Get("Location"), redirect)
	}
}
