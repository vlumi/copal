package online

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/pangaea/compile"
	"github.com/vlumi/copal/pangaea/layouts"
)

func TestCents(t *testing.T) {
	for text, want := range map[string]Cents{"12.34": 1234, "12.3": 1230, "12": 1200, "0.05": 5} {
		var c Cents
		if err := c.UnmarshalText([]byte(text)); err != nil || c != want {
			t.Errorf("%s read as %d, %v; want %d", text, c, err, want)
		}
	}
	for _, text := range []string{"", "-1.00", "1.234", "1,00", ".5"} {
		var c Cents
		if err := c.UnmarshalText([]byte(text)); err == nil {
			t.Errorf("%q read as %d", text, c)
		}
	}
	for text, want := range map[string]Cents{"-00000000150.00": -15000, "    12.00": 1200, "": 0} {
		if c, err := parseEdited(text); err != nil || c != want {
			t.Errorf("%q read as %d, %v; want %d", text, c, err, want)
		}
	}
	if s := Cents(-5).String(); s != "-0.05" {
		t.Errorf("-5 cents written %s", s)
	}
}

// The server answers over HTTP while the online day is open, and that the
// bank is closed when it is not.
func TestServer(t *testing.T) {
	programs := t.TempDir()
	if _, err := compile.Programs("../programs", "../copybooks", programs); err != nil {
		t.Fatal(err)
	}
	b := bank.Bank{Dir: t.TempDir(), Programs: programs, Copybooks: "../copybooks", Jobs: "../jobs"}
	if err := b.Reset(7, 300); err != nil {
		t.Fatal(err)
	}
	l, err := layouts.Load("../copybooks")
	if err != nil {
		t.Fatal(err)
	}
	s := &Server{Bank: b, Layouts: l, Clock: func() time.Time { return time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC) }}
	api := httptest.NewServer(s.Handler())
	defer api.Close()
	card := firstCard(t, b, l)

	post := func(path, body string) (int, Response) {
		t.Helper()
		resp, err := http.Post(api.URL+path, "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = resp.Body.Close() }()
		var r Response
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			t.Fatal(err)
		}
		return resp.StatusCode, r
	}
	if status, r := post("/balance", `{"card":"`+card+`","atm":"ATM 1"}`); status != http.StatusServiceUnavailable || r.Result != ResultClosed {
		t.Errorf("before opening: %d %+v", status, r)
	}
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}
	events, cancel := s.Subscribe()
	defer cancel()
	_, before := post("/balance", `{"card":"`+card+`","atm":"ATM 1"}`)
	if !before.Approved {
		t.Fatalf("balance enquiry: %+v", before)
	}
	_, bought := post("/purchase", `{"card":"`+card+`","amount":"1.50","merchant":"KIOSK"}`)
	if !bought.Approved || bought.AuthCode != "000001" || bought.Available != before.Available-150 {
		t.Errorf("purchase of 1.50 with %s available: %+v", before.Available, bought)
	}
	if e := <-events; e.Request.Function != Balance {
		t.Errorf("first event published: %+v", e)
	}
	if status, r := post("/purchase", `{"card":"12","amount":"1.50"}`); status != http.StatusBadRequest || r.Result != ResultInvalid {
		t.Errorf("a short card number: %d %+v", status, r)
	}
	<-events
	if _, err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if e := <-events; !e.Closed {
		t.Errorf("event published on closing: %+v", e)
	}
	if status, _ := post("/purchase", `{"card":"`+card+`","amount":"1.50","merchant":"KIOSK"}`); status != http.StatusServiceUnavailable {
		t.Errorf("after closing: %d", status)
	}
}

// firstCard is the number of the bank's first card that is active.
func firstCard(t *testing.T, b bank.Bank, l *layouts.Layouts) string {
	t.Helper()
	store, err := b.Datasets()
	if err != nil {
		t.Fatal(err)
	}
	cards, err := store.Read("CARD.MASTER.G0001V00")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cards {
		if status, _ := l.Card.Text(c, "CR-STATUS"); status == "A" {
			no, err := l.Card.Number(c, "CR-CARD-NO")
			if err != nil {
				t.Fatal(err)
			}
			return strconv.FormatInt(no, 10)
		}
	}
	t.Fatal("no active card")
	return ""
}
