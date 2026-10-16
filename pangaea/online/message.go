// Package online serves the bank's online day over HTTP: card terminals
// and ATMs send authorization requests as JSON, the server passes them to
// AUTHSRV as AUTHMSG records, one at a time, and answers with the
// response. Between online days, during the batch window, it answers that
// the bank is closed.
package online

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vlumi/copal/pangaea/layouts"
)

// Functions of a request.
const (
	Purchase   = "PURC"
	Reversal   = "REVS"
	Withdrawal = "WDRL"
	Balance    = "BALQ"
)

// Request is an authorization request: a purchase, the reversal of one
// by its authorization code, a withdrawal or a balance enquiry, made with
// a card at a terminal, a merchant's or an ATM.
type Request struct {
	Function string
	Card     int64
	Amount   Cents
	AuthCode int64
	Terminal string
}

// Response is AUTHSRV's answer. Amounts are in cents; those of a declined
// request are zero.
type Response struct {
	Result    string `json:"result"`
	Approved  bool   `json:"approved"`
	Message   string `json:"message"`
	AuthCode  string `json:"authCode,omitempty"`
	Amount    Cents  `json:"amount"`
	Fee       Cents  `json:"fee"`
	Balance   Cents  `json:"balance"`
	Available Cents  `json:"available"`
}

// Results that the server gives itself, not AUTHSRV.
const (
	ResultInvalid = "12" // a request that is not one
	ResultClosed  = "91" // the bank is closed: the batch window
)

// encode writes a request as an AUTHMSG request record, stamped with the
// business date and the time of day, hhmmss.
func encode(l *layouts.Layouts, r Request, date, timeOfDay int64) ([]byte, error) {
	record := l.Request.Blank()
	for name, value := range map[string]int64{
		"RQ-CARD-NO": r.Card, "RQ-AMOUNT": int64(r.Amount), "RQ-AUTH-CODE": r.AuthCode,
		"RQ-TIMESTAMP": date*1_000_000 + timeOfDay,
	} {
		if err := l.Request.SetNumber(record, name, value); err != nil {
			return nil, err
		}
	}
	for name, text := range map[string]string{"RQ-FUNCTION": r.Function, "RQ-TERMINAL": r.Terminal} {
		if err := l.Request.SetText(record, name, text); err != nil {
			return nil, err
		}
	}
	return record, nil
}

// decode reads an AUTHMSG response record.
func decode(l *layouts.Layouts, record []byte) (Response, error) {
	var r Response
	var err error
	if r.Result, err = l.Response.Text(record, "RS-RESULT"); err != nil {
		return Response{}, err
	}
	if r.Message, err = l.Response.Text(record, "RS-MESSAGE"); err != nil {
		return Response{}, err
	}
	if r.AuthCode, err = l.Response.Text(record, "RS-AUTH-CODE"); err != nil {
		return Response{}, err
	}
	r.AuthCode = strings.TrimSpace(r.AuthCode)
	r.Approved = r.Result == "00"
	for name, target := range map[string]*Cents{
		"RS-AMOUNT": &r.Amount, "RS-FEE": &r.Fee, "RS-BALANCE": &r.Balance, "RS-AVAILABLE": &r.Available,
	} {
		text, err := l.Response.Text(record, name)
		if err != nil {
			return Response{}, err
		}
		if *target, err = parseEdited(text); err != nil {
			return Response{}, fmt.Errorf("%s: %w", name, err)
		}
	}
	return r, nil
}

// Cents is an amount in cents, written in JSON as a decimal string.
type Cents int64

func (c Cents) String() string {
	sign := ""
	if c < 0 {
		sign, c = "-", -c
	}
	return fmt.Sprintf("%s%d.%02d", sign, c/100, c%100)
}

// MarshalText writes the amount as 12.34.
func (c Cents) MarshalText() ([]byte, error) { return []byte(c.String()), nil }

// UnmarshalText reads an amount written as 12.34, 12.3, 12 or -12.34.
func (c *Cents) UnmarshalText(text []byte) error {
	digits, negative := strings.CutPrefix(string(text), "-")
	units, fraction, _ := strings.Cut(digits, ".")
	if len(fraction) > 2 || units == "" || strings.ContainsAny(units, "+-") {
		return fmt.Errorf("amount %q is not units and cents", text)
	}
	n, err := strconv.ParseInt(units+(fraction + "00")[:2], 10, 64)
	if err != nil {
		return fmt.Errorf("amount %q is not units and cents", text)
	}
	if negative {
		n = -n
	}
	*c = Cents(n)
	return nil
}

// parseEdited reads a numeric-edited field, such as -00000000150.00 or
// "    12.00", as cents; a blank one is zero.
func parseEdited(text string) (Cents, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, nil
	}
	var c Cents
	err := c.UnmarshalText([]byte(text))
	return c, err
}
