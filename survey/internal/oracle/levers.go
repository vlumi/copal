package oracle

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/copal/survey/blueprint"
)

// lever is one change to a case: a value put into a field of an input, in
// each record that has the field, or the value of another input field,
// give or take one, or a change to an input's records: none at all, one
// of each type more, none of one type, those of one type of none, or as
// many as a count; or the input file failed.
type lever struct {
	input  string
	field  *layout.Field
	part   part
	value  value
	from   target
	delta  int64
	change recordChange
	kind   int
	count  int
	// failure is how the input fails, for a lever that fails it.
	failure failure
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
	untyped
	repeated
	failFile
)

// maxCount is the most records a lever repeats an input's to.
const maxCount = 100

// levers lists the changes worth trying, in a fixed order: for each
// condition of the program, the values it compares with and their
// boundaries, put into the input fields the condition depends on; then the
// changes to each input's records; then the failures of each input file.
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
					if to.field == from.field || to.part != (part{}) || from.part != (part{}) {
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
			if !numeric(p.shapes[in.name()].discriminator) {
				add(lever{input: in.name(), change: untyped, kind: k})
			}
		}
		for _, n := range p.counts() {
			add(lever{input: in.name(), change: repeated, count: n})
		}
	}
	for _, l := range p.failures() {
		add(l)
	}
	return levers
}

// counts are the numbers of records worth trying an input with: one more
// than each number, up to maxCount, that a condition compares an item the
// program only adds to with, such as a count of the records read.
func (p *program) counts() []int {
	var counts []int
	for _, u := range p.units {
		for _, c := range conditionsOf(u.source) {
			for _, cmp := range comparisons(c) {
				n, err := strconv.Atoi(cmp.Right.Text)
				if err == nil && n > 0 && n < maxCount && u.onlyAddedTo(cmp.Left.Name) && !slices.Contains(counts, n+1) {
					counts = append(counts, n+1)
				}
			}
		}
	}
	return counts
}

// comparisons lists the comparisons within a condition of an item with a
// number, the item on the left.
func comparisons(c *cobol.Condition) []*cobol.Condition {
	var list []*cobol.Condition
	for _, part := range c.Parts {
		list = append(list, comparisons(part)...)
	}
	switch c.Op {
	case "=", ">", ">=":
		if c.Left.Name != "" && c.Right.Name == "" && len(c.Right.Refs) == 0 {
			list = append(list, c)
		}
	}
	return list
}

