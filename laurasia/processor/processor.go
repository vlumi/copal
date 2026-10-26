// Package processor runs Laurasia Card Services day after day: it opens
// the processor on a generated portfolio, and for each day generates what
// the switch logged, what the networks cleared and what the disputes
// desk received, runs the night on a business day and the billing cycle
// on every day, and backs a day out whole when it fails.
package processor

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vlumi/cobolith/record"
	"github.com/vlumi/copal/bedrock/runlog"
	"github.com/vlumi/copal/bedrock/runner"
	"github.com/vlumi/copal/laurasia/generate"
	"github.com/vlumi/copal/laurasia/night"
	"github.com/vlumi/copal/laurasia/records"
)

// StartDate is the first day of every processor opened.
var StartDate = generate.NewDate(2026, time.October, 1)

// Processor is a card processor in a directory: its seed, its data as the
// scripts expect it under data/, and under failed/DATE/ what a day that
// did not go through left behind.
type Processor struct {
	Dir       string
	Programs  string // the compiled COBOL programs
	Scripts   string
	Copybooks string
}

// Step is a step of a script as its log records it: its name, the program
// or command it ran, and the return code.
type Step struct {
	Name    string
	Program string
	RC      int
}

// Day is the outcome of one day: the nightly run's steps, none on a
// weekend when no clearing file comes, the billing cycle's, and what the
// day amounted to.
type Day struct {
	Date         generate.Date
	Night        []Step
	Cycle        []Step
	Presentments int // in the network's clearing file and the on-us one
	Rejected     int // by the clearing, the posting and the card system's feed
	Postings     int
	Payments     int
	Statements   int
}

// RC is the highest return code of the day's steps.
func (d Day) RC() int {
	rc := 0
	for _, s := range append(d.Night, d.Cycle...) {
		rc = max(rc, s.RC)
	}
	return rc
}

const seedFile = "SEED"

// ErrNotOpen is for a directory where no processor has been opened yet.
var ErrNotOpen = errors.New("no processor here yet; open one first")

// Open replaces the processor with a new one: about the given number of
// cards and the fee schedule, generated from the seed, no disputes in
// flight, and the accounts loaded by the script that opens the processor.
func (p Processor) Open(seed uint64, cards int) error {
	if err := os.RemoveAll(p.Dir); err != nil {
		return err
	}
	g, err := p.generator(seed)
	if err != nil {
		return err
	}
	accounts, err := g.Cards(cards)
	if err != nil {
		return err
	}
	fees, err := g.Fees()
	if err != nil {
		return err
	}
	for _, f := range []struct {
		name    string
		records [][]byte
	}{
		{"inbound/accounts.txt", accounts},
		{"inbound/FEES.txt", fees},
		{"master/disputes.txt", nil},
	} {
		if err := p.write(f.name, f.records); err != nil {
			return err
		}
	}
	if err := p.runner().Run("open.sh"); err != nil {
		return err
	}
	opened, err := p.logSteps(p.data("logs", "open.log"))
	if err != nil {
		return err
	}
	if err := p.record(StartDate, "open", opened); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(p.Dir, seedFile), []byte(strconv.FormatUint(seed, 10)+"\n"), 0o644)
}

// RunDays runs the next n days, reporting each as it ends. It stops at a
// day that did not go through.
func (p Processor) RunDays(n int, report func(Day)) error {
	g, err := p.open()
	if err != nil {
		return err
	}
	for range n {
		date, err := p.nextDate()
		if err != nil {
			return err
		}
		day, err := p.runDay(g, date)
		if err != nil {
			return err
		}
		if err := errors.Join(p.record(date, "nightly", day.Night), p.record(date, "cycle", day.Cycle)); err != nil {
			return err
		}
		report(day)
	}
	return nil
}

// runDay runs a day, which goes through whole or not at all: if any of
// it fails, the masters are put back as they were, and everything the
// day wrote, the generated inputs included, moves to failed/DATE, so
// that the next run starts the day afresh.
func (p Processor) runDay(g *generate.Generator, date generate.Date) (Day, error) {
	before := filepath.Join(p.Dir, "master.before")
	if err := os.RemoveAll(before); err != nil {
		return Day{}, err
	}
	if err := os.CopyFS(before, os.DirFS(p.data("master"))); err != nil {
		return Day{}, err
	}
	day, err := p.tryDay(g, date)
	if err != nil {
		if backOutErr := p.backOut(before, date); backOutErr != nil {
			return day, errors.Join(err, backOutErr)
		}
		return day, fmt.Errorf("%w; the day was backed out, what it left is in %s", err, p.failed(date))
	}
	return day, os.RemoveAll(before)
}

