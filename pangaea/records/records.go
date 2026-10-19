// Package records loads the record layouts of the First Bank of Pangaea,
// computed from its copybooks, for the Go code around the bank to read and
// write its records by field, and says which layout describes each of the
// bank's datasets.
package records

import (
	"path/filepath"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/record"
)

// Layouts are the bank's record layouts.
type Layouts struct {
	Master, Tran, Control, Rate, Accrual, Maint, Audit, Reject, Card, Ledger *record.Layout
	// Request and Response are the online authorizations' messages.
	Request, Response *record.Layout
}

// Load reads the layouts from the bank's copybooks.
func Load(copybookDir string) (*Layouts, error) {
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
		{&l.Audit, "AUDITREC", "AUDIT-RECORD"},
		{&l.Reject, "REJREC", "REJECT-RECORD"},
		{&l.Card, "CARDREC", "CARD-RECORD"},
		{&l.Ledger, "GLREC", "GL-SUMMARY-RECORD"},
		{&l.Request, "AUTHMSG", "AUTH-REQUEST"},
		{&l.Response, "AUTHMSG", "AUTH-RESPONSE"},
	} {
		layout, err := record.FromCopybook(filepath.Join(copybookDir, c.copybook+".cpy"), []string{copybookDir}, c.name, cobol.IBM)
		if err != nil {
			return nil, err
		}
		*c.target = layout
	}
	return &l, nil
}

// Datasets maps each of the bank's datasets, by name without generation,
// to the layout of its records. The indexed online files, which only the
// programs read, are left out.
func (l *Layouts) Datasets() map[string]*record.Layout {
	return map[string]*record.Layout{
		"ACCT.MASTER":   l.Master,
		"TRAN.DAILY":    l.Tran,
		"TRAN.CARD":     l.Tran,
		"BANK.CONTROL":  l.Control,
		"BANK.RATES":    l.Rate,
		"ACCT.ACCRUAL":  l.Accrual,
		"ACCT.MAINT":    l.Maint,
		"TRAN.AUDIT":    l.Audit,
		"TRAN.REJECTS":  l.Reject,
		"MAINT.REJECTS": l.Reject,
		"CARD.MASTER":   l.Card,
		"GL.SUMMARY":    l.Ledger,
	}
}
