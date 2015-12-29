package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"
)

func Test_MainXunitNet(t *testing.T) {
	cmd := exec.Command("go", "build")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("can't build - %s", err)
	}

	base := "gocheck-pass.out"
	infile := fmt.Sprintf("data/%s", base)
	cmd = exec.Command("./go2xunit", "-gocheck", "-input", infile)
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		t.Fatalf("can't run - %s", err)
	}
	actual := out.Bytes()

	expected, err := ioutil.ReadAll(getOutputData("xunit", base))
	checkError(err)

	if !bytes.Equal(expected, actual) {
		t.Errorf("xUnit.net XML output %s differs from expected", base)
	}
}
