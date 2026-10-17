package blueprint

import "github.com/vlumi/copal/survey/internal/diagram"

// flowDiagram is the program's flow, each paragraph linking to its source.
func flowDiagram(dir string, p Program) diagram.Graph {
	return FlowGraph(p, FlowNotes{Link: func(para Paragraph) string { return sourceTarget(dir, para.Source) }})
}

// FlowNotes say where a program's paragraphs link to, and which are on a
// path being followed, to highlight. Either may be nil.
type FlowNotes struct {
	Link   func(Paragraph) string
	OnPath func(Paragraph) bool
}

// FlowGraph is the graph of a program's paragraphs and how control moves
// between them: PERFORMs solid, GO TOs dashed, and falling into the next
// paragraph dotted. The entry paragraph is emphasized. Each paragraph
// shows its comment on hover, and links and is highlighted as the notes
// say.
func FlowGraph(p Program, notes FlowNotes) diagram.Graph {
	g := diagram.Graph{Legend: []diagram.Key{
		diagram.NodeKey(diagram.Node{Emphasis: true}, "first paragraph"),
		diagram.EdgeKey(diagram.Solid, "PERFORM"),
		diagram.EdgeKey(diagram.Dashed, "GO TO"),
		diagram.EdgeKey(diagram.Dotted, "falls into the next"),
	}}
	if notes.OnPath != nil {
		g.Legend = append(g.Legend, diagram.NodeKey(diagram.Node{Highlight: true}, "moves the values followed"))
	}
	for i, para := range p.Paragraphs {
		n := diagram.Node{ID: para.Name, Label: para.Name, Title: para.Description, Emphasis: i == 0}
		if notes.Link != nil {
			n.Link = notes.Link(para)
		}
		if notes.OnPath != nil {
			n.Highlight = notes.OnPath(para)
		}
		g.Nodes = append(g.Nodes, n)
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
