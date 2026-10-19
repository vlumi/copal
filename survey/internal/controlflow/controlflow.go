// Package controlflow finds the paths control can take through a program's
// statements: in order within a paragraph, into the branches of IF,
// EVALUATE, SEARCH and the conditional phrases, round the loops of PERFORM,
// to the target of a GO TO, out of the paragraph at EXIT PARAGRAPH and
// out of the section at EXIT SECTION, through the DECLARATIVES for a file
// from an input or output statement that fails on it, on
// from one paragraph into the next, and into a performed range and back to
// the statement after the PERFORM.
//
// A performed range is laid out anew for every PERFORM of it, so that
// control returns only to the statement after the PERFORM that entered it.
// A loop may run no times, and a branch that is not one of a complete set,
// an IF with no ELSE, say, may be skipped; the graph has every path the
// program can take, and some it cannot, as a condition is not evaluated.
package controlflow

import (
	"slices"

	"github.com/vlumi/copal/bedrock/cobol"
)

// Graph is a program's statements and the paths between them. Node 0 is
// the program's entry and node 1 its end; a node with no statement joins
// paths, as the start of a paragraph does.
//
// Covered is false when the graph leaves out paths the program can take: a
// GO TO out of a performed range, a PERFORM of a range already being
// performed, a target the program does not declare, an EXIT SECTION that
// leaves the range performed, or an EXIT PERFORM, which it does not
// follow.
type Graph struct {
	Nodes   []Node
	Covered bool
}

// Node is a statement on the paths through the program, and the nodes
// control can go to next and come from.
type Node struct {
	Statement *cobol.Statement
	Next      []int
	Prev      []int
}

// Entry and End are the nodes where the program starts and ends.
const (
	Entry = 0
	End   = 1
)

// maxNodes bounds the graph, against performed ranges nested deep enough to
// lay out too many times.
const maxNodes = 100_000

// New lays out the paths through a program's procedure division, which
// starts at its first paragraph after the DECLARATIVES.
func New(p *cobol.Program) *Graph {
	b := &builder{program: p, graph: &Graph{Nodes: make([]Node, 2), Covered: true}}
	first := slices.IndexFunc(p.Paragraphs, func(para cobol.Paragraph) bool { return !para.Declarative })
	if first >= 0 {
		b.link(Entry, b.performed(first, len(p.Paragraphs)-1, End))
	} else {
		b.link(Entry, End)
	}
	return b.graph
}

// Reachable lists the nodes control can reach from the entry.
func (g *Graph) Reachable() []bool {
	reached := make([]bool, len(g.Nodes))
	queue := []int{Entry}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if reached[n] {
			continue
		}
		reached[n] = true
		queue = append(queue, g.Nodes[n].Next...)
	}
	return reached
}

type builder struct {
	program *cobol.Program
	graph   *Graph
	// active are the ranges being laid out, from the outermost PERFORM in.
	active [][2]int
	// paragraphEnd is where control goes when the paragraph being laid
	// out ends, which EXIT PARAGRAPH goes to, and sectionEnd when its
	// section ends, which EXIT SECTION goes to, or -1 if that is outside
	// the range being laid out.
	paragraphEnd, sectionEnd int
}

func (b *builder) node(s *cobol.Statement) int {
	b.graph.Nodes = append(b.graph.Nodes, Node{Statement: s})
	if len(b.graph.Nodes) > maxNodes {
		b.graph.Covered = false
	}
	return len(b.graph.Nodes) - 1
}

func (b *builder) link(from, to int) {
	if !slices.Contains(b.graph.Nodes[from].Next, to) {
		b.graph.Nodes[from].Next = append(b.graph.Nodes[from].Next, to)
		b.graph.Nodes[to].Prev = append(b.graph.Nodes[to].Prev, from)
	}
}

// performed lays out the paragraphs from first through last, control
// running on from each into the next and from the last to exit, and
// returns the node of the first one's start. A paragraph is laid out only
// once control can reach it.
func (b *builder) performed(first, last, exit int) int {
	b.active = append(b.active, [2]int{first, last})
	defer func() { b.active = b.active[:len(b.active)-1] }()
	starts := map[int]int{}
	var queue []int
	start := func(i int) int {
		if n, ok := starts[i]; ok {
			return n
		}
		starts[i] = b.node(nil)
		queue = append(queue, i)
		return starts[i]
	}
	r := scope{first: first, last: last, start: start}
	entry := start(first)
	for len(queue) > 0 && b.graph.Covered {
		i := queue[0]
		queue = queue[1:]
		next := exit
		if i < last {
			next = start(i + 1)
		}
		outerParagraph, outerSection := b.paragraphEnd, b.sectionEnd
		b.paragraphEnd, b.sectionEnd = next, b.endOfSection(i, last, exit, start)
		b.link(starts[i], b.sequence(r, b.program.Paragraphs[i].Statements, next))
		b.paragraphEnd, b.sectionEnd = outerParagraph, outerSection
	}
	return entry
}

