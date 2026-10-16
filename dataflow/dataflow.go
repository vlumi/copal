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
	storage    *storage
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
// record a WRITE or RELEASE sends to it. From is a literal, as written,
// for a value the statement gives itself, such as MOVE 'TXPOST' TO
// RJ-PROGRAM or SET TD-INTEREST TO TRUE, written as "condition
// TD-INTEREST".
type Edge struct {
	From, To         string
	FromFile, ToFile bool
	FromLiteral      bool
	Paragraph        string
	Verb             string
	Pos              cobol.Pos
}

// Analyze finds the program's data flow, given its records laid out.
func Analyze(p *cobol.Program, records *layout.ProgramRecords) Flow {
	names := newNames(records)
	flow := Flow{storage: newStorage(records)}
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
	return flow
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

func newNames(records *layout.ProgramRecords) *names {
	n := &names{conditions: map[string]string{}, records: map[string][]string{}, files: map[string]string{}}
	all := records.WorkingStorage
	for _, f := range records.Files {
		for _, r := range f.Records {
			n.records[f.File] = append(n.records[f.File], r.Name)
			n.files[r.Name] = f.File
		}
		all = append(all, f.Records...)
	}
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
	return n
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
		if e.FromFile || e.ToFile || e.FromLiteral || e.From != e.To {
			edges = append(edges, e)
		}
	}
	var sources []Edge
	for _, op := range s.Reads {
		switch {
		case op.Name != "":
			sources = append(sources, Edge{From: n.item(op.Name)})
		case len(op.Refs) == 0:
			sources = append(sources, Edge{From: op.Text, FromLiteral: true})
		}
		if op.Name == "" {
			for _, ref := range op.Refs {
				sources = append(sources, Edge{From: n.item(ref)})
			}
		}
	}
	for _, target := range s.Writes {
		if target.Name == "" {
			continue
		}
		if s.Verb == "SET" && len(s.Reads) == 0 {
			add(Edge{From: "condition " + target.Name, FromLiteral: true, To: n.item(target.Name)})
		}
		for _, source := range sources {
			source.To = n.item(target.Name)
			add(source)
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
