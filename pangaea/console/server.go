package console

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"sync"

	"github.com/vlumi/copal/blueprint"
	"github.com/vlumi/copal/pangaea/bank"
)

//go:embed templates static
var files embed.FS

// shown is the most accounts a search lists.
const shown = 200

// Server serves the console for one bank, with the blueprints of its
// programs and jobs.
type Server struct {
	bank       bank.Bank
	layouts    *Layouts
	system     *blueprint.System
	blueprints string
	pages      map[string]*template.Template
	// busy is locked while the bank runs or resets, and read-locked while
	// a page reads its files.
	busy sync.RWMutex
}

// New makes a console for the bank. The system is Copal's analysis of the
// bank's programs and jobs, whose blueprints are in the blueprints
// directory.
func New(b bank.Bank, system *blueprint.System, blueprints string) (*Server, error) {
	layouts, err := LoadLayouts(b.Copybooks)
	if err != nil {
		return nil, err
	}
	s := &Server{bank: b, layouts: layouts, system: system, blueprints: blueprints, pages: map[string]*template.Template{}}
	funcs := template.FuncMap{"money": money, "count": count, "date": date, "legend": func() string { return blueprint.NetworkLegend }}
	for _, page := range []string{"overview", "nobank", "accounts", "account", "day", "compare", "error"} {
		t, err := template.New("").Funcs(funcs).ParseFS(files, "templates/base.html", "templates/"+page+".html")
		if err != nil {
			return nil, err
		}
		s.pages[page] = t
	}
	return s, nil
}

// Handler routes the console's pages and actions.
func (s *Server) Handler() http.Handler {
	static, _ := fs.Sub(files, "static")
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /{$}", s.view("overview", s.overview))
	mux.HandleFunc("GET /accounts", s.view("accounts", s.accounts))
	mux.HandleFunc("GET /accounts/{no}", s.view("account", s.account))
	mux.HandleFunc("GET /days/{date}", s.view("day", s.day))
	mux.HandleFunc("GET /days/{date}/sysout/{file...}", s.sysout)
	mux.Handle("GET /blueprints/", http.StripPrefix("/blueprints/", http.FileServer(http.Dir(s.blueprints))))
	mux.HandleFunc("GET /compare", s.view("compare", s.compare))
	mux.HandleFunc("POST /reset", s.reset)
	mux.HandleFunc("POST /run", s.run)
	return mux
}

// page is what every template gets: the bank's history for the header, and
// the page's own view.
type page struct {
	History bank.History
	View    any
}

type viewFunc func(*Ledger, bank.History, *http.Request) (any, error)

// view serves a page that reads the bank. Without a bank, every page offers
// to reset one.
func (s *Server) view(name string, read viewFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.busy.RLock()
		defer s.busy.RUnlock()
		h, err := s.bank.History()
		if errors.Is(err, bank.ErrNoBank) {
			s.render(w, http.StatusOK, "nobank", page{})
			return
		}
		if err != nil {
			s.fail(w, err)
			return
		}
		store, err := s.bank.Datasets()
		if err != nil {
			s.fail(w, err)
			return
		}
		v, err := read(&Ledger{Layouts: s.layouts, Store: store, Jobs: s.system.Jobs}, h, r)
		if err != nil {
			s.fail(w, err)
			return
		}
		s.render(w, http.StatusOK, name, page{History: h, View: v})
	}
}

func (s *Server) overview(l *Ledger, h bank.History, _ *http.Request) (any, error) {
	return l.Overview(h)
}

// AccountsView is a search of the current master.
type AccountsView struct {
	Query    string
	Total    int
	Found    int
	Accounts []Account
}

func (s *Server) accounts(l *Ledger, h bank.History, r *http.Request) (any, error) {
	all, err := l.Accounts(h.Master)
	if err != nil {
		return nil, err
	}
	v := AccountsView{Query: r.FormValue("q"), Total: len(all)}
	found := Search(all, v.Query)
	v.Found, v.Accounts = len(found), found[:min(len(found), shown)]
	return v, nil
}

