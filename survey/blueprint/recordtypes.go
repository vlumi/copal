package blueprint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vlumi/cobolith/cobol"
)

// RecordType is a type of record a program tells apart: the condition name
// on the record's type field that finds it, its value, and the
// redefinition that holds its fields, where the survey can pair one.
type RecordType struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Fields string `json:"fields,omitempty"`
}

// recordTypes finds the types of the program's records: of a record laid
// out in redefinitions, and of a file described by several records, which
// share its record area as redefinitions do. The type field is the first
// field of the record, or of any of the file's records, with a condition
// name the program tests or sets, and its condition names, one for each
// value, are the types. A record with no redefinitions has no types,
// however many states a field of it has. When there are as many
// redefinitions, or records, as types, they are paired in order, as a
// copybook declares them alike. A file's types are its first record's.
func recordTypes(p *cobol.Program, records []Record) {
	tested := testedConditions(p)
	for i := range records {
		r := &records[i]
		var redefinitions, parts []string
		for _, f := range r.Fields {
			if f.Redefines != "" {
				redefinitions = append(redefinitions, f.Name)
			}
			if f.Redefines != "" || len(f.RedefinedBy) > 0 {
				parts = append(parts, f.Name)
			}
		}
		if len(redefinitions) > 0 {
			setTypes(r, tested, r.Fields, parts)
		}
	}
	for _, file := range filesOfSeveralRecords(records) {
		var first *Record
		var fields []Field
		var names []string
		for i := range records {
			if records[i].File == file {
				if first == nil {
					first = &records[i]
				}
				fields = append(fields, records[i].Fields...)
				names = append(names, records[i].Name)
			}
		}
		setTypes(first, tested, fields, names)
	}
}

// testedConditions are the names the program's conditions test, and the
// condition names it sets.
func testedConditions(p *cobol.Program) map[string]bool {
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
	return tested
}

// setTypes gives a record the types the first of the fields with a tested
// condition name tells apart, each paired with the part that holds its
// fields: the part named after the type's condition name, ZG-ﾍｯﾀﾞﾌﾞ for
// ZG-ﾍｯﾀﾞ, or, when there are as many parts as types, the part in the
// same place.
func setTypes(r *Record, tested map[string]bool, fields []Field, parts []string) {
	i := slices.IndexFunc(fields, func(f Field) bool {
		return slices.ContainsFunc(f.Conditions, func(c Condition) bool { return tested[c.Name] })
	})
	if i < 0 {
		return
	}
	r.TypeField = fields[i].Name
	for _, c := range fields[i].Conditions {
		if len(c.Values) > 0 && !slices.ContainsFunc(r.Types, func(t RecordType) bool { return t.Value == c.Values[0] }) {
			r.Types = append(r.Types, RecordType{Name: c.Name, Value: c.Values[0]})
		}
	}
	var pairedTypes, pairedParts []int
	for k, t := range r.Types {
		for i, part := range parts {
			if !slices.Contains(pairedParts, i) && (strings.HasPrefix(part, t.Name) || strings.HasPrefix(t.Name, part)) {
				r.Types[k].Fields = part
				pairedTypes, pairedParts = append(pairedTypes, k), append(pairedParts, i)
				break
			}
		}
	}
	if len(parts)-len(pairedParts) == len(r.Types)-len(pairedTypes) {
		i := 0
		for k := range r.Types {
			if slices.Contains(pairedTypes, k) {
				continue
			}
			for slices.Contains(pairedParts, i) {
				i++
			}
			r.Types[k].Fields = parts[i]
			i++
		}
	}
}

// filesOfSeveralRecords lists the files more than one record describes.
func filesOfSeveralRecords(records []Record) []string {
	count := map[string]int{}
	var files []string
	for _, r := range records {
		if r.File == "" {
			continue
		}
		if count[r.File]++; count[r.File] == 2 {
			files = append(files, r.File)
		}
	}
	return files
}

// typeLeaves lists the elementary fields a record of a type has: of a
// record laid out in redefinitions, those outside them, and those of the
// type's own part, the field the others redefine or one of them, or, for
// a type with none, of the field the others redefine;
// of a file of several records, the type field and the fields of the
// type's record.
func typeLeaves(records []Record, r Record, t RecordType) []Field {
	if i := slices.IndexFunc(records, func(s Record) bool { return s.File == r.File && s.Name == t.Fields }); i >= 0 && r.File != "" {
		out := []Field{}
		if f, ok := recordField(Program{Records: records}, r.TypeField); ok {
			out = append(out, f)
		}
		for _, f := range leaves(records[i].Field) {
			if f.Name != r.TypeField {
				out = append(out, f)
			}
		}
		return out
	}
	var out []Field
	for _, child := range r.Fields {
		redefined := slices.ContainsFunc(r.Fields, func(g Field) bool { return g.Redefines == child.Name })
		switch {
		case child.Redefines != "" && child.Name != t.Fields:
			continue
		case redefined && t.Fields != "" && t.Fields != child.Name:
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
