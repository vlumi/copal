package oracle

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/internal/golden"
	"github.com/vlumi/copal/survey/blueprint"
)

// TestConform checks two replacements for BANKDATE against its recorded
// cases: the program rewritten in Go, which conforms, and a script that
// writes the control record back unchanged, which does not; and the
// reports of both.
func TestConform(t *testing.T) {
	rewrite := filepath.Join(t.TempDir(), "bankdate")
	if out, err := exec.Command("go", "build", "-buildvcs=false", "-o", rewrite, "./testdata/rewrites/bankdate").CombinedOutput(); err != nil {
		t.Fatalf("building the rewrite: %v\n%s", err, out)
	}
	wrong, err := filepath.Abs("testdata/rewrites/wrong.sh")
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir("../../..")
	programs, err := filepath.Glob("pangaea/programs/*.cbl")
	if err != nil {
		t.Fatal(err)
	}
	system, err := blueprint.Analyze(programs, "pangaea/copybooks", cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	bp := &system.Programs[slices.IndexFunc(system.Programs, func(p blueprint.Program) bool { return p.Name == "BANKDATE" })]
	reports := map[string][]byte{}
	for name, with := range map[string]string{"bankdate.md": rewrite, "wrong.md": wrong} {
		c, err := Conform(system, bp, "pangaea/copybooks", "oracle/pangaea/BANKDATE", t.TempDir(), with)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if conforms := name == "bankdate.md"; c.Conforms() != conforms {
			t.Errorf("%s: conforms %v, want %v", name, c.Conforms(), conforms)
		}
		reports[name] = c.Markdown()
	}
	golden.CompareFiles(t, reports, "survey/internal/oracle/testdata/conform")
}

// TestInKeyOrder checks that an indexed file's records are given to a
// replacement in key order, the first of a key kept, as the loader keeps
// it: HOLDEXT's holds, from its first case, reversed and one of them
// repeated.
func TestInKeyOrder(t *testing.T) {
	t.Chdir("../../..")
	system, err := blueprint.Analyze([]string{"pangaea/programs/HOLDEXT.cbl"}, "pangaea/copybooks", cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	p, err := prepare(system, &system.Programs[0], "pangaea/copybooks")
	if err != nil {
		t.Fatal(err)
	}
	cases, err := p.readCases("oracle/pangaea/HOLDEXT")
	if err != nil {
		t.Fatal(err)
	}
	holds := cases[0]["HOLD-FILE"]
	if len(holds) != 2 {
		t.Fatalf("case 01 has %d holds, want 2", len(holds))
	}
	f := p.blueprint.Files[slices.IndexFunc(p.blueprint.Files, func(f blueprint.File) bool { return f.Name == "HOLD-FILE" })]
	repeated := slices.Clone(holds[1])
	repeated[len(repeated)-1] = 'X'
	got := p.inKeyOrder(f, [][]byte{repeated, holds[0], holds[1]})
	if want := [][]byte{holds[0], repeated}; !slices.EqualFunc(got, want, bytes.Equal) {
		t.Errorf("got %q, want %q", got, want)
	}
}
