// Package lineage follows a field's values back across a system: within a
// program, through the statements that move them, to the files it reads;
// and from a file, through the dataset its job step binds to it, to the
// step that wrote that dataset and the field of the record it wrote over
// the same bytes, and on, until the values start: at literals, initial
// values, or datasets no job writes.
//
// A dataset a step reads was written by an earlier step of its job that
// writes the same dataset, or, for a generation the job does not write
// itself, by any job's step that writes a new generation of it. A SORT
// that reads files with USING and writes them with GIVING passes their
// records on unchanged, and a sort with an input procedure gets its
// records from the program's RELEASE statements.
package lineage

import (
	"slices"
	"strings"

	"github.com/vlumi/copal/blueprint"
	"github.com/vlumi/copal/dataflow"
)

// Node is a field of a program, as a job step runs the program. Job and
// Step are empty for a program no job runs.
type Node struct {
	Job, Step, Program, Field string
}

// Hop is a node of a trail: where its values come from within the program,
// and where each of the files it reads them from leads.
type Hop struct {
	Node    Node
	Trace   dataflow.Trace
	Sources []Source
}

// Source is a field a program reads from a file, and where the bytes it
// reads were written: the datasets the step binds to the file, and the
// fields that other steps, or the program itself before a sort, wrote over
// those bytes. External is set for a dataset no step writes, which comes
// from outside the system.
type Source struct {
	File, Field string
	Datasets    []string
	Writers     []Node
	External    []string
}

// Trail is every hop back from a field, each node once, nearest first.
type Trail struct {
	Hops []Hop
}

// Back follows a field's values back from where a program has them. With
// no job step given, it starts from every step that runs the program.
func Back(s *blueprint.System, start Node) Trail {
	t := tracer{system: s, programs: map[string]*blueprint.Program{}}
	for i := range s.Programs {
		t.programs[s.Programs[i].Name] = &s.Programs[i]
	}
	queue := []Node{start}
	if start.Job == "" {
		queue = t.runs(start)
	}
	seen := map[Node]bool{}
	var trail Trail
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if seen[n] {
			continue
		}
		seen[n] = true
		hop, ok := t.hop(n)
		if !ok {
			continue
		}
		trail.Hops = append(trail.Hops, hop)
		for _, src := range hop.Sources {
			queue = append(queue, src.Writers...)
		}
	}
	return trail
}

type tracer struct {
	system   *blueprint.System
	programs map[string]*blueprint.Program
}

// runs is the node in every job step that runs the program, or the node
// alone if none does.
func (t tracer) runs(n Node) []Node {
	var nodes []Node
	for _, job := range t.system.Jobs {
		for _, step := range job.Steps {
			if step.Program == n.Program {
				nodes = append(nodes, Node{Job: job.Name, Step: step.Name, Program: n.Program, Field: n.Field})
			}
		}
	}
	if len(nodes) == 0 {
		return []Node{n}
	}
	return nodes
}

func (t tracer) hop(n Node) (Hop, bool) {
	p, ok := t.programs[n.Program]
	if !ok {
		return Hop{}, false
	}
	hop := Hop{Node: n, Trace: p.Trace(n.Field)}
	for _, o := range hop.Trace.Origins {
		if o.File == "" {
			continue
		}
		start, length, ok := bytesOf(p, o.File, o.Item)
		if !ok {
			continue
		}
		src := Source{File: o.File, Field: o.Item}
		t.fromFile(&src, n, p, o.File, start, length, 0)
		hop.Sources = append(hop.Sources, src)
	}
	return hop, true
}

// maxSorts is how many sorts deep a source is followed, against a sort
// whose GIVING file another sort reads with USING, round and round.
const maxSorts = 8

// fromFile finds the writers of the bytes a program reads from a file: a
// sort's input, or the datasets the job step binds to the file.
func (t tracer) fromFile(src *Source, n Node, p *blueprint.Program, file string, start, length, depth int) {
	if depth > maxSorts {
		return
	}
	if sort, ok := sortOf(p, file); ok {
		if len(sort.Using) == 0 {
			src.Writers = append(src.Writers, written(Node{Job: n.Job, Step: n.Step, Program: p.Name}, p, file, start, length)...)
		}
		for _, using := range sort.Using {
			t.fromFile(src, n, p, using, start, length, depth+1)
		}
		return
	}
	if n.Job == "" {
		return
	}
	for _, ds := range t.datasets(n, file) {
		src.Datasets = append(src.Datasets, ds)
		writers := t.writers(n, ds)
		if len(writers) == 0 {
			src.External = append(src.External, ds)
		}
		for _, w := range writers {
			q := t.programs[w.program]
			if q == nil {
				continue
			}
			at := Node{Job: w.job, Step: w.step, Program: q.Name}
			if sort, ok := givenBy(q, w.file); ok {
				for _, using := range sort.Using {
					t.fromFile(src, at, q, using, start, length, depth+1)
				}
				continue
			}
			src.Writers = append(src.Writers, written(at, q, w.file, start, length)...)
		}
	}
}

