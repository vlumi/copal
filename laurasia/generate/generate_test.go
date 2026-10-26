package generate

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/vlumi/cobolith/record"
	"github.com/vlumi/copal/laurasia/records"
)

func generator(t *testing.T) *Generator {
	t.Helper()
	l, err := records.Load("../copybooks")
	if err != nil {
		t.Fatal(err)
	}
	return &Generator{Layouts: l, Seed: 42}
}

func TestCalendar(t *testing.T) {
	monday := NewDate(2026, time.October, 5)
	if monday.Weekday() != time.Monday || !monday.BusinessDay() || monday.AddDays(-1).BusinessDay() {
		t.Errorf("%s: weekday %s, business %v, Sunday business %v", monday, monday.Weekday(), monday.BusinessDay(), monday.AddDays(-1).BusinessDay())
	}
	if friday := monday.PreviousBusinessDay(); friday.String() != "2026-10-02" || friday.Stamp() != "20261002" || friday.Number() != 20261002 {
		t.Errorf("the business day before %s is %s, stamp %s, number %d", monday, friday, friday.Stamp(), friday.Number())
	}
	if d, err := DateFromNumber(20261001); err != nil || d.Day() != 1 || !d.After(NewDate(2026, time.September, 30)) {
		t.Errorf("20261001 read as %s, %v", d, err)
	}
	if d, err := DateFromNumber(20260230); err == nil {
		t.Errorf("20260230 read as %s", d)
	}
}

func TestCards(t *testing.T) {
	g := generator(t)
	recs, err := g.Cards(500)
	if err != nil {
		t.Fatal(err)
	}
	again, _ := g.Cards(500)
	if !slices.EqualFunc(recs, again, slices.Equal) {
		t.Error("two portfolios from the same seed differ")
	}
	cards, err := g.ReadCards(recs)
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 500 {
		t.Fatalf("%d cards, want 500", len(cards))
	}
	products, statuses := map[string]int{}, map[string]int{}
	numbers := map[string]bool{}
	for _, c := range cards {
		products[c.Product]++
		statuses[c.Status]++
		if numbers[c.No] || len(c.No) != 16 || luhn(c.No[:15]) != c.No {
			t.Errorf("card %s: issued twice, or not a card number", c.No)
		}
		numbers[c.No] = true
		if c.Limit <= 0 || c.Balance < 0 || c.Balance > c.Limit || c.Statement > c.Balance || c.CycleDay < 1 || c.CycleDay > 28 {
			t.Errorf("card %s: limit %d, balance %d, statement %d, cycle day %d", c.No, c.Limit, c.Balance, c.Statement, c.CycleDay)
		}
	}
	if len(products) != 3 || statuses["A"] < 450 || statuses["B"] == 0 || statuses["C"] == 0 {
		t.Errorf("products %v, statuses %v", products, statuses)
	}
}

func TestFees(t *testing.T) {
	g := generator(t)
	recs, err := g.Fees()
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 27 {
		t.Errorf("%d fees, want 27: thirty slots less the three the schedule has never had", len(recs))
	}
	rate, _ := g.Layouts.Fee.Number(recs[0], "FD-RATE")
	if product, _ := g.Layouts.Fee.Text(recs[0], "FD-PRODUCT"); product != "CL" || rate != 110 {
		t.Errorf("first fee: %s at %d, want CL at 110", product, rate)
	}
}

