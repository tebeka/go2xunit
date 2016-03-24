package mmath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	x, y := 1, 2
	z := Add(x, y)
	assert.Equal(t, z, x+y, "%d + %d != %d\n", x, y, x+y)
}

func TestSub(t *testing.T) {
	x, y := 1, 2
	z := Sub(x, y)
	assert.Equal(t, z, x-y, "%d-%d != %d\n", x, y, x-y)
}

func TestMul(t *testing.T) {
	x, y := 2, 3
	z := Mul(x, y)
	assert.Equal(t, z, x*y, "%d*%d != %d\n", x, y, x*y)
}

func TestDiv(t *testing.T) {
	x, y := 2, 3
	z := Div(x, y)
	quot := float64(x) / float64(y)
	assert.Equal(t, float64(z), quot, "%d/%d != %f\n", x, y, quot)
}
