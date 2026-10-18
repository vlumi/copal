package blueprint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/survey/internal/dataflow"
)

// Link is a relation between the records of two inputs that a program
// relies on: a field of one that it moves into the record key of an
// indexed file, to read that file's record by it, "key"; or a field of
// one that it compares for equality with a field of another, "equal".
type Link struct {
	Kind   string    `json:"kind"`
	From   FileField `json:"from"`
	To     FileField `json:"to"`
	Source string    `json:"source"`
}

// FileField is a field of a record a program reads from a file, or SYSIN.
type FileField struct {
	File  string `json:"file"`
	Field string `json:"field"`
}

// Total is a control total a program checks: a field of an input that it
// compares with an item it adds a field of the same input's records up
// into, Sum, or 1 for each of them, with no Sum; over the records of a
// type, if it adds up only those.
type Total struct {
	File   string `json:"file"`
	Field  string `json:"field"`
	Sum    string `json:"sum,omitempty"`
	Type   string `json:"type,omitempty"`
	Source string `json:"source"`
}

// inputLinks finds the links a program relies on: each field of an
// input it moves into the record key of an indexed file it reads, and
// each pair of fields of different inputs it compares for equality, each
// side's values from that field alone.
func inputLinks(p Program, source *cobol.Program, flow dataflow.Flow) []Link {
	var links []Link
	add := func(l Link) {
		if !slices.ContainsFunc(links, func(m Link) bool {
			return m.Kind == l.Kind && (m.From == l.From && m.To == l.To || l.Kind == "equal" && m.From == l.To && m.To == l.From)
		}) {
			links = append(links, l)
		}
	}
	for _, f := range p.Files {
		if f.Organization != "INDEXED" || f.RecordKey == "" || !slices.ContainsFunc(f.Opened, func(m string) bool { return m == "INPUT" || m == "I-O" }) {
			continue
		}
		key, ok := recordField(p, f.RecordKey)
		if !ok {
			continue
		}
		parts := map[string]bool{}
		walk(key, func(g Field) { parts[g.Name] = true })
		for _, e := range flow.Lineage {
			if !parts[e.To] || e.FromFile || e.FromLiteral {
				continue
			}
			for _, from := range inputOrigins(flow, e.From) {
				if from.File != f.Name {
					add(Link{Kind: "key", From: from, To: FileField{File: f.Name, Field: e.To}, Source: e.Pos.String()})
				}
			}
		}
	}
	for _, c := range conditions(source) {
		for _, pair := range equalities(c.condition) {
			a, aok := soleInputOrigin(flow, pair[0])
			b, bok := soleInputOrigin(flow, pair[1])
			if aok && bok && a.File != b.File {
				add(Link{Kind: "equal", From: a, To: b, Source: c.pos.String()})
			}
		}
	}
	return links
}

// soleInputOrigin is the field of an input an item's values come from, if
// they come from that one alone of the inputs' fields, besides literals,
// such as HIGH-VALUES at the end of a file, and initial values.
func soleInputOrigin(flow dataflow.Flow, item string) (FileField, bool) {
	fields := inputOrigins(flow, item)
	if len(fields) != 1 {
		return FileField{}, false
	}
	return fields[0], true
}

// inputOrigins lists the fields of the inputs an item's values come from.
func inputOrigins(flow dataflow.Flow, item string) []FileField {
	var fields []FileField
	for _, o := range flow.Origins(item) {
		file := o.File
		if o.Stream != "" {
			file = o.Stream
		}
		if file != "" && !slices.Contains(fields, FileField{File: file, Field: o.Item}) {
			fields = append(fields, FileField{File: file, Field: o.Item})
		}
	}
	return fields
}

// condition is a condition a program tests, and where.
type condition struct {
	condition *cobol.Condition
	pos       cobol.Pos
}

// conditions lists every condition a program tests: of IF, of the WHENs
// of EVALUATE and SEARCH, and of PERFORM UNTIL.
func conditions(p *cobol.Program) []condition {
	var list []condition
	for _, para := range p.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if s.Condition != nil {
				list = append(list, condition{s.Condition, s.Pos})
			}
			for _, b := range s.Branches {
				if b.Condition != nil {
					list = append(list, condition{b.Condition, s.Pos})
				}
			}
		})
	}
	return list
}

// equalities lists the pairs of data items a condition compares for
// equality, or inequality, which tells as much of how they relate.
func equalities(c *cobol.Condition) [][2]string {
	var pairs [][2]string
	for _, part := range c.Parts {
		pairs = append(pairs, equalities(part)...)
	}
	if c.Op == "=" && c.Left.Name != "" && c.Right.Name != "" {
		pairs = append(pairs, [2]string{c.Left.Name, c.Right.Name})
	}
	return pairs
}

// programTotals finds the control totals a program checks: each
// comparison for equality of a field of an input with an item that is no
// input's,
// which the program adds a field of the same input up into, or adds 1
// to. A sum is over the records of the type whose fields the summed field
// is among; a count of records of no known type, over those a sum of the
// same input adds up.
func programTotals(p Program, source *cobol.Program, flow dataflow.Flow) []Total {
	var totals []Total
	for _, c := range conditions(source) {
		for _, pair := range equalities(c.condition) {
			t, ok := totalOf(p, flow, pair[0], pair[1])
			if !ok {
				t, ok = totalOf(p, flow, pair[1], pair[0])
			}
			if ok && !slices.ContainsFunc(totals, func(u Total) bool { return u.File == t.File && u.Field == t.Field }) {
				t.Source = c.pos.String()
				totals = append(totals, t)
			}
		}
	}
	for i, t := range totals {
		if t.Sum != "" || t.Type != "" {
			continue
		}
		if j := slices.IndexFunc(totals, func(u Total) bool { return u.File == t.File && u.Sum != "" }); j >= 0 {
			totals[i].Type = totals[j].Type
		}
	}
	return totals
}

