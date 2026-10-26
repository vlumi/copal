package processor

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vlumi/copal/laurasia/generate"
)

// History is what the processor has done since it opened: its seed, the
// next day it will run, and the days it ran, oldest first, read back from
// their logs and outbound files.
type History struct {
	Seed uint64
	Next generate.Date
	Days []Day
}

// History reads the processor's history.
func (p Processor) History() (History, error) {
	seed, err := p.seed()
	if err != nil {
		return History{}, err
	}
	h := History{Seed: seed}
	if h.Next, err = p.nextDate(); err != nil {
		return History{}, err
	}
	for date := StartDate; h.Next.After(date); date = date.AddDays(1) {
		day, err := p.Day(date)
		if err != nil {
			return History{}, err
		}
		h.Days = append(h.Days, day)
	}
	return h, nil
}

// Day reads back a day the processor ran: its steps from the scripts'
// logs, none of the night's on a day without one, and what it amounted to
// from its files.
func (p Processor) Day(date generate.Date) (Day, error) {
	day := Day{Date: date}
	var err error
	if _, statErr := os.Stat(p.data("logs", "night."+date.Stamp()+".log")); statErr == nil {
		if day.Night, err = p.steps("night", date); err != nil {
			return Day{}, err
		}
		if day.Presentments, err = p.presentments(date); err != nil {
			return Day{}, err
		}
	}
	if day.Cycle, err = p.steps("cycle", date); err != nil {
		return Day{}, err
	}
	return day, p.count(&day)
}

// presentments counts the presentments of a day's clearing file.
func (p Processor) presentments(date generate.Date) (int, error) {
	content, err := os.ReadFile(p.data("inbound", "CLR."+date.Stamp()+".txt"))
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	n := 0
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "PR") {
			n++
		}
	}
	return n, nil
}

// File is a file of the processor's data a day read or wrote, by its path
// under data/.
type File struct {
	Path string
	Size int64
}

// Files lists the files of a day: those under data/ named by its date, the
// inbound files it read, the outbound files and reports it wrote, in the
// order of their directories and then their names.
func (p Processor) Files(date generate.Date) ([]File, error) {
	stamp := date.Stamp()
	var files []File
	for _, dir := range []string{"inbound", "outbound", "reports"} {
		entries, err := os.ReadDir(p.data(dir))
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

// ReadFile reads a file of the processor's data by its path under data/,
// refusing one outside it.
func (p Processor) ReadFile(path string) ([]byte, error) {
	clean := filepath.Clean(filepath.FromSlash(path))
	if !filepath.IsLocal(clean) {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(p.data(clean))
}

// Cards reads the cards as the last cycle unloaded them, or as the
// processor opened with them before any cycle.
func (p Processor) Cards() ([]byte, error) {
	next, err := p.nextDate()
	if err != nil {
		return nil, err
	}
	if next.After(StartDate) {
		return os.ReadFile(p.data("outbound", "ACCOUNTS.CYCLE."+next.AddDays(-1).Stamp()+".txt"))
	}
	return os.ReadFile(p.data("inbound", "accounts.txt"))
}

// Postings reads every day's postings, oldest first, a file a day.
func (p Processor) Postings() ([][]byte, error) {
	paths, err := filepath.Glob(p.data("outbound", "POSTINGS.*.txt"))
	if err != nil {
		return nil, err
	}
	slices.Sort(paths)
	var files [][]byte
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		files = append(files, content)
	}
	return files, nil
}
