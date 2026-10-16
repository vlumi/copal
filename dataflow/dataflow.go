// Package dataflow finds what a program does with its data: which data
// items each paragraph reads and writes, and the lineage of every item,
// the statements through which values flow into it from other items and
// from files.
//
// Items are named as declared. A condition name stands for the field it
// is declared on: SET TD-INTEREST TO TRUE writes TD-TRAN-CODE. Names the
// program does not declare, such as RETURN-CODE or an index name, are kept
// as they are.
package dataflow

import (
	"slices"

	"github.com/vlumi/copal/cobol"
	"github.com/vlumi/copal/layout"
)

// Flow is what a program does with its data.
type Flow struct {
	Paragraphs []Paragraph
	Lineage    []Edge
}

// Paragraph is a paragraph's reads and writes, each item once, in the
// order the paragraph first uses it.
type Paragraph struct {
	Name   string
	Reads  []string
	Writes []string
}

// Edge is one statement moving a value from one item into another. From
// is a file for the records a READ or RETURN fills, and To a file for the
// record a WRITE or RELEASE sends to it.
type Edge struct {
	From, To         string
	FromFile, ToFile bool
	Paragraph        string
	Verb             string
	Pos              cobol.Pos
}

// Analyze finds the program's data flow.
func Analyze(p *cobol.Program) (Flow, error) {
	names, err := newNames(p)
	if err != nil {
		return Flow{}, err
	}
	var flow Flow
	for _, para := range p.Paragraphs {
		pf := Paragraph{Name: para.Name}
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			reads, writes := names.uses(s)
			for _, r := range reads {
				addOnce(&pf.Reads, r)
			}
			for _, w := range writes {
				addOnce(&pf.Writes, w)
			}
			flow.Lineage = append(flow.Lineage, names.edges(para.Name, s)...)
		})
		flow.Paragraphs = append(flow.Paragraphs, pf)
	}
	return flow, nil
}

func addOnce(names *[]string, name string) {
	if !slices.Contains(*names, name) {
		*names = append(*names, name)
	}
}

// names resolves the names a program uses: condition names to their
// fields, files to their records and records to their files.
type names struct {
	conditions map[string]string
	records    map[string][]string // file to records
	files      map[string]string   // record to file
}

func newNames(p *cobol.Program) (*names, error) {
	n := &names{conditions: map[string]string{}, records: map[string][]string{}, files: map[string]string{}}
	var all []*layout.Field
	for _, fd := range p.FileSection {
		records, err := layout.Records(fd.Entries)
		if err != nil {
			return nil, err
		}
		for _, r := range records {
			n.records[fd.Name] = append(n.records[fd.Name], r.Name)
			n.files[r.Name] = fd.Name
		}
		all = append(all, records...)
	}
	records, err := layout.Records(p.WorkingStorage)
	if err != nil {
		return nil, err
	}
	all = append(all, records...)
	var index func(f *layout.Field)
	index = func(f *layout.Field) {
		for _, c := range f.Conditions {
			n.conditions[c.Name] = f.Name
		}
		for _, child := range f.Children {
			index(child)
		}
	}
	for _, r := range all {
		index(r)
	}
	return n, nil
}

// item is the data item a name stands for.
func (n *names) item(name string) string {
	if field, ok := n.conditions[name]; ok {
		return field
	}
	return name
}

// uses lists the items a statement reads and writes: its operands and
// condition, a written operand's subscripts read, and the records of the
// files it reads.
func (n *names) uses(s cobol.Statement) (reads, writes []string) {
	for _, op := range s.Reads {
		for _, ref := range op.Refs {
			reads = append(reads, n.item(ref))
		}
	}
	for _, op := range s.Writes {
		if op.Name != "" {
			writes = append(writes, n.item(op.Name))
		}
		for _, ref := range op.Refs {
			if ref != op.Name {
				reads = append(reads, n.item(ref))
			}
		}
	}
	if s.Condition != nil {
		for _, ref := range s.Condition.Refs() {
			reads = append(reads, n.item(ref))
		}
	}
	for _, b := range s.Branches {
		if b.Condition != nil {
			for _, ref := range b.Condition.Refs() {
				reads = append(reads, n.item(ref))
			}
		}
	}
	if s.Verb == "READ" || s.Verb == "RETURN" {
		for _, f := range s.Files {
			writes = append(writes, n.records[f.Name]...)
		}
	}
	return reads, writes
}

// lineageVerbs move values from what they read into what they write.
var lineageVerbs = map[string]bool{
	"MOVE": true, "ADD": true, "SUBTRACT": true, "MULTIPLY": true, "DIVIDE": true, "COMPUTE": true,
	"STRING": true, "SET": true, "WRITE": true, "RELEASE": true, "READ": true, "RETURN": true,
}

// edges lists the values a statement moves. A data item operand moves its
// own value, not its subscripts'; an expression moves every item in it.
func (n *names) edges(paragraph string, s cobol.Statement) []Edge {
	if !lineageVerbs[s.Verb] {
		return nil
	}
	var edges []Edge
	add := func(e Edge) {
		e.Paragraph, e.Verb, e.Pos = paragraph, s.Verb, s.Pos
		if e.FromFile || e.ToFile || e.From != e.To {
			edges = append(edges, e)
		}
	}
	var sources []string
	for _, op := range s.Reads {
		if op.Name != "" {
			sources = append(sources, n.item(op.Name))
		} else {
			for _, ref := range op.Refs {
				sources = append(sources, n.item(ref))
			}
		}
	}
	for _, target := range s.Writes {
		if target.Name == "" {
			continue
		}
		for _, source := range sources {
			add(Edge{From: source, To: n.item(target.Name)})
		}
	}
	switch s.Verb {
	case "READ", "RETURN":
		for _, f := range s.Files {
			for _, record := range n.records[f.Name] {
				add(Edge{From: f.Name, FromFile: true, To: record})
				for _, into := range s.Writes {
					add(Edge{From: record, To: n.item(into.Name)})
				}
			}
		}
	case "WRITE", "RELEASE":
		for _, record := range s.Writes {
			if file, ok := n.files[record.Name]; ok {
				add(Edge{From: record.Name, To: file, ToFile: true})
			}
		}
	}
	return edges
}
