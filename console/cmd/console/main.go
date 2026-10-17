// Command console serves the bank console: a web window onto the running
// First Bank of Pangaea, its online channel, and the survey of the bank.
//
//	console [-bank dir] [-addr host:port] [-blueprints dir]
package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/vlumi/copal/console"
	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/survey/blueprint"
)

func main() {
	if err := serve(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "console:", err)
		os.Exit(1)
	}
}

func serve(args []string) error {
	flags := flag.NewFlagSet("console", flag.ExitOnError)
	b := bank.Bank{Copybooks: "pangaea/copybooks", Jobs: "pangaea/jobs"}
	flags.StringVar(&b.Dir, "bank", "pangaea/run/bank", "directory of the running bank")
	flags.StringVar(&b.Programs, "programs", "pangaea/build", "directory of the compiled programs")
	addr := flags.String("addr", "localhost:8080", "address to serve on")
	src := flags.String("src", "pangaea/programs", "directory of the COBOL programs")
	blueprints := flags.String("blueprints", "blueprints/pangaea", "directory of the bank's blueprints")
	_ = flags.Parse(args)
	programs, err := filepath.Glob(filepath.Join(*src, "*.cbl"))
	if err != nil {
		return err
	}
	system, err := blueprint.Analyze(append(programs, b.JobFiles()...), b.Copybooks)
	if err != nil {
		return err
	}
	c, err := console.New(b, system, *blueprints)
	if err != nil {
		return err
	}
	if err := c.Open(); err != nil {
		return err
	}
	fmt.Printf("the bank console is at http://%s/, the online channel at http://%s/api/; interrupt to close it\n", *addr, *addr)
	server := &http.Server{Addr: *addr, Handler: c.Handler(), ReadHeaderTimeout: 10 * time.Second}
	interrupted := make(chan os.Signal, 1)
	signal.Notify(interrupted, os.Interrupt)
	go func() {
		<-interrupted
		_ = server.Close()
	}()
	err = server.ListenAndServe()
	if closeErr := c.Close(); err == nil || errors.Is(err, http.ErrServerClosed) {
		err = closeErr
	}
	return err
}
