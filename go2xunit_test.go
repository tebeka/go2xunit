package main

import (
	"os"
	"testing"
)

func Test_numFailuers(t *testing.T) {
	tests := []*Test{}
	if numFailures(tests) != 0 {
		t.Fatal("found more than 1 failure in empty list")
	}

	tests = []*Test{
		&Test{Failed: false},
		&Test{Failed: true},
		&Test{Failed: false},
	}
	if numFailures(tests) != 1 {
		t.Fatal("can't count")
	}
}

func loadTests(filename string, t *testing.T) []*Suite {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("can't open %s - %s", filename, err)
	}

	suites, err := parseOutput(file)
	if err != nil {
		t.Fatalf("error parsing - %s", err)
	}

	return suites
}

func Test_parseOutput(t *testing.T) {
	suites := loadTests("gotest.out", t)
	if len(suites) != 2 {
		t.Fatalf("got %d suites instead of 2", len(suites))
	}
	if suites[0].Name != "_/home/miki/Projects/goroot/src/xunit" {
		t.Fatalf("bad Suite name %s, expected _/home/miki/Projects/goroot/src/xunit", suites[0].Name)
	}
	if suites[1].Name != "_/home/miki/Projects/goroot/src/anotherTest" {
		t.Fatalf("bad Suite name %s, expected _/home/miki/Projects/goroot/src/anotherTest", suites[1].Name)
	}

	tests := suites[0].Tests
	if len(tests) != 4 {
		t.Fatalf("got %d tests instead of 4", len(tests))
	}

	if numFailures(tests) != 1 {
		t.Fatalf("wrong number of failed %d, should be 1", numFailures(tests))
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
	tests = suites[1].Tests
	if len(tests) != 1 {
		t.Fatalf("got %d tests instead of 1", len(tests))
	}

	if numFailures(tests) != 0 {
		t.Fatalf("wrong number of failed %d, should be 0", numFailures(tests))
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
	suites := loadTests("go2xunit.go", t)
	if len(suites) != 0 {
		t.Fatalf("managed to find suites in junk")
	}
}
