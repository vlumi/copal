// Package dataflow finds what a program does with its data: which data
// items each paragraph reads and writes, and the lineage of every item,
// the statements through which values flow into it from other items and
// from files.
//
// Items are named as declared. A condition name stands for the field it
// is declared on: SET TD-INTEREST TO TRUE writes TD-TRAN-CODE. Names the
// program does not declare, such as RETURN-CODE or an index name, are kept
// as they are.
package dataflow

import (
	"regexp"
	"slices"
	"strconv"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/layout"
)

// Flow is what a program does with its data.
type Flow struct {
	Paragraphs []Paragraph
	Lineage    []Edge
	storage    *storage
	// parameters are the positions of the items the program is passed
	// with USING, 1-based, by name.
	parameters map[string]int
	// order is the order its statements run in, nil where the paths
	// through them are not known.
	order *order
	// program and names are the program analyzed and how its names
	// resolve, for what its statements change.
	program *cobol.Program
	names   *names
	// callWrites says what a CALL writes, where it is known.
	callWrites CallWrites
}

// Paragraph is a paragraph's reads and writes, each item once, in the
// order the paragraph first uses it.
type Paragraph struct {
	Name   string
	Reads  []string
	Writes []string
}

// Edge is one statement moving a value from one item into another. From
// is a file for the records a READ or RETURN fills, and To a file for the
// record a WRITE, REWRITE or RELEASE sends to it. From is a literal, as written,
// for a value the statement gives itself, such as MOVE 'TXPOST' TO
// RJ-PROGRAM or SET TD-INTEREST TO TRUE, written as "condition
// TD-INTEREST".
//
// A CALL moves each item it passes into the program called, To with
// ToCall, and may move a value back into each item it passes by
// reference, From with FromCall; Arg is the argument's position, 1-based.
// It moves what the program returns into its RETURNING item, From with
// FromCall and Arg 0.
// An ACCEPT moves a value from the job's SYSIN, From with FromStream, and
// a DISPLAY one to its SYSOUT, To with ToStream.
type Edge struct {
	From, To string
	// FromStart and FromLength are the bytes of From the statement moves,
	// 1-based within it, when it names them by reference modification,
	// X(1:8); ToStart and ToLength the same of To. Zero for the whole.
	FromStart, FromLength int
	ToStart, ToLength     int
	FromFile, ToFile      bool
	FromLiteral           bool
	FromCall, ToCall      bool
	FromStream, ToStream  bool
	Arg                   int
	Paragraph             string
	Verb                  string
	Pos                   cobol.Pos
}

// Analyze finds the program's data flow, given its records laid out.
func Analyze(p *cobol.Program, records *layout.ProgramRecords) Flow {
	names := newNames(p, records)
	flow := Flow{storage: newStorage(records), parameters: map[string]int{}, program: p, names: names}
	for i, item := range p.Using {
		flow.parameters[item] = i + 1
	}
	byStatement := map[cobol.Pos]statementFacts{}
	for _, para := range p.Paragraphs {
		pf := Paragraph{Name: para.Name}
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			reads, writes := names.uses(s)
			for _, r := range reads {
				addOnce(&pf.Reads, r)
			}
			for _, w := range writes {
				addOnce(&pf.Writes, w)
			}
			sf := statementFacts{reads: reads, kills: names.kills(s, reads)}
			for _, e := range names.edges(para.Name, s) {
				sf.edges = append(sf.edges, len(flow.Lineage))
				flow.Lineage = append(flow.Lineage, e)
			}
			byStatement[s.Pos] = sf
		})
		flow.Paragraphs = append(flow.Paragraphs, pf)
	}
	flow.order = newOrder(p, flow.storage, byStatement)
	return flow
}

// Write is a statement that writes a record to a file: the record named in
// its FROM, or the file's own record.
type Write struct {
	File, Record string
	Pos          cobol.Pos
}

