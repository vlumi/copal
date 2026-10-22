package rules

import (
	"maps"
	"slices"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/survey/internal/structure"
)

// program follows the paths of a program's paragraphs, each once, from
// what holds wherever the paragraph is entered: at the start of the run,
// the values the working storage starts with, for the first paragraph of
// a program entered once a run; and for a paragraph entered by PERFORM
// alone, the outcomes every PERFORM of it knows.
type program struct {
	source    *cobol.Program
	knowledge Knowledge
	never     map[string]bool
	start     int
	runOnce   bool
	performed map[int]bool
	followed  map[int]*extractor
	following map[int]bool
}

func newProgram(p *cobol.Program, knowledge Knowledge) *program {
	g := &program{source: p, knowledge: knowledge, never: neverReturning(p),
		followed: map[int]*extractor{}, following: map[int]bool{}}
	g.start = slices.IndexFunc(p.Paragraphs, func(para cobol.Paragraph) bool { return !para.Declarative })
	g.performed, g.runOnce = g.entries()
	return g
}

// follow follows the paths of a paragraph, once.
func (g *program) follow(i int) *extractor {
	if e, ok := g.followed[i]; ok {
		return e
	}
	g.following[i] = true
	given := g.entry(i)
	para := g.source.Paragraphs[i]
	e := &extractor{paragraph: para.Name, knowledge: g.knowledge, never: g.never}
	e.top = e.pathsKnowing(given, para.Statements)
	delete(g.following, i)
	g.followed[i] = e
	return e
}

// entry is what a paragraph knows when it starts.
func (g *program) entry(i int) []event {
	switch {
	case i == g.start && g.runOnce:
		return g.initialValues()
	case g.performed[i]:
		return g.atPerforms(i)
	}
	return nil
}

// initialValues are the values the working storage starts with, as
// statements had given them.
func (g *program) initialValues() []event {
	var events []event
	for _, field := range slices.Sorted(maps.Keys(g.knowledge.InitialValues)) {
		events = append(events, event{assign: &assignment{field: field, value: g.knowledge.InitialValues[field]}})
	}
	return events
}

// atPerforms is what a paragraph knows at every PERFORM of it: the
// outcome of each condition it tests that every path to every PERFORM of
// it knows the same. A PERFORM that loops knows only that its UNTIL
// condition is false; a PERFORM a path does not show, inside a decision
// collapsed or after one no path takes, leaves nothing known.
func (g *program) atPerforms(i int) []event {
	name := g.source.Paragraphs[i].Name
	var known [][]event
	for j, caller := range g.source.Paragraphs {
		var sites []cobol.Statement
		cobol.Walk(caller.Statements, func(s cobol.Statement) {
			if s.Verb == "PERFORM" && slices.ContainsFunc(s.Procedures, func(r cobol.ProcedureRef) bool { return r.Name == name }) {
				sites = append(sites, s)
			}
		})
		if len(sites) == 0 {
			continue
		}
		if g.following[j] {
			return nil
		}
		e := g.follow(j)
		for _, s := range sites {
			if s.Condition != nil || len(s.Reads) > 0 {
				known = append(known, passing(s))
				continue
			}
			before := e.before(s)
			if len(before) == 0 {
				return nil
			}
			known = append(known, before...)
		}
	}
	return g.agreed(i, known)
}

// before lists what each path of a paragraph knows just before a
// statement it runs.
func (e *extractor) before(s cobol.Statement) [][]event {
	var known [][]event
	for _, paths := range append([][]path{e.top}, e.bodies...) {
		for _, p := range paths {
			if k := slices.IndexFunc(p.events, func(ev event) bool { return ev.change != nil && ev.change.Pos == s.Pos }); k >= 0 {
				known = append(known, p.events[:k])
			}
		}
	}
	return known
}

