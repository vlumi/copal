package oracle

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vlumi/copal/bedrock/layout"
	"github.com/vlumi/copal/survey/blueprint"
)

// Files writes a program's cases and its coverage as files, by their path
// in the program's directory: README.md, the coverage, rule by rule, with
// links to the decision tables in rulesDir, a path relative to the
// program's directory; coverage.json, the same for tools; and for each
// case, NN/README.md, what it returned
// and the rules it reaches first, NN/input/FILE.txt, its input files,
// field by field, NN/input/failed.txt, the file it fails and how, if it
// fails one, and NN/output/FILE.txt, what it wrote, a record a line.
func (r *Result) Files(rulesDir string) map[string][]byte {
	files := map[string][]byte{"README.md": []byte(r.coverage(rulesDir)), "coverage.json": r.coverageJSON()}
	for i, c := range r.Cases {
		dir := caseName(i)
		files[dir+"/README.md"] = []byte(r.caseReadme(i, "../"+rulesDir))
		for _, in := range r.program.inputs {
			if _, fails := c.failsFile(in.name()); !fails {
				files[dir+"/input/"+in.name()+".txt"] = r.program.encode(in, c[in.name()])
			}
		}
		if len(c[failed]) > 0 {
			files[dir+"/input/"+failedName] = lines(c[failed])
		}
		for name, data := range r.runs[i].outputs {
			files[dir+"/output/"+name+".txt"] = data
		}
		if len(r.runs[i].stdout) > 0 {
			files[dir+"/output/STDOUT.txt"] = asLines(r.runs[i].stdout, true, 0)
		}
	}
	return files
}

func caseName(i int) string { return fmt.Sprintf("%02d", i+1) }

func (r *Result) coverage(rulesDir string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s: cases\n\n", r.Program)
	reached, total := r.Rules(0)
	fmt.Fprintf(&b, "%d cases reach %d of the %d rules of %s's decision tables", len(r.Cases), reached, total, r.Program)
	var called []string
	for u := 1; u < len(r.Units); u++ {
		reached, total := r.Rules(u)
		called = append(called, fmt.Sprintf("%d of the %d of %s", reached, total, r.Units[u].Program))
	}
	if len(called) > 0 {
		fmt.Fprintf(&b, ", and of the programs it calls, %s", strings.Join(called, ", "))
	}
	b.WriteString(". Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, " +
		"and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from " +
		"its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its " +
		"conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.\n")
	for u, unit := range r.Units {
		level := "##"
		if len(r.Units) > 1 {
			fmt.Fprintf(&b, "\n## %s\n", unit.Program)
			level = "###"
		}
		for t, table := range unit.Tables {
			title := unit.blueprint.Rules.Tables[t].Title()
			fmt.Fprintf(&b, "\n%s [%s](%s/%s.md#%s)\n\n| Rule | Reached by |\n| --: | --- |\n",
				level, title, rulesDir, unit.Program, blueprint.HeadingAnchor(title))
			for i, rl := range table.Rules {
				cell := "not reached: " + outcomes(table.Conditions, rl.Entries)
				if c, ok := r.First[rule{u, t, i}]; ok {
					cell = fmt.Sprintf("[case %s](%s/)", caseName(c), caseName(c))
				}
				fmt.Fprintf(&b, "| %d | %s |\n", i+1, cell)
			}
		}
	}
	return b.String()
}

// outcomes writes the outcome of each condition a rule tests.
func outcomes(conditions, entries []string) string {
	var parts []string
	for i, e := range entries {
		if e != "" {
			parts = append(parts, fmt.Sprintf("`%s` %s", strings.ReplaceAll(conditions[i], "|", `\|`), e))
		}
	}
	return strings.Join(parts, ", ")
}

func (r *Result) caseReadme(i int, rulesDir string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s: case %s\n\n", r.Program, caseName(i))
	run := r.runs[i]
	switch {
	case run.outcome.TimedOut:
		b.WriteString("The run did not end in time, and was killed.\n")
	case run.outcome.Abended:
		b.WriteString("The run abended.\n")
	default:
		fmt.Fprintf(&b, "The run ended with return code %d.\n", run.outcome.RC)
	}
	if file, how, ok := r.Cases[i].failedFile(); ok {
		fmt.Fprintf(&b, "\n%s\n", how.describe(file))
	}
	var first []string
	for u, unit := range r.Units {
		for t, table := range unit.Tables {
			title := unit.blueprint.Rules.Tables[t].Title()
			for j := range table.Rules {
				if r.First[rule{u, t, j}] == i && run.reached[rule{u, t, j}] {
					first = append(first, fmt.Sprintf("- [%s, %s, rule %d](%s/%s.md#%s)", unit.Program, title, j+1,
						rulesDir, unit.Program, blueprint.HeadingAnchor(title)))
				}
			}
		}
	}
	if len(first) > 0 {
		fmt.Fprintf(&b, "\nThe first case to reach:\n\n%s\n", strings.Join(first, "\n"))
	}
	return b.String()
}

