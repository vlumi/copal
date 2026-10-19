package oracle

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/dataset"
	"github.com/vlumi/copal/bedrock/jcl"
	"github.com/vlumi/copal/bedrock/run"
	jobs "github.com/vlumi/copal/bedrock/runner"
	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/coverage"
)

// Batch is one run of a job, after the jobs before it in a batch, on what
// they wrote, and on inputs taken from the cases of its steps' programs:
// each dataset the job reads that no job before it, and none of its steps
// before, wrote, with the records of the cleanest case of the first program
// that reads it, and the in-stream data of a step with the lines its
// program's cleanest case accepts. It records the steps' outcomes, the datasets the job
// wrote, its output, and the rules each program reached.
type Batch struct {
	Job     string
	inputs  []batchInput
	steps   []jobs.Result
	outputs map[string][]byte // by dataset
	sysout  map[string][]byte // by file in the job's output
	reached map[string]map[rule]bool
	results map[string]*Result
	about   string // how the batch was run
	ended   string // how it ended, if that is more than its steps say
}

// batchInput is a dataset a job reads from outside it, the program whose
// cleanest case gave its records, and the input of the program it is.
type batchInput struct {
	dsn     string
	reader  *Result
	input   input
	records [][]byte
}

// RunBatch runs a batch of jobs, or of scripts, as RunJobs or RunScripts
// runs them.
func RunBatch(system *blueprint.System, jobs []*blueprint.Job, results map[string]*Result, scratch string) ([]*Batch, error) {
	scripts := 0
	for _, job := range jobs {
		if job.Kind == "script" {
			scripts++
		}
	}
	switch scripts {
	case 0:
		return RunJobs(jobs, results, scratch)
	case len(jobs):
		return RunScripts(system, jobs, results, scratch)
	}
	return nil, errors.New("a batch of both jobs and scripts")
}

// RunJobs runs jobs, in the order given, on one catalog of datasets, in a
// scratch directory the programs' cases share: each job reads what the
// jobs before it wrote, and takes from the cases of its steps' programs
// only the datasets none of them did.
func RunJobs(jobs []*blueprint.Job, results map[string]*Result, scratch string) ([]*Batch, error) {
	root := filepath.Join(scratch, "jobs")
	if err := os.RemoveAll(root); err != nil {
		return nil, err
	}
	store, err := dataset.Open(filepath.Join(root, "datasets"))
	if err != nil {
		return nil, err
	}
	var batches []*Batch
	for _, job := range jobs {
		b, err := runJob(job, results, store, scratch, filepath.Join(root, job.Name))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", job.Name, err)
		}
		batches = append(batches, b)
	}
	return batches, nil
}

func runJob(job *blueprint.Job, results map[string]*Result, store *dataset.Store, scratch, dir string) (*Batch, error) {
	parsed, err := jcl.ParseFile(job.Source)
	if err != nil {
		return nil, err
	}
	b := &Batch{Job: job.Name, results: results, reached: map[string]map[rule]bool{}, about: jobAbout}
	if err := b.provide(parsed, store, dir); err != nil {
		return nil, err
	}
	before, err := snapshot(store)
	if err != nil {
		return nil, err
	}
	traces := filepath.Join(dir, "traces")
	if err := os.MkdirAll(traces, 0o755); err != nil {
		return nil, err
	}
	r := jobs.Runner{Programs: filepath.Join(scratch, "programs"), Datasets: store, Sysout: filepath.Join(dir, "sysout"),
		Traces: traces, InStream: bytes.NewReader(b.inStream(parsed))}
	if b.steps, err = r.Run(parsed); err != nil {
		return nil, err
	}
	if err := b.replay(parsed, traces); err != nil {
		return nil, err
	}
	if b.outputs, err = b.written(parsed, store, before, dir); err != nil {
		return nil, err
	}
	b.sysout, err = readTree(r.Sysout)
	return b, err
}

// reading is a DD of a step that reads a dataset.
type reading struct {
	step jcl.Step
	dd   jcl.DD
}

// provide catalogs the datasets the job reads from outside it: for each,
// as many generations as its readers go back, each holding the records of
// the cleanest case of the first program that reads it.
func (b *Batch) provide(job *jcl.Job, store *dataset.Store, dir string) error {
	written := map[string]bool{}
	readers := map[string][]reading{}
	var order []string
	for _, step := range job.Steps {
		for _, dd := range step.DDs {
			for _, d := range append([]jcl.DD{dd}, dd.Concatenated...) {
				name := datasetKey(d)
				if d.Sysout || d.Dummy || d.InStream || written[name] || cataloged(store, d) {
					continue
				}
				if d.Disp == jcl.Shr || d.Disp == jcl.Old {
					if _, ok := readers[name]; !ok {
						order = append(order, name)
					}
					readers[name] = append(readers[name], reading{step, d})
				}
			}
		}
		for _, dd := range step.DDs {
			if dd.Disp == jcl.New || dd.Disp == jcl.Mod {
				written[datasetKey(dd)] = true
			}
		}
	}
	for _, name := range order {
		if err := b.provideDataset(readers[name], store, dir); err != nil {
			return err
		}
	}
	return nil
}

