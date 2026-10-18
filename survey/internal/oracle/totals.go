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

// totals finds the control totals of the program's inputs: each condition
// that compares a field of an input with an item that is no input's, into
// which the program adds a field of the same input, or 1. A count of
// records of no known type counts those a sum of the same input adds up.
func (p *program) totals() []total {
	var totals []total
	for _, c := range p.conditions() {
		for _, pair := range relations(c) {
			if t, ok := p.totalOf(pair[0], pair[1]); ok && !slices.ContainsFunc(totals, func(u total) bool { return u.field == t.field }) {
				totals = append(totals, t)
			}
		}
	}
	for i, t := range totals {
		if t.sumOf != nil || t.counts >= 0 {
			continue
		}
		if j := slices.IndexFunc(totals, func(u total) bool { return u.input == t.input && u.sumOf != nil && u.counts >= 0 }); j >= 0 {
			totals[i].counts = totals[j].counts
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

// totalOf reads a comparison as a control total: a field of an input on
// one side, and on the other an item the program adds a field of the same
// input to, or adds 1 to.
func (p *program) totalOf(field, item string) (total, bool) {
	var in *input
	var f *layout.Field
	for i := range p.inputs {
		if found := find(p.inputs[i].root, field); found != nil {
			in, f = &p.inputs[i], found
		}
	}
	if in == nil || p.isInputField(item) {
		return total{}, false
	}
	s := p.shapes[in.name()]
	for _, e := range p.flow.Lineage {
		if e.To != item || e.Verb != "ADD" {
			continue
		}
		switch {
		case e.FromLiteral && e.From == "1":
			return total{input: in.name(), field: f, counts: p.typeUsing(map[string]bool{item: true}, s)}, true
		case !e.FromLiteral:
			if summed := find(in.root, e.From); summed != nil && summed != f {
				return total{input: in.name(), field: f, sumOf: summed, counts: s.typeOfField(summed)}, true
			}
		}
	}
	return total{}, false
}

func (p *program) isInputField(name string) bool {
	return slices.ContainsFunc(p.inputs, func(in input) bool { return find(in.root, name) != nil })
}

// typeOfField is the record type whose redefinition a field is in; -1 for
// a field outside the redefinitions.
func (s shape) typeOfField(f *layout.Field) int {
	for variant, kind := range s.variants {
		if find(variant, f.Name) == f {
			return kind
		}
	}
	return -1
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
