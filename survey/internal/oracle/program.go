// Package oracle finds input data that reaches the rules of a program's
// decision tables, by running the program, compiled with GnuCOBOL's trace
// code, on candidate inputs and keeping those that reach a rule no input
// before them reached. The candidates come from the program itself: the
// literals its conditions compare with, the values of its condition names,
// and their boundaries, in the input fields the lineage says each
// condition depends on.
package oracle

import (
	"fmt"
	"slices"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
	"github.com/vlumi/copal/bedrock/record"
	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/dataflow"
	"github.com/vlumi/copal/survey/internal/rules"
)

// program is a main program prepared for the search: its source, the
// dialect it is written in, its blueprint, its records laid out, its data
// flow, its decision tables, the
// inputs it reads, with the shape of each, by name, and the units whose
// rules a run of it can reach: itself, and the programs it calls.
type program struct {
	source    *cobol.Program
	dialect   cobol.Dialect
	blueprint *blueprint.Program
	records   *layout.ProgramRecords
	flow      dataflow.Flow
	tables    []rules.Table
	inputs    []input
	shapes    map[string]shape
	allTotals []total
	units     []Unit
}

// Unit is a program whose decision tables a run reaches: the program run,
// or one it calls, at any depth.
type Unit struct {
	Program   string
	Tables    []rules.Table
	source    *cobol.Program
	records   *layout.ProgramRecords
	blueprint *blueprint.Program
}

// input is a file the program reads, or its SYSIN, with the layout of its
// records.
type input struct {
	file   blueprint.File // empty for SYSIN
	root   *layout.Field
	layout *record.Layout
}

// sysin is the name an input that is the program's SYSIN goes by in a
// case.
const sysin = "SYSIN"

func (in input) name() string {
	if in.file.Name == "" {
		return sysin
	}
	return in.file.Name
}

func prepare(system *blueprint.System, bp *blueprint.Program, copybookDir string) (*program, error) {
	source, records, flow, err := analyze(bp, copybookDir, system.Dialect)
	if err != nil {
		return nil, err
	}
	p := &program{source: source, dialect: system.Dialect, blueprint: bp, records: records, flow: flow, tables: bp.Tables(), shapes: map[string]shape{}}
	p.units = []Unit{{Program: bp.Name, Tables: p.tables, source: source, records: records, blueprint: bp}}
	if err := p.addCallees(system, bp, copybookDir); err != nil {
		return nil, err
	}
	for _, f := range bp.Files {
		if !reads(f) {
			continue
		}
		if r := fileRecord(records, f.Name); r != nil {
			p.inputs = append(p.inputs, input{file: f, root: r, layout: record.New(r)})
		}
	}
	if r := accepted(source, records); r != nil {
		p.inputs = append(p.inputs, input{root: r, layout: record.New(r)})
	}
	for _, in := range p.inputs {
		p.shapes[in.name()] = p.shapeOf(in)
	}
	p.allTotals = p.totals()
	return p, nil
}

func analyze(bp *blueprint.Program, copybookDir string, dialect cobol.Dialect) (*cobol.Program, *layout.ProgramRecords, dataflow.Flow, error) {
	file, err := cobol.ReadProgram(bp.Source, []string{copybookDir})
	if err != nil {
		return nil, nil, dataflow.Flow{}, err
	}
	source, ok := file.Find(bp.Name)
	if !ok {
		return nil, nil, dataflow.Flow{}, fmt.Errorf("%s: no program %s", bp.Source, bp.Name)
	}
	records, err := layout.Program(source, dialect)
	if err != nil {
		return nil, nil, dataflow.Flow{}, err
	}
	return source, records, dataflow.Analyze(source, records), nil
}

// addCallees adds the programs a program calls, and those they call, each
// once, as units of the search.
func (p *program) addCallees(system *blueprint.System, bp *blueprint.Program, copybookDir string) error {
	for _, call := range bp.Calls {
		i := slices.IndexFunc(system.Programs, func(q blueprint.Program) bool { return q.Name == call.Program })
		if i < 0 || slices.ContainsFunc(p.units, func(u Unit) bool { return u.Program == call.Program }) {
			continue
		}
		callee := &system.Programs[i]
		source, records, _, err := analyze(callee, copybookDir, system.Dialect)
		if err != nil {
			return err
		}
		p.units = append(p.units, Unit{Program: callee.Name, Tables: callee.Tables(), source: source,
			records: records, blueprint: callee})
		if err := p.addCallees(system, callee, copybookDir); err != nil {
			return err
		}
	}
	return nil
}

// reads reports whether a program opens a file to read it.
func reads(f blueprint.File) bool {
	for _, mode := range f.Opened {
		if mode == "INPUT" || mode == "I-O" {
			return true
		}
	}
	return false
}

// writes reports whether a program opens a file to write it, from the
// start or at its end.
func writes(f blueprint.File) bool {
	for _, mode := range f.Opened {
		if mode == "OUTPUT" || mode == "EXTEND" {
			return true
		}
	}
	return false
}

// fileRecord is the first record of a file's description, which its other
// records, if any, redefine.
func fileRecord(records *layout.ProgramRecords, file string) *layout.Field {
	for _, f := range records.Files {
		if f.File == file && len(f.Records) > 0 {
			return f.Records[0]
		}
	}
	return nil
}

// find finds a field by name among a record's fields, at any level.
func find(f *layout.Field, name string) *layout.Field {
	if f.Name == name {
		return f
	}
	for _, child := range f.Children {
		if found := find(child, name); found != nil {
			return found
		}
	}
	return nil
}

// accepted is the record a program ACCEPTs from its SYSIN, if it does.
func accepted(p *cobol.Program, records *layout.ProgramRecords) *layout.Field {
	var name string
	for _, para := range p.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if s.Verb == "ACCEPT" && s.Device == "" && len(s.Writes) > 0 && name == "" {
				name = s.Writes[0].Name
			}
		})
	}
	for _, r := range records.All() {
		if r.Name == name {
			return r
		}
	}
	return nil
}
