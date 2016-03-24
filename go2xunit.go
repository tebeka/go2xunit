package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	Version = "1.2.4"
)

// Test data structure
type Test struct {
	Name, Time, Message string
	Failed              bool
	Skipped             bool
	Passed              bool
}

// Suite of tests (found in some unit testing frameworks)
type Suite struct {
	Name   string
	Time   string
	Status string
	Tests  []*Test
}

// SuiteStack is a stack of test suites
type SuiteStack struct {
	nodes []*Suite
	count int
}

// Push adds a node to the stack.
func (s *SuiteStack) Push(n *Suite) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *SuiteStack) Pop() *Suite {
	if s.count == 0 {
		return nil
	}
	s.count--
	return s.nodes[s.count]
}

// NumFailed return number of failed tests in suite
func (suite *Suite) NumFailed() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Failed {
			count++
		}
	}

	return count
}

// NumSkipped return number of skipped tests in suite
func (suite *Suite) NumSkipped() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Skipped {
			count++
		}
	}

	return count
}

// NumPassed return number of passed tests in the suite
func (suite *Suite) NumPassed() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Passed {
			count++
		}
	}

	return count
}

// Count return the number of tests in the suite
func (suite *Suite) Count() int {
	return len(suite.Tests)
}

// hasFailures return true is there's at least one failing test in the suite
func hasFailures(suites []*Suite) bool {
	for _, suite := range suites {
		if suite.NumFailed() > 0 {
			return true
		}
	}
	return false
}

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

func main() {
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

	// We'd like the test time to be the time of the generated file
	var testTime time.Time
	stat, err := input.Stat()
	if err != nil {
		testTime = time.Now()
	} else {
		testTime = stat.ModTime()
	}

	var parse func(rd io.Reader) ([]*Suite, error)

	if args.isGocheck {
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
		os.Exit(1)
	}

	xmlTemplate := xunitTemplate
	if args.xunitnetOut {
		xmlTemplate = xunitNetTemplate
	} else if args.bambooOut || (len(suites) > 1) {
		xmlTemplate = multiTemplate
	}

	writeXML(suites, output, xmlTemplate, testTime)
	if args.fail && hasFailures(suites) {
		os.Exit(1)
	}
}
