package oracle

import (
	"bytes"
	"slices"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/copal/survey/blueprint"
)

// total is a field of an input that the program checks against an item it
// accumulates from the input's other records: the sum of a field of
// theirs, or how many of them there are.
type total struct {
	input  string
	field  *layout.Field
	sumOf  *layout.Field // nil for a count
	counts int           // the type of the records added up; -1 for all
	guards []blueprint.Guard
}

// totals are the control totals the blueprint finds the program checks,
// with the fields they are of.
func (p *program) totals() []total {
	var totals []total
	for _, t := range p.blueprint.Totals {
		i := slices.IndexFunc(p.inputs, func(in input) bool { return in.file.Name == t.File })
		if i < 0 {
			continue
		}
		in := p.inputs[i]
		u := total{input: in.name(), field: find(in.root, t.Field), counts: slices.Index(p.shapes[in.name()].typeNames, t.Type), guards: t.Guards}
		if t.Sum != "" {
			u.sumOf = find(in.root, t.Sum)
		}
		if u.field != nil && (t.Sum == "" || u.sumOf != nil) {
			totals = append(totals, u)
		}
	}
	return totals
}

// relations lists the pairs of data items a condition compares, each pair
// both ways round.
func relations(c *cobol.Condition) [][2]string {
	var pairs [][2]string
	for _, part := range c.Parts {
		pairs = append(pairs, relations(part)...)
	}
	switch c.Op {
	case "=", "<", ">", "<=", ">=":
		if c.Left.Name != "" && c.Right.Name != "" {
			pairs = append(pairs, [2]string{c.Left.Name, c.Right.Name}, [2]string{c.Right.Name, c.Left.Name})
		}
	}
	return pairs
}

// balance sets each control total of a case to what the records it adds
// up come to, except the field a lever has just changed, so that a total
// that differs can be had too.
func (p *program) balance(c Case, except *layout.Field) {
	for _, t := range p.allTotals {
		if t.field == except {
			continue
		}
		in := p.input(t.input)
		s := p.shapes[t.input]
		var sum int64
		for _, r := range c[t.input] {
			if t.counts >= 0 && s.typeOfRecord(in, r) != t.counts || !p.addedUp(in, r, t.guards) {
				continue
			}
			if t.sumOf == nil {
				sum++
			} else if n, err := in.layout.Number(r, t.sumOf.Name); err == nil {
				sum += n * pow10(t.field.Scale-t.sumOf.Scale)
			}
		}
		for i, r := range c[t.input] {
			if s.has(in, r, t.field) && (t.counts < 0 || s.typeOfRecord(in, r) != t.counts) {
				set(r, in.root, t.field, in.layout, number(sum))
				c[t.input][i] = r
			}
		}
	}
}

// addedUp reports whether the program adds a record up into a total: each
// condition the total is guarded by has its outcome on the record, or one
// the record does not tell.
func (p *program) addedUp(in input, r []byte, guards []blueprint.Guard) bool {
	for _, g := range guards {
		if outcome, known := p.holdsOn(in, r, g.Condition); known && outcome != g.Outcome {
			return false
		}
	}
	return true
}

// holdsOn evaluates a condition on a record's fields, and whether the
// record tells: a condition name, a class test, and a comparison for
// equality with a literal, of the record's fields, and the conditions
// joined of those.
func (p *program) holdsOn(in input, r []byte, c *cobol.Condition) (outcome, known bool) {
	switch c.Op {
	case "NOT":
		outcome, known = p.holdsOn(in, r, c.Parts[0])
		return !outcome, known
	case "AND", "OR":
		all := true
		for _, part := range c.Parts {
			o, k := p.holdsOn(in, r, part)
			switch {
			case k && o == (c.Op == "OR"):
				return o, true
			case !k:
				all = false
			}
		}
		return c.Op == "AND", all
	case "NAME":
		f := conditionField(p.records, c.Left.Name)
		if f == nil || find(in.root, f.Name) == nil {
			return false, false
		}
		i := slices.IndexFunc(f.Conditions, func(cond layout.Condition) bool { return cond.Name == c.Left.Name })
		for _, v := range f.Conditions[i].Values {
			is, ok := p.fieldIs(in, r, f, v)
			if !ok {
				return false, false
			}
			if is {
				return true, true
			}
		}
		return false, true
	case "IS":
		f := find(in.root, c.Left.Name)
		if f == nil || c.Class != "NUMERIC" {
			return false, false
		}
		_, err := in.layout.Number(r, f.Name)
		return err == nil, true
	case "=":
		for _, side := range [][2]*cobol.Operand{{c.Left, c.Right}, {c.Right, c.Left}} {
			if f := find(in.root, side[0].Name); f != nil && side[1].Name == "" {
				return p.fieldIs(in, r, f, side[1].Text)
			}
		}
	}
	return false, false
}

// fieldIs reports whether a field of a record holds a literal, and
// whether that can be told: a range of values cannot.
func (p *program) fieldIs(in input, r []byte, f *layout.Field, lit string) (is, ok bool) {
	if strings.Contains(lit, " THRU ") {
		return false, false
	}
	v, ok := literal(f, lit)
	if !ok {
		return false, false
	}
	want := slices.Clone(r)
	if !set(want, in.root, f, in.layout, v) {
		return false, true
	}
	start := f.Start - in.root.Start
	return bytes.Equal(r[start:start+f.Size()], want[start:start+f.Size()]), true
}
