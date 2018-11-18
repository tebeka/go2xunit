package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

const (
	// Version is the current version
	Version = "1.4.8"
)

func inFile() (*os.File, error) {
	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		return os.Stdin, nil
	}

	return os.Open(flag.Arg(0))
}

func outFile(path string) (*os.File, error) {
	if path == "" || path == "-" {
		return os.Stdout, nil
	}

	return os.Create(flag.Arg(1))
}

// getInput return input io.File from file name, if file name is - it will
// return os.Stdin
func main() {
	if args.version {
		fmt.Printf("%s %s\n", Version, path.Base(os.Args[0]))
		os.Exit(0)
	}

	// No time ... prefix for error messages
	log.SetFlags(0)

	if err := parseArgs(); err != nil {
		log.Fatalf("error: %s", err)
	}

	input, err := inFile()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	/*
		output, err := outFile(args.output)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
	*/

	Parse(input)

	/*
		var parse func(rd io.Reader, suiteName string) (lib.Suites, error)

		if args.isGocheck {
			parse = lib.ParseGocheck
		} else {
			parse = lib.ParseGotest
		}

		suites, err := parse(input, args.suitePrefix)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		if len(suites) == 0 {
			log.Fatalf("error: no tests found")
			os.Exit(1)
		}

		xmlTemplate := lib.XUnitTemplate
		if args.xunitnetOut {
			xmlTemplate = lib.XUnitNetTemplate
		} else if args.bambooOut || (len(suites) > 1) {
			xmlTemplate = lib.XMLMultiTemplate
		}

		lib.WriteXML(suites, output, xmlTemplate, testTime)
		if args.fail && suites.HasFailures() {
			os.Exit(1)
		}
	*/
}
