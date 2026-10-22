package run

import (
	"errors"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/vlumi/cobolith/cobol"
)

func TestPrograms(t *testing.T) {
	out := t.TempDir()
	names, err := Build{Dialect: cobol.IBM, Sources: "testdata/programs", Copybooks: "testdata/copybooks", Out: out}.Programs()
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"HELLO"}; !slices.Equal(names, want) {
		t.Fatalf("compiled %v, want %v", names, want)
	}

	output, err := exec.Command(filepath.Join(out, "HELLO")).Output()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 4 {
		t.Errorf("exit %v, want return code 4", err)
	}
	if got, want := string(output), "HELLO 0012N\n"; got != want {
		t.Errorf("output %q, want %q with the EBCDIC overpunch sign", got, want)
	}
}

// A program called with USING is compiled as a module, which the runtime
// loads from COB_LIBRARY_PATH when another program calls it.
func TestCalledPrograms(t *testing.T) {
	out := t.TempDir()
	if _, err := (Build{Dialect: cobol.IBM, Sources: "testdata/called", Copybooks: "testdata/copybooks", Out: out}).Programs(); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(out, "CALLER"))
	cmd.Env = append(cmd.Environ(), "COB_LIBRARY_PATH="+out)
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(output), "HELLO WORLD\n"; got != want {
		t.Errorf("output %q, want %q", got, want)
	}
}

// A Fujitsu NetCOBOL program is in Shift_JIS, which its header is read in
// and GnuCOBOL compiles as it is, so its Japanese literals come out whole.
func TestFujitsuPrograms(t *testing.T) {
	out := t.TempDir()
	if _, err := (Build{Dialect: cobol.Fujitsu, Sources: "testdata/fujitsu", Copybooks: "testdata/copybooks", Out: out}).Programs(); err != nil {
		t.Fatal(err)
	}
	output, err := exec.Command(filepath.Join(out, "KONNICHI")).Output()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 4 {
		t.Errorf("exit %v, want return code 4", err)
	}
	// こんにちは世界 1,234,567円, in Shift_JIS.
	if got, want := string(output), "\x82\xb1\x82\xf1\x82\xc9\x82\xbf\x82\xcd\x90\xa2\x8aE 1,234,567\x89~\n"; got != want {
		t.Errorf("output %q, want %q", got, want)
	}
}
