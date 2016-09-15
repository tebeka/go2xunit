package main

import "testing"

func Test_NumFailed_Empty(t *testing.T) {
	suite := Suite{}
	if failures := suite.NumFailed(); failures != 0 {
		t.Fatal("Expected 0 failures, got:", failures)
	}
}

func Test_NumSkipped_Empty(t *testing.T) {
	suite := Suite{}
	if skipped := suite.NumSkipped(); skipped != 0 {
		t.Fatal("Expected 0 skipped, got:", skipped)
	}
}

func Test_NumFailed_Mixed(t *testing.T) {
	suite := &Suite{
		Tests: []*Test{
			&Test{Failed: false},
			&Test{Failed: true},
			&Test{Failed: false},
		},
	}

	if failures := suite.NumFailed(); failures != 1 {
		t.Fatal("Expected 1 failures, got:", failures)
	}
}

func Test_NumSkipped_Mixed(t *testing.T) {
	suite := &Suite{
		Tests: []*Test{
			&Test{Failed: false},
			&Test{Skipped: true},
			&Test{Failed: false},
		},
	}

	if skipped := suite.NumSkipped(); skipped != 1 {
		t.Fatal("Expected 1 skipped, got:", skipped)
	}
}
