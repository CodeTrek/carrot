package carrot_test

import (
	"fmt"
	"testing"

	"github.com/CodeTrek/carrot"
	"github.com/stretchr/testify/assert"
)

var f1 = func() int {
	return 1
}

var newF1 = func() int {
	return 2
}

var oldF1 = func() int { return 3 }

func TestSample(t *testing.T) {
	carrot.Patch(f1, newF1, oldF1)
	assert.Equal(t, 2, f1())
	assert.Equal(t, 1, oldF1())
	carrot.UnpatchAll()
	assert.Equal(t, 1, f1())
	assert.Equal(t, 3, oldF1())
}

var f2 = func(p1 [2000]byte) int {
	fmt.Printf("%s%s%s%s", p1[0:1], p1[0:1], p1[0:1], p1[0:1])
	return 1
}

var newF2 = func(p1 [2000]byte) int {
	return 2
}

var oldF2 = func(p1 [2000]byte) int { return 3 }

func TestComplex(t *testing.T) {
	var b = [2000]byte{0}
	assert.True(t, carrot.Patch(f2, newF2, oldF2))
	assert.Equal(t, 2, f2(b))
	carrot.UnpatchAll()
	assert.Equal(t, 1, f2(b))
}
