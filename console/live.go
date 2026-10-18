package console

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vlumi/copal/bedrock/dataset"
	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/pangaea/online"
)

// LiveView is the online channel as the page opens: whether an online day
// is open, and for which business date.
type LiveView struct {
	Open bool
	Date string
}

func (s *Server) live(_ *Ledger, h bank.History, _ *http.Request) (any, error) {
	return LiveView{Open: s.online.IsOpen(), Date: h.Date.String()}, nil
}

// liveEvent is an online event as the live page gets it.
type liveEvent struct {
	Time      string       `json:"time"`
	Kind      string       `json:"kind"` // opened, closed or authorization
	Function  string       `json:"function,omitempty"`
	Card      string       `json:"card,omitempty"`
	Account   int64        `json:"account,omitempty"`
	Terminal  string       `json:"terminal,omitempty"`
	Requested online.Cents `json:"requested,omitempty"`
	online.Response
}

// events streams the online channel's events to the live page as
// server-sent events, each an authorization with the account of its card.
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming is not supported", http.StatusInternalServerError)
		return
	}
	accounts, err := s.cardAccounts()
	if err != nil {
		s.fail(w, err)
		return
	}
	events, cancel := s.online.Subscribe()
	defer cancel()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	flusher.Flush()
	for {
		select {
		case <-r.Context().Done():
			return
		case e, ok := <-events:
			if !ok {
				return
			}
			data, err := json.Marshal(newLiveEvent(e, accounts))
			if err != nil {
				return
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func newLiveEvent(e online.Event, accounts map[int64]int64) liveEvent {
	le := liveEvent{Time: e.Time.Format("15:04:05")}
	switch {
	case e.Opened:
		le.Kind = "opened"
	case e.Closed:
		le.Kind = "closed"
	default:
		le.Kind = "authorization"
		le.Function = e.Request.Function
		le.Card = strconv.FormatInt(e.Request.Card, 10)
		le.Account = accounts[e.Request.Card]
		le.Terminal = e.Request.Terminal
		le.Requested = e.Request.Amount
		le.Response = e.Response
	}
	return le
}

// cardAccounts maps the current card master's cards to their accounts.
func (s *Server) cardAccounts() (map[int64]int64, error) {
	s.busy.RLock()
	defer s.busy.RUnlock()
	store, err := s.bank.Datasets()
	if err != nil {
		return nil, err
	}
	accounts := map[int64]int64{}
	n := store.LatestGeneration("CARD.MASTER")
	if n == 0 {
		return accounts, nil
	}
	records, err := store.Read(dataset.GenerationName("CARD.MASTER", n))
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		card, err := s.layouts.Card.Number(r, "CR-CARD-NO")
		if err != nil {
			return nil, err
		}
		if accounts[card], err = s.layouts.Card.Number(r, "CR-ACCT-NO"); err != nil {
			return nil, err
		}
	}
	return accounts, nil
}
