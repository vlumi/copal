package rules

import (
	"os"
	"strings"
	"testing"

	"github.com/vlumi/cobolith/cobol"
)

func TestEntry(t *testing.T) {
	knowledge := Knowledge{MayChange: byName, InitialValues: map[string]string{"TOTAL": "0"}}
	want := []string{
		`MAIN-PARA 
Y Y | TOTAL = 0
Y N | ERR = SPACES
1 1 | DISPLAY 'NONE'
2   | PERFORM CHECK-PARA
3   | PERFORM JUMPED-PARA
  2 | PERFORM LOOP-PARA UNTIL DONE = 1
  3 | GO TO JUMPED-PARA
`,
		`CHECK-PARA 
Y | ERR = SPACES
1 | DISPLAY 'GOOD'
`,
		`LOOP-PARA 
N | DONE = 1
`,
		`JUMPED-PARA 
Y N | ERR = SPACES
1   | DISPLAY 'GOOD'
2 1 | STOP RUN
`,
	}
	if got := entryTables(t, "testdata/ENTRY.cbl", knowledge); got != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", got, strings.Join(want, "\n"))
	}
}

// TestEntryOfCalled checks that a program a caller passes items to does
// not start knowing the values its working storage starts with, which a
// call before may have changed.
func TestEntryOfCalled(t *testing.T) {
	source, err := os.ReadFile("testdata/ENTRY.cbl")
	if err != nil {
		t.Fatal(err)
	}
	called := strings.Replace(string(source), "PROCEDURE DIVISION.", "PROCEDURE DIVISION USING LK-A.", 1)
	path := t.TempDir() + "/CALLED.cbl"
	if err := os.WriteFile(path, []byte(called), 0o644); err != nil {
		t.Fatal(err)
	}
	knowledge := Knowledge{MayChange: byName, InitialValues: map[string]string{"TOTAL": "0"}}
	got := entryTables(t, path, knowledge)
	if want := "Y Y N N | TOTAL = 0\n"; !strings.Contains(got, want) {
		t.Errorf("got\n%s\nwant a line %q", got, want)
	}
}

func entryTables(t *testing.T, path string, knowledge Knowledge) string {
	t.Helper()
	p, err := cobol.ReadProgram(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, table := range Tables(p, knowledge) {
		got = append(got, layout(table))
	}
	return strings.Join(got, "\n")
}
