// Package batchfile reads Windows batch files in the small subset that
// Gondwana's are written in, and runs them, standing in for cmd.exe as
// bedrock's runner stands in for JES: variables set, the DD_ bindings
// among them, programs run with their output redirected, one batch file
// calling another, echo into a log, move, and return codes tested with if
// errorlevel. A batch file also reads into the model a shell script reads
// into, so that the survey draws its job network from either.
package batchfile

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

// File is a parsed batch file: its name, the rem lines at its head that
// describe it, and its lines in order.
type File struct {
	Name    string
	Comment []string
	Lines   []Line
}

// Kind is what a line does.
type Kind string

// The kinds of line, with the fields of Line each uses.
const (
	Set   Kind = "set"   // set NAME=value: Var and Value
	Run   Kind = "run"   // a program run: Words, its path and arguments, and Output
	Call  Kind = "call"  // call BATCH args: Words, its path and arguments
	Move  Kind = "move"  // move /y SOURCE TARGET: Words
	Echo  Kind = "echo"  // >>FILE echo text: Text and Output
	Mkdir Kind = "mkdir" // if not exist DIR mkdir DIR: Words, the directory
	Test  Kind = "test"  // if errorlevel N exit /b CODE: Level and Code
	Exit  Kind = "exit"  // exit /b CODE: Code
)

// Line is one line of a batch file, as written, with its variables not
// expanded. A program run, a call or a move is a step, named by the rem
// STEP NAME: comment before it, whose text and the rem lines after it are
// its Comment. Code is as written, a number or %ERRORLEVEL%.
type Line struct {
	N       int
	Kind    Kind
	Step    string
	Comment []string
	Var     string
	Value   string
	Words   []string
	Text    string
	Output  string
	Level   int
	Code    string
}

// ParseFile reads a batch file, named after it.
func ParseFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	file, err := Parse(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return file, nil
}

// Parse reads a batch file from r.
func Parse(name string, r io.Reader) (*File, error) {
	p := &parser{file: &File{Name: name}, head: true}
	scanner := bufio.NewScanner(r)
	for n := 1; scanner.Scan(); n++ {
		if err := p.line(n, strings.TrimRight(scanner.Text(), "\r")); err != nil {
			return nil, fmt.Errorf("line %d: %w", n, err)
		}
	}
	return p.file, scanner.Err()
}

type parser struct {
	file    *File
	head    bool  // still in the comment at the file's head
	pending *Line // a step's name and comment, until its line
}

