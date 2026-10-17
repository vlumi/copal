// Package explore serves Copal's analysis of a system as a web app, for
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
	runs    map[string][]StepRef
	sources map[string]bool
	names   []Name
}

// StepRef names a step of a job.
type StepRef struct {
	Job, Step string
}

// DatasetUse is a job step's DD naming a dataset, and how the step's
// program uses it.
type DatasetUse struct {
	StepRef
	Dataset  string // as the job names it, such as ACCT.MASTER(+1)
	DD       string
	Disp     string
	Program  string
	File     string
	Copybook string
	Uses     []string
}

// Name is anything the search finds: its kind, its name, what it belongs
// to, and its page.
type Name struct {
	Kind, Name, In, Link string
}

func newIndex(s *blueprint.System) *index {
	x := &index{system: s, programs: map[string]*blueprint.Program{}, jobs: map[string]*blueprint.Job{},
		copybooks: map[string]*blueprint.Copybook{}, datasets: map[string][]DatasetUse{},
		runs: map[string][]StepRef{}, sources: map[string]bool{}}
	for i := range s.Programs {
		p := &s.Programs[i]
		x.programs[p.Name] = p
		x.sources[p.Source] = true
		x.add("program", p.Name, "", programLink(p.Name))
		for _, para := range p.Paragraphs {
			x.add("paragraph", para.Name, p.Name, paragraphLink(p.Name, para.Name))
		}
		for _, r := range p.Records {
			walkFields(r.Field, func(f blueprint.Field) {
				x.sources[sourceFile(f.Source)] = true
				x.add("field", f.Name, p.Name, fieldLink(p.Name, f.Name))
			})
		}
	}
	for i := range s.Copybooks {
		c := &s.Copybooks[i]
		x.copybooks[c.Name] = c
		x.sources[c.Source] = true
		x.add("copybook", c.Name, "", copybookLink(c.Name))
	}
	for i := range s.Jobs {
		j := &s.Jobs[i]
		x.jobs[j.Name] = j
		x.sources[j.Source] = true
		x.add("job", j.Name, "", jobLink(j.Name))
		for _, step := range j.Steps {
			ref := StepRef{Job: j.Name, Step: step.Name}
			x.add("step", step.Name, j.Name, stepLink(ref))
			x.runs[step.Program] = append(x.runs[step.Program], ref)
			for _, dd := range step.DDs {
				for _, ds := range dd.Datasets {
					base := datasetBase(ds)
					if _, ok := x.datasets[base]; !ok {
						x.add("dataset", base, "", datasetLink(base))
					}
					x.datasets[base] = append(x.datasets[base], DatasetUse{StepRef: ref, Dataset: ds, DD: dd.Name,
						Disp: dd.Disp, Program: step.Program, File: dd.File, Copybook: dd.Copybook, Uses: dd.Uses})
				}
			}
		}
	}
	slices.SortFunc(x.names, func(a, b Name) int {
		return strings.Compare(a.Name+" "+a.In+" "+a.Kind, b.Name+" "+b.In+" "+b.Kind)
	})
	return x
}

func (x *index) add(kind, name, in, link string) {
	x.names = append(x.names, Name{Kind: kind, Name: name, In: in, Link: link})
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
