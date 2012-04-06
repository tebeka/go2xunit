package main

import (
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
