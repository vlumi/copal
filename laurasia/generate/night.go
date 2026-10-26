package generate

import (
	"maps"
	"math/rand/v2"
	"slices"

	"github.com/vlumi/cobolith/record"
)

// Night is what a business day's night finds in its inbound directory:
// the switch's log of the authorizations since the last night, the
// network's clearing file of the purchases at other acquirers' merchants,
// and our own acquiring's on-us file of those at ours, in the same layout.
type Night struct {
	Log      [][]byte
	Clearing [][]byte
	OnUs     [][]byte
	// Presentments is how many presentments the two files carry.
	Presentments int
}

// network is the network the clearing file comes from, and onUs what the
// on-us file names in its place.
const (
	network = "LCN1"
	onUs    = "ONUS"
)

// maxBatch is the most presentments an acquirer sends in one batch.
const maxBatch = 500

// presentment is a purchase as the clearing file presents it.
type presentment struct {
	card     string
	product  string
	amount   int64
	auth     string
	merchant Merchant
	date     Date
}

// traffic collects what the switch and the networks saw between two
// nights.
type traffic struct {
	g            *Generator
	r            *rand.Rand
	cards        []Card
	merchants    []Merchant
	weights      []int
	balance      map[string]int64
	codes        map[string]bool
	log          [][]byte
	presentments []presentment
	err          error
}

// Night makes up the night of a business date: every card's purchases on
// each day since the last business day, as the switch logged them and
// as the networks present them. The switch approves against the cards as
// the last night left them, and against what it approved since; it
// stands in now and then, when the issuer could not answer. The networks
// present most of what was approved, some of it a day or two late with
// its authorization long gone from the log, and now and then a batch
// whose trailer does not add up.
func (g *Generator) Night(date, since Date, cards []Card, merchants []Merchant) (Night, error) {
	t := &traffic{
		g: g, r: g.rng("night", date.Number()), cards: cards, merchants: merchants,
		balance: map[string]int64{}, codes: map[string]bool{},
	}
	for _, m := range merchants {
		t.weights = append(t.weights, m.trade.weight*m.size)
	}
	for _, c := range cards {
		t.balance[c.No] = c.Balance
	}
	for day := since.AddDays(1); !day.After(date); day = day.AddDays(1) {
		t.day(day)
	}
	if t.err != nil {
		return Night{}, t.err
	}
	var ours, theirs []presentment
	for _, p := range t.presentments {
		if p.merchant.Acquirer == OwnAcquirer {
			ours = append(ours, p)
		} else {
			theirs = append(theirs, p)
		}
	}
	clearing, err := t.clearingFile(date, network, theirs)
	if err != nil {
		return Night{}, err
	}
	onUsFile, err := t.clearingFile(date, onUs, ours)
	return Night{Log: t.log, Clearing: clearing, OnUs: onUsFile, Presentments: len(t.presentments)}, err
}

// day is one calendar day's purchases, and a few presentments of earlier
// days' purchases arriving late.
func (t *traffic) day(day Date) {
	for _, c := range t.cards {
		for range t.purchases(c, day) {
			t.purchase(c, day)
		}
	}
	for range 1 + len(t.cards)/500 {
		t.latePresentment(day)
	}
}

// purchases is how many purchases a card makes on a day: by its product
// and the day of the week, fewer when it is blocked, hardly any when it
// is closed.
func (t *traffic) purchases(c Card, day Date) int {
	p := productNamed(c.Product)
	rate := p.weekdayRate
	if !day.BusinessDay() {
		rate = p.weekendRate
	}
	switch c.Status {
	case "B":
		rate /= 5
	case "C":
		rate /= 50
	}
	n := 0
	for t.r.Float64() < rate {
		n++
		rate *= 0.4
	}
	return n
}

// purchase is one purchase: its authorization at the switch, and its
// presentment if the merchant presents it. A merchant that takes tips
// presents more than was approved; a declined purchase is presented now
// and then all the same.
func (t *traffic) purchase(c Card, day Date) {
	m := t.merchants[pick(t.r, t.weights)]
	amount := m.trade.low + t.r.Int64N(m.trade.high-m.trade.low+1)
	authorized := amount
	if m.trade.preauth > 0 {
		authorized = m.trade.preauth
	}
	card := c.No
	standIn := t.r.IntN(100) < 2
	if standIn && t.r.IntN(10) == 0 {
		card = t.g.cardNumber(t.r, map[string]bool{}) // a card the processor has not been told of yet
	}
	approved := standIn || t.approves(c, authorized)
	auth := t.logEntry(card, authorized, approved, standIn, day)
	if approved {
		t.balance[c.No] += authorized
	}
	switch {
	case approved && t.r.IntN(100) < 96:
		if m.trade.tips {
			amount += amount * int64(t.r.IntN(23)) / 100
		}
		t.present(card, c.Product, amount, auth, m, day)
	case !approved && t.r.IntN(100) < 8:
		t.present(card, c.Product, amount, auth, m, day)
	}
}

// approves is the switch's decision on a card it knows: the card takes
// purchases, the purchase fits under the limit and what the product lets
// a card go over it, and nothing looks wrong.
func (t *traffic) approves(c Card, amount int64) bool {
	if !c.Active() {
		return false
	}
	limit := c.Limit * int64(100+productNamed(c.Product).limitPremium) / 100
	if t.balance[c.No]+amount > limit {
		return false
	}
	return t.r.IntN(100) >= 1
}

