package lib

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

// gotest regular expressions
var (

	// === RUN TestAdd
	gtStartRE = regexp.MustCompile(
		"^=== RUN:?[[:space:]]+([a-zA-Z_][^[:space:]]*)")

	// --- PASS: TestSub (0.00 seconds)
	// --- FAIL: TestSubFail (0.00 seconds)
	// --- SKIP: TestSubSkip (0.00 seconds)
	gtEndRE = regexp.MustCompile(
		"--- (PASS|FAIL|SKIP):[[:space:]]+([a-zA-Z_][^[:space:]]*) " +
			"\\((-?\\d+(.\\d+)?)")

	// FAIL	_/home/miki/Projects/goroot/src/xunit	0.004s
	// ok  	_/home/miki/Projects/goroot/src/anotherTest	0.000s
	gtSuiteRE = regexp.MustCompile(
		"^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(-?\\d+.\\d+)")

	// ?       alipay  [no test files]
	gtNoFilesRE = regexp.MustCompile("^\\?.*\\[no test files\\]$")
	// FAIL    node/config [build failed]
	gtBuildFailedRE = regexp.MustCompile(`^FAIL.*\[(build|setup) failed\]$`)

	// exit status - 0
	gtExitRE = regexp.MustCompile("^exit status -?\\d+")

	// fatal error: all goroutines are asleep - deadlock!
	gtErrorRE = regexp.MustCompile("^fatal error: ")

	// PASS
	// FAIL
	gtSuiteDelimRE = regexp.MustCompile("^(PASS|FAIL)$")

	// testing: warning: no tests to run
	gtWarningRE = regexp.MustCompile("^testing: warning: .*")
)

// gocheck regular expressions
var (
	// START: mmath_test.go:16: MySuite.TestAdd
	gcStartRE = regexp.MustCompile(
		"START: [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)")

	// PASS: mmath_test.go:16: MySuite.TestAdd	0.000s
	// FAIL: mmath_test.go:35: MySuite.TestDiv
	gcEndRE = regexp.MustCompile(
		"(PASS|FAIL|SKIP|PANIC|MISS): [^:]+:[^:]+: " +
			"([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)" +
			"[[:space:]]?(-?[0-9]+.[0-9]+)?")

	// FAIL	go2xunit/demo-gocheck	0.008s
	// ok  	go2xunit/demo-gocheck	0.008s
	gcSuiteRE = regexp.MustCompile("^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(-?\\d+.\\d+)")
)

// LineScanner scans lines and keep track of line numbers
type LineScanner struct {
	*bufio.Scanner
	lnum int
}

// NewLineScanner creates a new line scanner from r
func NewLineScanner(r io.Reader) *LineScanner {
	scan := bufio.NewScanner(r)
	return &LineScanner{scan, 0}
}

// Scan advances to next line
func (ls *LineScanner) Scan() bool {
	val := ls.Scanner.Scan()
	ls.lnum++
	return val
}

// Line returns the current line number
func (ls *LineScanner) Line() int {
	return ls.lnum
}

// Token types
const (
	StartToken       = "StartToken"
	EndToken         = "EndToken"
	SuiteToken       = "SuiteToken"
	NoFilesToken     = "NoFilesToken"
	BuildFailedToken = "BuildFailedToken"
	ExitToken        = "ExitToken"
	ErrorToken       = "ErrorToken"
	DataToken        = "DataToken"
	SuiteDelimToken  = "SuiteDelimToken"
	WarningToken     = "WarningToken"
)

// Token is a lex token (line oriented)
type Token struct {
	Line   int
	Type   string
	Data   string
	Fields []string
	//	Error error
}

func (tok *Token) String() string {
	return fmt.Sprintf("<%s Token line:%d>", tok.Type, tok.Line)
}

// Lexer generates tokens
type Lexer interface {
	Scan() bool
	Token() *Token
	Err() error
}

var gtTypes = []struct {
	typ string
	re  *regexp.Regexp
}{
	{StartToken, gtStartRE},
	{EndToken, gtEndRE},
	{SuiteToken, gtSuiteRE},
	{NoFilesToken, gtNoFilesRE},
	{BuildFailedToken, gtBuildFailedRE},
	{ExitToken, gtExitRE},
	{ErrorToken, gtErrorRE},
	{SuiteDelimToken, gtSuiteDelimRE},
	{WarningToken, gtWarningRE},
}

// GotestLexer is a lexer for gotest output
type GotestLexer struct {
	scanner *LineScanner
	tok     *Token
	err     error
}

// Scan scans the next token. Return true if there's a new one
func (lex *GotestLexer) Scan() bool {
	if !lex.scanner.Scan() {
		return false
	}
	if lex.scanner.Err() != nil {
		return false
	}

	line := lex.scanner.Text()
	lnum := lex.scanner.Line()
	for _, typ := range gtTypes {
		fields := typ.re.FindStringSubmatch(line)
		if len(fields) > 0 {
			lex.tok = &Token{lnum, typ.typ, line, fields}
			return true
		}
	}

	lex.tok = &Token{lnum, DataToken, line, nil}
	return true
}

// Token returns the current Token
func (lex *GotestLexer) Token() *Token {
	return lex.tok
}

// Err returns the current error (if any)
func (lex *GotestLexer) Err() error {
	if lex.err != nil {
		return lex.err
	}
	return lex.scanner.Err()
}

// NewGotestLexer returns a new lexer for gotest files
func NewGotestLexer(in io.Reader) Lexer {
	return &GotestLexer{
		scanner: NewLineScanner(in),
		tok:     nil,
		err:     nil,
	}
}
