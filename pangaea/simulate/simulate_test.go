package simulate

import (
	"context"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vlumi/copal/bedrock/run"
	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/pangaea/online"
	"github.com/vlumi/copal/pangaea/records"
)

// The simulator's requests reach an online day, which approves some and
// declines others.
func TestSimulator(t *testing.T) {
	programs := t.TempDir()
	if _, err := (run.Build{Sources: "../programs", Copybooks: "../copybooks", Out: programs}).Programs(); err != nil {
		t.Fatal(err)
	}
	b := bank.Bank{Dir: t.TempDir(), Programs: programs, Copybooks: "../copybooks", Jobs: "../jobs"}
	if err := b.Reset(7, 300); err != nil {
		t.Fatal(err)
	}
	l, err := records.Load("../copybooks")
	if err != nil {
		t.Fatal(err)
	}
	server := &online.Server{Bank: b, Layouts: l, Clock: time.Now}
	if err := server.Open(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := server.Close(); err != nil {
			t.Error(err)
		}
	}()
	api := httptest.NewServer(server.Handler())
	defer api.Close()
	store, err := b.Datasets()
	if err != nil {
		t.Fatal(err)
	}
	s, err := New(api.URL, store, l, 1)
	if err != nil {
		t.Fatal(err)
	}
	results := map[bool]int{}
	paths := map[string]int{}
	for range 200 {
		r, err := s.Send(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		results[r.Response.Approved]++
		paths[r.Path]++
	}
	if results[true] == 0 || results[false] == 0 {
		t.Errorf("%d approved, %d declined, want some of both", results[true], results[false])
	}
	for _, path := range []string{"/purchase", "/withdrawal", "/balance"} {
		if paths[path] == 0 {
			t.Errorf("no %s among %v", path, paths)
		}
	}
}

// A server that is not the online channel, answering JSON without a
// result, is an error, not a run of empty responses.
func TestSimulatorRefusesOtherServers(t *testing.T) {
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUpgradeRequired)
		_, _ = w.Write([]byte(`{"title":"Upgrade Required","status":426}`))
	}))
	defer other.Close()
	s := &Simulator{API: other.URL, Client: other.Client(), Cards: []Card{{No: 9100000000000018, Purchase: 1000}}, rand: rand.New(rand.NewPCG(1, 0))}
	if _, err := s.Send(context.Background()); err == nil || !strings.Contains(err.Error(), "not the online channel") {
		t.Errorf("sent to another server: %v", err)
	}
}
