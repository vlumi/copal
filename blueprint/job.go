package blueprint

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/vlumi/copal/cobol"
	"github.com/vlumi/copal/diagram"
	"github.com/vlumi/copal/jcl"
)

// Job is the blueprint of a job: its steps, the datasets each step's DD
// statements bind, how the step's program uses each, and what does not
// fit between the job and its programs.
type Job struct {
	Name        string   `json:"job"`
	Source      string   `json:"source"`
	Description []string `json:"description,omitempty"`
	Steps       []Step   `json:"steps"`
	Findings    []string `json:"findings,omitempty"`
}

// Step is one step of a job: the program it runs, the return code test
// that bypasses it, and its DD statements.
type Step struct {
	Name        string   `json:"step"`
	Program     string   `json:"program"`
	Cond        string   `json:"cond,omitempty"`
	Description []string `json:"description,omitempty"`
	Source      string   `json:"source"`
	DDs         []DD     `json:"dds"`
}

// DD is a DD statement: the datasets it binds, as the job names them, or
// the job's output, or nothing; and the program's file assigned to it with
// how the program uses it.
type DD struct {
	Name     string   `json:"dd"`
	Datasets []string `json:"datasets,omitempty"`
	Disp     string   `json:"disp,omitempty"`
	Sysout   bool     `json:"sysout,omitempty"`
	Dummy    bool     `json:"dummy,omitempty"`
	File     string   `json:"file,omitempty"`
	Use      string   `json:"use,omitempty"`
	Source   string   `json:"source"`
}

// Uses of a dataset, from the modes the program opens its file in.
const (
	useRead   = "read"
	useWrite  = "write"
	useAppend = "append"
	useUpdate = "update"
)

var uses = map[string]string{"INPUT": useRead, "OUTPUT": useWrite, "EXTEND": useAppend, "I-O": useUpdate}

func analyzeJob(path string, programs []Program) (Job, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Job{}, err
	}
	parsed, err := jcl.Parse(bytes.NewReader(data))
	if err != nil {
		return Job{}, fmt.Errorf("%s: %w", path, err)
	}
	job := Job{Name: parsed.Name, Source: path, Description: cobol.CommentText(parsed.Comment)}
	for _, s := range parsed.Steps {
		step := Step{Name: s.Name, Program: s.Program, Description: cobol.CommentText(s.Comment),
			Source: fmt.Sprintf("%s:%d", path, s.Line)}
		if s.Cond != nil {
			step.Cond = fmt.Sprintf("(%d,%s)", s.Cond.Code, s.Cond.Operator)
		}
		i := slices.IndexFunc(programs, func(p Program) bool { return p.Name == s.Program })
		if i < 0 {
			job.Findings = append(job.Findings, fmt.Sprintf("Step %s runs %s, which is not among the programs analyzed.", s.Name, s.Program))
		}
		for _, d := range s.DDs {
			dd := newDD(path, d)
			if i >= 0 {
				job.Findings = append(job.Findings, bind(&dd, s.Name, programs[i])...)
			}
			step.DDs = append(step.DDs, dd)
		}
		if i >= 0 {
			job.Findings = append(job.Findings, unbound(step, programs[i])...)
		}
		job.Steps = append(job.Steps, step)
	}
	job.Findings = append(job.Findings, temporaryFindings(job)...)
	return job, nil
}

func newDD(path string, d jcl.DD) DD {
	dd := DD{Name: d.Name, Disp: string(d.Disp), Sysout: d.Sysout, Dummy: d.Dummy, Source: fmt.Sprintf("%s:%d", path, d.Line)}
	for _, member := range append([]jcl.DD{d}, d.Concatenated...) {
		if member.DSN != "" {
			dd.Datasets = append(dd.Datasets, datasetName(member))
		}
	}
	return dd
}

// datasetName writes a DD's dataset as the job names it: ACCT.MASTER(+1),
// &&SORTED, or every generation of a group by its base, TRAN.AUDIT.
func datasetName(d jcl.DD) string {
	switch {
	case d.Temporary:
		return "&&" + d.DSN
	case d.Relative && d.Generation > 0:
		return fmt.Sprintf("%s(+%d)", d.DSN, d.Generation)
	case d.Relative:
		return fmt.Sprintf("%s(%d)", d.DSN, d.Generation)
	}
	return d.DSN
}

