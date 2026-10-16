// Package blueprint assembles the analysis of a COBOL system into
// blueprints: one per program and one per copybook, written as Markdown for
// people and as JSON for tools.
package blueprint

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vlumi/copal/cobol"
	"github.com/vlumi/copal/layout"
	"github.com/vlumi/copal/structure"
)

// System is the analysis of a set of programs and the copybooks they use.
type System struct {
	Programs  []Program
	Copybooks []Copybook
}

// Program is the blueprint of one program.
type Program struct {
	Name           string           `json:"program"`
	Source         string           `json:"source"`
	Description    []string         `json:"description,omitempty"`
	Identification []Identification `json:"identification,omitempty"`
	Files          []File           `json:"files"`
	Records        []Record         `json:"records"`
	Paragraphs     []Paragraph      `json:"paragraphs"`
}

// Identification is an identification paragraph such as AUTHOR.
type Identification struct {
	Paragraph string `json:"paragraph"`
	Text      string `json:"text"`
}

// File is a file the program declares, with the records of its FD or SD.
type File struct {
	Name         string   `json:"name"`
	Kind         string   `json:"kind"`
	Assign       string   `json:"assign"`
	Organization string   `json:"organization"`
	Status       string   `json:"status,omitempty"`
	Records      []string `json:"records"`
	Source       string   `json:"source"`
}

// Record is a level-01 or level-77 item and the fields under it. Section is
// "file" or "working-storage" in a program, empty in a copybook.
type Record struct {
	Section  string `json:"section,omitempty"`
	File     string `json:"file,omitempty"`
	Copybook string `json:"copybook,omitempty"`
	Field
}

// Field is one item of a record layout. Start is 1-based.
type Field struct {
	Level       int         `json:"level"`
	Name        string      `json:"name"`
	Start       int         `json:"start"`
	Length      int         `json:"length"`
	Picture     string      `json:"picture,omitempty"`
	Usage       string      `json:"usage,omitempty"`
	Category    string      `json:"category"`
	Digits      int         `json:"digits,omitempty"`
	Scale       int         `json:"scale,omitempty"`
	Signed      bool        `json:"signed,omitempty"`
	Occurs      int         `json:"occurs,omitempty"`
	OccursMin   int         `json:"occursMin,omitempty"`
	DependingOn string      `json:"dependingOn,omitempty"`
	IndexedBy   []string    `json:"indexedBy,omitempty"`
	Redefines   string      `json:"redefines,omitempty"`
	Values      []string    `json:"values,omitempty"`
	Description string      `json:"description,omitempty"`
	Conditions  []Condition `json:"conditions,omitempty"`
	Source      string      `json:"source"`
	Fields      []Field     `json:"fields,omitempty"`
}

// Condition is a level-88 condition name.
type Condition struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
	Source string   `json:"source"`
}

// Paragraph is a paragraph with what it performs, where it goes by GO TO
// or by falling into the next, and what performs it or comes to it.
type Paragraph struct {
	Name        string    `json:"name"`
	Section     string    `json:"section,omitempty"`
	Description string    `json:"description,omitempty"`
	Source      string    `json:"source"`
	Performs    []Perform `json:"performs,omitempty"`
	GoTos       []GoTo    `json:"goTos,omitempty"`
	FallsInto   string    `json:"fallsInto,omitempty"`
	PerformedBy []string  `json:"performedBy,omitempty"`
	ReachedBy   []string  `json:"reachedBy,omitempty"`
}

// GoTo is a GO TO statement's target.
type GoTo struct {
	Target string `json:"target"`
	Source string `json:"source"`
}

// Perform is a PERFORM statement naming a procedure or a THRU range, or a
// SORT statement's input or output procedure.
type Perform struct {
	Target string `json:"target"`
	Thru   string `json:"thru,omitempty"`
	Sort   bool   `json:"sort,omitempty"`
	Source string `json:"source"`
}

// Copybook is the blueprint of a copybook of record layouts.
type Copybook struct {
	Name        string   `json:"copybook"`
	Source      string   `json:"source"`
	Description []string `json:"description,omitempty"`
	Records     []Record `json:"records"`
	UsedBy      []string `json:"usedBy"`
}

// Analyze reads the programs and every copybook in copybookDir.
func Analyze(programPaths []string, copybookDir string) (*System, error) {
	s := &System{}
	for _, path := range programPaths {
		program, err := analyzeProgram(path, copybookDir)
		if err != nil {
			return nil, err
		}
		s.Programs = append(s.Programs, program)
	}
	slices.SortFunc(s.Programs, func(a, b Program) int { return strings.Compare(a.Name, b.Name) })
	copybookPaths, err := filepath.Glob(filepath.Join(copybookDir, "*.cpy"))
	if err != nil {
		return nil, err
	}
	for _, path := range copybookPaths {
		copybook, err := analyzeCopybook(path, copybookDir, s.Programs)
		if err != nil {
			return nil, err
		}
		s.Copybooks = append(s.Copybooks, copybook)
	}
	return s, nil
}

