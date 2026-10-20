// Package script reads batch scripts in the small subset of POSIX sh that
// Laurasia's are written in: variables, exported DD_ bindings, programs run
// with their output redirected, return codes tested with if, Unix sort and
// mv, and one script running another. The survey reads them, as it reads
// JCL, to draw the job network.
package script

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Script is a parsed script: its name, the comment at its head that
// describes it, its steps in order, and the variables that only stand for
// a parameter of the script or a variable of its environment, as DATE for
// $1 or DATA for $LCS_DATA, which its steps' paths are written in.
type Script struct {
	Name       string
	Comment    []string
	Steps      []Step
	Parameters map[string]string
}

// Kind is what a step runs.
type Kind string

// The kinds of step.
const (
	Program Kind = "program" // a COBOL program
	Sort    Kind = "sort"    // Unix sort, of Input into Output
	Move    Kind = "move"    // mv, of Input to Output
	Run     Kind = "script"  // another script
)

// Step is one step of a script, named by the # STEP NAME: comment before
// it, whose text after the name, and the comment lines after it, are its
// Comment. Command is what it runs, as written with the script's own
// variables expanded, and Name of a program or script its file's name
// without the directory or .sh. A program runs with the DD_ bindings
// exported before it, by name without DD_, and its standard output goes
// to Stdout. A sort and a move read Input and write Output; a sort orders
// by its Keys. StopsIf is the test of the step's return code after which
// the script exits, if it tests it.
type Step struct {
	Name     string
	Comment  []string
	Kind     Kind
	Command  string
	Program  string
	Args     []string
	Bindings map[string]string
	Stdout   string
	Input    []string
	Output   string
	Keys     []Key
	StopsIf  *Test
	Line     int
}

// Key is a sort key: the columns of each record, 1-based, from Start
// through End.
type Key struct {
	Start, End int
}

// Test is a test of a step's return code, RC Op Value, with Op one of
// the test command's -eq, -ne, -lt, -le, -gt, -ge.
type Test struct {
	Op    string
	Value int
}

// Holds reports whether the test holds for a return code.
func (t Test) Holds(rc int) bool {
	switch t.Op {
	case "-eq":
		return rc == t.Value
	case "-ne":
		return rc != t.Value
	case "-lt":
		return rc < t.Value
	case "-le":
		return rc <= t.Value
	case "-gt":
		return rc > t.Value
	}
	return rc >= t.Value
}

