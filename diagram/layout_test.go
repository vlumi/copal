package diagram

import (
	"bytes"
	"strings"
	"testing"
)

// program is shaped like a batch program: a main paragraph performing
// steps, which share helpers, one of them two layers down from main.
func program() Graph {
	nodes := []Node{{ID: "MAIN", Label: "0000-MAIN", Emphasis: true}}
	for _, id := range []string{"OPEN", "READ", "PROCESS", "REJECT", "REPORT", "CLOSE", "ABORT"} {
		nodes = append(nodes, Node{ID: id, Label: id, Title: "does " + strings.ToLower(id), Link: "#" + id})
	}
	edges := []Edge{
		{From: "MAIN", To: "OPEN"}, {From: "MAIN", To: "PROCESS"}, {From: "MAIN", To: "REPORT"}, {From: "MAIN", To: "CLOSE"},
		{From: "OPEN", To: "ABORT"}, {From: "PROCESS", To: "READ"}, {From: "PROCESS", To: "REJECT"},
		{From: "READ", To: "ABORT"}, {From: "ABORT", To: "CLOSE"}, {From: "PROCESS", To: "READ"},
	}
	return Graph{Nodes: nodes, Edges: edges}
}

func TestLayoutLayers(t *testing.T) {
	l, err := newLayout(program())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int{"MAIN": 0, "OPEN": 1, "PROCESS": 1, "REPORT": 1, "READ": 2, "REJECT": 2, "ABORT": 3, "CLOSE": 4}
	for i, s := range l.nodes {
		id := program().Nodes[i].ID
		if s.layer != want[id] {
			t.Errorf("%s on layer %d, want %d", id, s.layer, want[id])
		}
	}
	if len(l.edges) != 9 {
		t.Errorf("%d edges, want 9: the repeated PERFORM drawn once", len(l.edges))
	}
}

func TestLayoutKeepsBoxesApartAndEdgesDescending(t *testing.T) {
	l, err := newLayout(program())
	if err != nil {
		t.Fatal(err)
	}
	for i, layer := range l.layers {
		for j := 1; j < len(layer); j++ {
			if layer[j].x < layer[j-1].x+layer[j-1].width+nodeGap {
				t.Errorf("layer %d: slot %d overlaps the one before it", i, j)
			}
		}
	}
	for _, e := range l.edges {
		for i := 1; i < len(e.chain); i++ {
			if e.chain[i].layer != e.chain[i-1].layer+1 {
				t.Errorf("edge chain steps from layer %d to %d", e.chain[i-1].layer, e.chain[i].layer)
			}
		}
	}
	if l.width <= 0 || l.height != 2*margin+5*boxHeight+4*layerGap {
		t.Errorf("size %dx%d", l.width, l.height)
	}
}

func TestLayoutTurnsCyclesAround(t *testing.T) {
	g := Graph{
		Nodes: []Node{{ID: "A", Label: "A"}, {ID: "B", Label: "B"}},
		Edges: []Edge{{From: "A", To: "B"}, {From: "B", To: "A", Style: Dashed}},
	}
	l, err := newLayout(g)
	if err != nil {
		t.Fatal(err)
	}
	if l.nodes[0].layer != 0 || l.nodes[1].layer != 1 || !l.edges[1].reversed {
		t.Errorf("A on %d, B on %d, back edge reversed %v", l.nodes[0].layer, l.nodes[1].layer, l.edges[1].reversed)
	}
}

func TestSVGIsDeterministic(t *testing.T) {
	first, err := program().SVG()
	if err != nil {
		t.Fatal(err)
	}
	second, _ := program().SVG()
	if !bytes.Equal(first, second) {
		t.Error("two drawings of the same graph differ")
	}
	for _, want := range []string{`<title>does open</title>`, `<a href="#OPEN">`, `class="box emphasis"`, `marker-end="url(#arrow)"`} {
		if !bytes.Contains(first, []byte(want)) {
			t.Errorf("SVG lacks %s", want)
		}
	}
}

func TestSVGErrors(t *testing.T) {
	tests := map[string]Graph{
		"duplicate node": {Nodes: []Node{{ID: "A"}, {ID: "A"}}},
		"missing node":   {Nodes: []Node{{ID: "A"}}, Edges: []Edge{{From: "A", To: "B"}}},
	}
	for name, g := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := g.SVG(); err == nil {
				t.Error("drew without error")
			}
		})
	}
}
