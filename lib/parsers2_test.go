package lib

import (
	"os"
	"testing"
)

func TestGotestParser(t *testing.T) {
	file, err := os.Open("../data/in/gotest-0.out")
	if err != nil {
		t.Fatal(err)
	}
	pr := NewGtParser(file)
	nt := 0
	var s *Suite
	for pr.Scan() {
		s = pr.Suite()
		nt++
	}
	if pr.Err() != nil {
		t.Fatal(err)
	}
	if nt != 1 {
		t.Fatalf("wrong number of suites %d (expected 1)", nt)
	}
	if len(s.Tests) != 2 {
		t.Fatalf("wrong number of tests %d (expected 2)", len(s.Tests))
	}

	names := []string{"ExampleA", "ExampleOp"}
	for i, test := range s.Tests {
		if test.Name != names[i] {
			t.Fatalf("name mismatch %q != %q", test.Name, names[i])
		}
	}
}
