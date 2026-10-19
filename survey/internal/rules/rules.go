// Package rules extracts a program's business rules as decision tables:
// for each paragraph, every path through its decisions, with the
// conditions the path tests and the actions it takes.
//
// IF, EVALUATE and SEARCH decide, and so do the conditional phrases of
// READ, RETURN and the like: AT END, INVALID KEY, ON SIZE ERROR. The WHENs
// of an EVALUATE are tried in order, so a path through one is also a path
// past those before it; an EVALUATE without WHEN OTHER has a path for when
// nothing matches. A path ends at a statement that ends unconditionally,
// such as GO TO or STOP RUN, or at EXIT PARAGRAPH or EXIT SECTION, so
// that in the IF ...
// GO TO chains of old code the rest of the paragraph runs only on the
// paths that did not jump.
// The body of an inline PERFORM is a table of its own, for one pass
// through the loop, which starts knowing the loop's UNTIL condition false
// unless the loop tests it after each pass.
package rules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/survey/internal/structure"
)

// Table is a decision table: the conditions tested, with the data items
// and condition names each refers to, the actions taken, and the rules,
// one per path.
type Table struct {
	Paragraph  string
	Loop       string // the inline PERFORM whose body the table is, if any
	Pos        cobol.Pos
	Conditions []string
	Refs       [][]string // parallel to Conditions
	Actions    []string
	Rules      []Rule
	// Collapsed are the statements inside decisions shown as actions, past
	// maxPaths, whose paths the rules do not tell apart.
	Collapsed []cobol.Pos
}

// Rule is one path: the outcome of each condition, Y, N or empty when the
// path does not test it, the actions it takes, in order, as indexes into
// the table's actions, and the statements it runs, in order, by position,
// those inside the table's inline loops and collapsed decisions left out.
type Rule struct {
	Entries []string
	Actions []int
	Steps   []cobol.Pos
}

// maxPaths caps the paths of a table: a decision that would take a
// paragraph beyond it is shown as an action instead.
const maxPaths = 64

// Tables extracts the decision tables of a program, for each paragraph
// and each inline loop that decides anything. A condition tested again
// after a statement that may change what it refers to is a test of its
// own, and a path that would find a condition other than the values it
// knows decide is no rule: one that moves an error message into a field
// does not then find the field's condition name for no error true. A
// paragraph's paths start from what holds wherever it is entered.
func Tables(p *cobol.Program, knowledge Knowledge) []Table {
	var tables []Table
	g := newProgram(p, knowledge)
	for i, para := range p.Paragraphs {
		e := g.follow(i)
		if t, ok := newTable(para.Name, "", para.Pos, e.top, e.collapsed); ok {
			tables = append(tables, t)
		}
		tables = append(tables, e.loops...)
	}
	return tables
}

// test is a condition, the names it refers to, and its outcome on a path.
// A path that tests a condition again, after a statement that may have
// changed what it refers to, tests a new instance of it.
type test struct {
	condition string
	refs      []string
	outcome   bool
	instance  int
	cond      *cobol.Condition // nil for a conditional phrase
	// given is a test the path knows the outcome of when it starts, which
	// its table does not show.
	given bool
}

// label names a test's column: the condition, and for a test of it
// again, which time.
func (t test) label() string {
	switch t.instance {
	case 0, 1:
		return t.condition
	case 2:
		return t.condition + ", again"
	}
	return fmt.Sprintf("%s, again (%d)", t.condition, t.instance)
}

// MayChange reports whether a statement may change any of the data items
// named, which a condition refers to.
type MayChange func(s cobol.Statement, names []string) bool

// event is a test, a statement that may change what a later test refers
// to, or a value a statement gives a field, on a path, in the order the
// path meets them.
type event struct {
	test   *test
	change *cobol.Statement
	assign *assignment
}

// path is one way through a list of statements: the tests and changes it
// meets, and the actions it takes. An ended path has left them, by a jump
// or by stopping the run, and takes no more actions.
type path struct {
	events  []event
	actions []string
	steps   []cobol.Pos
	ended   bool
}

// tests lists the tests a path makes, in order.
func (p path) tests() []test {
	var tests []test
	for _, e := range p.events {
		if e.test != nil && !e.test.given {
			tests = append(tests, *e.test)
		}
	}
	return tests
}

