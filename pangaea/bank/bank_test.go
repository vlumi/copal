package bank

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
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
		b := Bank{Dir: t.TempDir(), Programs: programs, Copybooks: "../copybooks", DailyJob: "../jobs/DAILY.jcl"}
		if err := b.Reset(7, 300); err != nil {
			t.Fatal(err)
		}
		var days []Day
		if err := b.RunDays(3, func(d Day) { days = append(days, d) }); err != nil {
			t.Fatal(err)
		}
		checkDays(t, days)
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
