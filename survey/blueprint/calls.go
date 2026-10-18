package blueprint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/survey/internal/dataflow"
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

// callWrites tells, for a program's CALLs, which of the items passed by
// reference the program called writes: those of the caller's fields over
// the bytes of the called program's USING items it writes, at the same
// offset in the argument passed; the whole argument where no field of the
// caller's is over them. A program outside the system may write anything
// it is passed by reference.
func callWrites(caller Program, programs []Program) dataflow.CallWrites {
	return func(call cobol.Statement) ([]string, bool) {
		i := slices.IndexFunc(programs, func(p Program) bool { return p.Name == call.Calls })
		if i < 0 {
			return nil, false
		}
		called := programs[i]
		written := map[string]bool{}
		for _, para := range called.Paragraphs {
			for _, w := range para.Writes {
				written[w] = true
			}
		}
		var items []string
		for n, arg := range call.Reads {
			if n >= len(called.Using) || !slices.ContainsFunc(call.Writes, func(w cobol.Operand) bool { return w.Text == arg.Text }) {
				continue
			}
			param, ok := recordField(called, called.Using[n])
			argField, found := recordField(caller, arg.Name)
			if !ok || !found {
				continue
			}
			walk(param, func(f Field) {
				if !written[f.Name] {
					return
				}
				name, ok := fieldOver(argField, argField.Start+f.Start-param.Start, f.Length)
				if !ok {
					name = arg.Name
				}
				addOnce(&items, name)
			})
		}
		return items, true
	}
}

// recordField finds a field of a program's records by name.
func recordField(p Program, name string) (Field, bool) {
	for _, r := range p.Records {
		var found *Field
		walk(r.Field, func(f Field) {
			if found == nil && f.Name == name {
				found = &f
			}
		})
		if found != nil {
			return *found, true
		}
	}
	return Field{}, false
}

// fieldOver is the field of a group over exactly the bytes, if any.
func fieldOver(f Field, start, length int) (string, bool) {
	if f.Start == start && f.Length == length && f.Name != "FILLER" {
		return f.Name, true
	}
	for _, child := range f.Fields {
		if name, ok := fieldOver(child, start, length); ok {
			return name, true
		}
	}
	return "", false
}
