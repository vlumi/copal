// Command survey reads a COBOL system and writes its blueprints, serves
// its explorer, or searches for cases that reach the rules of its
// programs.
//
//	survey blueprints -copybooks dir -out dir [-dialect ibm|mf] [-title name] program.cbl... [job.jcl...]
//	survey explore -copybooks dir [-dialect ibm|mf] [-addr host:port] [-title name] [-cases dir] program.cbl... [job.jcl...]
//	survey oracle -copybooks dir [-dialect ibm|mf] [-out dir -rules dir] [-title name] [-program name] [-jobs JOB,...] program.cbl... [job.jcl...]
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/explore"
	"github.com/vlumi/copal/survey/internal/oracle"
)

const usage = `usage:
  survey blueprints -copybooks dir -out dir [-dialect ibm|mf] [-title name] program.cbl... [job.jcl...]
  survey explore -copybooks dir [-dialect ibm|mf] [-addr host:port] [-title name] [-cases dir] program.cbl... [job.jcl...]
  survey oracle -copybooks dir [-dialect ibm|mf] [-out dir -rules dir] [-title name] [-program name] [-jobs JOB,...] program.cbl... [job.jcl...]
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
	dialect := flags.String("dialect", cobol.IBM.Name, "the dialect the system is written in, ibm or mf")
	out := flags.String("out", "", "directory to write the blueprints to")
	title := flags.String("title", "", "name of the system, for the index")
	_ = flags.Parse(args)
	if *copybooks == "" || *out == "" || flags.NArg() == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	system, err := analyzeSystem(flags.Args(), *copybooks, *dialect)
	if err != nil {
		return err
	}
	return system.Write(*out, *title)
}

func serveExplorer(args []string) error {
	flags := flag.NewFlagSet("explore", flag.ExitOnError)
	copybooks := flags.String("copybooks", "", "directory of the copybooks")
	dialect := flags.String("dialect", cobol.IBM.Name, "the dialect the system is written in, ibm or mf")
	addr := flags.String("addr", "127.0.0.1:1192", "address to serve on")
	title := flags.String("title", "the system", "name of the system")
	cases := flags.String("cases", "", "directory of the cases the oracle wrote")
	_ = flags.Parse(args)
	if *copybooks == "" || flags.NArg() == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	system, err := analyzeSystem(flags.Args(), *copybooks, *dialect)
	if err != nil {
		return err
	}
	x, err := explore.New(system, *title, *cases)
	if err != nil {
		return err
	}
	fmt.Printf("the explorer is at http://%s/\n", *addr)
	server := &http.Server{Addr: *addr, Handler: x.Handler(), ReadHeaderTimeout: 10 * time.Second}
	return server.ListenAndServe()
}

// searchCases searches for cases that reach the rules of each main
// program, or of the one named, and of the programs it calls, reports how
// many rules they reach, and writes them, each program's to a directory
// of its own, with an index of them all.
func searchCases(args []string) error {
	flags := flag.NewFlagSet("oracle", flag.ExitOnError)
	copybooks := flags.String("copybooks", "", "directory of the copybooks")
	dialect := flags.String("dialect", cobol.IBM.Name, "the dialect the system is written in, ibm or mf")
	out := flags.String("out", "", "directory to write the cases to")
	rules := flags.String("rules", "", "directory of the blueprints' decision tables, relative to the cases' directory")
	title := flags.String("title", "the system", "name of the system, for the index")
	only := flags.String("program", "", "the main program to search for, instead of every one")
	order := flags.String("jobs", "", "jobs to run as a batch on the cases, in order, separated by commas")
	_ = flags.Parse(args)
	if *copybooks == "" || flags.NArg() == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	system, err := analyzeSystem(flags.Args(), *copybooks, *dialect)
	if err != nil {
		return err
	}
	scratch, err := os.MkdirTemp("", "survey-oracle-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(scratch) }()
	var results []*oracle.Result
	for i := range system.Programs {
		p := &system.Programs[i]
		if p.Kind != blueprint.KindMain || *only != "" && p.Name != *only {
			continue
		}
		result, err := oracle.Search(system, p, *copybooks, scratch)
		if err != nil {
			return err
		}
		for u, unit := range result.Units {
			reached, total := result.Rules(u)
			fmt.Printf("%-8s %-8s %3d of %3d rules\n", p.Name, unit.Program, reached, total)
		}
		fmt.Printf("%-8s %d cases\n", p.Name, len(result.Cases))
		results = append(results, result)
		if *out != "" {
			if err := writeCases(filepath.Join(*out, p.Name), result.Files(filepath.Join("..", *rules))); err != nil {
				return err
			}
		}
	}
	if *out == "" || *only != "" {
		return nil
	}
	byName := map[string]*oracle.Result{}
	for _, r := range results {
		byName[r.Program] = r
	}
	var jobs []*blueprint.Job
	for _, name := range strings.Split(*order, ",") {
		i := slices.IndexFunc(system.Jobs, func(j blueprint.Job) bool { return j.Name == name })
		if i < 0 && name != "" {
			return fmt.Errorf("no job %s", name)
		}
		if i >= 0 {
			jobs = append(jobs, &system.Jobs[i])
		}
	}
	batches, err := oracle.RunJobs(jobs, byName, scratch)
	if err != nil {
		return err
	}
	for _, batch := range batches {
		fmt.Printf("%-8s a batch\n", batch.Job)
		if err := writeCases(filepath.Join(*out, "jobs", batch.Job), batch.Files()); err != nil {
			return err
		}
	}
	return os.WriteFile(filepath.Join(*out, "README.md"), oracle.Index(*title, results, batches), 0o644)
}

// analyzeSystem reads the programs, jobs and copybooks of a system written
// in the dialect named.
func analyzeSystem(paths []string, copybooks, dialect string) (*blueprint.System, error) {
	d, ok := cobol.DialectNamed(dialect)
	if !ok {
		return nil, fmt.Errorf("no dialect %s", dialect)
	}
	return blueprint.Analyze(paths, copybooks, d)
}

// writeCases writes a program's cases to its directory, in place of those
// there before.
func writeCases(dir string, files map[string][]byte) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	for name, data := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
