package dataflow

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
	"github.com/vlumi/copal/survey/internal/controlflow"
)

// Origin is where a data item's values come from in the end: a field of a
// record a file is read into, a literal, an item's initial VALUE, a
// program called, which may change an item passed to it by reference, the
// caller of the program, which passes the items it is called with, or the
// job's SYSIN, which an ACCEPT reads.
type Origin struct {
	File      string // the file, for a field read from one
	Item      string // the item at the bytes, or the item with the initial value
	Literal   string // the literal, as written
	Call      string // the program called, for an item passed to it
	Arg       int    // the position of the argument passed, 1-based
	Parameter int    // the position of the USING item passed in, 1-based
	Stream    string // SYSIN, for an item accepted from it
	// Start and Length are the bytes of a file's record the values come
	// from, 1-based, when they are a part of Item and not all of it, as of
	// a record read whole into another.
	Start, Length int
}

// String writes the origin: "file TRAN-IN TD-AMOUNT", "literal 'TXPOST'",
// "initial value of WS-COUNT", "call AUTHPUR argument 6 HD-AMOUNT",
// "parameter 1 RQ-AMOUNT", "stream SYSIN AUTH-REQUEST".
func (o Origin) String() string {
	switch {
	case o.File != "" && o.Length > 0:
		return fmt.Sprintf("file %s %s %d-%d", o.File, o.Item, o.Start, o.Start+o.Length-1)
	case o.File != "":
		return "file " + o.File + " " + o.Item
	case o.Literal != "":
		return "literal " + o.Literal
	case o.Call != "" && o.Arg == 0:
		return fmt.Sprintf("call %s returns %s", o.Call, o.Item)
	case o.Call != "":
		return fmt.Sprintf("call %s argument %d %s", o.Call, o.Arg, o.Item)
	case o.Parameter > 0:
		return fmt.Sprintf("parameter %d %s", o.Parameter, o.Item)
	case o.Stream != "":
		return "stream " + o.Stream + " " + o.Item
	}
	return "initial value of " + o.Item
}

// place is where an item's bytes are: in which record, from where, how
// many; and whether it has an initial VALUE and fields of its own.
type place struct {
	record        string
	start, length int
	value, group  bool
}

func (p place) overlaps(q place) bool {
	return p.record == q.record && p.start < q.start+q.length && q.start < p.start+p.length
}

// storage is where every item the program declares lives, so that values
// can be followed by bytes, not just by name: into a group, out of a field
// of it, through a field that redefines another.
type storage struct {
	places map[string]place
	items  map[string][]string // each record's items, in declaration order
}

func newStorage(records *layout.ProgramRecords) *storage {
	s := &storage{places: map[string]place{}, items: map[string][]string{}}
	seen := map[string]int{}
	var walk func(record string, f *layout.Field)
	walk = func(record string, f *layout.Field) {
		seen[f.Name]++
		s.places[f.Name] = place{record: record, start: f.Start, length: f.Size(), value: len(f.Values) > 0, group: len(f.Children) > 0}
		s.items[record] = append(s.items[record], f.Name)
		for _, child := range f.Children {
			walk(record, child)
		}
	}
	for _, r := range records.All() {
		walk(r.Name, r)
	}
	for name, n := range seen {
		if n > 1 {
			delete(s.places, name)
		}
	}
	return s
}

// part is the place of a run of an item's bytes, from start, 1-based, for
// length.
func (s *storage) part(item string, start, length int) (place, bool) {
	p, ok := s.places[item]
	if !ok || start < 1 || start+length-1 > p.length {
		return place{}, false
	}
	return place{record: p.record, start: p.start + start - 1, length: length, group: p.group}, true
}

// at is the narrowest item that holds a record's bytes from start, for
// length; the record itself if no field does.
func (s *storage) at(record string, start, length int) string {
	best, size := record, s.places[record].length
	for _, name := range s.items[record] {
		p, ok := s.places[name]
		if ok && p.start <= start && start+length <= p.start+p.length && p.length < size {
			best, size = name, p.length
		}
	}
	return best
}

// Trace is where a data item's values come from: the origins at the
// start, and the statements that move them on to the item, as lineage
// edges, in the order the program declares them.
type Trace struct {
	Origins []Origin
	Edges   []Edge
}

// Origins is where a data item's values come from in the end, as Trace
// finds them.
func (f Flow) Origins(item string) []Origin {
	return f.Trace(item).Origins
}

// Trace traces where a data item's values come from: through every
// statement that moves a value into it or into bytes it shares, a group it
// is in, a field of it or a field that redefines it, and on through what
// those statements move from, to the files, literals and initial values
// at the start. A group moved into a group, or anything moved into an item
// of its own length, is followed byte for byte. A literal moved into a
// whole group around the item, such as MOVE SPACES TO A-RECORD, only
// clears it, and counts only for an item with no other origin.
//
// The values are followed back from every statement that uses the item,
// along the paths control can take to it, and stop where a statement
// overwrites the bytes. Where the paths are not known, every statement
// that moves a value into the bytes counts, whatever the order.
func (f Flow) Trace(item string) Trace {
	start, ok := f.storage.places[item]
	if !ok {
		return Trace{}
	}
	return f.traceFrom(start, nil)
}

// TraceBytes traces where the values of some of a data item's bytes come
// from, from start, 1-based, for length, as Trace does.
func (f Flow) TraceBytes(item string, start, length int) Trace {
	p, ok := f.storage.part(item, start, length)
	if !ok {
		return Trace{}
	}
	return f.traceFrom(p, nil)
}

