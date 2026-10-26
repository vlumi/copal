// Package console is a window onto the running bank: web pages that show
// its accounts, each business day's files and what the days changed, read
// through the record layouts computed from the copybooks. It never
// writes the bank's files; its open and run buttons call the bank's own
// operations, the same ones the make targets run.
package console

import (
	"fmt"
	"strings"

	"github.com/vlumi/cobolith/record"
	"github.com/vlumi/copal/bedrock/dataset"
	"github.com/vlumi/copal/pangaea/records"
)

// Ledger reads the bank's datasets through the layouts, and its jobs'
// runs from their logs.
type Ledger struct {
	Layouts *records.Layouts
	Store   *dataset.Store
}

// Account is an account master record.
type Account struct {
	No, Customer           int64
	Name, Type, Status     string
	OpenDate, LastActivity int64
	Balance, Limit         int64 // cents
	Record                 []byte
}

// Transaction is a detail record of a transaction file.
type Transaction struct {
	ID, Account int64
	Code        string
	Amount      int64
	Description string
	Counter     int64
}

// TransactionFile is a transaction file: the header's business date and
// file sequence number, the details, and the trailer's control totals.
type TransactionFile struct {
	Date, FileSeq int64
	Source        string
	Details       []Transaction
	Trailer       Totals
}

// Totals are a transaction file's control totals: the number of details,
// the sum of their amounts and the sum of their account numbers.
type Totals struct {
	Count, Amount, Hash int64
}

// Totals computes the control totals of the details.
func (f TransactionFile) Totals() Totals {
	t := Totals{Count: int64(len(f.Details))}
	for _, d := range f.Details {
		t.Amount += d.Amount
		t.Hash += d.Account
	}
	return t
}

// Posting is an audit trail record: one transaction posted to one account.
type Posting struct {
	Date, TranID, Account int64
	Code                  string
	Amount, Before, After int64
}

// Reject is a transaction a step could not take, and why.
type Reject struct {
	Program, Reason string
	Tran            Transaction
}

// Accounts reads an account master.
func (l *Ledger) Accounts(dsn string) ([]Account, error) {
	return read(l, l.Layouts.Master, dsn, func(d *decoder) Account {
		return Account{
			No: d.number("AM-ACCT-NO"), Customer: d.number("AM-CUST-NO"), Name: d.text("AM-SHORT-NAME"),
			Type: d.text("AM-ACCT-TYPE"), Status: d.text("AM-STATUS"),
			OpenDate: d.number("AM-OPEN-DATE"), LastActivity: d.number("AM-LAST-ACTIVITY-DATE"),
			Balance: d.number("AM-BALANCE"), Limit: d.number("AM-OVERDRAFT-LIMIT"), Record: d.r,
		}
	})
}

// Postings reads an audit trail.
func (l *Ledger) Postings(dsn string) ([]Posting, error) {
	return read(l, l.Layouts.Audit, dsn, func(d *decoder) Posting {
		return Posting{
			Date: d.number("AU-BUSINESS-DATE"), TranID: d.number("AU-TRAN-ID"), Account: d.number("AU-ACCT-NO"),
			Code: d.text("AU-TRAN-CODE"), Amount: d.number("AU-AMOUNT"),
			Before: d.number("AU-BALANCE-BEFORE"), After: d.number("AU-BALANCE-AFTER"),
		}
	})
}

// Rejects reads a rejects dataset, with the transaction each one holds.
func (l *Ledger) Rejects(dsn string) ([]Reject, error) {
	return read(l, l.Layouts.Reject, dsn, func(d *decoder) Reject {
		tran := &decoder{l: l.Layouts.Tran, r: []byte(fmt.Sprintf("%-*s", l.Layouts.Tran.Length, d.text("RJ-TRAN")))}
		r := Reject{Program: d.text("RJ-PROGRAM"), Reason: d.text("RJ-REASON"), Tran: tran.transaction()}
		if d.err == nil {
			d.err = tran.err
		}
		return r
	})
}

