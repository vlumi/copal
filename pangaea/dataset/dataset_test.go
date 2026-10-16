package dataset

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestLoadAndUnload(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	text := "0001ALICE\n0002BOB  \n"
	if err := s.Load("ACCT.MASTER", strings.NewReader(text)); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(s.Path("ACCT.MASTER"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "0001ALICE0002BOB  " {
		t.Errorf("stored %q, want the records without newlines", data)
	}

	reopened, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if lrecl, ok := reopened.LRECL("ACCT.MASTER"); !ok || lrecl != 9 {
		t.Errorf("cataloged LRECL %d, %v, want 9", lrecl, ok)
	}
	var unloaded strings.Builder
	if err := reopened.Unload("ACCT.MASTER", &unloaded); err != nil {
		t.Fatal(err)
	}
	if unloaded.String() != text {
		t.Errorf("unloaded %q, want %q", unloaded.String(), text)
	}
}

func TestLoadRejectsRaggedRecords(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Load("A", strings.NewReader("0001ALICE\n0002BOB\n")); err == nil {
		t.Error("loaded records of different lengths")
	}
	if err := s.Load("A", strings.NewReader("")); err == nil {
		t.Error("loaded an empty dataset")
	}
}

func TestCatalog(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Catalog("NEVER.OPENED", 80); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(s.Path("NEVER.OPENED")); err != nil || info.Size() != 0 {
		t.Errorf("a dataset never written should be cataloged empty, got %v, %v", info, err)
	}

	if err := os.WriteFile(s.Path("SHORT"), []byte("12345"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.Catalog("SHORT", 4); err == nil {
		t.Error("cataloged a partial record")
	}
	if err := s.Catalog("SHORT", 0); err == nil {
		t.Error("cataloged records of no length")
	}

	if got, want := s.Names(), []string{"NEVER.OPENED"}; !slices.Equal(got, want) {
		t.Errorf("names %v, want %v", got, want)
	}
}

// An indexed dataset is cataloged as written, survives reopening, and
// has no text form.
func TestCatalogIndexed(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.Path("ACCT.ONLINE"), []byte("not whole records"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.CatalogIndexed("ACCT.ONLINE", 300); err != nil {
		t.Fatal(err)
	}
	if err := s.Load("ACCT.MASTER", strings.NewReader("0001\n")); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if lrecl, _ := reopened.LRECL("ACCT.ONLINE"); lrecl != 300 || !reopened.Indexed("ACCT.ONLINE") || reopened.Indexed("ACCT.MASTER") {
		t.Errorf("ACCT.ONLINE LRECL %d, indexed %v; ACCT.MASTER indexed %v", lrecl, reopened.Indexed("ACCT.ONLINE"), reopened.Indexed("ACCT.MASTER"))
	}
	if _, err := reopened.Read("ACCT.ONLINE"); err == nil {
		t.Error("read an indexed dataset as records")
	}
	out := t.TempDir()
	if err := reopened.UnloadAll(out); err != nil {
		t.Fatal(err)
	}
	if entries, _ := os.ReadDir(out); len(entries) != 1 || entries[0].Name() != "ACCT.MASTER.txt" {
		t.Errorf("unloaded %v, want only ACCT.MASTER.txt", entries)
	}
}

func TestGenerations(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if n := s.LatestGeneration("ACCT.MASTER"); n != 0 {
		t.Errorf("latest of an empty group %d, want 0", n)
	}
	for _, dsn := range []string{GenerationName("ACCT.MASTER", 1), GenerationName("ACCT.MASTER", 2), "ACCT.MASTER.NEW"} {
		if err := s.Catalog(dsn, 100); err != nil {
			t.Fatal(err)
		}
	}
	if n := s.LatestGeneration("ACCT.MASTER"); n != 2 {
		t.Errorf("latest generation %d, want 2", n)
	}
	if err := s.Delete(GenerationName("ACCT.MASTER", 2)); err != nil {
		t.Fatal(err)
	}
	if n := s.LatestGeneration("ACCT.MASTER"); n != 1 {
		t.Errorf("latest generation after deleting G0002 %d, want 1", n)
	}
	if _, err := os.Stat(s.Path("ACCT.MASTER.G0002V00")); !os.IsNotExist(err) {
		t.Error("deleted dataset's file is still there")
	}
}

func TestWriteAndRead(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	records := [][]byte{[]byte("AB"), []byte("CD")}
	if err := s.Write("A", 2, records); err != nil {
		t.Fatal(err)
	}
	got, err := s.Read("A")
	if err != nil || len(got) != 2 || string(got[1]) != "CD" {
		t.Errorf("read %q, %v", got, err)
	}
	if err := s.Write("B", 3, records); err == nil {
		t.Error("wrote 2-byte records as 3-byte ones")
	}
}

// A record that is not printable is written in hex, and so is one that
// would otherwise read as hex.
func TestRecordsInHex(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	records := [][]byte{{'A', 0x12, 0x34, 0x5C, 'Z'}, []byte("X'41'"), []byte("HELLO")}
	if err := s.Write("PACKED", 5, records); err != nil {
		t.Fatal(err)
	}
	var text strings.Builder
	if err := s.Unload("PACKED", &text); err != nil {
		t.Fatal(err)
	}
	if want := "X'4112345C5A'\nX'5827343127'\nHELLO\n"; text.String() != want {
		t.Errorf("unloaded %q, want %q", text.String(), want)
	}
	if err := s.Load("RELOADED", strings.NewReader(text.String())); err != nil {
		t.Fatal(err)
	}
	if reloaded, err := s.Read("RELOADED"); err != nil || !slices.EqualFunc(reloaded, records, bytes.Equal) {
		t.Errorf("reloaded %q, %v, want %q", reloaded, err, records)
	}
	if err := s.Load("RAW", strings.NewReader("A\x01C\n")); err == nil {
		t.Error("loaded a line that is not printable")
	}
}

func TestCatalogRejectsNoRecordLength(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "CATALOG"), []byte("A 0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(dir); err == nil {
		t.Error("opened a catalog with a record length of 0")
	}
}