// TraceAt traces where a data item's values come from, as Trace does, as
// the statements at the positions use them: the WRITEs of a record, say.
func (f Flow) TraceAt(item string, positions []cobol.Pos) Trace {
	start, ok := f.storage.places[item]
	if !ok {
		return Trace{}
	}
	return f.traceFrom(start, positions)
}

func (f Flow) traceFrom(start place, positions []cobol.Pos) Trace {
	if f.order == nil {
		return f.traceAll(start)
	}
	nodes := f.order.uses(start, f.Lineage, f.storage)
	if positions != nil {
		nodes = f.order.at(positions)
	}
	if len(nodes) == 0 {
		return f.traceAll(start)
	}
	return f.traceOrdered(start, nodes)
}

// tracing collects the origins a trace finds and the edges it follows.
type tracing struct {
	flow              Flow
	found, clearing   map[Origin]bool
	followed, cleared map[int]bool
}

func newTracing(f Flow) *tracing {
	return &tracing{flow: f, found: map[Origin]bool{}, clearing: map[Origin]bool{}, followed: map[int]bool{}, cleared: map[int]bool{}}
}

// place notes the origins of bytes the trace reaches: the item's initial
// value and, in a program called, the caller's parameter.
func (t *tracing) place(here place) {
	name := t.flow.storage.at(here.record, here.start, here.length)
	if p := t.flow.storage.places[name]; p.value {
		t.found[Origin{Item: name}] = true
	}
	if n := t.flow.parameters[here.record]; n > 0 {
		t.found[Origin{Parameter: n, Item: name}] = true
	}
}

// edge follows a lineage edge back from the bytes it moves into: to its
// origin, or to the bytes of the item it moves from, which it returns.
func (t *tracing) edge(i int, here place) (from place, ok bool) {
	e := t.flow.Lineage[i]
	to, ok := t.flow.storage.places[e.To]
	if e.ToFile || e.ToCall || e.ToStream || !ok || !to.overlaps(here) {
		return place{}, false
	}
	name := t.flow.storage.at(here.record, here.start, here.length)
	switch {
	case e.FromFile:
		o := Origin{File: e.From, Item: name}
		if item := t.flow.storage.places[name]; here.length < item.length {
			o.Start, o.Length = here.start, here.length
		}
		t.found[o] = true
	case e.FromCall:
		t.found[Origin{Call: e.From, Arg: e.Arg, Item: name}] = true
	case e.FromStream:
		t.found[Origin{Stream: e.From, Item: name}] = true
	case e.FromLiteral && to.group && to.length > here.length:
		t.clearing[Origin{Literal: e.From}] = true
		t.cleared[i] = true
		return place{}, false
	case e.FromLiteral:
		t.found[Origin{Literal: e.From}] = true
	default:
		source, ok := t.flow.storage.places[e.From]
		if !ok {
			return place{}, false
		}
		t.followed[i] = true
		return follow(here, to, source), true
	}
	t.followed[i] = true
	return place{}, false
}

func (t *tracing) trace() Trace {
	found, followed := t.found, t.followed
	if len(found) == 0 {
		found, followed = t.clearing, t.cleared
	}
	trace := Trace{Origins: make([]Origin, 0, len(found))}
	for o := range found {
		trace.Origins = append(trace.Origins, o)
	}
	slices.SortFunc(trace.Origins, func(a, b Origin) int { return strings.Compare(a.String(), b.String()) })
	for i, e := range t.flow.Lineage {
		if followed[i] {
			trace.Edges = append(trace.Edges, e)
		}
	}
	return trace
}

// traceAll follows every statement that moves a value into the bytes,
// whatever the order.
func (f Flow) traceAll(start place) Trace {
	t := newTracing(f)
	visited := map[place]bool{}
	queue := []place{start}
	for len(queue) > 0 {
		here := queue[0]
		queue = queue[1:]
		if visited[here] {
			continue
		}
		visited[here] = true
		t.place(here)
		for i := range f.Lineage {
			if from, ok := t.edge(i, here); ok {
				queue = append(queue, from)
			}
		}
	}
	return t.trace()
}

// traceOrdered follows the bytes back from the statements at the nodes,
// along the paths control can take to them: through each statement's moves
// into them, to what it moves from, and past it to the statements before,
// for the bytes it does not overwrite. Bytes still to be followed at the
// program's entry hold their initial values or the caller's.
func (f Flow) traceOrdered(start place, nodes []int) Trace {
	t := newTracing(f)
	g := f.order.graph
	visited := map[state]bool{}
	var queue []state
	before := func(n int, here place) {
		for _, p := range g.Nodes[n].Prev {
			queue = append(queue, state{p, here})
		}
	}
	for _, n := range nodes {
		if usesAfter(g.Nodes[n].Statement) {
			queue = append(queue, state{n, start})
		} else {
			before(n, start)
		}
	}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		if visited[s] {
			continue
		}
		visited[s] = true
		if s.node == controlflow.Entry {
			t.place(s.at)
			continue
		}
		facts := f.order.facts[s.node]
		for _, i := range facts.edges {
			if from, ok := t.edge(i, s.at); ok {
				queue = append(queue, state{s.node, from})
			}
		}
		for _, left := range without(s.at, facts.kills) {
			before(s.node, left)
		}
	}
	return t.trace()
}

// follow maps bytes of a statement's target back to its source: the same
// bytes, by offset, when a group is moved into a group or an item into one
// of its own length, and the whole source otherwise.
func follow(here, to, from place) place {
	if (!to.group || !from.group) && to.length != from.length {
		return from
	}
	start := max(here.start, to.start)
	end := min(here.start+here.length, to.start+to.length)
	offset := start - to.start
	length := min(end-start, from.length-offset)
	if length <= 0 {
		return from
	}
	return place{record: from.record, start: from.start + offset, length: length, group: true}
}
