package explore

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/diagram"
)

//go:embed templates static
var files embed.FS

// errNotFound is for a page of something the system does not have.
var errNotFound = errors.New("not found")

// Server serves the explorer for one system.
type Server struct {
	title string
	index *index
	pages map[string]*template.Template
}

// New makes an explorer of the system, titled with its name.
func New(system *blueprint.System, title string) (*Server, error) {
	s := &Server{title: title, index: newIndex(system), pages: map[string]*template.Template{}}
	funcs := template.FuncMap{
		"program": programLink, "paragraph": paragraphLink, "field": fieldLink, "copybook": copybookLink,
		"job": jobLink, "entity": entityLink, "entityField": entityFieldLink, "step": stepLink, "dataset": datasetLink, "source": sourceLink,
		"sourceText": sourceText, "base": datasetBase, "join": strings.Join,
		"add":   func(a, b int) int { return a + b },
		"hop":   hopID,
		"trace": traceLink, "forward": forwardLink,
		"stepRef": func(job, step string) blueprint.StepRef { return blueprint.StepRef{Job: job, Step: step} },
		"lineage": func(program string, fields map[string]bool, rows []blueprint.Lineage) lineageTable {
			return lineageTable{Program: program, Fields: fields, Rows: rows}
		},
		"layoutOf":   func(record string, f blueprint.Field) layoutRow { return layoutRow{Record: record, Field: f} },
		"layoutRows": blueprint.LayoutRows,
		"notes":      func(f *blueprint.Field) string { return blueprint.FieldNotes(*f) },
		"decisions": func(program string, rows []blueprint.Decision) decisionTable {
			return decisionTable{Program: program, Rows: rows}
		},
	}
	for _, page := range []string{"index", "search", "job", "step", "dataset", "program", "paragraph", "field", "copybook",
		"entities", "entity", "entityfield", "ownership", "boundary", "schedule", "source", "trace", "error"} {
		t, err := template.New("").Funcs(funcs).ParseFS(files, "templates/base.html", "templates/"+page+".html")
		if err != nil {
			return nil, err
		}
		s.pages[page] = t
	}
	return s, nil
}

// Handler routes the explorer's pages.
func (s *Server) Handler() http.Handler {
	static, _ := fs.Sub(files, "static")
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /{$}", s.view("index", s.home))
	mux.HandleFunc("GET /search", s.view("search", s.search))
	mux.HandleFunc("GET /jobs/{job}", s.view("job", s.job))
	mux.HandleFunc("GET /jobs/{job}/steps/{step}", s.view("step", s.step))
	mux.HandleFunc("GET /datasets/{name}", s.view("dataset", s.dataset))
	mux.HandleFunc("GET /programs/{program}", s.view("program", s.program))
	mux.HandleFunc("GET /programs/{program}/paragraphs/{paragraph}", s.view("paragraph", s.paragraph))
	mux.HandleFunc("GET /programs/{program}/fields/{field}", s.view("field", s.field))
	mux.HandleFunc("GET /copybooks/{copybook}", s.view("copybook", s.copybook))
	mux.HandleFunc("GET /entities", s.view("entities", s.entities))
	mux.HandleFunc("GET /ownership", s.view("ownership", s.ownership))
	mux.HandleFunc("GET /boundary", s.view("boundary", s.boundary))
	mux.HandleFunc("GET /schedule", s.view("schedule", s.schedule))
	mux.HandleFunc("GET /entities/{record}", s.view("entity", s.entity))
	mux.HandleFunc("GET /entities/{record}/fields/{field}", s.view("entityfield", s.entityField))
	mux.HandleFunc("GET /source/{path...}", s.view("source", s.source))
	mux.HandleFunc("GET /trace", s.view("trace", s.trace))
	return mux
}

// page is what every template gets: the system's title, and the page's
// own view.
type page struct {
	Title string
	View  any
}

