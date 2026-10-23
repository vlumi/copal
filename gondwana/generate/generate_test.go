package generate

import (
	"slices"
	"testing"
	"time"

	"github.com/vlumi/copal/gondwana/records"
)

func generator(t *testing.T) *Generator {
	t.Helper()
	l, err := records.Load("../copybooks")
	if err != nil {
		t.Fatal(err)
	}
	return &Generator{Layouts: l, Seed: 42}
}

var opening = NewDate(2026, time.April, 1)

func TestExtract(t *testing.T) {
	g := generator(t)
	recs, err := g.Extract(500, opening)
	if err != nil {
		t.Fatal(err)
	}
	again, _ := g.Extract(500, opening)
	if !slices.EqualFunc(recs, again, slices.Equal) {
		t.Error("two extracts from the same seed differ")
	}
	policies, err := g.ReadPolicies(recs)
	if err != nil {
		t.Fatal(err)
	}
	if len(policies) != 500 {
		t.Fatalf("%d policies, want 500", len(policies))
	}
	eras, products, statuses, modes := map[byte]int{}, map[string]int{}, map[int]int{}, map[int]int{}
	numbers := map[string]bool{}
	for _, p := range policies {
		eras[p.Contract.Era()[0]]++
		products[p.Product]++
		statuses[p.Status]++
		modes[p.Mode]++
		if numbers[p.No] || len(p.No) != 10 {
			t.Errorf("policy %s: issued twice, or not a policy number", p.No)
		}
		numbers[p.No] = true
		if !p.Contract.Before(opening) || p.Year < 1 || p.Premium <= 0 || p.Sum <= 0 || p.EntryAge < 0 || p.PaidTo > opening.YearMonth() {
			t.Errorf("policy %s: contracted %s, year %d, sum %d, premium %d, entry age %d, paid to %d", p.No, p.Contract, p.Year, p.Sum, p.Premium, p.EntryAge, p.PaidTo)
		}
		if p.Status == 2 && (p.GraceStart.IsZero() || p.Arrears == 0) || p.Status == 3 && (p.LapseDate.IsZero() || p.Arrears == 0) || p.Covered() && p.Reserve == 0 {
			t.Errorf("policy %s in status %d: grace from %s, lapsed %s, %d arrears, reserve %d", p.No, p.Status, p.GraceStart, p.LapseDate, p.Arrears, p.Reserve)
		}
		if (p.Product == "EN" || p.Product == "TM") && p.Maturity.IsZero() || (p.Product == "WL" || p.Product == "MD") && !p.Maturity.IsZero() {
			t.Errorf("policy %s of %s matures %s", p.No, p.Product, p.Maturity)
		}
		if got := insuranceAge(p.Birth, p.Contract); got != p.EntryAge {
			t.Errorf("policy %s: entry age %d, but %d from the dates", p.No, p.EntryAge, got)
		}
	}
	if eras['S'] == 0 || eras['H'] == 0 || eras['R'] == 0 || len(products) != 4 || modes[2] == 0 {
		t.Errorf("eras %v, products %v, modes %v", eras, products, modes)
	}
	if statuses[1] < 350 || statuses[2] == 0 || statuses[3] == 0 || statuses[4] == 0 || statuses[5] == 0 || statuses[6] == 0 {
		t.Errorf("statuses %v", statuses)
	}
}

func TestTables(t *testing.T) {
	g := generator(t)
	for _, c := range []struct {
		name  string
		table func() ([][]byte, error)
		lines int
	}{
		{"products", g.Products, 4}, {"rates", g.Rates, 648}, {"cash values", g.CashValues, 960}, {"reserves", g.Reserves, 200}, {"agents", g.Agents, 10},
	} {
		recs, err := c.table()
		if err != nil {
			t.Fatal(err)
		}
		if len(recs) != c.lines {
			t.Errorf("%d %s, want %d", len(recs), c.name, c.lines)
		}
	}
	if r := products[0].rate(1, 40); r != 400 {
		t.Errorf("a man of 40 pays %d hundredths per 1,000 yen of whole life, want 400", r)
	}
	if r := products[2].rate(2, 3); r != 347 {
		t.Errorf("a girl of 3 pays %d hundredths per 1,000 yen of endowment, want 347", r)
	}
	if p := premium(10_000_000, 379, 2); p != 432_060 {
		t.Errorf("an annual premium of %d, want 432,060", p)
	}
	if age := insuranceAge(NewDate(1990, time.July, 15), NewDate(2027, time.March, 1)); age != 37 {
		t.Errorf("insurance age %d, want 37: 36 years and 7 months", age)
	}
}

