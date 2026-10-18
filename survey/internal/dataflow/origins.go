package dataflow

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/layout"
)

// Origin is where a data item's values come from in the end: a field of a
// record a file is read into, a literal, an item's initial VALUE, a
// program called, which may change an item passed to it by reference, the
// caller of the program, which passes the items it is called with, or the
// job's SYSIN, which an ACCEPT reads.
type Origin struct {
	File      string // the file, for a field read from one
	Item      string // the item at the bytes, or the item with the initial value
	Literal   string // the literal, as written
	Call      string // the program called, for an item passed to it
	Arg       int    // the position of the argument passed, 1-based
	Parameter int    // the position of the USING item passed in, 1-based
	Stream    string // SYSIN, for an item accepted from it
}

// String writes the origin: "file TRAN-IN TD-AMOUNT", "literal 'TXPOST'",
// "initial value of WS-COUNT", "call AUTHPUR argument 6 HD-AMOUNT",
// "parameter 1 RQ-AMOUNT", "stream SYSIN AUTH-REQUEST".
func (o Origin) String() string {
	switch {
	case o.File != "":
		return "file " + o.File + " " + o.Item
	case o.Literal != "":
		return "literal " + o.Literal
	case o.Call != "":
		return fmt.Sprintf("call %s argument %d %s", o.Call, o.Arg, o.Item)
	case o.Parameter > 0:
		return fmt.Sprintf("parameter %d %s", o.Parameter, o.Item)
	case o.Stream != "":
		return "stream " + o.Stream + " " + o.Item
	}
	return "initial value of " + o.Item
}

// place is where an item's bytes are: in which record, from where, how
// many; and whether it has an initial VALUE and fields of its own.
type place struct {
	record        string
	start, length int
	value, group  bool
}

func (p place) overlaps(q place) bool {
	return p.record == q.record && p.start < q.start+q.length && q.start < p.start+p.length
}

// storage is where every item the program declares lives, so that values
// can be followed by bytes, not just by name: into a group, out of a field
// of it, through a field that redefines another.
type storage struct {
	places map[string]place
	items  map[string][]string // each record's items, in declaration order
}

func newStorage(records *layout.ProgramRecords) *storage {
	s := &storage{places: map[string]place{}, items: map[string][]string{}}
	seen := map[string]int{}
	var walk func(record string, f *layout.Field)
	walk = func(record string, f *layout.Field) {
		seen[f.Name]++
		s.places[f.Name] = place{record: record, start: f.Start, length: f.Size(), value: len(f.Values) > 0, group: len(f.Children) > 0}
		s.items[record] = append(s.items[record], f.Name)
		for _, child := range f.Children {
			walk(record, child)
		}
	}
	for _, r := range records.All() {
		walk(r.Name, r)
	}
	for name, n := range seen {
		if n > 1 {
			delete(s.places, name)
		}
	}
	return s
}

// part is the place of a run of an item's bytes, from start, 1-based, for
// length.
func (s *storage) part(item string, start, length int) (place, bool) {
	p, ok := s.places[item]
	if !ok || start < 1 || start+length-1 > p.length {
		return place{}, false
	}
	return place{record: p.record, start: p.start + start - 1, length: length, group: p.group}, true
}

// at is the narrowest item that holds a record's bytes from start, for
// length; the record itself if no field does.
func (s *storage) at(record string, start, length int) string {
	best, size := record, s.places[record].length
	for _, name := range s.items[record] {
		p, ok := s.places[name]
		if ok && p.start <= start && start+length <= p.start+p.length && p.length < size {
			best, size = name, p.length
		}
	}
	return best
}

// Trace is where a data item's values come from: the origins at the
// start, and the statements that move them on to the item, as lineage
// edges, in the order the program declares them.
type Trace struct {
	Origins []Origin
	Edges   []Edge
}

// Origins is where a data item's values come from in the end, as Trace
// finds them.
func (f Flow) Origins(item string) []Origin {
	return f.Trace(item).Origins
}

// Trace traces where a data item's values come from: through every
// statement that moves a value into it or into bytes it shares, a group it
// is in, a field of it or a field that redefines it, and on through what
// those statements move from, to the files, literals and initial values
// at the start. A group moved into a group, or anything moved into an item
// of its own length, is followed byte for byte. A literal moved into a
// whole group around the item, such as MOVE SPACES TO A-RECORD, only
// clears it, and counts only for an item with no other origin.
func (f Flow) Trace(item string) Trace {
	start, ok := f.storage.places[item]
	if !ok {
		return Trace{}
	}
	return f.trace(start)
}

// TraceBytes traces where the values of some of a data item's bytes come
// from, from start, 1-based, for length, as Trace does.
func (f Flow) TraceBytes(item string, start, length int) Trace {
	p, ok := f.storage.part(item, start, length)
	if !ok {
		return Trace{}
	}
	return f.trace(p)
}

func (f Flow) trace(start place) Trace {
	found, clearing := map[Origin]bool{}, map[Origin]bool{}
	followed, cleared := map[int]bool{}, map[int]bool{}
	visited := map[place]bool{}
	queue := []place{start}
	for len(queue) > 0 {
		here := queue[0]
		queue = queue[1:]
		if visited[here] {
			continue
		}
		visited[here] = true
		name := f.storage.at(here.record, here.start, here.length)
		if p := f.storage.places[name]; p.value {
			found[Origin{Item: name}] = true
		}
		if n := f.parameters[here.record]; n > 0 {
			found[Origin{Parameter: n, Item: name}] = true
		}
		for i, e := range f.Lineage {
			to, ok := f.storage.places[e.To]
			if e.ToFile || e.ToCall || e.ToStream || !ok || !to.overlaps(here) {
				continue
			}
			switch {
			case e.FromFile:
				found[Origin{File: e.From, Item: name}] = true
				followed[i] = true
			case e.FromCall:
				found[Origin{Call: e.From, Arg: e.Arg, Item: name}] = true
				followed[i] = true
			case e.FromStream:
				found[Origin{Stream: e.From, Item: name}] = true
				followed[i] = true
			case e.FromLiteral && to.group && to.length > here.length:
				clearing[Origin{Literal: e.From}] = true
				cleared[i] = true
			case e.FromLiteral:
				found[Origin{Literal: e.From}] = true
				followed[i] = true
			default:
				if from, ok := f.storage.places[e.From]; ok {
					queue = append(queue, follow(here, to, from))
					followed[i] = true
				}
			}
		}
	}
	if len(found) == 0 {
		found, followed = clearing, cleared
	}
	t := Trace{Origins: make([]Origin, 0, len(found))}
	for o := range found {
		t.Origins = append(t.Origins, o)
	}
	slices.SortFunc(t.Origins, func(a, b Origin) int { return strings.Compare(a.String(), b.String()) })
	for i, e := range f.Lineage {
		if followed[i] {
			t.Edges = append(t.Edges, e)
		}
	}
	return t
}

// follow maps bytes of a statement's target back to its source: the same
// bytes, by offset, when a group is moved into a group or an item into one
// of its own length, and the whole source otherwise.
func follow(here, to, from place) place {
	if (!to.group || !from.group) && to.length != from.length {
		return from
	}
	start := max(here.start, to.start)
	end := min(here.start+here.length, to.start+to.length)
	offset := start - to.start
	length := min(end-start, from.length-offset)
	if length <= 0 {
		return from
	}
	return place{record: from.record, start: from.start + offset, length: length, group: true}
}
