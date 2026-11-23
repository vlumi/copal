package oracle

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/dataflow"
)

// lever is one change to a case: a value put into a field of an input, in
// each record that has the field, or the value of another input field,
// give or take one, or a change to an input's records: none at all, one
// of each type more, none of one type, those of one type of none, or as
// many as a count, or one more under each key of several parts; or the
// input file failed.
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
	continued
	failFile
)

// maxCount is the most records a lever repeats an input's to.
const maxCount = 100

// levers lists the changes worth trying, in a fixed order: for each
// condition of the program, the values it compares with and their
// boundaries, put into the input fields the condition depends on; then the
// changes to each input's records; then the failures of each input file.
// A field that holds a date as text takes the dates around the business
// date, in the form the system writes them.
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
	for _, u := range p.units {
		for _, c := range conditionsOf(u.source) {
			for _, pair := range relations(c) {
				for _, to := range p.callerOrigins(u, pair[0]) {
					for _, from := range p.callerOrigins(u, pair[1]) {
						digits := numeric(to.field) && from.part.length == to.field.Digits && to.field.Scale == 0
						if to.field == from.field || to.part != (part{}) || from.part != (part{}) && !digits {
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
	}
	for _, ref := range p.references() {
		add(lever{input: ref[0].input, field: ref[0].field, from: ref[1], change: copyField})
	}
	for _, l := range p.dateLevers() {
		add(l)
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
		if _, ok := in.lastOfKey(); ok {
			add(lever{input: in.name(), change: continued})
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
// condition of a unit depends on. A part of a field of text is tried as
// digits, which a date in it is, and as text, which a prefix a range of
// numbers is told by is.
func (p *program) fieldLevers(u Unit, c *cobol.Condition) []lever {
	var levers []lever
	for _, item := range c.Refs() {
		for _, target := range p.callerOrigins(u, item) {
			as := []*layout.Field{target.part.of(target.field)}
			if text, ok := target.part.textOf(target.field); ok {
				as = append(as, text)
			}
			compared := c
			if target.madeFrom {
				compared = &cobol.Condition{}
			}
			for _, f := range as {
				for _, v := range valuesFor(u.records, constantsOf(u.blueprint), p.presetsOf(u), compared, item, f) {
					l := lever{input: target.input, field: target.field, part: target.part, value: v}
					if !slices.Contains(levers, l) {
						levers = append(levers, l)
					}
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

// target is a field of an input, or a part of one. It is made from when
// the item it is a target for does not hold its value but one a program
// called made from it, a cash value from a sum insured: what a condition
// compares the item with is then nothing to put into the field.
type target struct {
	input    string
	field    *layout.Field
	part     part
	madeFrom bool
}

// part is bytes of a field, from offset, 0-based, for length, which a
// lever sets alone: digits of a field of digits, the month of a date, or
// characters of a field of text, the prefix of a number a range is told
// by. The zero part is the whole field.
type part struct{ offset, length int }

// of is a field as its part: a number of the part's digits.
func (pt part) of(f *layout.Field) *layout.Field {
	if pt.length == 0 {
		return f
	}
	return &layout.Field{Name: pt.name(f), Start: f.Start + pt.offset, Length: pt.length, Category: layout.Numeric, Digits: pt.length,
		Picture: fmt.Sprintf("9(%d)", pt.length)}
}

// textOf is a part of a field of text as text, its characters, or nothing
// for a part of a number, or the whole field.
func (pt part) textOf(f *layout.Field) (*layout.Field, bool) {
	if pt.length == 0 || numeric(f) {
		return nil, false
	}
	return &layout.Field{Name: pt.name(f), Start: f.Start + pt.offset, Length: pt.length, Category: f.Category,
		Picture: fmt.Sprintf("X(%d)", pt.length)}, true
}

func (pt part) name(f *layout.Field) string {
	return fmt.Sprintf("%s(%d:%d)", f.Name, pt.offset+1, pt.length)
}

// set writes a value into a part of a field of a record, if it fits: a
// number as digits, or into a field of text, text, padded with spaces.
func (pt part) set(r []byte, root, f *layout.Field, v value) bool {
	start := f.Start - root.Start + pt.offset
	switch {
	case v.kind == textValue && !numeric(f) && len(v.text) <= pt.length:
		copy(r[start:start+pt.length], fmt.Sprintf("%-*s", pt.length, v.text))
		return true
	case v.kind != numberValue || v.number < 0 || v.number > largest(pt.of(f)):
		return false
	}
	copy(r[start:start+pt.length], fmt.Sprintf("%0*d", pt.length, v.number))
	return true
}

// origins lists the input fields an item of the program run has its
// values from: the fields of the records it is read from, and of SYSIN.
// A condition name stands for its field.
func (p *program) origins(item string) []target { return p.originsIn(p.units[0], item) }

// direct lists the input fields an item of the program run has its values
// from without a program called between them: what the base case matches
// to each other, which a field a called program only may have made a
// value from is too little to go by.
func (p *program) direct(item string) []target { return p.originsBehind(p.units[0], item, nil) }

// originsIn lists the input fields an item of a unit has its values from
// within the unit: the fields of the records it reads them from, and for
// a value a program it calls gives back, what that program made it from.
func (p *program) originsIn(u Unit, item string) []target {
	return p.originsBehind(u, item, map[string]bool{})
}

// originsBehind is originsIn, with the calls already followed, or with
// no call followed at all when there is nowhere to note them.
func (p *program) originsBehind(u Unit, item string, followed map[string]bool) []target {
	if f := conditionField(u.records, item); f != nil {
		item = f.Name
	}
	var targets []target
	add := func(found []target) {
		for _, t := range found {
			if !slices.Contains(targets, t) {
				targets = append(targets, t)
			}
		}
	}
	for _, o := range u.flow.Origins(item) {
		if o.Call != "" && o.Arg > 0 {
			if followed != nil {
				for _, t := range p.behindCall(u, item, o, followed) {
					t.madeFrom = true
					add([]target{t})
				}
			}
			continue
		}
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
			add([]target{{input: name, field: digits, part: part{offset: o.Start - digits.Start, length: o.Length}}})
			continue
		}
		if f != nil {
			for _, leaf := range leaves(f) {
				add([]target{{input: name, field: leaf}})
			}
		}
	}
	return targets
}

// behindCall lists the input fields a program called makes a value from
// that it gives back to an item of the unit that calls it: the fields of
// the files it reads itself and, for the program run, those the other
// arguments have their values from at each call. The item is followed
// when it is the argument itself or a field of it, into the called
// program's field over the same bytes; each once.
func (p *program) behindCall(u Unit, item string, o dataflow.Origin, followed map[string]bool) []target {
	i := slices.IndexFunc(p.units, func(c Unit) bool { return c.Program == o.Call })
	if i < 0 || len(p.units[i].blueprint.Using) < o.Arg {
		return nil
	}
	callee := p.units[i]
	inner, ok := calleeItem(u, callee, item, o)
	key := fmt.Sprintf("%s %s %s", u.Program, o.Call, inner)
	if !ok || followed[key] {
		return nil
	}
	followed[key] = true
	targets := p.originsBehind(callee, inner, followed)
	if u.blueprint != p.blueprint {
		return targets
	}
	for _, q := range callee.flow.Origins(inner) {
		if q.Parameter == 0 || q.Item == inner {
			continue
		}
		for _, site := range p.inCaller(callee, q.Item) {
			for _, t := range within(p.originsBehind(u, site.name, followed), site.part) {
				if !slices.Contains(targets, t) {
					targets = append(targets, t)
				}
			}
		}
	}
	return targets
}

// calleeItem is the item of a program called that an item of its caller
// has its value from, by the argument a call passes: the parameter, for
// an argument of one field, whatever the item makes of it, or for a
// record passed, the field of the parameter over the same bytes as the
// item is of the argument. An item that is no field of a record passed
// has its value from some field of it the call does not tell, and is not
// followed.
func calleeItem(u, callee Unit, item string, o dataflow.Origin) (string, bool) {
	parameter := callee.blueprint.Using[o.Arg-1]
	argument, ok := fieldOf(u.blueprint, o.Item)
	if !ok {
		return "", false
	}
	if len(argument.Fields) == 0 {
		return parameter, true
	}
	field, ok := findField(argument, item)
	record, _, found := locate(callee.blueprint, parameter, parameter)
	if !ok || !found {
		return "", false
	}
	return over(record, record.Start+field.Start-argument.Start, field.Length)
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
// and minus one whole, which survive rounding to fewer. An item compared
// with a field of a table written with its values takes each of them.
func valuesFor(records *layout.ProgramRecords, constants map[string]string, presets map[string][]string, c *cobol.Condition, item string, f *layout.Field) []value {
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
				for _, preset := range presets[sides[1].Name] {
					if v, ok := literal(f, strconv.Quote(preset)); ok && !slices.Contains(values, v) {
						values = append(values, v)
					}
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

// presetsOf are the values a unit's tables are written with, found once.
func (p *program) presetsOf(u Unit) map[string][]string {
	if p.presets == nil {
		p.presets = map[string]map[string][]string{}
	}
	if _, ok := p.presets[u.Program]; !ok {
		p.presets[u.Program] = presets(u.records, p.dialect.Encoding)
	}
	return p.presets[u.Program]
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
// called that is no item passed to it, those of the files the program
// called reads itself and of the items passed that its values come from.
func (p *program) callerOrigins(u Unit, item string) []target {
	if u.blueprint == p.blueprint {
		return p.origins(item)
	}
	var targets []target
	add := func(sites []callerItem) {
		for _, site := range sites {
			for _, t := range within(p.origins(site.name), site.part) {
				if !slices.Contains(targets, t) {
					targets = append(targets, t)
				}
			}
		}
	}
	if sites := p.inCaller(u, item); len(sites) > 0 {
		add(sites)
		return targets
	}
	targets = p.originsIn(u, item)
	for _, o := range u.flow.Origins(item) {
		if o.Parameter == 0 {
			continue
		}
		sites := p.inCaller(u, o.Item)
		if moved, ok := u.movedFrom(item, o.Item); ok {
			for i := range sites {
				sites[i].part = moved.inside(sites[i].part)
			}
		}
		add(sites)
	}
	return targets
}

// movedFrom is the part of an item that another has its value from, when
// a statement moves the one into bytes the other is some of: the day of
// a date, from the last two digits of a month and day moved into the
// date's last four.
func (u Unit) movedFrom(item, from string) (part, bool) {
	held, source := u.declared(item), u.declared(from)
	if held == nil || source == nil {
		return part{}, false
	}
	for _, e := range u.flow.Lineage {
		into := u.declared(e.To)
		if e.From != from || e.FromFile || e.FromLiteral || into == nil {
			continue
		}
		start, length := into.Start+max(e.ToStart, 1)-1, e.ToLength
		if e.ToStart == 0 {
			length = into.Size()
		}
		at := held.Start - start + max(e.FromStart, 1) - 1
		if held.Start >= start && held.Start+held.Size() <= start+length && at+held.Size() <= source.Size() && held.Size() < source.Size() {
			return part{offset: at, length: held.Size()}, true
		}
	}
	return part{}, false
}

// declared finds an item of the unit by name.
func (u Unit) declared(name string) *layout.Field {
	for _, r := range u.records.All() {
		if f := find(r, name); f != nil {
			return f
		}
	}
	return nil
}

// inside is a part as a part of another: itself, of a whole field.
func (pt part) inside(outer part) part {
	if outer == (part{}) {
		return pt
	}
	return part{offset: outer.offset + pt.offset, length: pt.length}
}

// within narrows targets to a part of their fields: the bytes of an
// argument a called program's item is over, when no field of the caller
// is over exactly those. A target that is a part already is narrowed to
// the part of that part, and one the part reaches beyond is left out.
func within(targets []target, pt part) []target {
	if pt == (part{}) {
		return targets
	}
	var parts []target
	for _, t := range targets {
		room := t.field.Length
		if t.part != (part{}) {
			room = t.part.length
		}
		if pt.offset+pt.length <= room {
			t.part = pt.inside(t.part)
			parts = append(parts, t)
		}
	}
	return parts
}

// callerItem is an item of the program run, or a part of one, that an
// item of a program it calls is at one of its calls.
type callerItem struct {
	name string
	part part
}

// inCaller lists the items of the program run that an item of a program
// it calls is, one for each call that passes something for it: the
// caller's field over the same bytes of the argument passed where the
// item's USING record is, or, when no field is over exactly those bytes,
// the field they are within and the part of it they are. A condition
// name stands for its field.
func (p *program) inCaller(u Unit, item string) []callerItem {
	if f := conditionField(u.records, item); f != nil {
		item = f.Name
	}
	var items []callerItem
	for n, using := range u.blueprint.Using {
		record, field, ok := locate(u.blueprint, using, item)
		if !ok {
			continue
		}
		for _, call := range p.blueprint.Calls {
			if call.Program != u.Program || len(call.Using) <= n {
				continue
			}
			arg, ok := fieldOf(p.blueprint, call.Using[n])
			if !ok {
				continue
			}
			start := arg.Start + field.Start - record.Start
			var found callerItem
			if name, ok := over(arg, start, field.Length); ok {
				found = callerItem{name: name}
			} else if leaf, ok := around(arg, start, field.Length); ok {
				found = callerItem{name: leaf.Name, part: part{offset: start - leaf.Start, length: field.Length}}
			}
			if found.name != "" && !slices.Contains(items, found) {
				items = append(items, found)
			}
		}
	}
	return items
}

// around is the elementary field of a record that bytes are within.
func around(f blueprint.Field, start, length int) (blueprint.Field, bool) {
	for _, child := range f.Fields {
		if leaf, ok := around(child, start, length); ok {
			return leaf, true
		}
	}
	if len(f.Fields) == 0 && f.Name != "FILLER" && f.Start <= start && start+length <= f.Start+f.Length {
		return f, true
	}
	return blueprint.Field{}, false
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
	case continued:
		in := p.input(l.input)
		last, _ := in.lastOfKey()
		key := find(in.root, in.file.RecordKey)
		keyOf := func(r []byte) string { return string(r[key.Start-in.root.Start:][:key.Size()]) }
		for _, r := range c[l.input] {
			more := slices.Clone(r)
			n, err := in.layout.Number(more, last.Name)
			if err != nil || !set(more, in.root, last, in.layout, number(n+1)) ||
				slices.ContainsFunc(records, func(o []byte) bool { return keyOf(o) == keyOf(more) }) {
				continue
			}
			records = append(records, more)
		}
		if len(records) == len(c[l.input]) {
			return nil, false
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
// that has it, give or take a delta, for another field: the digits of a
// part of it, for a field of as many digits, a month a program takes from
// its command line for the month a record holds.
func (p *program) valueOf(c Case, from target, to *layout.Field, delta int64) (value, bool) {
	in := p.input(from.input)
	s := p.shapes[from.input]
	for _, r := range c[from.input] {
		if !s.has(in, r, from.field) {
			continue
		}
		if from.part != (part{}) {
			start := from.field.Start - in.root.Start + from.part.offset
			n, err := strconv.ParseInt(string(r[start:start+from.part.length]), 10, 64)
			return number(n + delta), err == nil
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

// lastOfKey is the last field of an indexed input's record key, when the
// key is several fields and that one a number: a coverage's number under
// its policy's, which one more record of the same policy counts on.
func (in input) lastOfKey() (*layout.Field, bool) {
	if in.file.RecordKey == "" {
		return nil, false
	}
	key := find(in.root, in.file.RecordKey)
	if key == nil {
		return nil, false
	}
	parts := leaves(key)
	if len(parts) < 2 || !numeric(parts[len(parts)-1]) {
		return nil, false
	}
	return parts[len(parts)-1], true
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