func analyzeProgram(path, copybookDir string) (Program, error) {
	p, err := cobol.ReadProgram(path, []string{copybookDir})
	if err != nil {
		return Program{}, err
	}
	program := Program{Name: p.ID, Source: path, Description: cobol.CommentText(p.Comment)}
	for _, e := range p.Identification {
		program.Identification = append(program.Identification, Identification{Paragraph: e.Paragraph, Text: e.Text})
	}

	for _, fd := range p.FileSection {
		records, err := layout.Records(fd.Entries)
		if err != nil {
			return Program{}, err
		}
		file, err := describeFile(p, fd)
		if err != nil {
			return Program{}, err
		}
		for _, r := range records {
			file.Records = append(file.Records, r.Name)
			program.Records = append(program.Records, newRecord(r, "file", fd.Name, path))
		}
		program.Files = append(program.Files, file)
	}

	records, err := layout.Records(p.WorkingStorage)
	if err != nil {
		return Program{}, err
	}
	for _, r := range records {
		program.Records = append(program.Records, newRecord(r, "working-storage", "", path))
	}

	paragraphs, err := structure.Paragraphs(p)
	if err != nil {
		return Program{}, err
	}
	for _, para := range paragraphs {
		program.Paragraphs = append(program.Paragraphs, newParagraph(para))
	}
	return program, nil
}

func describeFile(p *cobol.Program, fd cobol.FileDescription) (File, error) {
	i := slices.IndexFunc(p.Files, func(f cobol.FileControl) bool { return f.Name == fd.Name })
	if i < 0 {
		return File{}, fmt.Errorf("%s: %s %s has no SELECT", fd.Pos, fd.Kind, fd.Name)
	}
	selected := p.Files[i]
	return File{
		Name:         fd.Name,
		Kind:         fd.Kind,
		Assign:       selected.Assign,
		Organization: selected.Organization,
		Status:       selected.Status,
		Source:       selected.Pos.String(),
	}, nil
}

// newRecord converts a record layout. A record defined in a copybook names
// the copybook, so that blueprints can refer to its layout there.
func newRecord(f *layout.Field, section, file, programPath string) Record {
	r := Record{Section: section, File: file, Field: newField(f)}
	if f.Pos.File != programPath {
		r.Copybook = copybookName(f.Pos.File)
	}
	return r
}

func newField(f *layout.Field) Field {
	field := Field{
		Level:       f.Level,
		Name:        f.Name,
		Start:       f.Start,
		Length:      f.Length,
		Picture:     f.Picture,
		Category:    f.Category,
		Digits:      f.Digits,
		Scale:       f.Scale,
		Signed:      f.Signed,
		Occurs:      f.Occurs,
		OccursMin:   f.OccursMin,
		DependingOn: f.DependingOn,
		IndexedBy:   f.IndexedBy,
		Redefines:   f.Redefines,
		Values:      f.Values,
		Description: oneLine(f.Comment),
		Source:      f.Pos.String(),
	}
	if f.Category != layout.Group {
		field.Usage = f.Usage
	}
	for _, c := range f.Conditions {
		field.Conditions = append(field.Conditions, Condition{Name: c.Name, Values: c.Values, Source: c.Pos.String()})
	}
	for _, child := range f.Children {
		field.Fields = append(field.Fields, newField(child))
	}
	return field
}

func newParagraph(para structure.Paragraph) Paragraph {
	p := Paragraph{
		Name:        para.Name,
		Section:     para.Section,
		Description: oneLine(para.Comment),
		Source:      para.Pos.String(),
		FallsInto:   para.FallsInto,
		PerformedBy: para.PerformedBy,
		ReachedBy:   para.ReachedBy,
	}
	for _, perform := range para.Performs {
		p.Performs = append(p.Performs, Perform{Target: perform.Target, Thru: perform.Thru, Sort: perform.Sort, Source: perform.Pos.String()})
	}
	for _, goTo := range para.GoTos {
		p.GoTos = append(p.GoTos, GoTo{Target: goTo.Target, Source: goTo.Pos.String()})
	}
	return p
}

func analyzeCopybook(path, copybookDir string, programs []Program) (Copybook, error) {
	entries, err := cobol.ReadDataEntries(path, []string{copybookDir})
	if err != nil {
		return Copybook{}, err
	}
	records, err := layout.Records(entries)
	if err != nil {
		return Copybook{}, err
	}
	copybook := Copybook{Name: copybookName(path), Source: path, UsedBy: []string{}}
	for _, r := range records {
		copybook.Records = append(copybook.Records, Record{Field: newField(r)})
	}
	if len(records) > 0 {
		copybook.Description = cobol.CommentText(records[0].Comment)
		copybook.Records[0].Description = ""
	}
	for _, program := range programs {
		if slices.ContainsFunc(program.Records, func(r Record) bool { return r.Copybook == copybook.Name }) {
			copybook.UsedBy = append(copybook.UsedBy, program.Name)
		}
	}
	return copybook, nil
}

// oneLine joins the text of a short comment into one line.
func oneLine(comment []string) string {
	var words []string
	for _, line := range cobol.CommentText(comment) {
		if line != "" {
			words = append(words, strings.TrimSpace(line))
		}
	}
	return strings.Join(words, " ")
}

func copybookName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
