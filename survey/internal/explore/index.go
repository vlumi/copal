// Package explore serves the survey of a system as a web app, for
// moving through it rather than reading one blueprint at a time: from a
// job's network to a step, its program, a paragraph, a record and a
// field, where every name is a link and a search finds any of them. The
// pages are rendered on the server, with a little script.
package explore

import (
	"slices"
	"strings"

	"github.com/vlumi/copal/survey/blueprint"
)

// index looks up the system's parts by name.
type index struct {
	system    *blueprint.System
	programs  map[string]*blueprint.Program
	jobs      map[string]*blueprint.Job
	copybooks map[string]*blueprint.Copybook
	// datasets are the job steps that use each dataset, by the name of
	// the dataset without its generation.
	datasets map[string][]DatasetUse
	// runs are the job steps that run each program.
	runs map[string][]blueprint.StepRef
	// entities are the system's entities by record, and entityFields the
	// field of an entity each program's field is, by program and field.
	entities     map[string]*blueprint.Entity
	entityFields map[[2]string]EntityFieldRef
	// systemDatasets are the system's datasets by name.
	systemDatasets map[string]*blueprint.SystemDataset
	sources        map[string]bool
	names          []Name
}

// DatasetUse is a job step's DD naming a dataset, and how the step's
// program uses it.
type DatasetUse struct {
	blueprint.StepRef
	Dataset  string // as the job names it, such as ACCT.MASTER(+1)
	DD       string
	Disp     string
	Program  string
	File     string
	Copybook string
	Uses     []string
}

// EntityFieldRef names a field of an entity.
type EntityFieldRef struct {
	Record, Field string
}

// Name is anything the search finds: its kind, its name, what it belongs
// to, its page, and the badge of the kind of program, field or dataset it
// is.
type Name struct {
	Kind, Name, In, Link, Badge string
}

func newIndex(s *blueprint.System) *index {
	x := &index{system: s, programs: map[string]*blueprint.Program{}, jobs: map[string]*blueprint.Job{},
		copybooks: map[string]*blueprint.Copybook{}, datasets: map[string][]DatasetUse{},
		runs: map[string][]blueprint.StepRef{}, entities: map[string]*blueprint.Entity{},
		entityFields: map[[2]string]EntityFieldRef{}, systemDatasets: map[string]*blueprint.SystemDataset{}, sources: map[string]bool{}}
	for i := range s.Programs {
		p := &s.Programs[i]
		x.programs[p.Name] = p
		x.sources[p.Source] = true
		x.add("program", p.Name, "", programLink(p.Name), p.Kind)
		for _, para := range p.Paragraphs {
			x.add("paragraph", para.Name, p.Name, paragraphLink(p.Name, para.Name), "")
		}
		for _, r := range p.Records {
			walkFields(r.Field, func(f blueprint.Field) {
				x.sources[sourceFile(f.Source)] = true
				x.add("field", f.Name, p.Name, fieldLink(p.Name, f.Name), f.Kind)
			})
		}
	}
	for i := range s.Copybooks {
		c := &s.Copybooks[i]
		x.copybooks[c.Name] = c
		x.sources[c.Source] = true
		x.add("copybook", c.Name, "", copybookLink(c.Name), "")
	}
	for i := range s.Jobs {
		j := &s.Jobs[i]
		x.jobs[j.Name] = j
		x.sources[j.Source] = true
		x.add("job", j.Name, "", jobLink(j.Name), "")
		for _, step := range j.Steps {
			ref := blueprint.StepRef{Job: j.Name, Step: step.Name}
			x.add("step", step.Name, j.Name, stepLink(ref), "")
			x.runs[step.Program] = append(x.runs[step.Program], ref)
			for _, dd := range step.DDs {
				for _, ds := range dd.Datasets {
					base := datasetBase(ds)
					x.datasets[base] = append(x.datasets[base], DatasetUse{StepRef: ref, Dataset: ds, DD: dd.Name,
						Disp: dd.Disp, Program: step.Program, File: dd.File, Copybook: dd.Copybook, Uses: dd.Uses})
				}
			}
		}
	}
	for i := range s.Entities {
		e := &s.Entities[i]
		x.entities[e.Record] = e
		for _, f := range e.Fields {
			for _, u := range f.Uses {
				x.entityFields[[2]string{u.Program, u.Field}] = EntityFieldRef{Record: e.Record, Field: f.Name}
			}
		}
	}
	for i := range s.Datasets {
		d := &s.Datasets[i]
		x.systemDatasets[d.Name] = d
		x.add("dataset", d.Name, "", datasetLink(d.Name), d.Kind)
	}
	slices.SortFunc(x.names, func(a, b Name) int {
		return strings.Compare(a.Name+" "+a.In+" "+a.Kind, b.Name+" "+b.In+" "+b.Kind)
	})
	return x
}

func (x *index) add(kind, name, in, link, badge string) {
	x.names = append(x.names, Name{Kind: kind, Name: name, In: in, Link: link, Badge: badge})
}

// search finds the names that contain the query, ignoring case, those
// that start with it first.
func (x *index) search(query string) []Name {
	q := strings.ToUpper(strings.TrimSpace(query))
	if q == "" {
		return nil
	}
	var starts, contains []Name
	for _, n := range x.names {
		switch {
		case strings.HasPrefix(n.Name, q):
			starts = append(starts, n)
		case strings.Contains(n.Name, q):
			contains = append(contains, n)
		}
	}
	return append(starts, contains...)
}

// datasetBase is a dataset's name without its generation or the && of a
// temporary one, which keeps its name only within its job.
func datasetBase(name string) string {
	base, _, _ := strings.Cut(name, "(")
	return base
}

func walkFields(f blueprint.Field, visit func(blueprint.Field)) {
	visit(f)
	for _, child := range f.Fields {
		walkFields(child, visit)
	}
}

// sourceFile is the file of a source reference, path or path:line.
func sourceFile(source string) string {
	file, _, _ := strings.Cut(source, ":")
	return file
}
