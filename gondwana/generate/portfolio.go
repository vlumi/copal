package generate

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
)

// firstYear is the year the oldest policies in the extract date from.
const firstYear = 1976

// Extract makes up the extracts of about n policies and of their
// coverages the system opens with on the opening date: policies from the
// three eras and the four products, some with riders, most in force, some
// in grace, lapsed, surrendered, matured or paid out, each as its history
// would have left it in the old system.
func (g *Generator) Extract(n int, opening Date) (policyRecords, coverageRecords [][]byte, err error) {
	r := g.rng("extract", int64(n))
	issued := map[int]int{}
	policies := make([]Policy, 0, n)
	for range n {
		p := g.policy(r, opening)
		issued[p.Contract.Year()]++
		p.No = fmt.Sprintf("%04d%06d", p.Contract.Year(), issued[p.Contract.Year()])
		policies = append(policies, p)
	}
	slices.SortFunc(policies, func(a, b Policy) int { return strings.Compare(a.No, b.No) })
	for _, p := range policies {
		rec, err := g.record(p)
		if err != nil {
			return nil, nil, err
		}
		policyRecords = append(policyRecords, rec)
		recs, err := g.coverageRecords(p)
		if err != nil {
			return nil, nil, err
		}
		coverageRecords = append(coverageRecords, recs...)
	}
	return policyRecords, coverageRecords, nil
}

// policy makes up one policy as the old system holds it on the opening
// date.
func (g *Generator) policy(r *rand.Rand, opening Date) Policy {
	pr := aProduct(r)
	holder := somebody(r)
	yearsBack := int(float64(opening.Year()-firstYear) * r.Float64() * r.Float64())
	contract := opening.AddDays(-1 - r.IntN(365)).AddYears(-yearsBack)
	p := Policy{
		Product: pr.code, Mode: 1, Contract: contract, Term: pr.term,
		Name: holder.Kanji(), Kana: holder.Kana(), Sex: holder.sex,
		Address: addresses[r.IntN(len(addresses))], Beneficiary: relative(r, holder).Kanji(),
		Account: anAccount(r, holder), Agent: agents[r.IntN(len(agents))].code,
	}
	p.Mode = aMode(r)
	if yearsBack > 15 && r.IntN(100) < 4 {
		p.Agent = closedAgent
	}
	p.Birth, p.EntryAge = g.birth(r, pr, contract)
	p.Sum = g.sum(r, pr)
	p.Year = 1
	for p.Contract.AddYears(p.Year).Before(opening) {
		p.Year++
	}
	p.Status, p.StatusDate = 1, contract
	main := coverage(1, pr, p, p.Sum, opening)
	p.Term, p.Maturity = main.Term, main.Maturity
	if !main.InForce() {
		p.Status, p.StatusDate = 5, main.Ended
	}
	if p.Status == 1 {
		g.fate(r, &p, opening)
	}
	until := opening
	if p.Status > 3 {
		until = p.StatusDate
	}
	p.Coverages = append([]Coverage{main}, g.riders(r, p, pr, until)...)
	settle(&p)
	g.standing(r, &p, pr, opening)
	return p
}

// coverage makes up a coverage of a policy as its history left it by a
// day: issued with the policy, an endowment's matured at its term's end,
// and any other with a term renewed at each term's end at the attained
// age, or expired when that age passed the product's maximum.
func coverage(no int, pr product, p Policy, amount int64, until Date) Coverage {
	c := Coverage{No: no, Product: pr.code, Status: 1, Amount: amount, Contract: p.Contract, Term: pr.term, RatedAge: p.EntryAge}
	c.Premium = premium(amount, pr.rate(p.Sex, p.EntryAge), p.Mode)
	if pr.term == 0 {
		return c
	}
	c.Maturity = p.Contract.AddYears(pr.term)
	for c.Maturity.Before(until) {
		age := insuranceAge(p.Birth, c.Maturity)
		if pr.kind == 3 || age > pr.maxAge {
			c.Status, c.Ended = 2, c.Maturity
			return c
		}
		c.Premium, c.RatedAge = premium(amount, pr.rate(p.Sex, age), p.Mode), age
		c.Maturity = c.Maturity.AddYears(pr.term)
	}
	return c
}

// riders makes up the riders a policy was sold with: a hospitalization
// rider on about a third of the policies and a term rider on a fifth,
// where the rider is not of the main contract's kind and the entry age
// allows it, each as its history left it by a day.
func (g *Generator) riders(r *rand.Rand, p Policy, main product, until Date) []Coverage {
	var riders []Coverage
	for _, sold := range []struct {
		code  string
		share float64
	}{{"HR", 0.35}, {"TR", 0.2}} {
		pr := productOf(sold.code)
		if !chance(r, sold.share) || pr.kind == main.kind || p.EntryAge < pr.minAge || p.EntryAge > pr.maxAge {
			continue
		}
		riders = append(riders, coverage(len(riders)+2, pr, p, g.sum(r, pr), until))
	}
	return riders
}

// settle gives a policy the premium of its coverages, and ends those
// still in force of a policy that has ended, the way it ended: cancelled
// with a surrender, expired with a maturity, ended by a death.
func settle(p *Policy) {
	p.Premium = 0
	for _, c := range p.Coverages {
		if c.InForce() || c.No == 1 {
			p.Premium += c.Premium
		}
	}
	ended, ok := map[int]int{4: 3, 5: 2, 6: 4}[p.Status]
	if !ok {
		return
	}
	for i := range p.Coverages {
		if c := &p.Coverages[i]; c.InForce() {
			c.Status, c.Ended = ended, p.StatusDate
		}
	}
}

