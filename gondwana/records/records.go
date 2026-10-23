// Package records loads the record layouts of Gondwana Life Insurance,
// computed from its copybooks in its dialect, for the Go code around the
// insurer to read and write its files by field. The field names are the
// copybooks' half-width katakana, in UTF-8 as cobolith reads them; the
// records themselves stay Shift_JIS, as the programs read and write them.
package records

import (
	"path/filepath"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/record"
)

// Layouts are the insurer's record layouts: a policy of the master, a new
// policy from the sales system, a branch's change and benefit request,
// a record of the banks' transfer file, an agent, a product, a premium
// rate, a cash value and a reserve factor, a journal entry and a payment.
type Layouts struct {
	Policy, NewPolicy, Change, Claim, Transfer *record.Layout
	Agent, Product, Rate, CashValue, Reserve   *record.Layout
	Journal, Payment                           *record.Layout
}

// Load reads the layouts from the insurer's copybooks.
func Load(copybookDir string) (*Layouts, error) {
	var l Layouts
	for _, c := range []struct {
		target   **record.Layout
		copybook string
		name     string
	}{
		{&l.Policy, "POLMST", "PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ"},
		{&l.NewPolicy, "NEWPOL", "NP-ｼﾝｹｲﾔｸ-ﾚｺｰﾄﾞ"},
		{&l.Change, "MNTREQ", "MR-ﾍﾝｺｳ-ﾚｺｰﾄﾞ"},
		{&l.Claim, "CLMREQ", "CR-ｾｲｷｭｳ-ﾚｺｰﾄﾞ"},
		{&l.Transfer, "ZENGIN", "ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ"},
		{&l.Agent, "AGTMST", "AG-ﾀﾞｲﾘﾃﾝ-ﾚｺｰﾄﾞ"},
		{&l.Product, "PRODTB", "PR-ｼｮｳﾋﾝ-ﾚｺｰﾄﾞ"},
		{&l.Rate, "RATETB", "RT-ﾘｮｳﾘﾂ-ﾚｺｰﾄﾞ"},
		{&l.CashValue, "CSVTB", "CV-ｶｲﾔｸﾍﾝﾚｲｷﾝ-ﾚｺｰﾄﾞ"},
		{&l.Reserve, "RSVTB", "RS-ｼﾞｭﾝﾋﾞｷﾝ-ﾚｺｰﾄﾞ"},
		{&l.Journal, "GLJRNL", "JR-ｼﾜｹ-ﾚｺｰﾄﾞ"},
		{&l.Payment, "PAYOUT", "PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ"},
	} {
		layout, err := record.FromCopybook(filepath.Join(copybookDir, c.copybook+".cpy"), []string{copybookDir}, c.name, cobol.Fujitsu)
		if err != nil {
			return nil, err
		}
		*c.target = layout
	}
	return &l, nil
}
