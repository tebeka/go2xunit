package lib

import (
	"fmt"
	"io"
)

// Parser generates suites
type Parser interface {
	Scan() bool
	Suite() *Suite
	Err() error
}

// GtParser is a gotest output parser
type GtParser struct {
	suite *Suite
	lex   Lexer
	err   error
}

// Scan scans for the next suite
func (gtp *GtParser) Scan() bool {
	var tests []*Test
	byName := make(map[string]*Test)
	var curTest *Test
	hasTest := false

	for gtp.lex.Scan() {
		hasTest = true
		tok := gtp.lex.Token()
		switch tok.Type {
		case StartToken:
			curTest = &Test{Name: tok.Fields[1]}
			byName[curTest.Name] = curTest
			tests = append(tests, curTest)
		case EndToken:
			status, name, time := tok.Fields[1], tok.Fields[2], tok.Fields[2]
			test, ok := byName[name]
			if !ok {
				gtp.err = fmt.Errorf("%d: unknown test %s", tok.Line, name)
				return false
			}
			test.Status = Token2Status(status)
			test.Time = time
		case DataToken:
			curTest.Message += tok.Data
		case SuiteToken:
			status, name, time := tok.Fields[1], tok.Fields[2], tok.Fields[3]
			gtp.suite = &Suite{
				Name:   name,
				Time:   time,
				Status: status,
				Tests:  tests,
			}
			return true
		}
	}
	if gtp.lex.Err() != nil {
		return false
	}

	return hasTest
}

// Suite is the current suite
func (gtp *GtParser) Suite() *Suite {
	return gtp.suite
}

// Err is the current error
func (gtp *GtParser) Err() error {
	if err := gtp.lex.Err(); err != nil {
		return err
	}
	return gtp.err
}

// NewGtParser return a new gotest parser
func NewGtParser(in io.Reader) Parser {
	return &GtParser{
		lex: NewGotestLexer(in),
	}
}
