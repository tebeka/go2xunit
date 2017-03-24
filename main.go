package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/tebeka/go2xunit/lib"
)

const (
	// Version is the current version
	Version = "1.4.4"
)

// getInput return input io.File from file name, if file name is - it will
// return os.Stdin
func getInput(filename string) (*os.File, error) {
	if filename == "-" || filename == "" {
		return os.Stdin, nil
	}

	return os.Open(filename)
}

// getInput return output io.File from file name, if file name is - it will
// return os.Stdout
func getOutput(filename string) (*os.File, error) {
	if filename == "-" || filename == "" {
		return os.Stdout, nil
	}

	return os.Create(filename)
}

// getIO returns input and output streams from file names
func getIO(inFile, outFile string) (*os.File, io.Writer, error) {
	input, err := getInput(inFile)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open %s for reading: %s", inFile, err)
	}

	output, err := getOutput(outFile)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open %s for writing: %s", outFile, err)
	}

	return input, output, nil
}

func getTestTime(in *os.File) time.Time {
	stat, err := in.Stat()
	if err != nil {
		return time.Now()
	}
	return stat.ModTime()
}

func main() {
	flag.Parse()
	if args.showVersion {
		fmt.Printf("go2xunit %s\n", Version)
		os.Exit(0)
	}

	// No time ... prefix for error messages
	log.SetFlags(0)

	if err := validateArgs(); err != nil {
		log.Fatalf("error: %s", err)
	}

	input, output, err := getIO(args.inFile, args.outFile)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	parser, err := lib.GetParser(args.inType, input)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// We'd like the test time to be the time of the generated file
	testTime := getTestTime(input)

	var suites lib.Suites
	for parser.Scan() {
		suites = append(suites, parser.Suite())
	}
	if err := parser.Err(); err != nil {
		log.Fatalf("error: %s", err)
	}

	if len(suites) == 0 {
		log.Fatalf("error: no tests found")
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
}
