package blueprint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Annotation is prose a model wrote about a subject of the blueprints, a
// paragraph or section of a program: a name that says what it is for, and
// a sentence or two that say it. It is no fact of the survey's. Hash is
// the hash of the facts it was written from, as Facts gives them, and
// Writer and Prompt the model and the version of the prompt that wrote it.
type Annotation struct {
	Kind    string `json:"kind"`
	Program string `json:"program"`
	Subject string `json:"subject"`
	Hash    string `json:"hash"`
	Name    string `json:"name,omitempty"`
	Text    string `json:"text"`
	Writer  string `json:"writer"`
	Prompt  string `json:"prompt"`
}

// KindParagraph is an annotation of a paragraph, or of the statements of a
// section before its first paragraph.
const KindParagraph = "paragraph"

// Facts are what the survey established about a subject, as the text a
// model is given to write about it, and the hash of that text.
type Facts struct {
	Kind, Program, Subject string
	Text                   string
	Hash                   string
}

// ParagraphFacts are the facts of a paragraph of a program: where it is,
// its source, what it performs and what performs it, the items it reads
// and writes, its decision tables, and the program's own description.
func (s *System) ParagraphFacts(program, paragraph string) (Facts, error) {
	i := slices.IndexFunc(s.Programs, func(p Program) bool { return p.Name == program })
	if i < 0 {
		return Facts{}, fmt.Errorf("no program %s", program)
	}
	p := &s.Programs[i]
	j := slices.IndexFunc(p.Paragraphs, func(para Paragraph) bool { return para.Name == paragraph })
	if j < 0 {
		return Facts{}, fmt.Errorf("%s: no paragraph %s", program, paragraph)
	}
	para := p.Paragraphs[j]
	source, err := s.paragraphSource(para)
	if err != nil {
		return Facts{}, err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Program %s, paragraph %s", p.Name, para.Name)
	if para.Section != "" && para.Section != para.Name {
		fmt.Fprintf(&b, ", in section %s", para.Section)
	}
	b.WriteString(".\n")
	if len(p.Description) > 0 {
		fmt.Fprintf(&b, "\nThe program's description:\n%s\n", strings.Join(p.Description, "\n"))
	}
	fmt.Fprintf(&b, "\nIts source, from %s:\n%s\n", para.Source, source)
	list := func(label string, names []string) {
		if len(names) > 0 {
			fmt.Fprintf(&b, "%s: %s.\n", label, strings.Join(names, ", "))
		}
	}
	b.WriteString("\n")
	list("Performs", performTargets(para.Performs))
	list("Performed by", para.PerformedBy)
	list("Reached by", para.ReachedBy)
	list("Reads", para.Reads)
	list("Writes", para.Writes)
	for _, t := range p.Rules.Tables {
		if t.Paragraph != para.Name {
			continue
		}
		fmt.Fprintf(&b, "\n%s with %d rules. Conditions:\n- %s\nActions:\n- %s\n",
			t.Title(), len(t.Rules), strings.Join(t.Conditions, "\n- "), strings.Join(t.Actions, "\n- "))
	}
	return newFacts(KindParagraph, p.Name, para.Name, b.String()), nil
}

func newFacts(kind, program, subject, text string) Facts {
	sum := sha256.Sum256([]byte(text))
	return Facts{Kind: kind, Program: program, Subject: subject, Text: text, Hash: hex.EncodeToString(sum[:])}
}

// paragraphSource is the text of a paragraph's source: its lines from its
// header to the line before the next paragraph of the same file, or the
// file's end.
func (s *System) paragraphSource(para Paragraph) (string, error) {
	path, first, err := splitSource(para.Source)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	last := len(lines)
	for _, p := range s.Programs {
		for _, other := range p.Paragraphs {
			if otherPath, line, err := splitSource(other.Source); err == nil && otherPath == path && line > first && line-1 < last {
				last = line - 1
			}
		}
	}
	return strings.Join(lines[first-1:last], "\n"), nil
}

func splitSource(source string) (path string, line int, err error) {
	i := strings.LastIndexByte(source, ':')
	if i < 0 {
		return "", 0, fmt.Errorf("source %q has no line", source)
	}
	line, err = strconv.Atoi(source[i+1:])
	return source[:i], line, err
}

// codeName is a name a sentence writes as the code's: in backquotes, or in
// capitals with a hyphen, as COBOL names are.
var codeName = regexp.MustCompile("`([^`]+)`|\\b[A-Z0-9]+(?:-[A-Z0-9]+)+\\b")

// Check reports the names an annotation uses as the code's that its facts
// do not: a model's prose may name only what it was given.
func (f Facts) Check(a Annotation) error {
	known := strings.ToUpper(f.Text)
	var unknown []string
	for _, text := range []string{a.Name, a.Text} {
		for _, m := range codeName.FindAllStringSubmatch(text, -1) {
			name := m[0]
			if m[1] != "" {
				name = m[1]
			}
			if !containsName(known, strings.ToUpper(name)) && !slices.Contains(unknown, name) {
				unknown = append(unknown, name)
			}
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("%s %s %s: names %s, which its facts do not", f.Kind, f.Program, f.Subject, strings.Join(unknown, ", "))
	}
	return nil
}

// containsName reports whether text holds a name as a whole word of the
// code, not as a part of a longer one.
func containsName(text, name string) bool {
	for i := strings.Index(text, name); i >= 0; {
		before, after := i-1, i+len(name)
		if (before < 0 || !isNameByte(text[before])) && (after >= len(text) || !isNameByte(text[after])) {
			return true
		}
		next := strings.Index(text[i+1:], name)
		if next < 0 {
			break
		}
		i += 1 + next
	}
	return false
}

func isNameByte(c byte) bool {
	return c == '-' || c == '_' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}

// ReadAnnotations reads the annotations in a directory, a file for each
// program; none if there is no such directory.
func ReadAnnotations(dir string) ([]Annotation, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	var all []Annotation
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var annotations []Annotation
		if err := json.Unmarshal(data, &annotations); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		all = append(all, annotations...)
	}
	return all, nil
}

// WriteAnnotations writes annotations to a directory, a file for each
// program, in the order given, in place of those there before.
func WriteAnnotations(dir string, annotations []Annotation) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	old, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return err
	}
	for _, path := range old {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	byProgram := map[string][]Annotation{}
	var programs []string
	for _, a := range annotations {
		if _, ok := byProgram[a.Program]; !ok {
			programs = append(programs, a.Program)
		}
		byProgram[a.Program] = append(byProgram[a.Program], a)
	}
	for _, program := range programs {
		data, err := json.MarshalIndent(byProgram[program], "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, program+".json"), append(data, '\n'), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// Stale reports whether the facts an annotation was written from have
// changed since, or are no longer there.
func (s *System) Stale(a Annotation) bool {
	facts, err := s.ParagraphFacts(a.Program, a.Subject)
	return err != nil || facts.Hash != a.Hash
}

// Annotation is the annotation of a paragraph of a program, if a model
// wrote one.
func (s *System) Annotation(program, paragraph string) (Annotation, bool) {
	i := slices.IndexFunc(s.Annotations, func(a Annotation) bool {
		return a.Kind == KindParagraph && a.Program == program && a.Subject == paragraph
	})
	if i < 0 {
		return Annotation{}, false
	}
	return s.Annotations[i], true
}

// writeAnnotations writes what a model wrote about a program's paragraphs,
// if it wrote anything: each annotation whose facts are as they were, and
// for one whose facts have changed since, that it is stale.
func writeAnnotations(b *strings.Builder, s *System, p Program) {
	var mine []Annotation
	for _, a := range s.Annotations {
		if a.Kind == KindParagraph && a.Program == p.Name {
			mine = append(mine, a)
		}
	}
	if len(mine) == 0 {
		return
	}
	b.WriteString("\n## Prose by a model\n\nWhat a language model wrote about each paragraph from its facts and its source: " +
		"its reading of them, not a fact of the survey's. Prose whose facts have changed since it was written is stale, and not shown.\n\n" +
		"| Paragraph | Name | What it does | Written by |\n| --- | --- | --- | --- |\n")
	for _, a := range mine {
		name, text := a.Name, a.Text
		if s.Stale(a) {
			name, text = "", "stale: the facts have changed since it was written"
		}
		b.WriteString(row(a.Subject, name, text, a.Writer+", prompt "+a.Prompt))
	}
}
