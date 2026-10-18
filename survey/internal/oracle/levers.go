package oracle

import (
	"slices"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
	"github.com/vlumi/copal/survey/blueprint"
)

// lever is one change to a case: a value put into a field of an input, in
// each record that has the field, or the value of another input field,
// give or take one, or a change to an input's records: none at all, one
// of each type more, or none of one type.
type lever struct {
	input  string
	field  *layout.Field
	value  value
	from   target
	delta  int64
	change recordChange
	kind   int
	// also is a second field change made with this one, for a condition
	// that tests two fields together.
	also *lever
}

type recordChange int

const (
	setField recordChange = iota
	copyField
	emptyInput
	doubled
	withoutType
)

// levers lists the changes worth trying, in a fixed order: for each
// condition of the program, the values it compares with and their
// boundaries, put into the input fields the condition depends on; then the
// changes to each input's records.
func (p *program) levers() []lever {
	var levers []lever
	add := func(l lever) {
		if !slices.Contains(levers, l) {
			levers = append(levers, l)
		}
	}
	for _, u := range p.units {
		for _, c := range conditionsOf(u.source) {
			for _, l := range p.fieldLevers(u, c) {
				add(l)
			}
			parts := conjuncts(c)
			for i := range parts {
				for j := i + 1; j < len(parts); j++ {
					for _, a := range p.fieldLevers(u, parts[i]) {
						for _, b := range p.fieldLevers(u, parts[j]) {
							if a.input == b.input && a.field != b.field {
								both := a
								both.also = &b
								add(both)
							}
						}
					}
				}
			}
		}
	}
	for _, c := range p.conditions() {
		for _, pair := range relations(c) {
			for _, to := range p.origins(pair[0]) {
				for _, from := range p.origins(pair[1]) {
					if to.field == from.field {
						continue
					}
					for _, delta := range []int64{0, 1, -1} {
						if delta == 0 || numeric(to.field) {
							add(lever{input: to.input, field: to.field, from: from, delta: delta, change: copyField})
						}
					}
				}
			}
		}
	}
	for _, ref := range p.references() {
		add(lever{input: ref[0].input, field: ref[0].field, from: ref[1], change: copyField})
	}
	for _, in := range p.inputs {
		add(lever{input: in.name(), change: emptyInput})
		add(lever{input: in.name(), change: doubled})
		for k := range p.shapes[in.name()].types {
			add(lever{input: in.name(), change: withoutType, kind: k})
		}
	}
	return levers
}

// fieldLevers lists the values worth putting into the input fields a
// condition of a unit depends on.
func (p *program) fieldLevers(u Unit, c *cobol.Condition) []lever {
	var levers []lever
	for _, item := range c.Refs() {
		for _, target := range p.origins(p.inCaller(u, item)) {
			for _, v := range valuesFor(u.records, c, item, target.field) {
				l := lever{input: target.input, field: target.field, value: v}
				if !slices.Contains(levers, l) {
					levers = append(levers, l)
				}
			}
		}
	}
	return levers
}

// conjuncts lists the parts of a condition that must hold together: those
// joined by AND, in each of the alternatives joined by OR.
func conjuncts(c *cobol.Condition) []*cobol.Condition {
	switch c.Op {
	case "AND":
		var parts []*cobol.Condition
		for _, part := range c.Parts {
			parts = append(parts, conjuncts(part)...)
		}
		return parts
	case "OR":
		var parts []*cobol.Condition
		for _, part := range c.Parts {
			if part.Op == "AND" {
				parts = append(parts, conjuncts(part)...)
			}
		}
		return parts
	}
	return []*cobol.Condition{c}
}

// conditions lists every condition the program tests.
func (p *program) conditions() []*cobol.Condition { return conditionsOf(p.source) }

// conditionsOf lists every condition a program tests, in order: of IF,
// of the WHENs of EVALUATE and SEARCH, and of PERFORM UNTIL.
func conditionsOf(source *cobol.Program) []*cobol.Condition {
	var conditions []*cobol.Condition
	for _, para := range source.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if s.Condition != nil {
				conditions = append(conditions, s.Condition)
			}
			for _, b := range s.Branches {
				if b.Condition != nil {
					conditions = append(conditions, b.Condition)
				}
			}
		})
	}
	return conditions
}

