package generate

import "math/rand/v2"

// Fossil names for fictional customers.
var (
	surnames = []string{
		"AMMONITE", "BELEMNITE", "CRINOID", "DIATOM", "EURYPTERID", "FUSULINID", "GRAPTOLITE",
		"HYOLITH", "ICHNITE", "JAWLESS", "KOMODO", "LYCOPSID", "MEGALODON", "NAUTILOID",
		"OSTRACOD", "PTEROSAUR", "QUILLWORT", "RUDIST", "SPRIGGINA", "TRILOBITE", "URCHIN",
		"VOLUTE", "WAPTIA", "XIPHOSURAN", "YUNNANOZOON", "ZOSTEROPHYLL",
	}
	givenNames = []string{
		"ADA", "BARNABY", "CORA", "DAN", "ELLA", "FELIX", "GRETA", "HUGO", "IDA", "JONAS",
		"KAI", "LENA", "MILO", "NORA", "OTTO", "PIA", "RUBEN", "SAGA", "TOVE", "ULLA",
		"VERA", "WILMA", "XAVER", "YRJO", "ZELDA",
	}
)

// profile is how a customer uses the bank. It is not stored anywhere: it
// follows from the seed and the customer number, every day the same.
type profile struct {
	salary   int64 // cents a month; 0 for no salary
	rent     int64 // cents a month; 0 for no rent
	landlord int   // which of the landlords the rent goes to
	activity float64
}

func (g *Generator) profile(customer int64, customers int) profile {
	r := g.rng("profile", customer)
	p := profile{activity: 0.3 + 1.4*r.Float64()}
	if r.Float64() < 0.85 {
		p.salary = cents(1800 + r.IntN(4200))
	}
	if r.Float64() < 0.6 {
		p.rent = cents(450 + r.IntN(1100))
		p.landlord = r.IntN(max(1, customers/50))
	}
	return p
}

func cents(units int) int64 { return int64(units) * 100 }

// Master makes up the account master to start from: about n accounts, in
// account number order, for customers who mostly have a current account
// and often a savings account too.
func (g *Generator) Master(n int, date Date) ([][]byte, error) {
	r := g.rng("master", 0)
	var records [][]byte
	accountNo := int64(1000000000)
	for customer := int64(1); len(records) < n; customer++ {
		name := surnames[r.IntN(len(surnames))] + " " + givenNames[r.IntN(len(givenNames))]
		opened := NewDate(1987+r.IntN(39), 1, 1).AddDays(r.IntN(365))
		if opened.Number() >= date.Number() {
			opened = date.AddDays(-30)
		}
		kinds := []string{"C"}
		if r.Float64() < 0.5 {
			kinds = append(kinds, "S")
		}
		for _, kind := range kinds {
			accountNo += int64(7 + r.IntN(40))
			record, err := g.masterRecord(r, accountNo, customer, name, kind, opened, date)
			if err != nil {
				return nil, err
			}
			records = append(records, record)
		}
	}
	return records, nil
}

func (g *Generator) masterRecord(r *rand.Rand, no, customer int64, name, kind string, opened, date Date) ([]byte, error) {
	status := "O"
	switch x := r.Float64(); {
	case x < 0.015:
		status = "F"
	case x < 0.035:
		status = "X"
	}
	var balance, limit int64
	if kind == "S" {
		balance = cents(500 + r.IntN(60000))
	} else {
		limit = cents([]int{0, 0, 250, 500, 1000, 2000}[r.IntN(6)])
		balance = cents(r.IntN(6000)) - limit/2
	}
	if status == "X" {
		balance, limit = 0, 0
	}
	if len(name) > 20 {
		name = name[:20]
	}
	rec := g.Layouts.Master.Blank()
	return rec, set(g.Layouts.Master, rec, map[string]any{
		"AM-ACCT-NO":            no,
		"AM-CUST-NO":            customer,
		"AM-SHORT-NAME":         name,
		"AM-ACCT-TYPE":          kind,
		"AM-STATUS":             status,
		"AM-OPEN-DATE":          opened.Number(),
		"AM-LAST-ACTIVITY-DATE": date.AddDays(-1 - r.IntN(30)).Number(),
		"AM-BALANCE":            balance,
		"AM-OVERDRAFT-LIMIT":    limit,
	})
}
