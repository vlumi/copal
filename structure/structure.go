// Package structure finds a program's procedural structure: its paragraphs
// and sections, and which of them perform which. A SORT performs its input
// and output procedures, so they count as performed too.
package structure

import (
	"fmt"
	"slices"

	"github.com/vlumi/copal/cobol"
)

// Paragraph is a paragraph with the PERFORM statements in it and the
// paragraphs that perform it.
type Paragraph struct {
	Name        string
	Section     string
	Comment     []string
	Pos         cobol.Pos
	Performs    []Perform
	PerformedBy []string
}

// Perform is a PERFORM of a paragraph or section, or of the range from
// Target through Thru, or a SORT's input or output procedure.
type Perform struct {
	Target string
	Thru   string
	Sort   bool
	Pos    cobol.Pos
}

// Paragraphs analyzes the program's procedure division. A PERFORM followed
// by a procedure name performs it; any other PERFORM is an inline loop,
// whose body belongs to the paragraph around it.
func Paragraphs(p *cobol.Program) ([]Paragraph, error) {
	index := map[string]int{}
	for i, para := range p.Paragraphs {
		index[para.Name] = i
	}
	paragraphs := make([]Paragraph, len(p.Paragraphs))
	for i, para := range p.Paragraphs {
		performs, err := findPerforms(para.Body, index)
		if err != nil {
			return nil, err
		}
		paragraphs[i] = Paragraph{Name: para.Name, Section: para.Section, Comment: para.Comment, Pos: para.Pos, Performs: performs}
	}
	for _, para := range paragraphs {
		for _, perform := range para.Performs {
			target := &paragraphs[index[perform.Target]]
			if !slices.Contains(target.PerformedBy, para.Name) {
				target.PerformedBy = append(target.PerformedBy, para.Name)
			}
		}
	}
	return paragraphs, nil
}

// findPerforms finds the PERFORM statements that name a procedure in index,
// and the input and output procedures of SORT statements.
func findPerforms(body []cobol.Token, index map[string]int) ([]Perform, error) {
	var performs []Perform
	for i := 0; i+1 < len(body); i++ {
		var perform Perform
		switch {
		case body[i].Is("PERFORM"):
			perform = Perform{Pos: body[i].Pos}
		case body[i].Is("PROCEDURE") && i > 0 && (body[i-1].Is("INPUT") || body[i-1].Is("OUTPUT")):
			perform = Perform{Sort: true, Pos: body[i-1].Pos}
			if body[i+1].Is("IS") {
				i++
			}
		default:
			continue
		}
		if i+1 >= len(body) || !known(index, body[i+1]) {
			continue
		}
		perform.Target = body[i+1].Text
		if i+3 < len(body) && (body[i+2].Is("THRU") || body[i+2].Is("THROUGH")) {
			perform.Thru = body[i+3].Text
			if !known(index, body[i+3]) {
				return nil, fmt.Errorf("%s: %s THRU %s: no such paragraph", perform.Pos, perform.Target, perform.Thru)
			}
		}
		performs = append(performs, perform)
	}
	return performs, nil
}

func known(index map[string]int, t cobol.Token) bool {
	_, ok := index[t.Text]
	return ok && t.Kind == cobol.Word
}
