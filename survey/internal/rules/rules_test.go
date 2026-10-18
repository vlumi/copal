package rules

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
)

// layout writes a table as text: a line per condition and per action,
// with a column per rule.
func layout(t Table) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s\n", t.Paragraph, t.Loop)
	for i, c := range t.Conditions {
		var cells []string
		for _, r := range t.Rules {
			cells = append(cells, fmt.Sprintf("%-1s", r.Entries[i]))
		}
		fmt.Fprintf(&b, "%s | %s\n", strings.Join(cells, " "), c)
	}
	for i, a := range t.Actions {
		var cells []string
		for _, r := range t.Rules {
			cell := " "
			for order, action := range r.Actions {
				if action == i {
					cell = fmt.Sprint(order + 1)
				}
			}
			cells = append(cells, cell)
		}
		fmt.Fprintf(&b, "%s | %s\n", strings.Join(cells, " "), a)
	}
	return b.String()
}

// byName says a statement may change the names it writes, and anything
// at all if it performs a procedure.
func byName(s cobol.Statement, names []string) bool {
	changes := false
	cobol.Walk([]cobol.Statement{s}, func(n cobol.Statement) {
		changes = changes || len(n.Procedures) > 0 && n.Verb != "GO" ||
			slices.ContainsFunc(n.Writes, func(w cobol.Operand) bool { return slices.Contains(names, w.Name) })
	})
	return changes
}

func TestTables(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/RULES.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	var refs [][]string
	for _, table := range Tables(p, byName) {
		if table.Paragraph == "WIDE-PARA" {
			checkWide(t, table)
			continue
		}
		if table.Paragraph == "FEE-PARA" || table.Loop != "" {
			refs = append(refs, table.Refs...)
		}
		got = append(got, layout(table))
	}
	want := []string{
		`FEE-PARA 
Y N N N | AM-CLOSED
  Y Y N | AM-CURRENT AND AM-OPEN
  Y N   | AM-BALANCE < 1000
1       | SET REJECTED TO TRUE
  1     | MOVE 2.50 TO FEE
  2 1   | ADD 1 TO FEE-COUNT
`,
		`CHAIN-PARA 
Y N N | MR-OPEN
  N Y | HELD
1     | GO TO OPEN-PARA
  1   | SET EXISTS TO TRUE
  2 1 | GO TO REJECT-PARA
`,
		`READ-PARA 
Y N | READ IN-FILE AT END
1 1 | READ IN-FILE
2   | MOVE HIGH-VALUES TO IN-KEY
  2 | ADD 1 TO READ-COUNT
`,
		`LOOP-PARA PERFORM VARYING WS-I FROM 1 BY 1 UNTIL WS-I > 3
Y N | WS-RATE(WS-I) > 0
1   | MOVE WS-RATE(WS-I) TO WS-BEST
`,
		`EXIT-PARA 
Y N | DONE = 1
1   | EXIT PARAGRAPH
  1 | DISPLAY 'NOT DONE'
`,
		`AGAIN-PARA 
Y Y N | OK = 1
Y N   | OK = 1, again
1 1   | PERFORM OPEN-PARA
2     | DISPLAY 'STILL OK'
`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	wantRefs := [][]string{{"AM-CLOSED"}, {"AM-CURRENT", "AM-OPEN"}, {"AM-BALANCE"}, {"WS-RATE", "WS-I"}}
	if !slices.EqualFunc(refs, wantRefs, slices.Equal) {
		t.Errorf("condition refs %q, want %q", refs, wantRefs)
	}
}

// checkWide checks a paragraph whose decisions make more paths than a
// table holds: the decision that would go over is shown as an action, and
// what comes after it is still there.
func checkWide(t *testing.T, table Table) {
	t.Helper()
	if len(table.Rules) != maxPaths {
		t.Errorf("%d rules, want %d", len(table.Rules), maxPaths)
	}
	for _, want := range []string{"IF B = 1", "DISPLAY 'AFTER'"} {
		if !slices.Contains(table.Actions, want) {
			t.Errorf("no action %q in %q", want, table.Actions)
		}
	}
}
