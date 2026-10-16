// Command copal reads a COBOL system and writes its blueprints.
//
//	copal -copybooks dir -out dir [-title name] program.cbl...
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vlumi/copal/blueprint"
)

func main() {
	copybooks := flag.String("copybooks", "", "directory of the copybooks")
	out := flag.String("out", "", "directory to write the blueprints to")
	title := flag.String("title", "", "name of the system, for the index")
	flag.Parse()
	if *copybooks == "" || *out == "" || flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: copal -copybooks dir -out dir [-title name] program.cbl...")
		os.Exit(2)
	}
	system, err := blueprint.Analyze(flag.Args(), *copybooks)
	if err == nil {
		err = system.Write(*out, *title)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "copal:", err)
		os.Exit(1)
	}
}