// A night's log and clearing file: the log holds every day since the last
// night, the clearing file adds up batch by batch and in its trailer, and
// what the networks present is not quite what the switch approved.
func TestNight(t *testing.T) {
	g := generator(t)
	recs, err := g.Cards(500)
	if err != nil {
		t.Fatal(err)
	}
	cards, _ := g.ReadCards(recs)
	monday, friday := NewDate(2026, time.October, 5), NewDate(2026, time.October, 2)
	n, err := g.Night(monday, friday, cards, g.Merchants())
	if err != nil {
		t.Fatal(err)
	}
	again, _ := g.Night(monday, friday, cards, g.Merchants())
	if !slices.EqualFunc(n.Clearing, again.Clearing, slices.Equal) || !slices.EqualFunc(n.Log, again.Log, slices.Equal) {
		t.Error("two nights from the same seed differ")
	}

	days, decisions, standIns := map[int64]int{}, map[string]int{}, 0
	codes := map[string]bool{}
	for _, r := range n.Log {
		day, _ := g.Layouts.Auth.Number(r, "AL-DATE")
		decision, _ := g.Layouts.Auth.Text(r, "AL-DECISION")
		standIn, _ := g.Layouts.Auth.Text(r, "AL-STAND-IN")
		auth, _ := g.Layouts.Auth.Text(r, "AL-AUTH-CODE")
		days[day]++
		decisions[decision]++
		if standIn == "Y" {
			standIns++
		}
		if codes[auth] {
			t.Errorf("authorization code %s given twice", auth)
		}
		codes[auth] = true
	}
	if len(days) != 3 || days[20261003] == 0 || days[20261004] == 0 || days[20261005] == 0 {
		t.Errorf("the log covers %v, want Saturday to Monday", days)
	}
	if decisions["A"] < 500 || decisions["D"] == 0 || standIns == 0 {
		t.Errorf("decisions %v, %d in stand-in", decisions, standIns)
	}

	f := clearingFile{t: t, l: g.Layouts, network: "LCN1"}
	f.read(n.Clearing)
	ours := clearingFile{t: t, l: g.Layouts, network: "ONUS", own: true}
	ours.read(n.OnUs)
	if f.date != monday.Number() || ours.date != monday.Number() || f.presentments+ours.presentments != n.Presentments || f.batches < 4 || ours.batches < 1 {
		t.Errorf("files of %d and %d with %d and %d presentments in %d and %d batches, want %s, %d in all, 4 or more and 1 or more",
			f.date, ours.date, f.presentments, ours.presentments, f.batches, ours.batches, monday, n.Presentments)
	}
	for _, file := range []clearingFile{f, ours} {
		if file.trailerBatches != int64(file.batches) || file.trailerRecords != int64(file.presentments) || file.trailerAmount != file.amount {
			t.Errorf("file trailer says %d batches, %d records, %d; the file has %d, %d, %d", file.trailerBatches, file.trailerRecords, file.trailerAmount, file.batches, file.presentments, file.amount)
		}
	}
	if f.badBatches+ours.badBatches > 2 {
		t.Errorf("%d batches with a bad trailer, want one a file at most", f.badBatches+ours.badBatches)
	}
	for auth := range ours.auths {
		f.auths[auth] = true
	}
	f.presentments += ours.presentments
	if n.Presentments < decisions["A"]*8/10 || n.Presentments > decisions["A"] {
		t.Errorf("%d presentments of %d approvals", n.Presentments, decisions["A"])
	}
	unauthorized := 0
	for auth := range f.auths {
		if !codes[auth] {
			unauthorized++
		}
		delete(codes, auth)
	}
	if unauthorized == 0 || unauthorized > f.presentments/10 || len(codes) == 0 {
		t.Errorf("%d of %d presentments match no authorization, and %d authorizations lapsed", unauthorized, f.presentments, len(codes))
	}
}

// clearingFile reads a clearing file back, checking each batch against
// its trailer.
type clearingFile struct {
	t                                 *testing.T
	l                                 *records.Layouts
	date                              int64
	batches, presentments, badBatches int
	amount                            int64
	trailerBatches, trailerRecords    int64
	trailerAmount                     int64
	batchCount, batchAmount           int64
	auths                             map[string]bool
	// network is what the file header names, and own whether the file is
	// our own acquiring's, every batch and merchant ours, or the network's,
	// none of them.
	network string
	own     bool
}

func (f *clearingFile) read(file [][]byte) {
	f.auths = map[string]bool{}
	for _, r := range file {
		switch string(r[:2]) {
		case "FH":
			f.date = f.number(f.l.FileHeader, r, "FH-FILE-DATE")
			if network, _ := f.l.FileHeader.Text(r, "FH-NETWORK"); network != f.network {
				f.t.Errorf("a file from %s, want %s", network, f.network)
			}
		case "BH":
			f.batches++
			f.batchCount, f.batchAmount = 0, 0
			if acquirer, _ := f.l.BatchHeader.Text(r, "BH-ACQUIRER"); (acquirer == OwnAcquirer) != f.own {
				f.t.Errorf("a batch from %s in the %s file", acquirer, f.network)
			}
		case "PR":
			f.presentments++
			f.batchCount++
			amount := f.number(f.l.Presentment, r, "PR-AMOUNT")
			f.batchAmount += amount
			f.amount += amount
			auth, _ := f.l.Presentment.Text(r, "PR-AUTH-CODE")
			f.auths[auth] = true
			if merchant, _ := f.l.Presentment.Text(r, "PR-MERCHANT"); strings.HasPrefix(merchant, ownMerchants) != f.own {
				f.t.Errorf("merchant %s in the %s file", merchant, f.network)
			}
		case "BT":
			if f.batchCount > maxBatch {
				f.t.Errorf("a batch of %d presentments", f.batchCount)
			}
			if f.number(f.l.BatchTrailer, r, "BT-COUNT") != f.batchCount || f.number(f.l.BatchTrailer, r, "BT-AMOUNT") != f.batchAmount {
				f.badBatches++
			}
		case "FT":
			f.trailerBatches = f.number(f.l.FileTrailer, r, "FT-BATCH-COUNT")
			f.trailerRecords = f.number(f.l.FileTrailer, r, "FT-RECORD-COUNT")
			f.trailerAmount = f.number(f.l.FileTrailer, r, "FT-AMOUNT")
		default:
			f.t.Errorf("a record of type %q", r[:2])
		}
	}
}

