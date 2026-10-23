// Package insurer runs Gondwana Life Insurance business day after
// business day: it opens the system on a generated portfolio, and for
// each day writes what the sales system sold and what the branches
// requested, runs the daily batch, and around it the month's billing,
// the collection agent's results, the month-end close and the year-end
// dividends on their days, and backs a day out whole when it fails.
package insurer

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
	"github.com/vlumi/copal/gondwana/batch"
	"github.com/vlumi/copal/gondwana/generate"
	"github.com/vlumi/copal/gondwana/records"
)

// StartDate is the day every insurer opens on: the first day of a fiscal
// year, so that a year of days ends with the annual run.
var StartDate = generate.NewDate(2026, time.April, 1)

// Insurer is a life insurer in a directory: its seed, its data as the
// batch files expect it under data/, and under failed/DATE/ what a day
// that did not go through left behind.
type Insurer struct {
	Dir       string
	Programs  string // the compiled COBOL programs
	Scripts   string // the batch files
	Copybooks string
}

// Step is a step of a batch file as its log records it: its name, the
// program it ran, and the errorlevel.
type Step struct {
	Name    string
	Program string
	RC      int
}

// Batch is a batch file the day ran, with its steps.
type Batch struct {
	Name  string
	Steps []Step
}

// Day is the outcome of one business day: the batch files it ran, what
// came in and went out, and where the portfolio stands after it.
type Day struct {
	Date       generate.Date
	Batches    []Batch
	New        int // policies registered
	Benefits   int // payments made
	Changes    int // changes applied
	Rejected   int // records rejected by any step
	Billed     int // debits requested, on a billing day
	Collected  int // debits that went through, on a results day
	Unpaid     int // debits that failed
	InForce    int
	InGrace    int
	Lapsed     int
	Ended      int // surrendered, matured or paid out
	MonthClose bool
	YearEnd    bool
}

// RC is the highest errorlevel of the day's steps.
func (d Day) RC() int {
	rc := 0
	for _, b := range d.Batches {
		for _, s := range b.Steps {
			rc = max(rc, s.RC)
		}
	}
	return rc
}

const seedFile = "SEED"

// ErrNotOpen is for a directory where no insurer has been opened yet.
var ErrNotOpen = errors.New("no insurer here yet; open one first")

// Open replaces the insurer with a new one: the extract of about the
// given number of policies and the tables, generated from the seed, and
// the master loaded from them by the batch file that opens the system.
func (in Insurer) Open(seed uint64, policies int) error {
	if err := os.RemoveAll(in.Dir); err != nil {
		return err
	}
	g, err := in.generator(seed)
	if err != nil {
		return err
	}
	extract, err := g.Extract(policies, StartDate)
	if err != nil {
		return err
	}
	if err := in.write("IN/EXTRACT.TXT", generate.File(extract)); err != nil {
		return err
	}
	for name, table := range map[string]func() ([][]byte, error){
		"PRODTB": g.Products, "RATETB": g.Rates, "CSVTB": g.CashValues, "RSVTB": g.Reserves, "AGTMST": g.Agents,
	} {
		recs, err := table()
		if err != nil {
			return err
		}
		if err := in.write(filepath.Join("MASTER", name+".TXT"), generate.File(recs)); err != nil {
			return err
		}
	}
	if err := in.runner().Run("OPEN.BAT", StartDate.Stamp()); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(in.Dir, seedFile), []byte(strconv.FormatUint(seed, 10)+"\n"), 0o644)
}

// RunDays runs the next n business days, reporting each as it ends. It
// stops at a day that did not go through.
func (in Insurer) RunDays(n int, report func(Day)) error {
	g, err := in.open()
	if err != nil {
		return err
	}
	for range n {
		date, err := in.nextDate()
		if err != nil {
			return err
		}
		day, err := in.runDay(g, date)
		if err != nil {
			return err
		}
		report(day)
	}
	return nil
}

