package processor

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/copal/laurasia/night"
)

func programs(t *testing.T) string {
	t.Helper()
	bin := t.TempDir()
	if err := night.Build("..", bin); err != nil {
		t.Fatal(err)
	}
	return bin
}

func opened(t *testing.T, programs string) Processor {
	t.Helper()
	p := Processor{Dir: t.TempDir(), Programs: programs, Scripts: "../scripts", Copybooks: "../copybooks"}
	if err := p.Open(7, 300); err != nil {
		t.Fatal(err)
	}
	return p
}

// TestDaysAreReproducible opens two processors from the same seed, runs
// five days in each, Thursday to Monday, and checks that every file is
// the same, byte for byte, but for the indexed and relative files, which
// only the programs read.
func TestDaysAreReproducible(t *testing.T) {
	bin := programs(t)
	var trees []map[string][]byte
	for range 2 {
		p := opened(t, bin)
		var days []Day
		if err := p.RunDays(5, func(d Day) { days = append(days, d) }); err != nil {
			t.Fatal(err)
		}
		checkDays(t, days)
		trees = append(trees, readTree(t, p.Dir))
	}
	if len(trees[0]) != len(trees[1]) {
		t.Fatalf("%d files in one processor, %d in the other", len(trees[0]), len(trees[1]))
	}
	for name, data := range trees[0] {
		if !bytes.Equal(data, trees[1][name]) {
			t.Errorf("%s differs between the two processors", name)
		}
	}
}

func checkDays(t *testing.T, days []Day) {
	t.Helper()
	wantDates := []string{"2026-10-01", "2026-10-02", "2026-10-03", "2026-10-04", "2026-10-05"}
	statements, rejected, payments := 0, 0, 0
	for i, d := range days {
		if d.Date.String() != wantDates[i] {
			t.Errorf("day %d is %s, want %s", i+1, d.Date, wantDates[i])
		}
		business := d.Date.BusinessDay()
		if (d.Night != nil) != business || d.Cycle == nil {
			t.Errorf("%s: %d night steps and %d cycle steps on a day whose business is %v", d.Date, len(d.Night), len(d.Cycle), business)
		}
		if business && (d.Presentments < 100 || d.Postings < d.Presentments/2 || d.RC() != 4) {
			t.Errorf("%s: %d presentments, %d postings, RC %d", d.Date, d.Presentments, d.Postings, d.RC())
		}
		if !business && (d.Presentments != 0 || d.Postings != 0 || d.RC() != 0) {
			t.Errorf("%s: %d presentments, %d postings, RC %d on a weekend", d.Date, d.Presentments, d.Postings, d.RC())
		}
		statements += d.Statements
		rejected += d.Rejected
		payments += d.Payments
	}
	if statements < 20 || rejected == 0 || payments < 10 {
		t.Errorf("%d statements, %d rejects and %d payments in five days", statements, rejected, payments)
	}
}

// readTree reads every file of a processor but the indexed and relative
// ones, whose Berkeley DB files carry file IDs of their own.
func readTree(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	files := map[string][]byte{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.HasSuffix(path, ".idx") || strings.HasSuffix(path, ".rel") {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		data, err := os.ReadFile(path)
		files[rel] = data
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// TestFailedDayIsBackedOut runs a day whose billing cycle fails after the
// night posted, checks that the processor backs the day out whole, then
// runs the day again and ends up as a processor whose day never failed.
func TestFailedDayIsBackedOut(t *testing.T) {
	bin := programs(t)
	steady := opened(t, bin)
	if err := steady.RunDays(2, func(Day) {}); err != nil {
		t.Fatal(err)
	}

	broken := t.TempDir()
	if err := os.CopyFS(broken, os.DirFS(bin)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(broken, "cycle"), []byte("#!/bin/sh\nexit 16\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := opened(t, bin)
	if err := p.RunDays(1, func(Day) {}); err != nil {
		t.Fatal(err)
	}
	before := readTree(t, p.Dir)
	p.Programs = broken
	if err := p.RunDays(1, func(Day) {}); err == nil {
		t.Fatal("a day whose cycle fails went through")
	}
	after := readTree(t, p.Dir)
	for name := range after {
		if strings.HasPrefix(name, "failed") {
			delete(after, name)
		}
	}
	if len(after) != len(before) {
		t.Errorf("%d files before the failed day, %d after it was backed out", len(before), len(after))
	}
	for name, data := range before {
		if !bytes.Equal(data, after[name]) {
			t.Errorf("%s changed in a day that was backed out", name)
		}
	}
	for _, name := range []string{"logs/night.20261002.log", "logs/cycle.20261002.log", "inbound/CLR.20261002.txt", "outbound/POSTINGS.20261002.txt"} {
		if _, err := os.Stat(filepath.Join(p.Dir, "failed", "20261002", name)); err != nil {
			t.Errorf("the failed day did not keep %s: %v", name, err)
		}
	}

	p.Programs = bin
	if err := p.RunDays(1, func(Day) {}); err != nil {
		t.Fatal(err)
	}
	want, got := readTree(t, steady.Dir), readTree(t, p.Dir)
	for name := range got {
		if strings.HasPrefix(name, "failed") {
			delete(got, name)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("%d files after the rerun, %d in the processor that never failed", len(got), len(want))
	}
	for name, data := range want {
		if !bytes.Equal(data, got[name]) {
			t.Errorf("%s differs from the processor that never failed", name)
		}
	}
}
