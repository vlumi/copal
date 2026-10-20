// Command laurasia runs Laurasia Card Services' night on a fixture, and
// shows its logs and reports.
//
//	laurasia night [-fixture dir] [-out dir]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/vlumi/copal/laurasia/night"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "night" {
		fmt.Fprintln(os.Stderr, "usage: laurasia night [-fixture dir] [-out dir]")
		os.Exit(2)
	}
	flags := flag.NewFlagSet("night", flag.ExitOnError)
	fixture := flags.String("fixture", "laurasia/testdata/night", "the fixture to run")
	out := flags.String("out", "laurasia/run/night", "where the programs are compiled and the data directory made")
	_ = flags.Parse(os.Args[2:])
	if err := runNight(*fixture, *out); err != nil {
		fmt.Fprintln(os.Stderr, "laurasia:", err)
		os.Exit(1)
	}
}

// runNight runs the fixture afresh under out, then prints each log and
// report the night wrote.
func runNight(fixture, out string) error {
	if err := os.RemoveAll(out); err != nil {
		return err
	}
	bin, data := filepath.Join(out, "bin"), filepath.Join(out, "data")
	if err := night.Build("laurasia", bin); err != nil {
		return err
	}
	if err := night.Run("laurasia", fixture, bin, data); err != nil {
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
