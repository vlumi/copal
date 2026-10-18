// Package run compiles COBOL programs with GnuCOBOL and runs them one at a
// time, each with its DD names bound to files, as a job step would, and
// optionally with a trace of every paragraph and statement it executes.
package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/vlumi/copal/bedrock/cobol"
)

// flags make GnuCOBOL behave like an IBM compiler: IBM dialect, and signed
// display numbers with the EBCDIC overpunch (+5 as 'E', -5 as 'N') rather
// than the ASCII one. Warnings are errors.
var flags = []string{"-std=ibm", "-fsign=EBCDIC", "-Wall", "-Werror"}

// Build is the compiling of a directory of programs: where their sources
// and copybooks are, where the compiled programs go, and whether they
// carry the trace code that a run with a trace file turns on.
type Build struct {
	Sources, Copybooks, Out string
	Trace                   bool
}

// Programs compiles every .cbl file of the sources into the output
// directory, named after the file: an executable for a main program, and a
// module, which the runtime loads when a program calls it, for a program
// called with USING. It returns the program names.
func (b Build) Programs() ([]string, error) {
	sources, err := filepath.Glob(filepath.Join(b.Sources, "*.cbl"))
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(b.Out, 0o755); err != nil {
		return nil, err
	}
	options := flags
	if b.Trace {
		options = append(slices.Clone(flags), "-ftraceall")
	}
	var names []string
	for _, source := range sources {
		name := strings.TrimSuffix(filepath.Base(source), ".cbl")
		kind, err := kindFlag(source, b.Copybooks)
		if err != nil {
			return nil, err
		}
		args := append([]string{kind}, options...)
		args = append(args, "-I", b.Copybooks, "-o", filepath.Join(b.Out, name+moduleExtension(kind)), source)
		if output, err := exec.Command("cobc", args...).CombinedOutput(); err != nil {
			return nil, fmt.Errorf("cobc %s: %w\n%s", name, err, output)
		}
		names = append(names, name)
	}
	return names, nil
}

// kindFlag is -m for a program that is called with USING, that is, whose
// PROCEDURE DIVISION header has a USING phrase, and -x otherwise.
func kindFlag(source, copybookDir string) (string, error) {
	tokens, err := cobol.Read(source, []string{copybookDir})
	if err != nil {
		return "", err
	}
	for i := range len(tokens) - 2 {
		if tokens[i].Is("PROCEDURE") && tokens[i+1].Is("DIVISION") {
			if tokens[i+2].Is("USING") {
				return "-m", nil
			}
			break
		}
	}
	return "-x", nil
}

// moduleExtension is the file extension the runtime looks for when it
// loads a called program: a shared library's.
func moduleExtension(kind string) string {
	switch {
	case kind != "-m":
		return ""
	case runtime.GOOS == "darwin":
		return ".dylib"
	}
	return ".so"
}
