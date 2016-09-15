package main

import "testing"

func Test_NumFailed(t *testing.T) {
	suite := &Suite{
		Tests: []*Test{
			&Test{Failed: false},
			&Test{Failed: true},
			&Test{Failed: false},
		},
	}

	if suite.NumFailed() != 1 {
		t.Fatal("can't count")
	}
}

func Test_NumFailed_Empty(t *testing.T) {
	suite := Suite{}
	if failures := suite.NumFailed(); failures != 0 {
		t.Fatal("Expected 0 failures, got:", failures)
	}
}
