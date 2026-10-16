package blueprint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vlumi/copal/cobol"
)

// Call is a CALL statement naming the program it calls, with the items it
// passes. External is set when the program called is not one of those
// analyzed.
type Call struct {
	Program   string   `json:"program"`
	Using     []string `json:"using,omitempty"`
	External  bool     `json:"external,omitempty"`
	Paragraph string   `json:"paragraph"`
	Source    string   `json:"source"`
}

// programCalls lists a program's CALL statements that name the program
// they call.
func programCalls(p *cobol.Program) []Call {
	var calls []Call
	for _, para := range p.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if s.Calls == "" {
				return
			}
			call := Call{Program: s.Calls, Paragraph: para.Name, Source: s.Pos.String()}
			for _, op := range s.Reads {
				call.Using = append(call.Using, op.Text)
			}
			calls = append(calls, call)
		})
	}
	return calls
}

// linkCalls marks the calls of programs not analyzed as external, and
// lists for each program those that call it.
func linkCalls(programs []Program) {
	byName := map[string]*Program{}
	for i := range programs {
		byName[programs[i].Name] = &programs[i]
	}
	for i := range programs {
		for j, call := range programs[i].Calls {
			called, ok := byName[call.Program]
			if !ok {
				programs[i].Calls[j].External = true
				continue
			}
			if !slices.Contains(called.CalledBy, programs[i].Name) {
				called.CalledBy = append(called.CalledBy, programs[i].Name)
			}
		}
	}
}

// writeCalls writes what a program is passed, by whom it is called and
// what it calls.
func writeCalls(b *strings.Builder, dir string, p Program) {
	if len(p.Using) == 0 && len(p.CalledBy) == 0 && len(p.Calls) == 0 {
		return
	}
	b.WriteString("\n## Calls\n\n")
	if len(p.Using) > 0 || len(p.CalledBy) > 0 {
		b.WriteString("Called")
		if len(p.Using) > 0 {
			fmt.Fprintf(b, " with %s", strings.Join(p.Using, ", "))
		}
		fmt.Fprintf(b, " by %s.\n", programLinks(p.CalledBy, ""))
	}
	if len(p.Calls) == 0 {
		return
	}
	if len(p.Using) > 0 || len(p.CalledBy) > 0 {
		b.WriteString("\n")
	}
	b.WriteString("| Calls | Passing | Paragraph | Source |\n| --- | --- | --- | --- |\n")
	for _, c := range p.Calls {
		called := fmt.Sprintf("[%s](%s.md)", c.Program, c.Program)
		if c.External {
			called = c.Program + " (not analyzed)"
		}
		b.WriteString(row(called, strings.Join(c.Using, ", "), c.Paragraph, sourceLink(dir, c.Source, true)))
	}
}
