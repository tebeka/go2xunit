package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func checkRegression(t *testing.T, inFile, outFile string, args []string) {
	require := require.New(t)
	stdin, err := os.Open(inFile)
	require.NoErrorf(err, "open %s - %s", inFile, err)
	defer stdin.Close()

	cmd := exec.Command("./go2xunit", args...)
	cmd.Stdin = stdin
	data, err := cmd.Output()
	require.NoErrorf(err, "running on %s - %s", inFile, err)
	out := strings.TrimSpace(string(data))

	data, err = ioutil.ReadFile(outFile)
	require.NoErrorf(err, "read %s - %s", outFile, err)
	expected := strings.TrimSpace(string(data))

	require.Equalf(expected, out, "%s - output mismatch", inFile)
}

type testCase struct {
	inFile  string
	outFile string
	args    []string
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func genCases(t *testing.T) []testCase {
	require := require.New(t)
	var cases []testCase
	ext := ".json"
	files, err := filepath.Glob("testdata/*" + ext)
	require.NoError(err, "glob")
	for _, inFile := range files {
		for _, typ := range []string{"junit", "xunit"} {
			outFile := inFile[:len(inFile)-len(ext)]
			outFile += "." + typ + ".xml"

			if !fileExists(outFile) {
				continue
			}
			tc := testCase{inFile, outFile, nil}
			cases = append(cases, tc)
		}
	}

	return cases
}

func TestRegression(t *testing.T) {
	for _, tc := range genCases(t) {
		name := fmt.Sprintf("%+v", tc)
		t.Run(name, func(t *testing.T) {
			checkRegression(t, tc.inFile, tc.outFile, tc.args)
		})
	}
}
