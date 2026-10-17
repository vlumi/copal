// Package simulate makes up the card traffic of the bank's online day:
// customers paying by card and taking cash from ATMs, at a given rate,
// sent to the online channel's API as card terminals and ATMs would. Each
// card has a spending habit that follows from the seed and the card, so
// some cards run into their daily limits and some accounts out of money.
package simulate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/online"
	"github.com/vlumi/copal/pangaea/records"
)

// Card is a card the simulator uses, and how much its holder spends: a
// typical purchase, in cents.
type Card struct {
	No       int64
	Purchase int64
}

// Simulator sends made-up requests to the online channel's API.
type Simulator struct {
	API    string // the API's base URL, such as http://127.0.0.1:8080/api
	Client *http.Client
	Cards  []Card
	rand   *rand.Rand
	recent []purchase
}

// purchase is an approved purchase that may yet be reversed.
type purchase struct {
	card     int64
	authCode string
	merchant string
}

var (
	merchants = []string{
		"CORNER SHOP", "BAKERY", "FUEL STATION", "GROCER", "BOOKSHOP", "PHARMACY",
		"RESTAURANT", "TRAIN TICKETS", "HARDWARE STORE", "ELECTRONICS", "TRAVEL AGENT",
	}
	atms = []string{"ATM 0007 MARKET SQ", "ATM 0042 HARBOR", "ATM 0101 STATION", "ATM 0230 AIRPORT"}
)

// New makes a simulator for the cards of the bank's current card master,
// following from the seed.
func New(api string, store *dataset.Store, l *records.Layouts, seed uint64) (*Simulator, error) {
	n := store.LatestGeneration("CARD.MASTER")
	if n == 0 {
		return nil, fmt.Errorf("the bank has no card master")
	}
	records, err := store.Read(dataset.GenerationName("CARD.MASTER", n))
	if err != nil {
		return nil, err
	}
	s := &Simulator{API: api, Client: &http.Client{Timeout: 30 * time.Second}, rand: rand.New(rand.NewPCG(seed, 0))}
	for _, r := range records {
		no, err := l.Card.Number(r, "CR-CARD-NO")
		if err != nil {
			return nil, err
		}
		habit := rand.New(rand.NewPCG(seed, uint64(no)))
		typical := int64(500 + habit.IntN(4000))
		if habit.Float64() < 0.1 {
			typical *= 10
		}
		s.Cards = append(s.Cards, Card{No: no, Purchase: typical})
	}
	if len(s.Cards) == 0 {
		return nil, fmt.Errorf("the bank has no cards")
	}
	return s, nil
}

// Request is a request the simulator sent and the response it got.
type Request struct {
	Path     string
	Card     int64
	Amount   online.Cents
	Response online.Response
}

func (r Request) String() string {
	return fmt.Sprintf("%-11s %d %10s  %s %s", r.Path, r.Card, r.Amount, r.Response.Result, r.Response.Message)
}

// Run sends rate requests a second until the context ends, reporting each.
func (s *Simulator) Run(ctx context.Context, rate float64, report func(Request)) error {
	ticker := time.NewTicker(time.Duration(float64(time.Second) / rate))
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			r, err := s.Send(ctx)
			if err != nil {
				return err
			}
			report(r)
		}
	}
}

// Send makes up the next request, sends it, and returns it with its
// response.
func (s *Simulator) Send(ctx context.Context) (Request, error) {
	path, r, body := s.next()
	r.Path = path
	data, err := json.Marshal(body)
	if err != nil {
		return Request{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.API+path, bytes.NewReader(data))
	if err != nil {
		return Request{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.Client.Do(req)
	if err != nil {
		return Request{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if err := json.NewDecoder(resp.Body).Decode(&r.Response); err != nil {
		return Request{}, fmt.Errorf("%s: %s: %w", path, resp.Status, err)
	}
	if r.Response.Result == "" {
		return Request{}, fmt.Errorf("%s%s answered %s with no result: not the online channel; is the console at %s?",
			s.API, path, resp.Status, s.API)
	}
	if path == "/purchase" && r.Response.Approved {
		s.remember(purchase{card: r.Card, authCode: r.Response.AuthCode, merchant: body["merchant"]})
	}
	return r, nil
}

// next makes up a request: mostly purchases, then withdrawals, balance
// enquiries and the odd reversal of a recent purchase, and now and then
// one for a card the bank never issued.
func (s *Simulator) next() (path string, r Request, body map[string]string) {
	card := s.Cards[s.rand.IntN(len(s.Cards))]
	if s.rand.Float64() < 0.01 {
		card.No = 9199999999999990 + int64(s.rand.IntN(10))
	}
	body = map[string]string{"card": strconv.FormatInt(card.No, 10)}
	r.Card = card.No
	switch x := s.rand.Float64(); {
	case x < 0.05 && len(s.recent) > 0:
		i := s.rand.IntN(len(s.recent))
		p := s.recent[i]
		s.recent = append(s.recent[:i], s.recent[i+1:]...)
		r.Card, body["card"] = p.card, strconv.FormatInt(p.card, 10)
		body["authCode"], body["merchant"] = p.authCode, p.merchant
		return "/reversal", r, body
	case x < 0.70:
		r.Amount = online.Cents(card.Purchase/4 + s.rand.Int64N(card.Purchase*2))
		body["amount"], body["merchant"] = r.Amount.String(), merchants[s.rand.IntN(len(merchants))]
		return "/purchase", r, body
	case x < 0.90:
		r.Amount = online.Cents(2000 * (1 + s.rand.Int64N(max(1, card.Purchase/1000))))
		body["amount"], body["atm"] = r.Amount.String(), atms[s.rand.IntN(len(atms))]
		return "/withdrawal", r, body
	}
	body["atm"] = atms[s.rand.IntN(len(atms))]
	return "/balance", r, body
}

// remember keeps the latest approved purchases for reversals.
func (s *Simulator) remember(p purchase) {
	s.recent = append(s.recent, p)
	if len(s.recent) > 100 {
		s.recent = s.recent[1:]
	}
}
