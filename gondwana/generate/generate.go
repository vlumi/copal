// Package generate makes up Gondwana Life Insurance's portfolio and its
// days: the policies the system opens with, the tables the programs load,
// and for each business day what the sales system sold, what the branches
// requested and, on the day the collection agent answers, how each debit
// went. Everything follows from a seed, so the same seed and date always
// give the same records. The text is Shift_JIS, as the programs read it.
package generate

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/vlumi/cobolith/record"
	"github.com/vlumi/copal/gondwana/records"
	"golang.org/x/text/encoding/japanese"
)

// Generator makes up the insurer's records for one seed.
type Generator struct {
	Layouts *records.Layouts
	Seed    uint64
}

// rng returns a random source for one purpose, so that changing how one
// thing is generated does not change everything else.
func (g *Generator) rng(stream string, n int64) *rand.Rand {
	h := uint64(14695981039346656037)
	for _, c := range []byte(stream) {
		h = (h ^ uint64(c)) * 1099511628211
	}
	return rand.New(rand.NewPCG(g.Seed, h^uint64(n)))
}

// sjis encodes text in Shift_JIS, as every file of the insurer holds it.
// The text comes from this package's own tables, so a character the
// encoding lacks is a mistake in them.
func sjis(text string) string {
	encoded, err := japanese.ShiftJIS.NewEncoder().String(text)
	if err != nil {
		panic(fmt.Sprintf("%q: %v", text, err))
	}
	return encoded
}

// Records splits a line sequential file into records of a layout's
// length: a record a line, the trailing spaces the write dropped put back.
func Records(l *record.Layout, file []byte) [][]byte {
	var recs [][]byte
	for _, line := range bytes.Split(bytes.TrimRight(file, "\n"), []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		r := l.Blank()
		copy(r, line)
		recs = append(recs, r)
	}
	return recs
}

// Fixed splits a fixed-length file, one with no record delimiters, into
// records of a layout's length.
func Fixed(l *record.Layout, file []byte) [][]byte {
	var recs [][]byte
	for i := 0; i+l.Length <= len(file); i += l.Length {
		recs = append(recs, slices.Clone(file[i:i+l.Length]))
	}
	return recs
}

// File joins records into a line sequential file, a record a line, with
// the trailing spaces dropped as the programs drop them.
func File(records [][]byte) []byte {
	var b bytes.Buffer
	for _, r := range records {
		b.Write(bytes.TrimRight(r, " "))
		b.WriteByte('\n')
	}
	return b.Bytes()
}

// FixedFile joins records into a fixed-length file with no delimiters.
func FixedFile(records [][]byte) []byte {
	return bytes.Join(records, nil)
}

// get reads fields from a record into the pointers given for them,
// numbers and text alike.
func get(l *record.Layout, r []byte, fields map[string]any) error {
	for name, target := range fields {
		var err error
		switch v := target.(type) {
		case *int64:
			*v, err = l.Number(r, name)
		case *int:
			var n int64
			n, err = l.Number(r, name)
			*v = int(n)
		case *string:
			*v, err = l.Text(r, name)
		default:
			err = fmt.Errorf("%s: %T is neither a number nor text", name, target)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// set writes fields into a record, numbers and text alike.
func set(l *record.Layout, r []byte, fields map[string]any) error {
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		var err error
		switch v := fields[name].(type) {
		case int64:
			err = l.SetNumber(r, name, v)
		case int:
			err = l.SetNumber(r, name, int64(v))
		case string:
			err = l.SetText(r, name, v)
		default:
			err = fmt.Errorf("%s: %T is neither a number nor text", name, v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// made builds a record of a layout from its fields.
func made(l *record.Layout, fields map[string]any) ([]byte, error) {
	r := l.Blank()
	return r, set(l, r, fields)
}

// pick returns an index drawn by weight.
func pick(r *rand.Rand, weights []int) int {
	total := 0
	for _, w := range weights {
		total += w
	}
	n := r.IntN(total)
	for i, w := range weights {
		if n < w {
			return i
		}
		n -= w
	}
	return len(weights) - 1
}

// chance reports whether an event of the given probability happens.
func chance(r *rand.Rand, probability float64) bool {
	return r.Float64() < probability
}
