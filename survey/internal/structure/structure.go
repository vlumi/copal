// Package structure finds a program's procedural structure: its paragraphs
// and sections, which of them perform which, and where control goes from
// each by GO TO or by running on into the next. A SORT performs its input
// and output procedures, so they count as performed too.
package structure

import (
	"fmt"
	"slices"

	"github.com/vlumi/copal/bedrock/cobol"
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

// Paragraphs analyzes the program's procedure division.
func Paragraphs(p *cobol.Program) ([]Paragraph, error) {
	paragraphs := make([]Paragraph, len(p.Paragraphs))
	for i, para := range p.Paragraphs {
		paragraphs[i] = Paragraph{Name: para.Name, Section: para.Section, Comment: para.Comment, Pos: para.Pos}
		var err error
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			for _, ref := range s.Procedures {
				switch {
				case err != nil:
				case ref.Thru != "" && !known(p, ref.Thru):
					err = fmt.Errorf("%s: %s THRU %s: no such paragraph", s.Pos, ref.Name, ref.Thru)
				case !known(p, ref.Name):
					err = fmt.Errorf("%s: %s: no such paragraph", s.Pos, ref.Name)
				case s.Verb == "GO":
					paragraphs[i].GoTos = append(paragraphs[i].GoTos, GoTo{Target: ref.Name, Pos: s.Pos})
				default:
					paragraphs[i].Performs = append(paragraphs[i].Performs,
						Perform{Target: ref.Name, Thru: ref.Thru, Sort: s.Verb == "SORT", Pos: s.Pos})
				}
			}
		})
		if err != nil {
			return nil, err
		}
	}
	ends := map[int]bool{}
	for _, para := range paragraphs {
		for _, perform := range para.Performs {
			_, last, _ := p.Span(cobol.ProcedureRef{Name: perform.Target, Thru: perform.Thru})
			ends[last] = true
		}
	}
	for _, para := range paragraphs {
		for _, perform := range para.Performs {
			first, last, ok := p.Span(cobol.ProcedureRef{Name: perform.Target, Thru: perform.Thru})
			if !ok {
				return nil, fmt.Errorf("%s: %s THRU %s: the range runs backward", perform.Pos, perform.Target, perform.Thru)
			}
			for i := first; i <= last; i++ {
				addOnce(&paragraphs[i].PerformedBy, para.Name)
				if i < last && !ends[i] && !EndsUnconditionally(p.Paragraphs[i].Statements) {
					paragraphs[i].FallsInto = paragraphs[i+1].Name
				}
			}
		}
	}
	for _, para := range paragraphs {
		for _, goTo := range para.GoTos {
			first, _, _ := p.Span(cobol.ProcedureRef{Name: goTo.Target})
			addOnce(&paragraphs[first].ReachedBy, para.Name)
		}
		if para.FallsInto != "" {
			first, _, _ := p.Span(cobol.ProcedureRef{Name: para.FallsInto})
			addOnce(&paragraphs[first].ReachedBy, para.Name)
		}
	}
	return paragraphs, nil
}

func addOnce(names *[]string, name string) {
	if !slices.Contains(*names, name) {
		*names = append(*names, name)
	}
}

// known reports whether a name is a paragraph or section of the program.
func known(p *cobol.Program, name string) bool {
	_, _, ok := p.Span(cobol.ProcedureRef{Name: name})
	return ok
}

// EndsUnconditionally reports whether control never runs off the end of
// the statements: the last one is a GO TO with a single target, STOP RUN,
// GOBACK or EXIT PROGRAM, or an IF or EVALUATE all of whose branches end
// so, an ELSE or WHEN OTHER among them.
func EndsUnconditionally(statements []cobol.Statement) bool {
	if len(statements) == 0 {
		return false
	}
	last := statements[len(statements)-1]
	switch last.Verb {
	case "GO":
		return len(last.Procedures) == 1
	case "STOP", "GOBACK", "EXIT PROGRAM":
		return true
	case "IF", "EVALUATE":
		otherwise := false
		for _, b := range last.Branches {
			otherwise = otherwise || b.Kind == "ELSE" || b.Kind == "OTHER"
			if !EndsUnconditionally(b.Statements) {
				return false
			}
		}
		return otherwise
	}
	return false
}