// references lists the input fields the program moves into the record key
// of an indexed input, to read it by: each with the key, whose values it
// should take to find a record.
func (p *program) references() [][2]target {
	var refs [][2]target
	for _, in := range p.inputs {
		if in.file.Organization != "INDEXED" || find(in.root, in.file.RecordKey) == nil {
			continue
		}
		key := target{input: in.name(), field: find(in.root, in.file.RecordKey)}
		for _, e := range p.flow.Lineage {
			if e.To != in.file.RecordKey || e.FromFile || e.FromLiteral {
				continue
			}
			for _, from := range p.origins(e.From) {
				if from.input != in.name() {
					for _, leaf := range leaves(key.field) {
						refs = append(refs, [2]target{from, {input: key.input, field: leaf}})
					}
				}
			}
		}
	}
	return refs
}

// target is a field of an input.
type target struct {
	input string
	field *layout.Field
}

// origins lists the input fields an item's values come from: the fields
// of the records it is read from, and of SYSIN. A condition name stands
// for its field.
func (p *program) origins(item string) []target {
	if f := conditionField(p.records, item); f != nil {
		item = f.Name
	}
	var targets []target
	for _, o := range p.flow.Origins(item) {
		name := o.File
		if o.Stream != "" {
			name = sysin
		}
		i := slices.IndexFunc(p.inputs, func(in input) bool { return in.name() == name })
		if i < 0 || name == "" {
			continue
		}
		if f := find(p.inputs[i].root, o.Item); f != nil {
			for _, leaf := range leaves(f) {
				targets = append(targets, target{input: name, field: leaf})
			}
		}
	}
	return targets
}

// leaves lists a field's elementary fields, or the field itself.
func leaves(f *layout.Field) []*layout.Field {
	if len(f.Children) == 0 {
		return []*layout.Field{f}
	}
	var list []*layout.Field
	for _, child := range f.Children {
		if child.Name != "FILLER" {
			list = append(list, leaves(child)...)
		}
	}
	return list
}

// valuesFor lists the values worth putting into an input field for a
// condition on an item: the literals the condition compares the item with,
// and for a number one above and below; the values of the condition names
// it tests, and a value of none of them; for a class test, a value of the
// class and one not; and for a number that reaches the item through other
// statements, zero, one, minus one and its largest.
func valuesFor(records *layout.ProgramRecords, c *cobol.Condition, item string, f *layout.Field) []value {
	var values []value
	addLiteral := func(lit string) {
		v, ok := literal(f, lit)
		if !ok {
			return
		}
		values = append(values, v)
		if v.kind == numberValue {
			values = append(values, number(v.number+1), number(v.number-1))
		}
	}
	var visit func(c *cobol.Condition)
	visit = func(c *cobol.Condition) {
		for _, part := range c.Parts {
			visit(part)
		}
		switch c.Op {
		case "=", "<", ">", "<=", ">=":
			for _, sides := range [][2]*cobol.Operand{{c.Left, c.Right}, {c.Right, c.Left}} {
				if sides[0] != nil && sides[0].Name == item && sides[1] != nil && sides[1].Name == "" && len(sides[1].Refs) == 0 {
					addLiteral(sides[1].Text)
				}
			}
		case "NAME":
			if field := conditionField(records, c.Left.Text); field != nil {
				for _, cond := range field.Conditions {
					for _, v := range cond.Values {
						addLiteral(v)
					}
				}
				values = append(values, text("?"))
			}
		case "IS":
			values = append(values, text("X"), number(0), number(1), number(-1))
		}
	}
	visit(c)
	if numeric(f) {
		values = append(values, number(0), number(1), number(-1), number(largest(f)))
	} else {
		values = append(values, text(""))
	}
	return values
}

// conditionField is the field a condition name is declared on.
func conditionField(records *layout.ProgramRecords, name string) *layout.Field {
	var found *layout.Field
	for _, r := range records.All() {
		walk(r, func(f *layout.Field) {
			if found == nil && slices.ContainsFunc(f.Conditions, func(c layout.Condition) bool { return c.Name == name }) {
				found = f
			}
		})
	}
	return found
}

// inCaller is the item of the program run that a unit's item is: itself
// in the program run; in a program it calls, the caller's field over the
// same bytes of the argument passed where the item's USING record is. A
// condition name stands for its field; an item the caller passes nothing
// for is none.
func (p *program) inCaller(u Unit, item string) string {
	if u.blueprint == p.blueprint {
		return item
	}
	if f := conditionField(u.records, item); f != nil {
		item = f.Name
	}
	for n, using := range u.blueprint.Using {
		record, field, ok := locate(u.blueprint, using, item)
		if !ok {
			continue
		}
		for _, call := range p.blueprint.Calls {
			if call.Program != u.Program || len(call.Using) <= n {
				continue
			}
			if arg, _, ok := locate(p.blueprint, call.Using[n], call.Using[n]); ok {
				if name, ok := over(arg, arg.Start+field.Start-record.Start, field.Length); ok {
					return name
				}
			}
		}
	}
	return ""
}

