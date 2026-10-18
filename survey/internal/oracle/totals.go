package oracle

import (
	"slices"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
)

// total is a field of an input that the program checks against an item it
// accumulates from the input's other records: the sum of a field of
// theirs, or how many of them there are.
type total struct {
	input  string
	field  *layout.Field
	sumOf  *layout.Field // nil for a count
	counts int           // the type of the records added up; -1 for all
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
		u := total{input: in.name(), field: find(in.root, t.Field), counts: slices.Index(p.shapes[in.name()].typeNames, t.Type)}
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
			if t.counts >= 0 && s.typeOfRecord(in, r) != t.counts {
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