// fieldLevers lists the values worth putting into the input fields a
// condition of a unit depends on.
func (p *program) fieldLevers(u Unit, c *cobol.Condition) []lever {
	var levers []lever
	for _, item := range c.Refs() {
		for _, target := range p.callerOrigins(u, item) {
			for _, v := range valuesFor(u.records, constantsOf(u.blueprint), c, item, target.part.of(target.field)) {
				l := lever{input: target.input, field: target.field, part: target.part, value: v}
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

// references are the links the blueprint finds the program reads an
// indexed file by: each input field it moves into the file's record key,
// with the key, whose values it should take to find a record.
func (p *program) references() [][2]target {
	var refs [][2]target
	for _, l := range p.blueprint.Links {
		if l.Kind != "key" {
			continue
		}
		from, fok := p.inputField(l.From)
		to, tok := p.inputField(l.To)
		if fok && tok {
			refs = append(refs, [2]target{from, to})
		}
	}
	return refs
}

// inputField finds a field of a file, or of SYSIN, among the program's
// inputs.
func (p *program) inputField(f blueprint.FileField) (target, bool) {
	for _, in := range p.inputs {
		if in.name() == f.File || in.file.Name == f.File {
			if field := find(in.root, f.Field); field != nil {
				return target{input: in.name(), field: field}, true
			}
		}
	}
	return target{}, false
}

// target is a field of an input, or a part of one.
type target struct {
	input string
	field *layout.Field
	part  part
}

// part is digits of a field of digits, from offset, 0-based, for length,
// which a lever sets alone: the month of a date, say. The zero part is
// the whole field.
type part struct{ offset, length int }

// of is a field as its part: a number of the part's digits.
func (pt part) of(f *layout.Field) *layout.Field {
	if pt.length == 0 {
		return f
	}
	return &layout.Field{Name: fmt.Sprintf("%s(%d:%d)", f.Name, pt.offset+1, pt.length), Start: f.Start + pt.offset,
		Length: pt.length, Category: layout.Numeric, Digits: pt.length, Picture: fmt.Sprintf("9(%d)", pt.length)}
}

// set writes a number into a part of a field of a record, if it fits.
func (pt part) set(r []byte, root, f *layout.Field, v value) bool {
	if v.kind != numberValue || v.number < 0 || v.number > largest(pt.of(f)) {
		return false
	}
	start := f.Start - root.Start + pt.offset
	copy(r[start:start+pt.length], fmt.Sprintf("%0*d", pt.length, v.number))
	return true
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
			name = o.Stream
		}
		i := slices.IndexFunc(p.inputs, func(in input) bool { return in.name() == name })
		if i < 0 || name == "" {
			continue
		}
		f := find(p.inputs[i].root, o.Item)
		if over := fieldOverBytes(p.inputs[i].root, o.Start, o.Length); o.Length > 0 && over != nil {
			f = over
		} else if digits := digitsOver(p.inputs[i].root, o.Start, o.Length); o.Length > 0 && digits != nil {
			targets = append(targets, target{input: name, field: digits, part: part{offset: o.Start - digits.Start, length: o.Length}})
			continue
		}
		if f != nil {
			for _, leaf := range leaves(f) {
				targets = append(targets, target{input: name, field: leaf})
			}
		}
	}
	return targets
}

// fieldOverBytes is the field of a record over the bytes from start,
// 1-based, for length, if one is over them exactly.
func fieldOverBytes(f *layout.Field, start, length int) *layout.Field {
	if f.Start == start && f.Size() == length && f.Name != "FILLER" {
		return f
	}
	for _, child := range f.Children {
		if found := fieldOverBytes(child, start, length); found != nil {
			return found
		}
	}
	return nil
}

// digitsOver is the elementary field of unsigned digits whose bytes the
// bytes from start, 1-based, for length, are within, if one is.
func digitsOver(f *layout.Field, start, length int) *layout.Field {
	if len(f.Children) == 0 {
		within := f.Start <= start && start+length <= f.Start+f.Size()
		if within && (numeric(f) && !f.Signed && f.Size() == f.Digits || f.Category == layout.Alphanumeric) {
			return f
		}
		return nil
	}
	for _, child := range f.Children {
		if found := digitsOver(child, start, length); found != nil {
			return found
		}
	}
	return nil
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
// condition on an item: the literals and level-78 constants the condition
// compares the item with, and what it divides the item by, and for a
// number one above and below; the values of the condition names
// it tests, and a value of none of them; for a class test, a value of the
// class and one not; and for a number that reaches the item through other
// statements, zero, one, minus one and its largest, and with decimals, one
// and minus one whole, which survive rounding to fewer.
func valuesFor(records *layout.ProgramRecords, constants map[string]string, c *cobol.Condition, item string, f *layout.Field) []value {
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
				if d, ok := divisor(records, sides[0], item); ok {
					addLiteral(d)
				}
				if sides[0] == nil || sides[0].Name != item || sides[1] == nil {
					continue
				}
				if value, ok := constants[sides[1].Name]; ok {
					addLiteral(value)
				} else if sides[1].Name == "" && len(sides[1].Refs) == 0 {
					addLiteral(sides[1].Text)
				}
			}
		case "NAME":
			if field := conditionField(records, c.Left.Name); field != nil && (c.Left.Name == item || field.Name == item) {
				for _, cond := range field.Conditions {
					for _, v := range cond.Values {
						addLiteral(v)
					}
				}
				values = append(values, text("?"))
			}
		case "IS":
			values = append(values, text("X"), number(0), number(1), number(-1))
			if c.Class == "NUMERIC" && !numeric(f) {
				values = append(values, text(strings.Repeat("1", f.Size())))
			}
		}
	}
	visit(c)
	if numeric(f) {
		values = append(values, number(0), number(1), number(-1), number(largest(f)))
		if f.Scale > 0 {
			values = append(values, number(pow10(f.Scale)), number(-pow10(f.Scale)))
		}
	} else {
		values = append(values, text(""))
	}
	return values
}

