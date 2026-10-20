// Package dataset keeps a system's datasets: files of fixed-length records
// without delimiters, as on a mainframe, named by their dataset name and
// listed with their record length in a catalog.
//
// An indexed dataset, a key-sequenced dataset on a mainframe, is the file
// GnuCOBOL's indexed file handler keeps, a Berkeley DB database. Only
// programs read and write it, so it has no text form; the catalog lists it
// as KS.
//
// Fixtures and golden files hold the same records as text, one per line,
// so that they read and diff well. A record with bytes that are not
// printable, such as packed decimal fields, is written in hex the way
// COBOL writes a hex literal: X'F0F1C3'.
package dataset

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const catalogFile = "CATALOG"

// Store is a directory of datasets and the catalog that lists them.
type Store struct {
	dir     string
	catalog map[string]entry
}

// entry is a dataset's record length and whether it is indexed.
type entry struct {
	lrecl   int
	indexed bool
}

// Open opens the store in dir, creating it if needed.
func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{dir: dir, catalog: map[string]entry{}}
	return s, s.readCatalog()
}

// Path is the file that holds the dataset.
func (s *Store) Path(dsn string) string {
	return filepath.Join(s.dir, dsn)
}

// LRECL is the record length of a cataloged dataset.
func (s *Store) LRECL(dsn string) (int, bool) {
	e, ok := s.catalog[dsn]
	return e.lrecl, ok
}

