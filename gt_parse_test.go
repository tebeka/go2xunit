package main

import (
	"io"
	"os"
	"testing"
)

func loadAndParseGoTest(filename string, t *testing.T) []*Suite {
	file := loadGoTestFile(t, filename)
	suites, err := gt_Parse(file)
	if err != nil {
		t.Fatalf("error parsing %s - %s", filename, err)
	}
	return suites
}

func loadGoTestFile(t *testing.T, filename string) io.Reader {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("can't open %s - %s", filename, err)
	}
	return file
}

func Test_parseOutput(t *testing.T) {
	filename := "data/gotest.out"
	suites := loadAndParseGoTest(filename, t)

	numSuites := 2
	if len(suites) != numSuites {
		t.Fatalf("got %d suites instead of %d", len(suites), numSuites)
	}
	if suites[0].Name != "_/home/miki/Projects/goroot/src/xunit" {
		t.Fatalf("bad Suite name %s, expected _/home/miki/Projects/goroot/src/xunit", suites[0].Name)
	}
	if suites[1].Name != "_/home/miki/Projects/goroot/src/anotherTest" {
		t.Fatalf("bad Suite name %s, expected _/home/miki/Projects/goroot/src/anotherTest", suites[1].Name)
	}

	suite := suites[0]
	tests := suite.Tests
	numTests := 5
	if len(tests) != numTests {
		t.Fatalf("got %d tests instead of %d", len(tests), numTests)
	}

	numFailed := 1
	if suite.NumFailed() != numFailed {
		t.Fatalf("wrong number of failed %d, should be %d", suite.NumFailed(), numFailed)
	}

	numSkipped := 1
	if suite.NumSkipped() != numSkipped {
		t.Fatalf("wrong number of skipped %d, should be %d", suite.NumSkipped(), numSkipped)
	}

	test := tests[0]
	if test.Name != "TestAdd" {
		t.Fatalf("bad test name %s, expected TestAdd", test.Name)
	}

	if test.Time != "0.00" {
		t.Fatalf("bad test time %s, expected 0.00", test.Time)
	}

	if len(test.Message) != 0 {
		t.Fatalf("bad message (should be empty): %s", test.Message)
	}

	suite = suites[1]
	tests = suite.Tests
	numTests = 1
	if len(tests) != numTests {
		t.Fatalf("got %d tests instead of %d", len(tests), numTests)
	}

	numFailed = 0
	if suite.NumFailed() != numFailed {
		t.Fatalf("wrong number of failed %d, should be %d", suite.NumFailed(), numFailed)
	}

	test = tests[0]
	if test.Name != "TestAdd" {
		t.Fatalf("bad test name %s, expected TestAdd", test.Name)
	}

	if test.Time != "0.00" {
		t.Fatalf("bad test time %s, expected 0.00", test.Time)
	}

	if len(test.Message) != 0 {
		t.Fatalf("bad message (should be empty): %s", test.Message)
	}
}

func Test_parseOutputBad(t *testing.T) {
	filename := "main.go"
	suites := loadAndParseGoTest(filename, t)
	if len(suites) > 0 {
		t.Fatalf("managed to find suites in junk")
	}
}

func Test_parseDatarace(t *testing.T) {
	failOnRace = true
	filename := "data/gotest-datarace.out"
	suites := loadAndParseGoTest(filename, t)

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
	numFailed := 1
	if suite.NumFailed() != numFailed {
		t.Fatalf("wrong number of failed %d, should be %d", suite.NumFailed(), numFailed)
	}
	failOnRace = false
}

