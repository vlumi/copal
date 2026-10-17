// Package lineage follows a field's values across a system, back to where
// they come from and forward to where they go.
//
// Back follows them within a program, through the statements that move
// them, to the files it reads; and from a file, through the dataset its
// job step binds to it, to the step that wrote that dataset and the field
// of the record it wrote over the same bytes, and on, until the values
// start: at literals, initial values, or datasets no job writes.
//
// Forward is the mirror: through the statements that move the values on,
// to the files the program writes them to, through the datasets its step
// binds, to the steps that read them and the field of each reader's
// record over the same bytes, and on, until they reach the system's
// boundary: a report, a dataset no step reads, or a DISPLAY.
//
// A dataset a step reads was written by an earlier step of its job that
// writes the same dataset, or, for a generation the job does not write
// itself, by any job's step that writes a new generation of it; and the
// other way, a dataset a step writes is read by the later steps of its
// job that read the same dataset, and, unless it is temporary, in the
// runs that follow, by any job's step that reads it by its name, its
// current generation or an earlier one. A SORT that reads files
// with USING and writes them with GIVING passes their records on
// unchanged; a sort with an input procedure gets its records from the
// program's RELEASE statements, and one with an output procedure gives
// them to its RETURN statements.
package lineage

import (
	"slices"
	"strings"

	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/dataflow"
)

// Node is a field of a program, as a job step runs the program. Job and
// Step are empty for a program no job runs. Start and Length, when set,
// are the run of the field's bytes the node stands for, 1-based: those of
// a record that no field but a FILLER holds.
type Node struct {
	Job, Step, Program, Field string
	Start, Length             int
}

// Hop is a node of a trail. Going back, it has where its values come from
// within the program, and where each of the files it reads them from
// leads; going forward, where its values go within the program, where
// each of the files it writes them to leads, and the DISPLAYs that show
// them.
type Hop struct {
	Node     Node
	Trace    dataflow.Trace
	Sources  []Source
	Effect   dataflow.Effect
	Targets  []Target
	Displays []blueprint.Display
}

// Source is a field a program reads from a file, and where the bytes it
// reads were written: the datasets the step binds to the file, and the
// fields that other steps, or the program itself before a sort, wrote over
// those bytes. Sorts are the steps whose sorts passed the records on, with
// the datasets they read, and External the datasets no step writes, which
// come from outside the system.
type Source struct {
	File, Field string
	Datasets    []string
	Writers     []Node
	Sorts       []Node
	External    []string
}

// Target is a run of bytes a program writes to a file, and where they are
// read: the datasets the step binds to the file, and the fields that
// other steps, or the program itself after a sort, read over those bytes.
// Sorts are the steps whose sorts passed the records on, Outputs the
// datasets no step reads, which go outside the system, and Report the DD
// of the job's output the file is printed to.
type Target struct {
	File     string
	Datasets []string
	Readers  []Node
	Sorts    []Node
	Outputs  []string
	Report   string
}

// Trail is every hop from a field, each node once, nearest first. Capped
// is set when the trail was cut at maxHops.
type Trail struct {
	Hops   []Hop
	Capped bool
}

// maxHops is how many hops a trail follows: a field such as the business
// date reaches nearly every field of the system.
const maxHops = 400

// Back follows a field's values back from where a program has them. With
// no job step given, it starts from every step that runs the program.
func Back(s *blueprint.System, start Node) Trail {
	t := newTracer(s)
	return t.walk(start, func(n Node, p *blueprint.Program) (Hop, []Node) {
		hop := t.back(n, p)
		var next []Node
		for _, src := range hop.Sources {
			next = append(next, src.Writers...)
		}
		return hop, next
	})
}

// Forward follows a field's values on from where a program has them. With
// no job step given, it starts from every step that runs the program.
func Forward(s *blueprint.System, start Node) Trail {
	return newTracer(s).forwardFrom(start)
}

func (t tracer) forwardFrom(start Node) Trail {
	return t.walk(start, func(n Node, p *blueprint.Program) (Hop, []Node) {
		hop := t.forward(n, p)
		var next []Node
		for _, tgt := range hop.Targets {
			next = append(next, tgt.Readers...)
		}
		return hop, next
	})
}

