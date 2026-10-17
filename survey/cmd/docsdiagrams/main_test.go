package main

import (
	"testing"

	"github.com/vlumi/copal/internal/golden"
)

// The diagrams in docs/diagrams are what the command draws. Run with
// -update to redraw them.
func TestDiagramsAreCurrent(t *testing.T) {
	dir := t.TempDir()
	if err := Write(dir); err != nil {
		t.Fatal(err)
	}
	golden.CompareDir(t, dir, "../../../docs/diagrams")
}
