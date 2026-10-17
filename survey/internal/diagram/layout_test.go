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

// A source sits just above what it feeds, an edge spanning two layers
// keeps one free, and a source goes next to what it feeds, not where the
// walk left it.
func TestLayoutSpansAndSources(t *testing.T) {
	g := Graph{
		Nodes: []Node{{ID: "A", Label: "A"}, {ID: "B", Label: "B"}, {ID: "OUT", Label: "OUT"}, {ID: "IN", Label: "IN"}},
		Edges: []Edge{{From: "A", To: "OUT"}, {From: "A", To: "B", Span: 2}, {From: "IN", To: "B"}},
	}
	l, err := newLayout(g)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int{"A": 0, "OUT": 1, "IN": 1, "B": 2}
	for i, s := range l.nodes {
		if id := g.Nodes[i].ID; s.layer != want[id] {
			t.Errorf("%s on layer %d, want %d", id, s.layer, want[id])
		}
	}
	in, out := l.nodes[3], l.nodes[2]
	if in.x > out.x+out.width {
		t.Errorf("IN at %d, after OUT at %d: not next to B, which it feeds", in.x, out.x)
	}
}

// A layer is as tall as its tallest box, and shorter boxes are centered
// on it.
func TestLayoutMultilineLabels(t *testing.T) {
	g := Graph{
		Nodes: []Node{{ID: "A", Label: "A"}, {ID: "B", Label: "B\nprogram LONGNAME"}, {ID: "C", Label: "C"}},
		Edges: []Edge{{From: "A", To: "B"}, {From: "A", To: "C"}},
	}
	l, err := newLayout(g)
	if err != nil {
		t.Fatal(err)
	}
	b, c := l.nodes[1], l.nodes[2]
	if b.height != boxHeight+lineHeight || l.layerHeight[1] != b.height {
		t.Errorf("B %d high, its layer %d, want %d", b.height, l.layerHeight[1], boxHeight+lineHeight)
	}
	if b.width != len("program LONGNAME")*charWidth+2*boxPadding {
		t.Errorf("B %d wide, want its longest line's width", b.width)
	}
	if l.top(c)-l.top(b) != lineHeight/2 {
		t.Errorf("C at %d, B at %d: C not centered on the layer", l.top(c), l.top(b))
	}
}

// A source whose edge spans two layers keeps the layer between free.
func TestLayoutSourceKeepsSpan(t *testing.T) {
	g := Graph{
		Nodes: []Node{{ID: "ROOT", Label: "ROOT"}, {ID: "MID", Label: "MID"}, {ID: "END", Label: "END"}, {ID: "SRC", Label: "SRC"}},
		Edges: []Edge{{From: "ROOT", To: "MID"}, {From: "MID", To: "END"}, {From: "SRC", To: "END", Span: 2}},
	}
	l, err := newLayout(g)
	if err != nil {
		t.Fatal(err)
	}
	if src, end := l.nodes[3].layer, l.nodes[2].layer; end-src != 2 {
		t.Errorf("SRC on layer %d, END on %d: want 2 layers between", src, end)
	}
}

// A highlighted node and edge are drawn in the accent color, the edge
// with an arrowhead of its own.
func TestHighlight(t *testing.T) {
	g := Graph{
		Nodes: []Node{{ID: "a", Label: "A", Highlight: true}, {ID: "b", Label: "B"}},
		Edges: []Edge{{From: "a", To: "b", Highlight: true}},
	}
	svg, err := g.SVG()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`class="box highlight"`, `class="edge highlight"`, `marker-end="url(#arrow-highlight)"`} {
		if !strings.Contains(string(svg), want) {
			t.Errorf("no %s in\n%s", want, svg)
		}
	}
}

// A legend's sample of a node is drawn as the node is.
func TestLegendSamples(t *testing.T) {
	g := Graph{Nodes: []Node{{ID: "a", Label: "A"}}, Legend: []Key{NodeKey(Node{Emphasis: true}, "first")}}
	svg, err := g.SVG()
	if err != nil {
		t.Fatal(err)
	}
	if want := `<rect class="box emphasis" x="`; !strings.Contains(string(svg), want) {
		t.Errorf("no %s in\n%s", want, svg)
	}
}

// An external node has a look of its own, in the legend too, and a
// highlight still shows on it.
func TestExternal(t *testing.T) {
	g := Graph{
		Nodes:  []Node{{ID: "a", Label: "A", Data: true, External: true}, {ID: "b", Label: "B", External: true, Highlight: true}},
		Legend: []Key{NodeKey(Node{Data: true, External: true}, "from outside")},
	}
	svg, err := g.SVG()
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(svg), `class="box data external"`); n != 2 {
		t.Errorf("%d external data boxes, want the node and its legend sample:\n%s", n, svg)
	}
	style := string(svg[:strings.Index(string(svg), "@media")])
	if strings.Index(style, ".box.external") > strings.Index(style, ".box.highlight") || !strings.Contains(string(svg), `class="box external highlight"`) {
		t.Errorf("a highlight does not show on an external node:\n%s", svg)
	}
}

// Edges both ways between two nodes are drawn side by side, not over
// each other.
func TestEdgesBothWays(t *testing.T) {
	g := Graph{Nodes: []Node{{ID: "a", Label: "OPEN"}, {ID: "b", Label: "FROZEN"}},
		Edges: []Edge{{From: "a", To: "b"}, {From: "b", To: "a"}}}
	l, err := newLayout(g)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.edges) != 2 || l.edges[0].topLane == l.edges[1].topLane || l.edges[0].lowerLane == l.edges[1].lowerLane {
		t.Errorf("edges %+v", l.edges)
	}
	var b strings.Builder
	for _, e := range l.edges {
		l.writeEdge(&b, e)
	}
	d := func(path string) string {
		_, after, _ := strings.Cut(path, ` d="`)
		line, _, _ := strings.Cut(after, `"`)
		return line
	}
	paths := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(paths) != 2 || d(paths[0]) == d(paths[1]) {
		t.Errorf("edges drawn over each other:\n%s", b.String())
	}
}

// An edge's label is written at its midpoint, with room made for it: a
// gap tall enough between the layers it joins, or a waypoint as wide as
// the label where it passes through a layer.
func TestEdgeLabels(t *testing.T) {
	g := Graph{Nodes: []Node{{ID: "a", Label: "A"}, {ID: "b", Label: "B"}, {ID: "c", Label: "C"}},
		Edges: []Edge{
			{From: "a", To: "b", Label: "P 1000-ONE\nQ 2000-TWO\nR 3000-THREE"},
			{From: "b", To: "c"},
			{From: "a", To: "c", Label: "S 4000-FOUR-AROUND"},
		}}
	l, err := newLayout(g)
	if err != nil {
		t.Fatal(err)
	}
	if gap := l.layerTop[1] - l.layerTop[0] - l.layerHeight[0]; gap < 3*lineHeight+2*arrowLength {
		t.Errorf("gap %d below a three-line label", gap)
	}
	waypoint := l.edges[2].chain[1]
	if waypoint.node >= 0 || waypoint.width < len("S 4000-FOUR-AROUND")*labelCharWidth {
		t.Errorf("waypoint %+v too narrow for its label", waypoint)
	}
	svg, err := g.SVG()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{">Q 2000-TWO</text>", ">S 4000-FOUR-AROUND</text>"} {
		if !strings.Contains(string(svg), want) {
			t.Errorf("no %s in\n%s", want, svg)
		}
	}
	again, _ := g.SVG()
	if string(again) != string(svg) {
		t.Error("labels drawn differently twice")
	}
}
