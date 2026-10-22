// Package generate makes up Laurasia Card Services' portfolio and its
// days: the card accounts and the fee schedule the processor opens with,
// the merchants its cardholders buy from, and for each day what the
// authorization switch logged, the clearing file the networks sent and
// the day's dispute events. Everything follows from a seed, so the same
// seed and date always give the same records.
package generate

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/vlumi/cobolith/record"
	"github.com/vlumi/copal/laurasia/records"
)

// Generator makes up the processor's records for one seed.
type Generator struct {
	Layouts *records.Layouts
	Seed    uint64
}

// Card is what the generator knows of a card: its account record.
type Card struct {
	No      string
	Account string
	Product string // CL classic, GD gold, BZ business
	Status  string // A active, B blocked, C closed
	Limit   int64  // cents
	Balance int64  // cents
	// Statement is the balance the last statement showed, and PaidInFull
	// whether it was paid in full.
	Statement  int64
	PaidInFull bool
	CycleDay   int
}

// Active reports whether the card's purchases are approved.
func (c Card) Active() bool { return c.Status == "A" }

// Case is what the generator knows of a dispute: its case in the master.
type Case struct {
	No       int64
	Status   string // OP, CB, RP, PA, WN or LS
	Merchant string
	Amount   int64
}

// Settled reports whether the case is over, won or lost.
func (c Case) Settled() bool { return c.Status == "WN" || c.Status == "LS" }

// Sale is a purchase a cardholder might dispute: a presentment matched to
// its approval, as the night priced it.
type Sale struct {
	Card     string
	Merchant string
	Amount   int64
}

// rng returns a random source for one purpose, so that changing how one
// thing is generated does not change everything else.
func (g *Generator) rng(stream string, n int64) *rand.Rand {
	h := uint64(14695981039346656037)
	for _, c := range []byte(stream) {
		h = (h ^ uint64(c)) * 1099511628211
	}
	return rand.New(rand.NewPCG(g.Seed, h^uint64(n)))
}

// ReadCards reads cards from account records, as the card system's
// extract or the night's unload has them.
func (g *Generator) ReadCards(records [][]byte) ([]Card, error) {
	cards := make([]Card, 0, len(records))
	for _, r := range records {
		var c Card
		var cycleDay int64
		var paid string
		if err := get(g.Layouts.Account, r, map[string]any{
			"CA-CARD-NO": &c.No, "CA-ACCOUNT-NO": &c.Account, "CA-PRODUCT": &c.Product, "CA-STATUS": &c.Status,
			"CA-LIMIT": &c.Limit, "CA-BALANCE": &c.Balance, "CA-STATEMENT-BALANCE": &c.Statement,
			"CA-CYCLE-DAY": &cycleDay, "CA-PAID-IN-FULL": &paid,
		}); err != nil {
			return nil, err
		}
		c.CycleDay, c.PaidInFull = int(cycleDay), paid == "Y"
		cards = append(cards, c)
	}
	return cards, nil
}

// ReadCases reads the dispute master.
func (g *Generator) ReadCases(records [][]byte) ([]Case, error) {
	cases := make([]Case, 0, len(records))
	for _, r := range records {
		var c Case
		if err := get(g.Layouts.Case, r, map[string]any{
			"DC-CASE-NO": &c.No, "DC-STATUS": &c.Status, "DC-MERCHANT": &c.Merchant, "DC-AMOUNT": &c.Amount,
		}); err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, nil
}

// ReadSales reads the purchases among a night's priced presentments: the
// ones matched to their approval, which posted.
func (g *Generator) ReadSales(records [][]byte) ([]Sale, error) {
	var sales []Sale
	for _, r := range records {
		var s Sale
		var match string
		if err := get(g.Layouts.Pres, r, map[string]any{
			"PS-CARD-NO": &s.Card, "PS-MERCHANT": &s.Merchant, "PS-AMOUNT": &s.Amount, "PS-MATCH": &match,
		}); err != nil {
			return nil, err
		}
		if match == "M" {
			sales = append(sales, s)
		}
	}
	return sales, nil
}

// Records splits a line sequential file into records of a layout's
// length: a record a line, the trailing spaces the write dropped put back.
func Records(l *record.Layout, file []byte) [][]byte {
	var recs [][]byte
	for _, line := range bytes.Split(bytes.TrimRight(file, "\n"), []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		r := l.Blank()
		copy(r, line)
		recs = append(recs, r)
	}
	return recs
}

// File joins records into a line sequential file, a record a line.
func File(records [][]byte) []byte {
	var b bytes.Buffer
	for _, r := range records {
		b.Write(r)
		b.WriteByte('\n')
	}
	return b.Bytes()
}

// get reads fields from a record into the pointers given for them,
// numbers and text alike.
func get(l *record.Layout, r []byte, fields map[string]any) error {
	for name, target := range fields {
		var err error
		switch v := target.(type) {
		case *int64:
			*v, err = l.Number(r, name)
		case *string:
			*v, err = l.Text(r, name)
		default:
			err = fmt.Errorf("%s: %T is neither a number nor text", name, target)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// set writes fields into a record, numbers and text alike.
func set(l *record.Layout, r []byte, fields map[string]any) error {
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		var err error
		switch v := fields[name].(type) {
		case int64:
			err = l.SetNumber(r, name, v)
		case int:
			err = l.SetNumber(r, name, int64(v))
		case string:
			err = l.SetText(r, name, v)
		default:
			err = fmt.Errorf("%s: %T is neither a number nor text", name, v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// made builds a record of a layout from its fields.
func made(l *record.Layout, fields map[string]any) ([]byte, error) {
	r := l.Blank()
	return r, set(l, r, fields)
}

// luhn completes a card number: the check digit that makes the body pass
// the Luhn test.
func luhn(body string) string {
	sum := 0
	double := true
	for i := len(body) - 1; i >= 0; i-- {
		d := int(body[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return body + string(rune('0'+(10-sum%10)%10))
}

// pick returns an index drawn by weight.
func pick(r *rand.Rand, weights []int) int {
	total := 0
	for _, w := range weights {
		total += w
	}
	n := r.IntN(total)
	for i, w := range weights {
		if n < w {
			return i
		}
		n -= w
	}
	return len(weights) - 1
}

// code makes up an authorization code: six capitals and digits.
func code(r *rand.Rand) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b strings.Builder
	for range 6 {
		b.WriteByte(alphabet[r.IntN(len(alphabet))])
	}
	return b.String()
}
