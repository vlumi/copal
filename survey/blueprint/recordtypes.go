package blueprint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
)

// RecordType is a type of record a program tells apart: the condition name
// on the record's type field that finds it, its value, and the
// redefinition that holds its fields, where the survey can pair one.
type RecordType struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Fields string `json:"fields,omitempty"`
}

// recordTypes finds the types of a record laid out in redefinitions: the
// first field of it with a condition name the program tests or sets, whose
// condition names, one for each value, are the types. A record with no
// redefinitions has no types, however many states a field of it has. When
// there are as many redefinitions as types, they are paired in order, as a
// copybook declares them alike.
func recordTypes(p *cobol.Program, r *Record) {
	if !slices.ContainsFunc(r.Fields, func(f Field) bool { return f.Redefines != "" }) {
		return
	}
	tested := map[string]bool{}
	for _, para := range p.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if s.Verb == "SET" && len(s.Reads) == 0 {
				for _, w := range s.Writes {
					tested[w.Name] = true
				}
			}
			conditions := []*cobol.Condition{s.Condition}
			for _, b := range s.Branches {
				conditions = append(conditions, b.Condition)
			}
			for _, c := range conditions {
				if c != nil {
					for _, ref := range c.Refs() {
						tested[ref] = true
					}
				}
			}
		})
	}
	i := slices.IndexFunc(r.Fields, func(f Field) bool {
		return slices.ContainsFunc(f.Conditions, func(c Condition) bool { return tested[c.Name] })
	})
	if i < 0 {
		return
	}
	field := r.Fields[i]
	r.TypeField = field.Name
	for _, c := range field.Conditions {
		if len(c.Values) > 0 && !slices.ContainsFunc(r.Types, func(t RecordType) bool { return t.Value == c.Values[0] }) {
			r.Types = append(r.Types, RecordType{Name: c.Name, Value: c.Values[0]})
		}
	}
	var redefinitions []string
	for _, f := range r.Fields {
		if f.Redefines != "" {
			redefinitions = append(redefinitions, f.Name)
		}
	}
	if len(redefinitions) == len(r.Types) {
		for k := range r.Types {
			r.Types[k].Fields = redefinitions[k]
		}
	}
}

// typeLeaves lists a record's elementary fields a record of a type has:
// those outside its redefinitions, and those of the type's own, or, for a
// type with none, of the field the others redefine.
func typeLeaves(r Record, t RecordType) []Field {
	var out []Field
	for _, child := range r.Fields {
		redefined := slices.ContainsFunc(r.Fields, func(g Field) bool { return g.Redefines == child.Name })
		switch {
		case child.Redefines != "" && child.Name != t.Fields:
			continue
		case redefined && t.Fields != "":
			continue
		}
		out = append(out, leaves(Field{Fields: []Field{child}})...)
	}
	return out
}

// typesSentence writes a record's types, if it has any.
func typesSentence(r Record) string {
	if len(r.Types) == 0 {
		return ""
	}
	var types []string
	for _, t := range r.Types {
		s := fmt.Sprintf("`%s`, %s", t.Name, t.Value)
		if t.Fields != "" {
			s += fmt.Sprintf(", its fields in %s", t.Fields)
		}
		types = append(types, s)
	}
	return fmt.Sprintf(" Its types, told apart by %s: %s.", r.TypeField, strings.Join(types, "; "))
}
