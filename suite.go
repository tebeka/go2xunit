package main

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

// NumSkipped return number of skipped tests in suite
func (suite *Suite) NumSkipped() int {
	return suite.stats().skipped
}

// NumFailed return number of failed tests in suite
func (suite *Suite) NumFailed() int {
	return suite.stats().failures
}

// NumPassed return number of passed tests in the suite
func (suite *Suite) NumPassed() int {
	return suite.stats().passed
}

type statsReport struct {
	passed   int
	skipped  int
	failures int
}

func (suite *Suite) stats() (report statsReport) {
	for _, test := range suite.Tests {
		if test.Failed {
			report.failures++
		}
		if test.Skipped {
			report.skipped++
		}
		if test.Passed {
			report.passed++
		}
	}

	return report
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
