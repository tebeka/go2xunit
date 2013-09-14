package types

type Test struct {
	Name, Time, Message string
	Failed              bool
}

type Suite struct {
	Name   string
	Time   string
	Status string
	Tests  []*Test
}

func (suite *Suite) NumFailed() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Failed {
			count++
		}
	}

	return count
}

func (suite *Suite) Count() int {
	return len(suite.Tests)
}
