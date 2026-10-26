// Command console serves the bank console: a web window onto the running
// First Bank of Pangaea and its online channel, for its operators.
//
//	console [-bank dir] [-programs dir] [-addr host:port]
package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/vlumi/copal/pangaea/bank"
	"github.com/vlumi/copal/pangaea/console"
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
	addr := flags.String("addr", "127.0.0.1:4510", "address to serve on")
	_ = flags.Parse(args)
	c, err := console.New(b)
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
