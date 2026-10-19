package rules

import (
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
)

// ConditionNames resolves a condition name: the field it is declared on,
// and the values that make it true, as written.
type ConditionNames func(name string) (field string, values []string, ok bool)

// Knowledge is what the rules know of the program's data: what a statement
// may change, what a condition name stands for, and the values fields of
// the working storage start with, by name.
type Knowledge struct {
	MayChange      MayChange
	ConditionNames ConditionNames
	InitialValues  map[string]string
}

// assignment is a value a statement gives a field: a literal moved into
// it, or a condition name's first value, set.
type assignment struct {
	field string
	value string
}

// assignments lists the values a statement gives fields, as events after
// the statement's change: MOVE of a literal into fields by their plain
// names, and SET of condition names to true.
func (e *extractor) assignments(s cobol.Statement) []event {
	var events []event
	switch {
	case s.Verb == "MOVE" && len(s.Reads) == 1 && s.Reads[0].Name == "" && len(s.Reads[0].Refs) == 0:
		for _, w := range s.Writes {
			if w.Name != "" && w.Text == w.Name {
				events = append(events, event{assign: &assignment{field: w.Name, value: s.Reads[0].Text}})
			}
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

func (e *extractor) names(name string) (string, []string, bool) {
	if e.knowledge.ConditionNames == nil {
		return "", nil, false
	}
	return e.knowledge.ConditionNames(name)
}

// known is what a path knows of a field at its end: the value it was last
// given, by a statement or by a test found true, or the values tests
// found it not to be, since the last statement that may have changed it.
func (e *extractor) known(events []event, field string) (value string, ok bool, not []string) {
	for _, ev := range slices.Backward(events) {
		switch {
		case ev.assign != nil && ev.assign.field == field:
			return ev.assign.value, true, not
		case ev.change != nil && e.knowledge.MayChange(*ev.change, []string{field}):
			return "", false, not
		case ev.test != nil && ev.test.cond != nil:
			if v, isValue := e.valueOfTest(ev.test, field); isValue && ev.test.outcome {
				return v, true, not
			} else if isValue {
				not = append(not, v)
			}
		}
	}
	return "", false, not
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
		v, known, not := e.known(events, field)
		if known {
			return slices.ContainsFunc(values, func(w string) bool { return sameValue(v, w) }), true
		}
		excluded := !slices.ContainsFunc(values, func(w string) bool {
			return !slices.ContainsFunc(not, func(n string) bool { return sameValue(n, w) })
		})
		return false, excluded
	case "=", "<", ">", "<=", ">=":
		left, lok := e.operandValue(events, c.Left)
		right, rok := e.operandValue(events, c.Right)
		if !lok || !rok {
			return false, c.Op == "=" && (e.excludes(events, c.Left, c.Right) || e.excludes(events, c.Right, c.Left))
		}
		return compare(c.Op, left, right)
	case "TRUE":
		return true, true
	}
	return false, false
}

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

func (e *extractor) operandValue(events []event, op *cobol.Operand) (string, bool) {
	switch {
	case isLiteral(op):
		return op.Text, true
	case op != nil && op.Name != "" && op.Text == op.Name:
		v, ok, _ := e.known(events, op.Name)
		return v, ok
	}
	return "", false
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
