package blueprint

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/survey/internal/dataflow"
)

// Output is a record a program writes to a file, and where the values of
// each of its fields come from. A file of several record types, a header,
// details and a trailer, has an Output for each.
type Output struct {
	File   string        `json:"file"`
	Record string        `json:"record"`
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
// them. A record written FROM another is described by that one's fields.
func outputs(p Program, flow dataflow.Flow) []Output {
	var outs []Output
	writes := map[[2]string][]cobol.Pos{}
	for _, w := range flow.Writes() {
		key := [2]string{w.File, w.Record}
		if _, ok := writes[key]; !ok && slices.ContainsFunc(p.Records, func(r Record) bool { return r.Name == w.Record }) {
			outs = append(outs, Output{File: w.File, Record: w.Record})
		}
		writes[key] = append(writes[key], w.Pos)
	}
	for i, out := range outs {
		r := p.Records[slices.IndexFunc(p.Records, func(r Record) bool { return r.Name == out.Record })]
		for _, f := range leaves(r.Field) {
			fo := FieldOrigin{Name: f.Name, Start: f.Start, Length: f.Length}
			fo.origins = flow.TraceAt(f.Name, writes[[2]string{out.File, out.Record}]).Origins
			for _, o := range fo.origins {
				fo.Origins = append(fo.Origins, o.String())
			}
			outs[i].Fields = append(outs[i].Fields, fo)
		}
	}
	return outs
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
			if len(file) > 1 {
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
	links := make([]string, len(outs))
	for i, o := range outs {
		links[i] = fmt.Sprintf("[%s](#%s)", o.Record, HeadingAnchor(o.Record))
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