// Writes lists the statements that write records to files, in the order
// the program declares them.
func (f Flow) Writes() []Write {
	var writes []Write
	for _, e := range f.Lineage {
		if !e.ToFile {
			continue
		}
		w := Write{File: e.To, Record: e.From, Pos: e.Pos}
		for _, from := range f.Lineage {
			if from.Pos == e.Pos && from.To == e.From && !from.FromLiteral && !from.FromFile && !from.FromCall && !from.FromStream {
				w.Record = from.From
			}
		}
		writes = append(writes, w)
	}
	return writes
}

func addOnce(names *[]string, name string) {
	if !slices.Contains(*names, name) {
		*names = append(*names, name)
	}
}

// names resolves the names a program uses: condition names to their
// fields, files to their records and records to their files.
type names struct {
	conditions map[string]string
	records    map[string][]string // file to records
	files      map[string]string   // record to file
	status     map[string]string   // file to its FILE STATUS item
	// reports maps a report of the REPORT SECTION, and each of its
	// groups, to the file that prints it, and lines keeps that file's
	// records, the lines of its reports, with their columns: what a
	// GENERATE or TERMINATE reads through SOURCE and SUM and writes.
	reports map[string]string
	lines   map[string][]*layout.Field
}

func newNames(p *cobol.Program, records *layout.ProgramRecords) *names {
	n := &names{conditions: map[string]string{}, records: map[string][]string{}, files: map[string]string{},
		status: map[string]string{}, reports: map[string]string{}, lines: map[string][]*layout.Field{}}
	for _, fc := range p.Files {
		if fc.Status != "" {
			n.status[fc.Name] = fc.Status
		}
	}
	for _, f := range records.Files {
		for _, r := range f.Records {
			n.records[f.File] = append(n.records[f.File], r.Name)
			n.files[r.Name] = f.File
		}
	}
	for _, fd := range p.FileSection {
		for _, report := range fd.Reports {
			n.reports[report] = fd.Name
			for _, rd := range p.ReportSection {
				if rd.Name != report {
					continue
				}
				for _, group := range rd.Entries {
					if group.Level == 1 && group.Name != "FILLER" {
						n.reports[group.Name] = fd.Name
					}
				}
			}
			for _, f := range records.Files {
				if f.File == fd.Name {
					n.lines[fd.Name] = f.Records
				}
			}
		}
	}
	var index func(f *layout.Field)
	index = func(f *layout.Field) {
		for _, c := range f.Conditions {
			n.conditions[c.Name] = f.Name
		}
		for _, child := range f.Children {
			index(child)
		}
	}
	for _, r := range records.All() {
		index(r)
	}
	return n
}

// item is the data item a name stands for.
func (n *names) item(name string) string {
	if field, ok := n.conditions[name]; ok {
		return field
	}
	return name
}

// filesOf lists the files a statement operates on, whose status it sets:
// those it names, and the file of each record it writes.
func (n *names) filesOf(s cobol.Statement) []string {
	var files []string
	switch s.Verb {
	case "OPEN", "CLOSE", "READ", "START", "DELETE":
		for _, f := range s.Files {
			files = append(files, f.Name)
		}
	case "WRITE", "REWRITE":
		for _, record := range s.Writes {
			if file, ok := n.files[record.Name]; ok {
				files = append(files, file)
			}
		}
	}
	return files
}

// referenceModified matches an operand's reference modification by
// literal position and length, X(1:8).
var referenceModified = regexp.MustCompile(`\((\d+):(\d+)\)$`)

// referenceModification gives the bytes an operand names of its item by
// reference modification, 1-based, and 0 for the whole item or for bytes
// named by data items.
func referenceModification(op cobol.Operand) (start, length int) {
	m := referenceModified.FindStringSubmatch(op.Text)
	if m == nil {
		return 0, 0
	}
	start, _ = strconv.Atoi(m[1])
	length, _ = strconv.Atoi(m[2])
	return start, length
}

