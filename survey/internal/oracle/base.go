package oracle

import (
	"maps"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
	"github.com/vlumi/copal/survey/blueprint"
)

// businessDate is the date the base case's dates hold, as eight digits.
const businessDate = 20261016

// shape is how the records of an input are told apart and laid out: the
// field whose condition names the program tests to tell the record types
// apart, if any, the value of each type, and for each type, the
// redefinitions of the record that are its fields.
type shape struct {
	discriminator *layout.Field
	types         []value
	typeNames     []string
	variants      map[*layout.Field]int // redefinition to type, or anyType or noType
}

// The types of a redefinition that is no one type's: anyType is in every
// record whose type has no redefinition of its own, noType in none.
const (
	anyType = -1
	noType  = -2
)

// hasVariant reports whether a record of a type, -1 for none, has a
// redefinition of the given type.
func (s shape) hasVariant(variant, kind int) bool {
	return variant == kind || variant == anyType && !slices.Contains(slices.Collect(maps.Values(s.variants)), kind)
}

// shapeOf finds the shape of an input: the blueprint's record types, when
// it pairs each with a redefinition; or else the first field of the record
// with a condition name the program tests, whose condition names are then
// the types, and for each redefinition, the type whose rules use its
// fields most.
func (p *program) shapeOf(in input) shape {
	if s, ok := p.blueprintShape(in); ok {
		return s
	}
	tested := p.testedNames()
	s := shape{variants: map[*layout.Field]int{}}
	for _, f := range in.root.Children {
		if slices.ContainsFunc(f.Conditions, func(c layout.Condition) bool { return tested[c.Name] }) {
			s.discriminator = f
			break
		}
	}
	if s.discriminator != nil {
		for _, c := range s.discriminator.Conditions {
			if len(c.Values) == 0 {
				continue
			}
			if v, ok := literal(s.discriminator, c.Values[0]); ok && !slices.Contains(s.types, v) {
				s.types = append(s.types, v)
				s.typeNames = append(s.typeNames, c.Name)
			}
		}
	}
	for _, f := range in.root.Children {
		if f.Redefines != "" || slices.ContainsFunc(in.root.Children, func(g *layout.Field) bool { return g.Redefines == f.Name }) {
			s.variants[f] = anyType
		}
	}
	switch {
	case s.discriminator == nil && len(s.variants) > 0:
		var declared []*layout.Field
		for _, f := range in.root.Children {
			if _, ok := s.variants[f]; ok {
				declared = append(declared, f)
			}
		}
		used := p.mostUsed(declared)
		for f := range s.variants {
			if f != used {
				s.variants[f] = noType
			}
		}
	case s.discriminator == nil:
	default:
		for f := range s.variants {
			s.variants[f] = p.typeOf(f, s)
		}
		s.pairLast()
	}
	return s
}

// mostUsed is the field, of those given in the order they are declared,
// whose fields the paragraphs read and write most, the first on a tie.
func (p *program) mostUsed(fields []*layout.Field) *layout.Field {
	uses := map[string]int{}
	for _, para := range p.flow.Paragraphs {
		for _, name := range append(slices.Clone(para.Reads), para.Writes...) {
			uses[name]++
		}
	}
	var best *layout.Field
	bestUses := -1
	for _, f := range fields {
		n := 0
		walk(f, func(g *layout.Field) { n += uses[g.Name] })
		if n > bestUses {
			best, bestUses = f, n
		}
	}
	return best
}

// pairLast gives the one type no redefinition was found for the one
// redefinition left, when that is all that is left: the redefined field
// itself is no type's.
func (s shape) pairLast() {
	var types []int
	for k := range s.types {
		if !slices.ContainsFunc(slices.Collect(maps.Values(s.variants)), func(t int) bool { return t == k }) {
			types = append(types, k)
		}
	}
	var left []*layout.Field
	for f, k := range s.variants {
		if k < 0 && f.Redefines != "" {
			left = append(left, f)
		}
	}
	if len(types) == 1 && len(left) == 1 {
		s.variants[left[0]] = types[0]
	}
}