// endOfSection is where control goes when the section of paragraph i ends,
// within a range that ends at last and goes on to exit: the start of the
// paragraph after the section's last, or exit if that is the range's
// last; -1 if the section ends past the range, or there is none.
func (b *builder) endOfSection(i, last, exit int, start func(int) int) int {
	paragraphs := b.program.Paragraphs
	section := paragraphs[i].Section
	if section == "" {
		return -1
	}
	end := i
	for end+1 < len(paragraphs) && paragraphs[end+1].Section == section {
		end++
	}
	switch {
	case end > last:
		return -1
	case end == last:
		return exit
	}
	return start(end + 1)
}

// scope is the range of paragraphs being laid out, and the nodes of their
// starts, which a GO TO goes to.
type scope struct {
	first, last int
	start       func(int) int
}

// sequence lays out statements in order, the last one going on to next,
// and returns the node control enters them by.
func (b *builder) sequence(r scope, statements []cobol.Statement, next int) int {
	for i := len(statements) - 1; i >= 0 && b.graph.Covered; i-- {
		next = b.statement(r, &statements[i], next)
	}
	return next
}

func (b *builder) statement(r scope, s *cobol.Statement, next int) int {
	n := b.node(s)
	switch s.Verb {
	case "STOP", "GOBACK", "EXIT PROGRAM":
		b.link(n, End)
	case "EXIT PARAGRAPH":
		b.link(n, b.paragraphEnd)
	case "EXIT SECTION":
		if b.sectionEnd < 0 {
			b.graph.Covered = false
		} else {
			b.link(n, b.sectionEnd)
		}
	case "EXIT PERFORM":
		b.graph.Covered = false
	case "GO":
		b.goTo(r, n, s, next)
	case "PERFORM":
		b.perform(r, n, s, next)
	case "SORT":
		b.sort(n, s, next)
	default:
		b.branches(r, n, s, next)
		for _, section := range b.program.ErrorProcedures(*s) {
			b.link(n, b.procedure(cobol.ProcedureRef{Name: section}, next))
		}
	}
	return n
}

func (b *builder) goTo(r scope, n int, s *cobol.Statement, next int) {
	for _, ref := range s.Procedures {
		i, _, ok := b.program.Span(cobol.ProcedureRef{Name: ref.Name})
		if !ok || i < r.first || i > r.last {
			b.graph.Covered = false
			continue
		}
		b.link(n, r.start(i))
	}
	if len(s.Procedures) > 1 {
		b.link(n, next) // GO TO DEPENDING ON, out of range
	}
}

// perform lays out an inline PERFORM's body or the range performed, run
// once, or round a loop that may run no times when the PERFORM repeats.
func (b *builder) perform(r scope, n int, s *cobol.Statement, next int) {
	loops := s.Condition != nil || len(s.Reads) > 0
	exit := next
	if loops {
		exit = n
		b.link(n, next)
	}
	if len(s.Procedures) == 0 {
		var body []cobol.Statement
		if len(s.Branches) > 0 {
			body = s.Branches[0].Statements
		}
		b.link(n, b.sequence(r, body, exit))
		return
	}
	b.link(n, b.procedure(s.Procedures[0], exit))
}

// sort performs a SORT's input procedure, then its output procedure.
func (b *builder) sort(n int, s *cobol.Statement, next int) {
	for i := len(s.Procedures) - 1; i >= 0; i-- {
		next = b.procedure(s.Procedures[i], next)
	}
	b.link(n, next)
}

// procedure lays out a performed paragraph, section or THRU range, going
// on to exit when it ends, and returns the node control enters it by.
func (b *builder) procedure(ref cobol.ProcedureRef, exit int) int {
	first, last, ok := b.program.Span(ref)
	if !ok || slices.ContainsFunc(b.active, func(a [2]int) bool { return a == [2]int{first, last} }) {
		b.graph.Covered = false
		return exit
	}
	return b.performed(first, last, exit)
}

// branches lays out the statements nested in another: each branch, and
// the way past them all unless they are a complete set, an IF with its
// ELSE, an EVALUATE with its WHEN OTHER, or a phrase with its NOT form.
func (b *builder) branches(r scope, n int, s *cobol.Statement, next int) {
	for i := range s.Branches {
		b.link(n, b.sequence(r, s.Branches[i].Statements, next))
	}
	if !complete(s.Branches) {
		b.link(n, next)
	}
}

func complete(branches []cobol.Branch) bool {
	kinds := map[string]bool{}
	for _, br := range branches {
		kinds[br.Kind] = true
	}
	if kinds["ELSE"] || kinds["OTHER"] {
		return true
	}
	for kind := range kinds {
		if kinds["NOT "+kind] {
			return true
		}
	}
	return false
}
