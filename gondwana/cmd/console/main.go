// Command console serves the life insurer's console: a web window onto
// the running Gondwana Life Insurance for its operators.
//
//	console [-insurer dir] [-programs dir] [-addr host:port]
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vlumi/copal/gondwana/console"
	"github.com/vlumi/copal/gondwana/insurer"
)

func main() {
	if err := serve(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "console:", err)
		os.Exit(1)
	}
}

func serve(args []string) error {
	flags := flag.NewFlagSet("console", flag.ExitOnError)
	in := insurer.Insurer{Scripts: "gondwana/scripts", Copybooks: "gondwana/copybooks"}
	flags.StringVar(&in.Dir, "insurer", "gondwana/run/insurer", "directory of the running insurer")
	flags.StringVar(&in.Programs, "programs", "gondwana/build", "directory of the compiled programs")
	addr := flags.String("addr", "127.0.0.1:2525", "address to serve on")
	_ = flags.Parse(args)
	c, err := console.New(in)
	if err != nil {
		return err
	}
	fmt.Printf("the life insurer's console is at http://%s/\n", *addr)
	server := &http.Server{Addr: *addr, Handler: c.Handler(), ReadHeaderTimeout: 10 * time.Second}
	return server.ListenAndServe()
}
