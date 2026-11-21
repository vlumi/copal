package oracle

import (
	"slices"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/dataflow"
	"github.com/vlumi/copal/survey/internal/rules"
)

// marks finds, once, why no run of the program can reach some rules of
// the programs it calls: a rule that takes a condition name of a
// parameter to be what no call makes it, a direction of conversion the
// program never asks for, and the rules of the paragraphs performed only
// from such rules.
func (p *program) marks() map[rule]mark {
	if p.marked != nil {
		return p.marked
	}
	p.marked = map[rule]mark{}
	for u, unit := range p.units {
		for t, table := range unit.Tables {
			for i, r := range table.Rules {
				if why := unit.mark(table, r); why != unmarked {
					p.marked[rule{u, t, i}] = why
				} else if u > 0 && p.takesUnpassed(unit, table, r) {
					p.marked[rule{u, t, i}] = unpassed
				}
			}
		}
	}
	for u, unit := range p.units {
		if u > 0 {
			p.markUnperformed(u, unit)
		}
	}
	return p.marked
}

// takesUnpassed reports whether a rule of a called unit's table takes a
// condition name of a parameter, tested before anything could change it,
// to be true where no call passes a value of it, or false where every
// call does.
func (p *program) takesUnpassed(u Unit, t rules.Table, r rules.Rule) bool {
	for i, entry := range r.Entries {
		c := t.Tests[i].Condition
		if entry == "" || c == nil || c.Op != "NAME" || t.Tests[i].Time > 1 {
			continue
		}
		field := conditionField(u.records, c.Left.Name)
		if field == nil {
			continue
		}
		passed, ok := p.passed(u, field)
		if !ok {
			continue
		}
		named := field.Conditions[slices.IndexFunc(field.Conditions, func(n layout.Condition) bool { return n.Name == c.Left.Name })]
		holds := func(v value) bool {
			return slices.ContainsFunc(named.Values, func(lit string) bool {
				w, ok := literal(field, lit)
				return ok && w == v
			})
		}
		if entry == "Y" && !slices.ContainsFunc(passed, holds) || entry == "N" && !slices.ContainsFunc(passed, func(v value) bool { return !holds(v) }) {
			return true
		}
	}
	return false
}

// passed lists the values the calls of a unit pass in a field of its
// parameters, when they can be told: the unit is called by the program
// run alone and does not write the field itself, and at every call the
// argument's field over the same bytes holds only literals, the value it
// is declared with among them.
func (p *program) passed(u Unit, field *layout.Field) ([]value, bool) {
	calledElsewhere := slices.ContainsFunc(p.units[1:], func(v Unit) bool {
		return slices.ContainsFunc(v.blueprint.Calls, func(c blueprint.Call) bool { return c.Program == u.Program })
	})
	writes := slices.ContainsFunc(u.flow.Origins(field.Name), func(o dataflow.Origin) bool { return o.Parameter == 0 })
	if calledElsewhere || writes {
		return nil, false
	}
	var values []value
	located := false
	for n, using := range u.blueprint.Using {
		record, inner, ok := locate(u.blueprint, using, field.Name)
		if !ok {
			continue
		}
		located = true
		for _, call := range p.blueprint.Calls {
			if call.Program != u.Program {
				continue
			}
			if len(call.Using) <= n {
				return nil, false
			}
			argument, ok := fieldOf(p.blueprint, call.Using[n])
			if !ok {
				return nil, false
			}
			name, ok := over(argument, argument.Start+inner.Start-record.Start, inner.Length)
			if !ok {
				return nil, false
			}
			held, ok := p.literalsOf(name, field)
			if !ok {
				return nil, false
			}
			values = append(values, held...)
		}
	}
	return values, located && len(values) > 0
}

