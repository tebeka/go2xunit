package lib

import (
	"os"
	"testing"
)

func loadGotest(filename string, t *testing.T) ([]*Suite, error) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("can't open %s - %s", filename, err)
	}

	return ParseGotest(file, "")
}

func Test_parseOutputBad(t *testing.T) {
	filename := "parsers.go"
	if _, err := loadGotest(filename, t); err == nil {
		t.Fatalf("managed to find suites in junk")
	}
}

func Test_ignoreDatarace(t *testing.T) {
	filename := "../data/in/gotest-datarace.out"
	suites, err := loadGotest(filename, t)
	if err != nil {
		t.Fatalf("error loading %s - %s", filename, err)
	}

	suite := suites[0]
	tests := suite.Tests
	numTests := 1
	if len(tests) != numTests {
		t.Fatalf("got %d tests instead of %d", len(tests), numTests)
	}

	expectedMessage := "WARNING: DATA RACE"
	for i, test := range suite.Tests {
		if test.Message != expectedMessage {
			t.Errorf(
				"test %v message does not match expected result:\n\tGot: \"%v\"\n\tExpected: \"%v\"\n",
				i,
				test.Message,
				expectedMessage)
		}
	}
	numFailed := 0
	if suite.NumFailed() != numFailed {
		t.Fatalf("wrong number of failed %d, should be %d", suite.NumFailed(), numFailed)
	}
}