type tracer struct {
	system   *blueprint.System
	programs map[string]*blueprint.Program
	maxHops  int
}

func newTracer(s *blueprint.System) tracer {
	t := tracer{system: s, programs: map[string]*blueprint.Program{}, maxHops: maxHops}
	for i := range s.Programs {
		t.programs[s.Programs[i].Name] = &s.Programs[i]
	}
	return t
}

// walk visits the nodes breadth first from the start, each once, making a
// hop of each and going on to the nodes it leads to.
func (t tracer) walk(start Node, step func(Node, *blueprint.Program) (Hop, []Node)) Trail {
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
		p, ok := t.programs[n.Program]
		if !ok {
			continue
		}
		if len(trail.Hops) == t.maxHops {
			trail.Capped = true
			break
		}
		hop, next := step(n, p)
		trail.Hops = append(trail.Hops, hop)
		queue = append(queue, next...)
	}
	return trail
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

func (t tracer) back(n Node, p *blueprint.Program) Hop {
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
	return hop
}

func (t tracer) forward(n Node, p *blueprint.Program) Hop {
	hop := Hop{Node: n, Effect: p.Affects(n.Field)}
	if n.Length > 0 {
		hop.Effect = p.AffectsBytes(n.Field, n.Start, n.Length)
	}
	for _, w := range hop.Effect.Files {
		tgt := Target{File: w.File}
		t.toFile(&tgt, n, p, w.File, w.Start, w.Length, 0)
		hop.Targets = append(hop.Targets, tgt)
	}
	shown := append([]string{n.Field}, hop.Effect.Items...)
	for _, d := range p.Displays {
		if slices.ContainsFunc(d.Items, func(item string) bool {
			return slices.ContainsFunc(shown, func(s string) bool { return p.Share(item, s) })
		}) {
			hop.Displays = append(hop.Displays, d)
		}
	}
	return hop
}

// maxSorts is how many sorts deep a source or target is followed, against
// a sort whose GIVING file another sort reads with USING, round and round.
const maxSorts = 8

// fromFile finds the writers of the bytes a program reads from a file: a
// sort's input, or the datasets the job step binds to the file.
func (t tracer) fromFile(src *Source, n Node, p *blueprint.Program, file string, start, length, depth int) {
	if depth > maxSorts {
		return
	}
	if sort, ok := sortOf(p, file); ok {
		if len(sort.Using) == 0 {
			src.Writers = append(src.Writers, over(Node{Job: n.Job, Step: n.Step, Program: p.Name}, outputs(p, file), start, length)...)
		}
		for _, using := range sort.Using {
			t.fromFile(src, n, p, using, start, length, depth+1)
		}
		return
	}
	dd, ok := t.dd(n, file)
	if !ok {
		return
	}
	for _, ds := range dd.Datasets {
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
				src.Sorts = append(src.Sorts, at)
				for _, using := range sort.Using {
					t.fromFile(src, at, q, using, start, length, depth+1)
				}
				continue
			}
			src.Writers = append(src.Writers, over(at, outputs(q, w.file), start, length)...)
		}
	}
}

// toFile finds the readers of the bytes a program writes to a file: a
// sort's output, or the steps that read the datasets the job step binds
// to the file.
func (t tracer) toFile(tgt *Target, n Node, p *blueprint.Program, file string, start, length, depth int) {
	if depth > maxSorts {
		return
	}
	if sort, ok := sortOf(p, file); ok {
		t.sorted(tgt, Node{Job: n.Job, Step: n.Step, Program: p.Name}, p, sort, start, length, depth)
		return
	}
	dd, ok := t.dd(n, file)
	if !ok {
		return
	}
	if dd.Sysout {
		tgt.Report = dd.Name
	}
	for _, ds := range dd.Datasets {
		tgt.Datasets = append(tgt.Datasets, ds)
		readers := t.readers(n, ds)
		if len(readers) == 0 {
			tgt.Outputs = append(tgt.Outputs, ds)
		}
		for _, r := range readers {
			q := t.programs[r.program]
			if q == nil {
				continue
			}
			at := Node{Job: r.job, Step: r.step, Program: q.Name}
			if sort, ok := usedBy(q, r.file); ok {
				tgt.Sorts = append(tgt.Sorts, at)
				t.sorted(tgt, at, q, sort, start, length, depth)
				continue
			}
			tgt.Readers = append(tgt.Readers, over(at, inputs(q, r.file), start, length)...)
		}
	}
}

