package generate

import "math/rand/v2"

// product is a product the insurer sells, as the product table has it:
// its kind, the ages it is sold at, its term, the sums it insures, the
// daily hospitalization benefit per 1,000 yen and the waiting days of a
// medical policy, and the dividend rate of a participating one.
type product struct {
	code       string
	name       string
	kind       int // 1 whole life, 2 renewable term, 3 endowment, 4 medical
	minAge     int
	maxAge     int
	term       int
	minSum     int64
	maxSum     int64
	dailyRate  int // thousandths per 1,000 yen
	waiting    int
	dividend   int // ten-thousandths of the sum, a year
	weight     int // how much of the business it is
	baseRate   int // the monthly premium per 1,000 yen, in hundredths, at the base age
	baseAge    int
	ratePerAge int // how much the rate grows a year of entry age, in hundredths
}

var products = []product{
	{"WL", sjis("終身保険"), 1, 15, 70, 0, 1_000_000, 100_000_000, 0, 0, 20, 45, 260, 20, 7},
	{"TM", sjis("定期保険"), 2, 20, 60, 10, 1_000_000, 300_000_000, 0, 0, 0, 20, 80, 20, 3},
	{"EN", sjis("養老保険"), 3, 0, 55, 20, 500_000, 50_000_000, 0, 0, 15, 20, 380, 0, 2},
	{"MD", sjis("医療保険"), 4, 15, 75, 0, 1_000_000, 20_000_000, 1, 90, 0, 15, 90, 20, 2},
}

// aProduct draws a product by its share of the business.
func aProduct(r *rand.Rand) product {
	weights := make([]int, len(products))
	for i, p := range products {
		weights[i] = p.weight
	}
	return products[pick(r, weights)]
}

// rate is the monthly premium per 1,000 yen of the sum, in hundredths of
// a yen, for a product, a sex and an entry age: a line for each age, a
// woman's nine tenths of a man's, and nothing outside the product's ages.
func (p product) rate(sex, age int) int {
	if age < p.minAge || age > p.maxAge {
		return 0
	}
	r := p.baseRate + p.ratePerAge*(age-p.baseAge)
	if sex == 2 {
		r = r * 9 / 10
	}
	return r
}

// cashValue is the surrender value per 1,000 yen of the sum, in
// thousandths of a yen, by product, entry age band and policy year: a
// whole life policy's grows from its second year, more for an older
// entrant, an endowment's reaches the sum at maturity, and a term or
// medical policy has none.
func (p product) cashValue(band, year int) int {
	switch p.kind {
	case 1:
		if year < 2 {
			return 0
		}
		return (21 + 5*(band-1) + 8*(year-2)) * 1000
	case 3:
		return min(48*year, 1000) * 1000
	}
	return 0
}

// reserve is the policy reserve per 1,000 yen of the sum at the year's
// end, in thousandths of a yen, by product and policy year.
func (p product) reserve(year int) int {
	switch p.kind {
	case 1:
		return 9 * year * 1000
	case 2:
		return min(year, 10) * 1500
	case 3:
		return min(48*year, 1000) * 1000
	case 4:
		return 2 * year * 1000
	}
	return 0
}

// ageBand is the entry age band of the cash value table.
func ageBand(entryAge int) int {
	switch {
	case entryAge < 20:
		return 1
	case entryAge < 60:
		return entryAge / 10
	}
	return 6
}

// monthlyPremium is the premium for a sum at a rate, truncated to the
// yen, as the programs compute it.
func monthlyPremium(sum int64, rate int) int64 {
	return sum / 1000 * int64(rate) / 100
}

// premium is the premium a policy pays at a time: the monthly premium,
// or twelve months less five percent for an annual payer.
func premium(sum int64, rate int, mode int) int64 {
	monthly := monthlyPremium(sum, rate)
	if mode == 2 {
		return monthly * 12 * 95 / 100
	}
	return monthly
}