// divisor is what an operand divides an item by, as written, if it is a
// FUNCTION MOD or REM of the item: a literal, or the VALUE of the field it
// names.
func divisor(records *layout.ProgramRecords, op *cobol.Operand, item string) (string, bool) {
	if op == nil {
		return "", false
	}
	args, ok := strings.CutPrefix(op.Text, "FUNCTION MOD(")
	if !ok {
		args, ok = strings.CutPrefix(op.Text, "FUNCTION REM(")
	}
	words := strings.Fields(strings.TrimSuffix(args, ")"))
	if !ok || len(words) != 2 || words[0] != item {
		return "", false
	}
	var by *layout.Field
	for _, r := range records.All() {
		walk(r, func(f *layout.Field) {
			if by == nil && f.Name == words[1] {
				by = f
			}
		})
	}
	if by == nil {
		return words[1], true
	}
	if len(by.Values) == 1 {
		return by.Values[0], true
	}
	return "", false
}

// constantsOf are a program's level-78 constants, by name, with their
// values as written.
func constantsOf(bp *blueprint.Program) map[string]string {
	constants := map[string]string{}
	for _, c := range bp.Constants {
		constants[c.Name] = c.Value
	}
	return constants
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

// callerOrigins lists the input fields a unit's item's values come from:
// those of the item of the program run it is, or, for an item of a program
// called that is no item passed to it, those of the items passed that its
// values come from.
func (p *program) callerOrigins(u Unit, item string) []target {
	if name := p.inCaller(u, item); name != "" || u.blueprint == p.blueprint {
		return p.origins(name)
	}
	var targets []target
	for _, o := range u.flow.Origins(item) {
		if o.Parameter == 0 {
			continue
		}
		for _, t := range p.origins(p.inCaller(u, o.Item)) {
			if !slices.Contains(targets, t) {
				targets = append(targets, t)
			}
		}
	}
	return targets
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
			if arg, ok := fieldOf(p.blueprint, call.Using[n]); ok {
				if name, ok := over(arg, arg.Start+field.Start-record.Start, field.Length); ok {
					return name
				}
			}
		}
	}
	return ""
}

// fieldOf finds a field of any of a program's records by name: an argument
// passed, which may be a record or a field of one.
func fieldOf(bp *blueprint.Program, name string) (blueprint.Field, bool) {
	for _, r := range bp.Records {
		if f, ok := findField(r.Field, name); ok {
			return f, true
		}
	}
	return blueprint.Field{}, false
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
	if l.change == failFile {
		return c.fail(l.input, l.failure)
	}
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
	case repeated:
		if len(records) == 0 || len(records) >= l.count {
			return nil, false
		}
		for i := len(records); i < l.count; i++ {
			records = append(records, slices.Clone(records[i%len(c[l.input])]))
		}
		next[l.input] = records
	case untyped:
		in := p.input(l.input)
		changed := false
		for i, r := range records {
			if s.typeOfRecord(in, r) == l.kind {
				set(r, in.root, s.discriminator, in.layout, text("?"))
				records[i], changed = r, true
			}
		}
		if !changed {
			return nil, false
		}
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
			written := false
			if l.part != (part{}) {
				written = l.part.set(r, in.root, l.field, v)
			} else {
				written = set(r, in.root, l.field, in.layout, v)
			}
			if written && !slices.Equal(before, r) {
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

// has reports whether a record has a field: whether the field is the type
// field, which every record has wherever it is declared, or is outside the
// record's redefinitions, or in the one of the record's type, or in one of
// no type, when the record's type has none of its own.
func (s shape) has(in input, r []byte, f *layout.Field) bool {
	if f == s.discriminator {
		return true
	}
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
