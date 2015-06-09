package main

// Test suite
type Suite struct {
	Name   string
	Time   string
	Status string
	Tests  []*Test
}

// Number of failed tests in the suite
func (suite *Suite) NumFailed() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Failed {
			count++
		}
	}
	return count
}

// Number of skipped tests in the suite
func (suite *Suite) NumSkipped() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Skipped {
			count++
		}
	}
	return count
}

// Number of passed tests in the suite
func (suite *Suite) NumPassed() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Passed {
			count++
		}
	}
	return count
}

// Number of tests in the suite
func (suite *Suite) Count() int {
	return len(suite.Tests)
}

// Stack of test suites
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
