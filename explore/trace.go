package explore

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/vlumi/copal/blueprint"
	"github.com/vlumi/copal/dataflow"
	"github.com/vlumi/copal/lineage"
)

// TraceView is a field's values followed back across the system: every
// hop of the trail, nearest first, the job networks and the program flows
// with the trail on them, and where the values start.
type TraceView struct {
	Start    lineage.Node
	Hops     []TraceHop
	Jobs     []Diagram
	Programs []Diagram
	External []string
}

// TraceHop is a field of a program on the trail: the statements that move
// values into it, the literals and initial values it starts from, and the
// files it reads from.
type TraceHop struct {
	ID         string
	Node       lineage.Node
	Statements []blueprint.Lineage
	Fields     map[string]bool
	Starts     []string
	Sources    []lineage.Source
}

// Diagram is a job network or a program flow, by name.
type Diagram struct {
	Name string
	SVG  template.HTML
}

func (s *Server) trace(r *http.Request) (any, error) {
	start := lineage.Node{Job: r.FormValue("job"), Step: r.FormValue("step"), Program: r.FormValue("program"), Field: r.FormValue("field")}
	p, ok := s.index.programs[start.Program]
	if !ok || !fieldNames(p)[start.Field] {
		return nil, fmt.Errorf("field %s of %s: %w", start.Field, start.Program, errNotFound)
	}
	trail := lineage.Back(s.index.system, start)
	v := TraceView{Start: start}
	if len(trail.Hops) > 0 {
		v.Start = trail.Hops[0].Node
	}
	paths := map[string]*blueprint.JobPath{}
	paragraphs := map[string]map[string]bool{}
	for _, h := range trail.Hops {
		q := s.index.programs[h.Node.Program]
		hop := TraceHop{ID: hopID(h.Node), Node: h.Node, Fields: fieldNames(q), Sources: h.Sources}
		for _, e := range h.Trace.Edges {
			hop.Statements = append(hop.Statements, statement(e))
			if paragraphs[q.Name] == nil {
				paragraphs[q.Name] = map[string]bool{}
			}
			paragraphs[q.Name][e.Paragraph] = true
		}
		for _, o := range h.Trace.Origins {
			if o.File == "" {
				hop.Starts = append(hop.Starts, o.String())
			}
		}
		v.Hops = append(v.Hops, hop)
		s.markPath(paths, h)
		for _, src := range h.Sources {
			for _, ds := range src.External {
				if !slices.Contains(v.External, ds) {
					v.External = append(v.External, ds)
				}
			}
		}
	}
	var err error
	if v.Jobs, err = s.jobDiagrams(paths); err != nil {
		return nil, err
	}
	v.Programs, err = s.flowDiagrams(paragraphs)
	return v, err
}

// markPath adds a hop to the paths through the jobs: its step, the steps
// that wrote what it reads and the sorts that passed it on, and the
// datasets between them.
func (s *Server) markPath(paths map[string]*blueprint.JobPath, h lineage.Hop) {
	mark := func(n lineage.Node, datasets []string) {
		if n.Job == "" {
			return
		}
		path := paths[n.Job]
		if path == nil {
			path = &blueprint.JobPath{Steps: map[string]bool{}, Datasets: map[[2]string]bool{}}
			paths[n.Job] = path
		}
		path.Steps[n.Step] = true
		job := s.index.jobs[n.Job]
		i := slices.IndexFunc(job.Steps, func(step blueprint.Step) bool { return step.Name == n.Step })
		if i < 0 {
			return
		}
		for _, dd := range job.Steps[i].DDs {
			for _, ds := range dd.Datasets {
				if slices.ContainsFunc(datasets, func(d string) bool { return datasetBase(d) == datasetBase(ds) }) {
					path.Datasets[[2]string{n.Step, ds}] = true
				}
			}
		}
	}
	mark(h.Node, nil)
	for _, src := range h.Sources {
		mark(h.Node, src.Datasets)
		for _, n := range append(slices.Clone(src.Writers), src.Sorts...) {
			mark(n, src.Datasets)
		}
	}
}

func (s *Server) jobDiagrams(paths map[string]*blueprint.JobPath) ([]Diagram, error) {
	var diagrams []Diagram
	for _, j := range s.index.system.Jobs {
		path, ok := paths[j.Name]
		if !ok {
			continue
		}
		svg, err := svg(blueprint.JobGraph(j, blueprint.JobNotes{
			StepLink:    func(step blueprint.Step) string { return stepLink(StepRef{Job: j.Name, Step: step.Name}) },
			DatasetLink: func(ds string) string { return datasetLink(datasetBase(ds)) },
			Path:        path,
		}))
		if err != nil {
			return nil, err
		}
		diagrams = append(diagrams, Diagram{Name: j.Name, SVG: svg})
	}
	return diagrams, nil
}

func (s *Server) flowDiagrams(paragraphs map[string]map[string]bool) ([]Diagram, error) {
	var diagrams []Diagram
	for _, p := range s.index.system.Programs {
		on, ok := paragraphs[p.Name]
		if !ok {
			continue
		}
		svg, err := svg(blueprint.FlowGraph(p, blueprint.FlowNotes{
			Link:   func(para blueprint.Paragraph) string { return paragraphLink(p.Name, para.Name) },
			OnPath: func(para blueprint.Paragraph) bool { return on[para.Name] },
		}))
		if err != nil {
			return nil, err
		}
		diagrams = append(diagrams, Diagram{Name: p.Name, SVG: svg})
	}
	return diagrams, nil
}

func statement(e dataflow.Edge) blueprint.Lineage {
	return blueprint.Lineage{From: e.From, To: e.To, FromFile: e.FromFile, FromLiteral: e.FromLiteral, ToFile: e.ToFile,
		Paragraph: e.Paragraph, Verb: e.Verb, Source: e.Pos.String()}
}

// hopID is the anchor of a hop on the trace page.
func hopID(n lineage.Node) string {
	return "hop-" + strings.NewReplacer(" ", "-", "&", "", "(", "", ")", "", "+", "").Replace(
		strings.Join([]string{n.Job, n.Step, n.Program, n.Field}, "-"))
}

// traceLink is the trace page of a program's field.
func traceLink(program, field string) string {
	return "/trace?" + url.Values{"program": {program}, "field": {field}}.Encode()
}
