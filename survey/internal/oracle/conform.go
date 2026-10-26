package oracle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/copal/bedrock/run"
	"github.com/vlumi/copal/survey/blueprint"
)

// Conformance is how a replacement for a program did on the program's
// recorded cases: for each case, how what it returned, wrote and
// displayed differs from what was recorded, if it does.
type Conformance struct {
	Program string
	With    string
	Cases   []CaseConformance
}

// CaseConformance is one case's differences; none means the replacement
// conforms on it.
type CaseConformance struct {
	Name        string
	Differences []string
}

// Conforms reports whether the replacement conforms on every case.
func (c *Conformance) Conforms() bool {
	return !slices.ContainsFunc(c.Cases, func(cc CaseConformance) bool { return len(cc.Differences) > 0 })
}

// Conform runs a replacement for a program on the cases recorded in dir,
// the way the oracle ran the program, and compares what it returns,
// writes and displays with what was recorded. The replacement is an
// executable: it gets the case's command line as its argument, its SYSIN
// on its standard input and each file's path as DD_ASSIGN in its
// environment, and its exit status is its return code. A file comes as
// its records, back to back or a line each as the program's SELECT has
// it, an indexed one as a sequential file of its records in key order,
// and is read back the same way. What the runtime itself displayed in the
// recorded run, its warnings, is not compared.
func Conform(system *blueprint.System, bp *blueprint.Program, copybookDir, dir, scratch, with string) (*Conformance, error) {
	p, err := prepare(system, bp, copybookDir)
	if err != nil {
		return nil, err
	}
	cases, err := p.readCases(dir)
	if err != nil {
		return nil, err
	}
	recorded, err := readCoverage(filepath.Join(dir, "coverage.json"))
	if err != nil {
		return nil, err
	}
	if len(recorded.Cases) != len(cases) {
		return nil, fmt.Errorf("%s: %d cases, %d in coverage.json", dir, len(cases), len(recorded.Cases))
	}
	c := &Conformance{Program: bp.Name, With: filepath.Base(with)}
	for i, cs := range cases {
		caseDir := filepath.Join(scratch, caseName(i))
		if err := os.MkdirAll(caseDir, 0o755); err != nil {
			return nil, err
		}
		got, err := p.runReplacement(with, cs, caseDir)
		if err != nil {
			return nil, fmt.Errorf("case %s: %w", caseName(i), err)
		}
		outputs, err := readOutputs(filepath.Join(dir, caseName(i), "output"))
		if err != nil {
			return nil, err
		}
		c.Cases = append(c.Cases, CaseConformance{Name: caseName(i), Differences: p.differences(recorded.Cases[i], outputs, got)})
	}
	return c, nil
}

func readCoverage(path string) (*Coverage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Coverage
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &c, nil
}

// readOutputs reads what a recorded case's run wrote, by file name, with
// STDOUT for what it displayed.
func readOutputs(dir string) (map[string][]byte, error) {
	entries, err := os.ReadDir(dir)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	outputs := map[string][]byte{}
	for _, e := range entries {
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		outputs[strings.TrimSuffix(e.Name(), ".txt")] = data
	}
	return outputs, nil
}

// runReplacement runs a replacement on a case in a directory of its own,
// as the oracle runs the program, except that an indexed file is given
// and read back as a sequential file of its records in key order.
func (p *program) runReplacement(with string, c Case, dir string) (result, error) {
	dds := map[string]string{}
	for _, f := range p.blueprint.Files {
		path := filepath.Join(dir, f.Name)
		dds[f.Assign] = path
		var err error
		how, fails := c.failsFile(f.Name)
		switch {
		case fails:
			dds[f.Assign], err = prepareFailure(how, path)
		case reads(f):
			err = os.WriteFile(path, sequential(f, p.inKeyOrder(f, c[f.Name])), 0o644)
		case writes(f) && !keyed(f):
			err = os.WriteFile(path, nil, 0o644)
		}
		if err != nil {
			return result{}, err
		}
	}
	var stdout bytes.Buffer
	var args []string
	if line := c[commandLine]; len(line) > 0 {
		args = []string{strings.TrimRight(string(line[0]), " ")}
	}
	outcome, err := run.Step{Programs: filepath.Dir(with), Program: filepath.Base(with), Args: args, DDs: dds, Stdout: &stdout,
		Stdin: bytes.NewReader(lines(c[sysin])), Timeout: timeout}.Run()
	if err != nil {
		return result{}, err
	}
	res := result{outcome: outcome, stdout: stdout.Bytes(), outputs: map[string][]byte{}}
	for _, f := range p.blueprint.Files {
		if _, fails := c.failsFile(f.Name); fails || !writes(f) && !slices.Contains(f.Opened, "I-O") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, f.Name))
		if err != nil && !os.IsNotExist(err) {
			return result{}, err
		}
		res.outputs[f.Name] = asLines(data, f.Organization == "LINE SEQUENTIAL", recordLength(p.records, f))
	}
	return res, nil
}

