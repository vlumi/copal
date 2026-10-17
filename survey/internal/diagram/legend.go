package diagram

import (
	"fmt"
	"strings"
)

// Key is an entry of a graph's legend: what a kind of node or edge stands
// for, next to a sample of it, or a line of text alone.
type Key struct {
	text  string
	node  *Node
	style *Style
}

// NodeKey says what nodes drawn like the sample stand for.
func NodeKey(sample Node, text string) Key { return Key{text: text, node: &sample} }

// EdgeKey says what edges drawn in the style stand for.
func EdgeKey(style Style, text string) Key { return Key{text: text, style: &style} }

// TextKey is a line of the legend without a sample.
func TextKey(text string) Key { return Key{text: text} }

// Legend sizes in pixels.
const (
	legendPadding = 8
	legendRow     = 20
	sampleWidth   = 34
	sampleGap     = 8
)

// legendSize is the width and height of a legend's box.
func legendSize(keys []Key) (int, int) {
	longest := 0
	for _, k := range keys {
		longest = max(longest, len(k.text))
	}
	return 2*legendPadding + sampleWidth + sampleGap + longest*charWidth, 2*legendPadding + len(keys)*legendRow
}

// writeLegend draws the legend's box with its top left corner at x, y.
func writeLegend(b *strings.Builder, keys []Key, x, y int) {
	width, height := legendSize(keys)
	fmt.Fprintf(b, `<g class="legend-keys"><rect class="legend" x="%d" y="%d" width="%d" height="%d" rx="4"/>`+"\n", x, y, width, height)
	for i, k := range keys {
		left := x + legendPadding
		middle := y + legendPadding + i*legendRow + legendRow/2
		switch {
		case k.node != nil:
			radius := 3
			if k.node.Data {
				radius = 7
			}
			fmt.Fprintf(b, `<rect class="%s" x="%d" y="%d" width="%d" height="14" rx="%d"/>`+"\n",
				nodeClass(*k.node), left, middle-7, sampleWidth, radius)
		case k.style != nil:
			class := "edge"
			if *k.style != Solid {
				class += " " + string(*k.style)
			}
			fmt.Fprintf(b, `<path class="%s" d="M%d,%d L%d,%d" marker-end="url(#arrow)"/>`+"\n",
				class, left, middle, left+sampleWidth-arrowLength, middle)
		}
		textX := left + sampleWidth + sampleGap
		if k.node == nil && k.style == nil {
			textX = left
		}
		fmt.Fprintf(b, `<text class="label" x="%d" y="%d">%s</text>`+"\n", textX, middle+4, escape(k.text))
	}
	b.WriteString("</g>\n")
}
