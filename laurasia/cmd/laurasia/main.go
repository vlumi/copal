// Command laurasia builds and runs Laurasia Card Services.
//
//	laurasia build  compile the COBOL programs
//	laurasia night  run the night on a fixture and show its logs and reports
//	laurasia open   open the running processor on a generated portfolio
//	laurasia days   run the running processor's next days
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/vlumi/copal/laurasia/night"
	"github.com/vlumi/copal/laurasia/processor"
)

const usage = `usage:
  laurasia build [-out dir]
  laurasia night [-fixture dir] [-out dir]
  laurasia open [-processor dir] [-seed n] [-cards n]
  laurasia days [-processor dir] [-n days]
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	commands := map[string]func([]string) error{
		"build": build,
		"night": runNight,
		"open":  open,
		"days":  days,
	}
	command, ok := commands[os.Args[1]]
	if !ok {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	if err := command(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "laurasia:", err)
		os.Exit(1)
	}
}

func build(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	out := flags.String("out", "laurasia/build", "directory for the executables")
	_ = flags.Parse(args)
	return night.Build("laurasia", *out)
}

// runNight runs the fixture afresh under out, then prints each log and
// report the night wrote.
func runNight(args []string) error {
	flags := flag.NewFlagSet("night", flag.ExitOnError)
	fixture := flags.String("fixture", "laurasia/testdata/night", "the fixture to run")
	out := flags.String("out", "laurasia/run/night", "where the programs are compiled and the data directory made")
	_ = flags.Parse(args)
	if err := os.RemoveAll(*out); err != nil {
		return err
	}
	bin, data := filepath.Join(*out, "bin"), filepath.Join(*out, "data")
	if err := night.Build("laurasia", bin); err != nil {
		return err
	}
	if err := night.Run("laurasia", *fixture, bin, data); err != nil {
		return err
	}
	var shown []string
	for _, dir := range []string{"logs", "reports"} {
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
		fmt.Printf("\n== %s\n%s", rel, content)
	}
	return nil
}

func processorFlags(name string) (*flag.FlagSet, *processor.Processor) {
	flags := flag.NewFlagSet(name, flag.ExitOnError)
	p := &processor.Processor{Scripts: "laurasia/scripts", Copybooks: "laurasia/copybooks"}
	flags.StringVar(&p.Dir, "processor", "laurasia/run/processor", "directory of the running processor")
	flags.StringVar(&p.Programs, "programs", "laurasia/build", "directory of the compiled programs")
	return flags, p
}

func open(args []string) error {
	flags, p := processorFlags("open")
	seed := flags.Uint64("seed", 1, "seed the processor is generated from")
	cards := flags.Int("cards", 1000, "number of cards")
	_ = flags.Parse(args)
	if err := p.Open(*seed, *cards); err != nil {
		return err
	}
	fmt.Printf("processor opened in %s: seed %d, %d cards, first day %s\n", p.Dir, *seed, *cards, processor.StartDate)
	return nil
}

func days(args []string) error {
	flags, p := processorFlags("days")
	n := flags.Int("n", 1, "number of days to run")
	_ = flags.Parse(args)
	return p.RunDays(*n, func(d processor.Day) {
		night := fmt.Sprintf("%-47s", "no clearing file")
		if d.Night != nil {
			night = fmt.Sprintf("%5d presentments  %3d rejected  %5d postings", d.Presentments, d.Rejected, d.Postings)
		}
		fmt.Printf("%s %s  %s  %3d statements  RC=%d\n", d.Date, d.Date.Weekday().String()[:3], night, d.Statements, d.RC())
	})
}
