package rules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vlumi/copal/cobol"
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

func TestTables(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/RULES.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, table := range Tables(p) {
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
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
