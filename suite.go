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
	tests  []*Test
}

// AddTest adds a test case to the suite
func (suite *Suite) AddTest(test *Test) {
	suite.tests = append(suite.tests, test)
}

// Tests returns a collection of tests
func (suite *Suite) Tests() []*Test {
	return suite.tests
}

// LastTest returns the last test collected by the suite
func (suite *Suite) LastTest() *Test {
	return suite.tests[suite.Count()-1]
}

// SetLastTestMessage replaces the message of the last test collected by the suite
func (suite *Suite) SetLastTestMessage(message string) {
	suite.tests[suite.Count()-1].Message = message
}

// SetTests updates the suite with a new collection of tests
func (suite *Suite) SetTests(tests []*Test) {
	suite.tests = tests
}

// NumPassed return number of passed tests in the suite
func (suite *Suite) NumPassed() int {
	return suite.stats().passed
}

// NumSkipped return number of skipped tests in suite
func (suite *Suite) NumSkipped() int {
	return suite.stats().skipped
}

// NumFailed return number of failed tests in suite
func (suite *Suite) NumFailed() int {
	return suite.stats().failed
}

// report hold counts of the number of passed, skipped or failed
// tests.
type report struct {
	passed  int
	skipped int
	failed  int
}

// stats reports the number of passed, skipped or failed tests in a suite.
func (suite *Suite) stats() (r report) {
	for _, test := range suite.tests {
		if test.Passed {
			r.passed++
		}
		if test.Skipped {
			r.skipped++
		}
		if test.Failed {
			r.failed++
		}
	}
	return
}

// Count return the number of tests in the suite
func (suite *Suite) Count() int {
	return len(suite.tests)
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