// bind finds the program's file assigned to the DD and how the program
// uses it, and reports a DD no file is assigned to and a use its
// disposition does not allow.
func bind(dd *DD, step string, p Program) []string {
	i := slices.IndexFunc(p.Files, func(f File) bool { return f.Assign == dd.Name })
	if i < 0 {
		return []string{fmt.Sprintf("Step %s binds DD %s, which %s assigns no file to.", step, dd.Name, p.Name)}
	}
	f := p.Files[i]
	dd.File = f.Name
	var modes []string
	for _, mode := range f.Opened {
		modes = append(modes, uses[mode])
	}
	dd.Use = strings.Join(modes, ", ")
	if len(dd.Datasets) == 0 {
		return nil
	}
	writes := slices.Contains(modes, useWrite) || slices.Contains(modes, useAppend) || slices.Contains(modes, useUpdate)
	switch {
	case writes && (dd.Disp == "SHR" || dd.Disp == ""):
		return []string{fmt.Sprintf("Step %s: %s writes DD %s, a dataset it shares with DISP=SHR.", step, p.Name, dd.Name)}
	case !writes && len(modes) > 0 && (dd.Disp == "NEW" || dd.Disp == "MOD"):
		return []string{fmt.Sprintf("Step %s: %s only reads DD %s, which it creates with DISP=%s.", step, p.Name, dd.Name, dd.Disp)}
	case slices.Contains(modes, useAppend) && dd.Disp != "MOD":
		return []string{fmt.Sprintf("Step %s: %s appends to DD %s, which is not DISP=MOD.", step, p.Name, dd.Name)}
	}
	return nil
}

// unbound reports the program's files that the step binds no DD to. A
// sort's work file needs none.
func unbound(step Step, p Program) []string {
	var findings []string
	for _, f := range p.Files {
		if f.Kind == "SD" || slices.ContainsFunc(step.DDs, func(dd DD) bool { return dd.Name == f.Assign }) {
			continue
		}
		findings = append(findings, fmt.Sprintf("Step %s binds no DD %s for %s's file %s.", step.Name, f.Assign, p.Name, f.Name))
	}
	return findings
}

// temporaryFindings reports temporary datasets read before any step writes
// them, and those written but never read.
func temporaryFindings(job Job) []string {
	var findings []string
	written, read := map[string]bool{}, map[string]bool{}
	var order []string
	for _, step := range job.Steps {
		for _, dd := range step.DDs {
			for _, ds := range dd.Datasets {
				if !strings.HasPrefix(ds, "&&") {
					continue
				}
				if !slices.Contains(order, ds) {
					order = append(order, ds)
				}
				switch {
				case dd.Use == useRead && !written[ds]:
					findings = append(findings, fmt.Sprintf("Step %s reads %s before any step writes it.", step.Name, ds))
				case dd.Use == useRead:
					read[ds] = true
				case dd.Use != "":
					written[ds] = true
				}
			}
		}
	}
	for _, ds := range order {
		if written[ds] && !read[ds] {
			findings = append(findings, fmt.Sprintf("%s is written and never read.", ds))
		}
	}
	return findings
}

// jobDiagram is the graph of a job's steps and the datasets between them:
// an arrow from each dataset a step reads to the step, and from the step
// to each dataset it writes, and a dotted one from each step to the next.
// A dataset that one step writes and another reads is drawn once, between
// them; one the job only reads is drawn above each step that reads it, and
// one it only writes below each step that writes it, so that the flow
// through the job is not lost among arrows from far away. The job's output
// and dummy DDs are left out.
func jobDiagram(dir string, job Job) diagram.Graph {
	return JobGraph(job, JobNotes{Link: func(step Step) string { return sourceTarget(dir, step.Source) }})
}

// JobNotes add to a job's diagram what a run of it did: a note after each
// step's label, such as its return code, and after the label of each
// dataset a DD names alone, such as its records; and the link of each
// step. Any of them may be nil.
type JobNotes struct {
	Step func(Step) string
	DD   func(Step, DD) string
	Link func(Step) string
}