// encode writes an input's records field by field, a field a line, NAME
// VALUE, with a blank line between records: a number as a decimal, text
// quoted. A record whose fields do not give back its bytes is written
// whole, quoted, as RECORD.
func (p *program) encode(in input, records [][]byte) []byte {
	var blocks []string
	for _, r := range records {
		block := p.encodeFields(in, r)
		if back, err := p.decodeRecord(in, block); err != nil || !bytes.Equal(back, r) {
			block = fmt.Sprintf("RECORD %s\n", strconv.Quote(string(r)))
		}
		blocks = append(blocks, block)
	}
	return []byte(strings.Join(blocks, "\n"))
}

func (p *program) encodeFields(in input, r []byte) string {
	var b strings.Builder
	s := p.shapes[in.name()]
	walk(in.root, func(f *layout.Field) {
		if len(f.Children) > 0 || f.Name == "FILLER" || f.Occurs > 0 || !s.has(in, r, f) {
			return
		}
		start := f.Start - in.root.Start
		raw := r[start : start+f.Size()]
		value := strconv.Quote(strings.TrimRight(string(raw), " "))
		if numeric(f) {
			if n, err := in.layout.Number(r, f.Name); err == nil {
				value = decimal(n, f.Scale)
			}
		}
		fmt.Fprintf(&b, "%-24s %s\n", f.Name, value)
	})
	return b.String()
}

// decimal writes a number in units of a scale as a decimal.
func decimal(n int64, scale int) string {
	sign := ""
	if n < 0 {
		sign, n = "-", -n
	}
	digits := strconv.FormatInt(n, 10)
	if scale == 0 {
		return sign + digits
	}
	digits = strings.Repeat("0", max(0, scale+1-len(digits))) + digits
	return sign + digits[:len(digits)-scale] + "." + digits[len(digits)-scale:]
}

// decodeRecord reads a record back from its fields: a blank record, with
// each field set to its value, or the record written whole.
func (p *program) decodeRecord(in input, block string) ([]byte, error) {
	r := in.layout.Blank()
	for _, line := range strings.Split(strings.TrimSpace(block), "\n") {
		name, value, _ := strings.Cut(line, " ")
		value = strings.TrimSpace(value)
		if name == "RECORD" {
			whole, err := strconv.Unquote(value)
			return []byte(whole), err
		}
		f := find(in.root, name)
		if f == nil {
			return nil, fmt.Errorf("%s: no field %s", in.name(), name)
		}
		v, err := parseValue(f, value)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %w", in.name(), name, err)
		}
		if !set(r, in.root, f, in.layout, v) {
			return nil, fmt.Errorf("%s %s: %s does not fit", in.name(), name, value)
		}
	}
	return r, nil
}

func parseValue(f *layout.Field, value string) (value, error) {
	if strings.HasPrefix(value, `"`) {
		s, err := strconv.Unquote(value)
		return text(s), err
	}
	v, ok := scaled(f, value)
	if !ok {
		return v, fmt.Errorf("%q is no number", value)
	}
	return v, nil
}

// readCases reads the cases recorded in a program's directory, each from
// its input files.
func (p *program) readCases(dir string) ([]Case, error) {
	var cases []Case
	for i := 0; ; i++ {
		caseDir := filepath.Join(dir, caseName(i))
		if _, err := os.Stat(caseDir); os.IsNotExist(err) {
			return cases, nil
		}
		c := Case{}
		if data, err := os.ReadFile(filepath.Join(caseDir, "input", failedName)); err == nil {
			c[failed] = [][]byte{bytes.TrimSuffix(data, []byte("\n"))}
		} else if !os.IsNotExist(err) {
			return nil, err
		}
		for _, in := range p.inputs {
			if _, fails := c.failsFile(in.name()); fails {
				continue
			}
			data, err := os.ReadFile(filepath.Join(caseDir, "input", in.name()+".txt"))
			if err != nil {
				return nil, err
			}
			for _, block := range strings.Split(string(data), "\n\n") {
				if strings.TrimSpace(block) == "" {
					continue
				}
				r, err := p.decodeRecord(in, block)
				if err != nil {
					return nil, err
				}
				c[in.name()] = append(c[in.name()], r)
			}
		}
		cases = append(cases, c)
	}
}

// Index writes the index of a system's cases: for each main program, the
// rules its cases reach, its own and those of the programs it calls, and
// how many cases there are; and the batches, each job run on inputs from
// its programs' cases.
func Index(title string, results []*Result, batches []*Batch) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s: cases\n\n", title)
	b.WriteString("For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, " +
		"with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them " +
		"to check them.\n\n| Program | Rules reached | Cases |\n| --- | --: | --: |\n")
	for _, r := range results {
		var reached, total int
		for u := range r.Units {
			n, m := r.Rules(u)
			reached, total = reached+n, total+m
		}
		fmt.Fprintf(&b, "| [%s](%s/) | %d of %d | %d |\n", r.Program, r.Program, reached, total, len(r.Cases))
	}
	if len(batches) > 0 {
		b.WriteString("\nAnd for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, " +
			"with the datasets it wrote and its output.\n\n| Job | Steps |\n| --- | --- |\n")
		for _, batch := range batches {
			var steps []string
			for _, step := range batch.steps {
				steps = append(steps, step.Step+" "+step.Outcome())
			}
			fmt.Fprintf(&b, "| [%s](jobs/%s/) | %s |\n", batch.Job, batch.Job, strings.Join(steps, ", "))
		}
	}
	return []byte(b.String())
}
