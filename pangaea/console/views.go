package console

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/layout"
	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/pangaea/generate"
)

// Overview is the bank as it stands, and the days it has run, newest first.
type Overview struct {
	History  bank.History
	Balances Balances
	Days     []DaySummary
}

// DaySummary is one run day in a line.
type DaySummary struct {
	Day          bank.RunDay
	Transactions int
	Postings     int
	Rejects      int
	Steps        []Step
}

// Step is a job step's line in the job log.
type Step struct {
	Name, Program, Outcome string
}

// Worst is the outcome of the step that did worst: the highest return code
// (RC), or a step that was bypassed or abended.
func (s DaySummary) Worst() string {
	worst, rank := "", -1
	for _, step := range s.Steps {
		r := 10000
		if rc, ok := strings.CutPrefix(step.Outcome, "RC="); ok {
			r, _ = strconv.Atoi(rc)
		}
		if r > rank {
			worst, rank = step.Outcome, r
		}
	}
	return worst
}

// Overview reads the bank's state and a summary of each day.
func (l *Ledger) Overview(h bank.History) (Overview, error) {
	accounts, err := l.Accounts(h.Master)
	if err != nil {
		return Overview{}, err
	}
	o := Overview{History: h, Balances: balances(accounts)}
	for _, d := range slices.Backward(h.Days) {
		s, err := l.summary(d)
		if err != nil {
			return Overview{}, err
		}
		o.Days = append(o.Days, s)
	}
	return o, nil
}

func (l *Ledger) summary(d bank.RunDay) (DaySummary, error) {
	s := DaySummary{Day: d}
	for _, c := range []struct {
		count *int
		dsn   string
	}{
		{&s.Transactions, d.Transactions()},
		{&s.Postings, d.Audit()},
		{&s.Rejects, d.Rejects()},
	} {
		records, err := l.Store.Read(c.dsn)
		if err != nil {
			return DaySummary{}, err
		}
		*c.count = len(records)
	}
	s.Transactions -= 2 // the header and the trailer
	var err error
	s.Steps, err = jobSteps(d.Sysout)
	return s, err
}

// jobSteps reads the steps from a job log.
func jobSteps(sysout string) ([]Step, error) {
	data, err := os.ReadFile(filepath.Join(sysout, "JOBLOG"))
	if err != nil {
		return nil, err
	}
	var steps []Step
	for line := range strings.Lines(string(data)) {
		fields := strings.Fields(line)
		if len(fields) == 3 {
			steps = append(steps, Step{fields[0], fields[1], fields[2]})
		}
	}
	return steps, nil
}

// Search finds the accounts whose number starts with the query, whose
// customer number is the query, or whose name contains it.
func Search(accounts []Account, query string) []Account {
	query = strings.ToUpper(strings.TrimSpace(query))
	if query == "" {
		return accounts
	}
	var found []Account
	for _, a := range accounts {
		no := strconv.FormatInt(a.No, 10)
		if strings.HasPrefix(no, query) || strconv.FormatInt(a.Customer, 10) == query || strings.Contains(a.Name, query) {
			found = append(found, a)
		}
	}
	return found
}

// AccountPage is an account's current record, field by field, and its
// history: every posting and every reject since the reset.
type AccountPage struct {
	Account  Account
	Fields   []FieldValue
	Postings []Posting
	Rejects  []DatedReject
}

// FieldValue is a field of a record: where it is, its bytes, and what they
// mean.
type FieldValue struct {
	Field      *layout.Field
	Raw, Value string
}

// DatedReject is a reject and the day it happened.
type DatedReject struct {
	Date generate.Date
	Reject
}

// ErrNotFound is for an account or day the bank does not have.
var ErrNotFound = errors.New("not found")

// Account reads an account and its history.
func (l *Ledger) Account(h bank.History, no int64) (AccountPage, error) {
	accounts, err := l.Accounts(h.Master)
	if err != nil {
		return AccountPage{}, err
	}
	i := slices.IndexFunc(accounts, func(a Account) bool { return a.No == no })
	if i < 0 {
		return AccountPage{}, fmt.Errorf("account %d: %w", no, ErrNotFound)
	}
	p := AccountPage{Account: accounts[i]}
	if p.Fields, err = l.fields(accounts[i].Record); err != nil {
		return AccountPage{}, err
	}
	for _, d := range h.Days {
		postings, err := l.Postings(d.Audit())
		if err != nil {
			return AccountPage{}, err
		}
		for _, posting := range postings {
			if posting.Account == no {
				p.Postings = append(p.Postings, posting)
			}
		}
		rejects, err := l.Rejects(d.Rejects())
		if err != nil {
			return AccountPage{}, err
		}
		for _, r := range rejects {
			if r.Tran.Account == no || r.Tran.Counter == no {
				p.Rejects = append(p.Rejects, DatedReject{d.Date, r})
			}
		}
	}
	return p, nil
}

// fields shows a master record through its layout.
func (l *Ledger) fields(r []byte) ([]FieldValue, error) {
	var values []FieldValue
	for _, f := range l.Layouts.Master.Fields() {
		v := FieldValue{Field: f, Raw: string(r[f.Start-1 : f.Start-1+f.Length])}
		if f.Category == layout.Numeric {
			n, err := l.Layouts.Master.Number(r, f.Name)
			if err != nil {
				return nil, err
			}
			v.Value = decimal(n, f.Scale)
		} else if f.Name != "FILLER" {
			v.Value = strings.TrimRight(v.Raw, " ")
		}
		values = append(values, v)
	}
	return values, nil
}