// sorted follows records a sort has: to the files it writes with GIVING,
// or to the program's RETURN statements.
func (t tracer) sorted(tgt *Target, at Node, p *blueprint.Program, sort blueprint.Sort, start, length, depth int) {
	if len(sort.Giving) == 0 {
		tgt.Readers = append(tgt.Readers, over(at, inputs(p, sort.Work), start, length)...)
	}
	for _, giving := range sort.Giving {
		t.toFile(tgt, at, p, giving, start, length, depth+1)
	}
}

// dd is the DD statement the node's job step binds to a file.
func (t tracer) dd(n Node, file string) (blueprint.DD, bool) {
	step, ok := t.step(n.Job, n.Step)
	if !ok {
		return blueprint.DD{}, false
	}
	i := slices.IndexFunc(step.DDs, func(dd blueprint.DD) bool { return dd.File == file })
	if i < 0 {
		return blueprint.DD{}, false
	}
	return step.DDs[i], true
}

// access is a job step's program reading or writing a dataset through a
// file.
type access struct {
	job, step, program, file string
}

// writers lists the job steps that wrote a dataset a step reads: earlier
// steps of the job that write the same dataset, or, if none does, any
// job's steps that write a new generation of it, or it by its name.
func (t tracer) writers(n Node, ds string) []access {
	base, _, _ := strings.Cut(ds, "(")
	return t.earlier(n, ds, writes, func(d string) bool {
		return d == base || strings.HasPrefix(d, base+"(+")
	})
}

// readers lists the job steps that read a dataset a step writes: later
// steps of the job that read the same dataset, and, unless it is
// temporary, any job's steps that read it by its name, its current
// generation or an earlier one, in the runs that follow. Unlike a writer,
// which is the last one before the reader, every later reader reads what
// the step wrote.
func (t tracer) readers(n Node, ds string) []access {
	base, _, _ := strings.Cut(ds, "(")
	local := t.local(n, ds, false, reads)
	if strings.HasPrefix(ds, "&&") {
		return local
	}
	for _, a := range t.anyJob(reads, func(d string) bool {
		return d == base || strings.HasPrefix(d, base+"(0") || strings.HasPrefix(d, base+"(-")
	}) {
		if !slices.Contains(local, a) {
			local = append(local, a)
		}
	}
	return local
}

// earlier lists the steps that wrote a dataset the node's step reads:
// those of its job before it that write the same dataset, or, if there
// are none, any job's steps that write a generation that matches. A
// temporary dataset lives within its job.
func (t tracer) earlier(n Node, ds string, uses func(blueprint.DD) bool, generation func(string) bool) []access {
	local := t.local(n, ds, true, uses)
	if len(local) > 0 || strings.HasPrefix(ds, "&&") {
		return local
	}
	return t.anyJob(uses, generation)
}

// local lists the steps of the node's job, before it or after it, that use
// the same dataset in the way given.
func (t tracer) local(n Node, ds string, before bool, uses func(blueprint.DD) bool) []access {
	job, ok := t.job(n.Job)
	if !ok {
		return nil
	}
	i := slices.IndexFunc(job.Steps, func(s blueprint.Step) bool { return s.Name == n.Step })
	steps := job.Steps[i+1:]
	if before {
		steps = job.Steps[:max(i, 0)]
	}
	var list []access
	for _, step := range steps {
		list = append(list, bound(job.Name, step, uses, func(d string) bool { return d == ds })...)
	}
	return list
}

// anyJob lists every job's steps that use a dataset that matches in the
// way given.
func (t tracer) anyJob(uses func(blueprint.DD) bool, match func(string) bool) []access {
	var list []access
	for _, j := range t.system.Jobs {
		for _, step := range j.Steps {
			list = append(list, bound(j.Name, step, uses, match)...)
		}
	}
	return list
}

