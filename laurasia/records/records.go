// Package records loads the record layouts of Laurasia Card Services,
// computed from its copybooks in its dialect, for the Go code around the
// processor to read and write its files by field.
package records

import (
	"path/filepath"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/cobolith/record"
)

// Layouts are the processor's record layouts: a card's account, the
// switch's log entry, the five records of the clearing file, the card
// system's feed record, a dispute's event and case, a fee of the
// schedule, a presentment on its way through the night, and a posting to
// a card.
type Layouts struct {
	Account, Auth                                                   *record.Layout
	FileHeader, BatchHeader, Presentment, BatchTrailer, FileTrailer *record.Layout
	Feed, Event, Case, Fee, Pres, Posting                           *record.Layout
}

// Load reads the layouts from the processor's copybooks.
func Load(copybookDir string) (*Layouts, error) {
	var l Layouts
	for _, c := range []struct {
		target   **record.Layout
		copybook string
		name     string
	}{
		{&l.Account, "cardacct", "CARD-ACCOUNT"},
		{&l.Auth, "authlog", "AUTH-ENTRY"},
		{&l.FileHeader, "clrrec", "CLR-FILE-HEADER"},
		{&l.BatchHeader, "clrrec", "CLR-BATCH-HEADER"},
		{&l.Presentment, "clrrec", "CLR-PRESENTMENT"},
		{&l.BatchTrailer, "clrrec", "CLR-BATCH-TRAILER"},
		{&l.FileTrailer, "clrrec", "CLR-FILE-TRAILER"},
		{&l.Feed, "payfeed", "FEED-RECORD"},
		{&l.Event, "dispevt", "DISPUTE-EVENT"},
		{&l.Case, "dispcase", "DISPUTE-CASE"},
		{&l.Fee, "feedef", "FEE-DEFINITION"},
		{&l.Pres, "presrec", "PRES-RECORD"},
		{&l.Posting, "postrec", "POSTING"},
	} {
		layout, err := record.FromCopybook(filepath.Join(copybookDir, c.copybook+".cpy"), []string{copybookDir}, c.name, cobol.MicroFocus)
		if err != nil {
			return nil, err
		}
		*c.target = layout
	}
	return &l, nil
}
