package dataflow

import (
	"slices"

	"github.com/vlumi/copal/bedrock/cobol"
)

// MayChange reports whether a statement may change any of the data items
// named: whether it, a statement nested in it, or a statement of a
// paragraph it performs, at any depth, writes bytes one of them shares.
// A program called may change the items passed to it by reference, and a
// condition name stands for its field.
func (f Flow) MayChange(s cobol.Statement, names []string) bool {
	written := f.written(s, map[string]bool{})
	return slices.ContainsFunc(names, func(name string) bool {
		item := f.names.item(name)
		return slices.ContainsFunc(written, func(w string) bool { return w == item || f.Share(w, item) })
	})
}

// written lists the items a statement writes, with those nested in it and
// those of the paragraphs it performs, each paragraph once.
func (f Flow) written(s cobol.Statement, performed map[string]bool) []string {
	var items []string
	cobol.Walk([]cobol.Statement{s}, func(n cobol.Statement) {
		_, writes := f.names.uses(n)
		items = append(items, writes...)
		if n.Verb == "GO" {
			return
		}
		for _, ref := range n.Procedures {
			for _, para := range f.span(ref) {
				if performed[para.Name] {
					continue
				}
				performed[para.Name] = true
				for _, st := range para.Statements {
					items = append(items, f.written(st, performed)...)
				}
			}
		}
	})
	return items
}

// span lists the paragraphs a procedure reference names: a paragraph, the
// paragraphs of a section, or those from one through another.
func (f Flow) span(ref cobol.ProcedureRef) []cobol.Paragraph {
	paragraphs := f.program.Paragraphs
	first := slices.IndexFunc(paragraphs, func(p cobol.Paragraph) bool { return p.Name == ref.Name || p.Section == ref.Name })
	if first < 0 {
		return nil
	}
	last := first
	end := ref.Name
	if ref.Thru != "" {
		end = ref.Thru
	}
	for i := first; i < len(paragraphs); i++ {
		if paragraphs[i].Name == end || paragraphs[i].Section == end {
			last = i
		}
	}
	return paragraphs[first : last+1]
}
