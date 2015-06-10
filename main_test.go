package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func Test_MainXunitNet(t *testing.T) {
	filename := "gocheck-pass.out"
	args := []string{"-gocheck", "-input", "data/" + filename}
	os.Args = append(os.Args, args...)

	expected, err := ioutil.ReadAll(getOutputData("xunit", filename))
	checkError(err)

	// this can be done only once or test framework will panic
	actual := []byte(captureOutput(main))

	if !bytes.Equal(expected, actual) {
		t.Errorf("xUnit.net XML output %s differs from expected", filename)
	}
}

// captures Stdout and returns output of function f()
func captureOutput(f func()) string {
	// redirect output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	// reset output again
	w.Close()
	os.Stdout = old

	captured, _ := ioutil.ReadAll(r)
	return string(captured)
}