// runDay runs a day, which goes through whole or not at all: if any of
// it fails, the master is put back as it was, and everything the day
// wrote, the generated inputs included, moves to failed/DATE, so that
// the next run starts the day afresh.
func (in Insurer) runDay(g *generate.Generator, date generate.Date) (Day, error) {
	before := filepath.Join(in.Dir, "master.before")
	if err := os.RemoveAll(before); err != nil {
		return Day{}, err
	}
	if err := os.CopyFS(before, os.DirFS(in.data("MASTER"))); err != nil {
		return Day{}, err
	}
	day, err := in.tryDay(g, date)
	if err != nil {
		if backOutErr := in.backOut(before, date); backOutErr != nil {
			return day, errors.Join(err, backOutErr)
		}
		return day, fmt.Errorf("%w; the day was backed out, what it left is in %s", err, in.failed(date))
	}
	return day, os.RemoveAll(before)
}

// tryDay runs a business day in the order the operations guide gives:
// the agent's results when they arrive, the daily run, the billing on
// its day, the close on the month's last business day and the annual
// run after the close of March.
func (in Insurer) tryDay(g *generate.Generator, date generate.Date) (Day, error) {
	day := Day{Date: date}
	run := func(file string) error {
		if err := in.runner().Run(file, date.Stamp()); err != nil {
			return err
		}
		batches, err := in.batches(strings.TrimSuffix(file, ".BAT"), date)
		day.Batches = append(day.Batches, batches...)
		return err
	}
	if request, ok, err := in.requestAnswered(g, date); err != nil {
		return day, err
	} else if ok {
		if err := in.generateResults(g, date, request); err != nil {
			return day, err
		}
		if err := run("RESULTS.BAT"); err != nil {
			return day, err
		}
	}
	if err := in.generateDaily(g, date); err != nil {
		return day, err
	}
	if err := run("DAILY.BAT"); err != nil {
		return day, err
	}
	if date == billingDay(date) {
		if err := run("BILLING.BAT"); err != nil {
			return day, err
		}
	}
	if date.LastBusinessDayOfMonth() {
		day.MonthClose = true
		if err := run("CLOSE.BAT"); err != nil {
			return day, err
		}
		if date.Month() == time.March {
			day.YearEnd = true
			if err := run("ANNUAL.BAT"); err != nil {
				return day, err
			}
		}
	}
	return day, in.count(g, &day)
}

// billingDay is the day the month's premiums are billed: the 10th, or
// the business day after it, before the debit on the 27th.
func billingDay(inMonth generate.Date) generate.Date {
	return generate.NewDate(inMonth.Year(), inMonth.Month(), 10).BusinessDayOnOrAfter()
}

// resultsDay is the day the collection agent's results for a month's
// debit arrive: the second business day after the debit day, the 27th
// or the business day after it.
func resultsDay(inMonth generate.Date) generate.Date {
	return generate.NewDate(inMonth.Year(), inMonth.Month(), 27).BusinessDayOnOrAfter().NextBusinessDay().NextBusinessDay()
}

// requestAnswered is the transfer request whose results arrive on the
// date, if any: the request of the month whose debit the date follows,
// which may be the month before when the debit was late in it.
func (in Insurer) requestAnswered(g *generate.Generator, date generate.Date) ([][]byte, bool, error) {
	for _, month := range []generate.Date{date, date.AddMonths(-1)} {
		if resultsDay(month) != date {
			continue
		}
		request, err := in.readFixed(g.Layouts.Transfer, in.data("OUT", "PRMREQ."+billingDay(month).Stamp()+".DAT"))
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return request, err == nil, err
	}
	return nil, false, nil
}

// generateResults writes the agent's results for a request.
func (in Insurer) generateResults(g *generate.Generator, date generate.Date, request [][]byte) error {
	results, err := g.Results(date, request)
	if err != nil {
		return err
	}
	return in.write(filepath.Join("IN", "PRMRES."+date.Stamp()+".DAT"), generate.FixedFile(results))
}

// generateDaily writes the daily run's inbound files: the sales system's
// new policies and the branches' benefit and change requests, from the
// policies as the last unload has them.
func (in Insurer) generateDaily(g *generate.Generator, date generate.Date) error {
	policies, err := in.policies(g, date)
	if err != nil {
		return err
	}
	for _, f := range []struct {
		name     string
		generate func(generate.Date, []generate.Policy) ([][]byte, error)
	}{{"NEWPOL", g.NewBusiness}, {"CLMREQ", g.Claims}, {"MNTREQ", g.Changes}} {
		recs, err := f.generate(date, policies)
		if err != nil {
			return err
		}
		if err := in.write(filepath.Join("IN", f.name+"."+date.Stamp()+".TXT"), generate.File(recs)); err != nil {
			return err
		}
	}
	return nil
}

