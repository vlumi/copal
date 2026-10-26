package runlog

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/vlumi/copal/bedrock/runner"
)

// TestAppendAndRead appends runs in two goes and reads them back, a run
// for each date and job, an abend and a step with no program among them.
func TestAppendAndRead(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	if runs, err := Read(path); err != nil || runs != nil {
		t.Fatalf("a record not there: %v, %v", runs, err)
	}
	want := []Run{
		{Date: "2026-09-01", Job: "START", Steps: []runner.Result{{Step: "LOAD", Program: "ACCTLOAD"}}},
		{Date: "2026-09-02", Job: "DAILY", Steps: []runner.Result{{Step: "VALID", Program: "TXVALID", RC: 4}, {Step: "POST", Program: "TXPOST", Abended: true}}},
		{Date: "2026-09-02", Job: "nightly", Steps: []runner.Result{{Step: "SORT", Bypassed: true}}},
	}
	if err := Append(path, want[0]); err != nil {
		t.Fatal(err)
	}
	if err := Append(path, want[1:]...); err != nil {
		t.Fatal(err)
	}
	got, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("read %+v, want %+v", got, want)
	}
	if err := os.WriteFile(path, []byte("2026-09-02 DAILY VALID\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(path); err == nil {
		t.Error("a line of three fields read as a step's run")
	}
}
