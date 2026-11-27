package oracle

import (
	"slices"
	"testing"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/cobolith/record"
	"github.com/vlumi/copal/survey/blueprint"
)

// products is a program with a table of products for its one input, a
// file in no order or one indexed by the product's code.
func products(key string) (*program, input) {
	root := &layout.Field{Name: "PRODUCT", Start: 1, Length: 4, Children: []*layout.Field{
		{Name: "PR-CODE", Start: 1, Length: 2, Picture: "X(2)", Category: layout.Alphanumeric},
		{Name: "PR-KIND", Start: 3, Length: 1, Picture: "9", Category: layout.Numeric, Digits: 1,
			Conditions: []layout.Condition{{Name: "PR-WHOLE-LIFE", Values: []string{"1"}}, {Name: "PR-TERM", Values: []string{"2"}}}},
		{Name: "PR-CLASS", Start: 4, Length: 1, Picture: "X", Category: layout.Alphanumeric,
			Conditions: []layout.Condition{{Name: "PR-RIDER", Values: []string{`"T"`}}}},
	}}
	in := input{file: blueprint.File{Name: "PRODUCTS", RecordKey: key}, root: root, layout: record.New(root, cobol.IBM)}
	return &program{inputs: []input{in}, shapes: map[string]shape{"PRODUCTS": {}}}, in
}

// A record added to an input for another value of a field has the value
// after the last record's, and another of the named values of each field
// that has them; an input whose records already differ in the field gets
// none.
func TestAnother(t *testing.T) {
	p, in := products("")
	code := target{input: "PRODUCTS", field: find(in.root, "PR-CODE")}
	c := Case{"PRODUCTS": {[]byte("A 1 ")}}
	if at, ok := p.another(c, code); !ok || at != 1 {
		t.Fatalf("another record at %d, %v, want at 1", at, ok)
	}
	want := [][]byte{[]byte("A 1 "), []byte("B 2T")}
	if !slices.EqualFunc(c["PRODUCTS"], want, slices.Equal) {
		t.Errorf("records %q, want %q", c["PRODUCTS"], want)
	}
	if at, ok := p.another(c, code); !ok || at != 1 || len(c["PRODUCTS"]) != 2 {
		t.Errorf("records %q, want the two that differ kept", c["PRODUCTS"])
	}
	empty := Case{"PRODUCTS": {[]byte("  1 ")}}
	if _, ok := p.another(empty, code); ok {
		t.Errorf("records %q, of a code with none after it", empty["PRODUCTS"])
	}
}

// A record added for another value of a field outside an indexed input's
// key would have the key of the one it was copied from, and is not added.
func TestAnotherOfTheSameKey(t *testing.T) {
	p, in := products("PR-CODE")
	c := Case{"PRODUCTS": {[]byte("A 1 ")}}
	if _, ok := p.another(c, target{input: "PRODUCTS", field: find(in.root, "PR-KIND")}); ok {
		t.Errorf("records %q, two of one key", c["PRODUCTS"])
	}
	if _, ok := p.another(c, target{input: "PRODUCTS", field: find(in.root, "PR-CODE")}); !ok {
		t.Error("no record of another key")
	}
}

// The last record of an input differs when it is unlike the first of its
// type, not when it is the only one of its type.
func TestLastDiffers(t *testing.T) {
	_, in := products("")
	kind := find(in.root, "PR-KIND")
	typed := shape{discriminator: kind, types: []value{number(1), number(2)}}
	for _, tt := range []struct {
		name    string
		s       shape
		records []string
		want    bool
	}{
		{"one record", shape{}, []string{"A 1 "}, false},
		{"two alike", shape{}, []string{"A 1 ", "A 1 "}, false},
		{"two that differ", shape{}, []string{"A 1 ", "B 1 "}, true},
		{"one of each type", typed, []string{"A 1 ", "A 2 "}, false},
		{"two of a type that differ", typed, []string{"A 1 ", "A 2 ", "B 2 "}, true},
	} {
		var records [][]byte
		for _, r := range tt.records {
			records = append(records, []byte(r))
		}
		if got := tt.s.lastDiffers(in, records); got != tt.want {
			t.Errorf("%s: %v, want %v", tt.name, got, tt.want)
		}
	}
}

// A record is added after the last that has the field, before a trailer
// that has none, and the first record that differs from the first is the
// other one, not the last.
func TestAnotherAmongTypes(t *testing.T) {
	p, in := products("")
	kind := find(in.root, "PR-KIND")
	class := find(in.root, "PR-CLASS")
	p.shapes["PRODUCTS"] = shape{discriminator: kind, types: []value{number(1), number(9)},
		variants: map[*layout.Field]int{class: 0}}
	c := Case{"PRODUCTS": {[]byte("A 1 "), []byte("  9 ")}}
	at, ok := p.another(c, target{input: "PRODUCTS", field: class})
	if ok {
		t.Fatalf("a record at %d for a class with none after it", at)
	}
	c = Case{"PRODUCTS": {[]byte("A 1S"), []byte("  9 ")}}
	want := [][]byte{[]byte("A 1S"), []byte("A 1T"), []byte("  9 ")}
	if at, ok := p.another(c, target{input: "PRODUCTS", field: class}); !ok || at != 1 || !slices.EqualFunc(c["PRODUCTS"], want, slices.Equal) {
		t.Errorf("another record at %d, %v, in %q, want %q", at, ok, c["PRODUCTS"], want)
	}
	c = Case{"PRODUCTS": {[]byte("A 1S"), []byte("A 1T"), []byte("A 1U"), []byte("  9 ")}}
	if at, ok := p.another(c, target{input: "PRODUCTS", field: class}); !ok || at != 1 || len(c["PRODUCTS"]) != 4 {
		t.Errorf("another record at %d, %v, of %d, want the second of four", at, ok, len(c["PRODUCTS"]))
	}
}