// birth makes up a birth date giving an entry age within the product's
// range at the contract date.
func (g *Generator) birth(r *rand.Rand, pr product, contract Date) (Date, int) {
	age := int(35 + r.NormFloat64()*12)
	age = min(max(age, pr.minAge), pr.maxAge)
	for {
		birth := contract.AddMonths(-(12*age + r.IntN(12))).AddDays(-r.IntN(28))
		if got := insuranceAge(birth, contract); got >= pr.minAge && got <= pr.maxAge {
			return birth, got
		}
	}
}

// sum makes up a sum insured within the product's range, in the units a
// policy of it is usually sold in, smaller sums more often.
func (g *Generator) sum(r *rand.Rand, pr product) int64 {
	unit := map[string]int64{"WL": 1_000_000, "TM": 5_000_000, "EN": 500_000, "MD": 1_000_000, "TR": 1_000_000, "HR": 1_000_000}[pr.code]
	steps := int((pr.maxSum - pr.minSum) / unit)
	n := int(float64(steps+1) * r.Float64() * r.Float64() * r.Float64())
	return min(pr.minSum+int64(n)*unit, pr.maxSum)
}

// fate draws what became of a covered policy before the opening: most are
// in force, a few in grace, some lapsed, surrendered or paid out on a
// death.
func (g *Generator) fate(r *rand.Rand, p *Policy, opening Date) {
	months := MonthsBetween(p.Contract, opening)
	if months < 6 {
		return
	}
	switch n := r.IntN(1000); {
	case n < 30:
		p.Status, p.StatusDate = 4, opening.AddDays(-1-r.IntN(min(5*365, 30*months)))
	case n < 50:
		p.Status, p.StatusDate = 6, opening.AddDays(-1-r.IntN(min(3*365, 30*months)))
	case n < 110:
		p.Status = 3
		p.LapseDate = opening.AddDays(-30 - r.IntN(min(420, 30*months-30))).BusinessDayOnOrAfter()
		p.GraceStart = debitDay(p.LapseDate.AddMonths(-2))
		p.StatusDate = p.LapseDate
	case n < 130:
		p.Status = 2
		p.GraceStart = debitDay(opening.AddMonths(-1 - r.IntN(2)))
		p.StatusDate = p.GraceStart
	}
}

// debitDay is the day the collection agent debits the premiums of a
// month: the 27th, or the business day after it.
func debitDay(inMonth Date) Date {
	return NewDate(inMonth.Year(), inMonth.Month(), 27).BusinessDayOnOrAfter()
}

// standing fills in what a policy's state implies: its paid-to month and
// arrears, its reserve, dividends and loan, with interest added up to its
// last anniversary, and the last premium it paid.
func (g *Generator) standing(r *rand.Rand, p *Policy, pr product, opening Date) {
	lastBilling := opening.AddMonths(-1)
	switch p.Status {
	case 2:
		lastBilling = p.GraceStart.AddMonths(-1)
		p.Arrears = MonthsBetween(p.GraceStart.MonthStart(), opening.MonthStart()) + 1
	case 3:
		lastBilling = p.GraceStart.AddMonths(-1)
		p.Arrears = 2
	case 4, 5, 6:
		lastBilling = p.StatusDate.AddMonths(-1)
	}
	p.PaidTo = paidTo(p, lastBilling)
	p.LastPaidMonth, p.LastPaid = p.PaidTo, p.Premium
	if p.Status == 2 || p.Status == 3 {
		p.Arrears = min(p.Arrears, 2)
	}
	if p.Covered() {
		p.Reserve = p.reserve()
	}
	if p.Status <= 3 && pr.dividend > 0 {
		p.Dividends = g.dividends(p, pr, opening)
	}
	if cashValue := pr.cashValue(ageBand(p.EntryAge), min(p.Year, 40)); p.Status <= 3 && cashValue > 0 && p.Year >= 5 && r.IntN(100) < 8 {
		p.Loan = p.Sum / 1000 * int64(cashValue) / 1000 * 3 / 10 / 10000 * 10000
		p.LoanSince = p.Contract.AddYears(p.Year - 1)
	}
}

// paidTo is the month a policy is paid to after the billing of a month:
// that month for a monthly payer, and for an annual or a half-yearly one
// the latest month billed of those a year or half a year apart from the
// contract month on.
func paidTo(p *Policy, billing Date) int64 {
	if p.Mode == 1 {
		return max(billing.YearMonth(), p.Contract.YearMonth())
	}
	every := 12
	if p.Mode == 3 {
		every = 6
	}
	due := p.Contract.MonthStart()
	for next := due.AddMonths(every); !next.After(billing.MonthStart()); next = next.AddMonths(every) {
		due = next
	}
	return due.YearMonth()
}

// dividends is what the fiscal years' ends since the contract credited
// to a participating policy: a rate of the sum from its third year, for
// as long as it was covered.
func (g *Generator) dividends(p *Policy, pr product, opening Date) int64 {
	until := opening
	if p.Status == 3 {
		until = p.LapseDate
	}
	var total int64
	for yearEnd := NewDate(p.Contract.Year(), time.March, 31); yearEnd.Before(until); yearEnd = yearEnd.AddYears(1) {
		if yearEnd.Before(p.Contract) {
			continue
		}
		year := 1
		for p.Contract.AddYears(year).Before(yearEnd) || p.Contract.AddYears(year) == yearEnd {
			year++
		}
		if year >= 3 {
			total += p.Sum * int64(pr.dividend) / 10000
		}
	}
	return total
}