// policies reads the policies as the last unload before the date left
// them, or as the extract had them in an insurer just opened.
func (in Insurer) policies(g *generate.Generator, date generate.Date) ([]generate.Policy, error) {
	unloads, err := filepath.Glob(in.data("OUT", "POLMST.*.TXT"))
	if err != nil {
		return nil, err
	}
	latest := in.data("IN", "EXTRACT.TXT")
	for _, unload := range unloads {
		if stamp := stampOf(unload); stamp < date.Stamp() && (latest == in.data("IN", "EXTRACT.TXT") || stamp > stampOf(latest)) {
			latest = unload
		}
	}
	recs, err := in.readFile(g.Layouts.Policy, latest)
	if err != nil {
		return nil, err
	}
	return g.ReadPolicies(recs)
}

// stampOf is the date in a file's name, NAME.yyyymmdd.EXT.
func stampOf(path string) string {
	parts := strings.Split(filepath.Base(path), ".")
	if len(parts) < 3 {
		return ""
	}
	return parts[1]
}

// count reads what the day amounted to from what it wrote, and where the
// portfolio stands from the unload.
func (in Insurer) count(g *generate.Generator, day *Day) error {
	stamp := day.Date.Stamp()
	lines := func(name string) (int, error) {
		content, err := os.ReadFile(in.data("OUT", name+"."+stamp+".TXT"))
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		if err != nil {
			return 0, err
		}
		return strings.Count(string(content), "\n"), nil
	}
	inbound := func(name string) (int, error) {
		content, err := os.ReadFile(in.data("IN", name+"."+stamp+".TXT"))
		if err != nil {
			return 0, err
		}
		return strings.Count(string(content), "\n"), nil
	}
	var counts [8]int
	var err error
	for i, name := range []string{"NEWREJ", "CLMREJ", "MNTREJ", "RESREJ", "PAYOUT"} {
		if counts[i], err = lines(name); err != nil {
			return err
		}
	}
	for i, name := range []string{"NEWPOL", "MNTREQ"} {
		if counts[5+i], err = inbound(name); err != nil {
			return err
		}
	}
	day.Rejected = counts[0] + counts[1] + counts[2] + counts[3]
	day.Benefits = counts[4]
	day.New = counts[5] - counts[0]
	day.Changes = counts[6] - counts[2]
	if err := in.countTransfers(g, day); err != nil {
		return err
	}
	unload, err := in.readFile(g.Layouts.Policy, in.data("OUT", "POLMST."+stamp+".TXT"))
	if err != nil {
		return err
	}
	policies, err := g.ReadPolicies(unload)
	if err != nil {
		return err
	}
	for _, p := range policies {
		switch p.Status {
		case 1:
			day.InForce++
		case 2:
			day.InGrace++
		case 3:
			day.Lapsed++
		default:
			day.Ended++
		}
	}
	return nil
}

// countTransfers counts the debits the day requested or the results it
// applied, from the banks' files.
func (in Insurer) countTransfers(g *generate.Generator, day *Day) error {
	l := g.Layouts.Transfer
	for _, f := range []struct {
		path  string
		count func(kind, result string)
	}{
		{in.data("OUT", "PRMREQ."+day.Date.Stamp()+".DAT"), func(kind, _ string) {
			if kind == "2" {
				day.Billed++
			}
		}},
		{in.data("IN", "PRMRES."+day.Date.Stamp()+".DAT"), func(kind, result string) {
			switch {
			case kind != "2":
			case result == "0":
				day.Collected++
			default:
				day.Unpaid++
			}
		}},
	} {
		recs, err := in.readFixed(l, f.path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return err
		}
		for _, rec := range recs {
			kind, _ := l.Text(rec, "ZG-ﾃﾞｰﾀｸﾌﾞﾝ")
			result, _ := l.Text(rec, "ZG-D-ﾌﾘｶｴｹｯｶ")
			f.count(kind, result)
		}
	}
	return nil
}

