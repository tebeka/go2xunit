package main

import (
	"os"
	"testing"
)

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
		t.Fatalf("bad number of suites, expected %d but was %d", nsuites, len(suites))
	}

	suite := suites[0]
	ntests := 3
	if len(suite.Tests) != ntests {
		t.Fatalf("bad number of tests, expected %d but was %d", ntests, len(suite.Tests))
	}

	name := "MySuite"
	if suite.Name != name {
		t.Fatalf("bad suite name, expected %s but was %s", name, suite.Name)
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
		t.Fatalf("bad number of suites, expected %d but was %d", nsuites, len(suites))
	}

	suite := suites[1]
	ntests := 3
	if len(suite.Tests) != ntests {
		t.Fatalf("bad number of tests, expected %d but was %d", ntests, len(suite.Tests))
	}

	name := "MySuite"
	if suite.Name != name {
		t.Fatalf("bad suite name, expected %s but was %s", name, suite.Name)
	}

	nfailed := 1
	if suite.NumFailed() != nfailed {
		t.Fatalf("num failed differ, expected %d but was %d", nfailed, suite.NumFailed())
	}
}
