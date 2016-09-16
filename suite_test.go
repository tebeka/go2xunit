package main

import "testing"

func Test_EmptySuite(t *testing.T) {
	suite := Suite{}
	if count := suite.Count(); count != 0 {
		t.Fatal("Expected 0 tests; got:", count)
	}
	if failures := suite.NumFailed(); failures != 0 {
		t.Fatal("Expected 0 failures, got:", failures)
	}
	if skipped := suite.NumSkipped(); skipped != 0 {
		t.Fatal("Expected 0 skipped, got:", skipped)
	}
	if passed := suite.NumPassed(); passed != 0 {
		t.Fatal("Expected 0 passed, got:", passed)
	}
}

func Test_NumFailed_Mixed(t *testing.T) {
	suite := &Suite{
		Tests: []*Test{
			&Test{Skipped: false},
			&Test{Failed: false},
			&Test{Passed: false},
			&Test{Skipped: true},
			&Test{Failed: true},
			&Test{Passed: true},
		},
	}

	if count := suite.Count(); count != 6 {
		t.Fatal("Expected 6 tests; got:", count)
	}
	if failures := suite.NumFailed(); failures != 1 {
		t.Fatal("Expected 1 failures, got:", failures)
	}
	if skipped := suite.NumSkipped(); skipped != 1 {
		t.Fatal("Expected 1 skipped, got:", skipped)
	}
	if passed := suite.NumPassed(); passed != 1 {
		t.Fatal("Expected 1 passed, got:", passed)
	}
}

func TestMultipleSuitsWithoutFailures(t *testing.T) {
	suites := []*Suite{
		&Suite{},
		&Suite{Tests: []*Test{&Test{Passed: true}}},
		&Suite{Tests: []*Test{&Test{Skipped: true}}},
	}
	if hasFailures(suites) {
		t.Fatal("Expected false, got: true")
	}
}
