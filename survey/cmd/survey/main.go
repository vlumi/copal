// Command survey reads a COBOL system and writes its blueprints, serves
// its explorer, or searches for cases that reach the rules of its
// programs.
//
//	survey blueprints -copybooks dir -out dir [-title name] program.cbl... [job.jcl...]
//	survey explore -copybooks dir [-addr host:port] [-title name] program.cbl... [job.jcl...]
//	survey oracle -copybooks dir [-program name] program.cbl...
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/explore"
	"github.com/vlumi/copal/survey/internal/oracle"
)

const usage = `usage:
  survey blueprints -copybooks dir -out dir [-title name] program.cbl... [job.jcl...]
  survey explore -copybooks dir [-addr host:port] [-title name] program.cbl... [job.jcl...]
  survey oracle -copybooks dir [-program name] program.cbl...
`

func main() {
	commands := map[string]func([]string) error{"blueprints": writeBlueprints, "explore": serveExplorer, "oracle": searchCases}
	if len(os.Args) < 2 || commands[os.Args[1]] == nil {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	if err := commands[os.Args[1]](os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "survey:", err)
		os.Exit(1)
	}
}

func writeBlueprints(args []string) error {
	flags := flag.NewFlagSet("blueprints", flag.ExitOnError)
	copybooks := flags.String("copybooks", "", "directory of the copybooks")
	out := flags.String("out", "", "directory to write the blueprints to")
	title := flags.String("title", "", "name of the system, for the index")
	_ = flags.Parse(args)
	if *copybooks == "" || *out == "" || flags.NArg() == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	system, err := blueprint.Analyze(flags.Args(), *copybooks)
	if err != nil {
		return err
	}
	return system.Write(*out, *title)
}

func serveExplorer(args []string) error {
	flags := flag.NewFlagSet("explore", flag.ExitOnError)
	copybooks := flags.String("copybooks", "", "directory of the copybooks")
	addr := flags.String("addr", "127.0.0.1:1192", "address to serve on")
	title := flags.String("title", "the system", "name of the system")
	_ = flags.Parse(args)
	if *copybooks == "" || flags.NArg() == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	system, err := blueprint.Analyze(flags.Args(), *copybooks)
	if err != nil {
		return err
	}
	x, err := explore.New(system, *title)
	if err != nil {
		return err
	}
	fmt.Printf("the explorer is at http://%s/\n", *addr)
	server := &http.Server{Addr: *addr, Handler: x.Handler(), ReadHeaderTimeout: 10 * time.Second}
	return server.ListenAndServe()
}

// searchCases searches for cases that reach the rules of each main
// program, or of the one named, and of the programs it calls, and reports
// how many rules they reach.
func searchCases(args []string) error {
	flags := flag.NewFlagSet("oracle", flag.ExitOnError)
	copybooks := flags.String("copybooks", "", "directory of the copybooks")
	only := flags.String("program", "", "the main program to search for, instead of every one")
	_ = flags.Parse(args)
	if *copybooks == "" || flags.NArg() == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	system, err := blueprint.Analyze(flags.Args(), *copybooks)
	if err != nil {
		return err
	}
	scratch, err := os.MkdirTemp("", "survey-oracle-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(scratch) }()
	for i := range system.Programs {
		p := &system.Programs[i]
		if p.Kind != blueprint.KindMain || *only != "" && p.Name != *only {
			continue
		}
		result, err := oracle.Search(system, p, *copybooks, filepath.Join(scratch, p.Name))
		if err != nil {
			return err
		}
		for u, unit := range result.Units {
			reached, total := result.Rules(u)
			fmt.Printf("%-8s %-8s %3d of %3d rules\n", p.Name, unit.Program, reached, total)
		}
		fmt.Printf("%-8s %d cases\n", p.Name, len(result.Cases))
	}
	return nil
}
