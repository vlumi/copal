// Package layouts loads the record layouts of the First Bank of Pangaea,
// which Copal computes from its copybooks, for the Go code around the bank
// to read and write its records by field: the generator and the console.
package layouts

import (
	"path/filepath"

	"github.com/vlumi/copal/record"
)

// Layouts are the bank's record layouts.
type Layouts struct {
	Master, Tran, Control, Rate, Accrual, Maint, Audit, Reject *record.Layout
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
	} {
		layout, err := record.FromCopybook(filepath.Join(copybookDir, c.copybook+".cpy"), []string{copybookDir}, c.name)
		if err != nil {
			return nil, err
		}
		*c.target = layout
	}
	return &l, nil
}
