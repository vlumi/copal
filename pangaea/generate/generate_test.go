package generate

import (
	"slices"
	"testing"
	"time"

	"github.com/vlumi/copal/pangaea/layouts"
)

func generator(t *testing.T) *Generator {
	t.Helper()
	l, err := layouts.Load("../copybooks")
	if err != nil {
		t.Fatal(err)
	}
	return &Generator{Layouts: l, Seed: 42}
}

func TestCalendar(t *testing.T) {
	tests := []struct {
		date          Date
		business      bool
		firstBusiness bool
		payday        bool
	}{
		{NewDate(2026, time.September, 1), true, true, false},  // Tuesday
		{NewDate(2026, time.September, 25), true, false, true}, // Friday the 25th
		{NewDate(2026, time.October, 23), true, false, true},   // Friday before Sunday the 25th
		{NewDate(2026, time.October, 25), false, false, false}, // Sunday
		{NewDate(2026, time.November, 2), true, true, false},   // Monday after Sunday the 1st
		{NewDate(2026, time.November, 1), false, false, false}, // Sunday the 1st
		{NewDate(2026, time.December, 24), true, false, false}, // Thursday before Friday the 25th
	}
	for _, test := range tests {
		d := test.date
		if d.BusinessDay() != test.business || d.FirstBusinessDayOfMonth() != test.firstBusiness || d.Payday() != test.payday {
			t.Errorf("%s: business %v, first %v, payday %v; want %v, %v, %v", d,
				d.BusinessDay(), d.FirstBusinessDayOfMonth(), d.Payday(), test.business, test.firstBusiness, test.payday)
		}
	}
	if d, err := DateFromNumber(20260230); err == nil {
		t.Errorf("20260230 read as %s", d)
	}
}

func TestMaster(t *testing.T) {
	g := generator(t)
	records, err := g.Master(500, NewDate(2026, time.September, 1))
	if err != nil {
		t.Fatal(err)
	}
	again, _ := g.Master(500, NewDate(2026, time.September, 1))
	if !slices.EqualFunc(records, again, slices.Equal) {
		t.Error("two masters from the same seed differ")
	}
	accounts, err := g.ReadAccounts(records)
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) < 500 || len(accounts) > 501 {
		t.Errorf("%d accounts, want about 500", len(accounts))
	}
	statuses := map[string]int{}
	for i, a := range accounts {
		if i > 0 && a.No <= accounts[i-1].No {
			t.Fatalf("account %d after %d: not in order", a.No, accounts[i-1].No)
		}
		statuses[a.Status]++
	}
	if statuses["O"] < 450 || statuses["F"] == 0 || statuses["X"] == 0 {
		t.Errorf("statuses %v, want mostly open, a few frozen and closed", statuses)
	}
}

func TestDayHasMatchingControlTotals(t *testing.T) {
	g := generator(t)
	date := NewDate(2026, time.September, 25)
	master, err := g.Master(500, date)
	if err != nil {
		t.Fatal(err)
	}
	accounts, err := g.ReadAccounts(master)
	if err != nil {
		t.Fatal(err)
	}
	file, err := g.Day(date, 19, accounts)
	if err != nil {
		t.Fatal(err)
	}
	tran := g.Layouts.Tran
	var count, total, hash int64
	for _, r := range file[1 : len(file)-1] {
		amount, _ := tran.Number(r, "TD-AMOUNT")
		account, _ := tran.Number(r, "TD-ACCT-NO")
		count, total, hash = count+1, total+amount, hash+account
	}
	trailer := file[len(file)-1]
	for name, want := range map[string]int64{"TT-DETAIL-COUNT": count, "TT-AMOUNT-TOTAL": total, "TT-ACCT-HASH": hash} {
		if got, err := tran.Number(trailer, name); err != nil || got != want {
			t.Errorf("%s %d, %v; the details give %d", name, got, err, want)
		}
	}
	if seq, _ := tran.Number(file[0], "TH-FILE-SEQ"); seq != 19 {
		t.Errorf("file sequence %d, want 19", seq)
	}
	if count < 300 {
		t.Errorf("%d transactions on a payday for 500 accounts, want salaries for most", count)
	}
}

func TestRatesAndAccruals(t *testing.T) {
	g := generator(t)
	rates, err := g.Rates()
	if err != nil {
		t.Fatal(err)
	}
	if len(rates) != len(rateSheet) {
		t.Fatalf("%d rates, want %d", len(rates), len(rateSheet))
	}
	if floor, err := g.Layouts.Rate.Number(rates[1], "RT-FLOOR"); err != nil || floor != 1000000 {
		t.Errorf("second floor %d, %v, want 10,000.00", floor, err)
	}
	accruals, err := g.Accruals([]Account{{No: 1000000001}, {No: 1000000002}}, NewDate(2026, 9, 1))
	if err != nil {
		t.Fatal(err)
	}
	if through, err := g.Layouts.Accrual.Number(accruals[1], "AC-THROUGH-DATE"); err != nil || through != 20260901 {
		t.Errorf("through %d, %v", through, err)
	}
	if accrued, err := g.Layouts.Accrual.Number(accruals[1], "AC-ACCRUED"); err != nil || accrued != 0 {
		t.Errorf("accrued %d, %v", accrued, err)
	}
}

func TestMaintenanceIsInAccountOrder(t *testing.T) {
	g := generator(t)
	date := NewDate(2026, 9, 1)
	master, err := g.Master(500, date)
	if err != nil {
		t.Fatal(err)
	}
	accounts, err := g.ReadAccounts(master)
	if err != nil {
		t.Fatal(err)
	}
	var opened int
	for day := range 30 {
		requests, err := g.Maintenance(date.AddDays(day), accounts)
		if err != nil {
			t.Fatal(err)
		}
		var previous int64
		for _, r := range requests {
			no, err := g.Layouts.Maint.Number(r, "MR-ACCT-NO")
			if err != nil {
				t.Fatal(err)
			}
			if no < previous {
				t.Fatalf("%s: account %d after %d", date.AddDays(day), no, previous)
			}
			previous = no
			if kind, _ := g.Layouts.Maint.Text(r, "MR-REQUEST-TYPE"); kind == "O" && no > accounts[len(accounts)-1].No {
				opened++
			}
		}
	}
	if opened == 0 {
		t.Error("no accounts opened in a month")
	}
}

func TestCards(t *testing.T) {
	g := generator(t)
	date := NewDate(2026, time.September, 1)
	master, err := g.Master(500, date)
	if err != nil {
		t.Fatal(err)
	}
	accounts, err := g.ReadAccounts(master)
	if err != nil {
		t.Fatal(err)
	}
	cards, err := g.Cards(accounts, date)
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) < 200 {
		t.Fatalf("%d cards for %d accounts", len(cards), len(accounts))
	}
	previous := int64(0)
	for _, c := range cards {
		no, err := g.Layouts.Card.Number(c, "CR-CARD-NO")
		if err != nil {
			t.Fatal(err)
		}
		if no <= previous || no/1_000_000_000_000 != cardPrefix {
			t.Errorf("card %d after %d", no, previous)
		}
		previous = no
	}
	if got := cardNumber(1); got != 9100000000000018 {
		t.Errorf("card number of serial 1 is %d, want 9100000000000018, which passes the Luhn check", got)
	}
}
