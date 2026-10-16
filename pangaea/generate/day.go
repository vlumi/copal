package generate

import (
	"fmt"
	"math/rand/v2"
)

// detail is one transaction before it is written.
type detail struct {
	account     int64
	code        string
	amount      int64
	description string
	counter     int64
}

// day collects one business day's transactions, keeping a rough running
// balance per account so that most of what customers try would go through.
type day struct {
	r        *rand.Rand
	date     Date
	accounts []Account
	balance  map[int64]int64
	limit    map[int64]int64
	details  []detail
}

// Day makes up the transaction file the branches send for a business date:
// a header, the details in the order they arrived, and a trailer with the
// control totals. fileSeq numbers the file among all the days.
func (g *Generator) Day(date Date, fileSeq int64, accounts []Account) ([][]byte, error) {
	d := &day{
		r:        g.rng("day", date.Number()),
		date:     date,
		accounts: accounts,
		balance:  map[int64]int64{},
		limit:    map[int64]int64{},
	}
	for _, a := range accounts {
		d.balance[a.No], d.limit[a.No] = a.Balance, a.Limit
	}
	d.customerBusiness(g)
	d.mistakes()
	d.r.Shuffle(len(d.details), func(i, j int) { d.details[i], d.details[j] = d.details[j], d.details[i] })
	return g.transactionFile(date, fileSeq, d.details)
}

// customerBusiness adds what each customer does today with their accounts.
func (d *day) customerBusiness(g *Generator) {
	current, savings := map[int64]Account{}, map[int64]Account{}
	seen := map[int64]bool{}
	var customers, landlords []int64
	for _, a := range d.accounts {
		accounts := current
		if a.Savings {
			accounts = savings
		}
		if _, known := accounts[a.Customer]; !known {
			accounts[a.Customer] = a
			if !a.Savings && a.Open() {
				landlords = append(landlords, a.No)
			}
		}
		if !seen[a.Customer] {
			seen[a.Customer] = true
			customers = append(customers, a.Customer)
		}
	}
	afterPayday := d.previousBusinessDay().Payday()
	for _, c := range customers {
		a, ok := current[c]
		if !ok || a.Status == "X" {
			continue
		}
		p := g.profile(c, len(customers))
		if p.salary > 0 && d.date.Payday() {
			d.add(a.No, "DP", p.salary, "SALARY", 0)
		}
		if p.rent > 0 && d.date.FirstBusinessDayOfMonth() && len(landlords) > 0 {
			if landlord := landlords[p.landlord%len(landlords)]; landlord != a.No && d.affords(a, p.rent) {
				d.add(a.No, "TR", p.rent, "RENT", landlord)
			}
		}
		if s, ok := savings[c]; ok && afterPayday && p.salary > 0 {
			d.add(a.No, "TR", p.salary/10, "TO SAVINGS", s.No)
		}
		d.everyday(a, p, current, customers)
	}
	for _, a := range d.accounts {
		if a.Savings && a.Open() && d.r.Float64() < 0.01 {
			d.add(a.No, "WD", cents(100*(1+d.r.IntN(20))), "SAVINGS WITHDRAWAL", 0)
		}
	}
}

// everyday adds the business a customer may do on any day: cash, small
// deposits, transfers to other customers, and now and then a fee.
func (d *day) everyday(a Account, p profile, current map[int64]Account, customers []int64) {
	if d.r.Float64() < 0.25*p.activity {
		amount := cents([]int{20, 40, 50, 100, 200, 300}[d.r.IntN(6)])
		if d.affords(a, amount) {
			d.add(a.No, "WD", amount, "CASH", 0)
		}
	}
	if d.r.Float64() < 0.06*p.activity {
		d.add(a.No, "DP", cents(10+d.r.IntN(490)), "CASH DEPOSIT", 0)
	}
	if d.r.Float64() < 0.05*p.activity {
		if to, ok := current[customers[d.r.IntN(len(customers))]]; ok && to.No != a.No {
			d.add(a.No, "TR", cents(10+d.r.IntN(290)), "TRANSFER", to.No)
		}
	}
	if d.r.Float64() < 0.01 {
		d.add(a.No, "FE", 250, "SERVICE FEE", 0)
	}
}

