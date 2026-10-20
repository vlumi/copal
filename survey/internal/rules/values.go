package rules

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
)

// ConditionNames resolves a condition name: the field it is declared on,
// and the values that make it true, as written.
type ConditionNames func(name string) (field string, values []string, ok bool)

// Knowledge is what the rules know of the program's data: what a statement
// may change, and by which of the statements it runs; what a condition
// name stands for; the values fields of the working storage start with,
// by name; the fields of text a group holds; how many digits a field of
// digits has, and the other fields of digits over its bytes, through a
// redefinition too; and whether
// a MOVE from one field to another keeps its value whole.
type Knowledge struct {
	MayChange      MayChange
	Writers        func(s cobol.Statement, names []string) []cobol.Statement
	ConditionNames ConditionNames
	InitialValues  map[string]string
	TextFields     func(group string) []string
	DigitFields    func(field string) (digits int, over []DigitField)
	Copies         func(from, to string) bool
}

// DigitField is a field of digits over some of the digits of another:
// from offset, 0-based, for length.
type DigitField struct {
	Name           string
	Offset, Length int
}

// assignment is a value a statement gives a field: a literal moved into
// it, or a condition name's first value, set.
type assignment struct {
	field string
	value string
	from  string // the field whose value a MOVE copies, instead of a value
}

// assignments lists the values a statement gives fields, as events after
// the statement's change: MOVE of a literal into fields by their plain
// names, of SPACES into the fields of text of a group, and of a whole
// number into the fields of digits over a field of digits, a MOVE of a
// field that keeps its value whole, and SET of condition names to true.
func (e *extractor) assignments(s cobol.Statement) []event {
	var events []event
	switch {
	case s.Verb == "MOVE" && len(s.Reads) == 1 && s.Reads[0].Name != "" && s.Reads[0].Text == s.Reads[0].Name:
		for _, w := range s.Writes {
			if w.Name != "" && w.Text == w.Name && e.knowledge.Copies != nil && e.knowledge.Copies(s.Reads[0].Name, w.Name) {
				events = append(events, event{assign: &assignment{field: w.Name, from: s.Reads[0].Name}})
			}
		}
	case s.Verb == "MOVE" && len(s.Reads) == 1 && s.Reads[0].Name == "" && len(s.Reads[0].Refs) == 0:
		for _, w := range s.Writes {
			if w.Name == "" || w.Text != w.Name {
				continue
			}
			events = append(events, event{assign: &assignment{field: w.Name, value: s.Reads[0].Text}})
			if text(s.Reads[0].Text) == "" && e.knowledge.TextFields != nil {
				for _, f := range e.knowledge.TextFields(w.Name) {
					events = append(events, event{assign: &assignment{field: f, value: s.Reads[0].Text}})
				}
			}
			events = append(events, e.digitsOf(w.Name, s.Reads[0].Text)...)
		}
	case s.Verb == "SET" && len(s.Reads) == 0:
		for _, w := range s.Writes {
			if field, values, ok := e.names(w.Name); ok && len(values) > 0 {
				events = append(events, event{assign: &assignment{field: field, value: values[0]}})
			}
		}
	}
	return events
}

// digitsOf are the values a whole number moved into a field of digits
// gives the other fields of digits over its bytes: their digits of it.
func (e *extractor) digitsOf(field, literal string) []event {
	n, ok := number(literal)
	if !ok || n < 0 || n != float64(int64(n)) || e.knowledge.DigitFields == nil {
		return nil
	}
	width, over := e.knowledge.DigitFields(field)
	digits := fmt.Sprintf("%0*d", width, int64(n))
	if width == 0 || len(digits) != width {
		return nil
	}
	var events []event
	for _, d := range over {
		value := strings.TrimLeft(digits[d.Offset:d.Offset+d.Length], "0")
		if value == "" {
			value = "0"
		}
		events = append(events, event{assign: &assignment{field: d.Name, value: value}})
	}
	return events
}

func (e *extractor) names(name string) (string, []string, bool) {
	if e.knowledge.ConditionNames == nil {
		return "", nil, false
	}
	return e.knowledge.ConditionNames(name)
}

// known is what a path knows of a field at its end: the values it may
// have, or whether it does not know, and the values tests found it not to
// be, since. A field has the value it was last given, by a statement or by
// a test found true; after a statement that may change it, by statements
// that each give it a value, one of those or what it had before.
func (e *extractor) known(events []event, field string) (values []string, ok bool, not []string) {
	for i, ev := range slices.Backward(events) {
		switch {
		case ev.assign != nil && ev.assign.field == field && ev.assign.from != "":
			values, ok, _ := e.known(events[:i], ev.assign.from)
			return values, ok, not
		case ev.assign != nil && ev.assign.field == field:
			return []string{ev.assign.value}, true, not
		case ev.change != nil && e.knowledge.MayChange(*ev.change, []string{field}):
			given, gives := e.gives(*ev.change, field)
			before, known, _ := e.known(events[:i], field)
			if !gives || !known {
				return nil, false, not
			}
			return union(before, given), true, not
		case ev.test != nil && ev.test.cond != nil:
			if v, isValue := e.valueOfTest(ev.test, field); isValue && ev.test.outcome {
				return []string{v}, true, not
			} else if isValue {
				not = append(not, v)
			}
		}
	}
	return nil, false, not
}

