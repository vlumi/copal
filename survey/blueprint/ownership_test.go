package blueprint

import "testing"

// An entity's writers point at it with solid edges, and it points at its
// readers with dotted ones; a program is drawn once, and an entity no job keeps not at
// all.
func TestOwnershipGraph(t *testing.T) {
	g := OwnershipGraph([]Entity{
		{Record: "A-REC", Copybook: "A", WrittenBy: []string{"P"}, ReadBy: []string{"Q"}},
		{Record: "B-REC", Copybook: "B", WrittenBy: []string{"Q"}},
		{Record: "C-REC", Copybook: "C"},
	}, OwnershipNotes{})
	if len(g.Nodes) != 4 {
		t.Errorf("nodes %+v, want P, Q, A-REC and B-REC", g.Nodes)
	}
	want := map[[2]string]bool{{"program P", "entity A-REC"}: false, {"entity A-REC", "program Q"}: true, {"program Q", "entity B-REC"}: false}
	for _, e := range g.Edges {
		dotted, ok := want[[2]string{e.From, e.To}]
		if !ok || dotted != (e.Style == "dotted") {
			t.Errorf("edge %+v", e)
		}
	}
	if len(g.Edges) != len(want) {
		t.Errorf("edges %+v", g.Edges)
	}
}
