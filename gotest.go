// Parse "gotest -v" output
package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Since mucking with local package is a PITA, just prefix everything with gt_

// parseEnd parses "end of test" line and returns (name, time, error)
func gt_parseEnd(prefix, line string) (string, string, error) {
	// "end of test" regexp for name and time, examples:
	// --- PASS: TestSub (0.00 seconds)
	// --- FAIL: TestSubFail (0.00 seconds)
	var endRegexp *regexp.Regexp = regexp.MustCompile(`([^ ]+) \((\d+\.\d+)`)

	matches := endRegexp.FindStringSubmatch(line[len(prefix):])

	if len(matches) == 0 {
		return "", "", fmt.Errorf("can't parse %s", line)
	}

	return matches[1], matches[2], nil
}


// gt_parseEndTest parses "end of test file" line and returns (status, name, time, error)
func gt_parseEndTest(line string) (string, string, string, error) {
	// "end of tested file" regexp for parsing package & file name
	// ok  	teky/cointreau/gs1/deliver	0.015s
	// FAIL	teky/cointreau/gs1/deliver	0.010s
	var endTestRegexp *regexp.Regexp = regexp.MustCompile(`^(ok  |FAIL)\t([^ ]+)\t(\d+\.\d+)s$`)

	matches := endTestRegexp.FindStringSubmatch(line)

	if len(matches) == 0 {
		return "", "", "", fmt.Errorf("can't parse %s", line)
	}

	return matches[1], matches[2], matches[3], nil
}

// gt_Parse parses output of "go test -v", returns a list of tests
// See data/gotest.out for an example
func gt_Parse (rd io.Reader) ([]*Suite, error) {

	startPrefix := "=== RUN "
	passPrefix := "--- PASS: "
	failPrefix := "--- FAIL: "

	suites := []*Suite{}
	var test *Test = nil
	var suite *Suite = nil

	nextTest := func() {
		// We are switching to the next test, store the current one.
		if suite == nil {
			suite = &Suite{}
			suite.Tests = make([]*Test, 0, 1)
		}
		if test == nil {
			return
		}
		suite.Tests = append(suite.Tests, test)
		test = nil
	}
	nextSuite := func() {
		// We are switching to the next suite, store the current one.
		if suite == nil {
			return
		}

		suites = append(suites, suite)
		suite = nil
	}

	reader := bufio.NewReader(rd)
	for {
		buf, _, err := reader.ReadLine()

		switch err {
		case io.EOF:
			if suite != nil || test != nil {
				// if suite or test in progress EOF is an unexpected EOF
				return nil, fmt.Errorf("Unexpected EOF")
			}
			return suites, nil
		case nil:
			// nil is OK

		default: // Error other than io.EOF
			return nil, err
		}

		line := string(buf)
		switch {
		case strings.HasPrefix(line, startPrefix):
		case strings.HasPrefix(line, failPrefix):
			nextTest()

			// Extract the test name and the duration:
			name, time, err := gt_parseEnd(failPrefix, line)
			if err != nil {
				return nil, err
			}

			test = &Test{
				Name:   name,
				Time:   time,
				Failed: true,
			}

		case strings.HasPrefix(line, passPrefix):
			nextTest()
			// Extract the test name and the duration:
			name, time, err := gt_parseEnd(passPrefix, line)
			if err != nil {
				return nil, err
			}
			// Create the test structure and store it.
			suite.Tests = append(suite.Tests, &Test{
				Name:   name,
				Time:   time,
				Failed: false,
			})
			test = nil
		case line == "FAIL":
			nextTest()

		case strings.HasPrefix(line, "ok  \t") || strings.HasPrefix(line, "FAIL\t"):
			// End of suite, read data
			status, name, time, err := gt_parseEndTest(line)
			if err != nil {
				return nil, err
			}
			suite.Name = name
			suite.Count = len(suite.Tests)
			suite.Failed = numFailures(suite.Tests)
			suite.Time = time
			suite.Status = status
			nextSuite()
		default:
			if test != nil { // test != nil marks we're in the middle of a test
				test.Message += line + "\n"
			}
		}
	}

	// If we're here, it's an error
	return nil, fmt.Errorf("Error parsing")
}
