package main

import (
	"os"
	"testing"
)

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

func loadGotest(filename string, t *testing.T) ([]*Suite, error) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("can't open %s - %s", filename, err)
	}

	return gt_Parse(file)
}

func Test_parseOutput16(t *testing.T) {
	filename := "data/gotest-1.6.out"
	suites, err := loadGotest(filename, t)
	if err != nil {
		t.Fatalf("error loading %s - %s", filename, err)
	}
	numSuites := 1
	if len(suites) != numSuites {
		t.Fatalf("got %d suites instead of %d", len(suites), numSuites)
	}

	suiteName := "bitbucket.org/tebeka/go2xunit/demo"

	if suites[0].Name != suiteName {
		t.Fatalf("bad Suite name %s, expected %s", suites[0].Name, suiteName)
	}

}

func Test_parseOutput(t *testing.T) {
	checkOutput(t, "data/gotest.out")
}

func Test_parseOutputPrint(t *testing.T) {
	checkOutput(t, "data/gotest-print.out")
}

func checkOutput(t *testing.T, filename string) {
	suites, err := loadGotest(filename, t)
	if err != nil {
		t.Fatalf("error loading %s - %s", filename, err)
	}

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
	filename := "parsers.go"
	if _, err := loadGotest(filename, t); err == nil {
		t.Fatalf("managed to find suites in junk")
	}
}

func Test_parseDatarace(t *testing.T) {
	origFail := args.failOnRace
	args.failOnRace = true
	defer func() {
		args.failOnRace = origFail
	}()

	filename := "data/gotest-datarace.out"
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
	numFailed := 1
	if suite.NumFailed() != numFailed {
		t.Fatalf("wrong number of failed %d, should be %d", suite.NumFailed(), numFailed)
	}
}

func Test_ignoreDatarace(t *testing.T) {
	filename := "data/gotest-datarace.out"
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

func Test_parseLogOutput(t *testing.T) {
	filename := "data/gotest-log.out"
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
	suites, err := loadGotest(filename, t)
	if err != nil {
		t.Fatalf("error loading %s - %s", filename, err)
	}

	suite := suites[0]
	tests := suite.Tests
	numTests := 0
	if len(tests) != numTests {
		t.Fatalf("got %d tests instead of %d", len(tests), numTests)
	}
}

func Test_parseGocheckPass(t *testing.T) {
	filename := "data/gocheck-pass.out"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("can't open %s - %s", filename, err)
	}

	suites, err := gc_Parse(file)

	if err != nil {
		t.Fatalf("can't parse %s - %s", filename, err)
	}

	nsuites := 1
	if len(suites) != nsuites {
		t.Fatalf("bad number of suites %d != %d", len(suites), nsuites)
	}

	suite := suites[0]
	ntests := 3
	if len(suite.Tests) != ntests {
		t.Fatalf("bad number of tests %d != %d", len(suite.Tests), ntests)
	}

	name := "MySuite"
	if suite.Name != name {
		t.Fatalf("bad suite name %s != %s", suite.Name, name)
	}

	if suite.NumFailed() > 0 {
		t.Fatalf("suite failed, should have passed")
	}
}

func Test_parseGocheckFail(t *testing.T) {
	filename := "data/gocheck-fail.out"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("can't open %s - %s", filename, err)
	}

	suites, err := gc_Parse(file)

	if err != nil {
		t.Fatalf("can't parse %s - %s", filename, err)
	}

	nsuites := 2
	if len(suites) != nsuites {
		t.Fatalf("bad number of suites %d != %d", len(suites), nsuites)
	}

	suite := suites[1]
	ntests := 3
	if len(suite.Tests) != ntests {
		t.Fatalf("bad number of tests %d != %d", len(suite.Tests), ntests)
	}

	name := "MySuite"
	if suite.Name != name {
		t.Fatalf("bad suite name %s != %s", suite.Name, name)
	}

	nfailed := 1
	if suite.NumFailed() != nfailed {
		t.Fatalf("num failed differ %d != %d", suite.NumFailed(), nfailed)
	}
}

func Test_NoFiles(t *testing.T) {
	filename := "data/gotest-nofiles.out"
	suites, err := loadGotest(filename, t)

	if err != nil {
		t.Fatalf("error loading %s - %s", filename, err)
	}

	count := 1
	if len(suites) != count {
		t.Fatalf("bad number of suites. expected %d got %d", count, len(suites))
	}
}

func Test_Multi(t *testing.T) {
	filename := "data/gotest-multi.out"
	suites, err := loadGotest(filename, t)
	if err != nil {
		t.Fatalf("error loading %s - %s", filename, err)
	}

	count := 1
	if len(suites) != count {
		t.Fatalf("bad number of suites. expected %d got %d", count, len(suites))
	}
}

func Test_ThatMessageIsParsedCorrectly_WhenThereIsAnErrorWithinTheLastTestInSuite(t *testing.T) {
	filename := "data/gotest-multierror.out"
	suites, err := loadGotest(filename, t)
	if err != nil {
		t.Fatalf("error loading %s - %s", filename, err)
	}

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
	_, err := loadGotest("data/gotest-buildfailed.out", t)
	if err == nil {
		t.Fatalf("expected error when at least one package failed to build")
	}
}

func Test_nameWithNum(t *testing.T) {
	_, err := loadGotest("data/gotest-num.out", t)
	if err != nil {
		t.Fatalf("didn't parse name with number")
	}
}

func Test_0Time(t *testing.T) {
	suites, err := loadGotest("data/gotest-0.out", t)
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
	suites, err := loadGotest("data/gotest-testify-suite.out", t)
	if err != nil {
		t.Fatalf("Failed to parse")
	}
	if len(suites) != 2 {
		t.Fatalf("Wrong number of suites", len(suites))
	}
}
