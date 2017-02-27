package lib

import (
	"os"
	"testing"
)

// TODO: Regression (out file with list of tokens)

var testCases = []struct {
	fname    string
	expected []string
}{
	{
		"../data/in/gotest-1.7.out",
		[]string{
			StartToken,
			EndToken,
			StartToken,
			EndToken,
			StartToken,
			EndToken,
			StartToken,
			EndToken,
			DataToken,
			StartToken,
			StartToken,
			StartToken,
			EndToken,
			EndToken,
			EndToken,
			SuiteDelimToken,
			ExitToken,
			SuiteToken,
		},
	},
	{
		"../data/in/gotest-panic.out",
		[]string{
			StartToken,
			ErrorToken,
			DataToken,
			ExitToken,
			SuiteToken,
		},
	},
}

func CheckLexer(t *testing.T, fname string, expected []string) {
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	lex := NewGotestLexer(file)
	for lex.Scan() {
		out = append(out, lex.Token().Type)
	}
	if lex.Err() != nil {
		t.Fatalf("can't scan - %s", lex.Err())
	}

	if len(out) != len(expected) {
		t.Fatalf("wrong number of tokens (%d != %d)", len(out), len(expected))
	}
	for i, tt := range expected {
		if tt != out[i] {
			t.Fatalf("wrong type at %d %s != %s)", i, out[i], tt)
		}
	}

}

func TestGotestLexer(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.fname, func(t *testing.T) { CheckLexer(t, tc.fname, tc.expected) })
	}
}
