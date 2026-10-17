// Command survey reads a COBOL system and writes its blueprints, or serves
// its explorer.
//
//	survey blueprints -copybooks dir -out dir [-title name] program.cbl... [job.jcl...]
//	survey explore -copybooks dir [-addr host:port] [-title name] program.cbl... [job.jcl...]
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/explore"
)

const usage = `usage:
  survey blueprints -copybooks dir -out dir [-title name] program.cbl... [job.jcl...]
  survey explore -copybooks dir [-addr host:port] [-title name] program.cbl... [job.jcl...]
`

func main() {
	commands := map[string]func([]string) error{"blueprints": writeBlueprints, "explore": serveExplorer}
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
