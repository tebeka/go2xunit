package main

import (
	"regexp"
	"testing"
)

func Test_GcSuiteRePass(t *testing.T) {
	find_suite := regexp.MustCompile(gc_endRE).FindStringSubmatch
	tokens := find_suite("PASS: mmath_test.go:16: MySuite.TestAdd	0.000s")

	expected := "PASS"
	actual := tokens[1]
	checkExpected(t, expected, actual)

	expected = "MySuite"
	actual = tokens[2]
	checkExpected(t, expected, actual)

	expected = "TestAdd"
	actual = tokens[3]
	checkExpected(t, expected, actual)

	expected = "0.000"
	actual = tokens[4]
	checkExpected(t, expected, actual)
}

func Test_GcSuiteReFail(t *testing.T) {
	find_suite := regexp.MustCompile(gc_endRE).FindStringSubmatch
	tokens := find_suite("FAIL: mmath_test.go:35: MySuite.TestDiv")

	expected := "FAIL"
	actual := tokens[1]
	checkExpected(t, expected, actual)

	expected = "MySuite"
	actual = tokens[2]
	checkExpected(t, expected, actual)

	expected = "TestDiv"
	actual = tokens[3]
	checkExpected(t, expected, actual)

	expected = ""
	actual = tokens[4]
	checkExpected(t, expected, actual)
}

func checkExpected(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Fatalf("Expected %s but was %s\n", expected, actual)
	}
}
