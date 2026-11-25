package oracle

import (
	"slices"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/rules"
)

// failsLater reports whether a rule of a unit's table finds a file to
// have failed once it was open: it takes the file's status, as the
// statement before left it, to be none that the file's records make. A
// read in order leaves a record or the file's end, a read by key a record
// or none of that key, a write the record written or, by key, a key the
// file holds; and a rewrite or a delete of the record read, and a close,
// only succeed. An optional file opens whether it is there or not, where
// the open is not to make it new. Any other status is the device's or the
// file handler's doing, which the search does not make.
func (u Unit) failsLater(t rules.Table, r rules.Rule) bool {
	statements := statementsOf(u.source)
	for _, f := range u.blueprint.Files {
		status := u.statusOf(f)
		if status == nil {
			continue
		}
		tested := map[int][]int{} // the step that left the status, to the columns that test it
		for i, entry := range r.Entries {
			c := t.Tests[i].Condition
			if entry == "" || c == nil || t.Tests[i].Time > 1 {
				continue
			}
			if left, ok := u.leftBy(r.Steps, statements, c, status); ok {
				tested[left] = append(tested[left], i)
			}
		}
		for left, columns := range tested {
			left := statements[r.Steps[left]]
			made, ok := u.madeBy(f, left, status)
			if !ok {
				continue
			}
			takes := func(v string) bool {
				for _, i := range columns {
					settled := u.settledBefore(t, r, statements, i, status)
					if held := u.holds(t.Tests[i].Condition, status, v, settled); held == yes && r.Entries[i] == "N" || held == no && r.Entries[i] == "Y" {
						return false
					}
				}
				return true
			}
			if !slices.ContainsFunc(made, takes) {
				return true
			}
		}
	}
	return false
}

// statusOf is the item a file's FILE STATUS is kept in, if it has one.
func (u Unit) statusOf(f blueprint.File) *layout.Field {
	if f.Status == "" {
		return nil
	}
	for _, r := range u.records.All() {
		if status := find(r, f.Status); status != nil {
			return status
		}
	}
	return nil
}

// leftBy finds the step of a rule's path that left a file's status as a
// condition tests it: the last before the statement that tests the
// condition to operate on the file, itself or in what it performs.
func (u Unit) leftBy(steps []cobol.Pos, statements map[cobol.Pos]cobol.Statement, c *cobol.Condition, status *layout.Field) (int, bool) {
	if !slices.ContainsFunc(c.Refs(), func(ref string) bool { return u.item(ref) == status.Name }) {
		return 0, false
	}
	at := slices.IndexFunc(steps, func(pos cobol.Pos) bool { return tests(statements[pos], c) })
	for i := at - 1; i >= 0; i-- {
		if len(u.operations(statements[steps[i]], status)) > 0 {
			return i, true
		}
	}
	return 0, false
}

// item is the data item a name stands for: a condition name's field.
func (u Unit) item(name string) string {
	if f := conditionField(u.records, name); f != nil {
		return f.Name
	}
	return name
}

// tests reports whether a statement tests a condition, whole or as a part
// of the one it tests: an IF's, a PERFORM's UNTIL, a WHEN's. The
// condition is told by how it is written, a table's being of another
// reading of the source than the statements.
func tests(s cobol.Statement, c *cobol.Condition) bool {
	text := c.String()
	var within func(whole *cobol.Condition) bool
	within = func(whole *cobol.Condition) bool {
		return whole != nil && (whole.String() == text || slices.ContainsFunc(whole.Parts, within))
	}
	return within(s.Condition) || slices.ContainsFunc(s.Branches, func(b cobol.Branch) bool { return within(b.Condition) })
}

// operations lists the statements that set a file's status when a step
// of a path runs: the step, for one that operates on the file, and for a
// PERFORM, those of what it performs. The statements nested in an IF or
// an EVALUATE are steps of their own where the path takes them.
func (u Unit) operations(step cobol.Statement, status *layout.Field) []cobol.Statement {
	switch step.Verb {
	case "PERFORM", "OPEN", "CLOSE", "READ", "START", "DELETE", "WRITE", "REWRITE":
		return u.flow.Writers(step, []string{status.Name})
	}
	return nil
}

