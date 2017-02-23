package lib

import (
	"os"
	"testing"
)

func TestGtParser(t *testing.T) {
	fname := "../data/in/gotest-1.7.out"
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}

	parser := NewGtParser(file)
	numTests := 7
	sc := 0
	for parser.Scan() {
		sc++
		suite := parser.Suite()
		if len(suite.Tests) != numTests {
			t.Fatalf("wrong number of tests %d != %d", len(suite.Tests), numTests)
		}
	}
	if err := parser.Err(); err != nil {
		t.Fatal(err)
	}
	if sc != 1 {
		t.Fatalf("bad number of suites %d != 1", sc)
	}
}
