package dataflow

import (
	"slices"

	"github.com/vlumi/copal/survey/internal/controlflow"
)

// Effect is where a data item's values go: the items they reach, the
// bytes of the files they are written to, the bytes of the arguments they
// are passed in to programs called, the bytes of the items the program was
// called with that they are moved into, back to the caller, and the
// statements that move them there, as lineage edges, in the order the
// program declares them.
type Effect struct {
	Items      []string
	Files      []Written
	Calls      []Passed
	Parameters []Passed
	Edges      []Edge
}

// Passed is a run of bytes of an item passed between programs: to Program,
// the program called, as the argument at Arg, or, with no Program, as the
// USING item at Arg, back to the caller; Start is 1-based, within the
// item.
type Passed struct {
	Program       string
	Arg           int
	Start, Length int
}

// Written is a run of bytes of the record a file is written from: the
// file, and where in its record the bytes are, 1-based.
type Written struct {
	File          string
	Start, Length int
}

// Affects follows where a data item's values go, the mirror of Trace:
// through every statement that moves a value out of it or out of bytes it
// shares, a group it is in, a field of it or a field that redefines it,
// on to the items those statements move into, and to the files the
// records it reaches are written to. A group moved into a group, or an
// item moved into one of its own length, is followed byte for byte.
//
// The values are followed on from every statement that moves a value into
// the item, and from the program's entry, along the paths control can
// take, and stop where a statement overwrites the bytes. Where the paths
// are not known, every statement that moves a value out of the bytes
// counts, whatever the order.
func (f Flow) Affects(item string) Effect {
	start, ok := f.storage.places[item]
	if !ok {
		return Effect{}
	}
	return f.affectsFrom(start)
}

// AffectsBytes follows where the values of some of a data item's bytes go,
// from start, 1-based, for length, as Affects does: the bytes of a record
// that a FILLER holds, say, which no field names.
func (f Flow) AffectsBytes(item string, start, length int) Effect {
	p, ok := f.storage.part(item, start, length)
	if !ok {
		return Effect{}
	}
	return f.affectsFrom(p)
}

func (f Flow) affectsFrom(start place) Effect {
	if f.order == nil {
		return f.affectsAll(start)
	}
	return f.affectsOrdered(start, f.order.defs(start, f.Lineage, f.storage))
}

// affecting collects where an effect's values go and the edges it
// follows.
type affecting struct {
	flow     Flow
	start    place
	effect   Effect
	reached  map[string]bool
	written  map[Written]bool
	passed   map[Passed]bool
	followed map[int]bool
}

func newAffecting(f Flow, start place) *affecting {
	return &affecting{flow: f, start: start, reached: map[string]bool{}, written: map[Written]bool{}, passed: map[Passed]bool{}, followed: map[int]bool{}}
}

// item notes an item the values reach.
func (a *affecting) item(here place) {
	if name := a.flow.storage.at(here.record, here.start, here.length); here != a.start && !a.reached[name] {
		a.reached[name] = true
		a.effect.Items = append(a.effect.Items, name)
	}
}

// parameter notes bytes of an item the program was called with, which
// the values reach and so go back to the caller.
func (a *affecting) parameter(here place) {
	if n := a.flow.parameters[here.record]; n > 0 && here != a.start {
		a.pass(Passed{Arg: n, Start: here.start, Length: here.length}, &a.effect.Parameters)
	}
}

func (a *affecting) pass(p Passed, to *[]Passed) {
	if !a.passed[p] {
		a.passed[p] = true
		*to = append(*to, p)
	}
}

// edge follows a lineage edge on from the bytes it moves out of: to the
// file or program it sends them to, or to the bytes of the item it moves
// into, which it returns.
func (a *affecting) edge(i int, here place) (to place, ok bool) {
	e := a.flow.Lineage[i]
	from, ok := a.flow.storage.bytesOf(e.From, e.FromStart, e.FromLength)
	if e.FromFile || e.FromLiteral || e.FromCall || e.FromStream || e.ToStream || !ok || !from.overlaps(here) {
		return place{}, false
	}
	switch {
	case e.ToCall:
		w := writtenBytes(e.To, here, from)
		a.pass(Passed{Program: e.To, Arg: e.Arg, Start: w.Start, Length: w.Length}, &a.effect.Calls)
	case e.ToFile:
		if w := writtenBytes(e.To, here, from); !a.written[w] {
			a.written[w] = true
			a.effect.Files = append(a.effect.Files, w)
		}
	default:
		target, ok := a.flow.storage.bytesOf(e.To, e.ToStart, e.ToLength)
		if !ok {
			return place{}, false
		}
		a.followed[i] = true
		return follow(here, from, target), true
	}
	a.followed[i] = true
	return place{}, false
}

