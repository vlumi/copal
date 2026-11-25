package oracle

import (
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
)

// A condition on a file's status holds, does not hold, or is not settled
// by the status: a condition name of it by its values, NOT, and AND and
// OR as far as their parts tell, a part settled otherwise among them.
func TestHolds(t *testing.T) {
	status := &layout.Field{Name: "WS-ST", Length: 2, Conditions: []layout.Condition{
		{Name: "WS-OK", Values: []string{`"00"`}}, {Name: "WS-END", Values: []string{`"10"`}}}}
	u := Unit{records: &layout.ProgramRecords{WorkingStorage: []*layout.Field{status}}}
	name := func(n string) *cobol.Condition {
		return &cobol.Condition{Op: "NAME", Left: &cobol.Operand{Text: n, Name: n}}
	}
	not := func(c *cobol.Condition) *cobol.Condition {
		return &cobol.Condition{Op: "NOT", Parts: []*cobol.Condition{c}}
	}
	join := func(op string, parts ...*cobol.Condition) *cobol.Condition {
		return &cobol.Condition{Op: op, Parts: parts}
	}
	failed := join("AND", not(name("WS-OK")), not(name("WS-END")))
	other := &cobol.Condition{Op: "=", Left: &cobol.Operand{Text: "WS-COUNT", Name: "WS-COUNT"}, Right: &cobol.Operand{Text: "0"}}
	for _, c := range []struct {
		condition *cobol.Condition
		status    string
		want      held
	}{
		{name("WS-OK"), "00", yes},
		{name("WS-OK"), "10", no},
		{failed, "00", no},
		{failed, "10", no},
		{failed, "30", yes},
		{other, "00", unknown},
		{join("AND", name("WS-OK"), other), "00", unknown},
		{join("AND", name("WS-OK"), other), "10", no},
		{join("OR", name("WS-OK"), other), "00", yes},
		{join("OR", name("WS-OK"), other), "10", unknown},
	} {
		if got := u.holds(c.condition, status, c.status, nil); got != c.want {
			t.Errorf("%s with status %s: %v, want %v", c.condition, c.status, got, c.want)
		}
	}
	settled := func(c *cobol.Condition) held {
		if c == other {
			return yes
		}
		return unknown
	}
	if got := u.holds(join("AND", name("WS-OK"), other), status, "00", settled); got != yes {
		t.Errorf("WS-OK and a condition settled true, with status 00: %v, want %v", got, yes)
	}
}
