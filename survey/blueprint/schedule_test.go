package blueprint

import (
	"reflect"
	"testing"
)

// A job that writes NAME(+1) runs before one that reads NAME(0), NAME(-1)
// or every generation; one that reads what it wrote runs before its own
// next run; and the jobs that run before and after each other are a cycle.
func TestJobOrder(t *testing.T) {
	ds := func(name string, written, read bool) Dataset {
		d := Dataset{Name: name}
		if written {
			d.WrittenBy = []string{"S"}
		}
		if read {
			d.ReadBy = []string{"S"}
		}
		return d
	}
	jobs := []Job{
		{Name: "DAY", Datasets: []Dataset{ds("CTL(0)", false, true), ds("CTL(+1)", true, false), ds("LOG(+1)", true, false), ds("&&T", true, true)}},
		{Name: "MONTH", Datasets: []Dataset{ds("LOG", false, true), ds("CTL(-1)", false, true)}},
		{Name: "NIGHT", Datasets: []Dataset{ds("CTL(0)", false, true), ds("FEED(+1)", true, false)}},
		{Name: "OPEN", Datasets: []Dataset{ds("FEED(0)", false, true), ds("OUT(+1)", true, false)}},
	}
	order := jobOrder(jobs)
	want := []JobOrder{
		{Before: "DAY", After: "DAY", Datasets: []string{"CTL"}},
		{Before: "DAY", After: "MONTH", Datasets: []string{"CTL", "LOG"}},
		{Before: "DAY", After: "NIGHT", Datasets: []string{"CTL"}},
		{Before: "NIGHT", After: "OPEN", Datasets: []string{"FEED"}},
	}
	if !reflect.DeepEqual(order, want) {
		t.Errorf("order\n%+v\nwant\n%+v", order, want)
	}
	if cycles := JobCycles(jobs, order); cycles != nil {
		t.Errorf("cycles %q, want none: a job before its own next run is not one", cycles)
	}
	order = append(order, JobOrder{Before: "OPEN", After: "DAY", Datasets: []string{"X"}})
	if got, want := JobCycles(jobs, order), [][]string{{"DAY", "NIGHT", "OPEN"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("cycles %q, want %q", got, want)
	}
	g := ScheduleGraph(jobs, order, nil)
	if len(g.Nodes) != 4 || len(g.Edges) != 4 || g.Edges[0].Label != "CTL\nLOG" {
		t.Errorf("graph %+v", g)
	}
}
