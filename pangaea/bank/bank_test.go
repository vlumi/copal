package bank

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/bedrock/dataset"
	"github.com/vlumi/copal/bedrock/run"
	"github.com/vlumi/copal/bedrock/runlog"
)

// TestDaysAreReproducible resets two banks from the same seed, runs three
// business days in each, and checks that every dataset and report is the
// same, byte for byte, but for what readTree leaves out.
func TestDaysAreReproducible(t *testing.T) {
	programs := t.TempDir()
	if _, err := (run.Build{Dialect: cobol.IBM, Sources: "../programs", Copybooks: "../copybooks", Out: programs}).Programs(); err != nil {
		t.Fatal(err)
	}
	var trees []map[string][]byte
	for range 2 {
		b := Bank{Dir: t.TempDir(), Programs: programs, Copybooks: "../copybooks", Jobs: "../jobs"}
		if err := b.Reset(7, 300); err != nil {
			t.Fatal(err)
		}
		var days []Day
		if err := b.RunDays(3, func(d Day) { days = append(days, d) }); err != nil {
			t.Fatal(err)
		}
		checkDays(t, days)
		checkHistory(t, b)
		checkRuns(t, b, days)
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

// checkRuns checks the run record: START on the day the bank opened, and
// for each day its online day, its daily batch and, on the first business
// day of the month, the month-end, with the steps the day reported.
func checkRuns(t *testing.T, b Bank, days []Day) {
	t.Helper()
	runs, err := runlog.Read(filepath.Join(b.Dir, runlog.File))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range runs {
		got = append(got, r.Date+" "+r.Job)
	}
	want := []string{"2026-09-01 START", "2026-09-01 ONLINE", "2026-09-01 DAILY", "2026-09-01 MONTHLY",
		"2026-09-02 ONLINE", "2026-09-02 DAILY", "2026-09-03 ONLINE", "2026-09-03 DAILY"}
	if !slices.Equal(got, want) {
		t.Fatalf("runs %q, want %q", got, want)
	}
	if !reflect.DeepEqual(runs[2].Steps, days[0].Results) || !reflect.DeepEqual(runs[1].Steps, days[0].Online) {
		t.Errorf("the first day's runs %+v, want its steps as reported", runs[1:3])
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

// readTree reads every file of a bank. An indexed dataset is read as
// empty: its Berkeley DB file carries a file ID of its own, and the records
// in it are those of the master it was loaded from.
func readTree(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	store, err := dataset.Open(filepath.Join(dir, "datasets"))
	if err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{}
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		if filepath.Dir(rel) == "datasets" && store.Indexed(d.Name()) {
			files[rel] = nil
			return nil
		}
		data, err := os.ReadFile(path)
		files[rel] = data
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// TestFailedDayIsBackedOut runs a day whose posting fails, and checks that
// the bank backs its batch out whole, keeping its online day, then runs the
// batch again and ends up as a bank whose day never failed.
func TestFailedDayIsBackedOut(t *testing.T) {
	programs := t.TempDir()
	if _, err := (run.Build{Dialect: cobol.IBM, Sources: "../programs", Copybooks: "../copybooks", Out: programs}).Programs(); err != nil {
		t.Fatal(err)
	}
	newBank := func(programs string) Bank {
		b := Bank{Dir: t.TempDir(), Programs: programs, Copybooks: "../copybooks", Jobs: "../jobs"}
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
	before = append(before, "TRAN.CARD.G0002V00")
	slices.Sort(before)
	if after := store.Names(); !slices.Equal(after, before) {
		t.Errorf("after the failed day the catalog holds %q, want %q, with the online day's card file", after, before)
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

// An online day authorizes requests until it is closed, and runs once a
// business day: after it, the day's batch is next.
func TestOnlineDay(t *testing.T) {
	programs := t.TempDir()
	if _, err := (run.Build{Dialect: cobol.IBM, Sources: "../programs", Copybooks: "../copybooks", Out: programs}).Programs(); err != nil {
		t.Fatal(err)
	}
	b := Bank{Dir: t.TempDir(), Programs: programs, Copybooks: "../copybooks", Jobs: "../jobs"}
	if err := b.Reset(7, 300); err != nil {
		t.Fatal(err)
	}
	g, store, err := b.open()
	if err != nil {
		t.Fatal(err)
	}
	cards, err := store.Read("CARD.MASTER.G0001V00")
	if err != nil {
		t.Fatal(err)
	}
	card, err := g.Layouts.Card.Number(cards[0], "CR-CARD-NO")
	if err != nil {
		t.Fatal(err)
	}
	request := g.Layouts.Request.Blank()
	for name, value := range map[string]int64{"RQ-CARD-NO": card, "RQ-AMOUNT": 0, "RQ-AUTH-CODE": 0, "RQ-TIMESTAMP": 20260901120000} {
		if err := g.Layouts.Request.SetNumber(request, name, value); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.Layouts.Request.SetText(request, "RQ-FUNCTION", "BALQ"); err != nil {
		t.Fatal(err)
	}

	d, err := b.OpenDay()
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		response, err := d.Authorize(request)
		if err != nil {
			t.Fatal(err)
		}
		if result, _ := g.Layouts.Response.Text(response, "RS-RESULT"); result == "" {
			t.Errorf("response %q has no result", response)
		}
	}
	if _, err := d.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Authorize(request); !errors.Is(err, ErrOnlineDayClosed) {
		t.Errorf("authorized after the close: %v", err)
	}
	if _, err := b.OpenDay(); !errors.Is(err, ErrOnlineDayRun) {
		t.Errorf("opened the online day twice: %v", err)
	}
	if err := b.RunDays(1, func(Day) {}); err != nil {
		t.Fatal(err)
	}
	next, err := b.OpenDay()
	if err != nil {
		t.Fatalf("the next day's online day did not open: %v", err)
	}
	if _, err := next.Close(); err != nil {
		t.Fatal(err)
	}
}
