package main

import (
	"os"
	"testing"

	"bitbucket.org/tebeka/go2xunit/types"
	"bitbucket.org/tebeka/go2xunit/gocheck"
	"bitbucket.org/tebeka/go2xunit/gotest"
)

func Test_NumFailed(t *testing.T) {
	suite := &types.Suite{
		Tests: []*types.Test{},
	}
	if suite.NumFailed() != 0 {
		t.Fatal("found more than 1 failure in empty list")
	}

	suite.Tests = []*types.Test{
		&types.Test{Failed: false},
		&types.Test{Failed: true},
		&types.Test{Failed: false},
	}

	if suite.NumFailed() != 1 {
		t.Fatal("can't count")
	}
}

func loadTests(filename string, t *testing.T) ([]*types.Suite, error) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("can't open %s - %s", filename, err)
	}

	return gotest.Parse(file)
}

func Test_parseOutput(t *testing.T) {
	filename := "data/gotest.out"
	suites, err := loadTests(filename, t)
	if err != nil {
		t.Fatalf("error loading %s - %s", filename, err)
	}

	if len(suites) != 2 {
		t.Fatalf("got %d suites instead of 2", len(suites))
	}
	if suites[0].Name != "_/home/miki/Projects/goroot/src/xunit" {
		t.Fatalf("bad Suite name %s, expected _/home/miki/Projects/goroot/src/xunit", suites[0].Name)
	}
	if suites[1].Name != "_/home/miki/Projects/goroot/src/anotherTest" {
		t.Fatalf("bad Suite name %s, expected _/home/miki/Projects/goroot/src/anotherTest", suites[1].Name)
	}

	suite := suites[0]
	tests := suite.Tests
	if len(tests) != 4 {
		t.Fatalf("got %d tests instead of 4", len(tests))
	}

	if suite.NumFailed() != 1 {
		t.Fatalf("wrong number of failed %d, should be 1", suite.NumFailed())
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
	if len(tests) != 1 {
		t.Fatalf("got %d tests instead of 1", len(tests))
	}

	if suite.NumFailed() != 0 {
		t.Fatalf("wrong number of failed %d, should be 0", suite.NumFailed())
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
	_, err := loadTests("go2xunit.go", t)
	if err == nil {
		t.Fatalf("managed to find suites in junk")
	}
}

func Test_parseGocheckPass(t *testing.T) {
	filename := "data/gocheck-pass.out"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("can't open %s - %s", filename, err)
	}

	suites, err := gocheck.Parse(file)

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

	suites, err := gocheck.Parse(file)

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
