package main

import (
	//	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	//"path/filepath"
	//"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
// FIXME
/*
	ignored = map[string]bool{
		"gotest-buildfailed.out": true,
		"gocheck-nofiles.out":    true,
	}

	xTimeRe = regexp.MustCompile(`run-date="[^"]+" run-time="[^"]+"`)
	xTime   = []byte(`run-date="2015-06-05" run-time="18:34:41"`)
*/
)

/*
type fixFunc func([]byte) []byte

func fixXUnit(in []byte) []byte {
	return xTimeRe.ReplaceAll(in, xTime)
}
*/

func checkRegression(t *testing.T, inFile, outFile string, args []string) { // , fixer fixFunc) {
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
	/*
		if fixer != nil {
			out = fixer(out)
			expected = fixer(expected)
		}
	*/

	require.Equalf(expected, out, "%s - output mismatch", inFile)
}

/*
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
*/

type testCase struct {
	inFile  string
	outFile string
	args    []string
}

func genCases() []testCase {
	return []testCase{
		{"testdata/simple.json", "testdata/simple.junit.xml", nil},
	}
}

func TestRegression(t *testing.T) {
	cmd := exec.Command("go", "build")
	if err := cmd.Run(); err != nil {
		t.Fatalf("can't build - %s", err)
	}

	for _, tc := range genCases() {
		name := fmt.Sprintf("%+v", tc)
		t.Run(name, func(t *testing.T) {
			checkRegression(t, tc.inFile, tc.outFile, tc.args)
		})
	}

	/*
		iterCheck(t, "gotest", "xunit", nil, nil)
		iterCheck(t, "gocheck", "xunit", []string{"-gocheck"}, nil)
		iterCheck(t, "gotest", "xunit.net", []string{"-xunitnet"}, fixXUnit)
		iterCheck(t, "gocheck", "xunit.net", []string{"-gocheck", "-xunitnet"}, fixXUnit)
	*/
}
