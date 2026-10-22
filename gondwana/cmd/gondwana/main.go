// Command gondwana builds and runs Gondwana Life Insurance.
//
//	gondwana build  compile the COBOL programs
//	gondwana month  run a month on a fixture and show its logs and reports
//
// What it shows is in Shift_JIS, as the programs write it; pipe it
// through iconv for a terminal.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/vlumi/copal/gondwana/batch"
)

const usage = `usage:
  gondwana build [-out dir]
  gondwana month [-fixture dir] [-out dir]
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	commands := map[string]func([]string) error{
		"build": build,
		"month": runMonth,
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