// literalsOf lists the values an item of the program run holds, read as a
// field's, when all it is ever given are literals: those moved into it,
// the condition names set on it, and the VALUE it is declared with.
func (p *program) literalsOf(item string, as *layout.Field) ([]value, bool) {
	var values []value
	for _, o := range p.flow.Origins(item) {
		lit := o.Literal
		switch {
		case o.File != "" || o.Call != "" || o.Stream != "" || o.Parameter != 0:
			return nil, false
		case strings.HasPrefix(lit, "condition "):
			named := conditionField(p.records, strings.TrimPrefix(lit, "condition "))
			if named == nil {
				return nil, false
			}
			i := slices.IndexFunc(named.Conditions, func(n layout.Condition) bool { return "condition "+n.Name == lit })
			lit = named.Conditions[i].Values[0]
		case lit == "":
			declared := p.declared(o.Item)
			if declared == nil || len(declared.Values) != 1 {
				return nil, false
			}
			lit = declared.Values[0]
		}
		v, ok := literal(as, lit)
		if !ok {
			return nil, false
		}
		values = append(values, v)
	}
	return values, true
}

// declared finds an item of the program run by name.
func (p *program) declared(name string) *layout.Field {
	for _, r := range p.records.All() {
		if f := find(r, name); f != nil {
			return f
		}
	}
	return nil
}

// markUnperformed marks the rules of the paragraphs of a called unit that
// are performed only from rules already marked, and so on from those. It
// holds for a unit whose paragraphs are entered by PERFORM alone: no GO
// TO, no PERFORM THRU, and a first paragraph that ends the program.
func (p *program) markUnperformed(u int, unit Unit) {
	if !performedOnly(unit.source) {
		return
	}
	dead := map[string]bool{}
	for again := true; again; {
		again = false
		for _, para := range unit.source.Paragraphs[1:] {
			if dead[para.Name] || !p.onlyFromMarked(u, unit, para.Name, dead) {
				continue
			}
			dead[para.Name], again = true, true
			for t, table := range unit.Tables {
				for i := range table.Rules {
					if _, ok := p.marked[rule{u, t, i}]; !ok && table.Paragraph == para.Name {
						p.marked[rule{u, t, i}] = unperformed
					}
				}
			}
		}
	}
}

// onlyFromMarked reports whether every PERFORM of a paragraph, at least
// one, is in a paragraph already found dead, or on the paths of marked
// rules alone.
func (p *program) onlyFromMarked(u int, unit Unit, name string, dead map[string]bool) bool {
	performs, live := 0, false
	for _, para := range unit.source.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if !slices.ContainsFunc(s.Procedures, func(ref cobol.ProcedureRef) bool { return ref.Name == name }) {
				return
			}
			performs++
			if !dead[para.Name] && !p.onMarkedPathsOnly(u, unit, s.Pos) {
				live = true
			}
		})
	}
	return performs > 0 && !live
}

// onMarkedPathsOnly reports whether a statement is on the path of some
// rule of the unit's tables, and of marked rules alone.
func (p *program) onMarkedPathsOnly(u int, unit Unit, pos cobol.Pos) bool {
	on := 0
	for t, table := range unit.Tables {
		for i, r := range table.Rules {
			if !slices.Contains(r.Steps, pos) {
				continue
			}
			if _, ok := p.marked[rule{u, t, i}]; !ok {
				return false
			}
			on++
		}
	}
	return on > 0
}

// performedOnly reports whether a program's paragraphs after the first
// are entered by PERFORM alone.
func performedOnly(source *cobol.Program) bool {
	if len(source.Paragraphs) == 0 || len(source.Paragraphs[0].Statements) == 0 {
		return false
	}
	last := source.Paragraphs[0].Statements[len(source.Paragraphs[0].Statements)-1]
	ends := last.Verb == "GOBACK" || last.Verb == "STOP" || last.Verb == "EXIT" && strings.Contains(last.Text, "PROGRAM")
	jumps := false
	for _, para := range source.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			jumps = jumps || s.Verb == "GO" || s.Verb == "ALTER" ||
				slices.ContainsFunc(s.Procedures, func(ref cobol.ProcedureRef) bool { return ref.Thru != "" || ref.Section != "" })
		})
	}
	return ends && !jumps
}