// uses lists the items a statement reads and writes: its operands and
// condition, a written operand's subscripts read, and the records of the
// files it reads.
func (n *names) uses(s cobol.Statement) (reads, writes []string) {
	for _, op := range s.Reads {
		for _, ref := range op.Refs {
			reads = append(reads, n.item(ref))
		}
	}
	for _, op := range s.Writes {
		if op.Name != "" {
			writes = append(writes, n.item(op.Name))
		}
		for _, ref := range op.Refs {
			if ref != op.Name {
				reads = append(reads, n.item(ref))
			}
		}
	}
	if s.Condition != nil {
		for _, ref := range s.Condition.Refs() {
			reads = append(reads, n.item(ref))
		}
	}
	for _, b := range s.Branches {
		if b.Condition != nil {
			for _, ref := range b.Condition.Refs() {
				reads = append(reads, n.item(ref))
			}
		}
	}
	if s.Verb == "READ" || s.Verb == "RETURN" {
		for _, f := range s.Files {
			writes = append(writes, n.records[f.Name]...)
		}
	}
	for _, file := range n.filesOf(s) {
		if status, ok := n.status[file]; ok {
			writes = append(writes, n.item(status))
		}
	}
	if file, ok := n.reportFile(s); ok {
		for _, line := range n.lines[file] {
			for _, column := range line.Children {
				if column.Source != "" {
					reads = append(reads, n.item(column.Source))
				}
				for _, item := range column.Sum {
					reads = append(reads, n.item(item))
				}
			}
			writes = append(writes, line.Name)
		}
	}
	return reads, writes
}

// reportFile is the file a GENERATE or TERMINATE prints to, by the report
// or group it names.
func (n *names) reportFile(s cobol.Statement) (string, bool) {
	if s.Verb != "GENERATE" && s.Verb != "TERMINATE" || len(s.Reads) == 0 {
		return "", false
	}
	file, ok := n.reports[s.Reads[0].Name]
	return file, ok
}

// reportEdges lists the values a GENERATE or TERMINATE moves: each
// column's SOURCE item and SUM items into the column, and each line of the
// report into its file.
func (n *names) reportEdges(paragraph string, s cobol.Statement, file string) []Edge {
	var edges []Edge
	add := func(from, to string, toFile bool) {
		edges = append(edges, Edge{From: from, To: to, ToFile: toFile, Paragraph: paragraph, Verb: s.Verb, Pos: s.Pos})
	}
	for _, line := range n.lines[file] {
		for _, column := range line.Children {
			if column.Source != "" {
				add(n.item(column.Source), column.Name, false)
			}
			for _, item := range column.Sum {
				add(n.item(item), column.Name, false)
			}
		}
		add(line.Name, file, true)
	}
	return edges
}

// lineageVerbs move values from what they read into what they write.
var lineageVerbs = map[string]bool{
	"MOVE": true, "ADD": true, "SUBTRACT": true, "MULTIPLY": true, "DIVIDE": true, "COMPUTE": true,
	"STRING": true, "SET": true, "WRITE": true, "REWRITE": true, "RELEASE": true, "READ": true, "RETURN": true,
}

