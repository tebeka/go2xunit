package lib

// Status is test status
type Status int

// Test statuses
const (
	UnknownStatus Status = iota
	Failed
	Skipped
	Passed
)

// Test data structure
type Test struct {
	Name, Time, Message string
	Status              Status
	AppendedErrorOutput	bool
	isParentTest		bool
}

// Suite of tests (found in some unit testing frameworks)
type Suite struct {
	Name   string
	Time   string
	Status string
	Tests  []*Test
}

// NumPassed return number of passed tests in the suite
func (suite *Suite) NumPassed() int {
	return suite.numStatus(Passed)
}

// NumSkipped return number of skipped tests in suite
func (suite *Suite) NumSkipped() int {
	return suite.numStatus(Skipped)
}

// NumFailed return number of failed tests in suite
func (suite *Suite) NumFailed() int {
	return suite.numStatus(Failed)
}

// numStatus returns the number of tests in status
func (suite *Suite) numStatus(status Status) int {
	count := 0
	for _, test := range suite.Tests {
		if test.Status == status {
			count++
		}
	}
	return count
}

// Len return the number of tests in the suite
func (suite *Suite) Len() int {
	return len(suite.Tests)
}

// Suites is a list of suites
type Suites []*Suite

// HasFailures return true is there's at least one failing suite
func (s Suites) HasFailures() bool {
	for _, suite := range s {
		if suite.NumFailed() > 0 {
			return true
		}
	}
	return false
}

// SuiteStack is a stack of test suites
type SuiteStack struct {
	nodes []*Suite
	count int
}

// Push adds a node to the stack.
func (s *SuiteStack) Push(n *Suite) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *SuiteStack) Pop() *Suite {
	if s.count == 0 {
		return nil
	}
	s.count--
	return s.nodes[s.count]
}
