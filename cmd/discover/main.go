// discover provides node discovery on the command line.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	discover "github.com/hashicorp/go-discover"
)

func usage() {
	fmt.Println("Usage: discover addrs key=val key=val ...")
	fmt.Println(discover.HelpDiscoverAddrs)
}
func main() {
	var quiet bool
	var help bool
	flag.BoolVar(&quiet, "q", false, "no verbose output")
	flag.BoolVar(&help, "h", false, "print help")
	flag.Parse()

	args := flag.Args()
	if help || len(args) == 0 || args[0] != "addrs" {
		usage()
		os.Exit(0)
	}
	args = args[1:]

	var w io.Writer = os.Stderr
	if quiet {
		w = ioutil.Discard
	}
	l := log.New(w, "", 0)

	addrs, err := discover.Addrs(strings.Join(args, " "), l)
	if err != nil {
		l.Fatal(err)
	}
	fmt.Println(strings.Join(addrs, " "))
}
