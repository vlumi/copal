package blueprint

import (
	"slices"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/survey/internal/dataflow"
)

// callDepth is how many calls deep the origins are followed.
const callDepth = 4

// originsThroughCalls follows the origins of what the programs write
// through the programs they call. A field whose values come from an
// argument a program called writes keeps that origin, and gains those of
// what the called program gives back in it: what the called program
// reads, its literals and initial values, marked as its, and what is
// passed in to it, traced in the caller as the CALL passes it, so that a
// date converted by a called program is seen to come from the date it was
// converted from.
func originsThroughCalls(programs []Program) {
	t := &through{byName: map[string]*Program{}, given: map[parameterKey][]dataflow.Origin{}, passed: map[argumentKey][]dataflow.Origin{}}
	for i := range programs {
		t.byName[programs[i].Name] = &programs[i]
	}
	for i := range programs {
		p := &programs[i]
		for j := range p.Outputs {
			for k := range p.Outputs[j].Fields {
				fo := &p.Outputs[j].Fields[k]
				fo.origins = t.follow(*p, fo.origins, 0)
				fo.Origins = nil
				for _, o := range fo.origins {
					fo.Origins = append(fo.Origins, o.String())
				}
			}
		}
	}
}

// through follows origins through calls, remembering what each program
// gives back in each part of a parameter, and what each CALL passes in
// each argument, since many fields come through the same call.
type through struct {
	byName map[string]*Program
	given  map[parameterKey][]dataflow.Origin
	passed map[argumentKey][]dataflow.Origin
}

// parameterKey is a part of a program's parameter: the program, the
// parameter's position and the bytes, 0-based within it.
type parameterKey struct {
	program        string
	arg            int
	offset, length int
}

// argumentKey is an item a CALL of a program passes: the program and the
// CALL, and the item.
type argumentKey struct {
	program string
	call    cobol.Pos
	item    string
}

// follow adds to a program's origins those behind each call among them,
// through as many calls as callDepth allows.
func (t *through) follow(p Program, origins []dataflow.Origin, depth int) []dataflow.Origin {
	out := slices.Clone(origins)
	if depth == callDepth {
		return out
	}
	for _, o := range origins {
		called, ok := t.byName[o.Call]
		if o.Call == "" || o.Arg == 0 || !ok || o.Arg > len(called.Using) {
			continue
		}
		for _, call := range callsPassing(p, o) {
			for _, q := range t.givenBack(p, *called, o, call, depth) {
				if q.Parameter == 0 {
					if q.Through == "" {
						q.Through = called.Name
					}
					addOrigin(&out, q)
					continue
				}
				for _, r := range t.passedIn(p, *called, q, call, depth) {
					addOrigin(&out, r)
				}
			}
		}
	}
	return out
}

func addOrigin(origins *[]dataflow.Origin, o dataflow.Origin) {
	if !slices.Contains(*origins, o) {
		*origins = append(*origins, o)
	}
}

// givenBack is what the program called gives back in the bytes of its
// parameter that the item of a call origin was passed in, followed
// through the calls the called program makes.
func (t *through) givenBack(caller, called Program, o dataflow.Origin, call cobol.Statement, depth int) []dataflow.Origin {
	offset, _ := within(caller, o.Item, call.Arguments()[o.Arg-1].Name)
	item, _ := recordField(caller, o.Item)
	key := parameterKey{called.Name, o.Arg, offset, item.Length}
	if origins, ok := t.given[key]; ok {
		return origins
	}
	t.given[key] = nil
	origins := t.follow(called, parameterOrigins(called, called.Using[o.Arg-1], offset, item.Length), depth+1)
	t.given[key] = origins
	return origins
}

// passedIn traces, in the caller, what a CALL passes in the argument a
// parameter origin of the program called names, as the CALL passes it,
// followed through the calls the caller makes before it. What the CALL
// itself writes back into its arguments is not what it passes in.
func (t *through) passedIn(caller, called Program, q dataflow.Origin, call cobol.Statement, depth int) []dataflow.Origin {
	item := argumentOver(caller, called, q, call)
	key := argumentKey{caller.Name, call.Pos, item}
	if origins, ok := t.passed[key]; ok {
		return origins
	}
	t.passed[key] = nil
	before := slices.DeleteFunc(caller.flow.TraceAt(item, []cobol.Pos{call.Pos}).Origins, func(o dataflow.Origin) bool {
		return o.Call == called.Name && o.Arg > 0 && o.Arg <= len(call.Arguments()) && call.Arguments()[o.Arg-1].Name == o.Item
	})
	origins := t.follow(caller, before, depth+1)
	t.passed[key] = origins
	return origins
}

// callsPassing lists the CALLs of a program that pass the item of a call
// origin, or an item it is within, as the argument the origin names.
func callsPassing(p Program, o dataflow.Origin) []cobol.Statement {
	var calls []cobol.Statement
	for _, para := range p.source.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if s.Verb != "CALL" || s.Calls != o.Call || o.Arg > len(s.Arguments()) {
				return
			}
			if _, ok := within(p, o.Item, s.Arguments()[o.Arg-1].Name); ok {
				calls = append(calls, s)
			}
		})
	}
	return calls
}

// within gives where an item's bytes are within another of the program's
// items, 0-based, if the item is the other or a field of it.
func within(p Program, item, group string) (int, bool) {
	outer, ok := recordField(p, group)
	if !ok {
		return 0, false
	}
	offset, found := -1, false
	walk(outer, func(f Field) {
		if !found && f.Name == item {
			offset, found = f.Start-outer.Start, true
		}
	})
	return offset, found
}

// parameterOrigins are the origins, in a program called, of what it gives
// back in some bytes of a parameter: of each value a statement of it
// moves into those bytes, as the statement reads it.
func parameterOrigins(called Program, param string, offset, length int) []dataflow.Origin {
	var origins []dataflow.Origin
	for _, e := range called.flow.Lineage {
		at, ok := within(called, e.To, param)
		to, _ := recordField(called, e.To)
		if e.ToFile || e.ToCall || e.ToStream || !ok || at >= offset+length || at+to.Length <= offset {
			continue
		}
		switch {
		case e.FromLiteral:
			addOrigin(&origins, dataflow.Origin{Literal: e.From})
		case e.FromFile:
			addOrigin(&origins, dataflow.Origin{File: e.From, Item: e.To})
		case e.FromCall:
			addOrigin(&origins, dataflow.Origin{Call: e.From, Arg: e.Arg, Item: e.To})
		case e.FromStream:
			addOrigin(&origins, dataflow.Origin{Stream: e.From, Item: e.To})
		default:
			for _, q := range called.flow.TraceAt(e.From, []cobol.Pos{e.Pos}).Origins {
				addOrigin(&origins, q)
			}
		}
	}
	return origins
}

// argumentOver is the caller's item over the bytes a parameter origin of
// the program called names, in the argument the CALL passes for the
// parameter: the field at those bytes, or the whole argument.
func argumentOver(caller, called Program, q dataflow.Origin, call cobol.Statement) string {
	arg := call.Arguments()[q.Parameter-1].Name
	offset, ok := within(called, q.Item, called.Using[q.Parameter-1])
	item, found := recordField(called, q.Item)
	argField, known := recordField(caller, arg)
	if !ok || !found || !known {
		return arg
	}
	if name, ok := fieldOver(argField, argField.Start+offset, item.Length); ok {
		return name
	}
	return arg
}
