package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
)

var (
	// FIXME
	ignored = map[string]bool{
		"gotest-buildfailed.out": true,
		"gocheck-nofiles.out":    true,
	}

	xTimeRe = regexp.MustCompile(`run-date="[^"]+" run-time="[^"]+"`)
	xTime   = []byte(`run-date="2015-06-05" run-time="18:34:41"`)
)

type fixFunc func([]byte) []byte

func fixXUnit(in []byte) []byte {
	return xTimeRe.ReplaceAll(in, xTime)
}

func checkRegression(t *testing.T, inFile, outFile string, args []string, fixer fixFunc) {
	stdin, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("can't open %s - %s", inFile, err)
	}
	defer stdin.Close()
	cmd := exec.Command("./go2xunit", args...)
	cmd.Stdin = stdin
	out, err := cmd.Output()
	// TODO: Use -fail
	if err != nil {
		t.Fatalf("error running on %s - %s", inFile, err)
	}

	expected, err := ioutil.ReadFile(outFile)
	if err != nil {
		t.Fatalf("can't read %s - %s", outFile, err)
	}
	if fixer != nil {
		out = fixer(out)
		expected = fixer(expected)
	}

	if !bytes.Equal(out, expected) {
		diff := runDiff(out, expected)
		t.Fatalf("%s - output mismatch\n\n%s", inFile, diff)
	}

}

func runDiff(out, expected []byte) string {
	left, right := "/tmp/go2xunit-expected", "/tmp/go2xunit-out"
	ioutil.WriteFile(left, out, 0600)
	ioutil.WriteFile(right, expected, 0600)
	diffOut, _ := exec.Command("diff", left, right).Output()
	return string(diffOut)
}

const dataPath = "_data"

func iterCheck(t *testing.T, inType, outType string, args []string, fixer fixFunc) {
	pat := fmt.Sprintf(dataPath+"/in/%s*.out", inType)
	files, err := filepath.Glob(pat)
	if err != nil {
		t.Fatalf("no files matching %q", pat)
	}

	for _, inFile := range files {
		base := filepath.Base(inFile)
		if ignored[base] {
			continue
		}
		outFile := fmt.Sprintf(dataPath+"/out/%s/%s.xml", outType, base)
		name := fmt.Sprintf("%s-%s", outType, inFile)
		t.Run(name, func(t *testing.T) {
			checkRegression(t, inFile, outFile, args, fixer)
		})
	}
}

func TestRegression(t *testing.T) {
	cmd := exec.Command("go", "build")
	if err := cmd.Run(); err != nil {
		t.Fatalf("can't build - %s", err)
	}

	iterCheck(t, "gotest", "xunit", nil, nil)
	iterCheck(t, "gocheck", "xunit", []string{"-gocheck"}, nil)
	iterCheck(t, "gotest", "xunit.net", []string{"-xunitnet"}, fixXUnit)
	iterCheck(t, "gocheck", "xunit.net", []string{"-gocheck", "-xunitnet"}, fixXUnit)
}