// gives lists the values a statement may give a field, if each of the
// statements it runs that may change the field gives it a value.
func (e *extractor) gives(s cobol.Statement, field string) ([]string, bool) {
	if e.knowledge.Writers == nil {
		return nil, false
	}
	var values []string
	for _, w := range e.knowledge.Writers(s, []string{field}) {
		i := slices.IndexFunc(e.assignments(w), func(ev event) bool { return ev.assign.field == field })
		if i < 0 || e.assignments(w)[i].assign.from != "" {
			return nil, false
		}
		values = union(values, []string{e.assignments(w)[i].assign.value})
	}
	return values, true
}

// union adds the values of b that a does not have.
func union(a, b []string) []string {
	out := slices.Clone(a)
	for _, v := range b {
		if !slices.ContainsFunc(out, func(w string) bool { return sameValue(v, w) }) {
			out = append(out, v)
		}
	}
	return out
}

// possible are the values a path knows a field may have, less those
// tests found it not to be.
func possible(values, not []string) []string {
	return slices.DeleteFunc(slices.Clone(values), func(v string) bool {
		return slices.ContainsFunc(not, func(n string) bool { return sameValue(n, v) })
	})
}

// valueOfTest is the value a test compares a field with for equality, if
// it is one: FIELD = literal, or a condition name of one value on it.
func (e *extractor) valueOfTest(t *test, field string) (string, bool) {
	c := t.cond
	switch c.Op {
	case "=":
		if c.Left.Name == field && isLiteral(c.Right) {
			return c.Right.Text, true
		}
		if c.Right.Name == field && isLiteral(c.Left) {
			return c.Left.Text, true
		}
	case "NAME":
		if f, values, ok := e.names(c.Left.Text); ok && f == field && len(values) == 1 {
			return values[0], true
		}
	}
	return "", false
}

func isLiteral(op *cobol.Operand) bool { return op != nil && op.Name == "" && len(op.Refs) == 0 }

// evaluate works out a condition's outcome from what a path knows, if it
// can: the outcome a test of the same condition found, with nothing since
// that may have changed it; a comparison of fields whose values it knows, or with a literal, an
// equality with a literal the field is known not to be, a condition name on a field whose value it knows, or none of whose values
// the field can still have, and AND, OR and NOT of those.
func (e *extractor) evaluate(events []event, c *cobol.Condition) (outcome, ok bool) {
	if o, tested := e.tested(events, c); tested {
		return o, true
	}
	switch c.Op {
	case "AND", "OR":
		decided, all := c.Op == "OR", true
		for _, part := range c.Parts {
			o, known := e.evaluate(events, part)
			if known && o == decided {
				return decided, true
			}
			all = all && known
		}
		return !decided, all
	case "NOT":
		o, known := e.evaluate(events, c.Parts[0])
		return !o, known
	case "NAME":
		field, values, found := e.names(c.Left.Text)
		if !found {
			return false, false
		}
		has, known, not := e.known(events, field)
		if known {
			return agree(possible(has, not), func(v string) bool {
				return slices.ContainsFunc(values, func(w string) bool { return sameValue(v, w) })
			})
		}
		return false, len(possible(values, not)) == 0
	case "=", "<", ">", "<=", ">=":
		if o, same := e.sameItems(events, c); same {
			return o, true
		}
		left, lok := e.operandValues(events, c.Left)
		right, rok := e.operandValues(events, c.Right)
		if !lok || !rok {
			return false, c.Op == "=" && (e.excludes(events, c.Left, c.Right) || e.excludes(events, c.Right, c.Left))
		}
		return agreeAll(left, right, func(a, b string) (bool, bool) { return compare(c.Op, a, b) })
	case "TRUE":
		return true, true
	}
	return false, false
}

// sameItems decides a comparison of two fields by what the path knows
// they hold: the same value, when one was copied from the other, or from
// the same field, with neither changed since; or, for an equality, two
// values a test found one less than the other, with neither changed since.
func (e *extractor) sameItems(events []event, c *cobol.Condition) (outcome, ok bool) {
	if !isField(c.Left) || !isField(c.Right) {
		return false, false
	}
	a, b := e.alias(events, c.Left.Name), e.alias(events, c.Right.Name)
	switch {
	case a == b:
		return c.Op == "=" || c.Op == "<=" || c.Op == ">=", true
	case c.Op == "=" && e.ordered(events, a, b):
		return false, true
	}
	return false, false
}