// madeBy lists the statuses a file's records can make a step leave: those
// of the operations it runs, narrowed, after a loop, to those its UNTIL
// ends on. It has no answer for a step that opens a file that has to be
// there, which the search does fail, or that sets the status any other
// way.
func (u Unit) madeBy(f blueprint.File, step cobol.Statement, status *layout.Field) ([]string, bool) {
	var made []string
	for _, s := range u.operations(step, status) {
		var statuses []string
		switch {
		case s.Verb == "READ" && (f.Access == "" || strings.Contains(" "+s.Text+" ", " NEXT ")):
			statuses = []string{"00", "10"}
		case s.Verb == "READ", s.Verb == "START":
			statuses = []string{"00", "23"}
		case s.Verb == "WRITE" && keyed(f):
			statuses = []string{"00", "22"}
		case s.Verb == "WRITE", s.Verb == "REWRITE", s.Verb == "DELETE", s.Verb == "CLOSE":
			statuses = []string{"00"}
		case s.Verb == "OPEN" && f.Optional && !opensNew(s, f):
			statuses = []string{"00", "05"}
		default:
			return nil, false
		}
		for _, v := range statuses {
			if !slices.Contains(made, v) {
				made = append(made, v)
			}
		}
	}
	if step.Verb == "PERFORM" && step.Condition != nil {
		made = slices.DeleteFunc(made, func(v string) bool { return u.holds(step.Condition, status, v, nil) == no })
	}
	return made, len(made) > 0
}

// opensNew reports whether an OPEN opens a file for output, to make it
// new, which fails where it cannot be made.
func opensNew(s cobol.Statement, f blueprint.File) bool {
	return slices.ContainsFunc(s.Files, func(ref cobol.FileRef) bool { return ref.Name == f.Name && ref.Mode == "OUTPUT" })
}

// held is whether a condition holds: yes, no, or unknown, for one that
// depends on more than the status.
type held int

const (
	unknown held = iota
	yes
	no
)

// settledBefore gives what a rule's other columns settle of the
// conditions a column tests beside the status: the outcome the rule has
// for a condition it tested before, with nothing since that may have
// changed what it refers to. A switch found on, tested again with the
// status, is still on.
func (u Unit) settledBefore(t rules.Table, r rules.Rule, statements map[cobol.Pos]cobol.Statement, column int, status *layout.Field) func(*cobol.Condition) held {
	testing := func(c *cobol.Condition) int {
		return slices.IndexFunc(r.Steps, func(pos cobol.Pos) bool { return tests(statements[pos], c) })
	}
	at := testing(t.Tests[column].Condition)
	return func(c *cobol.Condition) held {
		if slices.ContainsFunc(c.Refs(), func(ref string) bool { return u.item(ref) == status.Name }) {
			return unknown
		}
		for j, entry := range r.Entries {
			before := t.Tests[j].Condition
			if j == column || entry == "" || before == nil || t.Tests[j].Time > 1 || before.String() != c.String() {
				continue
			}
			from := testing(before)
			if from < 0 || from >= at || slices.ContainsFunc(r.Steps[from+1:at], func(pos cobol.Pos) bool {
				return u.flow.MayChange(statements[pos], c.Refs())
			}) {
				continue
			}
			return map[string]held{"Y": yes, "N": no}[entry]
		}
		return unknown
	}
}

// holds tells whether a condition holds when a file's status is a value,
// as far as the status tells, and what is settled of its other parts.
func (u Unit) holds(c *cobol.Condition, status *layout.Field, v string, settled func(*cobol.Condition) held) held {
	if settled != nil {
		if known := settled(c); known != unknown {
			return known
		}
	}
	is := func(literals ...string) held {
		if slices.ContainsFunc(literals, func(lit string) bool { return strings.Trim(lit, `'"`) == v }) {
			return yes
		}
		return no
	}
	switch c.Op {
	case "NAME":
		if f := conditionField(u.records, c.Left.Name); f == status {
			named := f.Conditions[slices.IndexFunc(f.Conditions, func(n layout.Condition) bool { return n.Name == c.Left.Name })]
			return is(named.Values...)
		}
	case "=":
		if c.Left.Name == status.Name && c.Right.Name == "" {
			return is(c.Right.Text)
		}
	case "NOT":
		return map[held]held{yes: no, no: yes}[u.holds(c.Parts[0], status, v, settled)]
	case "AND", "OR":
		parts := make([]held, len(c.Parts))
		for i, part := range c.Parts {
			parts[i] = u.holds(part, status, v, settled)
		}
		decides, other := no, yes // a part that is no settles an AND
		if c.Op == "OR" {
			decides, other = yes, no
		}
		switch {
		case slices.Contains(parts, decides):
			return decides
		case !slices.Contains(parts, unknown):
			return other
		}
	}
	return unknown
}

// statementsOf lists a program's statements by where they are.
func statementsOf(source *cobol.Program) map[cobol.Pos]cobol.Statement {
	statements := map[cobol.Pos]cobol.Statement{}
	for _, para := range source.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) { statements[s.Pos] = s })
	}
	return statements
}
