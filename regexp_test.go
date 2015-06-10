package main

import (
	"regexp"
	"testing"
)

func Test_GcSuiteRePass(t *testing.T) {
	find_suite := regexp.MustCompile(gc_endRE).FindStringSubmatch
	tokens := find_suite("PASS: mmath_test.go:16: MySuite.TestAdd	0.000s")

	for i, expected := range []string{"PASS", "MySuite", "TestAdd", "0.000"} {
		checkExpected(t, expected, tokens[i+1])
	}
}

func Test_GcSuiteReFail(t *testing.T) {
	find_suite := regexp.MustCompile(gc_endRE).FindStringSubmatch
	tokens := find_suite("FAIL: mmath_test.go:35: MySuite.TestDiv")

	for i, expected := range []string{"FAIL", "MySuite", "TestDiv", ""} {
		checkExpected(t, expected, tokens[i+1])
	}
}

func Test_GcSuiteReSkip(t *testing.T) {
	find_suite := regexp.MustCompile(gc_endRE).FindStringSubmatch
	tokens := find_suite("SKIP: mmath_test.go:35: MySuite.TestMul (not implemented)")

	for i, expected := range []string{"SKIP", "MySuite", "TestMul", ""} {
		checkExpected(t, expected, tokens[i+1])
	}
}

func checkExpected(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Fatalf("Expected %s but was %s\n", expected, actual)
	}
}