// JobGraph is the graph of a job's network, as in its blueprint, with the
// notes added.
func JobGraph(job Job, notes JobNotes) diagram.Graph {
	read, written := map[string]bool{}, map[string]bool{}
	for _, step := range job.Steps {
		for _, dd := range step.DDs {
			for _, ds := range dd.Datasets {
				if dd.Use == useRead {
					read[ds] = true
				} else if dd.Use != "" {
					written[ds] = true
				}
			}
		}
	}
	node := func(ds, step string) string {
		if read[ds] && written[ds] {
			return "dataset " + ds
		}
		return "dataset " + ds + " at " + step
	}
	var g diagram.Graph
	added := map[string]bool{}
	addDataset := func(ds string, step Step, dd DD) string {
		id := node(ds, step.Name)
		if !added[id] {
			added[id] = true
			label := ds
			if notes.DD != nil && len(dd.Datasets) == 1 {
				label = withNote(label, notes.DD(step, dd))
			}
			g.Nodes = append(g.Nodes, diagram.Node{ID: id, Label: label, Rounded: true})
		}
		return id
	}
	for i, step := range job.Steps {
		id := "step " + step.Name
		for _, dd := range step.DDs {
			if dd.Use == useRead {
				for _, ds := range dd.Datasets {
					g.Edges = append(g.Edges, diagram.Edge{From: addDataset(ds, step, dd), To: id})
				}
			}
		}
		node := diagram.Node{ID: id, Label: step.Name + " " + step.Program, Title: strings.Join(step.Description, " ")}
		if notes.Step != nil {
			node.Label = withNote(node.Label, notes.Step(step))
		}
		if notes.Link != nil {
			node.Link = notes.Link(step)
		}
		g.Nodes = append(g.Nodes, node)
		if i > 0 {
			g.Edges = append(g.Edges, diagram.Edge{From: "step " + job.Steps[i-1].Name, To: id, Style: diagram.Dotted})
		}
		for _, dd := range step.DDs {
			if dd.Use != useRead && dd.Use != "" {
				for _, ds := range dd.Datasets {
					g.Edges = append(g.Edges, diagram.Edge{From: id, To: addDataset(ds, step, dd)})
				}
			}
		}
	}
	return g
}

func withNote(label, note string) string {
	if note == "" {
		return label
	}
	return label + " " + note
}

func jobMarkdown(dir string, job Job) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", job.Name)
	writeDescription(&b, job.Description)
	fmt.Fprintf(&b, "Job %s: %d steps.\n", sourceLink(dir, job.Source, false), len(job.Steps))
	fmt.Fprintf(&b, "\n## Network\n\n![The steps of %s and the datasets they read and write](%s.svg)\n", job.Name, job.Name)
	b.WriteString("\nBoxes are steps, with the program each runs; rounded boxes are datasets, as the job names them. " +
		"Arrows run from the datasets a step reads to the step, and from the step to the datasets it writes; " +
		"dotted ones from each step to the next.\n")

	b.WriteString("\n## Steps\n\n| Step | Program | Bypassed by | Description | Source |\n| --- | --- | --- | --- | --- |\n")
	for _, step := range job.Steps {
		cond := ""
		if step.Cond != "" {
			cond = "COND=" + step.Cond
		}
		b.WriteString(row(step.Name, fmt.Sprintf("[%s](../programs/%s.md)", step.Program, step.Program), cond,
			strings.Join(step.Description, " "), sourceLink(dir, step.Source, true)))
	}
	for _, step := range job.Steps {
		fmt.Fprintf(&b, "\n### %s\n\n| DD | Dataset | DISP | File | Use |\n| --- | --- | --- | --- | --- |\n", step.Name)
		for _, dd := range step.DDs {
			datasets := strings.Join(dd.Datasets, ", ")
			switch {
			case dd.Sysout:
				datasets = "job output"
			case dd.Dummy:
				datasets = "dummy"
			}
			b.WriteString(row(dd.Name, datasets, dd.Disp, dd.File, dd.Use))
		}
	}
	b.WriteString("\n## Findings\n\n")
	if len(job.Findings) == 0 {
		b.WriteString("None: every DD is assigned a file, every file a DD, and every disposition fits its use.\n")
	}
	for _, f := range job.Findings {
		b.WriteString("- " + f + "\n")
	}
	return b.String()
}