// inKeyOrder puts the records of an indexed file in the order the file
// holds them, by key, the first of a key kept as the loader keeps it;
// another file's records stay as the case has them.
func (p *program) inKeyOrder(f blueprint.File, records [][]byte) [][]byte {
	root := fileRecord(p.records, f.Name)
	if f.Organization != "INDEXED" || root == nil {
		return records
	}
	key := find(root, f.RecordKey)
	if key == nil {
		return records
	}
	keyOf := func(r []byte) []byte { return r[key.Start-root.Start : key.Start-root.Start+key.Size()] }
	sorted := slices.Clone(records)
	slices.SortStableFunc(sorted, func(a, b []byte) int { return bytes.Compare(keyOf(a), keyOf(b)) })
	return slices.CompactFunc(sorted, func(a, b []byte) bool { return bytes.Equal(keyOf(a), keyOf(b)) })
}

// differences says how a replacement's run of a case differs from the
// recorded one: in how it ended, in each file the program writes, record
// by record and field by field, and in what it displayed.
func (p *program) differences(recorded CaseOutcome, outputs map[string][]byte, got result) []string {
	var diffs []string
	switch {
	case got.outcome.TimedOut:
		diffs = append(diffs, "did not end in time")
	case got.outcome.Abended && !recorded.Abended:
		diffs = append(diffs, fmt.Sprintf("abended, recorded return code %d", recorded.ReturnCode))
	case !got.outcome.Abended && recorded.Abended:
		diffs = append(diffs, fmt.Sprintf("return code %d, recorded an abend", got.outcome.RC))
	case got.outcome.RC != recorded.ReturnCode:
		diffs = append(diffs, fmt.Sprintf("return code %d, recorded %d", got.outcome.RC, recorded.ReturnCode))
	}
	for _, f := range p.blueprint.Files {
		want, recorded := outputs[f.Name]
		data, written := got.outputs[f.Name]
		switch {
		case recorded && !written:
			diffs = append(diffs, f.Name+" not written")
		case written && !recorded:
			diffs = append(diffs, f.Name+" written, not recorded")
		case written:
			if d := p.recordsDiffer(f, data, want); d != "" {
				diffs = append(diffs, d)
			}
		}
	}
	if d := linesDiffer("displayed", programLines(asLines(got.stdout, true, 0)), programLines(outputs["STDOUT"])); d != "" {
		diffs = append(diffs, d)
	}
	return diffs
}

// programLines is what a run displayed less what the runtime itself said,
// its warnings of an implicit CLOSE, which are no replacement's to match.
func programLines(stdout []byte) [][]byte {
	var kept [][]byte
	for _, line := range splitLines(stdout) {
		if !bytes.HasPrefix(line, []byte("libcob: ")) {
			kept = append(kept, line)
		}
	}
	return kept
}

func splitLines(data []byte) [][]byte {
	if len(data) == 0 {
		return nil
	}
	return bytes.Split(bytes.TrimSuffix(data, []byte("\n")), []byte("\n"))
}

