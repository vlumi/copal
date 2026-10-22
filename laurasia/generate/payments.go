package generate

import (
	"fmt"
	"math/rand/v2"
	"strconv"
)

// payer is how a cardholder pays: in full, the minimum, a part, or not at
// all, and how many days after the cycle day.
type payer struct {
	kind  string
	delay int
	share int // percent of the statement a partial payer pays
}

// payerOf is a card's cardholder, the same for the card every day: most
// pay in full, as the product's share of them says, and of the rest most
// pay the minimum, some a part, and a few nothing.
func (g *Generator) payerOf(c Card) payer {
	n, _ := strconv.ParseInt(c.No, 10, 64)
	r := g.rng("payer", n)
	full := productNamed(c.Product).paidInFull
	roll := r.IntN(100)
	p := payer{delay: 10 + r.IntN(12), share: 30 + r.IntN(41)}
	switch {
	case roll < full:
		p.kind = "full"
	case roll < full+(100-full)*6/10:
		p.kind = "minimum"
	case roll < full+(100-full)*9/10:
		p.kind = "partial"
	default:
		p.kind = "none"
	}
	return p
}

// Payments makes up the card system's feed of a business day: each
// cardholder's payment on their day, some days after their statement,
// and the day's status changes, a card lost or found or an account
// closed; now and then a record the processor will reject.
func (g *Generator) Payments(date Date, cards []Card) ([][]byte, error) {
	r := g.rng("payments", date.Number())
	var records [][]byte
	add := func(card, kind string, amount int64, status string) error {
		rec, err := made(g.Layouts.Feed, map[string]any{
			"PF-CARD-NO": card, "PF-KIND": kind, "PF-AMOUNT": amount, "PF-STATUS": status,
			"PF-REFERENCE": fmt.Sprintf("P%05d", len(records)+1), "PF-DATE": date.Number(),
		})
		records = append(records, rec)
		return err
	}
	for _, c := range cards {
		if amount := g.payment(c, date); amount > 0 {
			if err := add(c.No, "P", amount, " "); err != nil {
				return nil, err
			}
		}
		if status := statusChange(r, c); status != "" {
			if err := add(c.No, "S", 0, status); err != nil {
				return nil, err
			}
		}
	}
	if r.IntN(10) == 0 {
		if err := add(g.cardNumber(r, map[string]bool{}), "P", 2500, " "); err != nil {
			return nil, err
		}
	}
	if r.IntN(20) == 0 {
		if err := add(cards[r.IntN(len(cards))].No, "S", 0, "X"); err != nil {
			return nil, err
		}
	}
	return records, nil
}

// payment is what a cardholder pays on a date: their payment on the
// first business day on or after their day of the month, some days
// after their cycle day, wrapping into the next month; nothing on any
// other day, and nothing when the statement is paid or there is none.
func (g *Generator) payment(c Card, date Date) int64 {
	if c.Statement <= 0 || c.PaidInFull || c.Status == "C" && c.Balance <= 0 {
		return 0
	}
	p := g.payerOf(c)
	day := (c.CycleDay+p.delay-1)%28 + 1
	previous := date.PreviousBusinessDay()
	if date.Day() < day || previous.Day() >= day && previous.Day() <= date.Day() {
		return 0
	}
	var amount int64
	switch p.kind {
	case "full":
		amount = c.Statement
	case "minimum":
		amount = max(2500, c.Statement/50)
	case "partial":
		amount = c.Statement * int64(p.share) / 100
	}
	return min(amount, max(c.Balance, 0))
}

// statusChange is a card lost or found, or an account closed, once in a
// long while.
func statusChange(r *rand.Rand, c Card) string {
	if r.IntN(1000) != 0 {
		return ""
	}
	switch c.Status {
	case "A":
		if r.IntN(4) == 0 {
			return "C"
		}
		return "B"
	case "B":
		if r.IntN(3) == 0 {
			return "C"
		}
		return "A"
	}
	return ""
}
