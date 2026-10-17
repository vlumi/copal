package simulate

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/pangaea/compile"
	"github.com/vlumi/copal/pangaea/online"
	"github.com/vlumi/copal/pangaea/records"
)

// The simulator's requests reach an online day, which approves some and
// declines others.
func TestSimulator(t *testing.T) {
	programs := t.TempDir()
	if _, err := compile.Programs("../programs", "../copybooks", programs); err != nil {
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
