package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func gt_Parse(rd io.Reader) ([]*Suite, error) {
	find_start := regexp.MustCompile(gt_startRE).FindStringSubmatch
	find_end := regexp.MustCompile(gt_endRE).FindStringSubmatch
	find_suite := regexp.MustCompile(gt_suiteRE).FindStringSubmatch
	is_nofiles := regexp.MustCompile(gt_noFiles).MatchString
	is_buildFailed := regexp.MustCompile(gt_buildFailed).MatchString
	is_exit := regexp.MustCompile("^exit status -?\\d+").MatchString

	suites := []*Suite{}
	var curTest *Test
	var curSuite *Suite
	var out []string
	suiteStack := SuiteStack{}
	// Handles a test that ended with a panic.
	handlePanic := func() {
		curTest.Failed = true
		curTest.Skipped = false
		curTest.Passed = false
		curTest.Time = "N/A"
		curSuite.Tests = append(curSuite.Tests, curTest)
		curTest = nil
	}

	// Appends output to the last test.
	appendError := func() error {
		if len(out) > 0 && curSuite != nil && len(curSuite.Tests) > 0 {
			message := strings.Join(out, "\n")
			if curSuite.Tests[len(curSuite.Tests)-1].Message == "" {
				curSuite.Tests[len(curSuite.Tests)-1].Message = message
			} else {
				curSuite.Tests[len(curSuite.Tests)-1].Message += "\n" + message
			}
		}
		out = []string{}
		return nil
	}

	scanner := bufio.NewScanner(rd)
	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()

		// TODO: Only outside a suite/test, report as empty suite?
		if is_nofiles(line) {
			continue
		}

		if is_buildFailed(line) {
			return nil, fmt.Errorf("%d: package build failed: %s", lnum, line)
		}

		if curSuite == nil {
			curSuite = &Suite{}
		}

		tokens := find_start(line)
		if tokens != nil {
			if curTest != nil {
				// This occurs when the last test ended with a panic.
				if suiteStack.count == 0 {
					suiteStack.Push(curSuite)
					curSuite = &Suite{Name: curTest.Name}
				} else {
					handlePanic()
				}
			}
			if e := appendError(); e != nil {
				return nil, e
			}
			curTest = &Test{
				Name: tokens[1],
			}
			continue
		}

		tokens = find_end(line)
		if tokens != nil {
			if curTest == nil {
				if suiteStack.count > 0 {
					prevSuite := suiteStack.Pop()
					suites = append(suites, curSuite)
					curSuite = prevSuite
					continue
				} else {
					return nil, fmt.Errorf("%d: orphan end test", lnum)
				}
			}
			if tokens[2] != curTest.Name {
				err := fmt.Errorf("%d: name mismatch (try disabling parallel mode)", lnum)
				return nil, err
			}
			curTest.Failed = (tokens[1] == "FAIL") || (failOnRace && hasDatarace(out))
			curTest.Skipped = (tokens[1] == "SKIP")
			curTest.Passed = (tokens[1] == "PASS")
			curTest.Time = tokens[3]
			curTest.Message = strings.Join(out, "\n")
			curSuite.Tests = append(curSuite.Tests, curTest)
			curTest = nil
			out = []string{}
			continue
		}

		tokens = find_suite(line)
		if tokens != nil {
			if curTest != nil {
				// This occurs when the last test ended with a panic.
				handlePanic()
			}
			if e := appendError(); e != nil {
				return nil, e
			}
			curSuite.Name = tokens[2]
			curSuite.Time = tokens[3]
			suites = append(suites, curSuite)
			curSuite = nil
			continue
		}

		if is_exit(line) || (line == "FAIL") || (line == "PASS") {
			continue
		}

		out = append(out, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}


func hasDatarace(lines []string) bool {
	has_datarace := regexp.MustCompile("^WARNING: DATA RACE$").MatchString
	for _, line := range lines {
		if has_datarace(line) {
			return true
		}
	}
	return false
}

