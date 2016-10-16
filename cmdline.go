// Command line parsing
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tebeka/go2xunit/lib"
)

var args struct {
	inFile      string
	outFile     string
	fail        bool
	showVersion bool
	bambooOut   bool
	xunitnetOut bool
	isGocheck   bool
	suitePrefix string
}

func init() {
	flag.StringVar(&args.inFile, "input", "", "input file (default to stdin)")
	flag.StringVar(&args.outFile, "output", "", "output file (default to stdout)")
	flag.BoolVar(&args.fail, "fail", false, "fail (non zero exit) if any test failed")
	flag.BoolVar(&args.showVersion, "version", false, "print version and exit")
	flag.BoolVar(&args.bambooOut, "bamboo", false,
		"xml compatible with Atlassian's Bamboo")
	flag.BoolVar(&args.xunitnetOut, "xunitnet", false, "xml compatible with xunit.net")
	flag.BoolVar(&args.isGocheck, "gocheck", false, "parse gocheck output")
	flag.BoolVar(&lib.Options.FailOnRace, "fail-on-race", false,
		"mark test as failing if it exposes a data race")
	flag.StringVar(&args.suitePrefix, "suite-name-prefix", "",
		"prefix to include before all suite names")

	flag.Parse()
}

// validateArgs validates command line arguments
func validateArgs() error {
	if flag.NArg() > 0 {
		return fmt.Errorf("%s does not take parameters (did you mean -input?)", os.Args[0])
	}

	if args.bambooOut && args.xunitnetOut {
		return fmt.Errorf("-bamboo and -xunitnet are mutually exclusive")
	}

	return nil
}
