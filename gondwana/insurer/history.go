package insurer

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/gondwana/generate"
)

// History is what the insurer has done since it opened: its seed, the
// next business day it will run, and the days it ran, oldest first, read
// back from the batch files' logs and the files the days wrote.
type History struct {
	Seed uint64
	Next generate.Date
	Days []Day
}

// History reads the insurer's history.
func (in Insurer) History() (History, error) {
	g, err := in.open()
	if err != nil {
		return History{}, err
	}
	h := History{Seed: g.Seed}
	if h.Next, err = in.nextDate(); err != nil {
		return History{}, err
	}
	logs, err := filepath.Glob(in.data("LOG", "DAILY.*.LOG"))
	if err != nil {
		return History{}, err
	}
	slices.Sort(logs)
	for _, log := range logs {
		n, err := strconv.ParseInt(stampOf(log), 10, 64)
		if err != nil {
			return History{}, err
		}
		date, err := generate.DateFromNumber(n)
		if err != nil {
			return History{}, err
		}
		day, err := in.day(g, date)
		if err != nil {
			return History{}, err
		}
		h.Days = append(h.Days, day)
	}
	return h, nil
}

// daysBatches are the batch files a day may run, in the order it runs
// them.
var daysBatches = []string{"RESULTS", "DAILY", "BILLING", "CLOSE", "ANNUAL"}

// day reads back a day the insurer ran: the batch files whose logs it
// wrote, and what it amounted to from its files.
func (in Insurer) day(g *generate.Generator, date generate.Date) (Day, error) {
	day := Day{Date: date}
	for _, name := range daysBatches {
		if _, err := os.Stat(in.data("LOG", name+"."+date.Stamp()+".LOG")); errors.Is(err, os.ErrNotExist) {
			continue
		}
		batches, err := in.batches(name, date)
		if err != nil {
			return Day{}, err
		}
		day.Batches = append(day.Batches, batches...)
		day.MonthClose = day.MonthClose || name == "CLOSE"
		day.YearEnd = day.YearEnd || name == "ANNUAL"
	}
	return day, in.count(g, &day)
}

// File is a file of the insurer's data a day read or wrote, by its path
// under data/.
type File struct {
	Path string
	Size int64
}

// Files lists the files of a day: those under data/ named by its date,
// the files it read from IN, and those it wrote to OUT and RPT, in the
// order of their directories and then their names.
func (in Insurer) Files(date generate.Date) ([]File, error) {
	stamp := "." + date.Stamp() + "."
	var files []File
	for _, dir := range []string{"IN", "OUT", "RPT"} {
		entries, err := os.ReadDir(in.data(dir))
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.IsDir() || !strings.Contains(e.Name(), stamp) {
				continue
			}
			info, err := e.Info()
			if err != nil {
				return nil, err
			}
			files = append(files, File{Path: dir + "/" + e.Name(), Size: info.Size()})
		}
	}
	return files, nil
}

// ReadFile reads a file of the insurer's data by its path under data/, as
// it is, in Shift_JIS; one outside the data is not there.
func (in Insurer) ReadFile(path string) ([]byte, error) {
	clean := filepath.Clean(filepath.FromSlash(path))
	if !filepath.IsLocal(clean) {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(in.data(clean))
}

// Policies reads the policy master as the last day's unload left it, or
// the extract the insurer opened with before any day.
func (in Insurer) Policies() ([]byte, error) {
	unloads, err := filepath.Glob(in.data("OUT", "POLMST.*.TXT"))
	if err != nil {
		return nil, err
	}
	if len(unloads) == 0 {
		return os.ReadFile(in.data("IN", "EXTRACT.TXT"))
	}
	slices.Sort(unloads)
	return os.ReadFile(unloads[len(unloads)-1])
}
