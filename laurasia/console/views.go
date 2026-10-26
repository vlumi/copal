package console

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strings"

	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/cobolith/record"
	"github.com/vlumi/copal/laurasia/generate"
	"github.com/vlumi/copal/laurasia/processor"
)

// errNotFound is for a card or a day the processor does not have.
var errNotFound = errors.New("not found")

// Card is a card's account record.
type Card struct {
	No, Account, Product, Status string
	Limit, Balance, Statement    int64 // cents
	CycleDay                     int64
	Record                       []byte
}

// Products and statuses by their codes, as the copybook's condition names
// say them.
var (
	products = map[string]string{"CL": "classic", "GD": "gold", "BZ": "business"}
	statuses = map[string]string{"A": "active", "B": "blocked", "C": "closed"}
)

// ProductName is the card's product in words.
func (c Card) ProductName() string { return products[c.Product] }

// StatusName is the card's status in words.
func (c Card) StatusName() string { return statuses[c.Status] }

// Overview is the processor as it stands: its cards in sum, and the days
// it has run, newest first.
type Overview struct {
	Cards    int
	ByStatus map[string]int
	Balance  int64
	Limit    int64
	Days     []processor.Day
}

func (s *Server) overview(h processor.History, _ *http.Request) (any, error) {
	cards, err := s.readCards()
	if err != nil {
		return nil, err
	}
	o := Overview{Cards: len(cards), ByStatus: map[string]int{}}
	for _, c := range cards {
		o.ByStatus[c.StatusName()]++
		o.Balance += c.Balance
		o.Limit += c.Limit
	}
	for _, d := range slices.Backward(h.Days) {
		o.Days = append(o.Days, d)
	}
	return o, nil
}

// readCards reads the cards as the last cycle left them.
func (s *Server) readCards() ([]Card, error) {
	content, err := s.processor.Cards()
	if err != nil {
		return nil, err
	}
	l := s.layouts.Account
	var cards []Card
	for _, r := range generate.Records(l, content) {
		c := Card{Record: r}
		for _, f := range []struct {
			name string
			text *string
		}{{"CA-CARD-NO", &c.No}, {"CA-ACCOUNT-NO", &c.Account}, {"CA-PRODUCT", &c.Product}, {"CA-STATUS", &c.Status}} {
			if *f.text, err = l.Text(r, f.name); err != nil {
				return nil, err
			}
		}
		for _, f := range []struct {
			name   string
			number *int64
		}{{"CA-LIMIT", &c.Limit}, {"CA-BALANCE", &c.Balance}, {"CA-STATEMENT-BALANCE", &c.Statement}, {"CA-CYCLE-DAY", &c.CycleDay}} {
			if *f.number, err = l.Number(r, f.name); err != nil {
				return nil, err
			}
		}
		cards = append(cards, c)
	}
	return cards, nil
}

// shown is the most cards a search lists.
const shown = 200

// CardList is the cards a search found, at most shown of them.
type CardList struct {
	Query string
	Found int
	Total int
	Cards []Card
}

func (s *Server) cards(_ processor.History, r *http.Request) (any, error) {
	cards, err := s.readCards()
	if err != nil {
		return nil, err
	}
	query := strings.TrimSpace(r.FormValue("q"))
	var found []Card
	for _, c := range cards {
		if query == "" || strings.HasPrefix(c.No, query) || strings.EqualFold(c.Account, query) {
			found = append(found, c)
		}
	}
	return CardList{Query: query, Found: len(found), Total: len(cards), Cards: found[:min(len(found), shown)]}, nil
}

// CardPage is a card's account record field by field, and every posting
// to it since the processor opened.
type CardPage struct {
	Card     Card
	Fields   []FieldValue
	Postings []Posting
}

// FieldValue is a field of a record: its bytes, and what they mean.
type FieldValue struct {
	Field      *layout.Field
	Raw, Value string
}

// Posting is a posting to a card.
type Posting struct {
	Date      string
	Kind      string
	Amount    int64
	Reference string
}

// kinds are the postings' kinds by their codes, as the copybook names them.
var kinds = map[string]string{"P": "purchase", "F": "over-limit fee", "I": "interest", "C": "payment"}

func (s *Server) card(_ processor.History, r *http.Request) (any, error) {
	cards, err := s.readCards()
	if err != nil {
		return nil, err
	}
	i := slices.IndexFunc(cards, func(c Card) bool { return c.No == r.PathValue("no") })
	if i < 0 {
		return nil, fmt.Errorf("card %s: %w", r.PathValue("no"), errNotFound)
	}
	p := CardPage{Card: cards[i], Fields: fields(s.layouts.Account, cards[i].Record)}
	files, err := s.processor.Postings()
	if err != nil {
		return nil, err
	}
	l := s.layouts.Posting
	for _, content := range files {
		for _, rec := range generate.Records(l, content) {
			if no, _ := l.Text(rec, "PO-CARD-NO"); no != p.Card.No {
				continue
			}
			var posting Posting
			kind, _ := l.Text(rec, "PO-KIND")
			posting.Kind = kinds[kind]
			posting.Reference, _ = l.Text(rec, "PO-REFERENCE")
			posting.Date, _ = l.Text(rec, "PO-DATE")
			if posting.Amount, err = l.Number(rec, "PO-AMOUNT"); err != nil {
				return nil, err
			}
			p.Postings = append(p.Postings, posting)
		}
	}
	return p, nil
}

// fields lays a record out field by field, numbers as decimals.
func fields(l *record.Layout, r []byte) []FieldValue {
	var values []FieldValue
	for _, f := range l.Fields() {
		v := FieldValue{Field: f, Raw: string(r[f.Start-1 : f.Start-1+f.Length])}
		if n, err := l.Number(r, f.Name); err == nil && f.Category == layout.Numeric {
			v.Value = decimal(n, f.Scale)
		} else if f.Name != "FILLER" {
			v.Value = strings.TrimRight(v.Raw, " ")
		}
		values = append(values, v)
	}
	return values
}

// decimal writes n in units of a scale with its decimals: 152000 with 2 is
// 1520.00.
func decimal(n int64, scale int) string {
	if scale == 0 {
		return fmt.Sprint(n)
	}
	sign := ""
	if n < 0 {
		sign, n = "-", -n
	}
	unit := int64(math.Pow10(scale))
	return fmt.Sprintf("%s%d.%0*d", sign, n/unit, scale, n%unit)
}

// DayPage is a day the processor ran: its steps, what it amounted to, and
// the files it read and wrote, a small one's text shown.
type DayPage struct {
	Day   processor.Day
	Files []DayFile
}

// DayFile is a file of a day, with its text, or Large when it is too large
// to show.
type DayFile struct {
	processor.File
	Text  string
	Large bool
}

// largeFile is the size from which a file is linked to, not shown.
const largeFile = 64 << 10

func (s *Server) day(h processor.History, r *http.Request) (any, error) {
	i := slices.IndexFunc(h.Days, func(d processor.Day) bool { return d.Date.String() == r.PathValue("date") })
	if i < 0 {
		return nil, fmt.Errorf("day %s: %w", r.PathValue("date"), errNotFound)
	}
	page := DayPage{Day: h.Days[i]}
	files, err := s.processor.Files(page.Day.Date)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		df := DayFile{File: f, Large: f.Size >= largeFile}
		if !df.Large {
			content, err := s.processor.ReadFile(f.Path)
			if err != nil {
				return nil, err
			}
			df.Text = string(content)
		}
		page.Files = append(page.Files, df)
	}
	return page, nil
}
