package insurer

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/vlumi/copal/gondwana/batch"
)

func programs(t *testing.T) string {
	t.Helper()
	bin := t.TempDir()
	if err := batch.Build("..", bin); err != nil {
		t.Fatal(err)
	}
	return bin
}

func opened(t *testing.T, programs string) Insurer {
	t.Helper()
	in := Insurer{Dir: t.TempDir(), Programs: programs, Scripts: "../scripts", Copybooks: "../copybooks"}
	if err := in.Open(7, 1000); err != nil {
		t.Fatal(err)
	}
	return in
}

// TestDaysAreReproducible opens two insurers from the same seed, runs
// eight business days in each, April 1 to April 10, 2026, the billing
// on the last, and checks that every file is the same, byte for byte,
// but for the indexed master, which only the programs read.
func TestDaysAreReproducible(t *testing.T) {
	bin := programs(t)
	var trees []map[string][]byte
	for range 2 {
		in := opened(t, bin)
		var days []Day
		if err := in.RunDays(8, func(d Day) { days = append(days, d) }); err != nil {
			t.Fatal(err)
		}
		checkDays(t, days)
		h, err := in.History()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(h.Days, days) || h.Seed != 7 || h.Next.String() != "2026-04-13" {
			t.Errorf("history %+v, want the days run, seed 7 and next 2026-04-13", h)
		}
		trees = append(trees, readTree(t, in.Dir))
	}
	if len(trees[0]) != len(trees[1]) {
		t.Fatalf("%d files in one insurer, %d in the other", len(trees[0]), len(trees[1]))
	}
	for name, data := range trees[0] {
		if !bytes.Equal(data, trees[1][name]) {
			t.Errorf("%s differs between the two insurers", name)
		}
	}
}

func checkDays(t *testing.T, days []Day) {
	t.Helper()
	wantDates := []string{"2026-04-01", "2026-04-02", "2026-04-03", "2026-04-06", "2026-04-07", "2026-04-08", "2026-04-09", "2026-04-10"}
	registered, benefits, changes, rejected, billed := 0, 0, 0, 0, 0
	for i, d := range days {
		if d.Date.String() != wantDates[i] {
			t.Errorf("day %d is %s, want %s", i+1, d.Date, wantDates[i])
		}
		if len(d.Batches) < 2 || d.Batches[0].Name != "DAILY" || len(d.Batches[0].Steps) != 5 || d.Batches[1].Name != "UNLOAD" {
			t.Errorf("%s: batches %+v", d.Date, d.Batches)
		}
		if d.RC() > 4 || d.InForce < 750 || d.InForce+d.InGrace+d.Lapsed+d.Ended < 1000 {
			t.Errorf("%s: RC %d, %d in force, %d in grace, %d lapsed, %d ended", d.Date, d.RC(), d.InForce, d.InGrace, d.Lapsed, d.Ended)
		}
		registered, benefits, changes, rejected, billed = registered+d.New, benefits+d.Benefits, changes+d.Changes, rejected+d.Rejected, billed+d.Billed
	}
	if registered == 0 || benefits == 0 || changes == 0 || billed == 0 {
		t.Errorf("%d policies registered, %d benefits, %d changes, %d rejects and %d debits billed in eight days", registered, benefits, changes, rejected, billed)
	}
}

// readTree reads every file of an insurer but the indexed master, whose
// file carries IDs of its own.
func readTree(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	files := map[string][]byte{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.HasSuffix(path, ".IDX") {
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

// TestFailedDayIsBackedOut runs a day whose lapse step fails after the
// new policies were registered, checks that the insurer backs the day
// out whole, then runs the day again and ends up as an insurer whose
// day never failed.
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
	if err := os.WriteFile(filepath.Join(broken, "GLD060"), []byte("#!/bin/sh\nexit 16\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	in := opened(t, bin)
	if err := in.RunDays(1, func(Day) {}); err != nil {
		t.Fatal(err)
	}
	before := readTree(t, in.Dir)
	in.Programs = broken
	if err := in.RunDays(1, func(Day) {}); err == nil {
		t.Fatal("a day whose lapse step fails went through")
	}
	after := readTree(t, in.Dir)
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
	for _, name := range []string{"LOG/DAILY.20260402.LOG", "IN/NEWPOL.20260402.TXT", "OUT/GLJRNL.20260402.TXT", "RPT/GLD010.20260402.RPT"} {
		if _, err := os.Stat(filepath.Join(in.Dir, "failed", "20260402", name)); err != nil {
			t.Errorf("the failed day did not keep %s: %v", name, err)
		}
	}

	in.Programs = bin
	if err := in.RunDays(1, func(Day) {}); err != nil {
		t.Fatal(err)
	}
	want, got := readTree(t, steady.Dir), readTree(t, in.Dir)
	for name := range got {
		if strings.HasPrefix(name, "failed") {
			delete(got, name)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("%d files after the rerun, %d in the insurer that never failed", len(got), len(want))
	}
	for name, data := range want {
		if !bytes.Equal(data, got[name]) {
			t.Errorf("%s differs from the insurer that never failed", name)
		}
	}
}
