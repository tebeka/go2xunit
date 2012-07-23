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

func loadTests(filename string, t *testing.T) []*Test {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("can't open %s - %s", filename, err)
	}

	tests, err := parseOutput(file)
	if err != nil {
		t.Fatalf("error parsing - %s", err)
	}

	return tests
}

func Test_parseOutput(t *testing.T) {
	tests := loadTests("gotest.out", t)
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
}

func Test_parseOutputBad(t *testing.T) {
	tests := loadTests("go2xunit.go", t)
	if len(tests) != 0 {
		t.Fatalf("managed to find tests in junk")
	}
}