// bound lists the files through which a step uses the datasets that
// match, in the way given.
func bound(job string, step blueprint.Step, uses func(blueprint.DD) bool, match func(string) bool) []access {
	var list []access
	for _, dd := range step.DDs {
		if dd.File == "" || !uses(dd) || !slices.ContainsFunc(dd.Datasets, match) {
			continue
		}
		list = append(list, access{job: job, step: step.Name, program: step.Program, file: dd.File})
	}
	return list
}

func writes(dd blueprint.DD) bool {
	return slices.ContainsFunc(dd.Uses, func(u string) bool { return u == "write" || u == "append" || u == "update" })
}

func reads(dd blueprint.DD) bool {
	return slices.ContainsFunc(dd.Uses, func(u string) bool { return u == "read" || u == "update" })
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

// span is a field of a record a program reads from or writes to a file,
// and where its bytes are.
type span struct {
	field         string
	start, length int
}

// record is the record a program writes to a file or reads from it, and
// its fields.
type record struct {
	whole  span
	fields []span
}

func outputs(p *blueprint.Program, file string) record {
	i := slices.IndexFunc(p.Outputs, func(out blueprint.Output) bool { return out.File == file })
	if i < 0 {
		return record{}
	}
	return newRecord(p, p.Outputs[i].Record, p.Outputs[i].Fields, func(f blueprint.FieldOrigin) span { return span{f.Name, f.Start, f.Length} })
}

func inputs(p *blueprint.Program, file string) record {
	i := slices.IndexFunc(p.Inputs, func(in blueprint.Input) bool { return in.File == file })
	if i < 0 {
		return record{}
	}
	return newRecord(p, p.Inputs[i].Record, p.Inputs[i].Fields, func(f blueprint.FieldEffect) span { return span{f.Name, f.Start, f.Length} })
}

func newRecord[F any](p *blueprint.Program, name string, fields []F, spanOf func(F) span) record {
	r := record{}
	if i := slices.IndexFunc(p.Records, func(r blueprint.Record) bool { return r.Name == name }); i >= 0 {
		r.whole = span{name, p.Records[i].Start, p.Records[i].Length}
	}
	for _, f := range fields {
		r.fields = append(r.fields, spanOf(f))
	}
	return r
}

// over lists the fields of a record over the bytes, as nodes: those over
// exactly the bytes, if there are any, as a field of the same layout is;
// otherwise every field that shares some of them; and if none does, as
// where the record has a FILLER, those bytes of the record.
func over(at Node, r record, start, length int) []Node {
	var exact, overlapping []Node
	for _, s := range r.fields {
		n := Node{Job: at.Job, Step: at.Step, Program: at.Program, Field: s.field}
		switch {
		case s.start == start && s.length == length:
			exact = append(exact, n)
		case s.start < start+length && start < s.start+s.length:
			overlapping = append(overlapping, n)
		}
	}
	switch {
	case len(exact) > 0:
		return exact
	case len(overlapping) > 0 || r.whole.field == "":
		return overlapping
	}
	from := max(start, r.whole.start)
	to := min(start+length, r.whole.start+r.whole.length)
	if from >= to {
		return nil
	}
	return []Node{{Job: at.Job, Step: at.Step, Program: at.Program, Field: r.whole.field, Start: from - r.whole.start + 1, Length: to - from}}
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
	return findSort(p, func(s blueprint.Sort) bool { return s.Work == file })
}

// givenBy is the sort that writes the file with GIVING.
func givenBy(p *blueprint.Program, file string) (blueprint.Sort, bool) {
	return findSort(p, func(s blueprint.Sort) bool { return slices.Contains(s.Giving, file) })
}

// usedBy is the sort that reads the file with USING.
func usedBy(p *blueprint.Program, file string) (blueprint.Sort, bool) {
	return findSort(p, func(s blueprint.Sort) bool { return slices.Contains(s.Using, file) })
}

func findSort(p *blueprint.Program, match func(blueprint.Sort) bool) (blueprint.Sort, bool) {
	i := slices.IndexFunc(p.Sorts, match)
	if i < 0 {
		return blueprint.Sort{}, false
	}
	return p.Sorts[i], true
}
