package explore

import (
	"net/http"
	"slices"

	"github.com/vlumi/copal/survey/blueprint"
)

// BoundaryView is the system's contracts with the outside: the datasets
// that come in and go out, the reports of the job steps, and the in-stream
// data programs ACCEPT, with what they DISPLAY.
type BoundaryView struct {
	Inputs   []BoundaryDataset
	Outputs  []BoundaryDataset
	Reports  []Report
	Streams  []Stream
	Entities map[string][]string // the entities of each copybook
}

// BoundaryDataset is a dataset at the boundary.
type BoundaryDataset struct {
	*blueprint.SystemDataset
	Entities []string
}

// Report is a job step's DD to the job's output, with the program's file
// bound to it and the records of the file.
type Report struct {
	blueprint.StepRef
	DD, Program, File string
	Records           []string
}

// Stream is a job step's SYSIN that its program ACCEPTs, and what the
// program DISPLAYs to the job's output.
type Stream struct {
	blueprint.StepRef
	Program  string
	InStream bool
	Datasets []string
	Displays []blueprint.Display
}

func (s *Server) boundary(*http.Request) (any, error) {
	system := s.index.system
	entities := map[string][]string{}
	for _, e := range system.Entities {
		entities[e.Copybook] = append(entities[e.Copybook], e.Record)
	}
	v := BoundaryView{Entities: entities}
	for i := range system.Datasets {
		d := &system.Datasets[i]
		at := BoundaryDataset{SystemDataset: d, Entities: entities[d.Copybook]}
		switch d.Kind {
		case blueprint.KindExternal:
			v.Inputs = append(v.Inputs, at)
		case blueprint.KindOutput:
			v.Outputs = append(v.Outputs, at)
		}
	}
	for _, job := range system.Jobs {
		for _, step := range job.Steps {
			ref := blueprint.StepRef{Job: job.Name, Step: step.Name}
			p := s.index.programs[step.Program]
			for _, dd := range step.DDs {
				switch {
				case dd.Sysout:
					r := Report{StepRef: ref, DD: dd.Name, Program: step.Program, File: dd.File}
					if p != nil {
						if i := slices.IndexFunc(p.Files, func(f blueprint.File) bool { return f.Name == dd.File }); i >= 0 {
							r.Records = p.Files[i].Records
						}
					}
					v.Reports = append(v.Reports, r)
				case dd.Name == "SYSIN" && p != nil && p.ReadsSysin:
					v.Streams = append(v.Streams, Stream{StepRef: ref, Program: p.Name, InStream: dd.InStream,
						Datasets: dd.Datasets, Displays: p.Displays})
				}
			}
		}
	}
	return v, nil
}
