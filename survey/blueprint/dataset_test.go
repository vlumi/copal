package blueprint

import (
	"reflect"
	"testing"
)

// A dataset no step writes is external, one no step reads an output, a
// temporary one temporary whatever its use, and the rest internal; each
// generation counts as its dataset.
func TestSystemDatasets(t *testing.T) {
	use := func(name string, uses ...string) DD {
		return DD{Name: name, Datasets: []string{name}, Uses: uses}
	}
	in, master, next := use("IN", "read"), use("MASTER(0)", "read"), use("MASTER(+1)", "write")
	next.Indexed = true
	jobs := []Job{{Name: "J", Steps: []Step{
		{Name: "A", DDs: []DD{in, master, next, use("&&T", "write"), use("OUT", "write")}},
		{Name: "B", DDs: []DD{use("&&T", "read")}},
	}}}
	got := systemDatasets(jobs)
	a, b := StepRef{Job: "J", Step: "A"}, StepRef{Job: "J", Step: "B"}
	want := []SystemDataset{
		{Name: "IN", Kind: KindExternal, ReadBy: []StepRef{a}},
		{Name: "MASTER", Kind: KindInternal, Indexed: true, WrittenBy: []StepRef{a}, ReadBy: []StepRef{a}},
		{Name: "OUT", Kind: KindOutput, WrittenBy: []StepRef{a}},
		{Name: "&&T", Kind: KindTemporary, Temporary: true, WrittenBy: []StepRef{a}, ReadBy: []StepRef{b}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("datasets\n%+v\nwant\n%+v", got, want)
	}
}
