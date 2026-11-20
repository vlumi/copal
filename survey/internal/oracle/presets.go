package oracle

import (
	"bytes"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
	"golang.org/x/text/encoding/japanese"
)

// presets are the values a program's tables are written with: for each
// field of an area that redefines one made of VALUE literals alone, the
// table of accounts a program carries in its working storage, what each
// occurrence of the field holds, by the field's name. They are values an
// input field compared with the table's has to take to be found in it.
func presets(records *layout.ProgramRecords, encoding string) map[string][]string {
	found := map[string][]string{}
	var among func(siblings []*layout.Field)
	among = func(siblings []*layout.Field) {
		for _, f := range siblings {
			among(f.Children)
			if f.Redefines == "" {
				continue
			}
			for _, written := range siblings {
				if written.Name != f.Redefines {
					continue
				}
				if content, ok := contentOf(written, encoding); ok {
					occurrences(f, f.Start, []int{0}, func(leaf *layout.Field, at int) {
						if at+leaf.Length <= len(content) {
							if t, err := cobol.ToUTF8(string(content[at:at+leaf.Length]), encoding); err == nil {
								found[leaf.Name] = append(found[leaf.Name], strings.TrimRight(t, " "))
							}
						}
					})
				}
			}
		}
	}
	among(records.All())
	return found
}

// contentOf is the bytes an area holds as written, when each of its
// elementary fields has one literal for its VALUE and none repeats.
func contentOf(area *layout.Field, encoding string) ([]byte, bool) {
	var content []byte
	ok := true
	walk(area, func(f *layout.Field) {
		if len(f.Children) > 0 {
			ok = ok && f.Occurs == 0
			return
		}
		if f.Occurs > 0 || len(f.Values) != 1 || len(f.Values[0]) < 2 || !strings.ContainsAny(f.Values[0][:1], `'"`) {
			ok = false
			return
		}
		lit, good := encode(f.Values[0][1:len(f.Values[0])-1], encoding)
		if !good || len(lit) > f.Length {
			ok = false
			return
		}
		content = append(append(content, lit...), bytes.Repeat([]byte(" "), f.Length-len(lit))...)
	})
	return content, ok
}

// occurrences visits each occurrence of each named elementary field of
// an area, with where its bytes are from the area's start.
func occurrences(f *layout.Field, start int, shifts []int, visit func(leaf *layout.Field, at int)) {
	if f.Occurs > 0 {
		var each []int
		for _, shift := range shifts {
			for i := range f.Occurs {
				each = append(each, shift+i*f.Length)
			}
		}
		shifts = each
	}
	if len(f.Children) == 0 {
		if f.Name != "FILLER" {
			for _, shift := range shifts {
				visit(f, f.Start-start+shift)
			}
		}
		return
	}
	for _, child := range f.Children {
		occurrences(child, start, shifts, visit)
	}
}

// encode gives text in the encoding a source is read in, as its bytes are
// in a record.
func encode(text, encoding string) ([]byte, bool) {
	if encoding != cobol.ShiftJIS {
		return []byte(text), true
	}
	b, err := japanese.ShiftJIS.NewEncoder().Bytes([]byte(text))
	return b, err == nil
}
