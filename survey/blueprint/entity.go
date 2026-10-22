package blueprint

import "slices"

// Entity is a record a copybook describes: the datasets whose records
// the copybook describes, the programs that write one of them and those
// that only read them, and what the programs that use the copybook for a
// record of their own do with each of its fields.
type Entity struct {
	Record    string        `json:"record"`
	Copybook  string        `json:"copybook"`
	Datasets  []string      `json:"datasets,omitempty"`
	WrittenBy []string      `json:"writtenBy,omitempty"`
	ReadBy    []string      `json:"readBy,omitempty"`
	Programs  []string      `json:"programs,omitempty"`
	Fields    []EntityField `json:"fields"`
	Links     []EntityLink  `json:"links,omitempty"`
	Totals    []EntityTotal `json:"totals,omitempty"`
}

// EntityField is a field of an entity, FILLER left out, with its kind and
// layout, its uses, one for each program that uses the entity, for a code
// field, its state machine, and the date fields of other layouts its
// values come from.
type EntityField struct {
	Name      string        `json:"name"`
	Kind      string        `json:"kind"`
	Layout    string        `json:"layout"`
	Uses      []FieldUse    `json:"uses,omitempty"`
	States    *StateMachine `json:"states,omitempty"`
	DatesFrom []FieldRef    `json:"datesFrom,omitempty"`
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
	held := heldRecords{}
	for _, c := range copybooks {
		for _, r := range c.Records {
			e := Entity{Record: r.Name, Copybook: c.Name}
			fields := entityFields(r.Field)
			for _, f := range fields {
				ef := EntityField{Name: f.Name, Kind: f.Kind, Layout: layoutOf(f)}
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
				held[[2]string{p.Name, p.Records[i].Name}] = len(list)
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
	kindsByUse(list, copybooks, programs, held)
	datesFrom(list, programs, held)
	datesInText(list, copybooks, programs, held)
	relate(list, programs, held)
	return list
}

// describeEntities lists the datasets each entity's copybook describes,
// the programs that write any of them, run by the job steps that do, and
// those that read them without writing any.
func describeEntities(entities []Entity, datasets []SystemDataset, jobs []Job) {
	programs := map[StepRef]string{}
	for _, job := range jobs {
		for _, step := range job.Steps {
			programs[StepRef{Job: job.Name, Step: step.Name}] = step.Program
		}
	}
	for i := range entities {
		e := &entities[i]
		var readers []string
		for _, d := range datasets {
			if d.Copybook != e.Copybook {
				continue
			}
			e.Datasets = append(e.Datasets, d.Name)
			for _, s := range d.WrittenBy {
				addOnce(&e.WrittenBy, programs[s])
			}
			for _, s := range d.ReadBy {
				addOnce(&readers, programs[s])
			}
		}
		for _, r := range readers {
			if !slices.Contains(e.WrittenBy, r) {
				e.ReadBy = append(e.ReadBy, r)
			}
		}
		slices.Sort(e.WrittenBy)
		slices.Sort(e.ReadBy)
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
