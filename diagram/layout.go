// Package diagram draws directed graphs as SVG, laid out in layers from the
// top down. The layout is deterministic: the same graph always gives the
// same bytes, so diagrams can be golden files.
package diagram

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

// Graph is the nodes and edges to draw. Their order matters: within a
// layer, nodes are ordered as a depth-first walk first reaches them,
// starting from the nodes in order and following the edges in order. For a
// call graph, with the edges in statement order, that is execution order.
type Graph struct {
	Nodes  []Node
	Edges  []Edge
	Legend []Key // drawn in a box of its own, to the right at the top
}

// Node is a box with a label, of one line or several separated by \n; the
// first of several is set in bold. Title is shown on hover and Link is followed
// on click, where the viewer supports it. An emphasized node is drawn with a
// heavier border; a data node with round ends and a tinted fill, and a
// dashed one with a dashed border.
type Node struct {
	ID       string
	Label    string
	Title    string
	Link     string
	Emphasis bool
	Data     bool
	Dashed   bool
}

// Edge is an arrow from one node to another, drawn in a style. Span is the
// fewest layers it goes down, 1 when not set: an edge that spans 2 keeps a
// layer free between its nodes, for others to use.
type Edge struct {
	From, To string
	Style    Style
	Span     int
}

// Style is how an edge's line is drawn.
type Style string

const (
	Solid  Style = ""
	Dashed Style = "dashed"
	Dotted Style = "dotted"
)

// Sizes in pixels. Labels are set in a monospace font, so their width
// follows from their length.
const (
	charWidth  = 7
	boxPadding = 10
	boxHeight  = 26 // of a box with a one-line label
	lineHeight = 14
	layerGap   = 46
	nodeGap    = 18
	dummyWidth = 8
	margin     = 12
)

// slot is a box or, for an edge that spans several layers, a waypoint the
// edge passes through on an intermediate layer.
type slot struct {
	node   int // index in Graph.Nodes, or -1 for a waypoint
	layer  int
	order  float64
	x      int // left edge
	width  int
	height int     // a box's; a waypoint takes its layer's
	above  []*slot // slots in the layer above with an edge to this one
	below  []*slot
}

func (s *slot) center() int { return s.x + s.width/2 }

// layout is the positioned graph. A layer is as tall as its tallest box,
// and the boxes in it are centered on it.
type layout struct {
	layers      [][]*slot
	layerTop    []int
	layerHeight []int
	nodes       []*slot
	edges       []routedEdge
	width       int
	height      int
}

// routedEdge is an edge as a chain of slots, from top to bottom through any
// waypoints. The arrow points at the bottom end unless the edge is reversed.
type routedEdge struct {
	chain    []*slot
	style    Style
	reversed bool
}

func newLayout(g Graph) (*layout, error) {
	index := map[string]int{}
	for i, n := range g.Nodes {
		if _, dup := index[n.ID]; dup {
			return nil, fmt.Errorf("diagram: node %s twice", n.ID)
		}
		index[n.ID] = i
	}
	edges, err := forwardEdges(g, index)
	if err != nil {
		return nil, err
	}
	l := &layout{}
	l.assignLayers(g, edges, walkOrder(len(g.Nodes), edges))
	l.route(edges)
	l.order()
	l.place()
	return l, nil
}

type indexedEdge struct {
	from, to int
	style    Style
	span     int
	reversed bool // turned around to break a cycle; the arrow points at from
}

// forwardEdges returns the graph's edges by node index, each once, without
// self-loops, and turned around where they would close a cycle, so that
// the layers can follow them downward.
func forwardEdges(g Graph, index map[string]int) ([]indexedEdge, error) {
	var edges []indexedEdge
	seen := map[[2]int]bool{}
	for _, e := range g.Edges {
		from, okFrom := index[e.From]
		to, okTo := index[e.To]
		if !okFrom || !okTo {
			return nil, fmt.Errorf("diagram: edge %s to %s names a missing node", e.From, e.To)
		}
		if from == to || seen[[2]int{from, to}] {
			continue
		}
		seen[[2]int{from, to}] = true
		edges = append(edges, indexedEdge{from: from, to: to, style: e.Style, span: max(e.Span, 1)})
	}
	back := backEdges(len(g.Nodes), edges)
	for i := range edges {
		if back[i] {
			edges[i].from, edges[i].to = edges[i].to, edges[i].from
			edges[i].reversed = true
		}
	}
	return edges, nil
}

