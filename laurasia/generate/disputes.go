package generate

// move is how a dispute may move on from where it stands: the event, and
// its chance in a hundred on any day.
type move struct {
	event  string
	chance int
}

// moves are the events a case at each status may take, as the dispute
// program allows them: the chargeback, then the merchant accepting or
// representing, then the issuer dropping it or going to pre-arbitration,
// where the network rules for one side.
var moves = map[string][]move{
	"OP": {{"CB", 50}},
	"CB": {{"AC", 15}, {"RP", 30}},
	"RP": {{"DR", 15}, {"PA", 25}},
	"PA": {{"RI", 30}, {"RM", 30}},
}

var events = []string{"OP", "CB", "RP", "PA", "AC", "DR", "RI", "RM"}

// Disputes makes up a day's dispute events: the cases in flight moving
// on, now and then by an event their status does not take, and new cases
// the cardholders open on last night's purchases, some charged back the
// same day.
func (g *Generator) Disputes(date Date, cases []Case, sales []Sale) ([][]byte, error) {
	r := g.rng("disputes", date.Number())
	var records [][]byte
	add := func(no int64, event, merchant string, amount int64) error {
		rec, err := made(g.Layouts.Event, map[string]any{
			"DE-CASE-NO": no, "DE-EVENT": event, "DE-MERCHANT": merchant, "DE-AMOUNT": amount, "DE-DATE": date.Number(),
		})
		records = append(records, rec)
		return err
	}
	var next int64 = 1
	for _, c := range cases {
		next = max(next, c.No+1)
		if c.Settled() {
			continue
		}
		if r.IntN(30) == 0 {
			if err := add(c.No, events[r.IntN(len(events))], c.Merchant, 0); err != nil {
				return nil, err
			}
			continue
		}
		roll := r.IntN(100)
		for _, m := range moves[c.Status] {
			if roll < m.chance {
				if err := add(c.No, m.event, c.Merchant, 0); err != nil {
					return nil, err
				}
				break
			}
			roll -= m.chance
		}
	}
	for _, s := range sales {
		if r.IntN(1000) >= 4 {
			continue
		}
		if err := add(next, "OP", s.Merchant, s.Amount); err != nil {
			return nil, err
		}
		if r.IntN(100) < 30 {
			if err := add(next, "CB", s.Merchant, 0); err != nil {
				return nil, err
			}
		}
		next++
	}
	return records, nil
}