func (p Processor) tryDay(g *generate.Generator, date generate.Date) (Day, error) {
	day := Day{Date: date}
	r := p.runner()
	if date.BusinessDay() {
		n, err := p.generateNight(g, date)
		if err != nil {
			return day, err
		}
		day.Presentments = n.Presentments
		if err := r.Run("nightly.sh", date.Stamp()); err != nil {
			return day, err
		}
		if day.Night, err = p.steps("night", date); err != nil {
			return day, err
		}
	}
	if err := r.Run("cycle.sh", date.Stamp()); err != nil {
		return day, err
	}
	var err error
	if day.Cycle, err = p.steps("cycle", date); err != nil {
		return day, err
	}
	return day, p.count(&day)
}

// generateNight writes the night's inbound files: the switch's log, the
// clearing file and the card system's feed, from the cards as the last
// cycle unloaded them, and the dispute events, from the cases in flight
// and last night's sales.
func (p Processor) generateNight(g *generate.Generator, date generate.Date) (generate.Night, error) {
	accounts, err := p.readFile(g.Layouts.Account, p.data("outbound", "ACCOUNTS.CYCLE."+date.AddDays(-1).Stamp()+".txt"))
	if errors.Is(err, os.ErrNotExist) {
		accounts, err = p.readFile(g.Layouts.Account, p.data("inbound", "accounts.txt"))
	}
	if err != nil {
		return generate.Night{}, err
	}
	cards, err := g.ReadCards(accounts)
	if err != nil {
		return generate.Night{}, err
	}
	n, err := g.Night(date, date.PreviousBusinessDay(), cards, g.Merchants())
	if err != nil {
		return n, err
	}
	payments, err := g.Payments(date, cards)
	if err != nil {
		return n, err
	}
	master, err := p.readFile(g.Layouts.Case, p.data("master", "disputes.txt"))
	if err != nil {
		return n, err
	}
	cases, err := g.ReadCases(master)
	if err != nil {
		return n, err
	}
	priced, err := p.readFile(g.Layouts.Pres, p.data("work", date.PreviousBusinessDay().Stamp(), "pres.priced"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return n, err
	}
	sales, err := g.ReadSales(priced)
	if err != nil {
		return n, err
	}
	events, err := g.Disputes(date, cases, sales)
	if err != nil {
		return n, err
	}
	for name, recs := range map[string][][]byte{"AUTH": n.Log, "CLR": n.Clearing, "ONUS": n.OnUs, "PAYMENTS": payments, "DISPEVT": events} {
		if err := p.write(filepath.Join("inbound", name+"."+date.Stamp()+".txt"), recs); err != nil {
			return n, err
		}
	}
	return n, nil
}

// count reads what the day amounted to from what it wrote.
func (p Processor) count(day *Day) error {
	stamp := day.Date.Stamp()
	for _, c := range []struct {
		file  string
		count *int
		only  func(string) bool
	}{
		{"CLRREJ", &day.Rejected, nil},
		{"ONUSREJ", &day.Rejected, nil},
		{"POSTREJ", &day.Rejected, nil},
		{"PAYREJ", &day.Rejected, nil},
		{"POSTINGS", &day.Postings, nil},
		{"CREDITS", &day.Payments, nil},
		{"STATEMENTS", &day.Statements, func(line string) bool { return !strings.HasPrefix(line, " ") }},
	} {
		content, err := os.ReadFile(p.data("outbound", c.file+"."+stamp+".txt"))
		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return err
		}
		for _, line := range strings.Split(strings.TrimRight(string(content), "\n"), "\n") {
			if line != "" && (c.only == nil || c.only(line)) {
				*c.count++
			}
		}
	}
	return nil
}

// steps reads a script's log of a date.
func (p Processor) steps(script string, date generate.Date) ([]Step, error) {
	return p.logSteps(p.data("logs", script+"."+date.Stamp()+".log"))
}

