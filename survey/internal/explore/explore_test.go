package explore

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/survey/blueprint"
)

// newPangaea serves the explorer of the First Bank of Pangaea, from the
// repository's root, where its paths start.
func newPangaea(t *testing.T) *httptest.Server {
	t.Helper()
	t.Chdir("../../..")
	paths, err := filepath.Glob("pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := filepath.Glob("pangaea/jobs/*.jcl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(append(paths, jobs...), "pangaea/copybooks", cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	x, err := New(system, "Pangaea", "oracle/pangaea", nil)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(x.Handler())
	t.Cleanup(server.Close)
	return server
}

// newGondwana serves the explorer of Gondwana Life Insurance, whose source
// is Shift_JIS and whose jobs are batch files, with its cases.
func newGondwana(t *testing.T) *httptest.Server {
	t.Helper()
	t.Chdir("../../..")
	paths, err := filepath.Glob("gondwana/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := filepath.Glob("gondwana/scripts/*.BAT")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(append(paths, jobs...), "gondwana/copybooks", cobol.Fujitsu)
	if err != nil {
		t.Fatal(err)
	}
	x, err := New(system, "Gondwana", "oracle/gondwana", nil)
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
	seen := crawl(t, newPangaea(t))
	for _, want := range []string{"/jobs/DAILY", "/jobs/DAILY/steps/POST", "/datasets/ACCT.MASTER", "/programs/TXPOST",
		"/programs/TXPOST/paragraphs/3200-POST-DEBIT", "/programs/TXPOST/fields/AM-BALANCE", "/copybooks/ACCTMAST",
		"/source/pangaea/programs/TXPOST.cbl"} {
		if !seen[want] {
			t.Errorf("no link to %s among the %d pages", want, len(seen))
		}
	}
}

// Every page of Gondwana's explorer is there too, its names in Japanese
// in the paths, its batch files as jobs, and its cases.
func TestEveryLinkResolvesInGondwana(t *testing.T) {
	seen := crawl(t, newGondwana(t))
	for _, want := range []string{"/jobs/DAILY", "/jobs/DAILY/steps/UNLOAD", "/programs/GLD060",
		"/programs/GLD060/paragraphs/" + url.PathEscape("ｼｯｺｳ-RTN"), "/programs/GLD010/fields/" + url.PathEscape("PM-ｼﾞｮｳﾀｲ"),
		"/copybooks/POLMST", "/source/gondwana/programs/GLD060.cbl"} {
		if !seen[want] {
			t.Errorf("no link to %s among the %d pages", want, len(seen))
		}
	}
}

// crawl follows every link from the overview on and returns the pages
// seen, failing on one that is not there.
func crawl(t *testing.T, server *httptest.Server) map[string]bool {
	t.Helper()
	seen := map[string]bool{"/": true}
	from := map[string]string{}
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
			t.Errorf("%s, linked from %s: %s", link, from[link], resp.Status)
			continue
		}
		for _, m := range href.FindAllStringSubmatch(string(body), -1) {
			next := strings.ReplaceAll(m[1], "&amp;", "&")
			if !seen[next] && !strings.HasPrefix(next, "/static/") {
				seen[next], from[next] = true, link
				queue = append(queue, next)
			}
		}
	}
	return seen
}

func TestPages(t *testing.T) {
	server := newPangaea(t)
	for path, want := range map[string][]string{
		"/search?q=balance":                 {"AM-BALANCE", `href="/programs/TXPOST/fields/AM-BALANCE"`},
		"/jobs/DAILY":                       {"<svg", `href="/jobs/DAILY/steps/POST"`, `href="/datasets/ACCT.MASTER"`},
		"/datasets/ACCT.MASTER":             {"ACCTMAST", "MAINT", "ACCT.MASTER(0)"},
		"/programs/AUTHSRV":                 {"<svg", `href="/programs/AUTHPUR"`, "run by"},
		"/programs/TXPOST/fields/AM-STATUS": {"<h2>Rules</h2>", `href="/programs/TXPOST/paragraphs/3000-APPLY-TRAN"`, "<code>AM-CLOSED</code>"},
		"/trace?program=HOLDEXT&field=TD-AMOUNT": {"Passed by reference to", ">AUTHPUR</a> as argument 6", "Passed in by the caller as parameter 1",
			"Accepted as RQ-AMOUNT from SYSIN", "the SYSIN of AUTH in ONLINE"},
		"/programs/TXPOST/fields/AM-BALANCE": {"S9(11)V99", "Values moved into it", "3200-POST-DEBIT", "Follow its values back", "follow where its values go"},
		"/programs/TXPOST/paragraphs/0000-MAIN": {"Performs", "1000-INITIALIZE", "<h2>Decision table</h2>", "reached by case", `href="/programs/TXPOST/cases/`,
			`<li>When <code>WS-REJECT-COUNT</code> is more than 0, it performs <code>1000-INITIALIZE</code>`},
		"/programs/TXPOST/cases/01":              {"The run ended with return code", "<h3>TRAN-IN</h3>", "<h3>ACCT-OUT</h3>"},
		"/source/pangaea/copybooks/ACCTMAST.cpy": {`id="L1"`, "ACCT-MASTER-RECORD"},
		"/":                                      {`TRAN.DAILY</a></td><td><span class="badge kind-external">`, `AUTHPUR</a> <span class="badge kind-called">`},
		"/copybooks/ACCTMAST": {`AM-BALANCE</a> <span class="badge kind-amount">`, `AM-STATUS</a> <span class="badge kind-code">`,
			`AM-ACCT-NO</a> <span class="badge kind-number">`, "Kinds of field"},
		"/copybooks/ACCTONL":                 {`AO-ACCT-NO</a> <span class="badge kind-key">`},
		"/programs/TXSORT/fields/SR-ACCT-NO": {`SR-ACCT-NO <span class="badge kind-key">`},
		"/entities/ACCT-MASTER-RECORD/fields/AM-STATUS": {`href="/programs/TXPOST/paragraphs/3000-APPLY-TRAN"`,
			`href="/programs/ACCTMNT/paragraphs/2230-FREEZE"`, `href="/programs/ACCTLOAD/paragraphs/2000-LOAD-ACCOUNT"`, "<code>AM-CLOSED</code>"},
		"/entities/ACCT-ONLINE-RECORD/fields/AO-OVERDRAFT-LIMIT": {`href="/programs/AUTHPUR/fields/AO-OVERDRAFT-LIMIT"`,
			`href="/programs/ATMWDL/fields/AO-OVERDRAFT-LIMIT"`, `href="/programs/ATMBAL/fields/AO-OVERDRAFT-LIMIT"`},
		"/entities/:P:-CODE-TABLE/fields/:P:-CODE-ENTRY": {"as GS-CODE-ENTRY", "as ST-CODE-ENTRY"},
		"/entities/TRAN-RECORD": {"<h2>Links</h2>", `>TD-ACCT-NO</a>`, `of <a href="/entities/ACCT-MASTER-RECORD">ACCT-MASTER-RECORD</a>`,
			"<h2>Control totals</h2>", `the sum of <a href="/entities/TRAN-RECORD/fields/TD-AMOUNT">TD-AMOUNT</a>`},
		"/entities/REJECT-RECORD/fields/RJ-REASON": {"<th>When</th>", "<code>AM-BALANCE = 0</code> N"},
		"/entities/REJECT-RECORD":                  {"Date layouts", "<code>X(10)</code>", `>RJ-REJECT-DATE</a>, from <a href="/entities/BANK-CONTROL-RECORD/fields/BC-BUSINESS-DATE">BC-BUSINESS-DATE</a> <code>9(8)</code>`},
		"/entities":                                {`href="/entities/ACCT-MASTER-RECORD"`, `href="/datasets/ACCT.MASTER"`},
		"/entities/ACCT-MASTER-RECORD": {`>AM-OPEN-DATE</a>, from <a href="/entities/MAINT-REQUEST/fields/MR-EFFECTIVE-DATE">MR-EFFECTIVE-DATE</a> <code>9(6)</code>`, `<span class="badge kind-internal">`, `href="/programs/TXPOST"`,
			`href="/entities/ACCT-MASTER-RECORD/fields/AM-STATUS"`, "State machines", ">ACCTMNT 2230-FREEZE</text>",
			">ACCTMNT 2240-UNFREEZE</text>", ">ACCTMNT 2210-OPEN</text>", ">ACCTMNT 2220-CLOSE</text>"},
		"/entities/HOLD-RECORD/fields/HD-STATUS": {"<svg", ">AUTHREV 2000-APPROVE</text>", ">any state</text>"},
		"/entities/CARD-RECORD":                  {"<strong>No transitions.</strong> No program sets CR-STATUS"},
		"/ownership": {"<svg", `href="/entities/ACCT-MASTER-RECORD"`, "service decomposition",
			`<a href="/programs/ACCTMNT">ACCTMNT</a>, <a href="/programs/TXPOST">TXPOST</a></td>`,
			`<a href="/entities/GL-SUMMARY-RECORD">`, `<a href="/programs/GLSUM">GLSUM</a></td>`},
		"/boundary": {`href="/datasets/TRAN.DAILY"`, `href="/entities/TRAN-RECORD"`, `href="/datasets/ACCT.MAINT"`,
			`href="/datasets/BANK.RATES"`, `href="/datasets/CARD.MASTER"`, `href="/datasets/GL.SUMMARY"`, `href="/datasets/MAINT.REJECTS"`,
			`href="/datasets/TRAN.REJECTS"`, "STMTOUT", "AUTHRPT", "<code>DISPLAY AUTH-RESPONSE</code>", "in-stream data"},
		"/trace?direction=forward&program=TXVALID&field=TD-AMOUNT": {"Where TD-AMOUNT goes", `id="hop-DAILY-POST-TXPOST-TD-AMOUNT"`,
			`href="/programs/TXPOST/fields/AM-BALANCE"`, `href="#hop-MONTHLY-STMT-STMTPRT-AU-AMOUNT"`, "DD STMTOUT", "GL.SUMMARY(+1)",
			`class="box highlight"`, "bytes 57 to 69"},
		"/trace?direction=forward&program=BANKDATE&field=BC-BUSINESS-DATE": {"Where BC-BUSINESS-DATE goes"},
		"/trace?direction=forward&program=AUTHSRV&field=RS-RESULT":         {"<code>DISPLAY AUTH-RESPONSE</code>"},
		"/schedule": {"<svg", ">TRAN.CARD</text>", `<a href="/jobs/START">START</a> runs before <a href="/jobs/ONLINE">ONLINE</a>`,
			`<a href="/jobs/ONLINE">ONLINE</a> runs before <a href="/jobs/DAILY">DAILY</a>: it writes a new generation of <a href="/datasets/TRAN.CARD">TRAN.CARD</a>`,
			`<a href="/jobs/DAILY">DAILY</a> runs before <a href="/jobs/MONTHLY">MONTHLY</a>`, `href="/datasets/TRAN.AUDIT"`,
			"run in a cycle", "cannot know the calendar"},
		"/programs/AUTHPUR":     {"called program", `cases of <a href="/programs/AUTHSRV">AUTHSRV</a>, which calls it,`},
		"/programs/TXPOST":      {"main program", `<span class="badge kind-file">`, "<h2>Cases</h2>", `href="/programs/TXPOST/cases/01"`},
		"/datasets/ACCT.ONLINE": {`<span class="badge kind-internal">`, `<span class="badge kind-indexed">`},
		"/datasets/&&SORTED":    {`<span class="badge kind-temporary">`},
		"/search?q=GL.SUMMARY":  {`<span class="badge kind-output">`},
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
