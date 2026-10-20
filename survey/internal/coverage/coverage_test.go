package coverage

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/copal/bedrock/run"
	"github.com/vlumi/copal/survey/internal/dataflow"
	"github.com/vlumi/copal/survey/internal/rules"
)

// traced compiles a program with trace code, runs it, and reads its trace.
func traced(t *testing.T, name string) (*cobol.Program, []Event) {
	t.Helper()
	programs, dir := t.TempDir(), t.TempDir()
	if _, err := (run.Build{Dialect: cobol.IBM, Sources: "testdata", Copybooks: "testdata", Out: programs, Trace: true}).Programs(); err != nil {
		t.Fatal(err)
	}
	trace := filepath.Join(dir, "trace")
	if _, err := (run.Step{Programs: programs, Program: name, Stdout: io.Discard, TraceFile: trace}).Run(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	events, err := Read(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	p, err := cobol.ReadProgram(filepath.Join("testdata", name+".cbl"), nil)
	if err != nil {
		t.Fatal(err)
	}
	return p, events
}

// Each pass through a paragraph or a loop reaches the rule whose
// statements it ran: the EVALUATE's first two WHENs once each, not WHEN
// OTHER, and the loop's IF once taken and twice not.
func TestReplay(t *testing.T) {
	p, events := traced(t, "COVER")
	got := map[string][]int{}
	records, err := layout.Program(p, cobol.IBM)
	if err != nil {
		t.Fatal(err)
	}
	for _, table := range Replay(p, rules.Tables(p, rules.Knowledge{MayChange: dataflow.Analyze(p, records).MayChange}), events) {
		if table.Unmatched > 0 {
			t.Errorf("%s %s: %d passes matched no rule", table.Table.Paragraph, table.Table.Loop, table.Unmatched)
		}
		got[table.Table.Paragraph+" "+table.Table.Loop] = table.Reached
	}
	want := map[string][]int{
		"CHECK-PARA ":                      {1, 1, 0},
		"MAIN-PARA PERFORM UNTIL WS-N = 3": {1, 2},
	}
	for table, reached := range want {
		if g := got[table]; !slices.Equal(g, reached) {
			t.Errorf("%q reached %v, want %v", table, g, reached)
		}
	}
}