// insuranceAge is the age the programs insure at: the age in full years,
// plus one when six months or more of the next have passed.
func insuranceAge(birth, at Date) int {
	months := MonthsBetween(birth, at)
	age := months / 12
	if months%12 >= 6 {
		age++
	}
	return age
}

// Products writes the product table.
func (g *Generator) Products() ([][]byte, error) {
	var recs [][]byte
	for _, p := range products {
		rec, err := made(g.Layouts.Product, map[string]any{
			"PR-ｺｰﾄﾞ": p.code, "PR-ﾒｲｼｮｳ": p.name, "PR-ｼｭﾙｲ": p.kind, "PR-ｻｲﾃｲﾈﾝﾚｲ": p.minAge, "PR-ｻｲｺｳﾈﾝﾚｲ": p.maxAge,
			"PR-ﾎｹﾝｷｶﾝ": p.term, "PR-ｻｲﾃｲﾎｹﾝｷﾝｶﾞｸ": p.minSum, "PR-ｻｲｺｳﾎｹﾝｷﾝｶﾞｸ": p.maxSum,
			"PR-ﾆｭｳｲﾝﾆﾁｶﾞｸﾘﾂ": p.dailyRate, "PR-ﾒﾝｾｷﾆｯｽｳ": p.waiting, "PR-ﾊｲﾄｳﾘﾂ": p.dividend,
		})
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

// Rates writes the premium rate table: a line for each product, sex and
// age from 0 to 80.
func (g *Generator) Rates() ([][]byte, error) {
	var recs [][]byte
	for _, p := range products {
		for sex := 1; sex <= 2; sex++ {
			for age := 0; age <= 80; age++ {
				rec, err := made(g.Layouts.Rate, map[string]any{"RT-ｼｮｳﾋﾝ": p.code, "RT-ｾｲﾍﾞﾂ": sex, "RT-ﾈﾝﾚｲ": age, "RT-ﾘｮｳﾘﾂ": p.rate(sex, age)})
				if err != nil {
					return nil, err
				}
				recs = append(recs, rec)
			}
		}
	}
	return recs, nil
}

// CashValues writes the surrender value table: a line for each product,
// entry age band and policy year to 40.
func (g *Generator) CashValues() ([][]byte, error) {
	var recs [][]byte
	for _, p := range products {
		for band := 1; band <= 6; band++ {
			for year := 1; year <= 40; year++ {
				rec, err := made(g.Layouts.CashValue, map[string]any{"CV-ｼｮｳﾋﾝ": p.code, "CV-ﾈﾝﾚｲｸﾌﾞﾝ": band, "CV-ｹｲｶﾈﾝｽｳ": year, "CV-ｹｲｽｳ": p.cashValue(band, year)})
				if err != nil {
					return nil, err
				}
				recs = append(recs, rec)
			}
		}
	}
	return recs, nil
}

// Reserves writes the reserve factor table: a line for each product and
// policy year to 50.
func (g *Generator) Reserves() ([][]byte, error) {
	var recs [][]byte
	for _, p := range products {
		for year := 1; year <= 50; year++ {
			rec, err := made(g.Layouts.Reserve, map[string]any{"RS-ｼｮｳﾋﾝ": p.code, "RS-ｹｲｶﾈﾝｽｳ": year, "RS-ｹｲｽｳ": p.reserve(year)})
			if err != nil {
				return nil, err
			}
			recs = append(recs, rec)
		}
	}
	return recs, nil
}

// Agents writes the agent master: every office but the one closed long
// ago.
func (g *Generator) Agents() ([][]byte, error) {
	var recs [][]byte
	for _, a := range agents {
		rec, err := made(g.Layouts.Agent, map[string]any{"AG-ｺｰﾄﾞ": a.code, "AG-ｼﾒｲ": a.name, "AG-ｼﾃﾝ": a.branch, "AG-ﾆｭｳｼｬﾈﾝｹﾞﾂ": a.joined, "AG-ｼｶｸ": a.grade})
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}