// edges lists the values a statement moves, each once. A data item operand
// moves its own value, not its subscripts'; an expression moves every item
// in it.
func (n *names) edges(paragraph string, s cobol.Statement) []Edge {
	switch {
	case s.Verb == "CALL" && s.Calls != "":
		return n.callEdges(paragraph, s)
	case s.Verb == "ACCEPT" && (s.Device == "" || s.Device == CommandLine), s.Verb == "DISPLAY" && s.Device == "":
		return n.streamEdges(paragraph, s)
	case !lineageVerbs[s.Verb]:
		if file, ok := n.reportFile(s); ok {
			return n.reportEdges(paragraph, s, file)
		}
		return nil
	}
	var edges []Edge
	add := func(e Edge) {
		e.Paragraph, e.Verb, e.Pos = paragraph, s.Verb, s.Pos
		if (e.FromFile || e.ToFile || e.FromLiteral || e.From != e.To) && !slices.Contains(edges, e) {
			edges = append(edges, e)
		}
	}
	var sources []Edge
	for _, op := range s.Reads {
		switch {
		case op.Name != "":
			start, length := referenceModification(op)
			sources = append(sources, Edge{From: n.item(op.Name), FromStart: start, FromLength: length})
		case len(op.Refs) == 0:
			sources = append(sources, Edge{From: op.Text, FromLiteral: true})
		}
		if op.Name == "" {
			for _, ref := range op.Refs {
				sources = append(sources, Edge{From: n.item(ref)})
			}
		}
	}
	for _, target := range s.Writes {
		if target.Name == "" {
			continue
		}
		if s.Verb == "SET" && len(s.Reads) == 0 {
			add(Edge{From: "condition " + target.Name, FromLiteral: true, To: n.item(target.Name)})
		}
		for _, source := range sources {
			source.To = n.item(target.Name)
			source.ToStart, source.ToLength = referenceModification(target)
			add(source)
		}
	}
	switch s.Verb {
	case "READ", "RETURN":
		for _, f := range s.Files {
			for _, record := range n.records[f.Name] {
				add(Edge{From: f.Name, FromFile: true, To: record})
				for _, into := range s.Writes {
					add(Edge{From: record, To: n.item(into.Name)})
				}
			}
		}
	case "WRITE", "REWRITE", "RELEASE":
		for _, record := range s.Writes {
			if file, ok := n.files[record.Name]; ok {
				add(Edge{From: record.Name, To: file, ToFile: true})
			}
		}
	}
	return edges
}

// callEdges lists the values a CALL moves: each item it passes into the
// program called, and back from it into each item passed by reference,
// which the program may change.
func (n *names) callEdges(paragraph string, s cobol.Statement) []Edge {
	var edges []Edge
	for i, op := range s.Arguments() {
		if op.Name == "" {
			continue
		}
		item := n.item(op.Name)
		edges = append(edges, Edge{From: item, To: s.Calls, ToCall: true, Arg: i + 1, Paragraph: paragraph, Verb: s.Verb, Pos: s.Pos})
		if slices.ContainsFunc(s.Writes, func(w cobol.Operand) bool { return w.Text == op.Text }) {
			edges = append(edges, Edge{From: s.Calls, FromCall: true, Arg: i + 1, To: item, Paragraph: paragraph, Verb: s.Verb, Pos: s.Pos})
		}
	}
	if s.Returning != "" {
		edges = append(edges, Edge{From: s.Calls, FromCall: true, To: n.item(s.Returning), Paragraph: paragraph, Verb: s.Verb, Pos: s.Pos})
	}
	return edges
}

// CommandLine is the stream of a program's command line, which ACCEPT ...
// FROM COMMAND-LINE reads.
const CommandLine = "COMMAND-LINE"

// streamEdges lists the values an ACCEPT moves from the job's SYSIN or
// the program's command line, or a DISPLAY to its SYSOUT.
func (n *names) streamEdges(paragraph string, s cobol.Statement) []Edge {
	var edges []Edge
	if s.Verb == "ACCEPT" {
		stream := "SYSIN"
		if s.Device == CommandLine {
			stream = CommandLine
		}
		for _, op := range s.Writes {
			if op.Name != "" {
				edges = append(edges, Edge{From: stream, FromStream: true, To: n.item(op.Name), Paragraph: paragraph, Verb: s.Verb, Pos: s.Pos})
			}
		}
		return edges
	}
	for _, op := range s.Reads {
		if op.Name != "" {
			edges = append(edges, Edge{From: n.item(op.Name), To: "SYSOUT", ToStream: true, Paragraph: paragraph, Verb: s.Verb, Pos: s.Pos})
		}
	}
	return edges
}
