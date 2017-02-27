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

	for gtp.lex.Scan() {
		tok := gtp.lex.Token()
		switch tok.Type {
		case StartToken:
			curTest = &Test{Name: tok.Fields[1]}
			byName[curTest.Name] = curTest
			tests = append(tests, curTest)
		case EndToken:
			status, name, time := tok.Fields[1], tok.Fields[2], tok.Fields[3]
			test, ok := byName[name]
			if !ok {
				gtp.err = fmt.Errorf("%d: unknown test %s", tok.Line, name)
				return false
			}
			test.Status = Token2Status(status)
			test.Time = time
		case DataToken:
			if curTest == nil {
				gtp.err = fmt.Errorf("%d: orphan data", tok.Line)
				return false
			}
			if len(curTest.Message) == 0 {
				curTest.Message = tok.Data
			} else {
				curTest.Message += ("\n" + tok.Data)
			}
		case ErrorToken:
			curTest.Message += tok.Data + "\n"
			curTest.Status = Failed
			curTest.Time = "N/A"
		case SuiteToken:
			status, name, time := tok.Fields[1], tok.Fields[2], tok.Fields[3]
			gtp.suite = &Suite{
				Name:   name,
				Time:   time,
				Status: status,
				Tests:  tests,
			}
			break
		}
	}
	if gtp.lex.Err() != nil {
		return false
	}

	// Missing summary
	if gtp.suite == nil && len(tests) > 0 {
		gtp.suite = &Suite{
			Name:   tests[0].Name,
			Time:   tests[0].Time,
			Status: "Unkonwn",
			Tests:  tests,
		}
	}

	return len(tests) > 0
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

// GetParser return parse for inType
func GetParser(inType string, in io.Reader) (Parser, error) {
	switch inType {
	case "gotest":
		return NewGtParser(in), nil
	}

	return nil, fmt.Errorf("unknown parser type - %s", inType)
}
