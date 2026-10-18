package blueprint

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/survey/internal/dataflow"
)

// Output is a record a program writes to a file, of one type if the
// program tells its types apart, and where the values of each of its
// fields come from. A file of several records, or of a record of several
// types, a header, details and a trailer, has an Output for each.
type Output struct {
	File   string        `json:"file"`
	Record string        `json:"record"`
	Type   string        `json:"type,omitempty"`
	Fields []FieldOrigin `json:"fields"`
}

// FieldOrigin is a field written to a file and its origins: the input
// fields, literals and initial values its values come from.
type FieldOrigin struct {
	Name    string   `json:"name"`
	Start   int      `json:"start"`
	Length  int      `json:"length"`
	Origins []string `json:"origins,omitempty"`
	origins []dataflow.Origin
}

// outputs traces the fields of every record the program writes to a file,
// once for each file and record, as the statements that write it use
// them; for a record of types, once for each type, as the statements that
// write a record of it use them, with the fields of that type. A record
// written FROM another is described by that one's fields.
func outputs(p Program, flow dataflow.Flow) []Output {
	var keys [][2]string
	writes := map[[2]string][]cobol.Pos{}
	for _, w := range flow.Writes() {
		key := [2]string{w.File, w.Record}
		if _, ok := writes[key]; !ok && slices.ContainsFunc(p.Records, func(r Record) bool { return r.Name == w.Record }) {
			keys = append(keys, key)
		}
		writes[key] = append(writes[key], w.Pos)
	}
	var outs []Output
	for _, key := range keys {
		r := p.Records[slices.IndexFunc(p.Records, func(r Record) bool { return r.Name == key[1] })]
		if len(r.Types) == 0 {
			outs = append(outs, traced(flow, Output{File: key[0], Record: key[1]}, leaves(r.Field), writes[key]))
			continue
		}
		byType := make([][]cobol.Pos, len(r.Types))
		for _, pos := range writes[key] {
			for _, k := range writtenTypes(r, flow, pos) {
				byType[k] = append(byType[k], pos)
			}
		}
		for k, t := range r.Types {
			if len(byType[k]) > 0 {
				outs = append(outs, traced(flow, Output{File: key[0], Record: key[1], Type: t.Name}, typeLeaves(r, t), byType[k]))
			}
		}
	}
	return outs
}

// traced is an output with the origins of each of its fields, as the
// statements at the positions write them.
func traced(flow dataflow.Flow, out Output, fields []Field, positions []cobol.Pos) Output {
	for _, f := range fields {
		fo := FieldOrigin{Name: f.Name, Start: f.Start, Length: f.Length}
		fo.origins = flow.TraceAt(f.Name, positions).Origins
		for _, o := range fo.origins {
			fo.Origins = append(fo.Origins, o.String())
		}
		out.Fields = append(out.Fields, fo)
	}
	return out
}

// writtenTypes are the types of record a statement may write: those its
// type field's values come from on the paths to it, a condition name set
// or a value moved; every type, when one comes from anything else.
func writtenTypes(r Record, flow dataflow.Flow, pos cobol.Pos) []int {
	all := make([]int, len(r.Types))
	for k := range all {
		all[k] = k
	}
	origins := flow.TraceAt(r.TypeField, []cobol.Pos{pos}).Origins
	if len(origins) == 0 {
		return all
	}
	var types []int
	for _, o := range origins {
		name, isName := strings.CutPrefix(o.Literal, "condition ")
		k := slices.IndexFunc(r.Types, func(t RecordType) bool { return isName && t.Name == name || !isName && t.Value == o.Literal })
		if o.Literal == "" || k < 0 {
			return all
		}
		if !slices.Contains(types, k) {
			types = append(types, k)
		}
	}
	slices.Sort(types)
	return types
}

// leaves lists a record's elementary fields, FILLER left out, or the record
// itself if it has none.
func leaves(f Field) []Field {
	if len(f.Fields) == 0 {
		return []Field{f}
	}
	var out []Field
	for _, child := range f.Fields {
		for _, leaf := range leaves(child) {
			if leaf.Name != "FILLER" {
				out = append(out, leaf)
			}
		}
	}
	return out
}

func writeOrigins(b *strings.Builder, p Program) {
	if len(p.Outputs) == 0 {
		return
	}
	b.WriteString("\n## Field origins\n\nWhere the values of each field the program writes to a file come from: " +
		"the fields of the files it reads, literals, and initial values, followed through every statement that moves a value, " +
		"and through groups and redefinitions by the bytes they share.\n")
	for _, file := range outputFiles(p.Outputs) {
		fmt.Fprintf(b, "\n### File %s\n\nWritten from %s.\n", file[0].File, recordLinks(file))
		for _, o := range file {
			switch {
			case o.Type != "":
				fmt.Fprintf(b, "\n#### Record %s, of type %s\n", o.Record, o.Type)
			case len(file) > 1:
				fmt.Fprintf(b, "\n#### Record %s\n", o.Record)
			}
			writeFieldOrigins(b, o.Fields)
		}
	}
}

// outputFiles groups the records written by the file they are written to,
// in the order the files are first written.
func outputFiles(outs []Output) [][]Output {
	var files [][]Output
	for _, o := range outs {
		i := slices.IndexFunc(files, func(file []Output) bool { return file[0].File == o.File })
		if i < 0 {
			files = append(files, []Output{o})
			continue
		}
		files[i] = append(files[i], o)
	}
	return files
}

func recordLinks(outs []Output) string {
	var links []string
	for _, o := range outs {
		if link := fmt.Sprintf("[%s](#%s)", o.Record, HeadingAnchor(o.Record)); !slices.Contains(links, link) {
			links = append(links, link)
		}
	}
	if len(links) == 1 {
		return links[0]
	}
	return strings.Join(links[:len(links)-1], ", ") + " and " + links[len(links)-1]
}

func writeFieldOrigins(b *strings.Builder, fields []FieldOrigin) {
	b.WriteString("\n| Field | Position | Origins |\n| --- | --: | --- |\n")
	for _, f := range fields {
		position := strconv.Itoa(f.Start)
		if f.Length > 1 {
			position += "-" + strconv.Itoa(f.Start+f.Length-1)
		}
		b.WriteString(row(f.Name, position, strings.Join(f.Origins, "; ")))
	}
}
