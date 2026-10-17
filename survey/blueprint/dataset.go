package blueprint

import (
	"fmt"
	"slices"
	"strings"
)

// SystemDataset is a dataset across every job, by its name without the
// generation: the copybook that describes it, whether it is indexed or
// temporary, its kind, and the job steps that write and read it.
type SystemDataset struct {
	Name      string    `json:"dataset"`
	Copybook  string    `json:"copybook,omitempty"`
	Indexed   bool      `json:"indexed,omitempty"`
	Temporary bool      `json:"temporary,omitempty"`
	Kind      string    `json:"kind"`
	WrittenBy []StepRef `json:"writtenBy,omitempty"`
	ReadBy    []StepRef `json:"readBy,omitempty"`
}

// StepRef names a step of a job.
type StepRef struct {
	Job  string `json:"job"`
	Step string `json:"step"`
}

// Dataset kinds: a temporary dataset lives within its job; an external
// one no step writes, so it comes from outside the system; an output no
// step reads, so it goes outside; and an internal one passes from step to
// step.
const (
	KindTemporary = "temporary"
	KindExternal  = "external"
	KindOutput    = "output"
	KindInternal  = "internal"
)

// systemDatasets lists the datasets the jobs name, by name, temporary ones
// last, each with its kind and the steps that use it.
func systemDatasets(jobs []Job) []SystemDataset {
	byName := map[string]*SystemDataset{}
	var names []string
	for _, job := range jobs {
		copybooks := map[string]string{}
		for _, d := range job.Datasets {
			copybooks[d.Name] = d.Copybook
		}
		for _, step := range job.Steps {
			ref := StepRef{Job: job.Name, Step: step.Name}
			for _, dd := range step.DDs {
				for _, name := range dd.Datasets {
					base := datasetBase(name)
					d, ok := byName[base]
					if !ok {
						d = &SystemDataset{Name: base, Temporary: strings.HasPrefix(base, "&&")}
						byName[base] = d
						names = append(names, base)
					}
					if d.Copybook == "" {
						d.Copybook = copybooks[name]
					}
					d.Indexed = d.Indexed || dd.Indexed
					if dd.writes() && !slices.Contains(d.WrittenBy, ref) {
						d.WrittenBy = append(d.WrittenBy, ref)
					}
					if dd.reads() && !slices.Contains(d.ReadBy, ref) {
						d.ReadBy = append(d.ReadBy, ref)
					}
				}
			}
		}
	}
	slices.SortFunc(names, func(a, b string) int {
		ta, tb := strings.HasPrefix(a, "&&"), strings.HasPrefix(b, "&&")
		switch {
		case ta && !tb:
			return 1
		case tb && !ta:
			return -1
		}
		return strings.Compare(a, b)
	})
	datasets := make([]SystemDataset, 0, len(names))
	for _, name := range names {
		d := byName[name]
		d.Kind = datasetKind(*d)
		datasets = append(datasets, *d)
	}
	return datasets
}

func datasetKind(d SystemDataset) string {
	switch {
	case d.Temporary:
		return KindTemporary
	case len(d.WrittenBy) == 0:
		return KindExternal
	case len(d.ReadBy) == 0:
		return KindOutput
	}
	return KindInternal
}

// writeDatasets writes the index's table of the system's datasets.
func writeDatasets(b *strings.Builder, datasets []SystemDataset) {
	if len(datasets) == 0 {
		return
	}
	b.WriteString("\n## Datasets\n\nEvery dataset the jobs name, by its name without the generation. An external dataset is " +
		"one no step writes, which comes from outside the system; an output one no step reads; a temporary one, &&NAME, " +
		"lives within its job; and an internal one passes from step to step.\n\n" +
		"| Dataset | Kind | Copybook | Written by | Read by |\n| --- | --- | --- | --- | --- |\n")
	for _, d := range datasets {
		kind := d.Kind
		if d.Indexed {
			kind += ", indexed"
		}
		b.WriteString(row(d.Name, kind, strings.ReplaceAll(copybookLink(d.Copybook), "../", ""),
			stepLinks(d.WrittenBy), stepLinks(d.ReadBy)))
	}
}

func stepLinks(steps []StepRef) string {
	var links []string
	for _, s := range steps {
		links = append(links, fmt.Sprintf("[%s %s](jobs/%s.md#%s)", s.Job, s.Step, s.Job, headingAnchor(s.Step)))
	}
	return strings.Join(links, ", ")
}
