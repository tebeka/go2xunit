// Command line parsing
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var args struct {
	noFail      bool
	version     bool
	format      string
	suitePrefix string
	failRace    bool
	output      string
}

var (
	// TODO: Populate
	outFormats = map[string]bool{
		"junit":  true,
		"xunit":  true,
		"bamboo": true,
	}
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [input]\nOptions:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.BoolVar(&args.noFail, "no-fail", false, "don't fail if tests failed")
	flag.BoolVar(&args.version, "version", false, "print version and exit")
	flag.StringVar(&args.format, "format", "junit", "output format: junit, xunit, bamboo")
	flag.BoolVar(&args.failRace, "fail-race", false, "mark test as failing if it exposes a data race")
	flag.StringVar(&args.suitePrefix, "suite-prefix", "", "prefix to include before all suite names")
	flag.StringVar(&args.output, "output", "", "output file")
}

// parseArgs parses and validates command line arguments
func parseArgs() error {
	flag.Parse()

	if flag.NArg() > 1 {
		return fmt.Errorf("too many arguments for %s", os.Args[0])
	}

	return nil
}
