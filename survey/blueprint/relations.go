package blueprint

import "slices"

// EntityLink is a link between the fields of two entities, as the
// programs that rely on it use it: "key", a field moved into the other's
// record key to read its record by; "equal", fields compared for equality;
// "merge", fields compared for order, the two read side by side in it.
type EntityLink struct {
	Kind     string   `json:"kind"`
	From     FieldRef `json:"from"`
	To       FieldRef `json:"to"`
	Programs []string `json:"programs"`
}

// EntityTotal is a control total of an entity's records, as the programs
// that check it check it: a field that is the sum of another over the
// records of a type, or with no Sum, how many of them there are.
type EntityTotal struct {
	Field    string   `json:"field"`
	Sum      string   `json:"sum,omitempty"`
	Type     string   `json:"type,omitempty"`
	Programs []string `json:"programs"`
}

// relate gives each entity the links its fields are in, and the control
// totals of its records, from the programs that hold its records.
func relate(entities []Entity, programs []Program, held heldRecords) {
	for _, p := range programs {
		for _, l := range p.Links {
			_, from, fok := heldField(entities, p, held, l.From.Field)
			_, to, tok := heldField(entities, p, held, l.To.Field)
			if !fok || !tok {
				continue
			}
			link := EntityLink{Kind: l.Kind, From: from, To: to}
			for i := range entities {
				if entities[i].Record == from.Record || entities[i].Record == to.Record {
					entities[i].Links = addLink(entities[i].Links, link, p.Name)
				}
			}
		}
		for _, t := range p.Totals {
			_, field, ok := heldField(entities, p, held, t.Field)
			if !ok {
				continue
			}
			total := EntityTotal{Field: field.Field, Type: t.Type}
			if t.Sum != "" {
				_, sum, found := heldField(entities, p, held, t.Sum)
				if !found {
					continue
				}
				total.Sum = sum.Field
			}
			e := &entities[slices.IndexFunc(entities, func(e Entity) bool { return e.Record == field.Record })]
			i := slices.IndexFunc(e.Totals, func(u EntityTotal) bool { return u.Field == total.Field && u.Sum == total.Sum && u.Type == total.Type })
			if i < 0 {
				e.Totals = append(e.Totals, total)
				i = len(e.Totals) - 1
			}
			addOnce(&e.Totals[i].Programs, p.Name)
		}
	}
}

func addLink(links []EntityLink, l EntityLink, program string) []EntityLink {
	i := slices.IndexFunc(links, func(m EntityLink) bool { return m.Kind == l.Kind && m.From == l.From && m.To == l.To })
	if i < 0 {
		links = append(links, l)
		i = len(links) - 1
	}
	addOnce(&links[i].Programs, program)
	return links
}
