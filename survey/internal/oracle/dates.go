package oracle

import (
	"bytes"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/copal/survey/blueprint"
)

// DateForms are the forms a system writes a date in: for each shape of
// field, text or digits of a length, what each of a few dates looks
// like. A year of an era, a date with separators and a date in digits
// are all text to a layout, and six digits are a month of a year in one
// system and a day without its century in another; only the system's
// own programs say which.
type DateForms map[dateShape]map[int64]string

// dateShape is the shape of a field that holds a date: digits or text,
// and how many bytes.
type dateShape struct {
	numeric bool
	length  int
}

// probeDates are the dates the forms are learned for: the business date
// first, then the days around it, a year before and five years before,
// so that a case can put one date before or after another, by a little
// or by years.
var probeDates = []int64{businessDate, businessDate - 1, businessDate + 1, businessDate - 10000, businessDate - 50000}

// of is a date in the form of a field, if the system was seen to write
// one in a field of its shape.
func (d DateForms) of(f *layout.Field, date int64) (string, bool) {
	t, ok := d[shapeOfDate(f)][date]
	return t, ok
}

// shapeOfDate is the shape a date field's form is learned by.
func shapeOfDate(f *layout.Field) dateShape {
	return dateShape{numeric: numeric(f), length: f.Length}
}

// LearnDates runs each main program of a system once for each probe
// date, on its base case with the dates of its files left out, so that
// the dates it writes are those it makes from the date it runs on, and
// reads them. A form is learned for a shape of field when every date
// written in a field of that shape for a probe date is the same, and the
// probe dates' are not all one: then it follows the date, and is how the
// system writes it.
func LearnDates(system *blueprint.System, copybookDir, scratch string) (DateForms, error) {
	seen := map[dateShape]map[int64][]string{}
	for i := range system.Programs {
		bp := &system.Programs[i]
		if bp.Kind != blueprint.KindMain {
			continue
		}
		s, err := newSearcher(system, bp, copybookDir, scratch)
		if err != nil {
			return nil, err
		}
		if len(s.program.writtenDates()) == 0 {
			continue
		}
		cases := make([]Case, len(probeDates))
		for j, date := range probeDates {
			cases[j] = s.program.probe(date)
		}
		results, err := s.runAll(cases, true)
		if err != nil {
			return nil, err
		}
		for j, res := range results {
			for shape, dates := range s.program.datesIn(res.outputs) {
				if seen[shape] == nil {
					seen[shape] = map[int64][]string{}
				}
				seen[shape][probeDates[j]] = append(seen[shape][probeDates[j]], dates...)
			}
		}
	}
	forms := DateForms{}
	for shape, byDate := range seen {
		if form, ok := formOf(byDate); ok {
			forms[shape] = form
		}
	}
	return forms, nil
}

// formOf is the form the dates written in fields of one shape show, by
// probe date: for each probe date the value most of them hold, more than
// half, when the probe dates' are not all one. A date a program makes
// partly from a date of its files, the year it runs in with the day of a
// contract, has nothing of the file's to take in a probe, and is the odd
// one out.
func formOf(byDate map[int64][]string) (map[int64]string, bool) {
	form := map[int64]string{}
	for date, dates := range byDate {
		count := map[string]int{}
		for _, t := range dates {
			if count[t]++; 2*count[t] > len(dates) {
				form[date] = t
			}
		}
	}
	distinct := slices.Compact(slices.Sorted(maps.Values(form)))
	return form, len(form) == len(probeDates) && len(distinct) > 1
}

// probe is the base case on a date with no date in the records of its
// files: only the date the run is given, on its command line or SYSIN.
func (p *program) probe(date int64) Case {
	c := p.baseOn(date)
	kinds := p.kinds()
	for _, in := range p.inputs {
		if in.file.Name == "" {
			continue
		}
		for _, r := range c[in.name()] {
			walk(in.root, func(leaf *layout.Field) {
				if !dateField(leaf, kinds) {
					return
				}
				if numeric(leaf) {
					set(r, in.root, leaf, in.layout, number(0))
				} else {
					set(r, in.root, leaf, in.layout, text(""))
				}
			})
		}
	}
	return c
}

// dateField reports whether a field holds a date whose form can be
// learned: text, or unsigned digits a byte each.
func dateField(f *layout.Field, kinds map[string]string) bool {
	if len(f.Children) > 0 || kinds[f.Name] != blueprint.KindDate {
		return false
	}
	return !numeric(f) || (f.Usage == "" || f.Usage == "DISPLAY") && !f.Signed && f.Scale == 0 && f.Digits == f.Length
}

// writtenDates are the date fields of the files the program writes or
// updates, by file.
func (p *program) writtenDates() map[string][]*layout.Field {
	kinds := p.kinds()
	fields := map[string][]*layout.Field{}
	for _, f := range p.files {
		if !writes(f) && !slices.Contains(f.Opened, "I-O") {
			continue
		}
		root := fileRecord(p.recordsOf(f), f.Name)
		if root == nil {
			continue
		}
		walk(root, func(leaf *layout.Field) {
			if dateField(leaf, kinds) {
				fields[f.Name] = append(fields[f.Name], leaf)
			}
		})
	}
	return fields
}

// datesIn reads the dates from what a run wrote, a record a line as the
// runner returns it, by the shape of their fields: text that is not
// blank, and digits that are not zero.
func (p *program) datesIn(outputs map[string][]byte) map[dateShape][]string {
	found := map[dateShape][]string{}
	for file, fields := range p.writtenDates() {
		for _, line := range bytes.Split(bytes.TrimSuffix(outputs[file], []byte("\n")), []byte("\n")) {
			record := string(line)
			if unquoted, err := strconv.Unquote(record); err == nil {
				record = unquoted
			}
			for _, f := range fields {
				if f.Start-1+f.Length > len(record) {
					continue
				}
				t := strings.TrimSpace(record[f.Start-1 : f.Start-1+f.Length])
				if n, err := strconv.ParseInt(t, 10, 64); numeric(f) && (err != nil || n == 0) || t == "" {
					continue
				}
				found[shapeOfDate(f)] = append(found[shapeOfDate(f)], t)
			}
		}
	}
	return found
}

// dateLevers are the changes that put each probe date, in the form the
// system writes it, into each date field of the inputs whose form was
// learned. Every such field is tried, not only those a condition names:
// a date reaches a condition through a called program's conversion more
// often than not, and the origins of a condition's items stop short of
// some of those.
func (p *program) dateLevers() []lever {
	kinds := p.kinds()
	var levers []lever
	for _, in := range p.inputs {
		walk(in.root, func(leaf *layout.Field) {
			if !dateField(leaf, kinds) {
				return
			}
			for _, date := range probeDates {
				if t, ok := p.forms.of(leaf, date); ok {
					levers = append(levers, lever{input: in.name(), field: leaf, value: text(t)})
				}
			}
		})
	}
	return levers
}
