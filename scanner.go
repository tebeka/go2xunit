package main

import (
	"bufio"
	"io"
)

// Scanner if line scanner with line numbers
type Scanner struct {
	scan *bufio.Scanner
	lnum int
}

// NewScanner returns a new scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scan: bufio.NewScanner(r),
		lnum: 0,
	}
}

// Scan scans for the next line
func (s *Scanner) Scan() bool {
	s.lnum++
	return s.scan.Scan()
}

// Err return last error
func (s *Scanner) Err() error {
	return s.scan.Err()
}

// Bytes return last line
func (s *Scanner) Bytes() []byte {
	return s.scan.Bytes()
}

// LineNum returns the current line number
func (s *Scanner) LineNum() int {
	return s.lnum
}