// MaintRequest is an account maintenance request. Its customer, name and
// account type are blank unless it opens an account.
type MaintRequest struct {
	ID          int64
	Type        string
	Account     int64
	Customer    string
	Name        string
	AccountType string
	Limit       int64
}

// MaintReject is a maintenance request the batch could not apply, and why.
type MaintReject struct {
	Program, Reason string
	Request         MaintRequest
}

// MaintRejects reads the rejected maintenance requests.
func (l *Ledger) MaintRejects(dsn string) ([]MaintReject, error) {
	return read(l, l.Layouts.Reject, dsn, func(d *decoder) MaintReject {
		m := &decoder{l: l.Layouts.Maint, r: []byte(fmt.Sprintf("%-*s", l.Layouts.Maint.Length, d.text("RJ-TRAN")))}
		r := MaintReject{Program: d.text("RJ-PROGRAM"), Reason: d.text("RJ-REASON"), Request: MaintRequest{
			ID: m.number("MR-REQUEST-ID"), Type: m.text("MR-REQUEST-TYPE"), Account: m.number("MR-ACCT-NO"),
			Customer: m.text("MR-CUST-NO"), Name: m.text("MR-SHORT-NAME"), AccountType: m.text("MR-ACCT-TYPE"),
			Limit: m.number("MR-OVERDRAFT-LIMIT"),
		}}
		if d.err == nil {
			d.err = m.err
		}
		return r
	})
}

// TransactionFile reads a transaction file.
func (l *Ledger) TransactionFile(dsn string) (TransactionFile, error) {
	var f TransactionFile
	_, err := read(l, l.Layouts.Tran, dsn, func(d *decoder) struct{} {
		switch kind := d.text("TR-REC-TYPE"); kind {
		case "H":
			f.Date, f.Source, f.FileSeq = d.number("TH-BUSINESS-DATE"), d.text("TH-SOURCE"), d.number("TH-FILE-SEQ")
		case "D":
			f.Details = append(f.Details, d.transaction())
		case "T":
			f.Trailer = Totals{d.number("TT-DETAIL-COUNT"), d.number("TT-AMOUNT-TOTAL"), d.number("TT-ACCT-HASH")}
		default:
			d.err = fmt.Errorf("record type %q", kind)
		}
		return struct{}{}
	})
	return f, err
}

// decoder reads the fields of one record and keeps the first error.
type decoder struct {
	l   *record.Layout
	r   []byte
	err error
}

func (d *decoder) number(name string) int64 {
	n, err := d.l.Number(d.r, name)
	if d.err == nil {
		d.err = err
	}
	return n
}

func (d *decoder) text(name string) string {
	s, err := d.l.Text(d.r, name)
	if d.err == nil {
		d.err = err
	}
	return s
}

// transaction reads a transaction detail. Its counter account is blank
// unless it is a transfer.
func (d *decoder) transaction() Transaction {
	t := Transaction{
		ID: d.number("TD-TRAN-ID"), Account: d.number("TD-ACCT-NO"), Code: d.text("TD-TRAN-CODE"),
		Amount: d.number("TD-AMOUNT"), Description: d.text("TD-DESCRIPTION"),
	}
	if strings.TrimSpace(d.text("TD-COUNTER-ACCT-NO")) != "" {
		t.Counter = d.number("TD-COUNTER-ACCT-NO")
	}
	return t
}

// read decodes every record of a dataset.
func read[T any](l *Ledger, layout *record.Layout, dsn string, decode func(*decoder) T) ([]T, error) {
	records, err := l.Store.Read(dsn)
	if err != nil {
		return nil, err
	}
	out := make([]T, 0, len(records))
	for i, r := range records {
		d := &decoder{l: layout, r: r}
		v := decode(d)
		if d.err != nil {
			return nil, fmt.Errorf("%s record %d: %w", dsn, i+1, d.err)
		}
		out = append(out, v)
	}
	return out, nil
}
