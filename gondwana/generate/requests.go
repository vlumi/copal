package generate

import (
	"math"
	"math/rand/v2"
	"strconv"
)

// businessDays is how many business days a year has, for the rates of
// what happens to a policy in a year.
const businessDays = 245

// perDay turns a yearly rate into a business day's.
func perDay(yearly float64) float64 { return yearly / businessDays }

// mortality is the chance a person of an age dies in a year.
func mortality(age int) float64 {
	return 0.0006 * math.Pow(1.09, float64(age-30))
}

// Claims makes up the branches' benefit requests of a business day:
// deaths by the insured's age, of illness, by accident or, one in fifty,
// by suicide, on policies covered or lapsed, so that some died before
// the lapse and some after; hospitalizations of the medical policies,
// one in fifteen keyed twice, and their surgeries, at ten, twenty or
// forty times the daily benefit; surrenders of covered and lapsed
// policies; and a request a branch
// should not have sent now and then: for no policy, of no kind, for an
// event before the contract or after today, or a hospitalization on a
// policy that pays none.
func (g *Generator) Claims(date Date, policies []Policy) ([][]byte, error) {
	r := g.rng("claims", date.Number())
	var recs [][]byte
	n := 0
	add := func(kind byte, p Policy, event Date, days int, payee string) error {
		n++
		cause := 1 + r.IntN(2)
		if kind == 'D' && r.IntN(50) == 0 {
			cause = 3
		}
		rec, err := made(g.Layouts.Claim, map[string]any{
			"CR-ｼﾞﾕｳｺｰﾄﾞ":    strconv.Itoa(cause),
			"CR-ｾｲｷｭｳﾊﾞﾝｺﾞｳ": requestNumber('C', date, n), "CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ": p.No, "CR-ｼｭﾍﾞﾂ": string(kind),
			"CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ": event.Era(), "CR-ﾆｯｽｳ": days, "CR-ｳｹﾄﾘﾆﾝｶﾅ": payee, "CR-ｼﾃﾝ": branchOf(p.Agent),
		})
		recs = append(recs, rec)
		return err
	}
	for _, p := range policies {
		if p.Status > 3 {
			continue
		}
		age := insuranceAge(p.Birth, date)
		var err error
		switch {
		case chance(r, perDay(mortality(age))):
			payee := kanaOf(p.Beneficiary)
			if payee == "" {
				payee = p.Kana
			}
			err = add('D', p, date.AddDays(-1-r.IntN(30)), 0, payee)
		case p.Covered() && p.Product == "MD" && chance(r, perDay(0.15)):
			event, days := date.AddDays(-10-r.IntN(50)), 3+r.IntN(28)
			if err = add('H', p, event, days, p.Kana); err == nil && r.IntN(15) == 0 {
				err = add('H', p, event, days, p.Kana)
			}
		case p.Covered() && p.Product == "MD" && chance(r, perDay(0.05)):
			err = add('O', p, date.AddDays(-5-r.IntN(40)), 10<<r.IntN(3), p.Kana)
		case p.Covered() && chance(r, perDay(0.04)), p.Lapsed() && chance(r, perDay(0.08)):
			err = add('S', p, date, 0, p.Kana)
		case chance(r, perDay(0.004)):
			err = g.wrongClaim(r, add, p, date)
		}
		if err != nil {
			return nil, err
		}
	}
	return recs, nil
}

// wrongClaim makes up a request the program rejects.
func (g *Generator) wrongClaim(r *rand.Rand, add func(byte, Policy, Date, int, string) error, p Policy, date Date) error {
	switch r.IntN(5) {
	case 0:
		p.No = "9999999999"
		return add('D', p, date, 0, p.Kana)
	case 4:
		return add('D', p, date.AddDays(7), 0, p.Kana)
	case 1:
		return add('X', p, date, 0, p.Kana)
	case 2:
		return add('D', p, p.Contract.AddYears(-1), 0, p.Kana)
	}
	return add('H', p, date.AddDays(-10), 5, p.Kana)
}

