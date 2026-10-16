package dataflow

import (
	"slices"
	"strings"

	"github.com/vlumi/copal/layout"
)

// Origin is where a data item's values come from in the end: a field of a
// record a file is read into, a literal, or an item's initial VALUE.
type Origin struct {
	File    string // the file, for a field read from one
	Item    string // the field read, or the item with the initial value
	Literal string // the literal, as written
}

// String writes the origin: "file TRAN-IN TD-AMOUNT", "literal 'TXPOST'",
// "initial value of WS-COUNT".
func (o Origin) String() string {
	switch {
	case o.File != "":
		return "file " + o.File + " " + o.Item
	case o.Literal != "":
		return "literal " + o.Literal
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
	all := records.WorkingStorage
	for _, f := range records.Files {
		all = append(all, f.Records...)
	}
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
	for _, r := range all {
		walk(r.Name, r)
	}
	for name, n := range seen {
		if n > 1 {
			delete(s.places, name)
		}
	}
	return s
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

// Origins traces where a data item's values come from: through every
// statement that moves a value into it or into bytes it shares, a group it
// is in, a field of it or a field that redefines it, and on through what
// those statements move from, to the files, literals and initial values
// at the start. A group moved into a group, or anything moved into an item
// of its own length, is followed byte for byte. A literal moved into a
// whole group around the item, such as MOVE SPACES TO A-RECORD, only
// clears it, and counts only for an item with no other origin.
func (f Flow) Origins(item string) []Origin {
	start, ok := f.storage.places[item]
	if !ok {
		return nil
	}
	found, clearing := map[Origin]bool{}, map[Origin]bool{}
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
		for _, e := range f.Lineage {
			to, ok := f.storage.places[e.To]
			if e.ToFile || !ok || !to.overlaps(here) {
				continue
			}
			switch {
			case e.FromFile:
				found[Origin{File: e.From, Item: name}] = true
			case e.FromLiteral && to.group && to.length > here.length:
				clearing[Origin{Literal: e.From}] = true
			case e.FromLiteral:
				found[Origin{Literal: e.From}] = true
			default:
				if from, ok := f.storage.places[e.From]; ok {
					queue = append(queue, follow(here, to, from))
				}
			}
		}
	}
	if len(found) == 0 {
		found = clearing
	}
	origins := make([]Origin, 0, len(found))
	for o := range found {
		origins = append(origins, o)
	}
	slices.SortFunc(origins, func(a, b Origin) int { return strings.Compare(a.String(), b.String()) })
	return origins
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
