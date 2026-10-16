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
