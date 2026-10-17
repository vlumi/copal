package blueprint

import "slices"

// Entity is a record a copybook describes: the datasets whose records
// the copybook describes, and what the programs that use the copybook for
// a record of their own do with each of its fields.
type Entity struct {
	Record   string        `json:"record"`
	Copybook string        `json:"copybook"`
	Datasets []string      `json:"datasets,omitempty"`
	Programs []string      `json:"programs,omitempty"`
	Fields   []EntityField `json:"fields"`
}

// EntityField is a field of an entity, FILLER left out, with its kind,
// its uses, one for each program that uses the entity, and for a code
// field, its state machine.
type EntityField struct {
	Name   string        `json:"name"`
	Kind   string        `json:"kind"`
	Uses   []FieldUse    `json:"uses,omitempty"`
	States *StateMachine `json:"states,omitempty"`
}

// FieldUse is what a program does with a field of an entity: the name it
// knows the field by, which differs where the copybook is copied with
// REPLACING, the paragraphs that read it and write it by that name, and
// the decisions on it.
type FieldUse struct {
	Program   string     `json:"program"`
	Field     string     `json:"field"`
	ReadBy    []string   `json:"readBy,omitempty"`
	WrittenBy []string   `json:"writtenBy,omitempty"`
	Decisions []Decision `json:"decisions,omitempty"`
}

// entities lists the records of every copybook as entities, in the order
// of the copybooks, with the programs that use each. A program's record
// is the copybook's when it is declared at the same place in the source,
// and its fields are the copybook's, in the same order.
func entities(copybooks []Copybook, programs []Program) []Entity {
	var list []Entity
	for _, c := range copybooks {
		for _, r := range c.Records {
			e := Entity{Record: r.Name, Copybook: c.Name}
			fields := entityFields(r.Field)
			for _, f := range fields {
				ef := EntityField{Name: f.Name, Kind: f.Kind}
				if f.Kind == KindCode {
					ef.States = newStateMachine(f)
				}
				e.Fields = append(e.Fields, ef)
			}
			for k := range programs {
				p := &programs[k]
				i := slices.IndexFunc(p.Records, func(pr Record) bool { return pr.Copybook == c.Name && pr.Source == r.Source })
				if i < 0 {
					continue
				}
				e.Programs = append(e.Programs, p.Name)
				for j, f := range entityFields(p.Records[i].Field) {
					e.Fields[j].Uses = append(e.Fields[j].Uses, fieldUse(*p, f.Name))
					if m := e.Fields[j].States; m != nil {
						m.add(newCodeField(p, f, fields[j]))
					}
				}
			}
			list = append(list, e)
		}
	}
	return list
}

// describeEntities lists the datasets each entity's copybook describes.
func describeEntities(entities []Entity, datasets []SystemDataset) {
	for i := range entities {
		for _, d := range datasets {
			if d.Copybook == entities[i].Copybook {
				entities[i].Datasets = append(entities[i].Datasets, d.Name)
			}
		}
	}
}

// entityFields lists a record's fields under it, in order, FILLER left
// out.
func entityFields(record Field) []Field {
	var fields []Field
	for _, child := range record.Fields {
		walk(child, func(f Field) {
			if f.Name != "FILLER" {
				fields = append(fields, f)
			}
		})
	}
	return fields
}

func fieldUse(p Program, field string) FieldUse {
	u := FieldUse{Program: p.Name, Field: field, Decisions: p.Decisions(field)}
	for _, para := range p.Paragraphs {
		if slices.Contains(para.Reads, field) {
			u.ReadBy = append(u.ReadBy, para.Name)
		}
		if slices.Contains(para.Writes, field) {
			u.WrittenBy = append(u.WrittenBy, para.Name)
		}
	}
	return u
}

// Field finds a field of the entity by name.
func (e Entity) Field(name string) (EntityField, bool) {
	i := slices.IndexFunc(e.Fields, func(f EntityField) bool { return f.Name == name })
	if i < 0 {
		return EntityField{}, false
	}
	return e.Fields[i], true
}