// DayPage is everything one business day read and wrote.
type DayPage struct {
	Day      bank.RunDay
	Next     generate.Date
	File     TransactionFile
	Computed Totals
	Postings []Posting
	Rejects  []Reject
	Steps    []Step
	Reports  []Report
}

// Report is a report a step wrote.
type Report struct {
	Name, Text string
}

// Day reads a run day's files.
func (l *Ledger) Day(h bank.History, date string) (DayPage, error) {
	i := slices.IndexFunc(h.Days, func(d bank.RunDay) bool { return d.Date.String() == date })
	if i < 0 {
		return DayPage{}, fmt.Errorf("business day %s: %w", date, ErrNotFound)
	}
	p := DayPage{Day: h.Days[i], Next: h.Date}
	if i+1 < len(h.Days) {
		p.Next = h.Days[i+1].Date
	}
	var err error
	if p.File, err = l.TransactionFile(p.Day.Transactions()); err != nil {
		return DayPage{}, err
	}
	p.Computed = p.File.Totals()
	if p.Postings, err = l.Postings(p.Day.Audit()); err != nil {
		return DayPage{}, err
	}
	if p.Rejects, err = l.Rejects(p.Day.Rejects()); err != nil {
		return DayPage{}, err
	}
	if p.Steps, err = jobSteps(p.Day.Sysout); err != nil {
		return DayPage{}, err
	}
	p.Reports, err = reports(p.Day.Sysout, p.Steps)
	return p, err
}

// reports reads the reports the steps wrote, in step order.
func reports(sysout string, steps []Step) ([]Report, error) {
	var out []Report
	for _, step := range steps {
		paths, err := filepath.Glob(filepath.Join(sysout, step.Name+".*"))
		if err != nil {
			return nil, err
		}
		for _, path := range paths {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			out = append(out, Report{filepath.Base(path), string(data)})
		}
	}
	return out, nil
}

// Opening is the bank at the start of a business day: the master the day
// starts from.
type Opening struct {
	Date   generate.Date
	Master string
}

// Openings lists every business day's opening, the next day's included.
func Openings(h bank.History) []Opening {
	var openings []Opening
	for _, d := range h.Days {
		openings = append(openings, Opening{d.Date, d.MasterBefore()})
	}
	return append(openings, Opening{h.Date, h.Master})
}

// Comparison is what changed in the master between two openings, and
// whether the postings in between account for it. Accounts are never
// removed from the master; closing one changes its status.
type Comparison struct {
	Openings      []Opening
	From, To      Opening
	Before, After Balances
	Posted        int64
	Postings      int
	Changes       []Change
	StatusChanges []Change
	Opened        []Account
}

// Balances are the accounts of one master and the sum of their balances.
type Balances struct {
	Accounts int
	Total    int64
}

// Reconciled reports whether the balances moved by what was posted.
func (c Comparison) Reconciled() bool {
	return c.After.Total-c.Before.Total == c.Posted
}

// Change is one account in two masters.
type Change struct {
	Before, After Account
}

// Delta is how much the balance moved.
func (c Change) Delta() int64 { return c.After.Balance - c.Before.Balance }

// Compare compares the masters at two openings, given by date. The
// accounts whose balance moved most come first.
func (l *Ledger) Compare(h bank.History, from, to string) (Comparison, error) {
	c := Comparison{Openings: Openings(h)}
	i := slices.IndexFunc(c.Openings, func(o Opening) bool { return o.Date.String() == from })
	j := slices.IndexFunc(c.Openings, func(o Opening) bool { return o.Date.String() == to })
	if i < 0 || j < 0 {
		return Comparison{}, fmt.Errorf("openings %s and %s: %w", from, to, ErrNotFound)
	}
	if i > j {
		i, j = j, i
	}
	c.From, c.To = c.Openings[i], c.Openings[j]
	before, err := l.Accounts(c.From.Master)
	if err != nil {
		return Comparison{}, err
	}
	after, err := l.Accounts(c.To.Master)
	if err != nil {
		return Comparison{}, err
	}
	c.Before, c.After = balances(before), balances(after)
	for _, d := range h.Days[i:j] {
		postings, err := l.Postings(d.Audit())
		if err != nil {
			return Comparison{}, err
		}
		c.Postings += len(postings)
		for _, p := range postings {
			c.Posted += p.Amount
		}
	}
	c.changes(before, after)
	return c, nil
}

func (c *Comparison) changes(before, after []Account) {
	was := map[int64]Account{}
	for _, a := range before {
		was[a.No] = a
	}
	for _, a := range after {
		b, ok := was[a.No]
		switch {
		case !ok:
			c.Opened = append(c.Opened, a)
		case b.Status != a.Status:
			c.StatusChanges = append(c.StatusChanges, Change{b, a})
		}
		if ok && b.Balance != a.Balance {
			c.Changes = append(c.Changes, Change{b, a})
		}
	}
	slices.SortStableFunc(c.Changes, func(x, y Change) int {
		return cmp.Compare(abs(y.Delta()), abs(x.Delta()))
	})
}

func balances(accounts []Account) Balances {
	b := Balances{Accounts: len(accounts)}
	for _, a := range accounts {
		b.Total += a.Balance
	}
	return b
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
