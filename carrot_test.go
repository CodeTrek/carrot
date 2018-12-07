package carrot_test

import (
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
