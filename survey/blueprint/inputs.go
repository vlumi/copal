package blueprint

import (
	"slices"

	"github.com/vlumi/copal/survey/internal/dataflow"
)

// Input is a record a program reads from a file, of one type if the
// program tells its types apart, and where the values of each of its
// fields go: the mirror of Output.
type Input struct {
	File   string        `json:"file"`
	Record string        `json:"record"`
	Type   string        `json:"type,omitempty"`
	Fields []FieldEffect `json:"fields"`
}

// FieldEffect is a field read from a file and what its values reach: the
// data items they are moved into and the files they are written to.
type FieldEffect struct {
	Name    string   `json:"name"`
	Start   int      `json:"start"`
	Length  int      `json:"length"`
	Reaches []string `json:"reaches,omitempty"`
}

// Affects follows where a data item's values go within the program,
// through the statements that move them.
func (p Program) Affects(item string) dataflow.Effect {
	return p.flow.Affects(item)
}

// AffectsBytes follows where the values of some of a data item's bytes
// go, from start, 1-based, for length.
func (p Program) AffectsBytes(item string, start, length int) dataflow.Effect {
	return p.flow.AffectsBytes(item, start, length)
}

// Share reports whether two data items of the program share any bytes.
func (p Program) Share(a, b string) bool {
	return p.flow.Share(a, b)
}

// inputs follows the fields of every record the program reads from a
// file, as outputs traces those it writes, for a record of types, those of
// each type. A record read INTO another is described by that one's fields.
func inputs(p Program, flow dataflow.Flow) []Input {
	var ins []Input
	for _, e := range p.Lineage {
		if !e.FromFile || slices.ContainsFunc(ins, func(in Input) bool { return in.File == e.From }) {
			continue
		}
		record := e.To
		for _, into := range p.Lineage {
			if (into.Verb == "READ" || into.Verb == "RETURN") && into.From == e.To && !into.FromFile {
				record = into.To
			}
		}
		i := slices.IndexFunc(p.Records, func(r Record) bool { return r.Name == record })
		if i < 0 {
			continue
		}
		r := p.Records[i]
		if len(r.Types) == 0 {
			ins = append(ins, affected(flow, Input{File: e.From, Record: record}, leaves(r.Field)))
			continue
		}
		for _, t := range r.Types {
			ins = append(ins, affected(flow, Input{File: e.From, Record: record, Type: t.Name}, typeLeaves(p.Records, r, t)))
		}
	}
	return ins
}

// affected is an input with where the values of each of its fields go.
func affected(flow dataflow.Flow, in Input, fields []Field) Input {
	for _, f := range fields {
		fe := FieldEffect{Name: f.Name, Start: f.Start, Length: f.Length}
		effect := flow.Affects(f.Name)
		fe.Reaches = append(fe.Reaches, effect.Items...)
		for _, w := range effect.Files {
			addOnce(&fe.Reaches, "file "+w.File)
		}
		in.Fields = append(in.Fields, fe)
	}
	return in
}
