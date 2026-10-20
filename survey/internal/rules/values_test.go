package rules

import (
	"slices"
	"strings"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
)

// TestGivenValues checks that a PERFORM of a paragraph that only gives a
// field values leaves it one of those or what it had, and that a MOVE of
// SPACES to a group blanks its fields of text.
func TestGivenValues(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/VALUES.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	writers := writersIn(p)
	knowledge := Knowledge{
		MayChange:  func(s cobol.Statement, names []string) bool { return len(writers(s, names)) > 0 },
		Writers:    writers,
		TextFields: func(group string) []string { return map[string][]string{"RESULT-GROUP": {"RESULT"}}[group] },
	}
	want := `GIVE-PARA 
Y Y N | KIND = 'A'
Y N Y | RESULT = SPACES
  N   | RESULT = 'OK'
1 1 1 | MOVE SPACES TO RESULT-GROUP
2 2   | PERFORM DECLINE-PARA
3   2 | DISPLAY 'NONE'
4 3 3 | STOP RUN
`
	if got := layout(Tables(p, knowledge)[0]); got != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
}

// writersIn finds the statements a statement runs that write any of the
// names: it, those nested in it, and those of the paragraphs it performs.
func writersIn(p *cobol.Program) func(cobol.Statement, []string) []cobol.Statement {
	return func(s cobol.Statement, names []string) []cobol.Statement {
		var writers []cobol.Statement
		visit := func(n cobol.Statement) {
			if slices.ContainsFunc(n.Writes, func(w cobol.Operand) bool { return slices.Contains(names, w.Name) }) {
				writers = append(writers, n)
			}
		}
		cobol.Walk([]cobol.Statement{s}, func(n cobol.Statement) {
			visit(n)
			for _, r := range n.Procedures {
				for _, para := range p.Paragraphs {
					if strings.EqualFold(para.Name, r.Name) {
						cobol.Walk(para.Statements, visit)
					}
				}
			}
		})
		return writers
	}
}

// TestCopiesAndDigits checks that a field copied from another holds its
// value, and that a number moved into a field of digits gives the fields
// of digits over its bytes their digits.
func TestCopiesAndDigits(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/COPIES.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	writers := writersIn(p)
	knowledge := Knowledge{
		MayChange: func(s cobol.Statement, names []string) bool { return len(writers(s, names)) > 0 },
		Writers:   writers,
		Copies:    func(from, to string) bool { return true },
		DigitFields: func(field string) (int, []DigitField) {
			if field != "MONTH-NUMBER" {
				return 0, nil
			}
			return 6, []DigitField{{Name: "CM-MONTH", Offset: 4, Length: 2}}
		},
	}
	want := []string{`COPY-PARA 
Y N | REQUEST-KEY < MASTER-KEY
N Y | MASTER-KEY = ACTIVE-KEY
1 1 | MOVE MASTER-KEY TO ACTIVE-KEY
2   | MOVE REQUEST-KEY TO ACTIVE-KEY
  2 | DISPLAY 'MASTER'
`, `DIGITS-PARA 
N | CM-MONTH = 12
1 | MOVE 0 TO MONTH-NUMBER
`}
	var got []string
	for _, table := range Tables(p, knowledge) {
		got = append(got, layout(table))
	}
	if !slices.Equal(got, want) {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
