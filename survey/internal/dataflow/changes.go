package dataflow

import (
	"slices"

	"github.com/vlumi/copal/bedrock/cobol"
)

// CallWrites tells which of the items a CALL passes the program called may
// write, and whether it knows: for a program outside the system, it does
// not.
type CallWrites func(call cobol.Statement) (items []string, known bool)

// WithCalls is the flow with a way to know what a CALL writes; without
// one, a program called may write every item passed to it by reference.
func (f Flow) WithCalls(writes CallWrites) Flow {
	f.callWrites = writes
	return f
}

// MayChange reports whether a statement may change any of the data items
// named: whether it, a statement nested in it, or a statement of a
// paragraph it performs, or of the DECLARATIVES its failure on a file
// runs, at any depth, writes bytes one of them shares.
// A program called may change what WithCalls says it writes, or else the
// items passed to it by reference, and a condition name stands for its
// field.
func (f Flow) MayChange(s cobol.Statement, names []string) bool {
	return len(f.Writers(s, names)) > 0
}

// Writers lists the statements a statement runs that may change any of
// the data items named, as MayChange finds them: it, those nested in it,
// and those of the paragraphs it performs and of the DECLARATIVES its
// failure on a file runs, at any depth.
func (f Flow) Writers(s cobol.Statement, names []string) []cobol.Statement {
	items := make([]string, len(names))
	for i, name := range names {
		items[i] = f.names.item(name)
	}
	var writers []cobol.Statement
	f.runs(s, map[string]bool{}, func(n cobol.Statement, writes []string) {
		if slices.ContainsFunc(writes, func(w string) bool {
			return slices.ContainsFunc(items, func(item string) bool { return w == item || f.Share(w, item) })
		}) {
			writers = append(writers, n)
		}
	})
	return writers
}

// runs visits the statements a statement runs, with the items each
// writes: it, those nested in it, and those of the paragraphs it
// performs, each paragraph once.
func (f Flow) runs(s cobol.Statement, performed map[string]bool, visit func(n cobol.Statement, writes []string)) {
	cobol.Walk([]cobol.Statement{s}, func(n cobol.Statement) {
		_, writes := f.names.uses(n)
		if n.Verb == "CALL" && f.callWrites != nil {
			if called, known := f.callWrites(n); known {
				writes = called
			}
		}
		visit(n, writes)
		if n.Verb == "GO" {
			return
		}
		refs := slices.Clone(n.Procedures)
		for _, section := range f.program.ErrorProcedures(n) {
			refs = append(refs, cobol.ProcedureRef{Name: section})
		}
		for _, ref := range refs {
			for _, para := range f.span(ref) {
				if performed[para.Name] {
					continue
				}
				performed[para.Name] = true
				for _, st := range para.Statements {
					f.runs(st, performed, visit)
				}
			}
		}
	})
}

// span lists the paragraphs a procedure reference names: a paragraph, the
// paragraphs of a section, or those from one through another.
func (f Flow) span(ref cobol.ProcedureRef) []cobol.Paragraph {
	first, last, ok := f.program.Span(ref)
	if !ok {
		return nil
	}
	return f.program.Paragraphs[first : last+1]
}
