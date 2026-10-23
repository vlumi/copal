package records

import (
	"testing"

	"github.com/vlumi/cobolith/record"
)

// TestLoad checks that the layouts come out of the copybooks with the
// lengths their banners state, and that a field is found by its katakana
// name, in UTF-8, where the copybook puts it.
func TestLoad(t *testing.T) {
	l, err := Load("../copybooks")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		layout *record.Layout
		length int
		field  string
		text   string
	}{
		{l.Policy, 300, "PM-ｼﾒｲ", "name"},
		{l.NewPolicy, 200, "NP-ｹｲﾔｸﾋﾞ", "20260401"},
		{l.Change, 180, "MR-ｼｭﾍﾞﾂ", "R"},
		{l.Claim, 80, "CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ", "R080401"},
		{l.Transfer, 120, "ZG-D-ｺｷｬｸﾊﾞﾝｺﾞｳ", "2026091001"},
		{l.Agent, 60, "AG-ｺｰﾄﾞ", "A00101"},
		{l.Product, 60, "PR-ｺｰﾄﾞ", "WL"},
		{l.Rate, 10, "RT-ｼｮｳﾋﾝ", "WL"},
		{l.CashValue, 12, "CV-ｼｮｳﾋﾝ", "EN"},
		{l.Reserve, 11, "RS-ｼｮｳﾋﾝ", "TM"},
		{l.Journal, 80, "JR-ﾌﾟﾛｸﾞﾗﾑ", "GLD010"},
		{l.Payment, 100, "PY-ｼｭﾍﾞﾂ", "D"},
	} {
		if c.layout.Length != c.length {
			t.Errorf("%s is %d bytes, want %d", c.layout.Name, c.layout.Length, c.length)
		}
		r := c.layout.Blank()
		if err := c.layout.SetText(r, c.field, c.text); err != nil {
			t.Errorf("%s: %v", c.layout.Name, err)
			continue
		}
		if got, _ := c.layout.Text(r, c.field); got != c.text {
			t.Errorf("%s.%s reads back as %q, want %q", c.layout.Name, c.field, got, c.text)
		}
	}
}