// backEdges marks the edges that point back to a node still being visited
// in a depth-first walk from the nodes in order.
func backEdges(n int, edges []indexedEdge) map[int]bool {
	const (
		unvisited = iota
		visiting
		done
	)
	state := make([]int, n)
	back := map[int]bool{}
	var visit func(v int)
	visit = func(v int) {
		state[v] = visiting
		for i, e := range edges {
			if e.from != v {
				continue
			}
			switch state[e.to] {
			case visiting:
				back[i] = true
			case unvisited:
				visit(e.to)
			}
		}
		state[v] = done
	}
	for v := range n {
		if state[v] == unvisited {
			visit(v)
		}
	}
	return back
}

// walkOrder numbers the nodes in the order a depth-first walk reaches them.
func walkOrder(n int, edges []indexedEdge) []int {
	order := make([]int, n)
	seen := make([]bool, n)
	next := 0
	var visit func(v int)
	visit = func(v int) {
		seen[v] = true
		order[v] = next
		next++
		for _, e := range edges {
			if e.from == v && !seen[e.to] {
				visit(e.to)
			}
		}
	}
	for v := range n {
		if !seen[v] {
			visit(v)
		}
	}
	return order
}

// assignLayers puts each node below the lowest node with an edge to it, by
// the edge's span: the length of the longest path to it. A node no edge
// points to goes as low as the spans of its edges allow, next to what it
// feeds.
func (l *layout) assignLayers(g Graph, edges []indexedEdge, order []int) {
	layer := make([]int, len(g.Nodes))
	for changed := true; changed; {
		changed = false
		for _, e := range edges {
			if layer[e.to] < layer[e.from]+e.span {
				layer[e.to] = layer[e.from] + e.span
				changed = true
			}
		}
	}
	pointedTo := make([]bool, len(g.Nodes))
	for _, e := range edges {
		pointedTo[e.to] = true
	}
	for v := range g.Nodes {
		if pointedTo[v] {
			continue
		}
		lowest := -1
		for _, e := range edges {
			if e.from == v && (lowest < 0 || layer[e.to]-e.span < lowest) {
				lowest = layer[e.to] - e.span
			}
		}
		layer[v] = max(layer[v], lowest)
	}
	for i, n := range g.Nodes {
		lines := strings.Split(n.Label, "\n")
		longest := 0
		for _, line := range lines {
			longest = max(longest, len(line))
		}
		s := &slot{node: i, layer: layer[i], order: float64(order[i]), width: longest*charWidth + 2*boxPadding,
			height: boxHeight + (len(lines)-1)*lineHeight}
		l.nodes = append(l.nodes, s)
		l.add(s)
	}
}

func (l *layout) add(s *slot) {
	for len(l.layers) <= s.layer {
		l.layers = append(l.layers, nil)
	}
	l.layers[s.layer] = append(l.layers[s.layer], s)
}

// route links the slots of each edge, adding a waypoint on every layer an
// edge passes through.
func (l *layout) route(edges []indexedEdge) {
	for _, e := range edges {
		from, to := l.nodes[e.from], l.nodes[e.to]
		chain := []*slot{from}
		for layer := from.layer + 1; layer < to.layer; layer++ {
			waypoint := &slot{node: -1, layer: layer, order: from.order, width: dummyWidth}
			l.add(waypoint)
			chain = append(chain, waypoint)
		}
		chain = append(chain, to)
		for i := 1; i < len(chain); i++ {
			chain[i-1].below = append(chain[i-1].below, chain[i])
			chain[i].above = append(chain[i].above, chain[i-1])
		}
		l.edges = append(l.edges, routedEdge{chain: chain, style: e.style, reversed: e.reversed})
	}
}

// order sorts each layer, from the top down, by the mean position of each
// slot's neighbors above, which keeps the walk order among slots with the
// same callers and reduces crossings between the rest.
//
// Then, from the bottom up, each slot with nothing above it moves to the
// mean position of the slots it points to, so that what feeds a node sits
// next to it instead of where the walk happened to leave it.
func (l *layout) order() {
	l.renumber()
	for i := 1; i < len(l.layers); i++ {
		sortByNeighbors(l.layers[i], func(s *slot) []*slot { return s.above })
	}
	for i := len(l.layers) - 2; i >= 0; i-- {
		sortByNeighbors(l.layers[i], func(s *slot) []*slot {
			if len(s.above) > 0 {
				return nil
			}
			return s.below
		})
	}
}

func (l *layout) renumber() {
	for _, layer := range l.layers {
		slices.SortStableFunc(layer, func(a, b *slot) int { return compare(a.order, b.order) })
		for i, s := range layer {
			s.order = float64(i)
		}
	}
}

