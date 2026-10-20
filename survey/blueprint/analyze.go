// Package blueprint assembles the analysis of a COBOL system into
// blueprints: one per program, per copybook and per job, written as
// Markdown for people and as JSON for tools.
package blueprint

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
	"github.com/vlumi/copal/survey/internal/dataflow"
	"github.com/vlumi/copal/survey/internal/fileio"
	"github.com/vlumi/copal/survey/internal/rules"
	"github.com/vlumi/copal/survey/internal/structure"
)

// System is the analysis of a set of programs, the copybooks they use and
// the jobs that run them.
type System struct {
	Dialect   cobol.Dialect
	Programs  []Program
	Copybooks []Copybook
	Jobs      []Job
	Datasets  []SystemDataset
	Entities  []Entity
	Schedule  []JobOrder
	// Annotations are the prose a model wrote about the system, which the
	// blueprints show beside the facts; Analyze leaves them to be read.
	Annotations []Annotation
}

// Program is the blueprint of one program. Its Kind is "called" for a
// program that is passed items when called, "main" for one that is not.
type Program struct {
	Name   string `json:"program"`
	Kind   string `json:"kind"`
	Source string `json:"source"`
	// Within is the program this one is nested in, if it is, which alone
	// can call it.
	Within         string           `json:"within,omitempty"`
	Description    []string         `json:"description,omitempty"`
	Identification []Identification `json:"identification,omitempty"`
	Files          []File           `json:"files"`
	Constants      []Constant       `json:"constants,omitempty"`
	Records        []Record         `json:"records"`
	Paragraphs     []Paragraph      `json:"paragraphs"`
	// Declaratives are the sections the runtime runs when an input or
	// output statement fails on a file.
	Declaratives []Declarative `json:"declaratives,omitempty"`
	// ReadsSysin is set for a program that ACCEPTs from the job's SYSIN.
	ReadsSysin bool `json:"readsSysin,omitempty"`
	// Displays are the DISPLAY statements to the job's output.
	Displays []Display `json:"displays,omitempty"`
	// Using are the linkage items the program is passed when called.
	Using    []string  `json:"using,omitempty"`
	Calls    []Call    `json:"calls,omitempty"`
	CalledBy []string  `json:"calledBy,omitempty"`
	Sorts    []Sort    `json:"sorts,omitempty"`
	Lineage  []Lineage `json:"lineage,omitempty"`
	Outputs  []Output  `json:"outputs,omitempty"`
	Inputs   []Input   `json:"inputs,omitempty"`
	// Links are the relations between its inputs' records it relies on,
	// Totals the control totals it checks, and Parameters, for a called
	// program, what it does with what it is passed.
	Links      []Link      `json:"links,omitempty"`
	Totals     []Total     `json:"totals,omitempty"`
	Parameters []Parameter `json:"parameters,omitempty"`
	// Rules are the program's decision tables, a blueprint of their own.
	Rules  Rules `json:"-"`
	flow   dataflow.Flow
	source *cobol.Program
	tables []rules.Table
}

// Sort is a SORT statement: its work file, the keys it sorts on, and the files it reads with
// USING and writes with GIVING, none on a side it has a procedure for.
type Sort struct {
	Work   string   `json:"work"`
	Keys   []string `json:"keys,omitempty"`
	Using  []string `json:"using,omitempty"`
	Giving []string `json:"giving,omitempty"`
	Source string   `json:"source"`
}

// Tables are the program's decision tables, as the rules analysis finds
// them, with the positions of the statements on each rule's path.
func (p Program) Tables() []rules.Table { return p.tables }

// Trace traces where a data item's values come from within the program,
// through the statements that move them.
func (p Program) Trace(item string) dataflow.Trace {
	return p.flow.Trace(item)
}

// TraceBytes traces where the values of some of a data item's bytes come
// from, from start, 1-based, for length.
func (p Program) TraceBytes(item string, start, length int) dataflow.Trace {
	return p.flow.TraceBytes(item, start, length)
}