func totalOf(p Program, flow dataflow.Flow, field, item string) (Total, bool) {
	r, ok := inputRecord(p, field)
	if !ok || isInputField(p, item) {
		return Total{}, false
	}
	for _, e := range flow.Lineage {
		if e.To != item || e.Verb != "ADD" {
			continue
		}
		switch {
		case e.FromLiteral && e.From == "1":
			return Total{File: r.File, Field: field}, true
		case !e.FromLiteral && e.From != field:
			if summed, found := inputRecord(p, e.From); found && summed.Name == r.Name {
				return Total{File: r.File, Field: field, Sum: e.From, Type: typeHolding(r, e.From)}, true
			}
		}
	}
	return Total{}, false
}

// inputRecord is the record of a file the program reads that a field is
// in.
func inputRecord(p Program, field string) (Record, bool) {
	for _, r := range p.Records {
		if r.File == "" {
			continue
		}
		if _, ok := findField(r.Field, field); ok {
			return r, true
		}
	}
	return Record{}, false
}

func isInputField(p Program, name string) bool {
	_, ok := inputRecord(p, name)
	return ok
}

// typeHolding is the record type whose redefinition holds a field, if any.
func typeHolding(r Record, field string) string {
	for _, t := range r.Types {
		if v, ok := findField(r.Field, t.Fields); ok && t.Fields != "" {
			if _, in := findField(v, field); in {
				return t.Name
			}
		}
	}
	return ""
}

// findField finds a field among a field's fields, at any level.
func findField(f Field, name string) (Field, bool) {
	if f.Name == name {
		return f, true
	}
	for _, child := range f.Fields {
		if found, ok := findField(child, name); ok {
			return found, true
		}
	}
	return Field{}, false
}

// Parameter is an item a called program is passed, and what it does with
// each of its fields: reads it, writes it, or both. A statement on the
// whole item, or a group of it, counts for each field in it.
type Parameter struct {
	Item   string           `json:"item"`
	Fields []ParameterField `json:"fields"`
}

// ParameterField is a field of a parameter and what the program does with
// it.
type ParameterField struct {
	Name   string `json:"name"`
	Reads  bool   `json:"reads,omitempty"`
	Writes bool   `json:"writes,omitempty"`
}

// parameters finds what a called program does with the fields of the
// items it is passed.
func parameters(p Program) []Parameter {
	reads, writes := map[string]bool{}, map[string]bool{}
	for _, para := range p.Paragraphs {
		for _, r := range para.Reads {
			reads[r] = true
		}
		for _, w := range para.Writes {
			writes[w] = true
		}
	}
	var params []Parameter
	for _, using := range p.Using {
		r, ok := recordField(p, using)
		if !ok {
			continue
		}
		param := Parameter{Item: using}
		var visit func(f Field, read, written bool)
		visit = func(f Field, read, written bool) {
			read, written = read || reads[f.Name], written || writes[f.Name]
			if len(f.Fields) == 0 {
				if f.Name != "FILLER" {
					param.Fields = append(param.Fields, ParameterField{Name: f.Name, Reads: read, Writes: written})
				}
				return
			}
			for _, child := range f.Fields {
				visit(child, read, written)
			}
		}
		visit(r, false, false)
		params = append(params, param)
	}
	return params
}

// writeRelations writes the links a program relies on, the control totals
// it checks and, for a called program, what it does with what it is
// passed.
func writeRelations(b *strings.Builder, dir string, p Program) {
	if len(p.Links) > 0 {
		b.WriteString("\n## Links between records\n\nFields of one input the program relies on to find or match the records of another: " +
			"moved into an indexed file's record key to read it by, or compared for equality, each side's values from that field alone.\n\n" +
			"| Field | Of | Links to | Of | How | Source |\n| --- | --- | --- | --- | --- | --- |\n")
		for _, l := range p.Links {
			how := "equal to"
			if l.Kind == "key" {
				how = "key to read by"
			}
			b.WriteString(row(l.From.Field, l.From.File, l.To.Field, l.To.File, how, sourceLink(dir, l.Source, true)))
		}
	}
	if len(p.Totals) > 0 {
		b.WriteString("\n## Control totals\n\nFields of an input the program checks against what it adds up from the input's other records.\n\n" +
			"| Field | Of | Is | Over the records of type | Source |\n| --- | --- | --- | --- | --- |\n")
		for _, t := range p.Totals {
			is := "their count"
			if t.Sum != "" {
				is = "the sum of " + t.Sum
			}
			b.WriteString(row(t.Field, t.File, is, t.Type, sourceLink(dir, t.Source, true)))
		}
	}
	if len(p.Parameters) > 0 {
		b.WriteString("\n## Parameters\n\nWhat the program does with each field of the items it is passed, by position: the contract its " +
			"callers rely on. A statement on a whole item, or a group of it, counts for each field in it; a field neither read nor " +
			"written is passed for nothing.\n")
	}
	for i, param := range p.Parameters {
		fmt.Fprintf(b, "\n### Parameter %d, %s\n\n| Field | Reads | Writes |\n| --- | :-: | :-: |\n", i+1, param.Item)
		for _, f := range param.Fields {
			b.WriteString(row(f.Name, check(f.Reads), check(f.Writes)))
		}
	}
}

func check(b bool) string {
	if b {
		return "yes"
	}
	return ""
}
