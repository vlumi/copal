package fixture

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/records"
)

// accruals is the bank's codec, whose accrual records hold packed decimal
// fields.
func accruals(t *testing.T) *Codec {
	t.Helper()
	l, err := records.Load("../copybooks")
	if err != nil {
		t.Fatal(err)
	}
	return NewCodec(l)
}

const accrual = `AC-ACCT-NO 10006
AC-THROUGH-DATE 20261030
AC-ACCRUED -0.048081
AC-DUE -1.53
AC-DUE-MONTH 202610
`

func TestFieldsRoundTrip(t *testing.T) {
	codec := accruals(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "ACCT.ACCRUAL.G0002V00.fields")
	if err := os.WriteFile(path, []byte(accrual+"\n"+strings.Replace(accrual, "10006", "10007", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	store, err := dataset.Open(filepath.Join(dir, "datasets"))
	if err != nil {
		t.Fatal(err)
	}
	if err := codec.LoadFile(store, path); err != nil {
		t.Fatal(err)
	}
	if records, err := store.Read("ACCT.ACCRUAL.G0002V00"); err != nil || len(records) != 2 {
		t.Fatalf("%d records, %v", len(records), err)
	}
	out := filepath.Join(dir, "out")
	if err := codec.UnloadAll(store, out); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(out, "ACCT.ACCRUAL.G0002V00.fields"))
	if err != nil {
		t.Fatal(err)
	}
	if want := accrual + "\n" + strings.Replace(accrual, "10006", "10007", 1); string(data) != want {
		t.Errorf("unloaded\n%s\nwant\n%s", data, want)
	}
}

func TestFieldsErrors(t *testing.T) {
	codec := accruals(t)
	l := codec.layout("ACCT.ACCRUAL.G0001V00")
	for _, lines := range [][]string{
		{"AC-DUE 1.5"},
		{"AC-NOTHING 1"},
		{"AC-ACCT-NO ten"},
	} {
		if _, err := decodeRecord(l, lines); err == nil {
			t.Errorf("%q decoded without error", lines)
		}
	}
	r, err := decodeRecord(l, strings.Split(strings.TrimSpace(accrual), "\n"))
	if err != nil {
		t.Fatal(err)
	}
	copy(r[l.Length-4:], "JUNK")
	if _, err := encodeFields(l, r); err == nil {
		t.Error("a record with something in its FILLER was written field by field")
	}
}
