// Package dataset keeps Pangaea's datasets: files of fixed-length records
// without delimiters, as on a mainframe, named by their dataset name and
// listed with their record length in a catalog.
//
// Fixtures and golden files hold the same records as text, one per line,
// so that they read and diff well.
package dataset

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

const catalogFile = "CATALOG"

// Store is a directory of datasets and the catalog that lists them.
type Store struct {
	dir     string
	catalog map[string]int
}

// Open opens the store in dir, creating it if needed.
func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{dir: dir, catalog: map[string]int{}}
	return s, s.readCatalog()
}

// Path is the file that holds the dataset.
func (s *Store) Path(dsn string) string {
	return filepath.Join(s.dir, dsn)
}

// LRECL is the record length of a cataloged dataset.
func (s *Store) LRECL(dsn string) (int, bool) {
	lrecl, ok := s.catalog[dsn]
	return lrecl, ok
}

// Names lists the cataloged datasets in order.
func (s *Store) Names() []string {
	names := make([]string, 0, len(s.catalog))
	for name := range s.catalog {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// Catalog records a dataset that a program has written, after checking that
// its file holds whole records. A dataset the program never opened is
// cataloged empty, as an allocated dataset would be.
func (s *Store) Catalog(dsn string, lrecl int) error {
	info, err := os.Stat(s.Path(dsn))
	if errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(s.Path(dsn), nil, 0o644); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if info.Size()%int64(lrecl) != 0 {
		return fmt.Errorf("%s: %d bytes is not a whole number of %d-byte records", dsn, info.Size(), lrecl)
	}
	s.catalog[dsn] = lrecl
	return s.writeCatalog()
}

// Load creates a dataset from text, one record per line. Every line must
// have the same length, which becomes the record length.
func (s *Store) Load(dsn string, r io.Reader) error {
	var records bytes.Buffer
	lrecl := 0
	scanner := bufio.NewScanner(r)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := scanner.Bytes()
		if lrecl == 0 {
			lrecl = len(line)
		}
		if len(line) != lrecl || lrecl == 0 {
			return fmt.Errorf("%s line %d: %d characters, the first record has %d", dsn, lineNo, len(line), lrecl)
		}
		records.Write(line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if lrecl == 0 {
		return fmt.Errorf("%s: no records, so no record length", dsn)
	}
	if err := os.WriteFile(s.Path(dsn), records.Bytes(), 0o644); err != nil {
		return err
	}
	return s.Catalog(dsn, lrecl)
}

// Unload writes a cataloged dataset as text, one record per line.
func (s *Store) Unload(dsn string, w io.Writer) error {
	lrecl, ok := s.LRECL(dsn)
	if !ok {
		return fmt.Errorf("%s: not cataloged", dsn)
	}
	data, err := os.ReadFile(s.Path(dsn))
	if err != nil {
		return err
	}
	for record := range slices.Chunk(data, lrecl) {
		if bytes.IndexByte(record, '\n') >= 0 {
			return fmt.Errorf("%s: a record contains a newline and cannot be unloaded as text", dsn)
		}
		if _, err := fmt.Fprintf(w, "%s\n", record); err != nil {
			return err
		}
	}
	return nil
}

// LoadFile loads a text file as the dataset named by the file name without
// its .txt extension.
func (s *Store) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return s.Load(strings.TrimSuffix(filepath.Base(path), ".txt"), bytes.NewReader(data))
}

// UnloadAll writes every cataloged dataset to outDir as DSN.txt.
func (s *Store) UnloadAll(outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	for _, dsn := range s.Names() {
		if err := s.unloadFile(dsn, filepath.Join(outDir, dsn+".txt")); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) unloadFile(dsn, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := s.Unload(dsn, f); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func (s *Store) readCatalog() error {
	data, err := os.ReadFile(filepath.Join(s.dir, catalogFile))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	for line := range strings.Lines(string(data)) {
		dsn, lrecl, ok := strings.Cut(strings.TrimSpace(line), " ")
		n, err := strconv.Atoi(lrecl)
		if !ok || err != nil {
			return fmt.Errorf("catalog: bad entry %q", line)
		}
		s.catalog[dsn] = n
	}
	return nil
}

func (s *Store) writeCatalog() error {
	var b strings.Builder
	for _, dsn := range s.Names() {
		fmt.Fprintf(&b, "%s %d\n", dsn, s.catalog[dsn])
	}
	return os.WriteFile(filepath.Join(s.dir, catalogFile), []byte(b.String()), 0o644)
}