// locate finds a field within a program's record by name.
func locate(bp *blueprint.Program, record, name string) (blueprint.Field, blueprint.Field, bool) {
	for _, r := range bp.Records {
		if r.Name != record {
			continue
		}
		if f, ok := findField(r.Field, name); ok {
			return r.Field, f, true
		}
	}
	return blueprint.Field{}, blueprint.Field{}, false
}

func findField(f blueprint.Field, name string) (blueprint.Field, bool) {
	if f.Name == name {
		return f, true
	}
	for _, child := range f.Fields {
		if found, ok := findField(child, name); ok {
			return found, true
		}
	}
	return blueprint.Field{}, false
}

// over is the field of a record over exactly the bytes.
func over(f blueprint.Field, start, length int) (string, bool) {
	if f.Start == start && f.Length == length && f.Name != "FILLER" {
		return f.Name, true
	}
	for _, child := range f.Fields {
		if name, ok := over(child, start, length); ok {
			return name, true
		}
	}
	return "", false
}

// apply makes a candidate from a case with the lever's change, if it
// changes anything.
func (l lever) apply(p *program, c Case) (Case, bool) {
	next := c.clone()
	records := next[l.input]
	s := p.shapes[l.input]
	switch l.change {
	case emptyInput:
		if len(records) == 0 {
			return nil, false
		}
		next[l.input] = nil
	case doubled:
		if len(records) == 0 {
			return nil, false
		}
		next[l.input] = append(records, cloneRecords(records)...)
	case withoutType:
		kept := slices.DeleteFunc(records, func(r []byte) bool { return s.typeOfRecord(p.input(l.input), r) == l.kind })
		if len(kept) == len(c[l.input]) {
			return nil, false
		}
		next[l.input] = kept
	default:
		in := p.input(l.input)
		v := l.value
		if l.change == copyField {
			var ok bool
			if v, ok = p.valueOf(next, l.from, l.field, l.delta); !ok {
				return nil, false
			}
		}
		changed := false
		for i, r := range records {
			if !s.has(in, r, l.field) {
				continue
			}
			before := slices.Clone(r)
			if set(r, in.root, l.field, in.layout, v) && !slices.Equal(before, r) {
				changed = true
			}
			records[i] = r
		}
		if !changed {
			return nil, false
		}
		if l.also != nil {
			if both, ok := l.also.apply(p, next); ok {
				next = both
			}
		}
		p.balance(next, l.field)
	}
	return next, true
}

// valueOf reads an input field's value from the first record of a case
// that has it, give or take a delta, for another field.
func (p *program) valueOf(c Case, from target, to *layout.Field, delta int64) (value, bool) {
	in := p.input(from.input)
	s := p.shapes[from.input]
	for _, r := range c[from.input] {
		if !s.has(in, r, from.field) {
			continue
		}
		if numeric(from.field) && numeric(to) {
			n, err := in.layout.Number(r, from.field.Name)
			return number(n*pow10(to.Scale)/pow10(from.field.Scale) + delta), err == nil
		}
		t, err := in.layout.Text(r, from.field.Name)
		return text(t), err == nil
	}
	return value{}, false
}

func (p *program) input(name string) input {
	i := slices.IndexFunc(p.inputs, func(in input) bool { return in.name() == name })
	return p.inputs[i]
}

// typeOfRecord is the type a record is of, by its discriminator; -1 for a
// record of no known type.
func (s shape) typeOfRecord(in input, r []byte) int {
	if s.discriminator == nil {
		return -1
	}
	for k, t := range s.types {
		probe := in.layout.Blank()
		set(probe, in.root, s.discriminator, in.layout, t)
		start := s.discriminator.Start - in.root.Start
		if slices.Equal(probe[start:start+s.discriminator.Size()], r[start:start+s.discriminator.Size()]) {
			return k
		}
	}
	return -1
}

// has reports whether a record has a field: whether the field is outside
// the record's redefinitions, or in the one of the record's type, or in
// one of no type, when the record's type has none of its own.
func (s shape) has(in input, r []byte, f *layout.Field) bool {
	for variant, k := range s.variants {
		if find(variant, f.Name) == f {
			return s.hasVariant(k, s.typeOfRecord(in, r))
		}
	}
	return true
}

func (c Case) clone() Case {
	next := Case{}
	for name, records := range c {
		next[name] = cloneRecords(records)
	}
	return next
}

func cloneRecords(records [][]byte) [][]byte {
	out := make([][]byte, len(records))
	for i, r := range records {
		out[i] = slices.Clone(r)
	}
	return out
}