// datasets lists the datasets the node's job step binds to a file.
func (t tracer) datasets(n Node, file string) []string {
	step, ok := t.step(n.Job, n.Step)
	if !ok {
		return nil
	}
	for _, dd := range step.DDs {
		if dd.File == file {
			return dd.Datasets
		}
	}
	return nil
}

// writer is a job step's program writing a dataset through a file.
type writer struct {
	job, step, program, file string
}

// writers lists the job steps that wrote a dataset a step reads: earlier
// steps of the job that write the same dataset, or, if none does, any
// job's steps that write a new generation of it, or it by its name.
func (t tracer) writers(n Node, ds string) []writer {
	job, ok := t.job(n.Job)
	if !ok || strings.HasPrefix(ds, "&&") && !t.writtenBefore(job, n.Step, ds) {
		return nil
	}
	var earlier []writer
	for _, step := range job.Steps {
		if step.Name == n.Step {
			break
		}
		earlier = append(earlier, writtenBy(job.Name, step, func(d string) bool { return d == ds })...)
	}
	if len(earlier) > 0 || strings.HasPrefix(ds, "&&") {
		return earlier
	}
	base, _, _ := strings.Cut(ds, "(")
	var writers []writer
	for _, j := range t.system.Jobs {
		for _, step := range j.Steps {
			writers = append(writers, writtenBy(j.Name, step, func(d string) bool {
				return d == base || strings.HasPrefix(d, base+"(+")
			})...)
		}
	}
	return writers
}

func (t tracer) writtenBefore(job *blueprint.Job, step, ds string) bool {
	for _, s := range job.Steps {
		if s.Name == step {
			return false
		}
		if len(writtenBy(job.Name, s, func(d string) bool { return d == ds })) > 0 {
			return true
		}
	}
	return false
}

// writtenBy lists the files through which a step writes the datasets that
// match.
func writtenBy(job string, step blueprint.Step, match func(string) bool) []writer {
	var writers []writer
	for _, dd := range step.DDs {
		if dd.File == "" || !writes(dd) || !slices.ContainsFunc(dd.Datasets, match) {
			continue
		}
		writers = append(writers, writer{job: job, step: step.Name, program: step.Program, file: dd.File})
	}
	return writers
}

func writes(dd blueprint.DD) bool {
	return slices.ContainsFunc(dd.Uses, func(u string) bool { return u == "write" || u == "append" || u == "update" })
}

func (t tracer) job(name string) (*blueprint.Job, bool) {
	i := slices.IndexFunc(t.system.Jobs, func(j blueprint.Job) bool { return j.Name == name })
	if i < 0 {
		return nil, false
	}
	return &t.system.Jobs[i], true
}

func (t tracer) step(job, name string) (blueprint.Step, bool) {
	j, ok := t.job(job)
	if !ok {
		return blueprint.Step{}, false
	}
	i := slices.IndexFunc(j.Steps, func(s blueprint.Step) bool { return s.Name == name })
	if i < 0 {
		return blueprint.Step{}, false
	}
	return j.Steps[i], true
}

// written lists the fields a program writes to a file over the bytes, as
// nodes: those over exactly the bytes, if there are any, as a field of
// the same layout is; otherwise every field that shares some of them.
func written(at Node, p *blueprint.Program, file string, start, length int) []Node {
	var exact, overlapping []Node
	for _, out := range p.Outputs {
		if out.File != file {
			continue
		}
		for _, f := range out.Fields {
			n := Node{Job: at.Job, Step: at.Step, Program: p.Name, Field: f.Name}
			switch {
			case f.Start == start && f.Length == length:
				exact = append(exact, n)
			case f.Start < start+length && start < f.Start+f.Length:
				overlapping = append(overlapping, n)
			}
		}
	}
	if len(exact) > 0 {
		return exact
	}
	return overlapping
}

// bytesOf is where a field is in the record of a file it belongs to.
func bytesOf(p *blueprint.Program, file, field string) (start, length int, ok bool) {
	for _, r := range p.Records {
		if r.File != file {
			continue
		}
		if f, found := find(r.Field, field); found {
			return f.Start, f.Length, true
		}
	}
	return 0, 0, false
}

func find(f blueprint.Field, name string) (blueprint.Field, bool) {
	if f.Name == name {
		return f, true
	}
	for _, child := range f.Fields {
		if found, ok := find(child, name); ok {
			return found, true
		}
	}
	return blueprint.Field{}, false
}

// sortOf is the sort whose work file the file is.
func sortOf(p *blueprint.Program, file string) (blueprint.Sort, bool) {
	i := slices.IndexFunc(p.Sorts, func(s blueprint.Sort) bool { return s.Work == file })
	if i < 0 {
		return blueprint.Sort{}, false
	}
	return p.Sorts[i], true
}

// givenBy is the sort that writes the file with GIVING.
func givenBy(p *blueprint.Program, file string) (blueprint.Sort, bool) {
	i := slices.IndexFunc(p.Sorts, func(s blueprint.Sort) bool { return slices.Contains(s.Giving, file) })
	if i < 0 {
		return blueprint.Sort{}, false
	}
	return p.Sorts[i], true
}