// ParseFile reads a script from a file, named after it.
func ParseFile(path string) (*Script, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	s, err := Parse(strings.TrimSuffix(filepath.Base(path), ".sh"), f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return s, nil
}

// Parse reads a script from r.
func Parse(name string, r io.Reader) (*Script, error) {
	p := &parser{script: &Script{Name: name, Parameters: map[string]string{}}, vars: map[string]string{}, bindings: map[string]string{}, head: true}
	scanner := bufio.NewScanner(r)
	for line := 1; scanner.Scan(); line++ {
		if err := p.line(line, scanner.Text()); err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
	}
	return p.script, scanner.Err()
}

type parser struct {
	script   *Script
	vars     map[string]string
	bindings map[string]string
	head     bool  // still in the comment at the script's head
	pending  *Step // a step's name and comment, until its command
	last     *Step // the step whose return code the next test tests
}

var stepComment = regexp.MustCompile(`^STEP ([A-Z0-9-]+):?\s*(.*)$`)

func (p *parser) line(n int, text string) error {
	trimmed := strings.TrimSpace(text)
	switch {
	case strings.HasPrefix(trimmed, "#!"):
		return nil
	case trimmed == "":
		p.head = false
		return nil
	case strings.HasPrefix(trimmed, "#"):
		p.comment(strings.TrimSpace(strings.TrimPrefix(trimmed, "#")))
		return nil
	}
	p.head = false
	words, err := p.words(trimmed)
	if err != nil {
		return err
	}
	return p.command(n, words)
}

// comment adds a comment line to the script's head, to the step it names,
// or to the step named before it.
func (p *parser) comment(text string) {
	switch m := stepComment.FindStringSubmatch(text); {
	case p.head:
		p.script.Comment = append(p.script.Comment, text)
	case m != nil:
		p.pending = &Step{Name: m[1], Comment: []string{m[2]}}
	case p.pending != nil:
		p.pending.Comment = append(p.pending.Comment, text)
	}
}

func (p *parser) command(n int, words []string) error {
	if name, value, ok := assignment(words[0]); ok && len(words) == 1 {
		p.assign(name, value)
		return nil
	}
	switch words[0] {
	case "export":
		for _, w := range words[1:] {
			if name, value, ok := assignment(w); ok {
				p.export(name, value)
			}
		}
		return nil
	case "mkdir", "echo":
		return nil
	case "if":
		return p.test(words)
	case "sort":
		return p.step(n, Sort, words, p.sort)
	case "mv":
		return p.step(n, Move, words, p.move)
	}
	kind := Program
	if strings.HasSuffix(words[0], ".sh") {
		kind = Run
	}
	return p.step(n, kind, words, p.run)
}

// assignment splits NAME=value.
func assignment(word string) (name, value string, ok bool) {
	name, value, ok = strings.Cut(word, "=")
	return name, value, ok && variableName.MatchString(name)
}

var (
	variableName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	parameter    = regexp.MustCompile(`^\$[A-Za-z0-9_?]+$`)
	rcTest       = regexp.MustCompile(`^if \[ \$rc (-eq|-ne|-lt|-le|-gt|-ge) (\d+) \] ; then exit \S+ ; fi$`)
)

// assign sets a variable. One that is only another parameter, such as
// DATE="$1" or DATA="$LCS_DATA", stays a name of its own, so that paths
// read $DATA/work/$DATE and not $LCS_DATA/work/$1.
func (p *parser) assign(name, value string) {
	if parameter.MatchString(value) {
		if value != "$?" {
			p.script.Parameters[name] = strings.TrimPrefix(value, "$")
		}
		value = "$" + name
	}
	p.vars[name] = value
}

func (p *parser) export(name, value string) {
	if binding, ok := strings.CutPrefix(name, "DD_"); ok {
		p.bindings[binding] = value
		return
	}
	p.assign(name, value)
}

// step reads a step's command into a step of the kind, with the name and
// comment before it.
func (p *parser) step(n int, kind Kind, words []string, read func(*Step, []string) error) error {
	s := Step{Kind: kind, Line: n}
	if p.pending != nil {
		s.Name, s.Comment = p.pending.Name, p.pending.Comment
		p.pending = nil
	}
	words, s.Stdout = redirection(words)
	if err := read(&s, words); err != nil {
		return err
	}
	if s.Name == "" {
		s.Name = strings.ToUpper(s.Program)
	}
	p.script.Steps = append(p.script.Steps, s)
	p.last = &p.script.Steps[len(p.script.Steps)-1]
	return nil
}

// redirection takes the redirection of standard output off a command's
// words, and the redirection of standard error into it.
func redirection(words []string) (rest []string, stdout string) {
	for i := 0; i < len(words); i++ {
		switch words[i] {
		case ">", ">>":
			if i+1 < len(words) {
				stdout = words[i+1]
				i++
			}
		case "2>&1":
		default:
			rest = append(rest, words[i])
		}
	}
	return rest, stdout
}

func (p *parser) run(s *Step, words []string) error {
	s.Command, s.Args = words[0], words[1:]
	s.Program = strings.TrimSuffix(filepath.Base(words[0]), ".sh")
	if s.Kind == Program {
		s.Bindings = map[string]string{}
		for name, path := range p.bindings {
			s.Bindings[name] = path
		}
	}
	return nil
}

var sortKey = regexp.MustCompile(`^(\d+)\.(\d+),(\d+)\.(\d+)$`)

// sort reads sort [-s] [-t c] [-k F.C,F.C]... input, with its output
// redirected. Keys are on the first field of a line, which a separator
// that never occurs in a record leaves whole.
func (p *parser) sort(s *Step, words []string) error {
	s.Command, s.Program = "sort", "sort"
	for i := 1; i < len(words); i++ {
		switch w := words[i]; {
		case w == "-s":
		case w == "-t":
			i++
		case strings.HasPrefix(w, "-k"):
			spec := strings.TrimPrefix(w, "-k")
			if spec == "" && i+1 < len(words) {
				i++
				spec = words[i]
			}
			m := sortKey.FindStringSubmatch(spec)
			if m == nil || m[1] != "1" || m[3] != "1" {
				return fmt.Errorf("sort key %s is not on the first field, by columns", spec)
			}
			start, _ := strconv.Atoi(m[2])
			end, _ := strconv.Atoi(m[4])
			s.Keys = append(s.Keys, Key{Start: start, End: end})
		case strings.HasPrefix(w, "-"):
			return fmt.Errorf("sort option %s is not supported", w)
		default:
			s.Input = append(s.Input, w)
		}
	}
	s.Output, s.Stdout = s.Stdout, ""
	if len(s.Input) == 0 || s.Output == "" {
		return fmt.Errorf("sort without an input file or an output redirected")
	}
	return nil
}

// move reads mv source target.
func (p *parser) move(s *Step, words []string) error {
	s.Command, s.Program = "mv", "mv"
	if len(words) != 3 {
		return fmt.Errorf("mv takes a source and a target")
	}
	s.Input, s.Output = []string{words[1]}, words[2]
	return nil
}

// test reads if [ "$rc" OP N ]; then exit ...; fi, after a step, which
// exits the script when the step's return code passes the test.
func (p *parser) test(words []string) error {
	text := strings.Join(words, " ")
	m := rcTest.FindStringSubmatch(text)
	if m == nil {
		return fmt.Errorf("%q is not a test of a step's return code", text)
	}
	if p.last == nil {
		return fmt.Errorf("a return code tested before any step")
	}
	value, _ := strconv.Atoi(m[2])
	p.last.StopsIf = &Test{Op: m[1], Value: value}
	return nil
}