// Lineage is one statement moving a value from one data item, or a file,
// into another, or into a file.
type Lineage struct {
	From        string `json:"from"`
	To          string `json:"to"`
	FromFile    bool   `json:"fromFile,omitempty"`
	FromLiteral bool   `json:"fromLiteral,omitempty"`
	FromCall    bool   `json:"fromCall,omitempty"`
	FromStream  bool   `json:"fromStream,omitempty"`
	ToFile      bool   `json:"toFile,omitempty"`
	ToCall      bool   `json:"toCall,omitempty"`
	ToStream    bool   `json:"toStream,omitempty"`
	Arg         int    `json:"arg,omitempty"`
	Paragraph   string `json:"paragraph"`
	Verb        string `json:"verb"`
	Source      string `json:"source"`
}

// NewLineage is a statement's move of a value, as data flow finds it.
func NewLineage(e dataflow.Edge) Lineage {
	return Lineage{From: e.From, To: e.To, FromFile: e.FromFile, FromLiteral: e.FromLiteral, FromCall: e.FromCall,
		FromStream: e.FromStream, ToFile: e.ToFile, ToCall: e.ToCall, ToStream: e.ToStream, Arg: e.Arg,
		Paragraph: e.Paragraph, Verb: e.Verb, Source: e.Pos.String()}
}

// FromItem reports whether the value comes from a data item, rather than
// a file, a literal, a program called or a stream.
func (e Lineage) FromItem() bool {
	return !e.FromFile && !e.FromLiteral && !e.FromCall && !e.FromStream
}

// ToItem reports whether the value goes into a data item, rather than a
// file, a program called or a stream.
func (e Lineage) ToItem() bool { return !e.ToFile && !e.ToCall && !e.ToStream }

// FromText writes where the value comes from when it is not a data item:
// "file TRAN-IN", "literal 1", "call AUTHPUR argument 6", "SYSIN".
func (e Lineage) FromText() string {
	switch {
	case e.FromFile:
		return "file " + e.From
	case e.FromLiteral:
		return "literal " + e.From
	case e.FromCall && e.Arg == 0:
		return fmt.Sprintf("call %s, what it returns", e.From)
	case e.FromCall:
		return fmt.Sprintf("call %s argument %d", e.From, e.Arg)
	}
	return e.From
}

// ToText writes where the value goes when it is not a data item: "file
// TRAN-OUT", "call AUTHPUR argument 6", "SYSOUT".
func (e Lineage) ToText() string {
	switch {
	case e.ToFile:
		return "file " + e.To
	case e.ToCall:
		return fmt.Sprintf("call %s argument %d", e.To, e.Arg)
	}
	return e.To
}

// Identification is an identification paragraph such as AUTHOR.
type Identification struct {
	Paragraph string `json:"paragraph"`
	Text      string `json:"text"`
}

// File is a file the program declares, with the records of its FD or SD
// and the modes the program opens it in: INPUT, OUTPUT, EXTEND or I-O,
// counting a SORT's USING files as input and its GIVING files as output.
type File struct {
	Name         string `json:"name"`
	Kind         string `json:"kind"`
	Assign       string `json:"assign"`
	Organization string `json:"organization"`
	Access       string `json:"access,omitempty"` // RANDOM or DYNAMIC; empty for SEQUENTIAL
	RecordKey    string `json:"recordKey,omitempty"`
	RelativeKey  string `json:"relativeKey,omitempty"`
	Status       string `json:"status,omitempty"`
	// MinLength and MaxLength are the record lengths the file's RECORD
	// clause gives, if it has one, and LengthIn the item a variable
	// length is in.
	MinLength int      `json:"minLength,omitempty"`
	MaxLength int      `json:"maxLength,omitempty"`
	LengthIn  string   `json:"lengthIn,omitempty"`
	Opened    []string `json:"opened,omitempty"`
	Records   []string `json:"records"`
	Source    string   `json:"source"`
}

// Record is a level-01 or level-77 item and the fields under it. Section is
// "file", "working-storage" or "linkage" in a program, empty in a copybook.
type Record struct {
	Section  string `json:"section,omitempty"`
	File     string `json:"file,omitempty"`
	Copybook string `json:"copybook,omitempty"`
	// TypeField is the field whose condition names tell the record's
	// types apart, and Types those types, as the program tests them.
	TypeField string       `json:"typeField,omitempty"`
	Types     []RecordType `json:"types,omitempty"`
	Field
}

