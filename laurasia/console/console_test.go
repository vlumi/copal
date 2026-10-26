package console

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/vlumi/copal/laurasia/night"
	"github.com/vlumi/copal/laurasia/processor"
)

// TestConsole opens a processor and runs four days from the console, as
// from the browser, then reads every page.
func TestConsole(t *testing.T) {
	bin := t.TempDir()
	if err := night.Build("..", bin); err != nil {
		t.Fatal(err)
	}
	c, err := New(processor.Processor{Dir: t.TempDir(), Programs: bin, Scripts: "../scripts", Copybooks: "../copybooks"})
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(c.Handler())
	defer server.Close()
	client := server.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

	get(t, server.URL+"/", http.StatusOK, "No processor yet")
	post(t, client, server.URL+"/open", url.Values{"seed": {"7"}, "cards": {"300"}}, "/")
	get(t, server.URL+"/", http.StatusOK, "The processor before 2026-10-01", "No days run yet", "active")
	post(t, client, server.URL+"/run", url.Values{"days": {"4"}}, "/")
	get(t, server.URL+"/", http.StatusOK, `href="/days/2026-10-04"`, "no clearing file", "RC=0004")
	get(t, server.URL+"/days/2026-10-02", http.StatusOK, "The night", "CLRIN", "The billing cycle", "inbound/CLR.20261002.txt", "outbound/POSTINGS.20261002.txt")
	get(t, server.URL+"/days/2026-10-03", http.StatusOK, "No clearing file on Saturday")
	get(t, server.URL+"/days/2026-10-09", http.StatusNotFound)
	get(t, server.URL+"/files/inbound/FEES.txt", http.StatusOK)
	get(t, server.URL+"/files/../SEED", http.StatusNotFound)
	cards := get(t, server.URL+"/cards", http.StatusOK, "of 300 cards")
	first := regexp.MustCompile(`/cards/(\d+)`).FindStringSubmatch(cards)
	if first == nil {
		t.Fatal("no card in the list")
	}
	get(t, server.URL+"/cards?q="+first[1][:8], http.StatusOK, first[1])
	get(t, server.URL+"/cards/"+first[1], http.StatusOK, "CA-CARD-NO", "CA-BALANCE", "Postings")
	get(t, server.URL+"/cards/0", http.StatusNotFound)
	post(t, client, server.URL+"/run", url.Values{"days": {"0"}}, "")
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
