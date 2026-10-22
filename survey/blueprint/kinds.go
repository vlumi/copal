package blueprint

import (
	"regexp"
	"slices"

	"github.com/vlumi/cobolith/cobol"
)

// kindsByUse carries the kinds the code shows along its data flow: a
// numeric field an amount is moved or computed into is an amount, a field
// a date is moved into, of a size a date takes, is a date, and so is an
// argument of INTEGER-OF-DATE; and a field of an entity is of
// the kind any program's record of it shows, which every program's is
// then. A date's arithmetic is not followed: the days between two dates
// are a number.
func kindsByUse(entities []Entity, copybooks []Copybook, programs []Program, held heldRecords) {
	for range 3 {
		changed := false
		for i := range programs {
			changed = kindsInProgram(&programs[i]) || changed
		}
		changed = kindsAcrossEntities(entities, copybooks, programs, held) || changed
		if !changed {
			return
		}
	}
}

// dateFunctions are the intrinsic functions whose argument is a date;
// DATE-OF-INTEGER's is a day number.
var dateFunctions = regexp.MustCompile(`(?:INTEGER-OF-DATE|INTEGER-OF-DAY)\s*\(`)

// kindsInProgram carries amounts and dates along a program's lineage,
// and reports whether any field's kind changed.
func kindsInProgram(p *Program) bool {
	kinds := map[string]*Field{}
	for i := range p.Records {
		indexFields(&p.Records[i].Field, kinds)
	}
	changed := false
	set := func(name, kind string) {
		if f, ok := kinds[name]; ok && f.Kind != kind {
			f.Kind = kind
			changed = true
		}
	}
	if p.source != nil {
		for _, para := range p.source.Paragraphs {
			cobol.Walk(para.Statements, func(s cobol.Statement) {
				for _, op := range s.Reads {
					if dateFunctions.MatchString(op.Text) {
						for _, ref := range op.Refs {
							if f, ok := kinds[ref]; ok && f.Kind == KindNumber {
								set(ref, KindDate)
							}
						}
					}
				}
			})
		}
	}
	for again := true; again; {
		again = false
		for _, e := range p.Lineage {
			from, to := kinds[e.From], kinds[e.To]
			if e.FromLiteral || e.FromFile || e.FromCall || e.FromStream || e.ToFile || e.ToCall || e.ToStream || from == nil || to == nil {
				continue
			}
			switch {
			case from.Kind == KindAmount && to.Kind == KindNumber:
				to.Kind = KindAmount
			case from.Kind == KindDate && e.Verb == "MOVE" && (to.Kind == KindNumber && slices.Contains([]int{6, 8, 14}, to.Digits) || to.Kind == KindText && to.Length >= 6 && to.Length <= 8):
				to.Kind = KindDate
			default:
				continue
			}
			again, changed = true, true
		}
	}
	return changed
}

// indexFields indexes a record's fields by name, the first of a name.
func indexFields(f *Field, kinds map[string]*Field) {
	if _, ok := kinds[f.Name]; !ok && f.Name != "FILLER" {
		kinds[f.Name] = f
	}
	for i := range f.Fields {
		indexFields(&f.Fields[i], kinds)
	}
}

// kindsAcrossEntities gives each field of an entity the kind, amount or
// date, that any program's record of the entity shows, and that kind to
// the entity's copybook and to every program's record of it, and reports
// whether any changed.
func kindsAcrossEntities(entities []Entity, copybooks []Copybook, programs []Program, held heldRecords) bool {
	changed := false
	for e := range entities {
		for j := range entities[e].Fields {
			field := &entities[e].Fields[j]
			if field.Kind != KindNumber && field.Kind != KindText {
				continue
			}
			kind := ""
			for key, k := range held {
				if k != e {
					continue
				}
				i := slices.IndexFunc(programs, func(p Program) bool { return p.Name == key[0] })
				r := slices.IndexFunc(programs[i].Records, func(r Record) bool { return r.Name == key[1] })
				if fields := entityFields(programs[i].Records[r].Field); j < len(fields) && (fields[j].Kind == KindAmount || fields[j].Kind == KindDate) {
					kind = fields[j].Kind
				}
			}
			if kind == "" {
				continue
			}
			field.Kind, changed = kind, true
			for c := range copybooks {
				if copybooks[c].Name == entities[e].Copybook {
					for r := range copybooks[c].Records {
						setKind(&copybooks[c].Records[r].Field, field.Name, kind)
					}
				}
			}
			for key, k := range held {
				if k != e {
					continue
				}
				i := slices.IndexFunc(programs, func(p Program) bool { return p.Name == key[0] })
				r := slices.IndexFunc(programs[i].Records, func(r Record) bool { return r.Name == key[1] })
				if fields := entityFields(programs[i].Records[r].Field); j < len(fields) {
					setKind(&programs[i].Records[r].Field, fields[j].Name, kind)
				}
			}
		}
	}
	return changed
}
