package blueprint

import (
	"slices"
	"strings"

	"github.com/vlumi/cobolith/layout"
)

// What a field holds, as its layout and the program's files and sorts say.
const (
	KindGroup  = "group"
	KindFiller = "filler"
	KindKey    = "key"
	KindCode   = "code"
	KindEdited = "edited"
	KindRate   = "rate"
	KindAmount = "amount"
	KindDate   = "date"
	KindNumber = "number"
	KindText   = "text"
)

// fieldKind classifies a field by the first of these that holds: a group
// has fields; a filler is named FILLER; a key is a record key of one of
// the program's files or a key one of its sorts sorts on; a code has
// condition names; an edited field has a numeric-edited picture; a rate
// is numeric and named for a rate; an amount is numeric with decimals, or
// without and named for an amount, as a currency without a minor unit
// leaves it; a date is numeric without decimals, of 6, 8 or 14 digits, or
// of 4 for a month or a day without its year, or text of 4 to 8 bytes,
// and named for a date; a number is any other numeric field; and text is
// anything else.
//
// Telling a rate, an amount or a date by its name is a heuristic, and the
// one place the survey reads meaning into names, in English or Japanese:
// a date's name contains DATE, MONTH, TIMESTAMP or YYMM, or ﾋﾂﾞｹ, ﾈﾝｹﾞﾂ,
// ｶﾞｯﾋﾟ, ﾜﾚｷ or ｾｲﾚｷ, or ends in ﾋﾞ, and a four-digit one's YYMM or MMDD,
// ﾈﾝｹﾞﾂ or ｶﾞｯﾋﾟ, or ends in ﾋﾞ; an amount's contains AMOUNT, AMT or
// BALANCE, or ｷﾝｶﾞｸ or ｻﾞﾝﾀﾞｶ, or ends in ｶﾞｸ, ﾘｮｳ or ｷﾝ; a rate's ends in
// RATE, FACTOR, PERCENT or PCT, or ﾘﾂ or ｹｲｽｳ. What the code does with a
// field says more, and kindsByUse carries these kinds along it.
func fieldKind(f *layout.Field, keys map[string]bool) string {
	numeric := f.Category == layout.Numeric
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
	case numeric && namedFor(f.Name, rateWords, rateSuffixes):
		return KindRate
	case numeric && f.Scale > 0:
		return KindAmount
	case numeric && slices.Contains([]int{6, 8, 14}, f.Digits) && namedForDate(f.Name),
		numeric && f.Digits == 4 && namedFor(f.Name, shortDateWords, dateSuffixes):
		return KindDate
	case numeric && namedFor(f.Name, amountWords, amountSuffixes):
		return KindAmount
	case numeric:
		return KindNumber
	case f.Length >= 4 && f.Length <= 8 && namedForDate(f.Name):
		return KindDate
	}
	return KindText
}

// The words and endings that name a date, an amount or a rate, in English
// and in Japanese.
var (
	dateWords      = []string{"DATE", "MONTH", "TIMESTAMP", "YYMM", "ﾋﾂﾞｹ", "ﾈﾝｹﾞﾂ", "ｶﾞｯﾋﾟ", "ﾜﾚｷ", "ｾｲﾚｷ"}
	shortDateWords = []string{"YYMM", "MMDD", "ﾈﾝｹﾞﾂ", "ｶﾞｯﾋﾟ"}
	dateSuffixes   = []string{"ﾋﾞ"}
	amountWords    = []string{"AMOUNT", "AMT", "BALANCE", "ｷﾝｶﾞｸ", "ｻﾞﾝﾀﾞｶ"}
	amountSuffixes = []string{"ｶﾞｸ", "ﾘｮｳ", "ｷﾝ"}
	rateWords      []string
	rateSuffixes   = []string{"RATE", "FACTOR", "PERCENT", "PCT", "ﾘﾂ", "ｹｲｽｳ"}
)

func namedForDate(name string) bool {
	return namedFor(name, dateWords, dateSuffixes)
}

// namedFor reports whether a name contains one of the words or ends in one
// of the suffixes.
func namedFor(name string, words, suffixes []string) bool {
	return slices.ContainsFunc(words, func(word string) bool { return strings.Contains(name, word) }) ||
		slices.ContainsFunc(suffixes, func(suffix string) bool { return strings.HasSuffix(name, suffix) })
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
