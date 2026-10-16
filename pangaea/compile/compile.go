// Package compile builds Pangaea's COBOL programs with GnuCOBOL.
package compile

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// flags make GnuCOBOL behave like an IBM compiler: IBM dialect, and signed
// display numbers with the EBCDIC overpunch (+5 as 'E', -5 as 'N') rather
// than the ASCII one. Warnings are errors.
var flags = []string{"-x", "-std=ibm", "-fsign=EBCDIC", "-Wall", "-Werror"}

// Programs compiles every .cbl file in srcDir into an executable in outDir,
// named after the file. It returns the program names.
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
		args := append(flags, "-I", copybookDir, "-o", filepath.Join(outDir, name), source)
		if output, err := exec.Command("cobc", args...).CombinedOutput(); err != nil {
			return nil, fmt.Errorf("cobc %s: %w\n%s", name, err, output)
		}
		names = append(names, name)
	}
	return names, nil
}