// batches reads the logs a batch file wrote on a date: its own, and the
// unload's when it called it.
func (in Insurer) batches(name string, date generate.Date) ([]Batch, error) {
	var batches []Batch
	names := []string{name}
	if name == "DAILY" || name == "CLOSE" {
		names = append(names, "UNLOAD")
	}
	for _, n := range names {
		steps, err := in.steps(n, date)
		if err != nil {
			return nil, err
		}
		batches = append(batches, Batch{Name: n, Steps: steps})
	}
	return batches, nil
}

// steps reads a batch file's log: a step a line, its name, its program
// and its errorlevel.
func (in Insurer) steps(name string, date generate.Date) ([]Step, error) {
	content, err := os.ReadFile(in.data("LOG", name+"."+date.Stamp()+".LOG"))
	if err != nil {
		return nil, err
	}
	var steps []Step
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("%s log of %s: %q is not a step", name, date, line)
		}
		rc, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("%s log of %s: %q is not a step", name, date, line)
		}
		steps = append(steps, Step{Name: fields[0], Program: fields[1], RC: rc})
	}
	return steps, nil
}

// backOut puts the master back as it was before the day and moves
// everything named by the day's date, under data/, to failed/DATE.
func (in Insurer) backOut(before string, date generate.Date) error {
	if err := os.RemoveAll(in.data("MASTER")); err != nil {
		return err
	}
	if err := os.Rename(before, in.data("MASTER")); err != nil {
		return err
	}
	failed := in.failed(date)
	if err := os.RemoveAll(failed); err != nil {
		return err
	}
	stamp := date.Stamp()
	return filepath.WalkDir(in.data(), func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(in.data(), path)
		if !strings.Contains(rel, stamp) {
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(failed, rel)), 0o755); err != nil {
			return err
		}
		return os.Rename(path, filepath.Join(failed, rel))
	})
}

// nextDate is the business day after the last the daily run ran, or
// the start date in an insurer just opened.
func (in Insurer) nextDate() (generate.Date, error) {
	logs, err := filepath.Glob(in.data("LOG", "DAILY.*.LOG"))
	if err != nil {
		return generate.Date{}, err
	}
	if len(logs) == 0 {
		return StartDate, nil
	}
	var last generate.Date
	for _, log := range logs {
		n, err := strconv.ParseInt(stampOf(log), 10, 64)
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
	return last.NextBusinessDay(), nil
}

// readFile reads a line sequential file as records of a layout.
func (in Insurer) readFile(l *record.Layout, path string) ([][]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return generate.Records(l, content), nil
}

// readFixed reads a fixed-length file as records of a layout.
func (in Insurer) readFixed(l *record.Layout, path string) ([][]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return generate.Fixed(l, content), nil
}

// write writes a file under data/.
func (in Insurer) write(name string, content []byte) error {
	path := in.data(name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func (in Insurer) runner() batch.Runner {
	return batch.Runner{Programs: in.Programs, Scripts: in.Scripts, Data: in.data()}
}

func (in Insurer) data(elem ...string) string {
	return filepath.Join(append([]string{in.Dir, "data"}, elem...)...)
}

func (in Insurer) failed(date generate.Date) string {
	return filepath.Join(in.Dir, "failed", date.Stamp())
}

func (in Insurer) seed() (uint64, error) {
	data, err := os.ReadFile(filepath.Join(in.Dir, seedFile))
	if errors.Is(err, os.ErrNotExist) {
		return 0, fmt.Errorf("%s: %w", in.Dir, ErrNotOpen)
	} else if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

// open opens the insurer: its generator, from its seed.
func (in Insurer) open() (*generate.Generator, error) {
	seed, err := in.seed()
	if err != nil {
		return nil, err
	}
	return in.generator(seed)
}

func (in Insurer) generator(seed uint64) (*generate.Generator, error) {
	layouts, err := records.Load(in.Copybooks)
	if err != nil {
		return nil, err
	}
	return &generate.Generator{Layouts: layouts, Seed: seed}, nil
}