// blueprintShape is the shape the blueprint gives a record of types, each
// paired with its redefinition, or of a file of several records, each
// paired with its record: the redefinition no type has is in no record of
// a type.
func (p *program) blueprintShape(in input) (shape, bool) {
	i := slices.IndexFunc(p.blueprint.Records, func(r blueprint.Record) bool {
		return r.Name == in.root.Name || r.File != "" && r.File == in.file.Name && len(r.Types) > 0
	})
	if i < 0 {
		return shape{}, false
	}
	r := p.blueprint.Records[i]
	if len(r.Types) == 0 || slices.ContainsFunc(r.Types, func(t blueprint.RecordType) bool { return t.Fields == "" }) {
		return shape{}, false
	}
	s := shape{discriminator: find(in.root, r.TypeField), variants: map[*layout.Field]int{}}
	for _, f := range in.root.Children {
		if f.Redefines != "" || slices.ContainsFunc(in.root.Children, func(g *layout.Field) bool { return g.Redefines == f.Name }) {
			s.variants[f] = anyType
		}
	}
	for k, t := range r.Types {
		v, ok := literal(s.discriminator, t.Value)
		if !ok {
			return shape{}, false
		}
		s.types, s.typeNames = append(s.types, v), append(s.typeNames, t.Name)
		s.variants[find(in.root, t.Fields)] = k
	}
	return s, true
}

// testedNames are the data items and condition names the program's
// conditions test.
func (p *program) testedNames() map[string]bool {
	tested := map[string]bool{}
	for _, c := range p.conditions() {
		for _, r := range c.Refs() {
			tested[r] = true
		}
	}
	return tested
}

// typeOf is the record type whose rules use the fields of a redefinition
// most: the rules that find the record of that type, by its condition
// name, and refer to one of the fields; -1 when no rule does.
func (p *program) typeOf(variant *layout.Field, s shape) int {
	fields := map[string]bool{}
	walk(variant, func(f *layout.Field) { fields[f.Name] = true })
	return p.typeUsing(fields, s)
}

// typeUsing is the record type whose rules refer to the names most.
func (p *program) typeUsing(fields map[string]bool, s shape) int {
	statements := p.statementsByPos()
	scores := make([]int, len(s.types))
	for _, t := range p.tables {
		for _, r := range t.Rules {
			kind := -1
			for i, entry := range r.Entries {
				if entry == "Y" {
					if k := slices.Index(s.typeNames, t.Conditions[i]); k >= 0 {
						kind = k
					}
				}
			}
			if kind >= 0 && slices.ContainsFunc(p.ruleNames(t.Refs, r.Entries, r.Steps, statements), func(n string) bool { return fields[n] }) {
				scores[kind]++
			}
		}
	}
	best := -1
	for k, score := range scores {
		if score > 0 && (best < 0 || score > scores[best]) {
			best = k
		}
	}
	return best
}

// ruleNames lists the names a rule refers to: in the conditions it tests,
// in the statements on its path, and in the paragraphs they perform.
func (p *program) ruleNames(refs [][]string, entries []string, steps []cobol.Pos, statements map[cobol.Pos]cobol.Statement) []string {
	var names []string
	for i, entry := range entries {
		if entry != "" {
			names = append(names, refs[i]...)
		}
	}
	for _, pos := range steps {
		s := statements[pos]
		for _, op := range append(slices.Clone(s.Reads), s.Writes...) {
			names = append(names, op.Refs...)
			if op.Name != "" {
				names = append(names, op.Name)
			}
		}
		if s.Condition != nil {
			names = append(names, s.Condition.Refs()...)
		}
		names = append(names, p.performedNames(s, map[string]bool{})...)
	}
	return names
}

// performedNames lists the names the paragraphs a statement performs read
// and write, and those of the paragraphs they perform in turn.
func (p *program) performedNames(s cobol.Statement, seen map[string]bool) []string {
	var names []string
	if s.Verb == "GO" {
		return nil
	}
	for _, ref := range s.Procedures {
		for i, para := range p.source.Paragraphs {
			if para.Name != ref.Name && para.Name != ref.Thru || seen[para.Name] {
				continue
			}
			seen[para.Name] = true
			names = append(names, p.flow.Paragraphs[i].Reads...)
			names = append(names, p.flow.Paragraphs[i].Writes...)
			cobol.Walk(para.Statements, func(n cobol.Statement) { names = append(names, p.performedNames(n, seen)...) })
		}
	}
	return names
}