func TestDays(t *testing.T) {
	g := generator(t)
	extract, err := g.Extract(2000, opening)
	if err != nil {
		t.Fatal(err)
	}
	policies, err := g.ReadPolicies(extract)
	if err != nil {
		t.Fatal(err)
	}
	newBusiness, claims, changes := 0, 0, 0
	for date := opening; date.Before(opening.AddDays(60)); date = date.NextBusinessDay() {
		n, err := g.NewBusiness(date, policies)
		if err != nil {
			t.Fatal(err)
		}
		c, err := g.Claims(date, policies)
		if err != nil {
			t.Fatal(err)
		}
		m, err := g.Changes(date, policies)
		if err != nil {
			t.Fatal(err)
		}
		newBusiness, claims, changes = newBusiness+len(n), claims+len(c), changes+len(m)
		for _, rec := range append(append(n, c...), m...) {
			if len(rec) == 0 {
				t.Errorf("%s: an empty record", date)
			}
		}
	}
	if newBusiness < 20 || claims < 10 || changes < 10 {
		t.Errorf("in two months: %d new policies, %d benefit requests, %d changes", newBusiness, claims, changes)
	}
}

func TestResults(t *testing.T) {
	g := generator(t)
	l := g.Layouts.Transfer
	var request [][]byte
	header, _ := made(l, map[string]any{"ZG-ﾃﾞｰﾀｸﾌﾞﾝ": 1, "ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ": 91, "ZG-H-ｷﾞﾝｺｳｺｰﾄﾞ": 1, "ZG-H-ｼﾃﾝｺｰﾄﾞ": 1})
	request = append(request, header)
	var total int64
	for i := range 200 {
		data, _ := made(l, map[string]any{"ZG-ﾃﾞｰﾀｸﾌﾞﾝ": 2, "ZG-D-ｷﾞﾝｺｳｺｰﾄﾞ": 5, "ZG-D-ｼﾃﾝｺｰﾄﾞ": 12, "ZG-D-ｷﾝｶﾞｸ": int64(1000 * (i + 1)), "ZG-D-ｺｷｬｸﾊﾞﾝｺﾞｳ": policyNumber(opening, i)})
		request = append(request, data)
		total += int64(1000 * (i + 1))
	}
	trailer, _ := made(l, map[string]any{"ZG-ﾃﾞｰﾀｸﾌﾞﾝ": 8, "ZG-T-ｺﾞｳｹｲｹﾝｽｳ": 200, "ZG-T-ｺﾞｳｹｲｷﾝｶﾞｸ": total})
	end, _ := made(l, map[string]any{"ZG-ﾃﾞｰﾀｸﾌﾞﾝ": 9})
	request = append(request, trailer, end)
	results, err := g.Results(opening, request)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != len(request) {
		t.Fatalf("%d results for %d records", len(results), len(request))
	}
	var paid, unpaid, paidAmount, unpaidAmount int64
	for _, rec := range results[1:201] {
		result, _ := l.Text(rec, "ZG-D-ﾌﾘｶｴｹｯｶ")
		amount, _ := l.Number(rec, "ZG-D-ｷﾝｶﾞｸ")
		if bank, _ := l.Text(rec, "ZG-D-ｷﾞﾝｺｳﾒｲ"); bank != banks[1].name {
			t.Errorf("the agent named bank 5 %q", bank)
		}
		if result == "0" {
			paid, paidAmount = paid+1, paidAmount+amount
		} else {
			unpaid, unpaidAmount = unpaid+1, unpaidAmount+amount
		}
	}
	gotPaid, _ := l.Number(results[201], "ZG-T-ﾌﾘｶｴｽﾞﾐｹﾝｽｳ")
	gotPaidAmount, _ := l.Number(results[201], "ZG-T-ﾌﾘｶｴｽﾞﾐｷﾝｶﾞｸ")
	gotUnpaid, _ := l.Number(results[201], "ZG-T-ﾌﾘｶｴﾌﾉｳｹﾝｽｳ")
	gotUnpaidAmount, _ := l.Number(results[201], "ZG-T-ﾌﾘｶｴﾌﾉｳｷﾝｶﾞｸ")
	if gotPaid != paid || gotPaidAmount != paidAmount || gotUnpaid != unpaid || gotUnpaidAmount != unpaidAmount || unpaid == 0 || paid < 180 {
		t.Errorf("trailer: %d paid %d, %d unpaid %d; counted %d paid %d, %d unpaid %d", gotPaid, gotPaidAmount, gotUnpaid, gotUnpaidAmount, paid, paidAmount, unpaid, unpaidAmount)
	}
}
