package lib

import "testing"

func TestEmptySuite(t *testing.T) {
	suite := Suite{}
	t.Run("NumberOfTests", func(t *testing.T) {
		if count := suite.Len(); count != 0 {
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
			&Test{Status: Skipped},
			&Test{Status: Failed},
			&Test{Status: Passed},
			&Test{Status: Skipped},
			&Test{Status: Failed},
			&Test{Status: Passed},
		},
	}
	t.Run("NumberOfTests", func(t *testing.T) {
		if count := suite.Len(); count != 6 {
			t.Fatal("Expected 6 tests; got:", count)
		}
	})
	t.Run("NumberOfPassed", func(t *testing.T) {
		if passed := suite.NumPassed(); passed != 2 {
			t.Fatal("Expected 1 passed, got:", passed)
		}
	})
	t.Run("NumberOfSkipped", func(t *testing.T) {
		if skipped := suite.NumSkipped(); skipped != 2 {
			t.Fatal("Expected 1 skipped, got:", skipped)
		}
	})
	t.Run("NumberOfFailed", func(t *testing.T) {
		if failures := suite.NumFailed(); failures != 2 {
			t.Fatal("Expected 1 failures, got:", failures)
		}
	})
}

func TestMultipleSuits(t *testing.T) {
	empty := Suite{}
	passed := Suite{Tests: []*Test{&Test{Status: Passed}}}
	skipped := Suite{Tests: []*Test{&Test{Status: Skipped}}}
	failed := Suite{Tests: []*Test{&Test{Status: Failed}}}
	golden := Suites([]*Suite{&empty, &passed, &skipped})

	t.Run("WithoutFailures", func(t *testing.T) {
		suites := golden
		if suites.HasFailures() {
			t.Fatal("Expected false, got: true")
		}
	})
	t.Run("WithFailures", func(t *testing.T) {
		suites := append(golden, &failed)
		if !suites.HasFailures() {
			t.Fatal("Expected true, got: false")
		}
	})
}
