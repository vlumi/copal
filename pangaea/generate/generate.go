// Package generate makes up the First Bank of Pangaea's customers and their
// daily business: an account master to start from, and for each business
// day the transaction file the branches send. Everything follows from a
// seed, so the same seed and date always give the same records.
package generate

import (
	"fmt"
	"math/rand/v2"
	"path/filepath"
	"slices"

	"github.com/vlumi/copal/record"
)

// Layouts are the record layouts the generator reads and writes through.
type Layouts struct {
	Master  *record.Layout
	Tran    *record.Layout
	Control *record.Layout
	Rate    *record.Layout
	Accrual *record.Layout
	Maint   *record.Layout
}

// LoadLayouts reads the layouts from Pangaea's copybooks.
func LoadLayouts(copybookDir string) (*Layouts, error) {
	var l Layouts
	for _, c := range []struct {
		target   **record.Layout
		copybook string
		name     string
	}{
		{&l.Master, "ACCTMAST", "ACCT-MASTER-RECORD"},
		{&l.Tran, "TRANREC", "TRAN-RECORD"},
		{&l.Control, "BANKCTL", "BANK-CONTROL-RECORD"},
		{&l.Rate, "RATEREC", "RATE-RECORD"},
		{&l.Accrual, "ACCRREC", "ACCRUAL-RECORD"},
		{&l.Maint, "MAINTREC", "MAINT-REQUEST"},
	} {
		layout, err := record.FromCopybook(filepath.Join(copybookDir, c.copybook+".cpy"), []string{copybookDir}, c.name)
		if err != nil {
			return nil, err
		}
		*c.target = layout
	}
	return &l, nil
}

// Generator makes up the bank's records for one seed.
type Generator struct {
	Layouts *Layouts
	Seed    uint64
}

// Account is what the generator knows of an account: what it reads from
// the account master.
type Account struct {
	No       int64
	Customer int64
	Name     string
	Savings  bool
	Status   string // O open, F frozen, X closed
	Opened   int64  // the open date, yyyymmdd
	Balance  int64  // cents
	Limit    int64  // overdraft limit, cents
}

// Open reports whether the account takes transactions.
func (a Account) Open() bool { return a.Status == "O" }

// rng returns a random source for one purpose, so that changing how one
// thing is generated does not change everything else.
func (g *Generator) rng(stream string, n int64) *rand.Rand {
	h := uint64(14695981039346656037)
	for _, c := range []byte(stream) {
		h = (h ^ uint64(c)) * 1099511628211
	}
	return rand.New(rand.NewPCG(g.Seed, h^uint64(n)))
}

// ControlRecord returns the bank control record for a business date.
func (g *Generator) ControlRecord(date Date) ([]byte, error) {
	r := g.Layouts.Control.Blank()
	return r, g.Layouts.Control.SetNumber(r, "BC-BUSINESS-DATE", date.Number())
}

// ReadAccounts reads accounts from account master records.
func (g *Generator) ReadAccounts(records [][]byte) ([]Account, error) {
	accounts := make([]Account, 0, len(records))
	for _, r := range records {
		var a Account
		var kind string
		fields := []struct {
			name   string
			number *int64
			text   *string
		}{
			{"AM-ACCT-NO", &a.No, nil},
			{"AM-CUST-NO", &a.Customer, nil},
			{"AM-OPEN-DATE", &a.Opened, nil},
			{"AM-BALANCE", &a.Balance, nil},
			{"AM-OVERDRAFT-LIMIT", &a.Limit, nil},
			{"AM-SHORT-NAME", nil, &a.Name},
			{"AM-STATUS", nil, &a.Status},
			{"AM-ACCT-TYPE", nil, &kind},
		}
		for _, f := range fields {
			var err error
			if f.number != nil {
				*f.number, err = g.Layouts.Master.Number(r, f.name)
			} else {
				*f.text, err = g.Layouts.Master.Text(r, f.name)
			}
			if err != nil {
				return nil, err
			}
		}
		a.Savings = kind == "S"
		accounts = append(accounts, a)
	}
	return accounts, nil
}

// set writes fields into a record, numbers and text alike.
func set(l *record.Layout, r []byte, fields map[string]any) error {
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		var err error
		switch v := fields[name].(type) {
		case int64:
			err = l.SetNumber(r, name, v)
		case string:
			err = l.SetText(r, name, v)
		default:
			err = fmt.Errorf("%s: %T is neither a number nor text", name, v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
