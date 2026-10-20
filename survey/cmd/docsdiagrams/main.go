// Command docsdiagrams draws the diagrams of docs/architecture.md with the
// survey's own diagram package, so that they are laid out like the
// blueprints' and stay in step with the code: the parts of Copal and how
// they depend on each other, and the survey's pipeline from source to
// blueprints.
//
//	docsdiagrams -out dir
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vlumi/copal/survey/internal/diagram"
)

func main() {
	out := flag.String("out", "docs/diagrams", "directory to write the diagrams to")
	flag.Parse()
	if err := Write(*out); err != nil {
		fmt.Fprintln(os.Stderr, "docsdiagrams:", err)
		os.Exit(1)
	}
}

// Write draws every diagram into dir.
func Write(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for name, g := range map[string]diagram.Graph{"architecture": architecture(), "pipeline": pipeline()} {
		svg, err := g.SVG()
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if err := os.WriteFile(filepath.Join(dir, name+".svg"), svg, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// architecture is the four parts of Copal, the second specimen and the
// blueprints between them.
func architecture() diagram.Graph {
	return diagram.Graph{
		Legend: []diagram.Key{
			diagram.NodeKey(diagram.Node{}, "Go, one part"),
			diagram.NodeKey(diagram.Node{Data: true}, "files"),
			diagram.EdgeKey(diagram.Solid, "imports"),
			diagram.EdgeKey(diagram.Dashed, "writes"),
			diagram.EdgeKey(diagram.Dotted, "reads"),
		},
		Nodes: []diagram.Node{
			{ID: "console", Label: "console\nruns the bank, shows its survey"},
			{ID: "survey", Label: "survey\nthe analyzer"},
			{ID: "pangaea", Label: "pangaea\nthe legacy bank"},
			{ID: "laurasia", Label: "laurasia\nthe card processor"},
			{ID: "blueprints", Label: "blueprints/pangaea\nthe survey of the bank", Data: true},
			{ID: "bedrock", Label: "bedrock\nrunning programs, jobs, scripts"},
			{ID: "cobolith", Label: "cobolith, a library\nCOBOL, JCL, layouts, records"},
		},
		Edges: []diagram.Edge{
			{From: "console", To: "survey"},
			{From: "console", To: "pangaea"},
			{From: "console", To: "bedrock", Span: 3},
			{From: "console", To: "blueprints", Style: diagram.Dotted, Span: 2},
			{From: "survey", To: "bedrock", Span: 2},
			{From: "survey", To: "pangaea", Style: diagram.Dotted},
			{From: "survey", To: "blueprints", Style: diagram.Dashed},
			{From: "pangaea", To: "bedrock"},
			{From: "laurasia", To: "bedrock"},
			{From: "bedrock", To: "cobolith"},
			{From: "console", To: "cobolith", Span: 4},
			{From: "survey", To: "cobolith", Span: 3},
			{From: "pangaea", To: "cobolith", Span: 2},
			{From: "laurasia", To: "cobolith", Span: 2},
		},
	}
}

// pipeline is the survey's analyses, from the source to the blueprints and
// the explorer. An arrow says what feeds what.
func pipeline() diagram.Graph {
	return diagram.Graph{
		Legend: []diagram.Key{
			diagram.NodeKey(diagram.Node{}, "package"),
			diagram.NodeKey(diagram.Node{Data: true}, "files"),
			diagram.EdgeKey(diagram.Solid, "feeds"),
		},
		Nodes: []diagram.Node{
			{ID: "source", Label: "programs, copybooks, jobs\nthe system's source", Data: true},
			{ID: "cobol", Label: "cobolith/cobol\nlexing, COPY, parsing"},
			{ID: "jcl", Label: "cobolith/jcl\njobs, steps, DDs"},
			{ID: "layout", Label: "cobolith/layout\nrecord layouts"},
			{ID: "structure", Label: "structure\nparagraphs, PERFORM, GO TO"},
			{ID: "controlflow", Label: "controlflow\npaths through the statements"},
			{ID: "fileio", Label: "fileio\nOPEN modes, sorts"},
			{ID: "dataflow", Label: "dataflow\nreads, writes, lineage, origins"},
			{ID: "rules", Label: "rules\ndecision tables"},
			{ID: "diagram", Label: "diagram\nlayered SVG"},
			{ID: "blueprint", Label: "blueprint\nprograms, copybooks, jobs, findings"},
			{ID: "lineage", Label: "lineage\nfields followed across programs"},
			{ID: "blueprints", Label: "blueprints\nMarkdown, JSON, SVG", Data: true},
			{ID: "explore", Label: "explore\nthe web app"},
			{ID: "run", Label: "bedrock/run\ncompiling and running"},
			{ID: "coverage", Label: "coverage\nrules a trace reached"},
			{ID: "oracle", Label: "oracle\ncases that reach the rules"},
		},
		Edges: []diagram.Edge{
			{From: "source", To: "cobol"},
			{From: "source", To: "jcl"},
			{From: "cobol", To: "layout"},
			{From: "cobol", To: "structure"},
			{From: "cobol", To: "controlflow"},
			{From: "controlflow", To: "dataflow"},
			{From: "cobol", To: "fileio"},
			{From: "cobol", To: "dataflow"},
			{From: "cobol", To: "rules"},
			{From: "layout", To: "dataflow"},
			{From: "jcl", To: "blueprint", Span: 3},
			{From: "structure", To: "blueprint"},
			{From: "fileio", To: "blueprint"},
			{From: "dataflow", To: "blueprint"},
			{From: "rules", To: "blueprint"},
			{From: "diagram", To: "blueprint"},
			{From: "blueprint", To: "blueprints"},
			{From: "blueprint", To: "lineage"},
			{From: "blueprint", To: "explore", Span: 2},
			{From: "lineage", To: "explore"},
			{From: "rules", To: "coverage"},
			{From: "coverage", To: "oracle"},
			{From: "blueprint", To: "oracle"},
			{From: "run", To: "oracle"},
		},
	}
}