var (
	stepComment = regexp.MustCompile(`^STEP ([A-Z0-9-]+):?\s*(.*)$`)
	assignment  = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)=(.*)$`)
)

func (p *parser) line(n int, text string) error {
	trimmed := strings.TrimPrefix(strings.TrimSpace(text), "@")
	if trimmed == "" {
		p.head = false
		return nil
	}
	words, err := split(trimmed)
	if err != nil {
		return err
	}
	keyword := strings.ToLower(words[0])
	if keyword == "rem" {
		p.comment(strings.TrimSpace(strings.TrimPrefix(trimmed[3:], " ")))
		return nil
	}
	line := Line{N: n}
	switch {
	case keyword == "echo" && len(words) == 2 && strings.EqualFold(words[1], "off"),
		keyword == "setlocal", keyword == "endlocal":
		return nil
	}
	p.head = false
	switch {
	case keyword == "set" && len(words) == 2:
		m := assignment.FindStringSubmatch(words[1])
		if m == nil {
			return fmt.Errorf("%q is not set NAME=value", trimmed)
		}
		line.Kind, line.Var, line.Value = Set, strings.ToUpper(m[1]), m[2]
	case keyword == "if":
		return p.conditional(line, words)
	case keyword == "exit":
		line.Kind, line.Code = Exit, exitCode(words[1:])
	case keyword == "call" && len(words) >= 2:
		line.Kind, line.Words = Call, words[1:]
	case keyword == "move":
		if len(words) > 1 && strings.EqualFold(words[1], "/y") {
			words = append(words[:1], words[2:]...)
		}
		if len(words) != 3 {
			return fmt.Errorf("move takes a source and a target")
		}
		line.Kind, line.Words = Move, words[1:]
	case words[0] == ">>" && len(words) >= 3 && strings.EqualFold(words[2], "echo"):
		line.Kind, line.Output, line.Text = Echo, words[1], strings.Join(words[3:], " ")
	case keyword == "echo":
		return fmt.Errorf("echo writes to a file: >>FILE echo text")
	default:
		line.Kind = Run
		line.Words, line.Output = redirection(words)
	}
	p.add(line)
	return nil
}

// conditional reads the two lines that start with if: a directory made
// if it is not there, and a step's return code tested.
func (p *parser) conditional(line Line, words []string) error {
	lower := make([]string, len(words))
	for i, w := range words {
		lower[i] = strings.ToLower(w)
	}
	switch {
	case len(words) == 6 && lower[1] == "not" && lower[2] == "exist" && lower[4] == "mkdir" && words[3] == words[5]:
		line.Kind, line.Words = Mkdir, words[3:4]
	case len(words) >= 5 && lower[1] == "errorlevel" && lower[3] == "exit":
		level, err := strconv.Atoi(words[2])
		if err != nil {
			return fmt.Errorf("if errorlevel %s: not a number", words[2])
		}
		line.Kind, line.Level, line.Code = Test, level, exitCode(words[4:])
	default:
		return fmt.Errorf("%q is neither a directory made if missing nor a test of errorlevel", strings.Join(words, " "))
	}
	p.add(line)
	return nil
}

// exitCode is the code after exit /b, or %ERRORLEVEL% when there is none,
// which is what cmd.exe exits with then.
func exitCode(words []string) string {
	if len(words) > 0 && strings.EqualFold(words[0], "/b") {
		words = words[1:]
	}
	if len(words) == 0 {
		return "%ERRORLEVEL%"
	}
	return words[0]
}

// comment adds a rem line to the file's head, to the step it names, or
// to the step named before it.
func (p *parser) comment(text string) {
	switch m := stepComment.FindStringSubmatch(text); {
	case p.head:
		p.file.Comment = append(p.file.Comment, text)
	case m != nil:
		p.pending = &Line{Step: m[1], Comment: []string{m[2]}}
	case p.pending != nil:
		p.pending.Comment = append(p.pending.Comment, text)
	}
}

// add adds a line, a step with the name and comment pending for it.
func (p *parser) add(line Line) {
	if p.pending != nil && (line.Kind == Run || line.Kind == Call || line.Kind == Move) {
		line.Step, line.Comment = p.pending.Step, p.pending.Comment
		p.pending = nil
	}
	p.file.Lines = append(p.file.Lines, line)
}

// redirection takes the redirection of standard output off a command's
// words, and the redirection of standard error into it.
func redirection(words []string) (rest []string, output string) {
	for i := 0; i < len(words); i++ {
		switch words[i] {
		case ">", ">>":
			if i+1 < len(words) {
				output = words[i+1]
				i++
			}
		case "2>&1":
		default:
			rest = append(rest, words[i])
		}
	}
	return rest, output
}

// split splits a line into words: on blanks outside double quotes, with
// the quotes removed, and a redirection glued to the word after it, as in
// >>"%LOG%", a word of its own.
func split(line string) ([]string, error) {
	var words []string
	var word strings.Builder
	inWord := false
	end := func() {
		if inWord {
			words = append(words, word.String())
			word.Reset()
			inWord = false
		}
	}
	for i := 0; i < len(line); i++ {
		switch c := line[i]; {
		case c == ' ' || c == '\t':
			end()
		case c == '"':
			j := strings.IndexByte(line[i+1:], '"')
			if j < 0 {
				return nil, fmt.Errorf("a quote not closed")
			}
			word.WriteString(line[i+1 : i+1+j])
			i, inWord = i+1+j, true
		case c == '>' && !inWord:
			end()
			if i+1 < len(line) && line[i+1] == '>' {
				i++
				words = append(words, ">>")
			} else {
				words = append(words, ">")
			}
		default:
			word.WriteByte(c)
			inWord = true
		}
	}
	end()
	if len(words) == 0 {
		return nil, fmt.Errorf("an empty line")
	}
	return words, nil
}