func (p *program) statementsByPos() map[cobol.Pos]cobol.Statement {
	statements := map[cobol.Pos]cobol.Statement{}
	for _, para := range p.source.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) { statements[s.Pos] = s })
	}
	return statements
}

func walk(f *layout.Field, visit func(*layout.Field)) {
	visit(f)
	for _, child := range f.Children {
		walk(child, visit)
	}
}

// base is the first case: for each input, a record of each of its types,
// or a single record, filled field by field by kind.
func (p *program) base() Case {
	kinds := p.kinds()
	c := Case{}
	for _, in := range p.inputs {
		s := p.shapes[in.name()]
		types := max(len(s.types), 1)
		for k := range types {
			index, kind := 0, k
			if len(s.variants) == 0 {
				index = k
			}
			if len(s.types) == 0 {
				kind = -1
			}
			c[in.name()] = append(c[in.name()], p.typical(in, s, kind, index, kinds))
		}
	}
	p.balance(c, nil)
	return c
}

// typical is a record of an input of the given type: its type's value,
// and its fields, those of its type's redefinition among them, filled by
// kind. Keys and numbers count up with the record's index: among the
// records of its type, for types that redefine the record, so that the
// first records of two files refer to each other; among all the records
// of the input otherwise, so that its keys are unique.
func (p *program) typical(in input, s shape, kind, index int, kinds map[string]string) []byte {
	r := in.layout.Blank()
	for _, f := range in.root.Children {
		if t, ok := s.variants[f]; ok && !s.hasVariant(t, kind) {
			continue
		}
		walk(f, func(leaf *layout.Field) {
			if len(leaf.Children) > 0 || leaf.Name == "FILLER" {
				return
			}
			if v, ok := typicalValue(leaf, kinds[leaf.Name], index); ok {
				set(r, in.root, leaf, in.layout, v)
			}
		})
	}
	if s.discriminator != nil {
		set(r, in.root, s.discriminator, in.layout, s.types[kind])
	}
	return r
}

// typicalValue is a field's value in the base case, by its kind.
func typicalValue(f *layout.Field, kind string, index int) (value, bool) {
	switch {
	case len(f.Conditions) > 0 && len(f.Conditions[0].Values) > 0:
		return literal(f, f.Conditions[0].Values[0])
	case kind == blueprint.KindDate:
		return dateValue(f), true
	case kind == blueprint.KindAmount:
		return number(10 * pow10(f.Scale)), true
	case kind == blueprint.KindKey, kind == blueprint.KindNumber:
		return number(int64(index + 1)), true
	case kind == blueprint.KindText:
		return text("A"), true
	}
	return value{}, false
}

// dateValue is the business date in a date field's form: eight digits, a
// month in six, or six digits, or fourteen with a time of day.
func dateValue(f *layout.Field) value {
	switch {
	case f.Digits == 14:
		return number(businessDate*1000000 + 120000)
	case f.Digits == 6 && (strings.Contains(f.Name, "MONTH") || strings.Contains(f.Name, "YYMM")):
		return number(businessDate / 100)
	case f.Digits == 6:
		return number(businessDate % 1000000)
	}
	return number(businessDate)
}

func pow10(n int) int64 {
	p := int64(1)
	for range n {
		p *= 10
	}
	return p
}

// kinds are the kinds of the program's fields, by name.
func (p *program) kinds() map[string]string {
	kinds := map[string]string{}
	var visit func(f blueprint.Field)
	visit = func(f blueprint.Field) {
		kinds[f.Name] = f.Kind
		for _, child := range f.Fields {
			visit(child)
		}
	}
	for _, r := range p.blueprint.Records {
		visit(r.Field)
	}
	return kinds
}
