/*
crater converts a manifest from gather into a RO-CRATE.

 1. read gather manifest.
 2. download all files and resources.
 3. output a ro-crate JSON.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ross-spencer/zenodocfl/internal/logformatter"
)

var (
	crate string
	vers  bool

	// app constants.
	version = "dev-0.0.0"
	commit  = "000000000000000000000000000000000baddeed"
	date    = "1970-01-01T00:00:01Z"
)

var agent string = fmt.Sprintf("INK-crater/%s", version)

// initFlags initializes the flags we use with this app.
func initFlags() {
	flag.StringVar(&crate, "crate", "", "manifest to convert to RO-CRATE")
	flag.BoolVar(&vers, "version", false, "Return version")
}

func main() {

	logformatter.Set("gather", true)

	initFlags()
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stderr, "%s (%s) commit: %s date: %s\n", agent, version, commit, date)
		os.Exit(0)
	} else if flag.NFlag() < 1 {
		fmt.Fprintln(os.Stderr, "Usage:  ")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-crate]  STRING")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-version] ")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Output: [STRING] {result JSON}")
		fmt.Fprintf(os.Stderr, "Output: [STRING] {version: '%s'}\n\n", agent)
		flag.Usage()
		os.Exit(0)
	}

	if vers {
		fmt.Println(agent)
		return
	}

	if crate == "" {
		log.Println("no arguments provided")
		os.Exit(1)
	}

	// we should never reach here.
	log.Println("no arguments provided")

}
