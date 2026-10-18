package dataflow

// Effect is where a data item's values go: the items they reach, the
// bytes of the files they are written to, the bytes of the arguments they
// are passed in to programs called, the bytes of the items the program was
// called with that they are moved into, back to the caller, and the
// statements that move them there, as lineage edges, in the order the
// program declares them.
type Effect struct {
	Items      []string
	Files      []Written
	Calls      []Passed
	Parameters []Passed
	Edges      []Edge
}

// Passed is a run of bytes of an item passed between programs: to Program,
// the program called, as the argument at Arg, or, with no Program, as the
// USING item at Arg, back to the caller; Start is 1-based, within the
// item.
type Passed struct {
	Program       string
	Arg           int
	Start, Length int
}

// Written is a run of bytes of the record a file is written from: the
// file, and where in its record the bytes are, 1-based.
type Written struct {
	File          string
	Start, Length int
}

// Affects follows where a data item's values go, the mirror of Trace:
// through every statement that moves a value out of it or out of bytes it
// shares, a group it is in, a field of it or a field that redefines it,
// on to the items those statements move into, and to the files the
// records it reaches are written to. A group moved into a group, or an
// item moved into one of its own length, is followed byte for byte.
func (f Flow) Affects(item string) Effect {
	start, ok := f.storage.places[item]
	if !ok {
		return Effect{}
	}
	return f.affects(start)
}

// AffectsBytes follows where the values of some of a data item's bytes go,
// from start, 1-based, for length, as Affects does: the bytes of a record
// that a FILLER holds, say, which no field names.
func (f Flow) AffectsBytes(item string, start, length int) Effect {
	p, ok := f.storage.part(item, start, length)
	if !ok {
		return Effect{}
	}
	return f.affects(p)
}

func (f Flow) affects(start place) Effect {
	var effect Effect
	reached, written, passed := map[string]bool{}, map[Written]bool{}, map[Passed]bool{}
	pass := func(p Passed, to *[]Passed) {
		if !passed[p] {
			passed[p] = true
			*to = append(*to, p)
		}
	}
	followed := map[int]bool{}
	visited := map[place]bool{}
	queue := []place{start}
	for len(queue) > 0 {
		here := queue[0]
		queue = queue[1:]
		if visited[here] {
			continue
		}
		visited[here] = true
		if name := f.storage.at(here.record, here.start, here.length); here != start && !reached[name] {
			reached[name] = true
			effect.Items = append(effect.Items, name)
		}
		if n := f.parameters[here.record]; n > 0 && here != start {
			pass(Passed{Arg: n, Start: here.start, Length: here.length}, &effect.Parameters)
		}
		for i, e := range f.Lineage {
			from, ok := f.storage.places[e.From]
			if e.FromFile || e.FromLiteral || e.FromCall || e.FromStream || e.ToStream || !ok || !from.overlaps(here) {
				continue
			}
			if e.ToCall {
				w := writtenBytes(e.To, here, from)
				pass(Passed{Program: e.To, Arg: e.Arg, Start: w.Start, Length: w.Length}, &effect.Calls)
				followed[i] = true
				continue
			}
			if e.ToFile {
				if w := writtenBytes(e.To, here, from); !written[w] {
					written[w] = true
					effect.Files = append(effect.Files, w)
				}
				followed[i] = true
				continue
			}
			if to, ok := f.storage.places[e.To]; ok {
				queue = append(queue, follow(here, from, to))
				followed[i] = true
			}
		}
	}
	for i, e := range f.Lineage {
		if followed[i] {
			effect.Edges = append(effect.Edges, e)
		}
	}
	return effect
}

// Share reports whether two data items share any bytes: one is in the
// other, they are the same, or one redefines the other.
func (f Flow) Share(a, b string) bool {
	pa, okA := f.storage.places[a]
	pb, okB := f.storage.places[b]
	return okA && okB && pa.overlaps(pb)
}

// writtenBytes is where bytes of a record are in the record written to a
// file.
func writtenBytes(file string, here, record place) Written {
	start := max(here.start, record.start)
	end := min(here.start+here.length, record.start+record.length)
	return Written{File: file, Start: start - record.start + 1, Length: end - start}
}
