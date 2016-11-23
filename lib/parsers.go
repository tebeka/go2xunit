// Package lib is exposing parsers and output generation
package lib

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	// gotest regular expressions

	// === RUN TestAdd
	gtStartRE = "^=== RUN:?[[:space:]]+([a-zA-Z_][^[:space:]]*)"

	// --- PASS: TestSub (0.00 seconds)
	// --- FAIL: TestSubFail (0.00 seconds)
	// --- SKIP: TestSubSkip (0.00 seconds)
	gtEndRE = "--- (PASS|FAIL|SKIP):[[:space:]]+([a-zA-Z_][^[:space:]]*) \\((-?\\d+(.\\d+)?)"

	// FAIL	_/home/miki/Projects/goroot/src/xunit	0.004s
	// ok  	_/home/miki/Projects/goroot/src/anotherTest	0.000s
	gtSuiteRE = "^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(-?\\d+.\\d+)"

	// ?       alipay  [no test files]
	gtNoFilesRE = "^\\?.*\\[no test files\\]$"
	// FAIL    node/config [build failed]
	gtBuildFailedRE = `^FAIL.*\[(build|setup) failed\]$`

	// gocheck regular expressions

	// START: mmath_test.go:16: MySuite.TestAdd
	gcStartRE = "START: [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)"

	// PASS: mmath_test.go:16: MySuite.TestAdd	0.000s
	// FAIL: mmath_test.go:35: MySuite.TestDiv
	gcEndRE = "(PASS|FAIL|SKIP|PANIC|MISS): [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)[[:space:]]?(-?[0-9]+.[0-9]+)?"

	// FAIL	go2xunit/demo-gocheck	0.008s
	// ok  	go2xunit/demo-gocheck	0.008s
	gcSuiteRE = "^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(-?\\d+.\\d+)"
)

var (
	matchDatarace = regexp.MustCompile("^WARNING: DATA RACE$").MatchString
)

// LineScanner scans lines and keep track of line numbers
type LineScanner struct {
	*bufio.Scanner
	lnum int
}

// NewLineScanner creates a new line scanner from r
func NewLineScanner(r io.Reader) *LineScanner {
	scan := bufio.NewScanner(r)
	return &LineScanner{scan, 0}
}

// Scan advances to next line
func (ls *LineScanner) Scan() bool {
	val := ls.Scanner.Scan()
	ls.lnum++
	return val
}

// Line returns the current line number
func (ls *LineScanner) Line() int {
	return ls.lnum
}

// hasDatarace checks if there's a data race warning in the line
func hasDatarace(lines []string) bool {
	for _, line := range lines {
		if matchDatarace(line) {
			return true
		}
	}
	return false
}

// Token2Status return matching status for token
func Token2Status(token string) Status {
	switch token {
	case "FAIL", "PANIC":
		return Failed
	case "PASS":
		return Passed
	case "SKIP", "MISS":
		return Skipped
	}
	return UnknownStatus
}

