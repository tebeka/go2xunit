// Parse "gocheck -vv" output
package gocheck

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"bitbucket.org/tebeka/go2xunit/types"
)

const (
	// START: mmath_test.go:16: MySuite.TestAdd
	startRE = "START: [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)"
	// PASS: mmath_test.go:16: MySuite.TestAdd	0.000s
	// FAIL: mmath_test.go:35: MySuite.TestDiv
	endRE = "(PASS|FAIL): [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)([[:space:]]+([0-9]+.[0-9]+))?"
)

func map2arr(m map[string]*types.Suite) []*types.Suite {
	arr := make([]*types.Suite, 0, len(m))
	for _, suite := range m {
		/* FIXME:
		suite.Status =
		suite.Time =
		*/
		arr = append(arr, suite)
	}

	return arr
}

// Parse parses output of "go test -gocheck.vv", returns a list of tests
// See data/gocheck.out for an example
func Parse(rd io.Reader) ([]*types.Suite, error) {
	find_start := regexp.MustCompile(startRE).FindStringSubmatch
	find_end := regexp.MustCompile(endRE).FindStringSubmatch

	scanner := bufio.NewScanner(rd)
	var test *types.Test
	var suites = make(map[string]*types.Suite)
	var suiteName string
	var out []string

	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()
		tokens := find_start(line)
		if len(tokens) > 0 {
			if test != nil {
				return nil, fmt.Errorf("%d: start in middle\n", lnum)
			}
			suiteName = tokens[1]
			test = &types.Test{Name: tokens[2]}
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

			suite, ok := suites[suiteName]
			if !ok {
				suite = &types.Suite{Name: suiteName}
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