// agreed gives, as tests a paragraph's paths start with, the outcome of
// each condition it tests, or part of one, that all of the known agree
// on; and of those the paragraphs it performs test, which carry on to
// them.
func (g *program) agreed(i int, known [][]event) []event {
	if len(known) == 0 {
		return nil
	}
	x := &extractor{knowledge: g.knowledge}
	var given []event
	for _, c := range conditionsIn(g.statementsReached(i)) {
		outcome, ok := x.evaluate(known[0], c)
		if ok && !slices.ContainsFunc(known[1:], func(events []event) bool {
			o, k := x.evaluate(events, c)
			return !k || o != outcome
		}) {
			t := testOf(c, outcome)
			t.given = true
			given = append(given, event{test: t})
		}
	}
	return given
}

// statementsReached lists a paragraph's statements and those of the
// paragraphs it performs, through them, each paragraph once.
func (g *program) statementsReached(i int) []cobol.Statement {
	index := map[string]int{}
	for j, para := range g.source.Paragraphs {
		index[para.Name] = j
	}
	var statements []cobol.Statement
	seen := map[int]bool{}
	var add func(i int)
	add = func(i int) {
		if seen[i] {
			return
		}
		seen[i] = true
		statements = append(statements, g.source.Paragraphs[i].Statements...)
		cobol.Walk(g.source.Paragraphs[i].Statements, func(s cobol.Statement) {
			if s.Verb != "PERFORM" {
				return
			}
			for _, r := range s.Procedures {
				from, ok := index[r.Name]
				if !ok {
					continue
				}
				to, ok := index[r.Thru]
				if !ok {
					to = from
				}
				for k := from; k <= to; k++ {
					add(k)
				}
			}
		})
	}
	add(i)
	return statements
}

// conditionsIn lists the conditions statements test, and their parts,
// each once.
func conditionsIn(statements []cobol.Statement) []*cobol.Condition {
	var conditions []*cobol.Condition
	seen := map[string]bool{}
	var add func(c *cobol.Condition)
	add = func(c *cobol.Condition) {
		if c == nil || c.Op == "TRUE" {
			return
		}
		for c.Op == "NOT" {
			c = c.Parts[0]
		}
		if !seen[c.String()] {
			seen[c.String()] = true
			conditions = append(conditions, c)
		}
		for _, part := range c.Parts {
			add(part)
		}
	}
	cobol.Walk(statements, func(s cobol.Statement) {
		add(s.Condition)
		for _, b := range s.Branches {
			add(b.Condition)
		}
	})
	return conditions
}

// entries finds the paragraphs entered by PERFORM alone, and whether the
// program is entered once a run, at its first paragraph. A paragraph is
// entered by PERFORM alone when nothing names it but a PERFORM of it
// alone, no PERFORM THRU spans it, it is in no section, and the paragraph
// before it is one too, or never runs into it. A program is entered once a
// run when nothing calls it with items, it never returns to a caller, and
// nothing names its first paragraph.
func (g *program) entries() (map[int]bool, bool) {
	paragraphs := g.source.Paragraphs
	index := map[string]int{}
	named := map[int]bool{}
	spanned := map[int]bool{}
	returns := false
	for i, para := range paragraphs {
		if _, dup := index[para.Name]; dup {
			spanned[index[para.Name]], spanned[i] = true, true
		}
		index[para.Name] = i
	}
	for _, para := range paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			returns = returns || s.Verb == "GOBACK" || s.Verb == "EXIT PROGRAM"
			for _, r := range s.Procedures {
				from, ok := index[r.Name]
				if !ok {
					continue
				}
				named[from] = true
				if s.Verb == "PERFORM" && r.Thru == "" {
					continue
				}
				to, ok := index[r.Thru]
				if !ok {
					to = from
				}
				for k := from; k <= to; k++ {
					spanned[k] = true
				}
			}
		})
	}
	performed := map[int]bool{}
	for i, para := range paragraphs {
		if i <= g.start || spanned[i] || para.Section != "" || para.Declarative {
			continue
		}
		prev := paragraphs[i-1]
		performed[i] = performed[i-1] || structure.EndsUnconditionally(prev.Statements) || stops(prev.Statements, g.never)
	}
	runOnce := g.start >= 0 && len(g.source.Using) == 0 && !returns && !named[g.start]
	return performed, runOnce
}
