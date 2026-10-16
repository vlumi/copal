// Package structure finds a program's procedural structure: its paragraphs
// and sections, which of them perform which, and where control goes from
// each by GO TO or by running on into the next. A SORT performs its input
// and output procedures, so they count as performed too.
//
// The analysis reads tokens, not statements: a paragraph runs on into the
// next unless its last sentence ends in GO TO, STOP RUN, GOBACK or EXIT
// PROGRAM outside any condition.
package structure

import (
	"fmt"
	"slices"

	"github.com/vlumi/copal/cobol"
)

// Paragraph is a paragraph with the PERFORM and GO TO statements in it,
// the paragraph it falls into, if any, and the paragraphs that come to it.
type Paragraph struct {
	Name     string
	Section  string
	Comment  []string
	Pos      cobol.Pos
	Performs []Perform
	GoTos    []GoTo
	// FallsInto is the next paragraph, when the paragraph is inside a
	// PERFORM THRU range, does not end a performed range itself, and
	// control can run on into it.
	FallsInto string
	// PerformedBy are the paragraphs that perform it, alone or as part of
	// a THRU range; ReachedBy those that go to it or fall into it.
	PerformedBy []string
	ReachedBy   []string
}

// GoTo is a GO TO statement's target.
type GoTo struct {
	Target string
	Pos    cobol.Pos
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
		paragraphs[i] = Paragraph{Name: para.Name, Section: para.Section, Comment: para.Comment, Pos: para.Pos,
			Performs: performs, GoTos: findGoTos(para.Body, index)}
	}
	ends := map[int]bool{}
	for _, para := range paragraphs {
		for _, perform := range para.Performs {
			if perform.Thru != "" {
				ends[index[perform.Thru]] = true
			} else {
				ends[index[perform.Target]] = true
			}
		}
	}
	for _, para := range paragraphs {
		for _, perform := range para.Performs {
			first, last := index[perform.Target], index[perform.Target]
			if perform.Thru != "" {
				last = index[perform.Thru]
			}
			if last < first {
				return nil, fmt.Errorf("%s: %s THRU %s: the range runs backward", perform.Pos, perform.Target, perform.Thru)
			}
			for i := first; i <= last; i++ {
				addOnce(&paragraphs[i].PerformedBy, para.Name)
				if i < last && !ends[i] && runsOn(p.Paragraphs[i].Body) {
					paragraphs[i].FallsInto = paragraphs[i+1].Name
				}
			}
		}
	}
	for _, para := range paragraphs {
		for _, goTo := range para.GoTos {
			addOnce(&paragraphs[index[goTo.Target]].ReachedBy, para.Name)
		}
		if para.FallsInto != "" {
			addOnce(&paragraphs[index[para.FallsInto]].ReachedBy, para.Name)
		}
	}
	return paragraphs, nil
}

func addOnce(names *[]string, name string) {
	if !slices.Contains(*names, name) {
		*names = append(*names, name)
	}
}

// findGoTos finds the GO TO statements' targets: one, or several for GO TO
// ... DEPENDING ON.
func findGoTos(body []cobol.Token, index map[string]int) []GoTo {
	var goTos []GoTo
	for i := 0; i < len(body); i++ {
		if !body[i].Is("GO") {
			continue
		}
		j := i + 1
		if j < len(body) && body[j].Is("TO") {
			j++
		}
		for ; j < len(body) && known(index, body[j]); j++ {
			goTos = append(goTos, GoTo{Target: body[j].Text, Pos: body[i].Pos})
		}
	}
	return goTos
}

// opens and closes are the words that open a condition, or a statement
// with conditional phrases, and those that close one. A period closes
// them all.
var (
	opens  = []string{"IF", "EVALUATE", "SEARCH", "AT", "INVALID", "ON"}
	closes = []string{"END-IF", "END-EVALUATE", "END-SEARCH", "END-READ", "END-RETURN", "END-WRITE",
		"END-REWRITE", "END-START", "END-DELETE", "END-COMPUTE", "END-ADD", "END-SUBTRACT", "END-MULTIPLY", "END-DIVIDE"}
)

// runsOn reports whether control can run off the end of a paragraph's
// body: unless its last sentence ends in GO TO, STOP RUN, GOBACK or EXIT
// PROGRAM, outside any condition.
func runsOn(body []cobol.Token) bool {
	end := len(body)
	if end > 0 && body[end-1].Kind == cobol.Period {
		end--
	}
	start := end
	for start > 0 && body[start-1].Kind != cobol.Period {
		start--
	}
	sentence := body[start:end]
	n := len(sentence)
	var last int
	switch {
	case n >= 3 && sentence[n-3].Is("GO") && sentence[n-2].Is("TO"):
		last = n - 3
	case n >= 2 && sentence[n-2].Is("GO"):
		last = n - 2
	case n >= 2 && (sentence[n-2].Is("STOP") && sentence[n-1].Is("RUN") ||
		sentence[n-2].Is("EXIT") && sentence[n-1].Is("PROGRAM")):
		last = n - 2
	case n >= 1 && sentence[n-1].Is("GOBACK"):
		last = n - 1
	default:
		return true
	}
	depth := 0
	for _, t := range sentence[:last] {
		switch {
		case t.Kind != cobol.Word:
		case slices.Contains(opens, t.Text):
			depth++
		case slices.Contains(closes, t.Text):
			depth = max(depth-1, 0)
		}
	}
	return depth > 0
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
