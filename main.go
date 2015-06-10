package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const version = "1.2.0"

var failOnRace = false

func main() {
	inputFile := flag.String("input", "", "input file (default to stdin)")
	outputFile := flag.String("output", "", "output file (default to stdout)")
	fail := flag.Bool("fail", false, "fail (non zero exit) if any test failed")
	showVersion := flag.Bool("version", false, "print version and exit")
	bamboo := flag.Bool("bamboo", false, "xml compatible with Atlassian's Bamboo")
	xunitnet := flag.Bool("xunitnet", false, "xml compatible with xunit.net")
	is_gocheck := flag.Bool("gocheck", false, "parse gocheck output")
	flag.BoolVar(&failOnRace, "fail-on-race", false, "mark test as failing if it exposes a data race")
	flag.Parse()

	if *showVersion {
		fmt.Println("tischda/go2xunit", version)
		return
	}

	// No time ... prefix for error messages
	log.SetFlags(0)

	if flag.NArg() > 0 {
		log.Fatalf("error: %s does not take parameters (did you mean -input?)", os.Args[0])
	}

	if *bamboo && *xunitnet {
		log.Fatalf("error: -bamboo and -xunitnet are mutually exclusive")
	}

	input, output := getIO(*inputFile, *outputFile)

	var parse func(rd io.Reader) ([]*Suite, error)

	if *is_gocheck {
		parse = gc_Parse
	} else {
		parse = gt_Parse
	}

	suites, err := parse(input)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	if len(suites) == 0 {
		log.Fatalf("error: no tests found")
	}

	xmlTemplate := xunitTemplate
	if *xunitnet {
		xmlTemplate = xunitNetTemplate
	} else if *bamboo || (len(suites) > 1) {
		xmlTemplate = multiTemplate
	}

	writeXML(suites, output, xmlTemplate)
	if *fail && hasFailures(suites) {
		os.Exit(1)
	}
}

func hasFailures(suites []*Suite) bool {
	for _, suite := range suites {
		if suite.NumFailed() > 0 {
			return true
		}
	}
	return false
}
