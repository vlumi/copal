package generate

import (
	"cmp"
	"fmt"
	"slices"
)

// request is one account maintenance request before it is written.
type request struct {
	account  int64
	kind     string // O open, C close, F freeze, U unfreeze, L limit
	customer int64
	name     string
	savings  bool
	limit    int64
}

// Maintenance makes up the branches' account maintenance requests for a
// business date, in account number order: new customers opening accounts,
// customers opening savings, empty accounts closing once they are a few
// months old, freezes and unfreezes, overdraft limit changes, and now and
// then a request the batch rejects.
func (g *Generator) Maintenance(date Date, accounts []Account) ([][]byte, error) {
	if len(accounts) == 0 {
		return nil, nil
	}
	r := g.rng("maintenance", date.Number())
	var requests []request
	add := func(q request) { requests = append(requests, q) }

	last := accounts[len(accounts)-1]
	nextNo := last.No
	newAccount := func() int64 {
		nextNo += int64(7 + r.IntN(40))
		return nextNo
	}
	saving := map[int64]bool{}
	var maxCustomer int64
	for _, a := range accounts {
		maxCustomer = max(maxCustomer, a.Customer)
		if a.Savings {
			saving[a.Customer] = true
		}
	}

	for i := range r.IntN(3) {
		customer := maxCustomer + int64(i) + 1
		name := surnames[r.IntN(len(surnames))] + " " + givenNames[r.IntN(len(givenNames))]
		add(request{account: newAccount(), kind: "O", customer: customer, name: name, limit: cents([]int{0, 500, 1000}[r.IntN(3)])})
		if r.Float64() < 0.4 {
			add(request{account: newAccount(), kind: "O", customer: customer, name: name, savings: true})
		}
	}
	settled := date.AddDays(-90).Number()
	for _, a := range accounts {
		switch {
		case a.Status == "O" && a.Balance == 0 && a.Opened < settled && r.Float64() < 0.2:
			add(request{account: a.No, kind: "C"})
		case a.Status == "F" && r.Float64() < 0.02:
			add(request{account: a.No, kind: "U"})
		}
	}
	pick := func() Account { return accounts[r.IntN(len(accounts))] }
	if a := pick(); !a.Savings && a.Open() && !saving[a.Customer] && r.Float64() < 0.3 {
		add(request{account: newAccount(), kind: "O", customer: a.Customer, name: a.Name, savings: true})
	}
	if a := pick(); a.Open() && r.Float64() < 0.15 {
		add(request{account: a.No, kind: "F"})
	}
	for range r.IntN(3) {
		if a := pick(); !a.Savings && a.Open() {
			add(request{account: a.No, kind: "L", limit: cents([]int{0, 500, 1000, 2000, 5000}[r.IntN(5)])})
		}
	}
	if r.Float64() < 0.25 {
		add(request{account: pick().No, kind: "C"})
	}
	if r.Float64() < 0.1 {
		add(request{account: last.No + 1000 + int64(r.IntN(1000)), kind: "F"})
	}
	slices.SortStableFunc(requests, func(x, y request) int { return cmp.Compare(x.account, y.account) })
	return g.maintenanceFile(date, requests)
}

// maintenanceFile writes the requests, numbered in order.
func (g *Generator) maintenanceFile(date Date, requests []request) ([][]byte, error) {
	if len(requests) > 9999 {
		return nil, fmt.Errorf("%s: %d maintenance requests, more than a day's IDs can number", date, len(requests))
	}
	l := g.Layouts.Maint
	idBase := date.Number() % 1000000 * 10000
	var records [][]byte
	for i, q := range requests {
		record := l.Blank()
		fields := map[string]any{
			"MR-REQUEST-ID": idBase + int64(i+1), "MR-REQUEST-TYPE": q.kind, "MR-ACCT-NO": q.account,
			"MR-OVERDRAFT-LIMIT": q.limit,
		}
		if q.kind == "O" {
			kind := "C"
			if q.savings {
				kind = "S"
			}
			fields["MR-CUST-NO"], fields["MR-SHORT-NAME"], fields["MR-ACCT-TYPE"] = q.customer, q.name, kind
		}
		if err := set(l, record, fields); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
