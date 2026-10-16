// Package structure finds a program's procedural structure: its paragraphs
// and sections, which of them perform which, and where control goes from
// each by GO TO or by running on into the next. A SORT performs its input
// and output procedures, so they count as performed too.
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

// Paragraphs analyzes the program's procedure division.
func Paragraphs(p *cobol.Program) ([]Paragraph, error) {
	index := map[string]int{}
	for i, para := range p.Paragraphs {
		index[para.Name] = i
	}
	paragraphs := make([]Paragraph, len(p.Paragraphs))
	for i, para := range p.Paragraphs {
		paragraphs[i] = Paragraph{Name: para.Name, Section: para.Section, Comment: para.Comment, Pos: para.Pos}
		var err error
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			for _, ref := range s.Procedures {
				switch {
				case err != nil:
				case ref.Thru != "" && !known(index, ref.Thru):
					err = fmt.Errorf("%s: %s THRU %s: no such paragraph", s.Pos, ref.Name, ref.Thru)
				case !known(index, ref.Name):
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
				if i < last && !ends[i] && !endsUnconditionally(p.Paragraphs[i].Statements) {
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

func known(index map[string]int, name string) bool {
	_, ok := index[name]
	return ok
}

// endsUnconditionally reports whether control never runs off the end of
// the statements: the last one is a GO TO with a single target, STOP RUN,
// GOBACK or EXIT PROGRAM, or an IF or EVALUATE all of whose branches end
// so, an ELSE or WHEN OTHER among them.
func endsUnconditionally(statements []cobol.Statement) bool {
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
			if !endsUnconditionally(b.Statements) {
				return false
			}
		}
		return otherwise
	}
	return false
}
