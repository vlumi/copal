package blueprint

import (
	"slices"
	"strings"

	"github.com/vlumi/copal/survey/internal/diagram"
)

// JobOrder is one job running before another, as the generations of the
// datasets between them say.
type JobOrder struct {
	Before   string   `json:"before"`
	After    string   `json:"after"`
	Datasets []string `json:"datasets"`
}

// jobOrder derives the order of the jobs from the generations they write
// and read: a job that writes a new generation of a dataset, NAME(+1),
// runs before a job that reads its current one, NAME(0), an earlier one,
// NAME(-1), or every generation by its name. A job that reads what it
// writes itself runs before its own next run. The JCL says nothing else
// about when a job runs.
func jobOrder(jobs []Job) []JobOrder {
	var order []JobOrder
	for _, writer := range jobs {
		for _, reader := range jobs {
			o := JobOrder{Before: writer.Name, After: reader.Name}
			for _, w := range writer.Datasets {
				base, generation, _ := strings.Cut(w.Name, "(")
				if len(w.WrittenBy) == 0 || !strings.HasPrefix(generation, "+") {
					continue
				}
				if slices.ContainsFunc(reader.Datasets, func(r Dataset) bool { return len(r.ReadBy) > 0 && readsWritten(r.Name, base) }) {
					addOnce(&o.Datasets, base)
				}
			}
			if len(o.Datasets) > 0 {
				order = append(order, o)
			}
		}
	}
	return order
}

// readsWritten reports whether a dataset read, as a job names it, is a
// generation written before: the current one, an earlier one, or all.
func readsWritten(name, base string) bool {
	return name == base || name == base+"(0)" || strings.HasPrefix(name, base+"(-")
}

// JobCycles lists the jobs that the order makes a cycle of, each running
// both before and after the others, in the order of the jobs; a job alone
// is not one.
func JobCycles(jobs []Job, order []JobOrder) [][]string {
	after := map[string][]string{}
	for _, o := range order {
		after[o.Before] = append(after[o.Before], o.After)
	}
	reaches := func(from, to string) bool {
		seen := map[string]bool{from: true}
		queue := []string{from}
		for len(queue) > 0 {
			next := queue[0]
			queue = queue[1:]
			for _, n := range after[next] {
				if n == to {
					return true
				}
				if !seen[n] {
					seen[n] = true
					queue = append(queue, n)
				}
			}
		}
		return false
	}
	var cycles [][]string
	placed := map[string]bool{}
	for _, j := range jobs {
		if placed[j.Name] {
			continue
		}
		cycle := []string{j.Name}
		for _, k := range jobs {
			if k.Name != j.Name && reaches(j.Name, k.Name) && reaches(k.Name, j.Name) {
				cycle = append(cycle, k.Name)
				placed[k.Name] = true
			}
		}
		if len(cycle) > 1 {
			cycles = append(cycles, cycle)
		}
	}
	return cycles
}

// ScheduleGraph draws the order of the jobs: a box for each job, and an
// edge from each job to each that runs after it, labeled with the datasets
// between them. A job that reads what it wrote itself has no edge, and a
// cycle is drawn as one.
func ScheduleGraph(jobs []Job, order []JobOrder, link func(string) string) diagram.Graph {
	g := diagram.Graph{Legend: []diagram.Key{
		diagram.NodeKey(diagram.Node{}, "job"),
		diagram.EdgeKey(diagram.Solid, "runs before, through the datasets"),
		diagram.TextKey("it writes NAME(+1) the next reads"),
	}}
	for _, j := range jobs {
		n := diagram.Node{ID: j.Name, Label: j.Name}
		if link != nil {
			n.Link = link(j.Name)
		}
		g.Nodes = append(g.Nodes, n)
	}
	for _, o := range order {
		if o.Before != o.After {
			g.Edges = append(g.Edges, diagram.Edge{From: o.Before, To: o.After, Label: strings.Join(o.Datasets, "\n")})
		}
	}
	return g
}
