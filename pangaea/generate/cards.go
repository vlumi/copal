package generate

// cardPrefix starts every card number the bank issues, followed by an
// 11-digit serial number and a check digit.
const cardPrefix = 9100

// Cards makes up the card master: a card for each current account that is
// not closed, and for some a second card, in card number order, which is
// account order. A few cards are blocked, and a few expired before date.
func (g *Generator) Cards(accounts []Account, date Date) ([][]byte, error) {
	r := g.rng("cards", 0)
	var records [][]byte
	serial := int64(0)
	for _, a := range accounts {
		if a.Savings || a.Status == "X" {
			continue
		}
		n := 1
		if r.Float64() < 0.15 {
			n = 2
		}
		for range n {
			serial++
			status := "A"
			if r.Float64() < 0.02 {
				status = "B"
			}
			expiry := monthNumber(date, 1+r.IntN(48))
			if r.Float64() < 0.01 {
				expiry = monthNumber(date, -1)
			}
			record := g.Layouts.Card.Blank()
			if err := set(g.Layouts.Card, record, map[string]any{
				"CR-CARD-NO":           cardNumber(serial),
				"CR-ACCT-NO":           a.No,
				"CR-CUST-NO":           a.Customer,
				"CR-STATUS":            status,
				"CR-EXPIRY-MONTH":      expiry,
				"CR-PURCHASE-LIMIT":    cents([]int{500, 1000, 1000, 2000, 5000}[r.IntN(5)]),
				"CR-CASH-LIMIT":        cents([]int{200, 400, 400, 1000}[r.IntN(4)]),
				"CR-TOTALS-DATE":       int64(0),
				"CR-DAY-PURCHASES":     int64(0),
				"CR-DAY-CASH":          int64(0),
				"CR-MONTH-WITHDRAWALS": int64(0),
			}); err != nil {
				return nil, err
			}
			records = append(records, record)
		}
	}
	return records, nil
}

// monthNumber is the month n months from date's, as yyyymm.
func monthNumber(date Date, n int) int64 {
	t := date.t.AddDate(0, n, 1-date.Day())
	return int64(t.Year()*100 + int(t.Month()))
}

// cardNumber is the 16-digit card number of a serial number: the bank's
// prefix, the serial and a Luhn check digit, as card numbers have.
func cardNumber(serial int64) int64 {
	body := cardPrefix*100_000_000_000 + serial
	sum := 0
	for n, double := body, true; n > 0; n, double = n/10, !double {
		d := int(n % 10)
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return body*10 + int64((10-sum%10)%10)
}
