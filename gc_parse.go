package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// gc_Parse parses output of "go test -gocheck.vv", returns a list of tests
// See data/gocheck.out for an example
func gc_Parse(rd io.Reader) ([]*Suite, error) {
	find_start := regexp.MustCompile(gc_startRE).FindStringSubmatch
	find_end := regexp.MustCompile(gc_endRE).FindStringSubmatch

	scanner := bufio.NewScanner(rd)
	var test *Test
	var suites = make([]*Suite, 0)
	var suiteName string
	var currentSuite *Suite
	var out []string

	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()
		tokens := find_start(line)
		if len(tokens) > 0 {
			if test != nil {
				return nil, fmt.Errorf("%d: start in middle\n", lnum)
			}
			suiteName = tokens[1]
			test = &Test{Name: tokens[2]}
			out = []string{}
			continue
		}

		tokens = find_end(line)
		if len(tokens) > 0 {
			if test == nil {
				return nil, fmt.Errorf("%d: orphan end", lnum)
			}
			if (tokens[2] != suiteName) || (tokens[3] != test.Name) {
				return nil, fmt.Errorf("%d: suite/name mismatch", lnum)
			}
			test.Message = strings.Join(out, "\n")
			test.Time = tokens[4]
			test.Failed = (tokens[1] == "FAIL")
			test.Passed = (tokens[1] == "PASS")
			test.Skipped = (tokens[1] == "SKIP")

			if currentSuite == nil || currentSuite.Name != suiteName {
				currentSuite = &Suite{Name: suiteName}
				suites = append(suites, currentSuite)
			}
			currentSuite.Tests = append(currentSuite.Tests, test)

			test = nil
			suiteName = ""
			out = []string{}

			continue
		}

		if test != nil {
			out = append(out, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}

/* FIXME:
suite.Status =
suite.Time =
*/