func (s *Server) view(name string, read func(*http.Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, err := read(r)
		status := http.StatusOK
		switch {
		case errors.Is(err, errNotFound):
			name, status, v = "error", http.StatusNotFound, err.Error()
		case err != nil:
			name, status, v = "error", http.StatusInternalServerError, err.Error()
		}
		var b bytes.Buffer
		if err := s.pages[name].ExecuteTemplate(&b, "base", page{Title: s.title, View: v}); err != nil {
			log.Printf("explore: %s: %v", name, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		_, _ = b.WriteTo(w)
	}
}

// lineageTable is the statements that move values, for the lineage
// template: those of a program whose fields link to their pages.
type lineageTable struct {
	Program string
	Fields  map[string]bool
	Rows    []blueprint.Lineage
}

// layoutRow is a field of a copybook's record, for the layout template.
type layoutRow struct {
	Record string
	Field  blueprint.Field
}

// decisionTable is a program's decisions on a field, for the decisions
// template.
type decisionTable struct {
	Program string
	Rows    []blueprint.Decision
}

// HomeView lists the system's parts.
type HomeView struct {
	System *blueprint.System
}

func (s *Server) home(*http.Request) (any, error) {
	return HomeView{System: s.index.system}, nil
}

// SearchView is what a search found.
type SearchView struct {
	Query string
	Names []Name
}

func (s *Server) search(r *http.Request) (any, error) {
	q := r.FormValue("q")
	return SearchView{Query: q, Names: s.index.search(q)}, nil
}

// JobView is a job and its network.
type JobView struct {
	Job     *blueprint.Job
	Network template.HTML
}

func (s *Server) job(r *http.Request) (any, error) {
	j, ok := s.index.jobs[r.PathValue("job")]
	if !ok {
		return nil, fmt.Errorf("job %s: %w", r.PathValue("job"), errNotFound)
	}
	network, err := svg(blueprint.JobGraph(*j, blueprint.JobNotes{
		StepLink:    func(step blueprint.Step) string { return stepLink(blueprint.StepRef{Job: j.Name, Step: step.Name}) },
		DatasetLink: func(ds string) string { return datasetLink(datasetBase(ds)) },
	}))
	return JobView{Job: j, Network: network}, err
}

// StepView is a step of a job, and whether its program was analyzed.
type StepView struct {
	Job      *blueprint.Job
	Step     blueprint.Step
	Analyzed bool
}

func (s *Server) step(r *http.Request) (any, error) {
	j, ok := s.index.jobs[r.PathValue("job")]
	if !ok {
		return nil, fmt.Errorf("job %s: %w", r.PathValue("job"), errNotFound)
	}
	i := slices.IndexFunc(j.Steps, func(step blueprint.Step) bool { return step.Name == r.PathValue("step") })
	if i < 0 {
		return nil, fmt.Errorf("step %s of job %s: %w", r.PathValue("step"), j.Name, errNotFound)
	}
	_, analyzed := s.index.programs[j.Steps[i].Program]
	return StepView{Job: j, Step: j.Steps[i], Analyzed: analyzed}, nil
}

// DatasetView is a dataset, its kind across the system, the copybooks that
// describe its records, and every job step that uses it.
type DatasetView struct {
	Name      string
	Dataset   *blueprint.SystemDataset
	Copybooks []string
	Uses      []DatasetUse
}

func (s *Server) dataset(r *http.Request) (any, error) {
	name := r.PathValue("name")
	uses, ok := s.index.datasets[name]
	if !ok {
		return nil, fmt.Errorf("dataset %s: %w", name, errNotFound)
	}
	v := DatasetView{Name: name, Dataset: s.index.systemDatasets[name], Uses: uses}
	for _, u := range uses {
		if u.Copybook != "" && !slices.Contains(v.Copybooks, u.Copybook) {
			v.Copybooks = append(v.Copybooks, u.Copybook)
		}
	}
	return v, nil
}

// ProgramView is a program, its flow, and the job steps that run it.
type ProgramView struct {
	Program *blueprint.Program
	Flow    template.HTML
	Runs    []blueprint.StepRef
}

func (s *Server) program(r *http.Request) (any, error) {
	p, err := s.lookupProgram(r)
	if err != nil {
		return nil, err
	}
	flow, err := svg(blueprint.FlowGraph(*p, blueprint.FlowNotes{
		Link: func(para blueprint.Paragraph) string { return paragraphLink(p.Name, para.Name) },
	}))
	return ProgramView{Program: p, Flow: flow, Runs: s.index.runs[p.Name]}, err
}

func (s *Server) lookupProgram(r *http.Request) (*blueprint.Program, error) {
	p, ok := s.index.programs[r.PathValue("program")]
	if !ok {
		return nil, fmt.Errorf("program %s: %w", r.PathValue("program"), errNotFound)
	}
	return p, nil
}

// ParagraphView is a paragraph of a program, and the values its statements
// move.
type ParagraphView struct {
	Program   *blueprint.Program
	Paragraph blueprint.Paragraph
	Lineage   []blueprint.Lineage
	Fields    map[string]bool
}

func (s *Server) paragraph(r *http.Request) (any, error) {
	p, err := s.lookupProgram(r)
	if err != nil {
		return nil, err
	}
	i := slices.IndexFunc(p.Paragraphs, func(para blueprint.Paragraph) bool { return para.Name == r.PathValue("paragraph") })
	if i < 0 {
		return nil, fmt.Errorf("paragraph %s of %s: %w", r.PathValue("paragraph"), p.Name, errNotFound)
	}
	v := ParagraphView{Program: p, Paragraph: p.Paragraphs[i], Fields: fieldNames(p)}
	for _, e := range p.Lineage {
		if e.Paragraph == v.Paragraph.Name {
			v.Lineage = append(v.Lineage, e)
		}
	}
	return v, nil
}

// FieldView is a data item of a program: where it is, the paragraphs that
// read and write it, the statements that move values into it and out of
// it, and the decisions on it.
type FieldView struct {
	Program   *blueprint.Program
	Record    blueprint.Record
	Path      []blueprint.Field // the groups it is in, from its record down
	Field     blueprint.Field
	ReadBy    []string
	WrittenBy []string
	Into      []blueprint.Lineage
	OutOf     []blueprint.Lineage
	Decisions []blueprint.Decision
	// Entity is the field of an entity it is, if its record is one.
	Entity *EntityFieldRef
	Fields map[string]bool
}

func (s *Server) field(r *http.Request) (any, error) {
	p, err := s.lookupProgram(r)
	if err != nil {
		return nil, err
	}
	name := r.PathValue("field")
	v := FieldView{Program: p, Fields: fieldNames(p)}
	found := false
	for _, rec := range p.Records {
		if path, f, ok := findField(rec.Field, name); ok {
			v.Record, v.Path, v.Field, found = rec, path, f, true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("field %s of %s: %w", name, p.Name, errNotFound)
	}
	for _, para := range p.Paragraphs {
		if slices.Contains(para.Reads, name) {
			v.ReadBy = append(v.ReadBy, para.Name)
		}
		if slices.Contains(para.Writes, name) {
			v.WrittenBy = append(v.WrittenBy, para.Name)
		}
	}
	v.Decisions = p.Decisions(name)
	if ref, ok := s.index.entityFields[[2]string{p.Name, name}]; ok {
		v.Entity = &ref
	}
	for _, e := range p.Lineage {
		if e.To == name && e.ToItem() {
			v.Into = append(v.Into, e)
		}
		if e.From == name && e.FromItem() {
			v.OutOf = append(v.OutOf, e)
		}
	}
	return v, nil
}

// findField finds a field by name in a record, with the groups it is in.
func findField(f blueprint.Field, name string) ([]blueprint.Field, blueprint.Field, bool) {
	if f.Name == name {
		return nil, f, true
	}
	for _, child := range f.Fields {
		if path, found, ok := findField(child, name); ok {
			return append([]blueprint.Field{f}, path...), found, true
		}
	}
	return nil, blueprint.Field{}, false
}

func fieldNames(p *blueprint.Program) map[string]bool {
	names := map[string]bool{}
	for _, r := range p.Records {
		walkFields(r.Field, func(f blueprint.Field) { names[f.Name] = true })
	}
	return names
}

func (s *Server) copybook(r *http.Request) (any, error) {
	c, ok := s.index.copybooks[r.PathValue("copybook")]
	if !ok {
		return nil, fmt.Errorf("copybook %s: %w", r.PathValue("copybook"), errNotFound)
	}
	return c, nil
}

func (s *Server) entities(*http.Request) (any, error) {
	return s.index.system.Entities, nil
}

// OwnershipView is the entities kept in datasets, with the programs that
// write them and those that only read them, drawn.
type OwnershipView struct {
	Entities []blueprint.Entity
	Diagram  template.HTML
}

func (s *Server) ownership(*http.Request) (any, error) {
	v := OwnershipView{}
	for _, e := range s.index.system.Entities {
		if len(e.WrittenBy)+len(e.ReadBy) > 0 {
			v.Entities = append(v.Entities, e)
		}
	}
	var err error
	v.Diagram, err = svg(blueprint.OwnershipGraph(v.Entities, blueprint.OwnershipNotes{EntityLink: entityLink, ProgramLink: programLink}))
	return v, err
}

// ScheduleView is the order of the jobs, drawn, and the cycles in it.
type ScheduleView struct {
	Order   []blueprint.JobOrder
	Cycles  [][]string
	Diagram template.HTML
}

func (s *Server) schedule(*http.Request) (any, error) {
	system := s.index.system
	v := ScheduleView{Order: system.Schedule, Cycles: blueprint.JobCycles(system.Jobs, system.Schedule)}
	var err error
	v.Diagram, err = svg(blueprint.ScheduleGraph(system.Jobs, system.Schedule, jobLink))
	return v, err
}

// EntityView is an entity: the datasets it describes, with their kinds,
// the programs that use it, its fields with the programs that read, write
// and decide on each, and when the entities hold dates in more than one
// layout, the layouts, with its fields in each and those of the others.
type EntityView struct {
	Entity   *blueprint.Entity
	Datasets []*blueprint.SystemDataset
	Programs []*blueprint.Program
	Fields   []EntityFieldUses
	Machines []StateView
	Dates    []DateLayoutView
}

// DateLayoutView is a layout the entities hold dates in: the entity's own
// fields in it, each with the dates of other layouts its values come
// from, and the other entities' fields.
type DateLayoutView struct {
	Layout    string
	Here      []DateFieldView
	Elsewhere []blueprint.FieldRef
}

// DateFieldView is a field of the entity in a date layout, and the date
// fields of other layouts its values come from, with their layouts.
type DateFieldView struct {
	Field string
	From  []DateSourceView
}

// DateSourceView is a date field its values come from, and its layout.
type DateSourceView struct {
	blueprint.FieldRef
	Layout string
}

// StateView is a code field's state machine, drawn.
type StateView struct {
	Record  string
	Field   blueprint.EntityField
	Diagram template.HTML
}

func stateView(record string, f blueprint.EntityField) (StateView, error) {
	v := StateView{Record: record, Field: f}
	if len(f.States.Transitions) == 0 {
		return v, nil
	}
	var err error
	v.Diagram, err = svg(blueprint.StateGraph(*f.States))
	return v, err
}

// EntityFieldUses is a field of an entity and the programs that use it,
// by what they do with it.
type EntityFieldUses struct {
	Field                      blueprint.EntityField
	Readers, Writers, Deciders []blueprint.FieldUse
}

func (s *Server) entity(r *http.Request) (any, error) {
	e, ok := s.index.entities[r.PathValue("record")]
	if !ok {
		return nil, fmt.Errorf("entity %s: %w", r.PathValue("record"), errNotFound)
	}
	v := EntityView{Entity: e}
	for _, name := range e.Datasets {
		v.Datasets = append(v.Datasets, s.index.systemDatasets[name])
	}
	for _, name := range e.Programs {
		v.Programs = append(v.Programs, s.index.programs[name])
	}
	for _, f := range e.Fields {
		row := EntityFieldUses{Field: f}
		for _, u := range f.Uses {
			if len(u.ReadBy) > 0 {
				row.Readers = append(row.Readers, u)
			}
			if len(u.WrittenBy) > 0 {
				row.Writers = append(row.Writers, u)
			}
			if len(u.Decisions) > 0 {
				row.Deciders = append(row.Deciders, u)
			}
		}
		v.Fields = append(v.Fields, row)
		if f.States != nil {
			m, err := stateView(e.Record, f)
			if err != nil {
				return nil, err
			}
			v.Machines = append(v.Machines, m)
		}
	}
	v.Dates = s.dateLayouts(e)
	return v, nil
}

// dateLayouts lays out the entities' date layouts for an entity's page:
// none unless there is more than one and the entity holds a date.
func (s *Server) dateLayouts(e *blueprint.Entity) []DateLayoutView {
	layouts := blueprint.DateLayouts(s.index.system.Entities)
	if len(layouts) < 2 || !slices.ContainsFunc(layouts, func(l blueprint.DateLayout) bool {
		return slices.ContainsFunc(l.Fields, func(f blueprint.FieldRef) bool { return f.Record == e.Record })
	}) {
		return nil
	}
	views := make([]DateLayoutView, len(layouts))
	for i, l := range layouts {
		views[i].Layout = l.Layout
		for _, ref := range l.Fields {
			if ref.Record != e.Record {
				views[i].Elsewhere = append(views[i].Elsewhere, ref)
				continue
			}
			f, _ := e.Field(ref.Field)
			view := DateFieldView{Field: f.Name}
			for _, from := range f.DatesFrom {
				source, _ := s.index.entities[from.Record].Field(from.Field)
				view.From = append(view.From, DateSourceView{FieldRef: from, Layout: source.Layout})
			}
			views[i].Here = append(views[i].Here, view)
		}
	}
	return views
}

// EntityFieldView is a field of an entity across every program that uses
// the entity: its layout in the copybook, and each program's use of it.
type EntityFieldView struct {
	Entity  *blueprint.Entity
	Layout  blueprint.Field
	Field   blueprint.EntityField
	Used    []blueprint.FieldUse
	Unused  []string
	Machine *StateView
}

func (s *Server) entityField(r *http.Request) (any, error) {
	e, ok := s.index.entities[r.PathValue("record")]
	if !ok {
		return nil, fmt.Errorf("entity %s: %w", r.PathValue("record"), errNotFound)
	}
	f, ok := e.Field(r.PathValue("field"))
	if !ok {
		return nil, fmt.Errorf("field %s of %s: %w", r.PathValue("field"), e.Record, errNotFound)
	}
	v := EntityFieldView{Entity: e, Field: f}
	if f.States != nil {
		m, err := stateView(e.Record, f)
		if err != nil {
			return nil, err
		}
		v.Machine = &m
	}
	c := s.index.copybooks[e.Copybook]
	for _, rec := range c.Records {
		if rec.Name == e.Record {
			_, v.Layout, _ = findField(rec.Field, f.Name)
		}
	}
	for _, u := range f.Uses {
		if len(u.ReadBy)+len(u.WrittenBy)+len(u.Decisions) > 0 {
			v.Used = append(v.Used, u)
		} else {
			v.Unused = append(v.Unused, u.Program)
		}
	}
	return v, nil
}

// SourceView is a source file, by line.
type SourceView struct {
	Path  string
	Lines []string
}

// source shows a source file the analysis read, and no other.
func (s *Server) source(r *http.Request) (any, error) {
	path := r.PathValue("path")
	if !s.index.sources[path] {
		return nil, fmt.Errorf("source %s: %w", path, errNotFound)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	v := SourceView{Path: path}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		v.Lines = append(v.Lines, scanner.Text())
	}
	return v, scanner.Err()
}

// sourceText writes a source reference as the file's name and the line.
func sourceText(source string) string {
	return source[strings.LastIndex(source, "/")+1:]
}

func svg(g diagram.Graph) (template.HTML, error) {
	b, err := g.SVG()
	return template.HTML(b), err //nolint:gosec // drawn by the survey, its labels escaped
}
