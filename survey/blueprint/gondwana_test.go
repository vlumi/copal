package blueprint

import (
	"path/filepath"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/internal/golden"
)

// TestGondwana generates the blueprints for Gondwana, read in Fujitsu
// COBOL from Shift_JIS, and compares them with the golden ones in
// blueprints/gondwana. Run with -update to rewrite them.
func TestGondwana(t *testing.T) {
	t.Chdir("../..")
	programs, err := filepath.Glob("gondwana/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	batchFiles, err := filepath.Glob("gondwana/scripts/*.BAT")
	if err != nil {
		t.Fatal(err)
	}
	system, err := Analyze(append(programs, batchFiles...), "gondwana/copybooks", cobol.Fujitsu)
	if err != nil {
		t.Fatal(err)
	}
	if system.Annotations, err = ReadAnnotations("annotations/gondwana"); err != nil {
		t.Fatal(err)
	}
	files, err := system.Render("blueprints/gondwana", "gondwana")
	if err != nil {
		t.Fatal(err)
	}
	golden.CompareFiles(t, files, "blueprints/gondwana")
}