// Changes makes up the branches' change requests of a business day: an
// address, a bank account, a name, a beneficiary, or the reinstatement
// of a lapsed policy with the arrears paid at the counter, usually in
// full; and a request a branch got wrong now and then.
func (g *Generator) Changes(date Date, policies []Policy) ([][]byte, error) {
	r := g.rng("changes", date.Number())
	var recs [][]byte
	n := 0
	add := func(kind byte, p Policy, fields map[string]any) error {
		n++
		all := map[string]any{"MR-ｲﾗｲﾊﾞﾝｺﾞｳ": requestNumber('M', date, n), "MR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ": p.No, "MR-ｼｭﾍﾞﾂ": string(kind),
			"MR-ｷﾞﾝｺｳｺｰﾄﾞ": 0, "MR-ｼﾃﾝｺｰﾄﾞ": 0, "MR-ﾖｷﾝｼｭﾓｸ": 0, "MR-ｺｳｻﾞﾊﾞﾝｺﾞｳ": 0, "MR-ﾆｭｳｷﾝｶﾞｸ": 0}
		for k, v := range fields {
			all[k] = v
		}
		rec, err := made(g.Layouts.Change, all)
		recs = append(recs, rec)
		return err
	}
	for _, p := range policies {
		if p.Status > 3 {
			continue
		}
		var err error
		switch {
		case p.Lapsed() && DaysBetween(p.LapseDate, date) <= 365 && chance(r, perDay(0.15)):
			err = add('R', p, map[string]any{"MR-ﾆｭｳｷﾝｶﾞｸ": g.counterPayment(r, p, date)})
		case chance(r, perDay(0.06)):
			err = add('A', p, map[string]any{"MR-ｼﾞｭｳｼｮ": addresses[r.IntN(len(addresses))]})
		case chance(r, perDay(0.03)):
			a := anAccount(r, person{})
			err = add('B', p, map[string]any{"MR-ｷﾞﾝｺｳｺｰﾄﾞ": a.bank, "MR-ｼﾃﾝｺｰﾄﾞ": a.branch, "MR-ﾖｷﾝｼｭﾓｸ": a.kind, "MR-ｺｳｻﾞﾊﾞﾝｺﾞｳ": a.number, "MR-ｺｳｻﾞﾒｲｷﾞ": p.Kana})
		case p.Sex == 2 && chance(r, perDay(0.02)):
			married := somebody(r)
			married.given = givenNames[2][r.IntN(len(givenNames[2]))]
			err = add('N', p, map[string]any{"MR-ｼﾒｲ": married.Kanji(), "MR-ｼﾒｲｶﾅ": married.Kana()})
		case chance(r, perDay(0.015)):
			err = add('K', p, map[string]any{"MR-ｳｹﾄﾘﾆﾝ": somebody(r).Kanji()})
		case chance(r, perDay(0.003)):
			switch r.IntN(3) {
			case 0:
				p.No = "8888888888"
				err = add('A', p, map[string]any{"MR-ｼﾞｭｳｼｮ": addresses[0]})
			case 1:
				err = add('Z', p, nil)
			default:
				err = add('A', p, nil)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return recs, nil
}

// counterPayment is what a policyholder pays at the counter to reinstate
// a lapsed policy: the arrears, premiums from the month after the paid-to
// month to this one, rounded up to the thousand as often as not, and
// short of them once in a while.
func (g *Generator) counterPayment(r *rand.Rand, p Policy, date Date) int64 {
	months := int((date.YearMonth()/100-p.PaidTo/100)*12 + date.YearMonth()%100 - p.PaidTo%100)
	switch p.Mode {
	case 2:
		months = (months + 11) / 12
	case 3:
		months = (months + 5) / 6
	}
	arrears := int64(months) * p.Premium
	switch r.IntN(10) {
	case 0:
		return arrears - 1000
	case 1, 2, 3, 4:
		return (arrears + 999) / 1000 * 1000
	}
	return arrears
}

// loanLimit is the share of a policy's cash value the insurer lends
// against, in percent.
const loanLimit = 80

// Loans makes up the branches' loan requests of a business day: a
// policyholder borrowing against a policy in force, in ten thousands of
// yen within what the limit leaves, a borrower repaying the balance as
// the last unload had it or a part of it, and a request a branch should
// not have sent now and then: more than the limit leaves, more than is
// owed, for no policy or of no kind.
func (g *Generator) Loans(date Date, policies []Policy) ([][]byte, error) {
	r := g.rng("loans", date.Number())
	var recs [][]byte
	n := 0
	add := func(kind byte, p Policy, amount int64) error {
		n++
		rec, err := made(g.Layouts.Loan, map[string]any{
			"LR-ｲﾗｲﾊﾞﾝｺﾞｳ": requestNumber('L', date, n), "LR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ": p.No, "LR-ｼｭﾍﾞﾂ": string(kind),
			"LR-ｷﾝｶﾞｸ": amount, "LR-ｼﾃﾝ": branchOf(p.Agent),
		})
		recs = append(recs, rec)
		return err
	}
	for _, p := range policies {
		if p.Status > 3 {
			continue
		}
		room := (p.cashValue()*loanLimit/100 - p.Loan) / 10_000
		var err error
		switch {
		case p.Status == 1 && room >= 10 && chance(r, perDay(0.05)):
			err = add('L', p, (1+r.Int64N(room*6/10))*10_000)
		case p.Covered() && p.Loan > 0 && chance(r, perDay(0.4)):
			err = add('P', p, repayment(r, p.Loan))
		case chance(r, perDay(0.003)):
			err = g.wrongLoan(r, add, p)
		}
		if err != nil {
			return nil, err
		}
	}
	return recs, nil
}

// repayment is what a borrower pays back: the balance, or a half, a
// third or a quarter of it in thousands of yen.
func repayment(r *rand.Rand, balance int64) int64 {
	if part := balance / int64(1+r.IntN(4)) / 1000 * 1000; part > 0 {
		return part
	}
	return balance
}

// wrongLoan makes up a request the program rejects.
func (g *Generator) wrongLoan(r *rand.Rand, add func(byte, Policy, int64) error, p Policy) error {
	switch r.IntN(4) {
	case 0:
		p.No = "7777777777"
		return add('L', p, 100_000)
	case 1:
		return add('X', p, 100_000)
	case 2:
		return add('P', p, p.Loan+1_000_000)
	}
	return add('L', p, p.Sum)
}
