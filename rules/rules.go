// Package rules extracts a program's business rules as decision tables:
// for each paragraph, every path through its decisions, with the
// conditions the path tests and the actions it takes.
//
// IF, EVALUATE and SEARCH decide, and so do the conditional phrases of
// READ, RETURN and the like: AT END, INVALID KEY, ON SIZE ERROR. The WHENs
// of an EVALUATE are tried in order, so a path through one is also a path
// past those before it; an EVALUATE without WHEN OTHER has a path for when
// nothing matches. A path ends at a statement that ends unconditionally,
// such as GO TO or STOP RUN, so that in the IF ... GO TO chains of old
// code the rest of the paragraph runs only on the paths that did not jump.
// The body of an inline PERFORM is a table of its own, for one pass
// through the loop.
package rules

import (
	"slices"
	"strings"

	"github.com/vlumi/copal/cobol"
	"github.com/vlumi/copal/structure"
)

// Table is a decision table: the conditions tested, the actions taken,
// and the rules, one per path.
type Table struct {
	Paragraph  string
	Loop       string // the inline PERFORM whose body the table is, if any
	Pos        cobol.Pos
	Conditions []string
	Actions    []string
	Rules      []Rule
}

// Rule is one path: the outcome of each condition, Y, N or empty when the
// path does not test it, and the actions it takes, in order, as indexes
// into the table's actions.
type Rule struct {
	Entries []string
	Actions []int
}

// maxPaths caps the paths of a table: a decision that would take a
// paragraph beyond it is shown as an action instead.
const maxPaths = 64

// Tables extracts the decision tables of a program, for each paragraph
// and each inline loop that decides anything.
func Tables(p *cobol.Program) []Table {
	var tables []Table
	for _, para := range p.Paragraphs {
		e := &extractor{paragraph: para.Name}
		if t, ok := newTable(para.Name, "", para.Pos, e.paths(para.Statements)); ok {
			tables = append(tables, t)
		}
		tables = append(tables, e.loops...)
	}
	return tables
}

// test is a condition and its outcome on a path.
type test struct {
	condition string
	outcome   bool
}

// path is one way through a list of statements. An ended path has left
// them, by a jump or by stopping the run, and takes no more actions.
type path struct {
	tests   []test
	actions []string
	ended   bool
}

// extractor follows the paths of one paragraph, collecting the tables of
// its inline loops on the way.
type extractor struct {
	paragraph string
	loops     []Table
}

// paths lists the ways through a list of statements.
func (e *extractor) paths(statements []cobol.Statement) []path {
	paths := []path{{}}
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
			alternatives = cross([]path{{actions: []string{s.Text}}}, e.cases(s))
		case s.Verb == "PERFORM" && len(s.Branches) > 0:
			if t, ok := newTable(e.paragraph, s.Text, s.Pos, e.paths(s.Branches[0].Statements)); ok {
				e.loops = append(e.loops, t)
			}
		case len(s.Branches) > 0:
			alternatives = e.phrases(s)
		}
		action := []path{{actions: []string{s.Text}, ended: structure.EndsUnconditionally([]cobol.Statement{s})}}
		if alternatives == nil {
			alternatives = action
		}
		next := cross(paths, alternatives)
		if len(next) > maxPaths {
			next = cross(paths, action)
		}
		paths = next
	}
	return paths
}

// under prefixes the paths through statements with a condition's outcome.
// NOT is turned around: the outcome of NOT A is the opposite of A's.
func (e *extractor) under(c *cobol.Condition, outcome bool, statements []cobol.Statement) []path {
	for c.Op == "NOT" {
		c, outcome = c.Parts[0], !outcome
	}
	return cross([]path{{tests: []test{{c.String(), outcome}}}}, e.paths(statements))
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
		paths = append(paths, cross(past, e.under(b.Condition, true, b.Statements))...)
		past = cross(past, e.under(b.Condition, false, nil))
	}
	return append(paths, cross(past, e.paths(otherwise))...)
}

// phrases lists the paths through a statement's conditional phrases, such
// as AT END and NOT AT END: the statement runs, then one of them, or, with
// only one, that one or nothing.
func (e *extractor) phrases(s cobol.Statement) []path {
	var paths []path
	for _, b := range s.Branches {
		kind, negated := strings.CutPrefix(b.Kind, "NOT ")
		start := []path{{tests: []test{{s.Text + " " + kind, !negated}}, actions: []string{s.Text}}}
		paths = append(paths, cross(start, e.paths(b.Statements))...)
		if len(s.Branches) == 1 {
			paths = append(paths, path{tests: []test{{s.Text + " " + kind, negated}}, actions: []string{s.Text}})
		}
	}
	return paths
}

// cross follows each path with each of the next, leaving out those that
// would test one condition both ways. An ended path stays as it is.
func cross(paths, next []path) []path {
	var out []path
	for _, p := range paths {
		if p.ended {
			out = append(out, p)
			continue
		}
	nextPath:
		for _, n := range next {
			tests := slices.Clone(p.tests)
			for _, t := range n.tests {
				i := slices.IndexFunc(tests, func(u test) bool { return u.condition == t.condition })
				switch {
				case i < 0:
					tests = append(tests, t)
				case tests[i].outcome != t.outcome:
					continue nextPath
				}
			}
			out = append(out, path{tests: tests, actions: append(slices.Clone(p.actions), n.actions...), ended: n.ended})
		}
	}
	return out
}

// newTable lays paths out as a table, if any of them tests a condition.
func newTable(paragraph, loop string, pos cobol.Pos, paths []path) (Table, bool) {
	t := Table{Paragraph: paragraph, Loop: loop, Pos: pos}
	for _, p := range paths {
		for _, test := range p.tests {
			if !slices.Contains(t.Conditions, test.condition) {
				t.Conditions = append(t.Conditions, test.condition)
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
		r := Rule{Entries: make([]string, len(t.Conditions))}
		for _, test := range p.tests {
			entry := "N"
			if test.outcome {
				entry = "Y"
			}
			r.Entries[slices.Index(t.Conditions, test.condition)] = entry
		}
		for _, a := range p.actions {
			r.Actions = append(r.Actions, slices.Index(t.Actions, a))
		}
		t.Rules = append(t.Rules, r)
	}
	return t, true
}