// extractor follows the paths of one paragraph, collecting the tables of
// its inline loops on the way: top are the paths through the paragraph,
// and bodies those through each pass of its inline loops.
type extractor struct {
	paragraph string
	top       []path
	bodies    [][]path
	loops     []Table
	knowledge Knowledge
	collapsed []cobol.Pos
	// never are the paragraphs that never return to a PERFORM of them.
	never map[string]bool
}

// paths lists the ways through a list of statements.
func (e *extractor) paths(statements []cobol.Statement) []path {
	return e.pathsKnowing(nil, statements)
}

// pathsKnowing lists the ways through a list of statements, from a start
// that knows the outcomes of some tests.
func (e *extractor) pathsKnowing(given []event, statements []cobol.Statement) []path {
	paths := []path{{events: given}}
	for _, s := range statements {
		var alternatives []path
		switch {
		case s.Verb == "IF":
			var otherwise []cobol.Statement
			if len(s.Branches) > 1 {
				otherwise = s.Branches[1].Statements
			}
			alternatives = append(e.under(s.Condition, true, s.Branches[0].Statements), e.under(s.Condition, false, otherwise)...)
		case s.Verb == "EVALUATE":
			alternatives = e.cases(s)
		case s.Verb == "SEARCH":
			alternatives = e.cross([]path{{actions: []string{s.Text}, events: changes(s)}}, e.cases(s))
		case s.Verb == "PERFORM" && len(s.Branches) > 0:
			outer := e.collapsed
			e.collapsed = nil
			body := e.pathsKnowing(passing(s), s.Branches[0].Statements)
			e.bodies = append(e.bodies, body)
			if t, ok := newTable(e.paragraph, s.Text, s.Pos, body, e.collapsed); ok {
				e.loops = append(e.loops, t)
			}
			e.collapsed = outer
		case len(s.Branches) > 0:
			alternatives = e.phrases(s)
		}
		action := []path{{actions: []string{s.Text}, events: append(changes(s), e.assignments(s)...), ended: e.endsPath(s)}}
		if alternatives == nil {
			alternatives = action
		}
		next := e.cross(paths, startingAt(s.Pos, alternatives))
		if len(next) > maxPaths {
			next = e.cross(paths, startingAt(s.Pos, action))
			for _, b := range s.Branches {
				cobol.Walk(b.Statements, func(n cobol.Statement) { e.collapsed = append(e.collapsed, n.Pos) })
			}
		}
		paths = next
	}
	return paths
}

// startingAt puts a statement at the start of the paths through it.
func startingAt(pos cobol.Pos, paths []path) []path {
	out := make([]path, len(paths))
	for i, p := range paths {
		p.steps = append([]cobol.Pos{pos}, p.steps...)
		out[i] = p
	}
	return out
}

// endsPath reports whether a path ends at a statement: one that ends
// unconditionally, EXIT PARAGRAPH or EXIT SECTION, which leave the
// paragraph, or a PERFORM of a paragraph that never returns.
func (e *extractor) endsPath(s cobol.Statement) bool {
	return s.Verb == "EXIT PARAGRAPH" || s.Verb == "EXIT SECTION" || structure.EndsUnconditionally([]cobol.Statement{s}) || performsNeverReturning(s, e.never)
}

// performsNeverReturning reports whether a statement performs, once, a
// paragraph that never returns.
func performsNeverReturning(s cobol.Statement, never map[string]bool) bool {
	return s.Verb == "PERFORM" && len(s.Branches) == 0 && s.Condition == nil && len(s.Reads) == 0 &&
		len(s.Procedures) == 1 && s.Procedures[0].Thru == "" && never[s.Procedures[0].Name]
}

// neverReturning finds the paragraphs that never return to a PERFORM of
// them: those whose every path ends in STOP RUN, GOBACK or EXIT PROGRAM,
// or in a PERFORM of a paragraph that never returns.
func neverReturning(p *cobol.Program) map[string]bool {
	never := map[string]bool{}
	for changed := true; changed; {
		changed = false
		for _, para := range p.Paragraphs {
			if !never[para.Name] && stops(para.Statements, never) {
				never[para.Name], changed = true, true
			}
		}
	}
	return never
}

