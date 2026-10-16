package generate

// rate is a line of the bank's rate sheet: for an account type, a credit
// rate from a balance floor up, or the debit rate below zero. Rates are
// annual, in millionths: 15000 is 1.5%.
type rate struct {
	accountType string
	kind        string // C credit, D debit
	floor       int64  // cents
	rate        int64
}

// rateSheet is the bank's interest rates. Current accounts earn little
// and charge for overdrafts; savings earn more the more they hold.
var rateSheet = []rate{
	{"C", "C", 0, 1000},
	{"C", "C", cents(10000), 5000},
	{"C", "D", 0, 120000},
	{"S", "C", 0, 15000},
	{"S", "C", cents(25000), 20000},
	{"S", "C", cents(100000), 25000},
	{"S", "D", 0, 150000},
}

// Rates returns the rate table.
func (g *Generator) Rates() ([][]byte, error) {
	var records [][]byte
	for _, r := range rateSheet {
		record := g.Layouts.Rate.Blank()
		if err := set(g.Layouts.Rate, record, map[string]any{
			"RT-ACCT-TYPE": r.accountType, "RT-KIND": r.kind, "RT-FLOOR": r.floor, "RT-ANNUAL-RATE": r.rate,
		}); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// Accruals returns an accrual record for each account, with nothing
// accrued yet: interest accrues from the given date on.
func (g *Generator) Accruals(accounts []Account, from Date) ([][]byte, error) {
	var records [][]byte
	for _, a := range accounts {
		record := g.Layouts.Accrual.Blank()
		if err := set(g.Layouts.Accrual, record, map[string]any{
			"AC-ACCT-NO": a.No, "AC-THROUGH-DATE": from.Number(), "AC-ACCRUED": int64(0),
			"AC-DUE": int64(0), "AC-DUE-MONTH": int64(0),
		}); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
