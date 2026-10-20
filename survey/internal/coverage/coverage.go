// Package coverage measures which rules of a program's decision tables a
// run reached, from the trace GnuCOBOL writes of a program compiled with
// trace code: every paragraph it enters and every statement it executes,
// by source line.
//
// The trace is cut into passes: through a paragraph, from entering it to
// entering it again, and through an inline loop's body, from one test of
// the loop to the next. A pass reached a rule when the statements it ran
// of its own, those of the loops inside it left out, are the rule's.
package coverage

import (
	"bufio"
	"io"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/survey/internal/rules"
)

// Event is a line of a trace: a program entering a paragraph or section,
// or executing the statement at a line of a source file. Names are
// upper-cased, as COBOL words are whatever their case in the source.
type Event struct {
	Program   string
	Paragraph string // the paragraph or section entered; empty for a statement
	File      string // the source file's name, without its directory
	Line      int
}

var (
	sourceLine = regexp.MustCompile(`^Source:\s+'(.+)'$`)
	traceLine  = regexp.MustCompile(`^Program-Id:\s+(\S+)\s+(.*?)\s*Line:\s+(\d+)$`)
	paragraph  = regexp.MustCompile(`^(?:Paragraph|Section):\s+(\S+)`)
	entry      = regexp.MustCompile(`^Entry:`)
)

// Read reads a trace into its events, leaving out entries to programs,
// which no rule needs, and the tests of WHEN, which are no statement but
// share the line of the one after them when it is written on the same
// line. A section entered starts its own statements, the paragraph named
// after it.
func Read(r io.Reader) ([]Event, error) {
	var events []Event
	file := ""
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if m := sourceLine.FindStringSubmatch(line); m != nil {
			file = filepath.Base(m[1])
			continue
		}
		m := traceLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		n, _ := strconv.Atoi(m[3])
		e := Event{Program: strings.ToUpper(m[1]), File: file, Line: n}
		switch p := paragraph.FindStringSubmatch(m[2]); {
		case p != nil:
			e.Paragraph = strings.ToUpper(p[1])
		case entry.MatchString(m[2]), strings.TrimSpace(m[2]) == "WHEN":
			continue
		}
		events = append(events, e)
	}
	return events, scanner.Err()
}

// Table is how often a run reached each rule of a decision table, and how
// many passes through the table's paragraph or loop matched no rule.
type Table struct {
	Table     *rules.Table
	Reached   []int // by rule
	Unmatched int
}

// line is a statement's place in the source, as a trace gives it.
type line struct {
	file string
	n    int
}

func lineOf(pos cobol.Pos) line { return line{filepath.Base(pos.File), pos.Line} }

// scopes knows which paragraph or inline loop each statement of a program
// is in, and the decision table of each, if it has one. A scope is named
// by its paragraph, or for a loop, by the line of its PERFORM.
type scopes struct {
	owner  map[line]string
	loops  map[line]string
	tables map[string]*Table
}

func newScopes(p *cobol.Program, tables []Table) scopes {
	s := scopes{owner: map[line]string{}, loops: map[line]string{}, tables: map[string]*Table{}}
	var walk func(scope string, statements []cobol.Statement)
	walk = func(scope string, statements []cobol.Statement) {
		for _, st := range statements {
			s.owner[lineOf(st.Pos)] = scope
			for _, b := range st.Branches {
				if st.Verb == "PERFORM" {
					loop := loopScope(st.Pos)
					s.loops[lineOf(st.Pos)] = loop
					walk(loop, b.Statements)
					continue
				}
				walk(scope, b.Statements)
			}
		}
	}
	for _, para := range p.Paragraphs {
		walk(para.Name, para.Statements)
	}
	for i := range tables {
		t := &tables[i]
		scope := t.Table.Paragraph
		if t.Table.Loop != "" {
			scope = loopScope(t.Table.Pos)
		}
		s.tables[scope] = t
	}
	return s
}

func loopScope(pos cobol.Pos) string {
	return "loop " + filepath.Base(pos.File) + ":" + strconv.Itoa(pos.Line)
}

// Replay measures which rules of the program's tables the trace's events
// reached.
func Replay(p *cobol.Program, tables []rules.Table, events []Event) []Table {
	coverage := make([]Table, len(tables))
	for i := range tables {
		coverage[i] = Table{Table: &tables[i], Reached: make([]int, len(tables[i].Rules))}
	}
	s := newScopes(p, coverage)
	passes := map[string][]line{}
	end := func(scope string) {
		if t, ok := s.tables[scope]; ok && len(passes[scope]) > 0 {
			t.match(passes[scope])
		}
		passes[scope] = nil
	}
	for _, e := range events {
		if e.Program != p.ID {
			continue
		}
		if e.Paragraph != "" {
			end(e.Paragraph)
			continue
		}
		l := line{e.File, e.Line}
		if loop, ok := s.loops[l]; ok {
			end(loop)
		}
		if scope, ok := s.owner[l]; ok {
			passes[scope] = append(passes[scope], l)
		}
	}
	for scope := range passes {
		end(scope)
	}
	return coverage
}

// match counts a pass for the rules whose statements are those it ran,
// leaving out the statements of collapsed decisions, which no rule tells
// apart.
func (t *Table) match(pass []line) {
	ran := map[line]bool{}
	for _, l := range pass {
		if !slices.ContainsFunc(t.Table.Collapsed, func(pos cobol.Pos) bool { return lineOf(pos) == l }) {
			ran[l] = true
		}
	}
	matched := false
	for i, r := range t.Table.Rules {
		steps := map[line]bool{}
		for _, pos := range r.Steps {
			steps[lineOf(pos)] = true
		}
		if sameLines(steps, ran) {
			t.Reached[i]++
			matched = true
		}
	}
	if !matched {
		t.Unmatched++
	}
}

func sameLines(a, b map[line]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for l := range a {
		if !b[l] {
			return false
		}
	}
	return true
}