// recordsDiffer says how a file's records differ from the recorded ones:
// in their count, or at the first record that differs, by the fields of
// the file's record that do.
func (p *program) recordsDiffer(f blueprint.File, got, want []byte) string {
	g, w := splitLines(got), splitLines(want)
	if len(g) != len(w) {
		return fmt.Sprintf("%s: %s, %s recorded", f.Name, records(len(g)), records(len(w)))
	}
	for i := range g {
		if bytes.Equal(g[i], w[i]) {
			continue
		}
		what := fmt.Sprintf("%s, record %d: ", f.Name, i+1)
		if fields := p.fieldsDiffer(f, g[i], w[i]); fields != "" {
			return what + fields
		}
		return what + fmt.Sprintf("%s, recorded %s", strconv.Quote(string(g[i])), strconv.Quote(string(w[i])))
	}
	return ""
}

func records(n int) string {
	switch n {
	case 0:
		return "no records"
	case 1:
		return "1 record"
	}
	return strconv.Itoa(n) + " records"
}

// fieldsDiffer names the elementary fields of a file's record whose bytes
// differ between a record written and the one recorded, with both values,
// or nothing when the file's record is not laid out or the difference
// falls in no named field.
func (p *program) fieldsDiffer(f blueprint.File, got, want []byte) string {
	root := fileRecord(p.records, f.Name)
	if root == nil {
		return ""
	}
	got, want = padded(got, root.Size()), padded(want, root.Size())
	var fields []string
	walk(root, func(field *layout.Field) {
		if len(field.Children) > 0 || field.Name == "FILLER" || field.Occurs > 0 {
			return
		}
		start, end := field.Start-root.Start, field.Start-root.Start+field.Size()
		if g, w := got[start:end], want[start:end]; !bytes.Equal(g, w) {
			fields = append(fields, fmt.Sprintf("%s %s, recorded %s", field.Name, trimmed(g), trimmed(w)))
		}
	})
	return strings.Join(fields, "; ")
}

// padded is a recorded line as a record: unquoted if it was quoted, and
// padded with spaces to the record's length.
func padded(line []byte, length int) []byte {
	if bytes.HasPrefix(line, []byte(`"`)) {
		if s, err := strconv.Unquote(string(line)); err == nil {
			line = []byte(s)
		}
	}
	if len(line) < length {
		line = append(slices.Clone(line), bytes.Repeat([]byte(" "), length-len(line))...)
	}
	return line
}

func trimmed(b []byte) string {
	return strconv.Quote(strings.TrimRight(string(b), " "))
}

// linesDiffer says where two sets of lines first differ, or that one has
// lines the other has not.
func linesDiffer(what string, got, want [][]byte) string {
	for i := range max(len(got), len(want)) {
		switch {
		case i >= len(got):
			return fmt.Sprintf("%s nothing, recorded %s", what, strconv.Quote(string(want[i])))
		case i >= len(want):
			return fmt.Sprintf("%s %s, recorded nothing", what, strconv.Quote(string(got[i])))
		case !bytes.Equal(got[i], want[i]):
			return fmt.Sprintf("%s %s, recorded %s", what, strconv.Quote(string(got[i])), strconv.Quote(string(want[i])))
		}
	}
	return ""
}

// Markdown writes the conformance as a report: a line for the whole, and a
// table of the cases with their differences.
func (c *Conformance) Markdown() []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s against %s\n\n", c.Program, c.With)
	conforming := 0
	for _, cc := range c.Cases {
		if len(cc.Differences) == 0 {
			conforming++
		}
	}
	switch {
	case conforming == len(c.Cases):
		fmt.Fprintf(&b, "%s conforms on all %d cases: it returns, writes and displays what %s did.\n", c.With, len(c.Cases), c.Program)
	default:
		fmt.Fprintf(&b, "%s conforms on %d of %d cases.\n", c.With, conforming, len(c.Cases))
	}
	b.WriteString("\n| Case | Differences |\n| ---: | --- |\n")
	for _, cc := range c.Cases {
		differences := "none"
		if len(cc.Differences) > 0 {
			differences = strings.Join(cc.Differences, "; ")
		}
		fmt.Fprintf(&b, "| %s | %s |\n", cc.Name, strings.ReplaceAll(differences, "|", "\\|"))
	}
	return []byte(b.String())
}
