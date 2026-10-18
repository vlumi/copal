package oracle

import (
	"math"
	"strconv"
	"strings"

	"github.com/vlumi/copal/bedrock/layout"
	"github.com/vlumi/copal/bedrock/record"
)

// value is what goes into a field: a number, in units of the field's
// scale, text, or one byte repeated, for HIGH-VALUES and LOW-VALUES.
type value struct {
	number int64
	text   string
	fill   byte
	kind   valueKind
}

type valueKind int

const (
	numberValue valueKind = iota
	textValue
	fillValue
)

func number(n int64) value { return value{number: n, kind: numberValue} }
func text(s string) value  { return value{text: s, kind: textValue} }
func filled(b byte) value  { return value{fill: b, kind: fillValue} }
func (v value) String() string {
	switch v.kind {
	case numberValue:
		return strconv.FormatInt(v.number, 10)
	case textValue:
		return strconv.Quote(v.text)
	}
	return "x" + strconv.FormatUint(uint64(v.fill), 16)
}

// literal reads a literal as a value for a field: a figurative constant,
// a quoted string, or a number, scaled to the field.
func literal(f *layout.Field, lit string) (value, bool) {
	switch lit {
	case "SPACE", "SPACES":
		return text(""), true
	case "ZERO", "ZEROS", "ZEROES":
		if numeric(f) {
			return number(0), true
		}
		return text(strings.Repeat("0", f.Length)), true
	case "HIGH-VALUE", "HIGH-VALUES":
		return filled(0xff), true
	case "LOW-VALUE", "LOW-VALUES":
		return filled(0x00), true
	}
	if strings.HasPrefix(lit, "'") || strings.HasPrefix(lit, `"`) {
		s := lit[1 : len(lit)-1]
		if numeric(f) {
			return scaled(f, s)
		}
		return text(s), true
	}
	if !numeric(f) {
		return text(lit), true
	}
	return scaled(f, lit)
}

// scaled reads a number as a value in units of the field's scale.
func scaled(f *layout.Field, s string) (value, bool) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return value{}, false
	}
	return number(int64(math.Round(n * math.Pow10(f.Scale)))), true
}

func numeric(f *layout.Field) bool { return f.Category == layout.Numeric }

// largest is the largest number a numeric field holds, in units of its
// scale.
func largest(f *layout.Field) int64 {
	return int64(math.Pow10(min(f.Digits, 18))) - 1
}

// set writes a value into a field of a record. A number that does not fit
// the field, or is negative in an unsigned one, is left out.
func set(r []byte, root, f *layout.Field, l *record.Layout, v value) bool {
	start := f.Start - root.Start
	switch {
	case v.kind == fillValue:
		for i := range f.Size() {
			r[start+i] = v.fill
		}
		return true
	case v.kind == textValue && !numeric(f):
		return l.SetText(r, f.Name, v.text) == nil
	case v.kind == textValue:
		copy(r[start:start+f.Size()], []byte((v.text + strings.Repeat(" ", f.Size()))[:f.Size()]))
		return true
	case !numeric(f):
		return l.SetText(r, f.Name, strconv.FormatInt(v.number, 10)) == nil
	case v.number < 0 && !f.Signed, v.number > largest(f), -v.number > largest(f):
		return false
	}
	return l.SetNumber(r, f.Name, v.number) == nil
}
