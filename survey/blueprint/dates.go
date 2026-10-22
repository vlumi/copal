package blueprint

import "slices"

// FieldRef names a field of an entity.
type FieldRef struct {
	Record string `json:"record"`
	Field  string `json:"field"`
}

// DateLayout is a layout the entities hold dates in, and the fields in it:
// the date fields, and the fields of any kind whose values come from a
// date field of another layout.
type DateLayout struct {
	Layout string
	Fields []FieldRef
}

// heldRecords are the entities programs hold records of, by program and
// the program's name for the record.
type heldRecords map[[2]string]int

// layoutOf writes how a field holds its value: its picture, and its usage
// unless it is DISPLAY.
func layoutOf(f Field) string {
	if f.Usage == "" || f.Usage == "DISPLAY" {
		return f.Picture
	}
	return f.Picture + " " + f.Usage
}

// datesFrom finds, for each date field of each entity, and each text field
// that is not a redefinition, the date fields of the entities in another
// layout that its values come from, by the origins of the field in each
// program that writes a record of the entity: one business date in a
// six-digit form, an eight-digit one and text, say.
func datesFrom(entities []Entity, programs []Program, held heldRecords) {
	for _, p := range programs {
		for _, out := range p.Outputs {
			e, ok := held[[2]string{p.Name, out.Record}]
			if !ok {
				continue
			}
			for _, fo := range out.Fields {
				j, f, ok := entityFieldIndex(p, out.Record, fo.Name)
				if !ok || !holdsDate(f) {
					continue
				}
				field := &entities[e].Fields[j]
				for _, o := range fo.origins {
					if o.File == "" || o.Through != "" {
						continue
					}
					source, ref, ok := heldField(entities, p, held, o.Item)
					if ok && source.Kind == KindDate && source.Layout != field.Layout && !slices.Contains(field.DatesFrom, ref) {
						field.DatesFrom = append(field.DatesFrom, ref)
					}
				}
			}
		}
	}
}

// heldField finds the entity field a program's item is, through the
// program's record of the entity that holds it.
func heldField(entities []Entity, p Program, held heldRecords, item string) (EntityField, FieldRef, bool) {
	for _, r := range p.Records {
		e, ok := held[[2]string{p.Name, r.Name}]
		if !ok {
			continue
		}
		if j, _, ok := entityFieldIndex(p, r.Name, item); ok {
			f := entities[e].Fields[j]
			return f, FieldRef{Record: entities[e].Record, Field: f.Name}, true
		}
	}
	return EntityField{}, FieldRef{}, false
}

// entityFieldIndex finds a field among the fields of a program's record of
// an entity, which are the entity's, in the same order: where it is, and
// its layout.
func entityFieldIndex(p Program, record, field string) (int, Field, bool) {
	i := slices.IndexFunc(p.Records, func(r Record) bool { return r.Name == record })
	if i < 0 {
		return 0, Field{}, false
	}
	fields := entityFields(p.Records[i].Field)
	j := slices.IndexFunc(fields, func(f Field) bool { return f.Name == field })
	if j < 0 {
		return 0, Field{}, false
	}
	return j, fields[j], true
}

// holdsDate reports whether a field can hold a date: it is one, or it is
// text whose bytes are not those of another record type, neither
// redefining another field nor redefined, nor within a group that is.
func holdsDate(f Field) bool {
	return f.Kind == KindDate || f.Kind == KindText && !f.sharesBytes
}

// DateLayouts groups the date fields of the entities by layout, with the
// fields of any kind whose values come from a date of another layout, in
// the order the entities first hold each layout.
func DateLayouts(entities []Entity) []DateLayout {
	var layouts []DateLayout
	add := func(layout string, ref FieldRef) {
		i := slices.IndexFunc(layouts, func(l DateLayout) bool { return l.Layout == layout })
		if i < 0 {
			layouts = append(layouts, DateLayout{Layout: layout})
			i = len(layouts) - 1
		}
		if !slices.Contains(layouts[i].Fields, ref) {
			layouts[i].Fields = append(layouts[i].Fields, ref)
		}
	}
	for _, e := range entities {
		for _, f := range e.Fields {
			if f.Kind == KindDate || len(f.DatesFrom) > 0 {
				add(f.Layout, FieldRef{Record: e.Record, Field: f.Name})
			}
		}
	}
	return layouts
}

// datesInText gives the kind of date to each text field of an entity whose
// values come from a date: in the entity, in its copybook, and in each
// program's record of it.
func datesInText(entities []Entity, copybooks []Copybook, programs []Program, held heldRecords) {
	for e := range entities {
		for j, f := range entities[e].Fields {
			if f.Kind != KindText || len(f.DatesFrom) == 0 {
				continue
			}
			entities[e].Fields[j].Kind = KindDate
			for _, c := range copybooks {
				if c.Name == entities[e].Copybook {
					for r := range c.Records {
						setKind(&c.Records[r].Field, f.Name, KindDate)
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
					setKind(&programs[i].Records[r].Field, fields[j].Name, KindDate)
				}
			}
		}
	}
}

// setKind sets the kind of a field among a field's fields.
func setKind(f *Field, name, kind string) {
	if f.Name == name {
		f.Kind = kind
	}
	for i := range f.Fields {
		setKind(&f.Fields[i], name, kind)
	}
}
