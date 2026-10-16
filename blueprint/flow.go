package blueprint

import "github.com/vlumi/copal/diagram"

// flowDiagram is the graph of a program's paragraphs and how control
// moves between them: PERFORMs solid, GO TOs dashed, and falling into the
// next paragraph dotted. The entry paragraph is emphasized. Each paragraph
// links to its source and shows its comment on hover.
func flowDiagram(dir string, p Program) diagram.Graph {
	g := diagram.Graph{Legend: []diagram.Key{
		diagram.NodeKey(diagram.Node{Emphasis: true}, "first paragraph"),
		diagram.EdgeKey(diagram.Solid, "PERFORM"),
		diagram.EdgeKey(diagram.Dashed, "GO TO"),
		diagram.EdgeKey(diagram.Dotted, "falls into the next"),
	}}
	for i, para := range p.Paragraphs {
		g.Nodes = append(g.Nodes, diagram.Node{
			ID:       para.Name,
			Label:    para.Name,
			Title:    para.Description,
			Link:     sourceTarget(dir, para.Source),
			Emphasis: i == 0,
		})
		for _, perform := range para.Performs {
			g.Edges = append(g.Edges, diagram.Edge{From: para.Name, To: perform.Target})
		}
		for _, goTo := range para.GoTos {
			g.Edges = append(g.Edges, diagram.Edge{From: para.Name, To: goTo.Target, Style: diagram.Dashed})
		}
		if para.FallsInto != "" {
			g.Edges = append(g.Edges, diagram.Edge{From: para.Name, To: para.FallsInto, Style: diagram.Dotted})
		}
	}
	return g
}
