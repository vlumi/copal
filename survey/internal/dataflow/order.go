package dataflow

import (
	"slices"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/survey/internal/controlflow"
)

// order is the program's statements in the order control can run them,
// each with the bytes it reads and the bytes it overwrites, so that values
// are followed only along paths that carry them: a value is not an origin
// of an item if every path from where it is moved to where the item is
// used moves another into the same bytes first.
type order struct {
	graph   *controlflow.Graph
	reached []bool
	facts   []facts // by node
}

// facts are what a statement does with data: the lineage edges it makes,
// by index, the bytes it reads, and the bytes it overwrites entirely.
type facts struct {
	edges []int
	reads []place
	kills []place
}

// statementFacts are a statement's facts by name, before the graph places
// the statement.
type statementFacts struct {
	edges        []int
	reads, kills []string
}

func newOrder(p *cobol.Program, s *storage, byStatement map[cobol.Pos]statementFacts) *order {
	g := controlflow.New(p)
	if !g.Covered {
		return nil
	}
	o := &order{graph: g, reached: g.Reachable(), facts: make([]facts, len(g.Nodes))}
	places := func(names []string) []place {
		var list []place
		for _, name := range names {
			if p, ok := s.places[name]; ok {
				list = append(list, p)
			}
		}
		return list
	}
	for i, n := range g.Nodes {
		if n.Statement == nil {
			continue
		}
		sf := byStatement[n.Statement.Pos]
		o.facts[i] = facts{edges: sf.edges, reads: places(sf.reads), kills: places(sf.kills)}
	}
	return o
}

// overwrites are the verbs that replace the whole of each item they write,
// unless they also read it, as ADD does the item it adds to.
var overwrites = map[string]bool{
	"MOVE": true, "COMPUTE": true, "ADD": true, "SUBTRACT": true, "MULTIPLY": true, "DIVIDE": true,
	"SET": true, "READ": true, "RETURN": true, "ACCEPT": true, "WRITE": true, "REWRITE": true, "RELEASE": true,
}

// kills lists the items a statement overwrites entirely: each item it
// writes by its plain name, not by a subscript or a reference
// modification, that it does not also read, and the records of the files
// it reads. A WRITE overwrites its record only to write FROM another, and
// a statement that leaves its targets as they were when a condition is
// raised, ON SIZE ERROR or AT END, overwrites nothing.
func (n *names) kills(s cobol.Statement, reads []string) []string {
	switch s.Verb {
	case "WRITE", "REWRITE", "RELEASE":
		if len(s.Reads) == 0 {
			return nil
		}
	}
	if !overwrites[s.Verb] || slices.ContainsFunc(s.Branches, func(b cobol.Branch) bool { return !strings.HasPrefix(b.Kind, "NOT ") }) {
		return nil
	}
	var kills []string
	for _, op := range s.Writes {
		item := n.item(op.Name)
		if op.Name != "" && op.Text == op.Name && !slices.Contains(reads, item) {
			kills = append(kills, item)
		}
	}
	if s.Verb == "READ" || s.Verb == "RETURN" {
		for _, f := range s.Files {
			kills = append(kills, n.records[f.Name]...)
		}
	}
	return kills
}

// state is bytes whose values are followed, at a node: just after its
// statement runs, going back, or just before, going forward.
type state struct {
	node int
	at   place
}

// uses lists the nodes control can reach whose statements read some of
// the bytes.
func (o *order) uses(bytes place, lineage []Edge, storage *storage) []int {
	var nodes []int
	for i, fs := range o.facts {
		if o.reached[i] && (slices.ContainsFunc(fs.reads, bytes.overlaps) || slices.ContainsFunc(fs.edges, func(e int) bool {
			from, ok := storage.places[lineage[e].From]
			return ok && from.overlaps(bytes)
		})) {
			nodes = append(nodes, i)
		}
	}
	return nodes
}

// at lists the nodes control can reach of the statements at the
// positions.
func (o *order) at(positions []cobol.Pos) []int {
	var nodes []int
	for i, n := range o.graph.Nodes {
		if o.reached[i] && n.Statement != nil && slices.Contains(positions, n.Statement.Pos) {
			nodes = append(nodes, i)
		}
	}
	return nodes
}

// defs lists the nodes control can reach whose statements move values
// into some of the bytes, and the entry, where the bytes hold their
// initial values or those the caller passed.
func (o *order) defs(bytes place, lineage []Edge, storage *storage) []int {
	nodes := []int{controlflow.Entry}
	for i, fs := range o.facts {
		if o.reached[i] && slices.ContainsFunc(fs.edges, func(e int) bool {
			to, ok := storage.places[lineage[e].To]
			return ok && !lineage[e].ToFile && !lineage[e].ToCall && !lineage[e].ToStream && to.overlaps(bytes)
		}) {
			nodes = append(nodes, i)
		}
	}
	return nodes
}

// without is what is left of the bytes once the statement's overwritten
// bytes are taken out: none, the bytes as they were, or the runs before
// and after what was overwritten.
func without(bytes place, kills []place) []place {
	left := []place{bytes}
	for _, k := range kills {
		var next []place
		for _, p := range left {
			if !p.overlaps(k) {
				next = append(next, p)
				continue
			}
			if p.start < k.start {
				next = append(next, place{record: p.record, start: p.start, length: k.start - p.start, group: p.group})
			}
			if end, kend := p.start+p.length, k.start+k.length; kend < end {
				next = append(next, place{record: p.record, start: kend, length: end - kend, group: p.group})
			}
		}
		left = next
	}
	return left
}

// usesAfter reports whether a statement uses the bytes it moves after it
// fills them itself: a WRITE FROM moves its record into the file's record
// before writing it, and a READ INTO reads the file's record before moving
// it into the item.
func usesAfter(s *cobol.Statement) bool {
	switch s.Verb {
	case "WRITE", "REWRITE", "RELEASE", "READ", "RETURN":
		return true
	}
	return false
}
