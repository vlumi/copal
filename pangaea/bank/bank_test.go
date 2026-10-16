package bank

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/vlumi/copal/pangaea/compile"
)

// TestDaysAreReproducible resets two banks from the same seed, runs three
// business days in each, and checks that every dataset and report is the
// same, byte for byte.
func TestDaysAreReproducible(t *testing.T) {
	programs := t.TempDir()
	if _, err := compile.Programs("../programs", "../copybooks", programs); err != nil {
		t.Fatal(err)
	}
	var trees []map[string][]byte
	for range 2 {
		b := Bank{Dir: t.TempDir(), Programs: programs, Copybooks: "../copybooks", DailyJob: "../jobs/DAILY.jcl", MonthlyJob: "../jobs/MONTHLY.jcl"}
		if err := b.Reset(7, 300); err != nil {
			t.Fatal(err)
		}
		var days []Day
		if err := b.RunDays(3, func(d Day) { days = append(days, d) }); err != nil {
			t.Fatal(err)
		}
		checkDays(t, days)
		checkHistory(t, b)
		trees = append(trees, readTree(t, b.Dir))
	}
	if len(trees[0]) != len(trees[1]) {
		t.Fatalf("%d files in one bank, %d in the other", len(trees[0]), len(trees[1]))
	}
	for name, data := range trees[0] {
		if !bytes.Equal(data, trees[1][name]) {
			t.Errorf("%s differs between the two banks", name)
		}
	}
}

func checkHistory(t *testing.T, b Bank) {
	t.Helper()
	h, err := b.History()
	if err != nil {
		t.Fatal(err)
	}
	if h.Seed != 7 || h.Date.String() != "2026-09-04" || h.Master != "ACCT.MASTER.G0004V00" || len(h.Days) != 3 {
		t.Fatalf("history: seed %d, date %s, master %s, %d days", h.Seed, h.Date, h.Master, len(h.Days))
	}
	last := h.Days[2]
	if last.Date.String() != "2026-09-03" || last.Audit() != "TRAN.AUDIT.G0003V00" || last.MasterAfter() != h.Master {
		t.Errorf("last day: %s, audit %s, master after %s", last.Date, last.Audit(), last.MasterAfter())
	}
	if _, err := os.Stat(filepath.Join(last.Sysout, "JOBLOG")); err != nil {
		t.Error(err)
	}
}

func checkDays(t *testing.T, days []Day) {
	t.Helper()
	wantDates := []string{"2026-09-01", "2026-09-02", "2026-09-03"}
	rejects := 0
	for i, d := range days {
		if d.Date.String() != wantDates[i] {
			t.Errorf("day %d is %s, want %s", i+1, d.Date, wantDates[i])
		}
		if d.Transactions < 50 || d.Rejects > d.Transactions/5 {
			t.Errorf("%s: %d transactions, %d rejects", d.Date, d.Transactions, d.Rejects)
		}
		rejects += d.Rejects
		if monthEnd := len(d.MonthEnd) > 0; monthEnd != (i == 0) {
			t.Errorf("%s: month-end job run %v, want it on the first business day of the month only", d.Date, monthEnd)
		}
	}
	if rejects == 0 {
		t.Error("three days without a single reject")
	}
}

func readTree(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	files := map[string][]byte{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		rel, _ := filepath.Rel(dir, path)
		files[rel] = data
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// TestFailedDayIsBackedOut runs a day whose posting fails, and checks that
// the bank backs it out whole, then runs it again and ends up as a bank
// whose day never failed.
func TestFailedDayIsBackedOut(t *testing.T) {
	programs := t.TempDir()
	if _, err := compile.Programs("../programs", "../copybooks", programs); err != nil {
		t.Fatal(err)
	}
	newBank := func(programs string) Bank {
		b := Bank{Dir: t.TempDir(), Programs: programs, Copybooks: "../copybooks", DailyJob: "../jobs/DAILY.jcl", MonthlyJob: "../jobs/MONTHLY.jcl"}
		if err := b.Reset(7, 300); err != nil {
			t.Fatal(err)
		}
		return b
	}
	steady := newBank(programs)
	if err := steady.RunDays(2, func(Day) {}); err != nil {
		t.Fatal(err)
	}

	broken := t.TempDir()
	if err := os.CopyFS(broken, os.DirFS(programs)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(broken, "TXPOST"), []byte("#!/bin/sh\nexit 16\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	b := newBank(programs)
	if err := b.RunDays(1, func(Day) {}); err != nil {
		t.Fatal(err)
	}
	store, err := b.Datasets()
	if err != nil {
		t.Fatal(err)
	}
	before := store.Names()
	b.Programs = broken
	if err := b.RunDays(1, func(Day) {}); err == nil {
		t.Fatal("a day whose posting fails went through")
	}
	if store, err = b.Datasets(); err != nil {
		t.Fatal(err)
	}
	if after := store.Names(); !slices.Equal(after, before) {
		t.Errorf("after the failed day the catalog holds %q, want %q", after, before)
	}
	if _, err := os.Stat(filepath.Join(b.Dir, "sysout", "2026-09-02.failed", "JOBLOG")); err != nil {
		t.Errorf("the failed day's job log was not kept: %v", err)
	}

	b.Programs = programs
	if err := b.RunDays(1, func(Day) {}); err != nil {
		t.Fatal(err)
	}
	want, got := readTree(t, steady.Dir), readTree(t, b.Dir)
	for name := range got {
		if strings.HasPrefix(name, filepath.Join("sysout", "2026-09-02.failed")) {
			delete(got, name)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("%d files after the rerun, %d in the bank that never failed", len(got), len(want))
	}
	for name, data := range want {
		if !bytes.Equal(data, got[name]) {
			t.Errorf("%s differs from the bank that never failed", name)
		}
	}
}