// cataloged reports whether a job before has written a dataset a DD
// reads: it, or a generation of it.
func cataloged(store *dataset.Store, dd jcl.DD) bool {
	if _, ok := store.LRECL(dd.DSN); ok {
		return true
	}
	return !dd.Temporary && len(store.Generations(dd.DSN)) > 0
}

// datasetKey names a dataset without its generation.
func datasetKey(dd jcl.DD) string {
	if dd.Temporary {
		return "&&" + dd.DSN
	}
	return dd.DSN
}

func (b *Batch) provideDataset(readers []reading, store *dataset.Store, dir string) error {
	first := readers[0]
	result := b.results[first.step.Program]
	if result == nil {
		return nil
	}
	in, f, ok := result.program.inputFor(first.dd.Name)
	if !ok {
		return nil
	}
	var records [][]byte
	if c, ok := result.cleanest(); ok {
		records = c[in.name()]
	}
	names := []string{first.dd.DSN}
	if first.dd.Relative || first.dd.Generation != 0 {
		back := 0
		for _, r := range readers {
			back = max(back, -r.dd.Generation)
		}
		names = nil
		for n := 1; n <= back+1; n++ {
			names = append(names, dataset.GenerationName(first.dd.DSN, n))
		}
	}
	for _, dsn := range names {
		var err error
		if f.Organization == "INDEXED" {
			if err = result.runner.load(f, records, store.Path(dsn), dir); err == nil {
				err = store.CatalogIndexed(dsn, in.root.Size())
			}
		} else {
			err = store.Write(dsn, in.root.Size(), records)
		}
		if err != nil {
			return fmt.Errorf("%s: %w", dsn, err)
		}
		b.inputs = append(b.inputs, batchInput{dsn: dsn, reader: result, input: in, records: records})
	}
	return nil
}

// cleanest is the case of a program whose run ended best: the first that
// ended with the lowest return code, an abend counting as worst.
func (r *Result) cleanest() (Case, bool) {
	best := -1
	for i, run := range r.runs {
		if run.outcome.Abended {
			continue
		}
		if best < 0 || run.outcome.RC < r.runs[best].outcome.RC {
			best = i
		}
	}
	if best < 0 {
		return nil, false
	}
	return r.Cases[best], true
}

// inputFor finds the input of a program a DD binds, and its file.
func (p *program) inputFor(dd string) (input, blueprint.File, bool) {
	for _, in := range p.inputs {
		if in.file.Assign == dd {
			return in, in.file, true
		}
	}
	return input{}, blueprint.File{}, false
}

// inStream is the job's in-stream data: the lines the cleanest case of the
// program of the step that reads it accepts.
func (b *Batch) inStream(job *jcl.Job) []byte {
	for _, step := range job.Steps {
		if !slices.ContainsFunc(step.DDs, func(dd jcl.DD) bool { return dd.InStream }) {
			continue
		}
		if r := b.results[step.Program]; r != nil {
			if c, ok := r.cleanest(); ok {
				return lines(c[sysin])
			}
		}
	}
	return nil
}

// replay replays each step's trace against its program's decision tables,
// and those of the programs it calls.
func (b *Batch) replay(job *jcl.Job, traces string) error {
	for _, step := range job.Steps {
		result := b.results[step.Program]
		data, err := os.ReadFile(filepath.Join(traces, step.Name+".trace"))
		if result == nil || os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		events, err := coverage.Read(bytes.NewReader(data))
		if err != nil {
			return err
		}
		if b.reached[step.Program] == nil {
			b.reached[step.Program] = map[rule]bool{}
		}
		for u, unit := range result.Units {
			for i, t := range coverage.Replay(unit.source, unit.Tables, events) {
				for j, n := range t.Reached {
					if n > 0 {
						b.reached[step.Program][rule{u, i, j}] = true
					}
				}
			}
		}
	}
	return nil
}

// snapshot is the contents of each dataset in a catalog.
func snapshot(store *dataset.Store) (map[string][]byte, error) {
	contents := map[string][]byte{}
	for _, dsn := range store.Names() {
		data, err := os.ReadFile(store.Path(dsn))
		if err != nil {
			return nil, err
		}
		contents[dsn] = data
	}
	return contents, nil
}

