package diagram

import (
	"fmt"
	"strings"
)

const style = `<style>
.box { fill: #ffffff; stroke: #57606a; stroke-width: 1 }
.box.emphasis { stroke-width: 2.5 }
.box.data { fill: #e8eef6; stroke: #8c959f }
.box.external { fill: #fbecd0; stroke: #a16207; stroke-width: 2 }
.box.dashed { stroke-dasharray: 4 3 }
.label { fill: #1f2328 }
.label.first { font-weight: 600 }
.edge { fill: none; stroke: #8c959f; stroke-width: 1.2 }
.edge.dashed { stroke-dasharray: 4 3 }
.edge.dotted { stroke-dasharray: 1 3 }
.head { fill: #8c959f }
.legend { fill: #f6f8fa; stroke: #d0d7de }
.edge-label { fill: #57606a; font-size: 11px; paint-order: stroke; stroke: #ffffff; stroke-width: 4px; stroke-linejoin: round }
.box.highlight { stroke: #c2410c; stroke-width: 2.5 }
.edge.highlight { stroke: #c2410c; stroke-width: 2 }
.head.highlight { fill: #c2410c }
@media (prefers-color-scheme: dark) {
  .box { fill: #161b22; stroke: #8b949e }
  .box.data { fill: #1f2a3a; stroke: #6e7681 }
  .box.external { fill: #3a2c12; stroke: #d6a64a }
  .label { fill: #e6edf3 }
  .edge { stroke: #6e7681 }
  .head { fill: #6e7681 }
  .legend { fill: #0d1117; stroke: #30363d }
  .edge-label { fill: #9da5ae; stroke: #0d1117 }
  .box.highlight, .edge.highlight { stroke: #fb923c }
  .head.highlight { fill: #fb923c }
}
</style>
`

// SVG draws the graph.
func (g Graph) SVG() ([]byte, error) {
	l, err := newLayout(g)
	if err != nil {
		return nil, err
	}
	legendX := l.width
	if len(g.Legend) > 0 {
		width, height := legendSize(g.Legend)
		l.width += width + margin
		l.height = max(l.height, height+2*margin)
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" `+
		`font-family="ui-monospace, SFMono-Regular, Menlo, Consolas, monospace" font-size="12">`+"\n",
		l.width, l.height, l.width, l.height)
	b.WriteString(style)
	b.WriteString("<defs>")
	for _, m := range []struct{ id, class string }{{"arrow", "head"}, {"arrow-highlight", "head highlight"}} {
		fmt.Fprintf(&b, `<marker id="%s" viewBox="0 0 %d 8" refX="0" refY="4" markerUnits="userSpaceOnUse" `+
			`markerWidth="%d" markerHeight="8" orient="auto-start-reverse"><path class="%s" d="M0,0 L%d,4 L0,8 z"/></marker>`,
			m.id, arrowLength, arrowLength, m.class, arrowLength)
	}
	b.WriteString("</defs>\n")
	for _, e := range l.edges {
		l.writeEdge(&b, e)
	}
	for _, s := range l.nodes {
		l.writeNode(&b, s, g.Nodes[s.node])
	}
	for _, e := range l.edges {
		l.writeLabel(&b, e)
	}
	if len(g.Legend) > 0 {
		writeLegend(&b, g.Legend, legendX, margin)
	}
	b.WriteString("</svg>\n")
	return []byte(b.String()), nil
}

// arrowLength is how far an arrowhead reaches from the end of its line to
// the box it points at.
const arrowLength = 8

// writeEdge draws an edge as curves between layers, straight down through
// its waypoints. The line stops short of the box at its arrow end, where
// the arrowhead's base takes over, and arrives there vertically.
func (l *layout) writeEdge(b *strings.Builder, e routedEdge) {
	from, last := e.chain[0], len(e.chain)-1
	x, y := port(from, from.below, e.chain[1], e.topLane), l.bottom(from)
	if e.reversed {
		y += arrowLength
	}
	path := fmt.Sprintf("M%d,%d", x, y)
	for i := 1; i <= last; i++ {
		s := e.chain[i]
		lane := 0
		if i == last {
			lane = e.lowerLane
		}
		nx, ny := port(s, s.above, e.chain[i-1], lane), l.top(s)
		if i == last && !e.reversed {
			ny -= arrowLength
		}
		mid := (y + ny) / 2
		path += fmt.Sprintf(" C%d,%d %d,%d %d,%d", x, mid, nx, mid, nx, ny)
		x, y = nx, ny
		if s.node < 0 {
			y = l.bottom(s)
			path += fmt.Sprintf(" L%d,%d", x, y)
		}
	}
	class, marker, arrow := "edge", "marker-end", "arrow"
	if e.style != Solid {
		class += " " + string(e.style)
	}
	if e.highlight {
		class, arrow = class+" highlight", "arrow-highlight"
	}
	if e.reversed {
		marker = "marker-start"
	}
	fmt.Fprintf(b, `<path class="%s" d="%s" %s="url(#%s)"/>`+"\n", class, path, marker, arrow)
}

// nodeClass is the CSS class of a node's box, which draws it as it is.
func nodeClass(n Node) string {
	class := "box"
	for _, c := range []struct {
		set  bool
		name string
	}{{n.Emphasis, "emphasis"}, {n.Data, "data"}, {n.External, "external"}, {n.Dashed, "dashed"}, {n.Highlight, "highlight"}} {
		if c.set {
			class += " " + c.name
		}
	}
	return class
}

// writeLabel writes an edge's label on its middle waypoint.
func (l *layout) writeLabel(b *strings.Builder, e routedEdge) {
	if len(e.label) == 0 {
		return
	}
	waypoint := e.chain[len(e.chain)/2]
	x, y := waypoint.center(), l.layerTop[waypoint.layer]+l.layerHeight[waypoint.layer]/2
	top := y - (len(e.label)-1)*lineHeight/2 + 4
	for i, line := range e.label {
		fmt.Fprintf(b, `<text class="edge-label" x="%d" y="%d" text-anchor="middle">%s</text>`+"\n", x, top+i*lineHeight, escape(line))
	}
}

func (l *layout) writeNode(b *strings.Builder, s *slot, n Node) {
	class := nodeClass(n)
	open, end := "<g>", "</g>"
	if n.Link != "" {
		open, end = fmt.Sprintf(`<a href="%s">`, escape(n.Link)), "</a>"
	}
	b.WriteString(open)
	if n.Title != "" {
		fmt.Fprintf(b, "<title>%s</title>", escape(n.Title))
	}
	top, radius := l.top(s), 4
	if n.Data {
		radius = boxHeight / 2
	}
	fmt.Fprintf(b, `<rect class="%s" x="%d" y="%d" width="%d" height="%d" rx="%d"/>`, class, s.x, top, s.width, s.height, radius)
	lines := strings.Split(n.Label, "\n")
	for i, line := range lines {
		class := "label"
		if i == 0 && len(lines) > 1 {
			class += " first"
		}
		fmt.Fprintf(b, `<text class="%s" x="%d" y="%d" text-anchor="middle">%s</text>`, class, s.center(), top+boxHeight/2+4+i*lineHeight, escape(line))
	}
	b.WriteString(end + "\n")
}