// sortByNeighbors reorders a layer by the mean order of each slot's
// neighbors; a slot without neighbors keeps its place.
func sortByNeighbors(layer []*slot, neighbors func(*slot) []*slot) {
	key := map[*slot]float64{}
	for _, s := range layer {
		key[s] = s.order
		if ns := neighbors(s); len(ns) > 0 {
			sum := 0.0
			for _, n := range ns {
				sum += n.order
			}
			key[s] = sum / float64(len(ns))
		}
	}
	slices.SortStableFunc(layer, func(a, b *slot) int {
		if c := compare(key[a], key[b]); c != 0 {
			return c
		}
		return compare(a.order, b.order)
	})
	for i, s := range layer {
		s.order = float64(i)
	}
}

func compare(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// place sets x coordinates. Each layer starts side by side; then, from the
// top down, every slot moves as close as it can to centered under its
// neighbors above, and from the bottom up, every slot over its neighbors
// below, keeping the layer's order and gaps. Then everything shifts to the
// margin.
func (l *layout) place() {
	for _, layer := range l.layers {
		next := 0
		for _, s := range layer {
			s.x = next
			next += s.width + nodeGap
		}
	}
	above := func(s *slot) []*slot { return s.above }
	below := func(s *slot) []*slot { return s.below }
	for i := 1; i < len(l.layers); i++ {
		center(l.layers[i], above)
	}
	for i := len(l.layers) - 2; i >= 0; i-- {
		center(l.layers[i], below)
	}
	left := math.MaxInt
	for _, layer := range l.layers {
		if len(layer) > 0 {
			left = min(left, layer[0].x)
		}
	}
	for _, layer := range l.layers {
		for _, s := range layer {
			s.x += margin - left
			l.width = max(l.width, s.x+s.width+margin)
		}
	}
	top := margin
	for _, layer := range l.layers {
		height := boxHeight
		for _, s := range layer {
			height = max(height, s.height)
		}
		l.layerTop = append(l.layerTop, top)
		l.layerHeight = append(l.layerHeight, height)
		top += height + layerGap
	}
	l.height = top - layerGap + margin
}

// center moves each slot of a layer toward the mean center of its
// neighbors, never closer to the next slot than the gap allows. It places
// the layer twice, pushing crowded slots rightward in one pass and leftward
// in the other, and takes the mean of the two, so that crowding spreads
// both ways instead of drifting to one side. A slot without neighbors stays
// as close to its place as it can.
func center(layer []*slot, neighbors func(*slot) []*slot) {
	want := make([]int, len(layer))
	for i, s := range layer {
		want[i] = s.x
		if ns := neighbors(s); len(ns) > 0 {
			sum := 0
			for _, n := range ns {
				sum += n.center()
			}
			want[i] = sum/len(ns) - s.width/2
		}
	}
	right := make([]int, len(layer))
	for i := range layer {
		right[i] = want[i]
		if i > 0 {
			right[i] = max(right[i], right[i-1]+layer[i-1].width+nodeGap)
		}
	}
	left := make([]int, len(layer))
	for i := len(layer) - 1; i >= 0; i-- {
		left[i] = want[i]
		if i < len(layer)-1 {
			left[i] = min(left[i], left[i+1]-layer[i].width-nodeGap)
		}
	}
	for i, s := range layer {
		s.x = (left[i] + right[i]) / 2
	}
}

// top and bottom are where a box's top and bottom sides are, or where a
// waypoint's line enters and leaves its layer.
func (l *layout) top(s *slot) int {
	if s.node < 0 {
		return l.layerTop[s.layer]
	}
	return l.layerTop[s.layer] + (l.layerHeight[s.layer]-s.height)/2
}

func (l *layout) bottom(s *slot) int {
	if s.node < 0 {
		return l.layerTop[s.layer] + l.layerHeight[s.layer]
	}
	return l.top(s) + s.height
}

// port spreads the edges that leave or enter a box along its bottom or top
// side, in the order of the boxes at their other end.
func port(s *slot, others []*slot, other *slot) int {
	if s.node < 0 || len(others) < 2 {
		return s.center()
	}
	sorted := slices.Clone(others)
	slices.SortStableFunc(sorted, func(a, b *slot) int { return a.center() - b.center() })
	k := slices.Index(sorted, other)
	return s.x + s.width*(k+1)/(len(sorted)+1)
}

// escape makes text safe in XML, in an element or an attribute.
func escape(s string) string { return xmlEscaper.Replace(s) }

var xmlEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