// written reads the datasets the job created or changed, a record a line,
// an indexed one through the unloader of a program that writes it.
func (b *Batch) written(job *jcl.Job, store *dataset.Store, before map[string][]byte, dir string) (map[string][]byte, error) {
	outputs := map[string][]byte{}
	for _, dsn := range store.Names() {
		data, err := os.ReadFile(store.Path(dsn))
		if err != nil {
			return nil, err
		}
		if old, ok := before[dsn]; ok && bytes.Equal(old, data) {
			continue
		}
		lrecl, _ := store.LRECL(dsn)
		if !store.Indexed(dsn) {
			outputs[dsn] = asLines(data, false, lrecl)
			continue
		}
		data, err = b.unload(job, dsn, store.Path(dsn), dir)
		if err != nil {
			return nil, err
		}
		outputs[dsn] = asLines(data, false, lrecl)
	}
	return outputs, nil
}

// unload reads an indexed dataset in key order through the unloader of a
// program whose step binds it.
func (b *Batch) unload(job *jcl.Job, dsn, path, dir string) ([]byte, error) {
	for _, step := range job.Steps {
		result := b.results[step.Program]
		if result == nil {
			continue
		}
		for _, dd := range step.DDs {
			if dsn != dd.DSN && !strings.HasPrefix(dsn, dd.DSN+".G") {
				continue
			}
			i := slices.IndexFunc(result.program.blueprint.Files, func(f blueprint.File) bool { return f.Assign == dd.Name })
			if i < 0 || result.program.blueprint.Files[i].Organization != "INDEXED" {
				continue
			}
			out := filepath.Join(dir, "unloaded")
			outcome, err := run.Step{Programs: result.runner.loaders, Program: unloaderName(result.program.blueprint.Files[i]),
				DDs: map[string]string{"UNLOADIN": path, "UNLOADOUT": out}, Stdout: io.Discard, Timeout: timeout}.Run()
			if err == nil && (outcome.RC != 0 || outcome.Abended) {
				err = fmt.Errorf("unloading %s: %+v", dsn, outcome)
			}
			if err != nil {
				return nil, err
			}
			return os.ReadFile(out)
		}
	}
	return nil, fmt.Errorf("%s: no program to unload it with", dsn)
}

// readTree reads the files of a directory, by their path in it.
func readTree(dir string) (map[string][]byte, error) {
	files := map[string][]byte{}
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		files[filepath.ToSlash(rel)] = data
		return nil
	})
	if os.IsNotExist(err) {
		return files, nil
	}
	return files, err
}

// Files writes the batch as files, by their path in the job's directory:
// README.md, how each step ended and the rules each program reached;
// input/DSN.txt, each dataset from outside the job, field by field;
// output/DSN.txt, each dataset it wrote, a record a line; and sysout/, the
// job's output, its log among it.
func (b *Batch) Files() map[string][]byte {
	files := map[string][]byte{"README.md": []byte(b.readme())}
	for _, in := range b.inputs {
		files["input/"+in.dsn+".txt"] = in.reader.program.encode(in.input, in.records)
	}
	for dsn, data := range b.outputs {
		files["output/"+dsn+".txt"] = data
	}
	for name, data := range b.sysout {
		files["sysout/"+name+".txt"] = data
	}
	return files
}

func (b *Batch) readme() string {
	var s strings.Builder
	fmt.Fprintf(&s, "# %s: a batch\n\n%s\n\n", b.Job, b.about)
	if b.ended != "" {
		s.WriteString(b.ended + "\n\n")
	}
	s.WriteString("| Step | Program | Outcome | Rules reached |\n| --- | --- | --- | --: |\n")
	for _, step := range b.steps {
		reached := ""
		if r := b.results[step.Program]; r != nil && !step.Bypassed {
			n, total := 0, 0
			for u, unit := range r.Units {
				for i, t := range unit.Tables {
					total += len(t.Rules)
					for j := range t.Rules {
						if b.reached[step.Program][rule{u, i, j}] {
							n++
						}
					}
				}
			}
			reached = fmt.Sprintf("%d of %d", n, total)
		}
		program := step.Program
		if b.results[step.Program] != nil {
			program = fmt.Sprintf("[%s](../../%s/)", step.Program, step.Program)
		}
		if reached == "" {
			fmt.Fprintf(&s, "| %s | %s | %s | |\n", step.Step, program, step.Outcome())
			continue
		}
		fmt.Fprintf(&s, "| %s | %s | %s | %s |\n", step.Step, program, step.Outcome(), reached)
	}
	return s.String()
}

const jobAbout = "The job run after the jobs before it in the batch, on what they wrote, and on inputs taken from the cases of " +
	"its steps' programs for the rest: each dataset it reads that no job before it, and no step of it before, wrote, with the records of the cleanest case of the first program that reads it, the first that ended with " +
	"the lowest return code, and its in-stream data with the lines its program's cleanest case accepts. The datasets it wrote and its output are recorded beside."
