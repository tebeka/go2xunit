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

	// FIXME: This randomly fails with len(suite.Tests) = 1
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
