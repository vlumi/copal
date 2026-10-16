package blueprint

import (
	"path/filepath"
	"testing"

	"github.com/vlumi/copal/internal/golden"
)

// TestPangaea generates the blueprints for Pangaea and compares them with
// the golden ones in blueprints/pangaea. Run with -update to rewrite them.
func TestPangaea(t *testing.T) {
	t.Chdir("..")
	programs, err := filepath.Glob("pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := Analyze(programs, "pangaea/copybooks")
	if err != nil {
		t.Fatal(err)
	}
	files, err := system.Render("blueprints/pangaea", "pangaea")
	if err != nil {
		t.Fatal(err)
	}
	golden.CompareFiles(t, files, "blueprints/pangaea")
}