// Field is one item of a record layout. Start is 1-based. Kind is what
// it holds, as fieldKind tells it.
type Field struct {
	Level       int         `json:"level"`
	Name        string      `json:"name"`
	Start       int         `json:"start"`
	Length      int         `json:"length"`
	Picture     string      `json:"picture,omitempty"`
	Usage       string      `json:"usage,omitempty"`
	Category    string      `json:"category"`
	Kind        string      `json:"kind"`
	Digits      int         `json:"digits,omitempty"`
	Scale       int         `json:"scale,omitempty"`
	Signed      bool        `json:"signed,omitempty"`
	Sign        string      `json:"sign,omitempty"` // LEADING or TRAILING, and SEPARATE, as a SIGN clause says
	Occurs      int         `json:"occurs,omitempty"`
	OccursMin   int         `json:"occursMin,omitempty"`
	DependingOn string      `json:"dependingOn,omitempty"`
	IndexedBy   []string    `json:"indexedBy,omitempty"`
	Redefines   string      `json:"redefines,omitempty"`
	RedefinedBy []string    `json:"redefinedBy,omitempty"` // the fields that redefine this one, in order
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
	Reads       []string  `json:"reads,omitempty"`
	Writes      []string  `json:"writes,omitempty"`
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

// Analyze reads the programs, the jobs, which are the paths ending in
// .jcl, and the batch scripts, in .sh, and every copybook in copybookDir.
func Analyze(paths []string, copybookDir string, dialect cobol.Dialect) (*System, error) {
	s := &System{Dialect: dialect}
	var programPaths, jobPaths []string
	for _, path := range paths {
		if ext := strings.ToLower(filepath.Ext(path)); ext == ".jcl" || ext == ".sh" {
			jobPaths = append(jobPaths, path)
		} else {
			programPaths = append(programPaths, path)
		}
	}
	for _, path := range programPaths {
		programs, err := analyzeFile(path, copybookDir, dialect)
		if err != nil {
			return nil, err
		}
		s.Programs = append(s.Programs, programs...)
	}
	slices.SortFunc(s.Programs, func(a, b Program) int { return strings.Compare(a.Name, b.Name) })
	linkCalls(s.Programs)
	for i := range s.Programs {
		p := &s.Programs[i]
		flow := p.flow.WithCalls(callWrites(*p, s.Programs))
		p.tables = rules.Tables(p.source, rules.Knowledge{
			MayChange:      flow.MayChange,
			Writers:        flow.Writers,
			ConditionNames: conditionNames(p.Records),
			InitialValues:  initialValues(p.Records),
			TextFields:     textFields(p.Records),
			DigitFields:    digitFields(p.Records),
			Copies:         copies(p.Records),
		})
		p.Rules = newRules(p.Name, p.tables, p.Records)
	}
	var copybookPaths []string
	for _, ext := range cobol.CopybookExtensions {
		paths, err := filepath.Glob(filepath.Join(copybookDir, "*"+ext))
		if err != nil {
			return nil, err
		}
		copybookPaths = append(copybookPaths, paths...)
	}
	slices.Sort(copybookPaths)
	for _, path := range copybookPaths {
		copybook, err := analyzeCopybook(path, copybookDir, dialect, s.Programs)
		if err != nil {
			return nil, err
		}
		s.Copybooks = append(s.Copybooks, copybook)
	}
	s.Entities = entities(s.Copybooks, s.Programs)
	for _, path := range jobPaths {
		analyze := analyzeJob
		if strings.EqualFold(filepath.Ext(path), ".sh") {
			analyze = analyzeScript
		}
		job, err := analyze(path, s.Programs)
		if err != nil {
			return nil, err
		}
		s.Jobs = append(s.Jobs, job)
	}
	slices.SortFunc(s.Jobs, func(a, b Job) int { return strings.Compare(a.Name, b.Name) })
	describeDatasets(s.Jobs)
	s.Datasets = systemDatasets(s.Jobs)
	classifyJobDatasets(s.Jobs, s.Datasets)
	describeEntities(s.Entities, s.Datasets, s.Jobs)
	s.Schedule = jobOrder(s.Jobs)
	return s, nil
}

// analyzeFile reads a program's source and describes the program and the
// programs nested in it.
func analyzeFile(path, copybookDir string, dialect cobol.Dialect) ([]Program, error) {
	p, err := cobol.ReadProgram(path, []string{copybookDir})
	if err != nil {
		return nil, err
	}
	return describeNested(p, "", path, dialect)
}

// describeNested describes a program, nested in another if within names
// one, and the programs nested in it.
func describeNested(p *cobol.Program, within, path string, dialect cobol.Dialect) ([]Program, error) {
	program, err := describeProgram(p, path, dialect)
	if err != nil {
		return nil, err
	}
	program.Within = within
	programs := []Program{program}
	for _, inner := range p.Contains {
		nested, err := describeNested(inner, p.ID, path, dialect)
		if err != nil {
			return nil, err
		}
		programs = append(programs, nested...)
	}
	return programs, nil
}

func describeProgram(p *cobol.Program, path string, dialect cobol.Dialect) (Program, error) {
	program := Program{Name: p.ID, Source: path, Description: cobol.CommentText(p.Comment)}
	for _, e := range p.Identification {
		program.Identification = append(program.Identification, Identification{Paragraph: e.Paragraph, Text: e.Text})
	}

	program.Constants = constants(p)
	program.Declaratives = declaratives(p)
	laidOut, err := layout.Program(p, dialect)
	if err != nil {
		return Program{}, err
	}
	uses := fileio.Uses(p)
	for i, fd := range p.FileSection {
		file, err := describeFile(p, fd)
		if err != nil {
			return Program{}, err
		}
		for _, u := range uses {
			if u.File == fd.Name {
				file.Opened = append(file.Opened, u.Mode)
			}
		}
		for _, r := range laidOut.Files[i].Records {
			file.Records = append(file.Records, r.Name)
		}
		program.Files = append(program.Files, file)
	}
	for _, s := range fileio.Sorts(p) {
		program.Sorts = append(program.Sorts, Sort{Work: s.Work, Keys: s.Keys, Using: s.Using, Giving: s.Giving, Source: s.Pos.String()})
	}

	keys := programKeys(program.Files, program.Sorts)
	for _, f := range laidOut.Files {
		for _, r := range f.Records {
			program.Records = append(program.Records, newRecord(r, "file", f.File, path, keys))
		}
	}
	for _, r := range laidOut.WorkingStorage {
		program.Records = append(program.Records, newRecord(r, "working-storage", "", path, keys))
	}
	for _, r := range laidOut.Linkage {
		program.Records = append(program.Records, newRecord(r, "linkage", "", path, keys))
	}
	recordTypes(p, program.Records)
	program.Using = p.Using
	program.Kind = programKind(p.Using)
	program.ReadsSysin = readsSysin(p)
	program.Displays = displays(p)

	paragraphs, err := structure.Paragraphs(p)
	if err != nil {
		return Program{}, err
	}
	flow := dataflow.Analyze(p, laidOut)
	program.flow = flow
	program.Calls = programCalls(p, program, flow)
	for i, para := range paragraphs {
		bp := newParagraph(para)
		bp.Reads, bp.Writes = flow.Paragraphs[i].Reads, flow.Paragraphs[i].Writes
		program.Paragraphs = append(program.Paragraphs, bp)
	}
	for _, e := range flow.Lineage {
		program.Lineage = append(program.Lineage, NewLineage(e))
	}
	program.Outputs = outputs(program, flow)
	program.Inputs = inputs(program, flow)
	program.Links = inputLinks(program, p, flow)
	program.Totals = programTotals(program, p, flow)
	program.Parameters = parameters(program)
	program.source = p
	return program, nil
}

func readsSysin(p *cobol.Program) bool {
	reads := false
	for _, para := range p.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			reads = reads || s.Verb == "ACCEPT" && s.Device == ""
		})
	}
	return reads
}