// latePresentment is a presentment of an earlier day's purchase, whose
// authorization was in an earlier log, so that tonight it matches
// nothing.
func (t *traffic) latePresentment(day Date) {
	c := t.cards[t.r.IntN(len(t.cards))]
	m := t.merchants[pick(t.r, t.weights)]
	amount := m.trade.low + t.r.Int64N(m.trade.high-m.trade.low+1)
	t.present(c.No, c.Product, amount, t.code(), m, day.AddDays(-1-t.r.IntN(2)))
}

// logEntry adds an authorization to the switch's log.
func (t *traffic) logEntry(card string, amount int64, approved, standIn bool, day Date) string {
	auth := t.code()
	decision, stood := "D", "N"
	if approved {
		decision = "A"
	}
	if standIn {
		stood = "Y"
	}
	rec, err := made(t.g.Layouts.Auth, map[string]any{
		"AL-AUTH-CODE": auth, "AL-CARD-NO": card, "AL-AMOUNT": amount,
		"AL-DECISION": decision, "AL-STAND-IN": stood, "AL-DATE": day.Number(),
	})
	t.fail(err)
	t.log = append(t.log, rec)
	return auth
}

// code is an authorization code no entry of the log has yet.
func (t *traffic) code() string {
	for {
		c := code(t.r)
		if !t.codes[c] {
			t.codes[c] = true
			return c
		}
	}
}

func (t *traffic) present(card, product string, amount int64, auth string, m Merchant, day Date) {
	t.presentments = append(t.presentments, presentment{card, product, amount, auth, m, day})
}

func (t *traffic) fail(err error) {
	if t.err == nil {
		t.err = err
	}
}

// clearingFile lays presentments out as a clearing file, the network's or
// the on-us one: a batch or more an acquirer, each with its trailer, and
// the file's. A presentment in a few thousand is mistyped, and a batch in
// a few dozen files has a trailer a cent out.
func (t *traffic) clearingFile(date Date, from string, presentments []presentment) ([][]byte, error) {
	byAcquirer := map[string][]presentment{}
	for _, p := range presentments {
		byAcquirer[p.merchant.Acquirer] = append(byAcquirer[p.merchant.Acquirer], p)
	}
	var batches [][]presentment
	for _, acquirer := range slices.Sorted(maps.Keys(byAcquirer)) {
		ps := byAcquirer[acquirer]
		t.r.Shuffle(len(ps), func(i, j int) { ps[i], ps[j] = ps[j], ps[i] })
		for len(ps) > 0 {
			n := min(len(ps), maxBatch)
			batches = append(batches, ps[:n])
			ps = ps[n:]
		}
	}
	badBatch := -1
	if len(batches) > 0 && t.r.IntN(25) == 0 {
		badBatch = t.r.IntN(len(batches))
	}
	l := t.g.Layouts
	var file [][]byte
	add := func(rec []byte, err error) {
		t.fail(err)
		file = append(file, rec)
	}
	add(clearing("FH", l.FileHeader, map[string]any{"FH-NETWORK": from, "FH-FILE-DATE": date.Number(), "FH-FILE-SEQ": 1}))
	var records int
	var total int64
	for i, batch := range batches {
		add(clearing("BH", l.BatchHeader, map[string]any{"BH-BATCH-NO": i + 1, "BH-ACQUIRER": batch[0].merchant.Acquirer}))
		var amount int64
		for j := range batch {
			add(t.presentmentRecord(&batch[j]))
			amount += batch[j].amount
		}
		records += len(batch)
		total += amount
		if i == badBatch {
			amount++
		}
		add(clearing("BT", l.BatchTrailer, map[string]any{"BT-COUNT": len(batch), "BT-AMOUNT": amount}))
	}
	add(clearing("FT", l.FileTrailer, map[string]any{"FT-BATCH-COUNT": len(batches), "FT-RECORD-COUNT": records, "FT-AMOUNT": total}))
	return file, t.err
}

// presentmentRecord is a presentment as the file carries it, now and then
// mistyped: a letter in the card number, or no amount.
func (t *traffic) presentmentRecord(p *presentment) ([]byte, error) {
	card := p.card
	switch t.r.IntN(3000) {
	case 0:
		at := t.r.IntN(len(card))
		card = card[:at] + "O" + card[at+1:]
	case 1:
		p.amount = 0
	}
	return clearing("PR", t.g.Layouts.Presentment, map[string]any{
		"PR-CARD-NO": card, "PR-AMOUNT": p.amount, "PR-AUTH-CODE": p.auth,
		"PR-MERCHANT": p.merchant.ID, "PR-MCC": p.merchant.MCC, "PR-REGION": p.merchant.Region,
		"PR-PRODUCT": p.product, "PR-TXN-DATE": p.date.Number(),
	})
}

func productNamed(code string) product {
	for _, p := range products {
		if p.code == code {
			return p
		}
	}
	return products[0]
}

// clearing builds a record of the clearing file: its fields, and the two
// characters every record of the file starts with, which tell its type.
func clearing(kind string, l *record.Layout, fields map[string]any) ([]byte, error) {
	rec, err := made(l, fields)
	copy(rec, kind)
	return rec, err
}
