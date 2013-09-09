// Parse "gotest -v" output
package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// parseEnd parses "end of test" line and returns (name, time, error)
func parseEnd(prefix, line string) (string, string, error) {
	// "end of test" regexp for name and time, examples:
	// PASS: mmath_test.go:16: MySuite.TestAdd	0.000s
	// FAIL: mmath_test.go:35: MySuite.TestDiv
	var endRegexp *regexp.Regexp = regexp.MustCompile(
		`(PASS|FAIL): [^ ]+:\d+: [^ ]+ \((\d+\.\d+)?`)

	matches := endRegexp.FindStringSubmatch(line[len(prefix):])

	if len(matches) == 0 {
		return "", "", fmt.Errorf("can't parse %s", line)
	}

	return matches[1], matches[2], nil
}

// parseOutput parses output of "go test -gocheck.vv", returns a list of tests
// See data/gocheck.out for an example
func parseGocheckOutput(rd io.Reader) ([]*Suite, error) {
	startPrefix := "START: "
	passPrefix := "PASS: "
	failPrefix := "FAIL: "

	var suite *Suite = &Suite{}
	suites := []*Suite{suite} // FIXME: Just one suite in gocheck
	var test *Test = nil
	inTest := false
	output := []string{}

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
			inTest = true
		case strings.HasPrefix(line, passPrefix):
		case strings.HasPrefix(line, failPrefix):
		default:
			if inTest {
				output = append(output, line)
			}
		}

	}
	return nil, nil
}