// Display is a DISPLAY statement to the job's output: the items and
// literals it writes, as written.
type Display struct {
	Items     []string `json:"items"`
	Paragraph string   `json:"paragraph"`
	Source    string   `json:"source"`
}

// displays lists the program's DISPLAY statements to the job's output,
// SYSOUT, leaving out those UPON another device.
func displays(p *cobol.Program) []Display {
	var list []Display
	for _, para := range p.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if s.Verb != "DISPLAY" || s.Device != "" {
				return
			}
			d := Display{Paragraph: para.Name, Source: s.Pos.String()}
			for _, op := range s.Reads {
				d.Items = append(d.Items, op.Text)
			}
			list = append(list, d)
		})
	}
	return list
}

func describeFile(p *cobol.Program, fd cobol.FileDescription) (File, error) {
	i := slices.IndexFunc(p.Files, func(f cobol.FileControl) bool { return f.Name == fd.Name })
	if i < 0 {
		return File{}, fmt.Errorf("%s: %s %s has no SELECT", fd.Pos, fd.Kind, fd.Name)
	}
	selected := p.Files[i]
	f := File{
		Name:         fd.Name,
		Kind:         fd.Kind,
		Assign:       selected.Assign,
		Organization: selected.Organization,
		RecordKey:    selected.RecordKey,
		RelativeKey:  selected.RelativeKey,
		Status:       selected.Status,
		MinLength:    fd.MinLength,
		MaxLength:    fd.MaxLength,
		LengthIn:     fd.LengthIn,
		Source:       selected.Pos.String(),
	}
	if selected.Access != "SEQUENTIAL" {
		f.Access = selected.Access
	}
	return f, nil
}

