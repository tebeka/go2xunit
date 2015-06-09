package main

import "testing"

func Test_NumFailed(t *testing.T) {
	suite := &Suite{
		Tests: []*Test{},
	}
	if suite.NumFailed() != 0 {
		t.Fatal("found more than 1 failure in empty list")
	}

	suite.Tests = []*Test{
		&Test{Failed: false},
		&Test{Failed: true},
		&Test{Failed: false},
	}

	if suite.NumFailed() != 1 {
		t.Fatal("can't count")
	}
}