// Indexed tells whether a cataloged dataset is indexed.
func (s *Store) Indexed(dsn string) bool {
	return s.catalog[dsn].indexed
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

// GenerationName is the name of generation n of a generation data group,
// as on a mainframe: ACCT.MASTER.G0003V00.
func GenerationName(base string, n int) string {
	return fmt.Sprintf("%s.G%04dV00", base, n)
}

// Generations lists the numbers of a generation data group's cataloged
// generations, oldest first.
func (s *Store) Generations(base string) []int {
	var generations []int
	for dsn := range s.catalog {
		suffix, ok := strings.CutPrefix(dsn, base+".G")
		if !ok || !strings.HasSuffix(suffix, "V00") {
			continue
		}
		if n, err := strconv.Atoi(strings.TrimSuffix(suffix, "V00")); err == nil {
			generations = append(generations, n)
		}
	}
	slices.Sort(generations)
	return generations
}

// LatestGeneration is the number of the newest cataloged generation of a
// generation data group, or 0 if it has none.
func (s *Store) LatestGeneration(base string) int {
	generations := s.Generations(base)
	if len(generations) == 0 {
		return 0
	}
	return generations[len(generations)-1]
}

// Delete removes a dataset and its catalog entry.
func (s *Store) Delete(dsn string) error {
	if err := os.Remove(s.Path(dsn)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	delete(s.catalog, dsn)
	return s.writeCatalog()
}

// Catalog records a dataset that a program has written, after checking that
// its file holds whole records. A dataset the program never opened is
// cataloged empty, as an allocated dataset would be.
func (s *Store) Catalog(dsn string, lrecl int) error {
	if lrecl < 1 {
		return fmt.Errorf("%s: a record length of %d", dsn, lrecl)
	}
	info, err := os.Stat(s.Path(dsn))
	switch {
	case errors.Is(err, os.ErrNotExist):
		if err := os.WriteFile(s.Path(dsn), nil, 0o644); err != nil {
			return err
		}
	case err != nil:
		return err
	case info.Size()%int64(lrecl) != 0:
		return fmt.Errorf("%s: %d bytes is not a whole number of %d-byte records", dsn, info.Size(), lrecl)
	}
	s.catalog[dsn] = entry{lrecl: lrecl}
	return s.writeCatalog()
}

// CatalogIndexed records an indexed dataset that a program has written. A
// program that never opened it leaves no file, as a mainframe leaves a
// cluster defined but never loaded, which cannot be opened for input.
func (s *Store) CatalogIndexed(dsn string, lrecl int) error {
	if lrecl < 1 {
		return fmt.Errorf("%s: a record length of %d", dsn, lrecl)
	}
	s.catalog[dsn] = entry{lrecl: lrecl, indexed: true}
	return s.writeCatalog()
}

// Write creates a dataset from records of the given length and catalogs it.
func (s *Store) Write(dsn string, lrecl int, records [][]byte) error {
	var data bytes.Buffer
	for _, r := range records {
		if len(r) != lrecl {
			return fmt.Errorf("%s: a record of %d bytes, want %d", dsn, len(r), lrecl)
		}
		data.Write(r)
	}
	if err := os.WriteFile(s.Path(dsn), data.Bytes(), 0o644); err != nil {
		return err
	}
	return s.Catalog(dsn, lrecl)
}

// Read returns the records of a cataloged dataset.
func (s *Store) Read(dsn string) ([][]byte, error) {
	lrecl, ok := s.LRECL(dsn)
	switch {
	case !ok:
		return nil, fmt.Errorf("%s: not cataloged", dsn)
	case s.Indexed(dsn):
		return nil, fmt.Errorf("%s: indexed, read only by programs", dsn)
	}
	data, err := os.ReadFile(s.Path(dsn))
	if err != nil {
		return nil, err
	}
	return slices.Collect(slices.Chunk(data, lrecl)), nil
}

// Load creates a dataset from text, one record per line, each line the
// record itself or the record in hex. Every record must have the same
// length, which becomes the record length.
func (s *Store) Load(dsn string, r io.Reader) error {
	var records bytes.Buffer
	lrecl := 0
	scanner := bufio.NewScanner(r)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		record, err := fromText(scanner.Text())
		if err != nil {
			return fmt.Errorf("%s line %d: %w", dsn, lineNo, err)
		}
		if lrecl == 0 {
			lrecl = len(record)
		}
		if len(record) != lrecl || lrecl == 0 {
			return fmt.Errorf("%s line %d: a record of %d bytes, the first has %d", dsn, lineNo, len(record), lrecl)
		}
		records.Write(record)
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
	records, err := s.Read(dsn)
	if err != nil {
		return err
	}
	for _, record := range records {
		if _, err := fmt.Fprintln(w, toText(record)); err != nil {
			return err
		}
	}
	return nil
}

// hexRecord is a record written in hex.
var hexRecord = regexp.MustCompile(`^X'((?:[0-9A-F]{2})+)'$`)

// toText writes a record as itself if it is printable, and in hex if it
// is not, or if it would read as hex.
func toText(record []byte) string {
	if bytes.IndexFunc(record, notPrintable) < 0 && !hexRecord.Match(record) {
		return string(record)
	}
	return "X'" + strings.ToUpper(hex.EncodeToString(record)) + "'"
}

func fromText(line string) ([]byte, error) {
	if m := hexRecord.FindStringSubmatch(line); m != nil {
		return hex.DecodeString(m[1])
	}
	if i := strings.IndexFunc(line, notPrintable); i >= 0 {
		return nil, fmt.Errorf("character %d is not printable; write the record in hex", i+1)
	}
	return []byte(line), nil
}

func notPrintable(r rune) bool {
	return r < ' ' || r > '~'
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

// UnloadAll writes every cataloged dataset but the indexed ones to outDir
// as DSN.txt.
func (s *Store) UnloadAll(outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	for _, dsn := range s.Names() {
		if s.Indexed(dsn) {
			continue
		}
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
		fields := strings.Fields(line)
		if len(fields) < 2 || len(fields) > 3 || len(fields) == 3 && fields[2] != "KS" {
			return fmt.Errorf("catalog: bad entry %q", line)
		}
		n, err := strconv.Atoi(fields[1])
		if err != nil || n < 1 {
			return fmt.Errorf("catalog: bad entry %q", line)
		}
		s.catalog[fields[0]] = entry{lrecl: n, indexed: len(fields) == 3}
	}
	return nil
}

func (s *Store) writeCatalog() error {
	var b strings.Builder
	for _, dsn := range s.Names() {
		e := s.catalog[dsn]
		fmt.Fprintf(&b, "%s %d", dsn, e.lrecl)
		if e.indexed {
			b.WriteString(" KS")
		}
		b.WriteString("\n")
	}
	return os.WriteFile(filepath.Join(s.dir, catalogFile), []byte(b.String()), 0o644)
}