// stops reports whether every path through the statements ends the run.
func stops(statements []cobol.Statement, never map[string]bool) bool {
	if len(statements) == 0 {
		return false
	}
	last := statements[len(statements)-1]
	switch last.Verb {
	case "STOP", "GOBACK", "EXIT PROGRAM":
		return true
	case "IF", "EVALUATE":
		return complete(last.Branches) && !slices.ContainsFunc(last.Branches, func(b cobol.Branch) bool { return !stops(b.Statements, never) })
	}
	return performsNeverReturning(last, never)
}

// complete reports whether a set of branches is taken whatever the
// outcome: an IF with its ELSE, an EVALUATE with its WHEN OTHER.
func complete(branches []cobol.Branch) bool {
	return slices.ContainsFunc(branches, func(b cobol.Branch) bool { return b.Kind == "ELSE" || b.Kind == "OTHER" })
}

// changes is a statement as an event that may change what later tests
// refer to.
func changes(s cobol.Statement) []event {
	return []event{{change: &s}}
}

// passing is what a pass through a loop knows when it starts: that the
// PERFORM's UNTIL condition is false, unless it is tested after each pass.
func passing(s cobol.Statement) []event {
	if s.Condition == nil || s.TestAfter {
		return nil
	}
	given := testOf(s.Condition, false)
	given.given = true
	return []event{{test: given}}
}

// under prefixes the paths through statements with a condition's outcome.
func (e *extractor) under(c *cobol.Condition, outcome bool, statements []cobol.Statement) []path {
	return e.cross([]path{{events: []event{{test: testOf(c, outcome)}}}}, e.paths(statements))
}

// testOf is a test of a condition with an outcome. NOT is turned around:
// the outcome of NOT A is the opposite of A's.
func testOf(c *cobol.Condition, outcome bool) *test {
	for c.Op == "NOT" {
		c, outcome = c.Parts[0], !outcome
	}
	return &test{condition: c.String(), refs: c.Refs(), outcome: outcome, cond: c}
}

// cases lists the paths through the WHENs of an EVALUATE or SEARCH: each
// WHEN's, past the WHENs before it, then WHEN OTHER or AT END, or nothing
// when no WHEN matches.
func (e *extractor) cases(s cobol.Statement) []path {
	var paths []path
	past := []path{{}}
	var otherwise []cobol.Statement
	for _, b := range s.Branches {
		if b.Condition == nil {
			otherwise = b.Statements
			continue
		}
		paths = append(paths, e.cross(past, e.under(b.Condition, true, b.Statements))...)
		past = e.cross(past, e.under(b.Condition, false, nil))
	}
	return append(paths, e.cross(past, e.paths(otherwise))...)
}

// phrases lists the paths through a statement's conditional phrases, such
// as AT END and NOT AT END: the statement runs, then one of them, or, with
// only one, that one or nothing.
func (e *extractor) phrases(s cobol.Statement) []path {
	var paths []path
	for _, b := range s.Branches {
		kind, negated := strings.CutPrefix(b.Kind, "NOT ")
		start := []path{{events: phraseEvents(s, kind, !negated), actions: []string{s.Text}}}
		paths = append(paths, e.cross(start, e.paths(b.Statements))...)
		if len(s.Branches) == 1 {
			paths = append(paths, path{events: phraseEvents(s, kind, negated), actions: []string{s.Text}})
		}
	}
	return paths
}

// phraseEvents is a statement running and its phrase's condition, such
// as AT END, being raised or not, then what the statement changes.
func phraseEvents(s cobol.Statement, kind string, outcome bool) []event {
	return append([]event{{test: &test{condition: s.Text + " " + kind, outcome: outcome}}}, changes(s)...)
}

// cross follows each path with each of the next, leaving out those that
// would test one condition both ways with nothing between the tests that
// may change it. An ended path stays as it is.
func (e *extractor) cross(paths, next []path) []path {
	var out []path
	for _, p := range paths {
		if p.ended {
			out = append(out, p)
			continue
		}
		for _, n := range next {
			if events, ok := e.follow(p.events, n.events); ok {
				out = append(out, path{events: events, actions: append(slices.Clone(p.actions), n.actions...),
					steps: append(slices.Clone(p.steps), n.steps...), ended: n.ended})
			}
		}
	}
	return out
}

