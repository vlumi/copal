// Command console serves the card processor's console: a web window onto
// the running Laurasia Card Services for its operators.
//
//	console [-processor dir] [-programs dir] [-addr host:port]
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vlumi/copal/laurasia/console"
	"github.com/vlumi/copal/laurasia/processor"
)

func main() {
	if err := serve(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "console:", err)
		os.Exit(1)
	}
}

func serve(args []string) error {
	flags := flag.NewFlagSet("console", flag.ExitOnError)
	p := processor.Processor{Scripts: "laurasia/scripts", Copybooks: "laurasia/copybooks"}
	flags.StringVar(&p.Dir, "processor", "laurasia/run/processor", "directory of the running processor")
	flags.StringVar(&p.Programs, "programs", "laurasia/build", "directory of the compiled programs")
	addr := flags.String("addr", "127.0.0.1:4649", "address to serve on")
	_ = flags.Parse(args)
	c, err := console.New(p)
	if err != nil {
		return err
	}
	fmt.Printf("the card processor's console is at http://%s/\n", *addr)
	server := &http.Server{Addr: *addr, Handler: c.Handler(), ReadHeaderTimeout: 10 * time.Second}
	return server.ListenAndServe()
}
