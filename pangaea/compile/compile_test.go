package compile

import (
	"errors"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
)

func TestPrograms(t *testing.T) {
	out := t.TempDir()
	names, err := Programs("testdata/programs", "testdata/copybooks", out)
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
