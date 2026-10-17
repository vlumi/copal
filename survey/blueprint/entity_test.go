package blueprint

import (
	"reflect"
	"testing"
)

// A program's record is a copybook's by where it is declared, so that a
// copy with REPLACING still counts, under the program's names.
func TestEntities(t *testing.T) {
	record := func(prefix string) Field {
		return Field{Name: prefix + "-REC", Source: "C.cpy:1", Fields: []Field{
			{Name: prefix + "-NO", Kind: KindNumber, Source: "C.cpy:2"},
			{Name: "FILLER", Kind: KindFiller, Source: "C.cpy:3"},
		}}
	}
	copybooks := []Copybook{{Name: "C", Records: []Record{{Field: record(":P:")}}}}
	programs := []Program{
		{Name: "P", Records: []Record{{Copybook: "C", Field: record("XX")}},
			Paragraphs: []Paragraph{{Name: "READ-IT", Reads: []string{"XX-NO"}}, {Name: "SET-IT", Writes: []string{"XX-NO"}}}},
		{Name: "Q", Records: []Record{{Copybook: "OTHER", Field: record("YY")}}},
	}
	want := []Entity{{Record: ":P:-REC", Copybook: "C", Programs: []string{"P"}, Fields: []EntityField{
		{Name: ":P:-NO", Kind: KindNumber, Uses: []FieldUse{{Program: "P", Field: "XX-NO", ReadBy: []string{"READ-IT"}, WrittenBy: []string{"SET-IT"}}}},
	}}}
	if got := entities(copybooks, programs); !reflect.DeepEqual(got, want) {
		t.Errorf("entities\n%+v\nwant\n%+v", got, want)
	}
}
