package generate

import (
	"math"
	"math/rand/v2"
)

// newBusinessRate is how many policies the sales system sells a business
// day, for each thousand the master holds.
const newBusinessRate = 0.8

// NewBusiness makes up the sales system's file of a business day: the
// applications that became policies, with the riders they were sold
// with, dated in the Western calendar as the sales system writes them, a
// mistaken one now and then.
func (g *Generator) NewBusiness(date Date, policies []Policy) ([][]byte, error) {
	r := g.rng("new business", date.Number())
	count := poisson(r, newBusinessRate*float64(len(policies))/1000)
	var recs [][]byte
	for i := range count {
		p := g.application(r, date, i+1)
		fields := map[string]any{
			"NP-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ": p.No, "NP-ｼｮｳﾋﾝ": p.Product, "NP-ﾎｹﾝｷﾝｶﾞｸ": p.Sum, "NP-ﾊﾗｲｺﾐﾎｳﾎｳ": string(rune('0' + p.Mode)),
			"NP-ｹｲﾔｸﾋﾞ": p.Contract.Number(), "NP-ｼﾒｲ": p.Name, "NP-ｼﾒｲｶﾅ": p.Kana, "NP-ｾｲﾈﾝｶﾞｯﾋﾟ": p.Birth.Number(),
			"NP-ｾｲﾍﾞﾂ": p.Sex, "NP-ｼﾞｭｳｼｮ": p.Address, "NP-ｳｹﾄﾘﾆﾝ": p.Beneficiary,
			"NP-ｷﾞﾝｺｳｺｰﾄﾞ": p.Account.bank, "NP-ｼﾃﾝｺｰﾄﾞ": p.Account.branch, "NP-ﾖｷﾝｼｭﾓｸ": p.Account.kind,
			"NP-ｺｳｻﾞﾊﾞﾝｺﾞｳ": p.Account.number, "NP-ｺｳｻﾞﾒｲｷﾞ": p.Account.holder, "NP-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ": p.Agent,
		}
		if r.IntN(100) < 5 {
			name, value := mistake(r, p, policies)
			fields[name] = value
		}
		rec, err := made(g.Layouts.NewPolicy, fields)
		if err != nil {
			return nil, err
		}
		riders := p.Riders()
		if r.IntN(100) < 2 {
			riders = []Coverage{{Product: wrongRider(r, p), Amount: 1_000_000}}
		}
		for i, c := range riders {
			if err := g.setRider(rec, i, c); err != nil {
				return nil, err
			}
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

// setRider writes the n-th rider of an application, from 0. The record
// library addresses the first occurrence of a repeated field, so a rider
// is written there in a record of its own and copied to its place.
func (g *Generator) setRider(rec []byte, n int, c Coverage) error {
	l := g.Layouts.NewPolicy
	first, err := made(l, map[string]any{"NP-ﾄｸﾔｸ-ｼｮｳﾋﾝ": c.Product, "NP-ﾄｸﾔｸ-ｷﾝｶﾞｸ": c.Amount})
	if err != nil {
		return err
	}
	start, length := -1, 0
	for _, f := range l.Fields() {
		if f.Name == "NP-ﾄｸﾔｸ-ｼｮｳﾋﾝ" || f.Name == "NP-ﾄｸﾔｸ-ｷﾝｶﾞｸ" {
			if start < 0 {
				start = f.Start - 1
			}
			length += f.Length
		}
	}
	copy(rec[start+n*length:], first[start:start+length])
	return nil
}

// wrongRider is a rider the program refuses on an application: a product
// the insurer does not sell, one it sells on its own, or a rider of the
// main contract's kind.
func wrongRider(r *rand.Rand, p Policy) string {
	switch r.IntN(3) {
	case 0:
		return "ZZ"
	case 1:
		return "WL"
	}
	for _, pr := range products {
		if pr.rider && pr.kind == productOf(p.Product).kind {
			return pr.code
		}
	}
	return "EN"
}

// application makes up a policy as the sales system sends it: numbered
// by the day, contracted within the last few days, sold by an office
// still open, with the riders it was sold with.
func (g *Generator) application(r *rand.Rand, date Date, n int) Policy {
	pr := aProduct(r)
	holder := somebody(r)
	p := Policy{
		No: policyNumber(date, n), Product: pr.code, Mode: 1, Contract: date.AddDays(-r.IntN(4)),
		Name: holder.Kanji(), Kana: holder.Kana(), Sex: holder.sex,
		Address: addresses[r.IntN(len(addresses))], Beneficiary: relative(r, holder).Kanji(),
		Account: anAccount(r, holder), Agent: agents[r.IntN(len(agents))].code,
	}
	p.Mode = aMode(r)
	p.Birth, p.EntryAge = g.birth(r, pr, p.Contract)
	p.Sum = g.sum(r, pr)
	p.Coverages = append([]Coverage{{No: 1, Product: pr.code, Status: 1, Amount: p.Sum}}, g.riders(r, p, pr, p.Contract)...)
	return p
}

// mistake spoils one field of an application the way the sales system's
// data is spoiled now and then: a product the insurer does not sell, a
// policy already on the master, an age or a sum outside the product, a
// year mistyped without its century, or a sex that is neither.
func mistake(r *rand.Rand, p Policy, policies []Policy) (string, any) {
	switch r.IntN(6) {
	case 0:
		return "NP-ｼｮｳﾋﾝ", "XX"
	case 1:
		if len(policies) > 0 {
			return "NP-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ", policies[r.IntN(len(policies))].No
		}
	case 2:
		return "NP-ｾｲﾈﾝｶﾞｯﾋﾟ", p.Contract.AddYears(-85).Number()
	case 3:
		return "NP-ﾎｹﾝｷﾝｶﾞｸ", int64(999_000_000)
	case 4:
		return "NP-ｹｲﾔｸﾋﾞ", p.Contract.Number() - 20000000
	}
	return "NP-ｾｲﾍﾞﾂ", 3
}

// poisson draws how many of something happen when the expected number
// is lambda.
func poisson(r *rand.Rand, lambda float64) int {
	limit, n, p := math.Exp(-lambda), 0, 1.0
	for {
		p *= r.Float64()
		if p < limit {
			return n
		}
		n++
	}
}
