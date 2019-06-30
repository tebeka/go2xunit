package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	require := require.New(t)
	lines := []string{
		"line one",
		"line two",
		"line three",
	}

	r := strings.NewReader(strings.Join(lines, "\n"))
	s := NewScanner(r)
	var i int
	for s.Scan() {
		i = s.LineNum() - 1
		require.Equalf(lines[i], string(s.Bytes()), "%d", i)
	}

	require.Equal(len(lines), i+1, "total lines")
}
