package carrot_test

import (
	"fmt"
	"testing"

	"github.com/CodeTrek/carrot"
	"github.com/stretchr/testify/assert"
)

func f1() int {
	return 1
}

func newF1() int {
	return 2
}

var oldF1 = func() int { return 3 }

func TestSample(t *testing.T) {
	assert.True(t, carrot.Patch(f1, newF1, oldF1))
	assert.Equal(t, 2, f1())
	assert.Equal(t, 1, oldF1())
}

func f2(p1, p2, p3, p4 [2000]byte) int {
	fmt.Printf("%s%s%s%s", p1, p3, p3, p4)
	return 1
}

func newF2(p1, p2, p3, p4 [2000]byte) int {
	return 2
}

var oldF2 = func(p1, p2, p3, p4 [2000]byte) int { return 3 }

func TestComplex(t *testing.T) {
	assert.True(t, carrot.Patch(f2, newF2, oldF2))
}
