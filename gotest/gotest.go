// Parse "gotest -v" output
package gotest

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"bitbucket.org/tebeka/go2xunit/types"
)

const (
	// === RUN TestAdd
	startRE = "^=== RUN ([a-zA-Z_][[:word:]]*)"

	// --- PASS: TestSub (0.00 seconds)
	// --- FAIL: TestSubFail (0.00 seconds)
	endRE = "^--- (PASS|FAIL): ([a-zA-Z_][[:word:]]*) \\((\\d+.\\d+)"

	// FAIL	_/home/miki/Projects/goroot/src/xunit	0.004s
	// ok  	_/home/miki/Projects/goroot/src/anotherTest	0.000s
	suiteRE = "^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(\\d+.\\d+)"
)

func Parse(rd io.Reader) ([]*types.Suite, error) {
	find_start := regexp.MustCompile(startRE).FindStringSubmatch
	find_end := regexp.MustCompile(endRE).FindStringSubmatch
	find_suite := regexp.MustCompile(suiteRE).FindStringSubmatch
	is_exit := regexp.MustCompile("^exit status -?\\d+").MatchString

	suites := []*types.Suite{}
	var curTest *types.Test
	var curSuite *types.Suite
	var out []string

	scanner := bufio.NewScanner(rd)
	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()

		tokens := find_start(line)
		if tokens != nil {
			if curTest != nil {
				return nil, fmt.Errorf("%d: test in middle of other", lnum)
			}
			curTest = &types.Test{
				Name: tokens[1],
			}
			if len(out) > 0 {
				message := strings.Join(out, "\n")
				if (curSuite == nil) {
					return nil, fmt.Errorf("orphan output: %s", message)
				}
				curSuite.Tests[len(curSuite.Tests)-1].Message = message
			}
			out = []string{}
			continue
		}

		tokens = find_end(line)
		if tokens != nil {
			if curTest == nil {
				return nil, fmt.Errorf("%d: orphan end test", lnum)
			}
			if tokens[2] != curTest.Name {
				return nil, fmt.Errorf("%d: name mismatch", lnum)
			}

			curTest.Failed = (tokens[1] == "FAIL")
			curTest.Time = tokens[3]
			curTest.Message = strings.Join(out, "\n")
			if curSuite == nil {
				curSuite = &types.Suite{}
			}
			curSuite.Tests = append(curSuite.Tests, curTest)
			curTest = nil
			continue
		}

		tokens = find_suite(line)
		if tokens != nil {
			if curSuite == nil {
				return nil, fmt.Errorf("%d: orphan end suite", lnum)
			}
			curSuite.Name = tokens[2]
			curSuite.Time = tokens[3]
			suites = append(suites, curSuite)
			curSuite = nil

			continue
		}

		if is_exit(line) || (line == "FAIL") {
			continue
		}

		if curSuite == nil {
			return nil, fmt.Errorf("%d: orphan line", lnum)
		}

		out = append(out, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}