// newRecord converts a record layout, with the program's keys. A record
// defined in a copybook names the copybook, so that blueprints can refer
// to its layout there.
func newRecord(f *layout.Field, section, file, programPath string, keys map[string]bool) Record {
	r := Record{Section: section, File: file, Field: newField(f, keys)}
	if f.Pos.File != programPath {
		r.Copybook = copybookName(f.Pos.File)
	}
	return r
}

// Declarative is a section of the DECLARATIVES, and what it is for: the
// files named, or every file opened in the modes named.
type Declarative struct {
	Section string   `json:"section"`
	Files   []string `json:"files,omitempty"`
	Modes   []string `json:"modes,omitempty"`
	Source  string   `json:"source"`
}

// declaratives lists a program's declaratives.
func declaratives(p *cobol.Program) []Declarative {
	var list []Declarative
	for _, d := range p.Declaratives {
		list = append(list, Declarative{Section: d.Section, Files: d.Files, Modes: d.Modes, Source: d.Pos.String()})
	}
	return list
}

// Constant is a level-78 constant: a name for a value.
type Constant struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Source string `json:"source"`
}

// constants lists a program's level-78 constants.
func constants(p *cobol.Program) []Constant {
	var list []Constant
	for _, e := range p.WorkingStorage {
		if e.Level == 78 && len(e.Values) == 1 {
			list = append(list, Constant{Name: e.Name, Value: e.Values[0], Source: e.Pos.String()})
		}
	}
	return list
}

// signClause is where a field's sign is, as its SIGN clause says.
func signClause(f *layout.Field) string {
	if f.SignSeparate {
		return f.Sign + " SEPARATE"
	}
	return f.Sign
}

func newField(f *layout.Field, keys map[string]bool) Field {
	field := Field{
		Level:       f.Level,
		Name:        f.Name,
		Start:       f.Start,
		Length:      f.Length,
		Picture:     f.Picture,
		Category:    f.Category,
		Kind:        fieldKind(f, keys),
		Digits:      f.Digits,
		Scale:       f.Scale,
		Signed:      f.Signed,
		Sign:        signClause(f),
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
		field.Fields = append(field.Fields, newField(child, keys))
	}
	for _, child := range field.Fields {
		if child.Redefines == "" {
			continue
		}
		for i := range field.Fields {
			if field.Fields[i].Name == child.Redefines {
				field.Fields[i].RedefinedBy = append(field.Fields[i].RedefinedBy, child.Name)
			}
		}
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

func analyzeCopybook(path, copybookDir string, dialect cobol.Dialect, programs []Program) (Copybook, error) {
	entries, err := cobol.ReadDataEntries(path, []string{copybookDir})
	if err != nil {
		return Copybook{}, err
	}
	records, err := layout.Records(entries, dialect)
	if err != nil {
		return Copybook{}, err
	}
	copybook := Copybook{Name: copybookName(path), Source: path, UsedBy: []string{}}
	keys := copybookKeys(copybook.Name, programs)
	for _, r := range records {
		copybook.Records = append(copybook.Records, Record{Field: newField(r, keys)})
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