// follow appends a path's events to those of the path before it. A test
// whose outcome the values the path knows decide must have that outcome. A
// test of a condition the path before tested, with no change between that
// may have changed it, must have the same outcome, and is the same test,
// unless the path was given the outcome; a test after such a change is a
// new instance.
func (e *extractor) follow(before, after []event) ([]event, bool) {
	events := slices.Clone(before)
	for _, ev := range after {
		if ev.test == nil {
			events = append(events, ev)
			continue
		}
		t := *ev.test
		if t.cond != nil {
			if outcome, known := e.evaluate(events, t.cond); known && outcome != t.outcome {
				return nil, false
			}
		}
		last, changed := e.lastTest(events, t)
		switch {
		case last != nil && !changed && last.outcome != t.outcome:
			return nil, false
		case last != nil && !changed && !last.given:
			continue
		case last != nil:
			t.instance = last.instance + 1
		default:
			t.instance = 1
		}
		events = append(events, event{test: &t})
		if !e.holds(events) {
			return nil, false
		}
	}
	return events, true
}

// holds reports whether the tests of a path that join conditions with AND
// or OR still have their outcomes, by what the path has learned since: a
// test whose condition nothing since may have changed, and each part of an
// AND found true, or of an OR found false, that nothing since may have
// changed, which has the same outcome. An AND found true, one of whose
// parts a later test finds false, is no path.
func (e *extractor) holds(events []event) bool {
	for i, ev := range events {
		t, later := ev.test, events[i+1:]
		others := slices.Delete(slices.Clone(events), i, i+1)
		contradicts := func(c *cobol.Condition, outcome bool) bool {
			o, known := e.evaluate(others, c)
			return known && o != outcome
		}
		if t == nil || t.cond == nil || t.cond.Op != "AND" && t.cond.Op != "OR" {
			continue
		}
		if !e.changedSince(later, t.refs) {
			if contradicts(t.cond, t.outcome) {
				return false
			}
			continue
		}
		if (t.cond.Op == "AND") != t.outcome {
			continue
		}
		for _, part := range t.cond.Parts {
			if !e.changedSince(later, part.Refs()) && contradicts(part, t.outcome) {
				return false
			}
		}
	}
	return true
}

// changedSince reports whether any of the events may change the data
// items named.
func (e *extractor) changedSince(events []event, refs []string) bool {
	return slices.ContainsFunc(events, func(ev event) bool { return ev.change != nil && e.knowledge.MayChange(*ev.change, refs) })
}

// lastTest finds the latest test of the condition among the events, and
// whether a change after it may have changed what the condition refers to.
func (e *extractor) lastTest(events []event, t test) (*test, bool) {
	changed := false
	for _, ev := range slices.Backward(events) {
		switch {
		case ev.change != nil && !changed && e.knowledge.MayChange(*ev.change, t.refs):
			changed = true
		case ev.test != nil && ev.test.condition == t.condition:
			return ev.test, changed
		}
	}
	return nil, false
}

// newTable lays paths out as a table, if any of them tests a condition.
func newTable(paragraph, loop string, pos cobol.Pos, paths []path, collapsed []cobol.Pos) (Table, bool) {
	t := Table{Paragraph: paragraph, Loop: loop, Pos: pos, Collapsed: collapsed}
	for _, p := range paths {
		for _, test := range p.tests() {
			if !slices.Contains(t.Conditions, test.label()) {
				t.Conditions = append(t.Conditions, test.label())
				t.Refs = append(t.Refs, test.refs)
			}
		}
		for _, a := range p.actions {
			if !slices.Contains(t.Actions, a) {
				t.Actions = append(t.Actions, a)
			}
		}
	}
	if len(t.Conditions) == 0 {
		return Table{}, false
	}
	for _, p := range paths {
		r := Rule{Entries: make([]string, len(t.Conditions)), Steps: p.steps}
		for _, test := range p.tests() {
			entry := "N"
			if test.outcome {
				entry = "Y"
			}
			r.Entries[slices.Index(t.Conditions, test.label())] = entry
		}
		for _, a := range p.actions {
			r.Actions = append(r.Actions, slices.Index(t.Actions, a))
		}
		t.Rules = append(t.Rules, r)
	}
	return t, true
}
