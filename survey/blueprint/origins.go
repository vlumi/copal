package blueprint

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

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
}

// outputs traces the fields of every record the program writes to a file,
// once for each file and record. A record written FROM another is
// described by that one's fields.
func outputs(p Program, flow dataflow.Flow) []Output {
	var outs []Output
	for _, e := range p.Lineage {
		if !e.ToFile {
			continue
		}
		for _, record := range writtenFrom(p, e) {
			i := slices.IndexFunc(p.Records, func(r Record) bool { return r.Name == record })
			if i < 0 || slices.ContainsFunc(outs, func(o Output) bool { return o.File == e.To && o.Record == record }) {
				continue
			}
			out := Output{File: e.To, Record: record}
			for _, f := range leaves(p.Records[i].Field) {
				fo := FieldOrigin{Name: f.Name, Start: f.Start, Length: f.Length}
				for _, o := range flow.Origins(f.Name) {
					fo.Origins = append(fo.Origins, o.String())
				}
				out.Fields = append(out.Fields, fo)
			}
			outs = append(outs, out)
		}
	}
	return outs
}

// writtenFrom is the record a statement writes to a file: the record named
// in its FROM, or the file's own record.
func writtenFrom(p Program, write Lineage) []string {
	var from []string
	for _, e := range p.Lineage {
		if e.Source == write.Source && e.To == write.From && e.FromItem() {
			from = append(from, e.From)
		}
	}
	if len(from) == 0 {
		return []string{write.From}
	}
	return from
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
		links[i] = fmt.Sprintf("[%s](#%s)", o.Record, headingAnchor(o.Record))
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
