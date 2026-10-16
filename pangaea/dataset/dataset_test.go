package dataset

import (
	"os"
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

	if got, want := s.Names(), []string{"NEVER.OPENED"}; !slices.Equal(got, want) {
		t.Errorf("names %v, want %v", got, want)
	}
}

func TestUnloadRejectsUnprintableBytes(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.Path("A"), []byte("12\x0034"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.Catalog("A", 5); err != nil {
		t.Fatal(err)
	}
	if err := s.Unload("A", &strings.Builder{}); err == nil {
		t.Error("unloaded a NUL byte as text")
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
