// Package bank runs the First Bank of Pangaea day after day: it resets the
// bank to a generated starting point, and for each business day generates
// the branches' transaction file and runs the daily batch on it.
package bank

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/generate"
	"github.com/vlumi/copal/pangaea/jcl"
	"github.com/vlumi/copal/pangaea/runner"
	"github.com/vlumi/copal/record"
)

// StartDate is the first business day of every reset bank.
var StartDate = generate.NewDate(2026, 9, 1)

// Bank is a bank in a directory: its datasets, every generation of them,
// under datasets/, and each day's job log and reports under sysout/DATE/.
type Bank struct {
	Dir       string
	Programs  string // the compiled COBOL programs
	Copybooks string
	DailyJob  string
}

// Day is the outcome of one business day.
type Day struct {
	Date         generate.Date
	Transactions int
	Results      []runner.Result
	Rejects      int
}

const seedFile = "SEED"

// ErrNoBank is for a directory where no bank has been reset yet.
var ErrNoBank = errors.New("no bank here yet; reset one first")

// Reset replaces the bank with a new one: about the given number of
// accounts, generated from the seed, on the start date, with the rate
// table and nothing accrued yet.
func (b Bank) Reset(seed uint64, accounts int) error {
	if err := os.RemoveAll(b.Dir); err != nil {
		return err
	}
	store, err := dataset.Open(b.datasets())
	if err != nil {
		return err
	}
	g, err := b.generator(seed)
	if err != nil {
		return err
	}
	master, err := g.Master(accounts, StartDate)
	if err != nil {
		return err
	}
	control, err := g.ControlRecord(StartDate)
	if err != nil {
		return err
	}
	opened, err := g.ReadAccounts(master)
	if err != nil {
		return err
	}
	accruals, err := g.Accruals(opened, StartDate)
	if err != nil {
		return err
	}
	rates, err := g.Rates()
	if err != nil {
		return err
	}
	for _, d := range []struct {
		dsn     string
		layout  *record.Layout
		records [][]byte
	}{
		{dataset.GenerationName("ACCT.MASTER", 1), g.Layouts.Master, master},
		{dataset.GenerationName("BANK.CONTROL", 1), g.Layouts.Control, [][]byte{control}},
		{dataset.GenerationName("ACCT.ACCRUAL", 1), g.Layouts.Accrual, accruals},
		{"BANK.RATES", g.Layouts.Rate, rates},
	} {
		if err := store.Write(d.dsn, d.layout.Length, d.records); err != nil {
			return err
		}
	}
	return os.WriteFile(filepath.Join(b.Dir, seedFile), []byte(strconv.FormatUint(seed, 10)+"\n"), 0o644)
}

// RunDays runs the next n business days, reporting each as it ends. It
// stops at a day whose batch did not go through.
func (b Bank) RunDays(n int, report func(Day)) error {
	seed, err := b.seed()
	if err != nil {
		return err
	}
	g, err := b.generator(seed)
	if err != nil {
		return err
	}
	job, err := b.job()
	if err != nil {
		return err
	}
	for range n {
		day, err := b.runDay(g, job)
		if err != nil {
			return err
		}
		report(day)
	}
	return nil
}

func (b Bank) runDay(g *generate.Generator, job *jcl.Job) (Day, error) {
	store, err := dataset.Open(b.datasets())
	if err != nil {
		return Day{}, err
	}
	date, err := b.businessDate(g, store)
	if err != nil {
		return Day{}, err
	}
	master, err := store.Read(latest(store, "ACCT.MASTER"))
	if err != nil {
		return Day{}, err
	}
	accounts, err := g.ReadAccounts(master)
	if err != nil {
		return Day{}, err
	}
	fileSeq := store.LatestGeneration("TRAN.DAILY") + 1
	file, err := g.Day(date, int64(fileSeq), accounts)
	if err != nil {
		return Day{}, err
	}
	if err := store.Write(dataset.GenerationName("TRAN.DAILY", fileSeq), g.Layouts.Tran.Length, file); err != nil {
		return Day{}, err
	}

	r := runner.Runner{Programs: b.Programs, Datasets: store, Sysout: b.sysout(date)}
	results, err := r.Run(job)
	day := Day{Date: date, Transactions: len(file) - 2, Results: results}
	if err != nil {
		return day, fmt.Errorf("%s: %w", date, err)
	}
	for _, result := range results {
		if result.RC > 4 || result.Bypassed || result.Abended {
			return day, fmt.Errorf("%s: the daily batch did not go through, at step %s; see %s", date, result.Step, r.Sysout)
		}
	}
	rejects, err := store.Read(latest(store, "TRAN.REJECTS"))
	day.Rejects = len(rejects)
	return day, err
}

func (b Bank) businessDate(g *generate.Generator, store *dataset.Store) (generate.Date, error) {
	return b.controlDate(g, store, latest(store, "BANK.CONTROL"))
}

func (b Bank) controlDate(g *generate.Generator, store *dataset.Store, dsn string) (generate.Date, error) {
	control, err := store.Read(dsn)
	if err != nil {
		return generate.Date{}, err
	}
	if len(control) != 1 {
		return generate.Date{}, fmt.Errorf("%s: %d records, want 1", dsn, len(control))
	}
	n, err := g.Layouts.Control.Number(control[0], "BC-BUSINESS-DATE")
	if err != nil {
		return generate.Date{}, err
	}
	return generate.DateFromNumber(n)
}

func latest(store *dataset.Store, base string) string {
	return dataset.GenerationName(base, store.LatestGeneration(base))
}

func (b Bank) datasets() string { return filepath.Join(b.Dir, "datasets") }

func (b Bank) seed() (uint64, error) {
	data, err := os.ReadFile(filepath.Join(b.Dir, seedFile))
	if errors.Is(err, os.ErrNotExist) {
		return 0, fmt.Errorf("%s: %w", b.Dir, ErrNoBank)
	} else if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

func (b Bank) generator(seed uint64) (*generate.Generator, error) {
	layouts, err := generate.LoadLayouts(b.Copybooks)
	if err != nil {
		return nil, err
	}
	return &generate.Generator{Layouts: layouts, Seed: seed}, nil
}

func (b Bank) job() (*jcl.Job, error) {
	data, err := os.ReadFile(b.DailyJob)
	if err != nil {
		return nil, err
	}
	return jcl.Parse(bytes.NewReader(data))
}
