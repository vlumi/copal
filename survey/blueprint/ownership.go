package blueprint

import "github.com/vlumi/copal/survey/internal/diagram"

// OwnershipNotes say where an ownership diagram's entities and programs
// link to. Either may be nil.
type OwnershipNotes struct {
	EntityLink  func(string) string
	ProgramLink func(string) string
}

// OwnershipGraph draws who owns the entities that are kept in datasets:
// each program that writes a dataset of an entity, with a solid edge to
// the entity, and each program that only reads one, with a dotted edge
// from it. The edges run the way the data does, so that an entity sits
// between the programs that write it and those that read it.
func OwnershipGraph(entities []Entity, notes OwnershipNotes) diagram.Graph {
	g := diagram.Graph{Legend: []diagram.Key{
		diagram.NodeKey(diagram.Node{}, "program"),
		diagram.NodeKey(diagram.Node{Data: true}, "entity, its copybook"),
		diagram.EdgeKey(diagram.Solid, "writes a dataset of it"),
		diagram.EdgeKey(diagram.Dotted, "only read by the program"),
	}}
	added := map[string]bool{}
	program := func(name string) string {
		id := "program " + name
		if !added[id] {
			added[id] = true
			n := diagram.Node{ID: id, Label: name}
			if notes.ProgramLink != nil {
				n.Link = notes.ProgramLink(name)
			}
			g.Nodes = append(g.Nodes, n)
		}
		return id
	}
	for _, e := range entities {
		if len(e.WrittenBy)+len(e.ReadBy) == 0 {
			continue
		}
		id := "entity " + e.Record
		n := diagram.Node{ID: id, Label: e.Record + "\ncopybook " + e.Copybook, Data: true}
		if notes.EntityLink != nil {
			n.Link = notes.EntityLink(e.Record)
		}
		for _, p := range e.WrittenBy {
			g.Edges = append(g.Edges, diagram.Edge{From: program(p), To: id})
		}
		for _, p := range e.ReadBy {
			g.Edges = append(g.Edges, diagram.Edge{From: id, To: program(p), Style: diagram.Dotted})
		}
		g.Nodes = append(g.Nodes, n)
	}
	return g
}
