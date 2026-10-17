// Package fileio finds how a program uses its files: the mode each OPEN
// opens it in, and the files a SORT reads with USING and writes with
// GIVING.
package fileio

import (
	"slices"

	"github.com/vlumi/copal/cobol"
)

// Modes of use. A SORT's USING files are read as input, its GIVING files
// written as output.
const (
	Input  = "INPUT"
	Output = "OUTPUT"
	Extend = "EXTEND"
	InOut  = "I-O"
)

// Use is a file used in a mode.
type Use struct {
	File string
	Mode string
	Pos  cobol.Pos
}

// Uses lists the ways the program uses its files, in the order they first
// appear, each file and mode once.
func Uses(p *cobol.Program) []Use {
	var uses []Use
	for _, para := range p.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			for _, f := range s.Files {
				u := Use{File: f.Name, Mode: f.Mode, Pos: s.Pos}
				if u.Mode != "" && !slices.ContainsFunc(uses, func(v Use) bool { return v.File == u.File && v.Mode == u.Mode }) {
					uses = append(uses, u)
				}
			}
		})
	}
	return uses
}

// Sort is a SORT statement: its work file, the files it reads with USING
// and those it writes with GIVING. A sort with an input or output
// procedure has no files on that side.
type Sort struct {
	Work   string
	Using  []string
	Giving []string
	Pos    cobol.Pos
}

// Sorts lists the program's SORT statements.
func Sorts(p *cobol.Program) []Sort {
	var sorts []Sort
	for _, para := range p.Paragraphs {
		cobol.Walk(para.Statements, func(s cobol.Statement) {
			if s.Verb != "SORT" {
				return
			}
			sort := Sort{Work: s.Files[0].Name, Pos: s.Pos}
			for _, f := range s.Files[1:] {
				if f.Mode == Input {
					sort.Using = append(sort.Using, f.Name)
				} else {
					sort.Giving = append(sort.Giving, f.Name)
				}
			}
			sorts = append(sorts, sort)
		})
	}
	return sorts
}
