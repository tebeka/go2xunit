// Package lib is exposing parsers and output generation
package lib

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	matchDatarace = regexp.MustCompile("^WARNING: DATA RACE$").MatchString
)

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

// Returns previous test in a suite, for a given test. Returns error if previous
// test doesn't exist.
func getPreviousFailTest(suite *Suite, curTest *Test) (*Test, error) {
	previousFailTestIndex := -1
	for testIndex, test := range suite.Tests {
		if test.Name == curTest.Name {
			break
		} else {
			if test.Status == Failed {
				previousFailTestIndex = testIndex
			}
		}
	}

	if previousFailTestIndex >= 0 {
		return suite.Tests[previousFailTestIndex], nil
	}
	return nil, fmt.Errorf("Not found previous test of %s in suite %s", curTest.Name, suite.Name)
}

// ParseGocheck parses output of "go test -gocheck.vv", returns a list of tests
// See data/gocheck.out for an example
// TODO: Refactor to shorter ones
func ParseGocheck(rd io.Reader, suitePrefix string) (Suites, error) {
	findStart := gcStartRE.FindStringSubmatch
	findEnd := gcEndRE.FindStringSubmatch
	findSuite := gcSuiteRE.FindStringSubmatch

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
				return nil, fmt.Errorf("%d: start in middle of test", scanner.Line())
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
	findStart := gtStartRE.FindStringSubmatch
	findEnd := gtEndRE.FindStringSubmatch
	findSuite := gtSuiteRE.FindStringSubmatch
	isNoFiles := gtNoFilesRE.MatchString
	isBuildFailed := gtBuildFailedRE.MatchString
	isExit := gtExitRE.MatchString
	isErrorOutput := gcTestErrorRE.MatchString

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
		curTest.Time = "0"
		curSuite.Tests = append(curSuite.Tests, curTest)
		curTest = nil
	}

	// Appends output to the last test.
	appendError := func() {
		if len(out) > 0 && curSuite != nil && len(curSuite.Tests) > 0 {
			message := strings.Join(out, "\n")
			test := curSuite.Tests[len(curSuite.Tests)-1]
			if test.isParentTest == false {
				if test.Message == "" {
					test.Message = message
				} else {
					test.Message += "\n" + message
				}
				test.AppendedErrorOutput = isErrorOutput(message)
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
					parentTest.isParentTest = true
					parentTest.AppendedErrorOutput = true
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

			if len(out) > 0 {
				message := strings.Join(out, "\n")
				prevTest, err := getPreviousFailTest(curSuite, curTest)
				var test *Test
				if err == nil && prevTest.AppendedErrorOutput == false {
					test = prevTest
				} else {
					test = curTest
				}
				if test.isParentTest == false {
					test.Message += message
					test.AppendedErrorOutput = isErrorOutput(message)
				}
			}

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
	}

	if curSuite != nil && len(curSuite.Tests) > 0 {
		suites = append(suites, curSuite)
	}

	return Suites(suites), nil
}
