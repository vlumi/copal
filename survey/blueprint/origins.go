package blueprint

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/survey/internal/dataflow"
)

// Output is a file a program writes, and where the values of each field of
// what it writes come from.
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

// outputs traces the fields of every record the program writes to a file.
// A record written FROM another is described by that one's fields.
func outputs(p Program, flow dataflow.Flow) []Output {
	var outs []Output
	for _, e := range p.Lineage {
		if !e.ToFile || slices.ContainsFunc(outs, func(o Output) bool { return o.File == e.To }) {
			continue
		}
		record := e.From
		for _, from := range p.Lineage {
			if (from.Verb == "WRITE" || from.Verb == "REWRITE" || from.Verb == "RELEASE") && from.To == e.From && !from.FromLiteral {
				record = from.From
			}
		}
		i := slices.IndexFunc(p.Records, func(r Record) bool { return r.Name == record })
		if i < 0 {
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
	for _, o := range p.Outputs {
		fmt.Fprintf(b, "\n### File %s\n\nWritten from [%s](#%s).\n\n| Field | Position | Origins |\n| --- | --: | --- |\n",
			o.File, o.Record, headingAnchor(o.Record))
		for _, f := range o.Fields {
			position := strconv.Itoa(f.Start)
			if f.Length > 1 {
				position += "-" + strconv.Itoa(f.Start+f.Length-1)
			}
			b.WriteString(row(f.Name, position, strings.Join(f.Origins, "; ")))
		}
	}
}
