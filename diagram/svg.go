package diagram

import (
	"fmt"
	"strings"
)

const style = `<style>
.box { fill: #ffffff; stroke: #57606a; stroke-width: 1 }
.box.emphasis { stroke-width: 2.5 }
.box.data { fill: #e8eef6; stroke: #8c959f }
.box.dashed { stroke-dasharray: 4 3 }
.label { fill: #1f2328 }
.label.first { font-weight: 600 }
.edge { fill: none; stroke: #8c959f; stroke-width: 1.2 }
.edge.dashed { stroke-dasharray: 4 3 }
.edge.dotted { stroke-dasharray: 1 3 }
.head { fill: #8c959f }
.legend { fill: #f6f8fa; stroke: #d0d7de }
@media (prefers-color-scheme: dark) {
  .box { fill: #161b22; stroke: #8b949e }
  .box.data { fill: #1f2a3a; stroke: #6e7681 }
  .label { fill: #e6edf3 }
  .edge { stroke: #6e7681 }
  .head { fill: #6e7681 }
  .legend { fill: #0d1117; stroke: #30363d }
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
	fmt.Fprintf(&b, `<defs><marker id="arrow" viewBox="0 0 %d 8" refX="0" refY="4" markerUnits="userSpaceOnUse" `+
		`markerWidth="%d" markerHeight="8" orient="auto-start-reverse"><path class="head" d="M0,0 L%d,4 L0,8 z"/></marker></defs>`+"\n",
		arrowLength, arrowLength, arrowLength)
	for _, e := range l.edges {
		l.writeEdge(&b, e)
	}
	for _, s := range l.nodes {
		l.writeNode(&b, s, g.Nodes[s.node])
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
	x, y := port(from, from.below, e.chain[1]), l.bottom(from)
	if e.reversed {
		y += arrowLength
	}
	path := fmt.Sprintf("M%d,%d", x, y)
	for i := 1; i <= last; i++ {
		s := e.chain[i]
		nx, ny := port(s, s.above, e.chain[i-1]), l.top(s)
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
	class, marker := "edge", "marker-end"
	if e.style != Solid {
		class += " " + string(e.style)
	}
	if e.reversed {
		marker = "marker-start"
	}
	fmt.Fprintf(b, `<path class="%s" d="%s" %s="url(#arrow)"/>`+"\n", class, path, marker)
}

func (l *layout) writeNode(b *strings.Builder, s *slot, n Node) {
	class := "box"
	for _, c := range []struct {
		set  bool
		name string
	}{{n.Emphasis, "emphasis"}, {n.Data, "data"}, {n.Dashed, "dashed"}} {
		if c.set {
			class += " " + c.name
		}
	}
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
