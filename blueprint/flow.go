package blueprint

import "github.com/vlumi/copal/diagram"

// flowDiagram is the graph of a program's paragraphs and the PERFORMs
// between them, with the entry paragraph emphasized. Each paragraph links
// to its source and shows its comment on hover.
func flowDiagram(dir string, p Program) diagram.Graph {
	var g diagram.Graph
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
	}
	return g
}
