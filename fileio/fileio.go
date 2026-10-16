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
	files := map[string]bool{}
	for _, f := range p.Files {
		files[f.Name] = true
	}
	var uses []Use
	add := func(u Use) {
		if !slices.ContainsFunc(uses, func(v Use) bool { return v.File == u.File && v.Mode == u.Mode }) {
			uses = append(uses, u)
		}
	}
	for _, para := range p.Paragraphs {
		body := para.Body
		for i := 0; i < len(body); i++ {
			var modes map[string]string
			switch {
			case body[i].Is("OPEN"):
				modes = map[string]string{"INPUT": Input, "OUTPUT": Output, "EXTEND": Extend, "I-O": InOut}
			case body[i].Is("SORT"):
				modes = map[string]string{"USING": Input, "GIVING": Output}
			default:
				continue
			}
			mode := ""
			for j := i + 1; j < len(body) && body[j].Kind == cobol.Word; j++ {
				t := body[j]
				if m, ok := modes[t.Text]; ok {
					mode = m
				} else if mode != "" && files[t.Text] {
					add(Use{File: t.Text, Mode: mode, Pos: t.Pos})
				} else if mode != "" {
					break
				}
			}
		}
	}
	return uses
}
