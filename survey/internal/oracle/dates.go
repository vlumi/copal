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

// DateForms are the forms a system writes a date in when it keeps it as
// text: for each length of field, what each of a few dates looks like.
// A year of an era, a date with separators and a date in digits are all
// text to a layout, and only the system's own programs say which it is.
type DateForms map[int]map[int64]string

// probeDates are the dates the forms are learned for: the business date
// first, then the days around it, a year before and five years before,
// so that a case can put one date before or after another, by a little
// or by years.
var probeDates = []int64{businessDate, businessDate - 1, businessDate + 1, businessDate - 10000, businessDate - 50000}

// text is a date in the form of a text field of a length, if the system
// was seen to write one.
func (d DateForms) text(length int, date int64) (string, bool) {
	t, ok := d[length][date]
	return t, ok
}

// LearnDates runs each main program of a system on its base case, once
// for each probe date, and reads the dates the runs wrote as text. A
// form is learned for a length of field when every date written in a
// field of that length for a probe date is the same text, and the probe
// dates' texts all differ: then the text follows the date, and is how
// the system writes it.
func LearnDates(system *blueprint.System, copybookDir, scratch string) (DateForms, error) {
	seen := map[int]map[int64][]string{}
	for i := range system.Programs {
		bp := &system.Programs[i]
		if bp.Kind != blueprint.KindMain {
			continue
		}
		s, err := newSearcher(system, bp, copybookDir, scratch)
		if err != nil {
			return nil, err
		}
		if len(s.program.writtenTextDates()) == 0 {
			continue
		}
		cases := make([]Case, len(probeDates))
		for j, date := range probeDates {
			cases[j] = s.program.baseOn(date)
		}
		results, err := s.runAll(cases, true)
		if err != nil {
			return nil, err
		}
		for j, res := range results {
			for length, texts := range s.program.textDatesIn(res.outputs) {
				if seen[length] == nil {
					seen[length] = map[int64][]string{}
				}
				seen[length][probeDates[j]] = append(seen[length][probeDates[j]], texts...)
			}
		}
	}
	forms := DateForms{}
	for length, byDate := range seen {
		form := map[int64]string{}
		for date, texts := range byDate {
			if len(slices.Compact(slices.Sorted(slices.Values(texts)))) == 1 {
				form[date] = texts[0]
			}
		}
		distinct := slices.Compact(slices.Sorted(maps.Values(form)))
		if len(form) == len(probeDates) && len(distinct) == len(probeDates) {
			forms[length] = form
		}
	}
	return forms, nil
}

// textDate reports whether a field holds a date as text.
func textDate(f *layout.Field, kinds map[string]string) bool {
	return len(f.Children) == 0 && kinds[f.Name] == blueprint.KindDate && !numeric(f)
}

// writtenTextDates are the fields of the files the program writes or
// updates that hold a date as text, by file.
func (p *program) writtenTextDates() map[string][]*layout.Field {
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
			if textDate(leaf, kinds) {
				fields[f.Name] = append(fields[f.Name], leaf)
			}
		})
	}
	return fields
}

// textDatesIn reads the dates kept as text from what a run wrote, a
// record a line as the runner returns it, by the length of their fields.
func (p *program) textDatesIn(outputs map[string][]byte) map[int][]string {
	found := map[int][]string{}
	for file, fields := range p.writtenTextDates() {
		for _, line := range bytes.Split(bytes.TrimSuffix(outputs[file], []byte("\n")), []byte("\n")) {
			record := string(line)
			if unquoted, err := strconv.Unquote(record); err == nil {
				record = unquoted
			}
			for _, f := range fields {
				if f.Start-1+f.Length > len(record) {
					continue
				}
				if t := strings.TrimSpace(record[f.Start-1 : f.Start-1+f.Length]); t != "" {
					found[f.Length] = append(found[f.Length], t)
				}
			}
		}
	}
	return found
}

// dateLevers are the changes that put each probe date, in the form the
// system writes it, into each field of the inputs that holds a date as
// text. Every such field is tried, not only those a condition names: a
// date reaches a condition through a called program's conversion more
// often than not, and the origins of a condition's items stop short of
// some of those.
func (p *program) dateLevers() []lever {
	kinds := p.kinds()
	var levers []lever
	for _, in := range p.inputs {
		walk(in.root, func(leaf *layout.Field) {
			if !textDate(leaf, kinds) {
				return
			}
			for _, date := range probeDates {
				if t, ok := p.forms.text(leaf.Length, date); ok {
					levers = append(levers, lever{input: in.name(), field: leaf, value: text(t)})
				}
			}
		})
	}
	return levers
}
