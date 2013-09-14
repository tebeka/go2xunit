package mmath

import (
	. "launchpad.net/gocheck"
	"testing"
)

type MySuite struct{}
var _ = Suite(&MySuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

func (s *MySuite) TestAdd(c *C) {
	x, y := 1, 2
	z := Add(x, y)
	c.Assert(z, Equals, x+y)
}

func (s *MySuite) TestSub(c *C) {
	x, y := 1, 2
	z := Sub(x, y)

	c.Assert(z, Equals, x-y)
}

func (s *MySuite) TestMul(c *C) {
	x, y := 2, 3
	z := Mul(x, y)
	c.Assert(z, Equals, x*y)
}

func (s *MySuite) TestDiv(c *C) {
	x, y := 2, 3
	z := Div(x, y)
	c.Assert(z, Equals, float64(x)/float64(y))
}
