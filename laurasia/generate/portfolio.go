package generate

import (
	"fmt"
	"math/rand/v2"
)

// product is a card product and how its cards are issued.
type product struct {
	code         string
	weight       int
	limitUnit    int64 // the limit is a multiple of this, in cents
	limitMin     int
	limitMax     int
	weekdayRate  float64 // purchases a card makes a weekday, on average
	weekendRate  float64
	paidInFull   int // percent of the cards that paid their last statement in full
	limitPremium int // percent the switch lets a card go over its limit by
}

var products = []product{
	{"CL", 60, 10000, 10, 50, 0.3, 0.4, 35, 0},
	{"GD", 30, 50000, 10, 40, 0.35, 0.45, 55, 5},
	{"BZ", 10, 100000, 10, 50, 0.6, 0.1, 70, 10},
}

// Cards makes up the card system's extract of about n cards: their
// accounts as the processor opens with them, by card number.
func (g *Generator) Cards(n int) ([][]byte, error) {
	r := g.rng("cards", int64(n))
	weights := make([]int, len(products))
	for i, p := range products {
		weights[i] = p.weight
	}
	issued := map[string]bool{}
	var records [][]byte
	for i := range n {
		p := products[pick(r, weights)]
		no := g.cardNumber(r, issued)
		limit := p.limitUnit * int64(p.limitMin+r.IntN(p.limitMax-p.limitMin+1))
		status := "A"
		switch s := r.IntN(100); {
		case s < 1:
			status = "C"
		case s < 4:
			status = "B"
		}
		balance := int64(float64(limit) * r.Float64() * 0.7)
		statement := balance * int64(50+r.IntN(50)) / 100
		paid := "N"
		if r.IntN(100) < p.paidInFull {
			paid = "Y"
		}
		if status == "C" {
			balance, statement, paid = 0, 0, "N"
		}
		rec, err := made(g.Layouts.Account, map[string]any{
			"CA-CARD-NO": no, "CA-ACCOUNT-NO": fmt.Sprintf("ACC%07d", i+1), "CA-PRODUCT": p.code, "CA-STATUS": status,
			"CA-LIMIT": limit, "CA-BALANCE": balance, "CA-STATEMENT-BALANCE": statement,
			"CA-CYCLE-DAY": 1 + r.IntN(28), "CA-PAID-IN-FULL": paid,
		})
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

// bin is the issuer's bank identification number, the first digits of
// every card it issues.
const bin = "400001"

// cardNumber makes up a card number nobody holds yet.
func (g *Generator) cardNumber(r *rand.Rand, issued map[string]bool) string {
	for {
		no := luhn(fmt.Sprintf("%s%09d", bin, r.Int64N(1_000_000_000)))
		if !issued[no] {
			issued[no] = true
			return no
		}
	}
}

// Merchant is a merchant the cardholders buy from: its category, where
// it is, who acquires for it, and what a purchase there costs.
type Merchant struct {
	ID       string
	MCC      int
	Region   string // D domestic, I international
	Acquirer string
	trade    trade
	size     int // how busy it is, 1 to 5
}

// trade is a line of business: its merchant category code, its typical
// tickets, how many merchants are in it, and how it authorizes.
type trade struct {
	mcc     int
	low     int64 // a ticket, in cents
	high    int64
	weight  int
	tips    bool  // the presentment adds a tip to the approval
	preauth int64 // the approval is for this much, and the presentment for what was taken
}

var trades = []trade{
	{5411, 800, 18000, 20, false, 0},    // grocery
	{5499, 300, 4000, 6, false, 0},      // convenience stores
	{5812, 1500, 12000, 14, true, 0},    // restaurants
	{5814, 500, 2500, 10, false, 0},     // fast food
	{5813, 800, 6000, 3, true, 0},       // bars
	{5541, 2000, 9000, 8, false, 10000}, // fuel
	{5983, 5000, 40000, 1, false, 0},    // fuel dealers
	{4111, 200, 1500, 5, false, 0},      // transit
	{4121, 800, 6000, 4, true, 0},       // taxis
	{4511, 8000, 90000, 2, false, 0},    // airlines
	{3012, 10000, 120000, 1, false, 0},  // an airline of its own
	{4722, 20000, 250000, 1, false, 0},  // travel agencies
	{7011, 9000, 60000, 3, false, 0},    // hotels
	{5311, 1500, 30000, 5, false, 0},    // department stores
	{5732, 3000, 150000, 2, false, 0},   // electronics
	{5912, 500, 8000, 6, false, 0},      // drugstores
	{5999, 500, 20000, 5, false, 0},     // other retail
	{5942, 800, 5000, 2, false, 0},      // books
	{7832, 1000, 4000, 2, false, 0},     // cinemas
	{4900, 3000, 25000, 3, false, 0},    // utilities
}

// acquirers are the acquirers of the merchants: the first three domestic,
// the others abroad. The first is our own acquiring, whose merchants'
// presentments come on-us, not through the network.
var acquirers = []string{OwnAcquirer, "ACQ002", "ACQ003", "ACQ004", "ACQ005"}

// OwnAcquirer is our own acquiring business, and ownMerchants the range
// its merchants are numbered in.
const (
	OwnAcquirer  = "LCSACQ"
	ownMerchants = "LCS"
)

// Merchants makes up the merchants: the same for a seed, whatever the
// portfolio.
func (g *Generator) Merchants() []Merchant {
	const count = 60
	r := g.rng("merchants", 0)
	weights := make([]int, len(trades))
	for i, t := range trades {
		weights[i] = t.weight
	}
	merchants := make([]Merchant, 0, count)
	for i := range count {
		m := Merchant{ID: fmt.Sprintf("MER%07d", 1000+i), trade: trades[pick(r, weights)], size: 1 + r.IntN(5)}
		m.MCC = m.trade.mcc
		if r.IntN(100) < 85 {
			m.Region, m.Acquirer = "D", acquirers[r.IntN(3)]
		} else {
			m.Region, m.Acquirer = "I", acquirers[3+r.IntN(2)]
		}
		if m.Acquirer == OwnAcquirer {
			m.ID = fmt.Sprintf("%s%07d", ownMerchants, 1000+i)
		}
		merchants = append(merchants, m)
	}
	return merchants
}

// Fees makes up the fee schedule: a rate and a fixed fee for each card
// product, region and merchant group, but for a few combinations the
// schedule has never had, which take the default fee.
func (g *Generator) Fees() ([][]byte, error) {
	groupRates := []int64{110, 150, 180, 120, 160} // basis points, by merchant group
	premium := map[string]int64{"CL": 0, "GD": 30, "BZ": 50}
	missing := map[string]bool{"BZ I 1": true, "BZ I 2": true, "GD I 4": true}
	var records [][]byte
	for _, product := range []string{"CL", "GD", "BZ"} {
		for _, region := range []string{"D", "I"} {
			for group := 1; group <= 5; group++ {
				if missing[fmt.Sprintf("%s %s %d", product, region, group)] {
					continue
				}
				rate := groupRates[group-1] + premium[product]
				fixed := int64(10 + 5*(group-1))
				if region == "I" {
					rate, fixed = rate+40, fixed+10
				}
				rec, err := made(g.Layouts.Fee, map[string]any{
					"FD-PRODUCT": product, "FD-REGION": region, "FD-GROUP": group, "FD-RATE": rate, "FD-FIXED": fixed,
				})
				if err != nil {
					return nil, err
				}
				records = append(records, rec)
			}
		}
	}
	return records, nil
}
