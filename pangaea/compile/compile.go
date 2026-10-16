// Package compile builds Pangaea's COBOL programs with GnuCOBOL.
package compile

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/vlumi/copal/cobol"
)

// flags make GnuCOBOL behave like an IBM compiler: IBM dialect, and signed
// display numbers with the EBCDIC overpunch (+5 as 'E', -5 as 'N') rather
// than the ASCII one. Warnings are errors.
var flags = []string{"-std=ibm", "-fsign=EBCDIC", "-Wall", "-Werror"}

// Programs compiles every .cbl file in srcDir into outDir, named after the
// file: an executable for a main program, and a module, which the runtime
// loads when a program calls it, for a program called with USING. It
// returns the program names.
func Programs(srcDir, copybookDir, outDir string) ([]string, error) {
	sources, err := filepath.Glob(filepath.Join(srcDir, "*.cbl"))
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	var names []string
	for _, source := range sources {
		name := strings.TrimSuffix(filepath.Base(source), ".cbl")
		kind, err := kindFlag(source, copybookDir)
		if err != nil {
			return nil, err
		}
		args := append([]string{kind}, flags...)
		args = append(args, "-I", copybookDir, "-o", filepath.Join(outDir, name+moduleExtension(kind)), source)
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
