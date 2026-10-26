// Package console is the life insurer's console: web pages for its
// operators that show its policies, each business day's batch files with
// their steps, and the files each day read and wrote, read through the
// record layouts computed from the copybooks and shown in UTF-8 from the
// Shift_JIS they are kept in. Its buttons open the insurer and run days by
// calling the insurer's own operations, the same ones the make targets
// run; the console itself never writes the insurer's files.
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

	"github.com/vlumi/copal/gondwana/insurer"
	"github.com/vlumi/copal/gondwana/records"
)

//go:embed templates static
var files embed.FS

// Server serves the console for one insurer.
type Server struct {
	insurer insurer.Insurer
	layouts *records.Layouts
	pages   map[string]*template.Template
	// busy is locked while the insurer opens or runs, and read-locked
	// while a page reads its files.
	busy sync.RWMutex
}

// New makes a console for the insurer.
func New(in insurer.Insurer) (*Server, error) {
	layouts, err := records.Load(in.Copybooks)
	if err != nil {
		return nil, err
	}
	s := &Server{insurer: in, layouts: layouts, pages: map[string]*template.Template{}}
	funcs := template.FuncMap{"yen": yen, "count": count, "size": size}
	for _, page := range []string{"overview", "noinsurer", "policies", "policy", "day", "error"} {
		t, err := template.New("").Funcs(funcs).ParseFS(files, "templates/base.html", "templates/"+page+".html")
		if err != nil {
			return nil, err
		}
		s.pages[page] = t
	}
	return s, nil
}

// Handler routes the console's pages and actions. Its forms open and run
// the insurer, so it refuses them from pages of other origins.
func (s *Server) Handler() http.Handler {
	static, _ := fs.Sub(files, "static")
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /{$}", s.view("overview", s.overview))
	mux.HandleFunc("GET /policies", s.view("policies", s.policies))
	mux.HandleFunc("GET /policies/{no}", s.view("policy", s.policy))
	mux.HandleFunc("GET /days/{date}", s.view("day", s.day))
	mux.HandleFunc("GET /files/{path...}", s.file)
	mux.HandleFunc("POST /open", s.open)
	mux.HandleFunc("POST /run", s.run)
	return http.NewCrossOriginProtection().Handler(mux)
}

// page is what every template gets: the insurer's history for the
// header, and the page's own view.
type page struct {
	History insurer.History
	View    any
}

type viewFunc func(insurer.History, *http.Request) (any, error)

// view serves a page that reads the insurer. Without an insurer, every
// page offers to open one.
func (s *Server) view(name string, read viewFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.busy.RLock()
		defer s.busy.RUnlock()
		h, err := s.insurer.History()
		if errors.Is(err, insurer.ErrNotOpen) {
			s.render(w, http.StatusOK, "noinsurer", page{})
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

// file serves a file of the insurer's data as text, in UTF-8.
func (s *Server) file(w http.ResponseWriter, r *http.Request) {
	s.busy.RLock()
	defer s.busy.RUnlock()
	content, err := s.insurer.ReadFile(r.PathValue("path"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(s.text(r.PathValue("path"), content)))
}

// open opens a new insurer from the form's seed and number of policies.
func (s *Server) open(w http.ResponseWriter, r *http.Request) {
	seed, err := strconv.ParseUint(r.FormValue("seed"), 10, 64)
	policies, policiesErr := strconv.Atoi(r.FormValue("policies"))
	if err != nil || policiesErr != nil || policies < 10 || policies > 100000 {
		s.refuse(w, "a seed, and between 10 and 100,000 policies")
		return
	}
	s.busy.Lock()
	err = s.insurer.Open(seed, policies)
	s.busy.Unlock()
	if err != nil {
		s.fail(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// run runs the form's number of business days and shows the overview.
func (s *Server) run(w http.ResponseWriter, r *http.Request) {
	days, err := strconv.Atoi(r.FormValue("days"))
	if err != nil || days < 1 || days > 366 {
		s.refuse(w, "between 1 and 366 days")
		return
	}
	s.busy.Lock()
	err = s.insurer.RunDays(days, func(insurer.Day) {})
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
