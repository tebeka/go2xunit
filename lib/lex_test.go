package lib

import (
	"os"
	"testing"
)

// TODO: Regression (out file with list of tokens)

func TestGotestLexer(t *testing.T) {
	fname := "../data/in/gotest-1.7.out"
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}

	expected := []TokenType{
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
		DataToken,
		ExitToken,
		SuiteToken,
	}

	var out []TokenType
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