// alias is the field whose value a field holds, as far as the path knows:
// the field it was last copied from, if neither has changed since, and so
// on; or itself.
func (e *extractor) alias(events []event, field string) string {
	for i, ev := range slices.Backward(events) {
		switch {
		case ev.assign != nil && ev.assign.field == field:
			if ev.assign.from != "" && !e.changedSince(events[i+1:], []string{ev.assign.from}) {
				return e.alias(events[:i], ev.assign.from)
			}
			return field
		case ev.change != nil && e.knowledge.MayChange(*ev.change, []string{field}):
			return field
		}
	}
	return field
}

// ordered reports whether a test of the path found one of two fields less
// than the other, with neither changed since.
func (e *extractor) ordered(events []event, a, b string) bool {
	for i, ev := range slices.Backward(events) {
		t := ev.test
		if t == nil || t.cond == nil || !t.outcome || t.cond.Op != "<" && t.cond.Op != ">" || !isField(t.cond.Left) || !isField(t.cond.Right) {
			continue
		}
		pair := []string{t.cond.Left.Name, t.cond.Right.Name}
		if slices.Contains(pair, a) && slices.Contains(pair, b) && !e.changedSince(events[i+1:], []string{a, b}) {
			return true
		}
	}
	return false
}

func isField(op *cobol.Operand) bool { return op != nil && op.Name != "" && op.Text == op.Name }

// tested is the outcome of the path's latest test of a condition, if
// nothing since may have changed what it refers to.
func (e *extractor) tested(events []event, c *cobol.Condition) (outcome, ok bool) {
	condition, refs := c.String(), c.Refs()
	for _, ev := range slices.Backward(events) {
		switch {
		case ev.change != nil && e.knowledge.MayChange(*ev.change, refs):
			return false, false
		case ev.test != nil && ev.test.cond != nil && ev.test.condition == condition:
			return ev.test.outcome, true
		}
	}
	return false, false
}

// excludes reports whether a path knows a field not to have a literal's
// value.
func (e *extractor) excludes(events []event, field, lit *cobol.Operand) bool {
	if !isLiteral(lit) || field == nil || field.Name == "" || field.Text != field.Name {
		return false
	}
	_, _, not := e.known(events, field.Name)
	return slices.ContainsFunc(not, func(n string) bool { return sameValue(n, lit.Text) })
}

// operandValues are the values an operand may have: a literal's, or
// those a path knows a field by its plain name may have.
func (e *extractor) operandValues(events []event, op *cobol.Operand) ([]string, bool) {
	switch {
	case isLiteral(op):
		return []string{op.Text}, true
	case op != nil && op.Name != "" && op.Text == op.Name:
		values, ok, not := e.known(events, op.Name)
		return possible(values, not), ok
	}
	return nil, false
}

// agree is the outcome a test has for every one of the values, if it has
// the same for all.
func agree(values []string, test func(string) bool) (outcome, ok bool) {
	return agreeAll(values, []string{""}, func(v, _ string) (bool, bool) { return test(v), true })
}

// agreeAll is the outcome a comparison has for every pair of the values,
// if it decides each the same way.
func agreeAll(left, right []string, cmp func(a, b string) (bool, bool)) (outcome, ok bool) {
	first := true
	for _, a := range left {
		for _, b := range right {
			o, decided := cmp(a, b)
			if !decided || !first && o != outcome {
				return false, false
			}
			outcome, first = o, false
		}
	}
	return outcome, !first
}

// compare compares two values as written: as numbers when both are, as
// text otherwise, which only equality decides.
func compare(op, a, b string) (bool, bool) {
	x, xok := number(a)
	y, yok := number(b)
	if xok && yok {
		switch op {
		case "=":
			return x == y, true
		case "<":
			return x < y, true
		case ">":
			return x > y, true
		case "<=":
			return x <= y, true
		}
		return x >= y, true
	}
	if op == "=" {
		return text(a) == text(b), true
	}
	return false, false
}

// sameValue reports whether two values as written are the same value.
func sameValue(a, b string) bool {
	same, _ := compare("=", a, b)
	return same
}

func number(v string) (float64, bool) {
	switch v {
	case "ZERO", "ZEROS", "ZEROES":
		return 0, true
	}
	n, err := strconv.ParseFloat(strings.Trim(v, `'"`), 64)
	return n, err == nil && !strings.HasPrefix(v, "'") && !strings.HasPrefix(v, `"`)
}

func text(v string) string {
	switch v {
	case "SPACE", "SPACES":
		return ""
	case "ZERO", "ZEROS", "ZEROES":
		return "0"
	}
	if len(v) >= 2 && (v[0] == '\'' || v[0] == '"') {
		v = v[1 : len(v)-1]
	}
	return strings.TrimRight(v, " ")
}
