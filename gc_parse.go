package main

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// gc_Parse parses output of "go test -gocheck.vv", returns a list of tests
// See data/gocheck.out for an example
func gc_Parse(rd io.Reader) ([]*Suite, error) {
	find_end := regexp.MustCompile(gc_endRE).FindStringSubmatch

	scanner := bufio.NewScanner(rd)
	var test *Test
	var suites = make(map[string]*Suite)
	var suiteName string
	var out []string

	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()
		tokens := find_end(line)
		if len(tokens) > 0 {
			suiteName = tokens[2]
			test = &Test{Name: tokens[3]}
			test.Message = strings.Join(out, "\n")
			test.Time = tokens[4]
			test.Failed = (tokens[1] == "FAIL")
			test.Passed = (tokens[1] == "PASS")
			test.Skipped = (tokens[1] == "SKIP")

			suite, ok := suites[suiteName]
			if !ok {
				suite = &Suite{Name: suiteName}
			}
			suite.Tests = append(suite.Tests, test)
			suites[suiteName] = suite

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

	return map2arr(suites), nil
}

func map2arr(m map[string]*Suite) []*Suite {
	arr := make([]*Suite, 0, len(m))
	for _, suite := range m {
		/* FIXME:
		suite.Status =
		suite.Time =
		*/
		arr = append(arr, suite)
	}

	return arr
}
