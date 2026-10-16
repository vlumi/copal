// Command pangaea builds and runs the First Bank of Pangaea.
//
//	pangaea build   compile the COBOL programs
//	pangaea load    load text files, one record per line, as datasets
//	pangaea unload  write every cataloged dataset as text
//	pangaea run     run a job
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vlumi/copal/pangaea/compile"
	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/jcl"
	"github.com/vlumi/copal/pangaea/runner"
)

const usage = `usage:
  pangaea build [-src dir] [-copybooks dir] [-out dir]
  pangaea load -datasets dir file.txt...
  pangaea unload -datasets dir outdir
  pangaea run -datasets dir -sysout dir [-programs dir] job.jcl
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	commands := map[string]func([]string) error{
		"build":  build,
		"load":   load,
		"unload": unload,
		"run":    run,
	}
	command, ok := commands[os.Args[1]]
	if !ok {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	if err := command(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "pangaea:", err)
		os.Exit(1)
	}
}

func build(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	src := flags.String("src", "pangaea/programs", "directory of the COBOL programs")
	copybooks := flags.String("copybooks", "pangaea/copybooks", "directory of the copybooks")
	out := flags.String("out", "pangaea/build", "directory for the executables")
	_ = flags.Parse(args)
	names, err := compile.Programs(*src, *copybooks, *out)
	if err != nil {
		return err
	}
	fmt.Println("compiled", strings.Join(names, " "))
	return nil
}

func load(args []string) error {
	flags := flag.NewFlagSet("load", flag.ExitOnError)
	dir := flags.String("datasets", "", "dataset directory")
	_ = flags.Parse(args)
	store, err := dataset.Open(*dir)
	if err != nil {
		return err
	}
	for _, path := range flags.Args() {
		if err := store.LoadFile(path); err != nil {
			return err
		}
	}
	return nil
}

func unload(args []string) error {
	flags := flag.NewFlagSet("unload", flag.ExitOnError)
	dir := flags.String("datasets", "", "dataset directory")
	_ = flags.Parse(args)
	if flags.NArg() != 1 {
		return fmt.Errorf("unload needs an output directory")
	}
	store, err := dataset.Open(*dir)
	if err != nil {
		return err
	}
	return store.UnloadAll(flags.Arg(0))
}

func run(args []string) error {
	flags := flag.NewFlagSet("run", flag.ExitOnError)
	dir := flags.String("datasets", "", "dataset directory")
	sysout := flags.String("sysout", "", "directory for the job log and reports")
	programs := flags.String("programs", "pangaea/build", "directory of the compiled programs")
	_ = flags.Parse(args)
	if flags.NArg() != 1 {
		return fmt.Errorf("run needs one job file")
	}
	job, err := parseJob(flags.Arg(0))
	if err != nil {
		return err
	}
	store, err := dataset.Open(*dir)
	if err != nil {
		return err
	}
	r := runner.Runner{Programs: *programs, Datasets: store, Sysout: *sysout}
	results, err := r.Run(job)
	for _, result := range results {
		fmt.Println(result)
	}
	return err
}

func parseJob(path string) (*jcl.Job, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	job, err := jcl.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return job, nil
}