func (f *clearingFile) number(l *record.Layout, r []byte, name string) int64 {
	n, err := l.Number(r, name)
	if err != nil {
		f.t.Fatal(err)
	}
	return n
}

// A day's dispute events move the cases in flight on, leave the settled
// ones alone, and open new cases numbered after the last.
func TestDisputes(t *testing.T) {
	g := generator(t)
	cases := []Case{
		{No: 3, Status: "WN", Merchant: "MER0001000", Amount: 1000},
		{No: 7, Status: "LS", Merchant: "MER0001001", Amount: 2000},
	}
	for i := range 40 {
		cases = append(cases, Case{No: int64(10 + i), Status: []string{"OP", "CB", "RP", "PA"}[i%4], Merchant: "MER0001002", Amount: 500})
	}
	var sales []Sale
	for range 2000 {
		sales = append(sales, Sale{Card: "4000010000000000", Merchant: "MER0001003", Amount: 4242})
	}
	date := NewDate(2026, time.October, 6)
	events, err := g.Disputes(date, cases, sales)
	if err != nil {
		t.Fatal(err)
	}
	again, _ := g.Disputes(date, cases, sales)
	if !slices.EqualFunc(events, again, slices.Equal) {
		t.Error("two days of events from the same seed differ")
	}
	moved, opened := 0, 0
	var next int64 = 50
	for _, r := range events {
		no, _ := g.Layouts.Event.Number(r, "DE-CASE-NO")
		event, _ := g.Layouts.Event.Text(r, "DE-EVENT")
		amount, _ := g.Layouts.Event.Number(r, "DE-AMOUNT")
		switch {
		case no == 3 || no == 7:
			t.Errorf("case %d is settled and got %s", no, event)
		case no < 50:
			moved++
		case event == "OP":
			opened++
			if no != next || amount != 4242 {
				t.Errorf("case %d opened for %d, want case %d for 4242", no, amount, next)
			}
			next++
		}
	}
	if moved < 5 || moved > 40 || opened < 2 || opened > 30 {
		t.Errorf("%d cases moved on, %d opened", moved, opened)
	}
}

// A business day's feed from the card system: each cardholder pays on
// their day and not on others, a paid statement is left alone, and a
// status change comes now and then.
func TestPayments(t *testing.T) {
	g := generator(t)
	recs, err := g.Cards(500)
	if err != nil {
		t.Fatal(err)
	}
	cards, _ := g.ReadCards(recs)
	owing := map[string]Card{}
	for _, c := range cards {
		if c.Statement > 0 && !c.PaidInFull {
			owing[c.No] = c
		}
	}
	paid := map[string]int{}
	kinds := map[string]int{}
	for date := NewDate(2026, time.October, 1); date.Number() <= 20261031; date = date.AddDays(1) {
		if !date.BusinessDay() {
			continue
		}
		feed, err := g.Payments(date, cards)
		if err != nil {
			t.Fatal(err)
		}
		again, _ := g.Payments(date, cards)
		if !slices.EqualFunc(feed, again, slices.Equal) {
			t.Fatalf("%s: two feeds from the same seed differ", date)
		}
		for _, r := range feed {
			card, _ := g.Layouts.Feed.Text(r, "PF-CARD-NO")
			kind, _ := g.Layouts.Feed.Text(r, "PF-KIND")
			amount, _ := g.Layouts.Feed.Number(r, "PF-AMOUNT")
			kinds[kind]++
			if kind != "P" {
				continue
			}
			paid[card]++
			if c, known := owing[card]; known && (amount <= 0 || amount > c.Balance) {
				t.Errorf("%s: card %s pays %d of a balance of %d", date, card, amount, c.Balance)
			}
		}
	}
	twice := 0
	for card, n := range paid {
		if n > 1 {
			twice++
		}
		if _, known := owing[card]; !known && paid[card] > 0 && slices.ContainsFunc(cards, func(c Card) bool { return c.No == card }) {
			t.Errorf("card %s owes nothing on its statement and paid", card)
		}
	}
	if len(paid) < len(owing)*8/10 || twice > 2 || kinds["S"] == 0 {
		t.Errorf("%d of %d cards owing paid in a month, %d of them twice; %d status changes", len(paid), len(owing), twice, kinds["S"])
	}
}