// ParseGocheck parses output of "go test -gocheck.vv", returns a list of tests
// See data/gocheck.out for an example
// TODO: Refactor to shorter ones
func ParseGocheck(rd io.Reader, suitePrefix string) (Suites, error) {
	findStart := regexp.MustCompile(gcStartRE).FindStringSubmatch
	findEnd := regexp.MustCompile(gcEndRE).FindStringSubmatch
	findSuite := regexp.MustCompile(gcSuiteRE).FindStringSubmatch

	scanner := NewLineScanner(rd)
	var suites = make([]*Suite, 0)
	var suiteName string
	var suite *Suite

	var testName string
	var out []string

	for scanner.Scan() {
		line := scanner.Text()

		tokens := findStart(line)
		if len(tokens) > 0 {
			if tokens[2] == "SetUpTest" || tokens[2] == "TearDownTest" {
				continue
			}
			if testName != "" {
				return nil, fmt.Errorf("%d: start in middle\n", scanner.Line())
			}
			suiteName = tokens[1]
			testName = tokens[2]
			out = []string{}
			continue
		}

		tokens = findEnd(line)
		if len(tokens) > 0 {
			if tokens[3] == "SetUpTest" || tokens[3] == "TearDownTest" {
				continue
			}
			if testName == "" {
				return nil, fmt.Errorf("%d: orphan end", scanner.Line())
			}
			if (tokens[2] != suiteName) || (tokens[3] != testName) {
				return nil, fmt.Errorf("%d: suite/name mismatch", scanner.Line())
			}
			test := &Test{Name: testName}
			test.Message = strings.Join(out, "\n")
			test.Time = tokens[4]
			test.Status = Token2Status(tokens[1])
			if test.Status == UnknownStatus {
				return nil, fmt.Errorf("%d: unknown status %s", scanner.Line(), tokens[1])
			}

			if suite == nil || suite.Name != suiteName {
				suite = &Suite{Name: suitePrefix + suiteName}
				suites = append(suites, suite)
			}
			suite.Tests = append(suite.Tests, test)

			testName = ""
			suiteName = ""
			out = []string{}

			continue
		}

		// last "suite" is test summary
		tokens = findSuite(line)
		if tokens != nil {
			if suite == nil {
				suite = &Suite{Name: tokens[2], Status: tokens[1], Time: tokens[3]}
				suites = append(suites, suite)
			} else {
				suite.Status = tokens[1]
				suite.Time = tokens[3]
			}

			testName = ""
			suiteName = ""
			out = []string{}

			continue
		}

		if testName != "" {
			out = append(out, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return Suites(suites), nil
}

// ParseGotest parser output of gotest
// TODO: Make it shorter
func ParseGotest(rd io.Reader, suitePrefix string) (Suites, error) {
	findStart := regexp.MustCompile(gtStartRE).FindStringSubmatch
	findEnd := regexp.MustCompile(gtEndRE).FindStringSubmatch
	findSuite := regexp.MustCompile(gtSuiteRE).FindStringSubmatch
	isNoFiles := regexp.MustCompile(gtNoFilesRE).MatchString
	isBuildFailed := regexp.MustCompile(gtBuildFailedRE).MatchString
	isExit := regexp.MustCompile("^exit status -?\\d+").MatchString

	suites := []*Suite{}
	subTests := map[string]*Test{}
	var parentTest *Test
	var curTest *Test
	var curSuite *Suite
	var out []string
	suiteStack := SuiteStack{}
	// Handles a test that ended with a panic.
	handlePanic := func() {
		curTest.Status = Failed
		curTest.Time = "N/A"
		curSuite.Tests = append(curSuite.Tests, curTest)
		curTest = nil
	}

	// Appends output to the last test.
	appendError := func() {
		if len(out) > 0 && curSuite != nil && len(curSuite.Tests) > 0 {
			message := strings.Join(out, "\n")
			if curSuite.Tests[len(curSuite.Tests)-1].Message == "" {
				curSuite.Tests[len(curSuite.Tests)-1].Message = message
			} else {
				curSuite.Tests[len(curSuite.Tests)-1].Message += "\n" + message
			}
		}
		out = []string{}
	}

	scanner := NewLineScanner(rd)
	for scanner.Scan() {
		line := scanner.Text()

		// TODO: Only outside a suite/test, report as empty suite?
		if isNoFiles(line) {
			continue
		}

		if isBuildFailed(line) {
			return nil, fmt.Errorf("%d: package build failed: %s", scanner.Line(), line)
		}

		if curSuite == nil {
			curSuite = &Suite{}
		}

		tokens := findStart(line)
		if tokens != nil {
			subTest := false
			if curTest != nil {
				// This occurs when the last test ended with a panic, or when subtests are found
				if parentTest == nil && strings.HasPrefix(tokens[1], curTest.Name+"/") {
					// First subtest after parent
					parentTest = curTest
					curSuite.Tests = append(curSuite.Tests, curTest)
					subTests = map[string]*Test{}
					subTest = true
				} else if parentTest != nil && strings.HasPrefix(tokens[1], parentTest.Name+"/") {
					subTest = true
				} else if suiteStack.count == 0 {
					suiteStack.Push(curSuite)
					curSuite = &Suite{Name: curTest.Name}
				} else {
					handlePanic()
				}
			}
			appendError()
			curTest = &Test{
				Name: tokens[1],
			}
			if subTest {
				curSuite.Tests = append(curSuite.Tests, curTest)
				subTests[curTest.Name] = curTest
			}
			continue
		}

		tokens = findEnd(line)
		if tokens != nil {
			appendTest := true
			if parentTest != nil && tokens[2] == parentTest.Name {
				curTest = parentTest
				parentTest = nil
				appendTest = false
			} else if subTest, ok := subTests[tokens[2]]; ok {
				curTest = subTest
				appendTest = false
			} else {
				parentTest = nil
				subTests = map[string]*Test{}
				if curTest == nil {
					if suiteStack.count > 0 {
						prevSuite := suiteStack.Pop()
						suites = append(suites, curSuite)
						curSuite = prevSuite
						continue
					} else {
						return nil, fmt.Errorf("%d: orphan end test", scanner.Line())
					}
				}
				if tokens[2] != curTest.Name {
					err := fmt.Errorf("%d: name mismatch (try disabling parallel mode)", scanner.Line())
					return nil, err
				}
			}
			curTest.Status = Token2Status(tokens[1])
			if curTest.Status == UnknownStatus {
				return nil, fmt.Errorf("%d: unknown status - %s", scanner.Line(), tokens[1])
			}
			if Options.FailOnRace && hasDatarace(out) {
				curTest.Status = Failed
			}
			curTest.Time = tokens[3]
			curTest.Message = strings.Join(out, "\n")
			if appendTest {
				curSuite.Tests = append(curSuite.Tests, curTest)
			}
			curTest = nil
			out = []string{}
			continue
		}

		tokens = findSuite(line)
		if tokens != nil {
			if curTest != nil {
				// This occurs when the last test ended with a panic.
				handlePanic()
			}
			appendError()
			curSuite.Name = suitePrefix + tokens[2]
			curSuite.Time = tokens[3]
			suites = append(suites, curSuite)
			curSuite = nil
			continue
		}

		if isExit(line) || (line == "FAIL") || (line == "PASS") {
			continue
		}

		out = append(out, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if curTest != nil {
		// This occurs when the last test fatal'd outside of the `go test` runner.
		handlePanic()
	}

	// If there were no suites found, but everything else went OK, return a
	// generic suite.
	if len(suites) == 0 && curSuite != nil {
		if curSuite.Name == "" {
			curSuite.Name = suitePrefix
		}
		// Catch any post-failure messages from the last test
		appendError()
		suites = append(suites, curSuite)
	}

	return Suites(suites), nil
}