// logSteps reads a script's log: a step a line, its name, its program and
// its return code.
func (p Processor) logSteps(path string) ([]Step, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var steps []Step
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("%s: %q is not a step", path, line)
		}
		rc, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("%s: %q is not a step", path, line)
		}
		steps = append(steps, Step{Name: fields[0], Program: fields[1], RC: rc})
	}
	return steps, nil
}

// record adds the steps a script ran on a date to the processor's run
// record, under the script's name; the night's log holds the steps of the
// scripts the nightly one calls, and so does its record. A script that did
// not run, the night on a weekend, is not added.
func (p Processor) record(date generate.Date, script string, steps []Step) error {
	if len(steps) == 0 {
		return nil
	}
	run := runlog.Run{Date: date.String(), Job: script}
	for _, s := range steps {
		run.Steps = append(run.Steps, runner.Result{Step: s.Name, Program: s.Program, RC: s.RC})
	}
	return runlog.Append(filepath.Join(p.Dir, runlog.File), run)
}

// backOut puts the masters back as they were before the day and moves
// everything named by the day's date, under data/, to failed/DATE.
func (p Processor) backOut(before string, date generate.Date) error {
	if err := os.RemoveAll(p.data("master")); err != nil {
		return err
	}
	if err := os.Rename(before, p.data("master")); err != nil {
		return err
	}
	failed := p.failed(date)
	if err := os.RemoveAll(failed); err != nil {
		return err
	}
	stamp := date.Stamp()
	err := filepath.WalkDir(p.data(), func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(p.data(), path)
		if !strings.Contains(rel, stamp) {
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(failed, rel)), 0o755); err != nil {
			return err
		}
		return os.Rename(path, filepath.Join(failed, rel))
	})
	if err != nil {
		return err
	}
	return os.RemoveAll(p.data("work", stamp))
}

// nextDate is the day after the last the cycle ran, which it does every
// day, or the start date in a processor just opened.
func (p Processor) nextDate() (generate.Date, error) {
	logs, err := filepath.Glob(p.data("logs", "cycle.*.log"))
	if err != nil {
		return generate.Date{}, err
	}
	last := StartDate.AddDays(-1)
	for _, log := range logs {
		stamp := strings.TrimSuffix(strings.TrimPrefix(filepath.Base(log), "cycle."), ".log")
		n, err := strconv.ParseInt(stamp, 10, 64)
		if err != nil {
			return generate.Date{}, fmt.Errorf("%s: %w", log, err)
		}
		date, err := generate.DateFromNumber(n)
		if err != nil {
			return generate.Date{}, fmt.Errorf("%s: %w", log, err)
		}
		if date.After(last) {
			last = date
		}
	}
	return last.AddDays(1), nil
}

// readFile reads a line sequential file as records of a layout.
func (p Processor) readFile(l *record.Layout, path string) ([][]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return generate.Records(l, content), nil
}

// write writes records as a line sequential file under data/.
func (p Processor) write(name string, records [][]byte) error {
	path := p.data(name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, generate.File(records), 0o644)
}

func (p Processor) runner() night.Runner {
	return night.Runner{Programs: p.Programs, Scripts: p.Scripts, Data: p.data()}
}

func (p Processor) data(elem ...string) string {
	return filepath.Join(append([]string{p.Dir, "data"}, elem...)...)
}

func (p Processor) failed(date generate.Date) string {
	return filepath.Join(p.Dir, "failed", date.Stamp())
}

func (p Processor) seed() (uint64, error) {
	data, err := os.ReadFile(filepath.Join(p.Dir, seedFile))
	if errors.Is(err, os.ErrNotExist) {
		return 0, fmt.Errorf("%s: %w", p.Dir, ErrNotOpen)
	} else if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

// open opens the processor: its generator, from its seed.
func (p Processor) open() (*generate.Generator, error) {
	seed, err := p.seed()
	if err != nil {
		return nil, err
	}
	return p.generator(seed)
}

func (p Processor) generator(seed uint64) (*generate.Generator, error) {
	layouts, err := records.Load(p.Copybooks)
	if err != nil {
		return nil, err
	}
	return &generate.Generator{Layouts: layouts, Seed: seed}, nil
}
