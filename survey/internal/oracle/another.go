package oracle

import (
	"bytes"
	"slices"

	"github.com/vlumi/cobolith/layout"
)

// another makes sure the last record of an input holds another value in a
// field than its first, and reports whether it does: if they hold the
// same, it adds a record unlike the others, a second product for a rider
// to name beside the one the policy is of.
func (p *program) another(c Case, t target) bool {
	in := p.input(t.input)
	records := c[t.input]
	if len(records) == 0 || t.part != (part{}) {
		return false
	}
	last := records[len(records)-1]
	if !bytes.Equal(bytesOf(in, records[0], t.field), bytesOf(in, last, t.field)) {
		return true
	}
	more := slices.Clone(last)
	if !p.next(in, more, t.field) {
		return false
	}
	p.vary(in, more, t.field)
	if !p.distinct(in, records, more) {
		return false
	}
	c[t.input] = append(records, more)
	return true
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
