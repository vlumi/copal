package diagram

import (
	"fmt"
	"strings"
)

const style = `<style>
.box { fill: #ffffff; stroke: #57606a; stroke-width: 1 }
.box.emphasis { stroke-width: 2.5 }
.label { fill: #1f2328 }
.edge { fill: none; stroke: #8c959f; stroke-width: 1.2 }
.edge.dashed { stroke-dasharray: 4 3 }
.head { fill: #8c959f }
@media (prefers-color-scheme: dark) {
  .box { fill: #161b22; stroke: #8b949e }
  .label { fill: #e6edf3 }
  .edge { stroke: #6e7681 }
  .head { fill: #6e7681 }
}
</style>
`

// SVG draws the graph.
func (g Graph) SVG() ([]byte, error) {
	l, err := newLayout(g)
	if err != nil {
		return nil, err
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
	x, y := port(from, from.below, e.chain[1]), l.top(from)+boxHeight
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
			y += boxHeight
			path += fmt.Sprintf(" L%d,%d", x, y)
		}
	}
	class, marker := "edge", "marker-end"
	if e.dashed {
		class += " dashed"
	}
	if e.reversed {
		marker = "marker-start"
	}
	fmt.Fprintf(b, `<path class="%s" d="%s" %s="url(#arrow)"/>`+"\n", class, path, marker)
}

func (l *layout) writeNode(b *strings.Builder, s *slot, n Node) {
	class := "box"
	if n.Emphasis {
		class += " emphasis"
	}
	open, end := "<g>", "</g>"
	if n.Link != "" {
		open, end = fmt.Sprintf(`<a href="%s">`, label(n.Link)), "</a>"
	}
	b.WriteString(open)
	if n.Title != "" {
		fmt.Fprintf(b, "<title>%s</title>", label(n.Title))
	}
	top := l.top(s)
	fmt.Fprintf(b, `<rect class="%s" x="%d" y="%d" width="%d" height="%d" rx="4"/>`, class, s.x, top, s.width, boxHeight)
	fmt.Fprintf(b, `<text class="label" x="%d" y="%d" text-anchor="middle">%s</text>`, s.center(), top+boxHeight/2+4, label(n.Label))
	b.WriteString(end + "\n")
}
