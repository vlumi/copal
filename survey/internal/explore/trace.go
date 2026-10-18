package explore

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/dataflow"
	"github.com/vlumi/copal/survey/internal/lineage"
)

// TraceView is a field's values followed back across the system, or
// forward: every hop of the trail, nearest first, the job networks and the
// program flows with the trail on them, and where the values start, or
// where they end.
type TraceView struct {
	Forward  bool
	Start    lineage.Node
	Hops     []TraceHop
	Jobs     []Diagram
	Programs []Diagram
	External []string // back, the datasets no job writes
	Outputs  []string // forward, the datasets no step reads
	Reports  []string // forward, the DDs of the job's output
	Capped   bool
}

// TraceHop is a field of a program on the trail: the statements that move
// its values, and going back, the literals and initial values it starts
// from and the files it reads from; going forward, the items its values
// reach, the files it writes them to and the DISPLAYs that show them.
type TraceHop struct {
	ID         string
	Node       lineage.Node
	Statements []blueprint.Lineage
	Fields     map[string]bool
	Starts     []string
	Sources    []lineage.Source
	Reaches    []string
	Targets    []lineage.Target
	Displays   []blueprint.Display
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
	v := TraceView{Start: start, Forward: r.FormValue("direction") == "forward"}
	var trail lineage.Trail
	if v.Forward {
		trail = lineage.Forward(s.index.system, start)
	} else {
		trail = lineage.Back(s.index.system, start)
	}
	v.Capped = trail.Capped
	if len(trail.Hops) > 0 {
		v.Start = trail.Hops[0].Node
	}
	paths := map[string]*blueprint.JobPath{}
	paragraphs := map[string]map[string]bool{}
	for _, h := range trail.Hops {
		q := s.index.programs[h.Node.Program]
		hop := TraceHop{ID: hopID(h.Node), Node: h.Node, Fields: fieldNames(q), Sources: h.Sources,
			Reaches: h.Effect.Items, Targets: h.Targets, Displays: h.Displays}
		for _, e := range append(slices.Clone(h.Trace.Edges), h.Effect.Edges...) {
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
		for _, src := range h.Sources {
			s.markPath(paths, h.Node, src.Datasets, append(slices.Clone(src.Writers), src.Sorts...))
			addAll(&v.External, src.External)
		}
		for _, tgt := range h.Targets {
			s.markPath(paths, h.Node, tgt.Datasets, append(slices.Clone(tgt.Readers), tgt.Sorts...))
			addAll(&v.Outputs, tgt.Outputs)
			if tgt.Report != "" {
				addAll(&v.Reports, []string{tgt.Report})
			}
		}
		s.markPath(paths, h.Node, nil, nil)
	}
	var err error
	if v.Jobs, err = s.jobDiagrams(paths); err != nil {
		return nil, err
	}
	v.Programs, err = s.flowDiagrams(paragraphs)
	return v, err
}

func addAll(list *[]string, names []string) {
	for _, name := range names {
		if !slices.Contains(*list, name) {
			*list = append(*list, name)
		}
	}
}

// markPath adds a hop's step to the paths through the jobs, and with it,
// the steps at the other end of the datasets between them: those that
// wrote what it reads, or read what it writes, and the sorts that passed
// the records on.
func (s *Server) markPath(paths map[string]*blueprint.JobPath, at lineage.Node, datasets []string, others []lineage.Node) {
	mark := func(n lineage.Node) {
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
	mark(at)
	for _, n := range others {
		mark(n)
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
			StepLink:    func(step blueprint.Step) string { return stepLink(blueprint.StepRef{Job: j.Name, Step: step.Name}) },
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
	return blueprint.NewLineage(e)
}

// hopID is the anchor of a hop on the trace page.
func hopID(n lineage.Node) string {
	parts := []string{n.Job, n.Step, n.Program, n.Field}
	if n.Length > 0 {
		parts = append(parts, strconv.Itoa(n.Start), strconv.Itoa(n.Length))
	}
	return "hop-" + strings.NewReplacer(" ", "-", "&", "", "(", "", ")", "", "+", "").Replace(strings.Join(parts, "-"))
}

// traceLink is the trace page of a program's field, back to where its
// values come from.
func traceLink(program, field string) string {
	return "/trace?" + url.Values{"program": {program}, "field": {field}}.Encode()
}

// forwardLink is the trace page of a program's field, on to where its
// values go.
func forwardLink(program, field string) string {
	return "/trace?" + url.Values{"direction": {"forward"}, "program": {program}, "field": {field}}.Encode()
}
