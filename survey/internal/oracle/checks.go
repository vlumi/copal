package oracle

import (
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
	"github.com/vlumi/copal/survey/internal/rules"
)

// checksOwnTotals reports whether a rule of a unit's table finds two
// totals the program adds up itself to differ: an equality, found false,
// of two items of its working storage that it only ever adds to, from a
// start of zero. No input field sets either, so the search has nothing to
// change to reach the rule: it takes the program's own arithmetic
// disagreeing with itself.
func (u Unit) checksOwnTotals(t rules.Table, r rules.Rule) bool {
	for i, entry := range r.Entries {
		if entry != "N" {
			continue
		}
		c := u.condition(t.Conditions[i])
		if c != nil && c.Op == "=" && c.Left.Name != "" && c.Right.Name != "" &&
			u.onlyAddedTo(c.Left.Name) && u.onlyAddedTo(c.Right.Name) {
			return true
		}
	}
	return false
}

// condition finds the condition a table's column tests, by its label.
func (u Unit) condition(label string) *cobol.Condition {
	label, _, _ = strings.Cut(label, ", again")
	var found *cobol.Condition
	var visit func(c *cobol.Condition)
	visit = func(c *cobol.Condition) {
		if found == nil && c.String() == label {
			found = c
		}
		for _, part := range c.Parts {
			visit(part)
		}
	}
	for _, c := range conditionsOf(u.source) {
		visit(c)
	}
	return found
}

// onlyAddedTo reports whether an item of the unit's working storage is
// only ever added to, from a start of zero: an ADD TO it, a MOVE of zero to
// it, or an INITIALIZE of it, at least one ADD among them.
func (u Unit) onlyAddedTo(item string) bool {
	var f *layout.Field
	for _, r := range u.records.WorkingStorage {
		if f == nil {
			f = find(r, item)
		}
	}
	if f == nil || len(f.Values) > 0 && !isZero(f.Values[0]) {
		return false
	}
	adds, others := 0, 0
	for _, para := range u.source.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if !slices.ContainsFunc(s.Writes, func(w cobol.Operand) bool { return w.Name == item }) {
				return
			}
			switch {
			case s.Verb == "ADD" && !strings.Contains(" "+s.Text+" ", " GIVING "):
				adds++
			case s.Verb == "INITIALIZE", s.Verb == "MOVE" && len(s.Reads) == 1 && isZero(s.Reads[0].Text):
			default:
				others++
			}
		})
	}
	return adds > 0 && others == 0
}

func isZero(text string) bool {
	return slices.Contains([]string{"0", "ZERO", "ZEROS", "ZEROES"}, text)
}
