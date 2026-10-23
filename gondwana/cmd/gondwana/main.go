// Command gondwana builds and runs Gondwana Life Insurance.
//
//	gondwana build  compile the COBOL programs
//	gondwana month  run a month on a fixture and show its logs and reports
//	gondwana open   open the running insurer on a generated portfolio
//	gondwana days   run the running insurer's next business days
//
// What month shows is in Shift_JIS, as the programs write it; pipe it
// through iconv for a terminal.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vlumi/copal/gondwana/batch"
	"github.com/vlumi/copal/gondwana/insurer"
)

const usage = `usage:
  gondwana build [-out dir]
  gondwana month [-fixture dir] [-out dir]
  gondwana open [-insurer dir] [-seed n] [-policies n]
  gondwana days [-insurer dir] [-n days]
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	commands := map[string]func([]string) error{
		"build": build,
		"month": runMonth,
		"open":  open,
		"days":  days,
	}
	command, ok := commands[os.Args[1]]
	if !ok {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	if err := command(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "gondwana:", err)
		os.Exit(1)
	}
}

func build(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	out := flags.String("out", "gondwana/build", "directory for the executables")
	_ = flags.Parse(args)
	return batch.Build("gondwana", *out)
}

// runMonth runs the fixture afresh under out, then prints each log and
// report the month wrote.
func runMonth(args []string) error {
	flags := flag.NewFlagSet("month", flag.ExitOnError)
	fixture := flags.String("fixture", "gondwana/testdata/march", "the fixture to run")
	out := flags.String("out", "gondwana/run/month", "where the programs are compiled and the data directory made")
	_ = flags.Parse(args)
	if err := os.RemoveAll(*out); err != nil {
		return err
	}
	bin, data := filepath.Join(*out, "bin"), filepath.Join(*out, "data")
	if err := batch.Build("gondwana", bin); err != nil {
		return err
	}
	if err := batch.Run("gondwana", *fixture, bin, data); err != nil {
		return err
	}
	var shown []string
	for _, dir := range []string{"LOG", "RPT"} {
		found, err := filepath.Glob(filepath.Join(data, dir, "*"))
		if err != nil {
			return err
		}
		shown = append(shown, found...)
	}
	sort.Strings(shown)
	for _, path := range shown {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(data, path)
		fmt.Printf("== %s\n%s\n", rel, content)
	}
	return nil
}

func insurerFlags(name string) (*flag.FlagSet, *insurer.Insurer) {
	flags := flag.NewFlagSet(name, flag.ExitOnError)
	in := &insurer.Insurer{Scripts: "gondwana/scripts", Copybooks: "gondwana/copybooks"}
	flags.StringVar(&in.Dir, "insurer", "gondwana/run/insurer", "directory of the running insurer")
	flags.StringVar(&in.Programs, "programs", "gondwana/build", "directory of the compiled programs")
	return flags, in
}

func open(args []string) error {
	flags, in := insurerFlags("open")
	seed := flags.Uint64("seed", 1, "seed the insurer is generated from")
	policies := flags.Int("policies", 2000, "number of policies")
	_ = flags.Parse(args)
	if err := in.Open(*seed, *policies); err != nil {
		return err
	}
	fmt.Printf("insurer opened in %s: seed %d, %d policies, first business day %s\n", in.Dir, *seed, *policies, insurer.StartDate)
	return nil
}

// days runs the next business days, a line each: what came in, where
// the portfolio stands, and the month's billing, results, close and
// year end on their days.
func days(args []string) error {
	flags, in := insurerFlags("days")
	n := flags.Int("n", 1, "number of business days to run")
	_ = flags.Parse(args)
	return in.RunDays(*n, func(d insurer.Day) {
		var events []string
		if d.Billed > 0 {
			events = append(events, fmt.Sprintf("billed %d", d.Billed))
		}
		if d.Collected+d.Unpaid > 0 {
			events = append(events, fmt.Sprintf("collected %d, unpaid %d", d.Collected, d.Unpaid))
		}
		if d.MonthClose {
			events = append(events, "close")
		}
		if d.YearEnd {
			events = append(events, "year end")
		}
		fmt.Printf("%s %s  %3d new  %3d benefits  %3d changes  %3d rejects  in force %5d  grace %3d  lapsed %3d  RC=%d  %s\n",
			d.Date, d.Date.Weekday().String()[:3], d.New, d.Benefits, d.Changes, d.Rejected, d.InForce, d.InGrace, d.Lapsed, d.RC(), strings.Join(events, ", "))
	})
}
