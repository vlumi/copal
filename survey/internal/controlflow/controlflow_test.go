package controlflow

import (
	"slices"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
)

func graph(t *testing.T, path string) *Graph {
	t.Helper()
	p, err := cobol.ReadProgram(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	return New(p)
}

// at lists the nodes of the statements written as text that control can
// reach, one for each time the graph lays them out.
func at(g *Graph, text string) []int {
	reached := g.Reachable()
	var nodes []int
	for i, n := range g.Nodes {
		if reached[i] && n.Statement != nil && n.Statement.Text == text {
			nodes = append(nodes, i)
		}
	}
	return nodes
}

// reaches reports whether control can go from any of the nodes to any of
// the others, by at least one step.
func reaches(g *Graph, from, to []int) bool {
	seen := map[int]bool{}
	var queue []int
	for _, n := range from {
		queue = append(queue, g.Nodes[n].Next...)
	}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if seen[n] {
			continue
		}
		seen[n] = true
		for _, m := range to {
			if n == m {
				return true
			}
		}
		queue = append(queue, g.Nodes[n].Next...)
	}
	return false
}

func TestPaths(t *testing.T) {
	g := graph(t, "testdata/PATHS.cbl")
	if !g.Covered {
		t.Fatal("graph does not cover the program")
	}
	if len(at(g, "GO TO MAIN-PARA")) > 0 {
		t.Error("a paragraph nothing performs or falls into is reached")
	}
	if work := at(g, "MOVE WS-A TO WS-B"); len(work) != 2 {
		t.Fatalf("WORK-PARA laid out %d times for its two PERFORMs, want 2", len(work))
	}
	for _, move := range []string{"MOVE 1 TO WS-A", "MOVE 2 TO WS-A"} {
		if !reaches(g, at(g, move), at(g, "MOVE WS-A TO WS-B")) {
			t.Errorf("WORK-PARA is not performed after %s", move)
		}
		if reaches(g, at(g, move), at(g, move)) {
			t.Errorf("control returns from a PERFORM of WORK-PARA to before %s", move)
		}
	}
	for _, n := range at(g, "EXIT PARAGRAPH") {
		if slices.ContainsFunc(g.Nodes[n].Next, func(m int) bool { return slices.Contains(at(g, "MOVE WS-A TO WS-B"), m) }) {
			t.Error("EXIT PARAGRAPH runs on to the paragraph's next statement")
		}
	}
	if !reaches(g, at(g, "EXIT PARAGRAPH"), at(g, "MOVE 2 TO WS-A")) {
		t.Error("EXIT PARAGRAPH does not return to the statement after the PERFORM")
	}
	if reaches(g, at(g, "GO TO CHECK-EXIT"), at(g, "MOVE 5 TO WS-B")) {
		t.Error("the GO TO runs on into the statements it skips")
	}
	if !reaches(g, at(g, "MOVE 4 TO WS-B"), at(g, "MOVE 5 TO WS-B")) {
		t.Error("the ELSE does not run on to the next statement")
	}
	add := at(g, "ADD 1 TO WS-A")
	if !reaches(g, add, add) {
		t.Error("the inline PERFORM does not loop")
	}
	if !reaches(g, at(g, "PERFORM UNTIL WS-A = 9"), at(g, "STOP RUN")) {
		t.Error("the loop cannot run no times")
	}
}

func TestGoToDependingOn(t *testing.T) {
	g := graph(t, "../structure/testdata/GOTO.cbl")
	if !g.Covered {
		t.Fatal("graph does not cover the program")
	}
	if !reaches(g, at(g, "GO TO CHECK-ONE CHECK-TWO DEPENDING ON WS-I"), at(g, "MOVE 1 TO WS-I")) {
		t.Error("GO TO DEPENDING ON does not go to its first target")
	}
	if !reaches(g, at(g, "GO TO CHECK-ONE CHECK-TWO DEPENDING ON WS-I"), at(g, "IF WS-I = 2")) {
		t.Error("GO TO DEPENDING ON does not go to its second target")
	}
}

func TestGoToOutOfRange(t *testing.T) {
	if g := graph(t, "testdata/ESCAPE.cbl"); g.Covered {
		t.Error("a GO TO out of a performed range is covered")
	}
}
