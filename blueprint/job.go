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
	Name        string    `json:"job"`
	Source      string    `json:"source"`
	Description []string  `json:"description,omitempty"`
	Steps       []Step    `json:"steps"`
	Datasets    []Dataset `json:"datasets"`
	Findings    []string  `json:"findings,omitempty"`
}

// Dataset is a dataset a job names, as it names it, with the copybook
// whose record describes it, as far as the programs that use it say, and
// the steps that write and read it.
type Dataset struct {
	Name      string   `json:"dataset"`
	Copybook  string   `json:"copybook,omitempty"`
	WrittenBy []string `json:"writtenBy,omitempty"`
	ReadBy    []string `json:"readBy,omitempty"`
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
	Uses     []string `json:"uses,omitempty"`
	// Copybook is the copybook of the record the program declares for the
	// file, when it declares it with one.
	Copybook string `json:"copybook,omitempty"`
	Source   string `json:"source"`
}

// Uses of a dataset, from the modes the program opens its file in.
const (
	useRead   = "read"
	useWrite  = "write"
	useAppend = "append"
	useUpdate = "update"
)

var uses = map[string]string{"INPUT": useRead, "OUTPUT": useWrite, "EXTEND": useAppend, "I-O": useUpdate}

// reads reports whether the step reads the DD's datasets.
func (dd DD) reads() bool {
	return slices.Contains(dd.Uses, useRead) || slices.Contains(dd.Uses, useUpdate)
}

// writes reports whether the step writes the DD's datasets, from scratch,
// by appending or by updating them.
func (dd DD) writes() bool {
	return slices.ContainsFunc(dd.Uses, func(u string) bool { return u != useRead })
}

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
	job.Datasets = jobDatasets(job)
	return job, nil
}

// jobDatasets lists the datasets the job names, in the order it first
// names them, with the steps that write and read each.
func jobDatasets(job Job) []Dataset {
	var datasets []Dataset
	for _, step := range job.Steps {
		for _, dd := range step.DDs {
			for _, name := range dd.Datasets {
				i := slices.IndexFunc(datasets, func(d Dataset) bool { return d.Name == name })
				if i < 0 {
					datasets = append(datasets, Dataset{Name: name})
					i = len(datasets) - 1
				}
				if dd.writes() {
					addOnce(&datasets[i].WrittenBy, step.Name)
				}
				if dd.reads() {
					addOnce(&datasets[i].ReadBy, step.Name)
				}
			}
		}
	}
	return datasets
}

// datasetBase is the dataset a name stands for, whatever its generation:
// ACCT.MASTER for ACCT.MASTER(+1).
func datasetBase(name string) string {
	base, _, _ := strings.Cut(name, "(")
	return base
}

// DatasetCopybooks maps each dataset the system's jobs name, by its name
// without generation, to the copybook that describes its records, where it
// is known. Temporary datasets are left out.
func (s *System) DatasetCopybooks() map[string]string {
	copybooks := map[string]string{}
	for _, job := range s.Jobs {
		for _, d := range job.Datasets {
			if d.Copybook != "" && !strings.HasPrefix(d.Name, "&&") {
				copybooks[datasetBase(d.Name)] = d.Copybook
			}
		}
	}
	return copybooks
}

// describeDatasets names each dataset's record layout from every job: the
// copybook of the record a program declares for a file bound to the
// dataset, in any generation. Two programs that describe one dataset with
// different copybooks are a finding of each job that names it.
func describeDatasets(jobs []Job) {
	copybooks := map[string][]string{}
	for _, job := range jobs {
		for _, step := range job.Steps {
			for _, dd := range step.DDs {
				for _, name := range dd.Datasets {
					if dd.Copybook != "" {
						list := copybooks[datasetBase(name)]
						addOnce(&list, dd.Copybook)
						copybooks[datasetBase(name)] = list
					}
				}
			}
		}
	}
	for j := range jobs {
		reported := map[string]bool{}
		for i, d := range jobs[j].Datasets {
			base := datasetBase(d.Name)
			list := copybooks[base]
			if len(list) > 0 {
				jobs[j].Datasets[i].Copybook = list[0]
			}
			if len(list) > 1 && !reported[base] {
				reported[base] = true
				jobs[j].Findings = append(jobs[j].Findings,
					fmt.Sprintf("%s is described by more than one copybook: %s.", base, strings.Join(list, ", ")))
			}
		}
	}
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
	dd.Copybook = fileCopybook(p, f)
	for _, mode := range f.Opened {
		dd.Uses = append(dd.Uses, uses[mode])
	}
	if len(dd.Datasets) == 0 {
		return nil
	}
	switch {
	case dd.writes() && (dd.Disp == "SHR" || dd.Disp == ""):
		return []string{fmt.Sprintf("Step %s: %s writes DD %s, a dataset it shares with DISP=SHR.", step, p.Name, dd.Name)}
	case dd.reads() && !dd.writes() && (dd.Disp == "NEW" || dd.Disp == "MOD"):
		return []string{fmt.Sprintf("Step %s: %s only reads DD %s, which it creates with DISP=%s.", step, p.Name, dd.Name, dd.Disp)}
	case slices.Contains(dd.Uses, useAppend) && dd.Disp != "MOD":
		return []string{fmt.Sprintf("Step %s: %s appends to DD %s, which is not DISP=MOD.", step, p.Name, dd.Name)}
	}
	return nil
}