func Test_ignoreDatarace(t *testing.T) {
	filename := "data/gotest-datarace.out"
	suites := loadAndParseGoTest(filename, t)

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

func Test_parseLogOutput(t *testing.T) {
	filename := "data/gotest-log.out"
	suites := loadAndParseGoTest(filename, t)

	suite := suites[0]
	tests := suite.Tests
	numTests := 1
	if len(tests) != numTests {
		t.Fatalf("got %d tests instead of %d", len(tests), numTests)
	}

	expectedMessage := "Log output."
	for i, test := range suite.Tests {
		if test.Message != expectedMessage {
			t.Errorf(
				"test %v message does not match expected result:\n\tGot: \"%v\"\n\tExpected: \"%v\"\n",
				i,
				test.Message,
				expectedMessage)
		}
	}
}

func Test_parsePanicOutput(t *testing.T) {
	filename := "data/gotest-panic.out"
	suites := loadAndParseGoTest(filename, t)

	suite := suites[0]
	tests := suite.Tests
	numTests := 1
	if len(tests) != numTests {
		t.Fatalf("got %d tests instead of %d", len(tests), numTests)
	}

	expectedMessage := "fatal error: all goroutines are asleep - deadlock!\n..."
	test := tests[0]
	if test.Message != expectedMessage {
		t.Errorf(
			"test message does not match expected result:\n\tGot: \"%v\"\n\tExpected: \"%v\"\n",
			test.Message,
			expectedMessage)
	}
}

func Test_parseEmptySuite(t *testing.T) {
	filename := "data/gotest-empty.out"
	suites := loadAndParseGoTest(filename, t)

	suite := suites[0]
	tests := suite.Tests
	numTests := 0
	if len(tests) != numTests {
		t.Fatalf("got %d tests instead of %d", len(tests), numTests)
	}
}

func Test_NoFiles(t *testing.T) {
	filename := "data/gotest-nofiles.out"
	suites := loadAndParseGoTest(filename, t)

	count := 1
	if len(suites) != count {
		t.Fatalf("bad number of suites. expected %d got %d", count, len(suites))
	}
}

func Test_Multi(t *testing.T) {
	filename := "data/gotest-multi.out"
	suites := loadAndParseGoTest(filename, t)

	count := 1
	if len(suites) != count {
		t.Fatalf("bad number of suites. expected %d got %d", count, len(suites))
	}
}

func Test_ThatMessageIsParsedCorrectly_WhenThereIsAnErrorWithinTheLastTestInSuite(t *testing.T) {
	filename := "data/gotest-multierror.out"
	suites := loadAndParseGoTest(filename, t)

	count := 1
	if len(suites) != count {
		t.Fatalf("bad number of suites. expected %d got %d", count, len(suites))
	}

	suite := suites[0]

	if len(suite.Tests) != 2 {
		t.Fatalf("bad number of tests. expected %d got %d", 2, len(suite.Tests))
	}

	expectedMessages := []string{
		"\tmain_test.go:10: something went wrong",
		"\tmain_test.go:14: something new went wrong",
	}

	for i, test := range suite.Tests {
		if test.Message != expectedMessages[i] {
			t.Errorf(
				"test %v message does not match expected result:\n\tGot: \"%v\"\n\tExpected: \"%v\"\n",
				i,
				test.Message,
				expectedMessages[i])
		}
	}
}

func Test_parseBuildFailed(t *testing.T) {
	file := loadGoTestFile(t, "data/gotest-buildfailed.out")
	_, err := gt_Parse(file)
	if err == nil {
		t.Fatalf("expected error when at least one package failed to build")
	}
}

func Test_nameWithNum(t *testing.T) {
	file := loadGoTestFile(t, "data/gotest-num.out")
	_, err := gt_Parse(file)
	if err != nil {
		t.Fatalf("didn't parse name with number")
	}
}

func Test_0Time(t *testing.T) {
	file := loadGoTestFile(t, "data/gotest-0.out")
	suites, err := gt_Parse(file)
	if err != nil {
		t.Fatalf("didn't parse name with number")
	}

	if len(suites) != 1 {
		t.Fatalf("bad number of suites - %d", len(suites))
	}

	suite := suites[0]
	if suite.Count() != 2 {
		t.Fatalf("bad number of tests. expected %d got %d", 2, len(suite.Tests))
	}
	if suite.NumFailed() != 0 {
		t.Fatalf("unexpected failure")
	}
}

func Test_TestifySuite(t *testing.T) {
	suites := loadAndParseGoTest("data/gotest-testify-suite.out", t)
	if len(suites) != 2 {
		t.Fatalf("Wrong number of suites", len(suites))
	}
}
