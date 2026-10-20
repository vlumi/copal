package blueprint

import (
	"path/filepath"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/internal/golden"
)

// TestLaurasia generates the blueprints for Laurasia and compares them
// with the golden ones in blueprints/laurasia. Run with -update to rewrite
// them.
func TestLaurasia(t *testing.T) {
	t.Chdir("../..")
	programs, err := filepath.Glob("laurasia/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	scripts, err := filepath.Glob("laurasia/scripts/*.sh")
	if err != nil {
		t.Fatal(err)
	}
	system, err := Analyze(append(programs, scripts...), "laurasia/copybooks", cobol.MicroFocus)
	if err != nil {
		t.Fatal(err)
	}
	if system.Annotations, err = ReadAnnotations("annotations/laurasia"); err != nil {
		t.Fatal(err)
	}
	files, err := system.Render("blueprints/laurasia", "laurasia")
	if err != nil {
		t.Fatal(err)
	}
	golden.CompareFiles(t, files, "blueprints/laurasia")
}
