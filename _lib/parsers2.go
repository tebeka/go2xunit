package lib

import "io"

// Parser generates suites
type Parser interface {
	Scan() bool
	Suite() *Suite
	Err() error
}

// GtParser is a gotest output parser
type GtParser struct {
	suite *Suite
	tests map[string]*Test
	lex   Lexer
}

// Scan scans for the next suite
func (gtp *GtParser) Scan() bool {
	return true
}

// Suite is the current suite
func (gtp *GtParser) Suite() *Suite {
	return nil
}

// Err is the current error
func (gtp *GtParser) Err() error {
	return nil
}

// NewGtParser return a new gotest parser
func NewGtParser(in io.Reader) Parser {
	return &GtParser{
		suite: nil,
		tests: make(map[string]*Test),
		lex:   NewGotestLexer(in),
	}
}
