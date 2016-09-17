package main

import "testing"

func TestEmptySuite(t *testing.T) {
	suite := Suite{}
	t.Run("NumberOfTests", func(t *testing.T) {
		if count := suite.Count(); count != 0 {
			t.Fatal("Expected 0 tests; got:", count)
		}
	})
	t.Run("NumberOfPassed", func(t *testing.T) {
		if passed := suite.NumPassed(); passed != 0 {
			t.Fatal("Expected 0 passed, got:", passed)
		}
	})
	t.Run("NumberOfSkipped", func(t *testing.T) {
		if skipped := suite.NumSkipped(); skipped != 0 {
			t.Fatal("Expected 0 skipped, got:", skipped)
		}
	})
	t.Run("NumberOfFailed", func(t *testing.T) {
		if failures := suite.NumFailed(); failures != 0 {
			t.Fatal("Expected 0 failures, got:", failures)
		}
	})
}

func TestMixedSuite(t *testing.T) {
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
	t.Run("NumberOfTests", func(t *testing.T) {
		if count := suite.Count(); count != 6 {
			t.Fatal("Expected 6 tests; got:", count)
		}
	})
	t.Run("NumberOfPassed", func(t *testing.T) {
		if passed := suite.NumPassed(); passed != 1 {
			t.Fatal("Expected 1 passed, got:", passed)
		}
	})
	t.Run("NumberOfSkipped", func(t *testing.T) {
		if skipped := suite.NumSkipped(); skipped != 1 {
			t.Fatal("Expected 1 skipped, got:", skipped)
		}
	})
	t.Run("NumberOfFailed", func(t *testing.T) {
		if failures := suite.NumFailed(); failures != 1 {
			t.Fatal("Expected 1 failures, got:", failures)
		}
	})
}

func TestMultipleSuits(t *testing.T) {
	empty := Suite{}
	passed := Suite{Tests: []*Test{&Test{Passed: true}}}
	skipped := Suite{Tests: []*Test{&Test{Skipped: true}}}
	failed := Suite{Tests: []*Test{&Test{Failed: true}}}
	golden := []*Suite{&empty, &passed, &skipped}

	t.Run("WithoutFailures", func(t *testing.T) {
		suites := golden
		if hasFailures(suites) {
			t.Fatal("Expected false, got: true")
		}
	})
	t.Run("WithFailures", func(t *testing.T) {
		suites := append(golden, &failed)
		if !hasFailures(suites) {
			t.Fatal("Expected true, got: false")
		}
	})
}