func (s *Server) account(l *Ledger, h bank.History, r *http.Request) (any, error) {
	no, err := strconv.ParseInt(r.PathValue("no"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("account %q: %w", r.PathValue("no"), ErrNotFound)
	}
	return l.Account(h, no)
}

func (s *Server) day(l *Ledger, h bank.History, r *http.Request) (any, error) {
	return l.Day(h, r.PathValue("date"))
}

// compare compares two openings, by default the last day's. Given only
// where to start, it shows that one day.
func (s *Server) compare(l *Ledger, h bank.History, r *http.Request) (any, error) {
	openings := Openings(h)
	from, to := r.FormValue("from"), r.FormValue("to")
	i := len(openings) - 2
	if from != "" {
		i = slices.IndexFunc(openings, func(o Opening) bool { return o.Date.String() == from })
	}
	if i < 0 {
		i = 0
	}
	if from == "" {
		from = openings[i].Date.String()
	}
	if to == "" {
		to = openings[min(i+1, len(openings)-1)].Date.String()
	}
	return l.Compare(h, from, to)
}

// sysout serves a file a day's jobs wrote, such as a report too large to
// show on the day's page.
func (s *Server) sysout(w http.ResponseWriter, r *http.Request) {
	s.busy.RLock()
	defer s.busy.RUnlock()
	h, err := s.bank.History()
	if err != nil {
		s.fail(w, err)
		return
	}
	i := slices.IndexFunc(h.Days, func(d bank.RunDay) bool { return d.Date.String() == r.PathValue("date") })
	if i < 0 {
		s.fail(w, fmt.Errorf("business day %s: %w", r.PathValue("date"), ErrNotFound))
		return
	}
	file := r.Clone(r.Context())
	file.URL.Path = "/" + r.PathValue("file")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.FileServer(http.Dir(h.Days[i].Sysout)).ServeHTTP(w, file)
}

// reset resets the bank from the form's seed and number of accounts.
func (s *Server) reset(w http.ResponseWriter, r *http.Request) {
	seed, err := strconv.ParseUint(r.FormValue("seed"), 10, 64)
	if err != nil {
		s.refuse(w, fmt.Errorf("seed %q: %w", r.FormValue("seed"), err))
		return
	}
	accounts, err := strconv.Atoi(r.FormValue("accounts"))
	if err != nil || accounts < 10 || accounts > 100000 {
		s.refuse(w, fmt.Errorf("accounts %q: between 10 and 100000", r.FormValue("accounts")))
		return
	}
	s.busy.Lock()
	err = s.bank.Reset(seed, accounts)
	s.busy.Unlock()
	if err != nil {
		s.fail(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// run runs the form's number of business days and shows what they changed.
func (s *Server) run(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.FormValue("days"))
	if err != nil || n < 1 || n > 25 {
		s.refuse(w, fmt.Errorf("days %q: between 1 and 25", r.FormValue("days")))
		return
	}
	s.busy.Lock()
	from, to, err := s.runDays(n)
	s.busy.Unlock()
	if err != nil {
		s.fail(w, err)
		return
	}
	http.Redirect(w, r, "/compare?"+url.Values{"from": {from}, "to": {to}}.Encode(), http.StatusSeeOther)
}

// runDays runs n business days and returns the openings before and after
// them.
func (s *Server) runDays(n int) (from, to string, err error) {
	before, err := s.bank.History()
	if err != nil {
		return "", "", err
	}
	if err := s.bank.RunDays(n, func(bank.Day) {}); err != nil {
		return "", "", err
	}
	after, err := s.bank.History()
	return before.Date.String(), after.Date.String(), err
}

func (s *Server) fail(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, ErrNotFound) {
		status = http.StatusNotFound
	}
	s.render(w, status, "error", page{View: err.Error()})
}

// refuse answers a form that asked for something the console does not do.
func (s *Server) refuse(w http.ResponseWriter, err error) {
	s.render(w, http.StatusBadRequest, "error", page{View: err.Error()})
}

func (s *Server) render(w http.ResponseWriter, status int, name string, p page) {
	var b bytes.Buffer
	if err := s.pages[name].ExecuteTemplate(&b, "base", p); err != nil {
		log.Printf("console: %s: %v", name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = b.WriteTo(w)
}
