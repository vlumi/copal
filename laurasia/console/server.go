// Package console is the card processor's console: web pages for its
// operators that show its cards, each day's night and billing cycle with
// their steps, and the files each day read and wrote, read through the
// record layouts computed from the copybooks. Its buttons open the
// processor and run days by calling the processor's own operations, the
// same ones the make targets run; the console itself never writes the
// processor's files.
package console

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/vlumi/copal/laurasia/processor"
	"github.com/vlumi/copal/laurasia/records"
)

//go:embed templates static
var files embed.FS

// Server serves the console for one processor.
type Server struct {
	processor processor.Processor
	layouts   *records.Layouts
	pages     map[string]*template.Template
	// busy is locked while the processor opens or runs, and read-locked
	// while a page reads its files.
	busy sync.RWMutex
}

// New makes a console for the processor.
func New(p processor.Processor) (*Server, error) {
	layouts, err := records.Load(p.Copybooks)
	if err != nil {
		return nil, err
	}
	s := &Server{processor: p, layouts: layouts, pages: map[string]*template.Template{}}
	funcs := template.FuncMap{"money": money, "count": count, "size": size}
	for _, page := range []string{"overview", "noprocessor", "cards", "card", "day", "error"} {
		t, err := template.New("").Funcs(funcs).ParseFS(files, "templates/base.html", "templates/"+page+".html")
		if err != nil {
			return nil, err
		}
		s.pages[page] = t
	}
	return s, nil
}

// Handler routes the console's pages and actions. Its forms open and run
// the processor, so it refuses them from pages of other origins.
func (s *Server) Handler() http.Handler {
	static, _ := fs.Sub(files, "static")
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /{$}", s.view("overview", s.overview))
	mux.HandleFunc("GET /cards", s.view("cards", s.cards))
	mux.HandleFunc("GET /cards/{no}", s.view("card", s.card))
	mux.HandleFunc("GET /days/{date}", s.view("day", s.day))
	mux.HandleFunc("GET /files/{path...}", s.file)
	mux.HandleFunc("POST /open", s.open)
	mux.HandleFunc("POST /run", s.run)
	return http.NewCrossOriginProtection().Handler(mux)
}

// page is what every template gets: the processor's history for the
// header, and the page's own view.
type page struct {
	History processor.History
	View    any
}

type viewFunc func(processor.History, *http.Request) (any, error)

// view serves a page that reads the processor. Without a processor, every
// page offers to open one.
func (s *Server) view(name string, read viewFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.busy.RLock()
		defer s.busy.RUnlock()
		h, err := s.processor.History()
		if errors.Is(err, processor.ErrNotOpen) {
			s.render(w, http.StatusOK, "noprocessor", page{})
			return
		}
		if err != nil {
			s.fail(w, err)
			return
		}
		v, err := read(h, r)
		if errors.Is(err, errNotFound) {
			s.render(w, http.StatusNotFound, "error", page{History: h, View: err.Error()})
			return
		}
		if err != nil {
			s.fail(w, err)
			return
		}
		s.render(w, http.StatusOK, name, page{History: h, View: v})
	}
}

// file serves a file of the processor's data as text.
func (s *Server) file(w http.ResponseWriter, r *http.Request) {
	s.busy.RLock()
	defer s.busy.RUnlock()
	content, err := s.processor.ReadFile(r.PathValue("path"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write(content)
}

// open opens a new processor from the form's seed and number of cards.
func (s *Server) open(w http.ResponseWriter, r *http.Request) {
	seed, err := strconv.ParseUint(r.FormValue("seed"), 10, 64)
	cards, cardsErr := strconv.Atoi(r.FormValue("cards"))
	if err != nil || cardsErr != nil || cards < 10 || cards > 100000 {
		s.refuse(w, "a seed, and between 10 and 100,000 cards")
		return
	}
	s.busy.Lock()
	err = s.processor.Open(seed, cards)
	s.busy.Unlock()
	if err != nil {
		s.fail(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// run runs the form's number of days and shows the overview.
func (s *Server) run(w http.ResponseWriter, r *http.Request) {
	days, err := strconv.Atoi(r.FormValue("days"))
	if err != nil || days < 1 || days > 366 {
		s.refuse(w, "between 1 and 366 days")
		return
	}
	s.busy.Lock()
	err = s.processor.RunDays(days, func(processor.Day) {})
	s.busy.Unlock()
	if err != nil {
		s.fail(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) refuse(w http.ResponseWriter, want string) {
	s.render(w, http.StatusBadRequest, "error", page{View: fmt.Sprintf("The console takes %s.", want)})
}

func (s *Server) fail(w http.ResponseWriter, err error) {
	log.Print(err)
	s.render(w, http.StatusInternalServerError, "error", page{View: err.Error()})
}

func (s *Server) render(w http.ResponseWriter, status int, name string, p page) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := s.pages[name].ExecuteTemplate(w, "base", p); err != nil {
		log.Print(err)
	}
}
