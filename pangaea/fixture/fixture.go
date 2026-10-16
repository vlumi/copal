// Package fixture loads Pangaea's fixtures and writes its golden datasets
// as text that reads and diffs well. Most datasets are text already, one
// record per line. A dataset whose records hold packed decimal fields is
// written field by field instead, decoded through the record layout Copal
// finds for it: NAME.fields, a line per field with its value, and a blank
// line between records.
package fixture

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/vlumi/copal/blueprint"
	"github.com/vlumi/copal/layout"
	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/record"
)

// FieldsExtension marks a dataset written field by field.
const FieldsExtension = ".fields"

// Codec writes and reads datasets as text, by the layouts of those with
// packed decimal fields, keyed by dataset name without generation.
type Codec struct {
	layouts map[string]*record.Layout
}

// NewCodec finds the datasets with packed decimal fields among those the
// system's jobs name, and their layouts from the system's copybooks, which
// copy others from dir.
func NewCodec(system *blueprint.System, dir string) (*Codec, error) {
	c := &Codec{layouts: map[string]*record.Layout{}}
	for base, name := range system.DatasetCopybooks() {
		for _, cb := range system.Copybooks {
			if cb.Name != name || len(cb.Records) == 0 {
				continue
			}
			l, err := record.FromCopybook(cb.Source, []string{dir}, cb.Records[0].Name)
			if err != nil {
				return nil, err
			}
			if packed(l) {
				c.layouts[base] = l
			}
		}
	}
	return c, nil
}

func packed(l *record.Layout) bool {
	for _, f := range l.Fields() {
		if f.Usage == "COMP-3" {
			return true
		}
	}
	return false
}

var generation = regexp.MustCompile(`\.G\d{4}V00$`)

func (c *Codec) layout(dsn string) *record.Layout {
	return c.layouts[generation.ReplaceAllString(dsn, "")]
}

// LoadFile loads a fixture file as the dataset it names: DSN.txt as text,
// DSN.fields field by field.
func (c *Codec) LoadFile(store *dataset.Store, path string) error {
	dsn, ok := strings.CutSuffix(filepath.Base(path), FieldsExtension)
	if !ok {
		return store.LoadFile(path)
	}
	l := c.layout(dsn)
	if l == nil {
		return fmt.Errorf("%s: no record layout for %s", path, dsn)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	records, err := decodeRecords(l, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return store.Write(dsn, l.Length, records)
}

// UnloadAll writes every cataloged dataset but the indexed ones, which
// have no text form, to dir: as DSN.fields if its records hold packed
// decimal fields, as DSN.txt otherwise.
func (c *Codec) UnloadAll(store *dataset.Store, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, dsn := range store.Names() {
		if store.Indexed(dsn) {
			continue
		}
		l := c.layout(dsn)
		var b bytes.Buffer
		var err error
		name := dsn + ".txt"
		if l == nil {
			err = store.Unload(dsn, &b)
		} else {
			name = dsn + FieldsExtension
			err = c.unloadFields(store, dsn, l, &b)
		}
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, name), b.Bytes(), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func (c *Codec) unloadFields(store *dataset.Store, dsn string, l *record.Layout, b *bytes.Buffer) error {
	records, err := store.Read(dsn)
	if err != nil {
		return err
	}
	for i, r := range records {
		if i > 0 {
			b.WriteString("\n")
		}
		lines, err := encodeFields(l, r)
		if err != nil {
			return fmt.Errorf("%s record %d: %w", dsn, i+1, err)
		}
		for _, line := range lines {
			b.WriteString(line + "\n")
		}
	}
	return nil
}

// encodeFields writes a record's fields, FILLER left out, and checks that
// they make the record again, byte for byte.
func encodeFields(l *record.Layout, r []byte) ([]string, error) {
	var lines []string
	for _, f := range l.Fields() {
		if f.Name == "FILLER" {
			continue
		}
		value, err := fieldText(l, r, f)
		if err != nil {
			return nil, err
		}
		lines = append(lines, f.Name+" "+value)
	}
	again, err := decodeRecord(l, lines)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(again, r) {
		return nil, fmt.Errorf("%q does not come back from its fields", r)
	}
	return lines, nil
}

func fieldText(l *record.Layout, r []byte, f *layout.Field) (string, error) {
	if f.Category != layout.Numeric {
		return l.Text(r, f.Name)
	}
	n, err := l.Number(r, f.Name)
	if err != nil {
		return "", err
	}
	if f.Scale == 0 {
		return strconv.FormatInt(n, 10), nil
	}
	sign, digits := "", strconv.FormatInt(n, 10)
	if n < 0 {
		sign, digits = "-", digits[1:]
	}
	digits = fmt.Sprintf("%0*s", f.Scale+1, digits)
	return sign + digits[:len(digits)-f.Scale] + "." + digits[len(digits)-f.Scale:], nil
}

// decodeRecords reads records written field by field.
func decodeRecords(l *record.Layout, r io.Reader) ([][]byte, error) {
	var records [][]byte
	var lines []string
	flush := func() error {
		if len(lines) == 0 {
			return nil
		}
		record, err := decodeRecord(l, lines)
		records, lines = append(records, record), nil
		return err
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if scanner.Text() == "" {
			if err := flush(); err != nil {
				return nil, err
			}
			continue
		}
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, flush()
}

// decodeRecord makes a record from its fields' lines.
func decodeRecord(l *record.Layout, lines []string) ([]byte, error) {
	r := l.Blank()
	for _, line := range lines {
		name, value, _ := strings.Cut(line, " ")
		var f *layout.Field
		for _, field := range l.Fields() {
			if field.Name == name {
				f = field
			}
		}
		if f == nil {
			return nil, fmt.Errorf("no field %s", name)
		}
		if f.Category != layout.Numeric {
			if err := l.SetText(r, name, value); err != nil {
				return nil, err
			}
			continue
		}
		whole, fraction, _ := strings.Cut(value, ".")
		units, err := strconv.ParseInt(whole+fraction, 10, 64)
		if err != nil || len(fraction) != f.Scale {
			return nil, fmt.Errorf("%s: %q is not a number with %d decimals", name, value, f.Scale)
		}
		if err := l.SetNumber(r, name, units); err != nil {
			return nil, err
		}
	}
	return r, nil
}
