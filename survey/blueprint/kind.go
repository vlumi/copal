package blueprint

import (
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/layout"
)

// What a field holds, as its layout and the program's files and sorts say.
const (
	KindGroup  = "group"
	KindFiller = "filler"
	KindKey    = "key"
	KindCode   = "code"
	KindEdited = "edited"
	KindAmount = "amount"
	KindDate   = "date"
	KindNumber = "number"
	KindText   = "text"
)

// fieldKind classifies a field by the first of these that holds: a group
// has fields; a filler is named FILLER; a key is a record key of one of
// the program's files or a key one of its sorts sorts on; a code has
// condition names; an edited field has a numeric-edited picture; an
// amount is numeric with decimals; a date is numeric without decimals, of
// 6, 8 or 14 digits, and named for a date; a number is any other numeric
// field; and text is anything else.
//
// Telling a date by its name is a heuristic: a field is a date if its name
// contains DATE, MONTH, TIMESTAMP or YYMM. Nothing else in the survey
// reads meaning into names.
func fieldKind(f *layout.Field, keys map[string]bool) string {
	switch {
	case len(f.Children) > 0:
		return KindGroup
	case f.Name == "FILLER":
		return KindFiller
	case keys[f.Name]:
		return KindKey
	case len(f.Conditions) > 0:
		return KindCode
	case f.Category == layout.NumericEdited:
		return KindEdited
	case f.Category == layout.Numeric && f.Scale > 0:
		return KindAmount
	case f.Category == layout.Numeric && f.Scale == 0 && slices.Contains([]int{6, 8, 14}, f.Digits) && namedForDate(f.Name):
		return KindDate
	case f.Category == layout.Numeric:
		return KindNumber
	}
	return KindText
}

func namedForDate(name string) bool {
	return slices.ContainsFunc([]string{"DATE", "MONTH", "TIMESTAMP", "YYMM"}, func(word string) bool {
		return strings.Contains(name, word)
	})
}

// Program kinds.
const (
	KindMain   = "main"
	KindCalled = "called"
)

// programKind is called for a program with a USING list, main otherwise.
func programKind(using []string) string {
	if len(using) > 0 {
		return KindCalled
	}
	return KindMain
}

// programKeys are the fields a program names as keys: its files' record
// keys and its sorts' keys.
func programKeys(files []File, sorts []Sort) map[string]bool {
	keys := map[string]bool{}
	for _, f := range files {
		if f.RecordKey != "" {
			keys[f.RecordKey] = true
		}
	}
	for _, s := range sorts {
		for _, k := range s.Keys {
			keys[k] = true
		}
	}
	return keys
}

// copybookKeys are the fields of a copybook that any program using it
// names as keys.
func copybookKeys(copybook string, programs []Program) map[string]bool {
	keys := map[string]bool{}
	for _, p := range programs {
		for _, r := range p.Records {
			if r.Copybook == copybook {
				walk(r.Field, func(f Field) {
					if f.Kind == KindKey {
						keys[f.Name] = true
					}
				})
			}
		}
	}
	return keys
}

// walk visits a field and every field under it, in order.
func walk(f Field, visit func(Field)) {
	visit(f)
	for _, child := range f.Fields {
		walk(child, visit)
	}
}