// result is the effect found, with the bytes that lie within others sent
// to the same place left out: values followed from where they are set in
// part, as well as from where they are set whole, are written in part.
func (a *affecting) result() Effect {
	for i, e := range a.flow.Lineage {
		if a.followed[i] {
			a.effect.Edges = append(a.effect.Edges, e)
		}
	}
	a.effect.Files = outermost(a.effect.Files, func(w, v Written) bool {
		return w.File == v.File && within(w.Start, w.Length, v.Start, v.Length)
	})
	inside := func(p, q Passed) bool {
		return p.Program == q.Program && p.Arg == q.Arg && within(p.Start, p.Length, q.Start, q.Length)
	}
	a.effect.Calls = outermost(a.effect.Calls, inside)
	a.effect.Parameters = outermost(a.effect.Parameters, inside)
	return a.effect
}

// outermost leaves out the runs of bytes that lie inside another.
func outermost[T comparable](runs []T, inside func(T, T) bool) []T {
	all := slices.Clone(runs)
	return slices.DeleteFunc(runs, func(r T) bool {
		return slices.ContainsFunc(all, func(o T) bool { return o != r && inside(r, o) })
	})
}

func within(start, length, outerStart, outerLength int) bool {
	return outerStart <= start && start+length <= outerStart+outerLength
}

// affectsAll follows every statement that moves a value out of the bytes,
// whatever the order.
func (f Flow) affectsAll(start place) Effect {
	a := newAffecting(f, start)
	visited := map[place]bool{}
	queue := []place{start}
	for len(queue) > 0 {
		here := queue[0]
		queue = queue[1:]
		if visited[here] {
			continue
		}
		visited[here] = true
		a.item(here)
		a.parameter(here)
		for i := range f.Lineage {
			if to, ok := a.edge(i, here); ok {
				queue = append(queue, to)
			}
		}
	}
	return a.result()
}

// onward is bytes whose values are followed forward, at a node: just
// before its statement runs, or, for bytes it has just moved, after; and
// whether a statement has moved them since the start.
type onward struct {
	node         int
	at           place
	after, moved bool
}

// affectsOrdered follows the bytes on from the statements at the nodes,
// along the paths control can take from them: through each statement's
// moves out of them, to what it moves into, and past it to the statements
// after, for the bytes it does not overwrite. Bytes moved into an item the
// program was called with that are still to be followed at its end go
// back to the caller.
func (f Flow) affectsOrdered(start place, nodes []int) Effect {
	a := newAffecting(f, start)
	g := f.order.graph
	visited := map[onward]bool{}
	var queue []onward
	for _, n := range nodes {
		queue = append(queue, onward{node: n, at: start, after: true})
	}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		if visited[s] {
			continue
		}
		visited[s] = true
		if s.node == controlflow.End {
			if s.moved {
				a.parameter(s.at)
			}
			continue
		}
		facts := f.order.facts[s.node]
		for _, i := range facts.edges {
			if to, ok := a.edge(i, s.at); ok {
				a.item(to)
				queue = append(queue, onward{node: s.node, at: to, after: true, moved: true})
			}
		}
		left := []place{s.at}
		if !s.after {
			left = without(s.at, facts.kills)
		}
		for _, l := range left {
			for _, n := range g.Nodes[s.node].Next {
				queue = append(queue, onward{node: n, at: l, moved: s.moved})
			}
		}
	}
	return a.result()
}

// Share reports whether two data items share any bytes: one is in the
// other, they are the same, or one redefines the other.
func (f Flow) Share(a, b string) bool {
	pa, okA := f.storage.places[a]
	pb, okB := f.storage.places[b]
	return okA && okB && pa.overlaps(pb)
}

// writtenBytes is where bytes of a record are in the record written to a
// file.
func writtenBytes(file string, here, record place) Written {
	start := max(here.start, record.start)
	end := min(here.start+here.length, record.start+record.length)
	return Written{File: file, Start: start - record.start + 1, Length: end - start}
}