// mistakes adds what branches get wrong now and then, about one detail in
// two hundred: an unknown code, a zero amount, an account that does not
// exist, a transfer to the account it comes from or to one that does not
// exist, a deposit to a closed account.
func (d *day) mistakes() {
	if len(d.accounts) == 0 {
		return
	}
	unknown := d.accounts[len(d.accounts)-1].No + 1000
	for range len(d.details)/200 + d.r.IntN(3) {
		a := d.accounts[d.r.IntN(len(d.accounts))]
		amount := cents(10 + d.r.IntN(200))
		switch d.r.IntN(6) {
		case 0:
			d.add(a.No, "XX", amount, "UNKNOWN CODE", 0)
		case 1:
			d.add(a.No, "WD", 0, "NO AMOUNT", 0)
		case 2:
			d.add(unknown+int64(d.r.IntN(1000)), "DP", amount, "NO SUCH ACCOUNT", 0)
		case 3:
			d.add(a.No, "TR", amount, "TO ITSELF", a.No)
		case 4:
			d.add(a.No, "TR", amount, "TO NO SUCH ACCOUNT", unknown+int64(d.r.IntN(1000)))
		case 5:
			for _, c := range d.accounts {
				if c.Status == "X" {
					d.add(c.No, "DP", amount, "TO A CLOSED ACCOUNT", 0)
					break
				}
			}
		}
	}
}

// affords reports whether a customer goes ahead with a debit: always when
// the balance covers it, sometimes anyway.
func (d *day) affords(a Account, amount int64) bool {
	return d.balance[a.No]+d.limit[a.No] >= amount || d.r.Float64() < 0.3
}

func (d *day) add(account int64, code string, amount int64, description string, counter int64) {
	d.details = append(d.details, detail{account, code, amount, description, counter})
	switch code {
	case "DP":
		d.balance[account] += amount
	case "WD", "TR", "FE":
		d.balance[account] -= amount
	}
	if code == "TR" {
		d.balance[counter] += amount
	}
}

func (d *day) previousBusinessDay() Date {
	prev := d.date.AddDays(-1)
	for !prev.BusinessDay() {
		prev = prev.AddDays(-1)
	}
	return prev
}

// transactionFile writes the header, the details numbered in order, and
// the trailer with the details' count, amount total and account hash.
func (g *Generator) transactionFile(date Date, fileSeq int64, details []detail) ([][]byte, error) {
	if len(details) > 9999 {
		return nil, fmt.Errorf("%s: %d transactions, more than a day's IDs can number", date, len(details))
	}
	header := g.Layouts.Tran.Blank()
	if err := set(g.Layouts.Tran, header, map[string]any{
		"TR-REC-TYPE": "H", "TH-BUSINESS-DATE": date.Number(), "TH-SOURCE": "BRANCHES", "TH-FILE-SEQ": fileSeq,
	}); err != nil {
		return nil, err
	}
	records := [][]byte{header}
	var total, hash int64
	idBase := date.Number() % 1000000 * 10000
	for i, dt := range details {
		r := g.Layouts.Tran.Blank()
		fields := map[string]any{
			"TR-REC-TYPE": "D", "TD-TRAN-ID": idBase + int64(i+1), "TD-ACCT-NO": dt.account,
			"TD-TRAN-CODE": dt.code, "TD-AMOUNT": dt.amount, "TD-DESCRIPTION": dt.description,
		}
		if dt.counter != 0 {
			fields["TD-COUNTER-ACCT-NO"] = dt.counter
		}
		if err := set(g.Layouts.Tran, r, fields); err != nil {
			return nil, err
		}
		records = append(records, r)
		total += dt.amount
		hash += dt.account
	}
	trailer := g.Layouts.Tran.Blank()
	if err := set(g.Layouts.Tran, trailer, map[string]any{
		"TR-REC-TYPE": "T", "TT-DETAIL-COUNT": int64(len(details)), "TT-AMOUNT-TOTAL": total, "TT-ACCT-HASH": hash,
	}); err != nil {
		return nil, err
	}
	return append(records, trailer), nil
}
