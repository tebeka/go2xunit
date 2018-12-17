package lib

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var (
	// gotest regular expressions

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
		"^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+((-?\\d+.\\d+)|\\(cached\\))")

	// ?       alipay  [no test files]
	gtNoFilesRE = regexp.MustCompile("^\\?.*\\[no test files\\]$")
	// FAIL    node/config [build failed]
	gtBuildFailedRE = regexp.MustCompile(`^FAIL.*\[(build|setup) failed\]$`)

	// exit status - 0
	gtExitRE = regexp.MustCompile("^exit status -?\\d+")

	// gocheck regular expressions

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
	// ok  	sisu.sh/go/code/catalog/transformer	(cached)
	gcSuiteRE = regexp.MustCompile("^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(-?\\d+.\\d+)")

	// Match error output from tests
	// Error Trace:    samples_test.go:27
	// Error:      	Should be true
	// Messages:   	Should be true
	gcTestErrorRE = regexp.MustCompile("(?s).*Error Trace:.*\\.go.*Error:.*Messages:.*")
)

// LineScanner scans lines and keep track of line numbers
type LineScanner struct {
	*bufio.Reader
	lnum int    // Current line number.
	text []byte // Content of current line of text.
	err  error  // Error from latest operation.
}

// NewLineScanner creates a new line scanner from r
func NewLineScanner(r io.Reader) *LineScanner {
	br := bufio.NewReader(r)
	ls := &LineScanner{
		Reader: br,
	}
	return ls
}

// Scan advances to next line.
func (ls *LineScanner) Scan() bool {
	if ls.text, _, ls.err = ls.Reader.ReadLine(); ls.err != nil {
		if ls.err == io.EOF {
			ls.err = nil
		}
		return false
	}
	ls.lnum++
	return true
}

// Text returns the current line
func (ls *LineScanner) Text() string {
	return string(ls.text)
}

// Err returns the current error (nil if no error)
func (ls *LineScanner) Err() error {
	return ls.err
}

// Line returns the current line number
func (ls *LineScanner) Line() int {
	return ls.lnum
}

// TokenType is the type of token
type TokenType int

// Token types
const (
	StartToken TokenType = iota
	EndToken
	SuiteToken
	NoFilesToken
	BuildFailedToken
	ExitToken
	DataToken
)

// Token is a lex token (line oriented)
type Token struct {
	Line int
	Type TokenType
	Data string
	//	Error error
}

func (typ TokenType) String() string {
	switch typ {
	case StartToken:
		return "Start"
	case EndToken:
		return "End"
	case SuiteToken:
		return "Suite"
	case NoFilesToken:
		return "NoFiles"
	case BuildFailedToken:
		return "BuildFailed"
	case DataToken:
		return "Data"
	case ExitToken:
		return "Exit"
	}

	return "???"
}

func (tok *Token) String() string {
	return fmt.Sprintf("<%s Token line:%03d>", tok.Type, tok.Line)
}

// Lexer generates tokens
type Lexer interface {
	Scan() bool
	Token() *Token
	Err() error
}

var gtTypes = []struct {
	typ TokenType
	re  *regexp.Regexp
}{
	{StartToken, gtStartRE},
	{EndToken, gtEndRE},
	{SuiteToken, gtSuiteRE},
	{NoFilesToken, gtNoFilesRE},
	{BuildFailedToken, gtBuildFailedRE},
	{ExitToken, gtExitRE},
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
	found := false
	for _, typ := range gtTypes {
		if typ.re.MatchString(line) {
			found = true
			lex.tok = &Token{lnum, typ.typ, line}
			break
		}
	}
	if !found {
		lex.tok = &Token{lnum, DataToken, line}
	}

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
