package online

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/pangaea/layouts"
	"github.com/vlumi/copal/pangaea/runner"
)

// Server is the bank's online channel: open for authorizations while an
// online day runs, closed during the batch window between them.
type Server struct {
	Bank    bank.Bank
	Layouts *layouts.Layouts
	// Clock gives the time of day that stamps each request.
	Clock func() time.Time

	mu          sync.Mutex
	day         *bank.OnlineDay
	subscribers map[chan Authorization]bool
}

// Authorization is a request the server answered, and when.
type Authorization struct {
	Time     time.Time
	Request  Request
	Response Response
}

// Open opens the online day of the bank's business date, unless one is
// open.
func (s *Server) Open() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.day != nil {
		return nil
	}
	day, err := s.Bank.OpenDay()
	if err != nil {
		return err
	}
	s.day = day
	return nil
}

// Close closes the online day, if one is open, for the batch window, and
// returns how the ONLINE job's steps ended.
func (s *Server) Close() ([]runner.Result, error) {
	s.mu.Lock()
	day := s.day
	s.day = nil
	s.mu.Unlock()
	if day == nil {
		return nil, nil
	}
	return day.Close()
}

// IsOpen tells whether an online day is open.
func (s *Server) IsOpen() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.day != nil
}

// Authorize passes a request to the online day and returns its response,
// or the bank-closed response when no online day is open.
func (s *Server) Authorize(r Request) (Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.day == nil {
		return Response{Result: ResultClosed, Message: "BANK CLOSED"}, nil
	}
	now := s.Clock()
	record, err := encode(s.Layouts, r, s.day.Date.Number(), int64(now.Hour()*10000+now.Minute()*100+now.Second()))
	if err != nil {
		return Response{Result: ResultInvalid, Message: "INVALID REQUEST"}, nil
	}
	answer, err := s.day.Authorize(record)
	if errors.Is(err, bank.ErrOnlineDayClosed) {
		s.day = nil
		return Response{Result: ResultClosed, Message: "BANK CLOSED"}, nil
	} else if err != nil {
		return Response{}, err
	}
	response, err := decode(s.Layouts, answer)
	if err != nil {
		return Response{}, err
	}
	s.publish(Authorization{Time: now, Request: r, Response: response})
	return response, nil
}

// Subscribe returns the authorizations answered from now on, until cancel
// is called. A subscriber that falls behind misses some.
func (s *Server) Subscribe() (authorizations <-chan Authorization, cancel func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.subscribers == nil {
		s.subscribers = map[chan Authorization]bool{}
	}
	c := make(chan Authorization, 64)
	s.subscribers[c] = true
	return c, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.subscribers[c] {
			delete(s.subscribers, c)
			close(c)
		}
	}
}

func (s *Server) publish(a Authorization) {
	for c := range s.subscribers {
		select {
		case c <- a:
		default:
		}
	}
}

// Handler serves the online channel's API: POST /purchase, /reversal,
// /withdrawal and /balance, each with a JSON request, and GET /status.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /purchase", s.serve(func(m message) Request {
		return Request{Function: Purchase, Amount: m.Amount, Terminal: m.Merchant}
	}))
	mux.HandleFunc("POST /reversal", s.serve(func(m message) Request {
		return Request{Function: Reversal, AuthCode: m.AuthCode, Terminal: m.Merchant}
	}))
	mux.HandleFunc("POST /withdrawal", s.serve(func(m message) Request {
		return Request{Function: Withdrawal, Amount: m.Amount, Terminal: m.ATM}
	}))
	mux.HandleFunc("POST /balance", s.serve(func(m message) Request {
		return Request{Function: Balance, Terminal: m.ATM}
	}))
	mux.HandleFunc("GET /status", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"open": s.IsOpen()})
	})
	return mux
}

// message is a request as a terminal sends it.
type message struct {
	Card     string `json:"card"`
	Amount   Cents  `json:"amount"`
	AuthCode int64  `json:"authCode,string"`
	Merchant string `json:"merchant"`
	ATM      string `json:"atm"`
}

func (s *Server) serve(request func(message) Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m message
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&m); err != nil {
			writeJSON(w, http.StatusBadRequest, Response{Result: ResultInvalid, Message: "INVALID REQUEST"})
			return
		}
		req := request(m)
		card, err := strconv.ParseInt(m.Card, 10, 64)
		if err != nil || len(m.Card) != 16 || len(req.Terminal) > 20 {
			writeJSON(w, http.StatusBadRequest, Response{Result: ResultInvalid, Message: "INVALID REQUEST"})
			return
		}
		req.Card = card
		response, err := s.Authorize(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		status := http.StatusOK
		if response.Result == ResultClosed {
			status = http.StatusServiceUnavailable
		}
		writeJSON(w, status, response)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