// fileCopybook is the copybook of a file's record: the record declared
// for the file, or else the one the program writes it from or reads it into.
func fileCopybook(p Program, f File) string {
	copybooks := map[string]string{}
	for _, r := range p.Records {
		if r.Copybook == "" {
			continue
		}
		if r.File == f.Name {
			return r.Copybook
		}
		copybooks[r.Name] = r.Copybook
	}
	for _, e := range p.Lineage {
		switch {
		case (e.Verb == "WRITE" || e.Verb == "RELEASE") && slices.Contains(f.Records, e.To) && copybooks[e.From] != "":
			return copybooks[e.From]
		case (e.Verb == "READ" || e.Verb == "RETURN") && slices.Contains(f.Records, e.From) && copybooks[e.To] != "":
			return copybooks[e.To]
		}
	}
	return ""
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
				if dd.reads() {
					if !written[ds] {
						findings = append(findings, fmt.Sprintf("Step %s reads %s before any step writes it.", step.Name, ds))
					}
					read[ds] = true
				}
				if dd.writes() {
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
// to each dataset it writes, and a dotted one from each step to the next
// where no dataset leads from one to the other.
// A dataset that one step writes and another reads is drawn once, between
// them; one the job only reads is drawn above each step that reads it, and
// one it only writes below each step that writes it, so that the flow
// through the job is not lost among arrows from far away. The job's output
// and dummy DDs are left out.
func jobDiagram(job Job) diagram.Graph {
	return JobGraph(job, JobNotes{
		ProgramLink:  func(name string) string { return "../programs/" + name + ".md" },
		CopybookLink: func(name string) string { return "../copybooks/" + name + ".md" },
	})
}

// networkKeys say how to read a job's network.
var networkKeys = []diagram.Key{
	diagram.NodeKey(diagram.Node{}, "step, the program it runs"),
	diagram.NodeKey(diagram.Node{Data: true}, "dataset, its record's copybook"),
	diagram.NodeKey(diagram.Node{Data: true, Dashed: true}, "temporary dataset, &&NAME"),
	diagram.EdgeKey(diagram.Solid, "read or written"),
	diagram.EdgeKey(diagram.Dotted, "runs next"),
	diagram.TextKey("NAME(0)  current generation"),
	diagram.TextKey("NAME(+1) generation written"),
}

// JobNotes add to a job's diagram what a run of it did, a line under each
// step's label, such as its return code, and under the label of each
// dataset a DD names alone, such as its records; and say where steps and
// copybooks link to, by name. Any of them may be nil.
type JobNotes struct {
	Step         func(Step) string
	DD           func(Step, DD) string
	ProgramLink  func(string) string
	CopybookLink func(string) string
}

// JobGraph is the graph of a job's network, as in its blueprint, with the
// notes added. A step shows its program, and a dataset the copybook that
// describes its records, when it is known.
func JobGraph(job Job, notes JobNotes) diagram.Graph {
	copybooks := map[string]string{}
	for _, d := range job.Datasets {
		copybooks[d.Name] = d.Copybook
	}
	g := diagram.Graph{Legend: networkKeys}
	read, written := map[string]bool{}, map[string]bool{}
	for _, step := range job.Steps {
		for _, dd := range step.DDs {
			for _, ds := range dd.Datasets {
				read[ds] = read[ds] || dd.reads()
				written[ds] = written[ds] || dd.writes()
			}
		}
	}
	node := func(ds, step string) string {
		if read[ds] && written[ds] {
			return "dataset " + ds
		}
		return "dataset " + ds + " at " + step
	}
	added := map[string]bool{}
	addDataset := func(ds string, step Step, dd DD) string {
		id := node(ds, step.Name)
		if !added[id] {
			added[id] = true
			n := diagram.Node{ID: id, Label: ds, Data: true, Dashed: strings.HasPrefix(ds, "&&")}
			if c := copybooks[ds]; c != "" {
				n.Label += "\ncopybook " + c
				if notes.CopybookLink != nil {
					n.Link = notes.CopybookLink(c)
				}
			}
			if notes.DD != nil && len(dd.Datasets) == 1 {
				n.Label = withNote(n.Label, notes.DD(step, dd))
			}
			g.Nodes = append(g.Nodes, n)
		}
		return id
	}
	for i, step := range job.Steps {
		id := "step " + step.Name
		for _, dd := range step.DDs {
			if dd.reads() {
				for _, ds := range dd.Datasets {
					g.Edges = append(g.Edges, diagram.Edge{From: addDataset(ds, step, dd), To: id})
				}
			}
		}
		node := diagram.Node{ID: id, Label: step.Name + "\nprogram " + step.Program, Title: strings.Join(step.Description, " ")}
		if notes.Step != nil {
			node.Label = withNote(node.Label, notes.Step(step))
		}
		if notes.ProgramLink != nil {
			node.Link = notes.ProgramLink(step.Program)
		}
		g.Nodes = append(g.Nodes, node)
		if i > 0 && !feeds(job.Steps[i-1], step, read, written) {
			g.Edges = append(g.Edges, diagram.Edge{From: "step " + job.Steps[i-1].Name, To: id, Style: diagram.Dotted, Span: 2})
		}
		for _, dd := range step.DDs {
			if dd.writes() {
				for _, ds := range dd.Datasets {
					g.Edges = append(g.Edges, diagram.Edge{From: id, To: addDataset(ds, step, dd)})
				}
			}
		}
	}
	return g
}

// feeds reports whether a step writes a dataset that the next one reads,
// one drawn once between them.
func feeds(step, next Step, read, written map[string]bool) bool {
	for _, out := range step.DDs {
		if !out.writes() {
			continue
		}
		for _, ds := range out.Datasets {
			if !read[ds] || !written[ds] {
				continue
			}
			for _, in := range next.DDs {
				if in.reads() && slices.Contains(in.Datasets, ds) {
					return true
				}
			}
		}
	}
	return false
}

func copybookLink(name string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("[%s](../copybooks/%s.md)", name, name)
}

func addOnce(names *[]string, name string) {
	if !slices.Contains(*names, name) {
		*names = append(*names, name)
	}
}

func withNote(label, note string) string {
	if note == "" {
		return label
	}
	return label + "\n" + note
}

func jobMarkdown(dir string, job Job) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", job.Name)
	writeDescription(&b, job.Description)
	fmt.Fprintf(&b, "Job %s: %d steps.\n", sourceLink(dir, job.Source, false), len(job.Steps))
	fmt.Fprintf(&b, "\n## Network\n\n![The steps of %s and the datasets they read and write](%s.svg)\n", job.Name, job.Name)

	b.WriteString("\n## Steps\n\n| Step | Program | Reads | Writes | Bypassed by | Description | Source |\n" +
		"| --- | --- | --- | --- | --- | --- | --- |\n")
	for _, step := range job.Steps {
		cond := ""
		if step.Cond != "" {
			cond = "COND=" + step.Cond
		}
		var reads, writes []string
		for _, dd := range step.DDs {
			if dd.reads() {
				reads = append(reads, dd.Datasets...)
			}
			if dd.writes() {
				writes = append(writes, dd.Datasets...)
			}
		}
		b.WriteString(row(step.Name, fmt.Sprintf("[%s](../programs/%s.md)", step.Program, step.Program),
			strings.Join(reads, ", "), strings.Join(writes, ", "), cond,
			strings.Join(step.Description, " "), sourceLink(dir, step.Source, true)))
	}
	for _, step := range job.Steps {
		fmt.Fprintf(&b, "\n### %s\n\n| DD | Dataset | DISP | File | Copybook | Use |\n| --- | --- | --- | --- | --- | --- |\n", step.Name)
		for _, dd := range step.DDs {
			datasets := strings.Join(dd.Datasets, ", ")
			switch {
			case dd.Sysout:
				datasets = "job output"
			case dd.Dummy:
				datasets = "dummy"
			}
			b.WriteString(row(dd.Name, datasets, dd.Disp, dd.File, copybookLink(dd.Copybook), strings.Join(dd.Uses, ", ")))
		}
	}
	b.WriteString("\n## Datasets\n\nEach dataset the job names, the copybook that describes its records, " +
		"from the programs that read or write it in any job, and the steps that write and read it.\n\n" +
		"| Dataset | Copybook | Written by | Read by |\n| --- | --- | --- | --- |\n")
	for _, d := range job.Datasets {
		b.WriteString(row(d.Name, copybookLink(d.Copybook), strings.Join(d.WrittenBy, ", "), strings.Join(d.ReadBy, ", ")))
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
