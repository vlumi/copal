// Command pangaea builds and runs the First Bank of Pangaea.
//
//	pangaea build   compile the COBOL programs
//	pangaea load    load text files, one record per line, as datasets
//	pangaea unload  write every cataloged dataset as text
//	pangaea run     run a job
//	pangaea reset   reset the running bank to a generated starting point
//	pangaea days    run the running bank's next business days
//	pangaea online  serve the running bank's online day until interrupted
//	pangaea simulate send card traffic to the online day until interrupted
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/vlumi/copal/bedrock/jcl"
	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/pangaea/compile"
	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/online"
	"github.com/vlumi/copal/pangaea/records"
	"github.com/vlumi/copal/pangaea/runner"
	"github.com/vlumi/copal/pangaea/simulate"
)

const usage = `usage:
  pangaea build [-src dir] [-copybooks dir] [-out dir]
  pangaea load -datasets dir file.txt...
  pangaea unload -datasets dir outdir
  pangaea run -datasets dir -sysout dir [-programs dir] job.jcl
  pangaea reset [-bank dir] [-seed n] [-accounts n]
  pangaea days [-bank dir] [-n days]
  pangaea online [-bank dir] [-addr host:port]
  pangaea simulate [-bank dir] [-api url] [-rate n] [-seed n]
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	commands := map[string]func([]string) error{
		"build":    build,
		"load":     load,
		"unload":   unload,
		"run":      run,
		"reset":    reset,
		"days":     days,
		"online":   serveOnline,
		"simulate": simulateTraffic,
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
	job, err := jcl.ParseFile(flags.Arg(0))
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

func bankFlags(name string) (*flag.FlagSet, *bank.Bank) {
	flags := flag.NewFlagSet(name, flag.ExitOnError)
	b := &bank.Bank{Copybooks: "pangaea/copybooks", Jobs: "pangaea/jobs"}
	flags.StringVar(&b.Dir, "bank", "pangaea/run/bank", "directory of the running bank")
	flags.StringVar(&b.Programs, "programs", "pangaea/build", "directory of the compiled programs")
	return flags, b
}

func reset(args []string) error {
	flags, b := bankFlags("reset")
	seed := flags.Uint64("seed", 1, "seed the bank is generated from")
	accounts := flags.Int("accounts", 3000, "number of accounts, about")
	_ = flags.Parse(args)
	if err := b.Reset(*seed, *accounts); err != nil {
		return err
	}
	fmt.Printf("bank reset in %s: seed %d, about %d accounts, first business day %s\n", b.Dir, *seed, *accounts, bank.StartDate)
	return nil
}

func days(args []string) error {
	flags, b := bankFlags("days")
	n := flags.Int("n", 1, "number of business days to run")
	_ = flags.Parse(args)
	return b.RunDays(*n, func(d bank.Day) {
		rc := 0
		for _, r := range d.Results {
			rc = max(rc, r.RC)
		}
		monthEnd := ""
		if len(d.MonthEnd) > 0 {
			monthEnd = "  and the month-end job"
		}
		fmt.Printf("%s  %5d transactions  %4d rejects  RC=%04d%s\n", d.Date, d.Transactions, d.Rejects, rc, monthEnd)
	})
}

// serveOnline opens the business day's online day and serves its API,
// until an interrupt closes the day, extracting its holds for the batch.
func serveOnline(args []string) error {
	flags, b := bankFlags("online")
	addr := flags.String("addr", "127.0.0.1:8081", "address to serve on")
	_ = flags.Parse(args)
	bankLayouts, err := records.Load(b.Copybooks)
	if err != nil {
		return err
	}
	s := &online.Server{Bank: *b, Layouts: bankLayouts, Clock: time.Now}
	if err := s.Open(); err != nil {
		return err
	}
	fmt.Printf("the online day is open at http://%s/; interrupt to close it\n", *addr)
	return serveUntilInterrupted(*addr, s.Handler(), func() error {
		results, err := s.Close()
		for _, r := range results {
			fmt.Println(r)
		}
		return err
	})
}

// simulateTraffic sends made-up card traffic to the online channel until
// interrupted, printing every request and its response.
func simulateTraffic(args []string) error {
	flags, b := bankFlags("simulate")
	api := flags.String("api", "http://127.0.0.1:8080/api", "the online channel's API")
	rate := flags.Float64("rate", 5, "requests a second")
	seed := flags.Uint64("seed", 1, "seed the traffic is made up from")
	_ = flags.Parse(args)
	if *rate <= 0 {
		return fmt.Errorf("rate %v: more than none a second", *rate)
	}
	bankLayouts, err := records.Load(b.Copybooks)
	if err != nil {
		return err
	}
	store, err := b.Datasets()
	if err != nil {
		return err
	}
	s, err := simulate.New(*api, store, bankLayouts, *seed)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	fmt.Printf("sending %v requests a second to %s; interrupt to stop\n", *rate, *api)
	return s.Run(ctx, *rate, func(r simulate.Request) { fmt.Println(r) })
}

// serveUntilInterrupted serves until an interrupt, and then calls done.
func serveUntilInterrupted(addr string, handler http.Handler, done func() error) error {
	server := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	interrupted := make(chan os.Signal, 1)
	signal.Notify(interrupted, os.Interrupt)
	go func() {
		<-interrupted
		_ = server.Close()
	}()
	err := server.ListenAndServe()
	if doneErr := done(); err == nil || errors.Is(err, http.ErrServerClosed) {
		err = doneErr
	}
	return err
}
