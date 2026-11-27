package oracle

import (
	"bytes"
	"slices"

	"github.com/vlumi/cobolith/layout"
)

// another finds the first record of an input that holds another value
// in a field than the first that has the field: a second product for a
// rider to name beside the one the policy is of, a second policy for a
// payment. Where they all hold the same, it adds one unlike the others
// after the last of them. It gives the record's place among the input's.
func (p *program) another(c Case, t target) (int, bool) {
	in, s := p.input(t.input), p.shapes[t.input]
	records := c[t.input]
	if t.part != (part{}) {
		return 0, false
	}
	first, last := -1, -1
	for i, r := range records {
		if !s.has(in, r, t.field) {
			continue
		}
		if first < 0 {
			first = i
		}
		if last = i; !bytes.Equal(bytesOf(in, records[first], t.field), bytesOf(in, r, t.field)) {
			return i, true
		}
	}
	if first < 0 {
		return 0, false
	}
	more := slices.Clone(records[last])
	if !p.next(in, more, t.field) {
		return 0, false
	}
	p.vary(in, more, t.field)
	if !p.distinct(in, records, more) {
		return 0, false
	}
	c[t.input] = slices.Insert(records, last+1, more)
	return last + 1, true
}

func bytesOf(in input, r []byte, f *layout.Field) []byte {
	return r[f.Start-in.root.Start:][:f.Size()]
}

// next puts the value after its own into a field of a record: a number
// one more, or text with its last letter or digit the one after.
func (p *program) next(in input, r []byte, f *layout.Field) bool {
	if numeric(f) {
		n, err := in.layout.Number(r, f.Name)
		return err == nil && set(r, in.root, f, in.layout, number(n+1))
	}
	held := bytes.TrimRight(bytesOf(in, r, f), " ")
	if len(held) == 0 {
		return false
	}
	switch b := held[len(held)-1]; {
	case b >= '0' && b < '9', b >= 'A' && b < 'Z':
		held[len(held)-1]++
		return true
	}
	return false
}

// vary gives a record another of the named values of each field that has
// them, the first it does not hold, so that it is unlike the record it
// was copied from in what the program tells its records apart by. The
// field that says what type the record is of, the fields of its key and
// the one field excepted keep their values.
func (p *program) vary(in input, r []byte, except *layout.Field) {
	s := p.shapes[in.name()]
	key := in.key()
	walk(in.root, func(f *layout.Field) {
		if f == except || f == s.discriminator || len(f.Conditions) == 0 || !s.has(in, r, f) || key != nil && find(key, f.Name) == f {
			return
		}
		for _, cond := range f.Conditions {
			v, ok := literal(f, cond.Values[0])
			probe := slices.Clone(r)
			if ok && set(probe, in.root, f, in.layout, v) && !bytes.Equal(probe, r) {
				copy(r, probe)
				return
			}
		}
	})
}

// key is the record key of an indexed input, or nil.
func (in input) key() *layout.Field {
	if in.file.RecordKey == "" {
		return nil
	}
	return find(in.root, in.file.RecordKey)
}

// distinct makes sure a record to add to an indexed input's has a key none
// of them has, taking the next value of the last field of a key of
// several, and reports whether it has.
func (p *program) distinct(in input, records [][]byte, r []byte) bool {
	key := in.key()
	if key == nil {
		return true
	}
	taken := func() bool {
		return slices.ContainsFunc(records, func(o []byte) bool { return bytes.Equal(bytesOf(in, o, key), bytesOf(in, r, key)) })
	}
	if !taken() {
		return true
	}
	last, ok := in.lastOfKey()
	return ok && p.next(in, r, last) && !taken()
}
