package carrot_test

import (
	"fmt"
	"testing"

	"github.com/CodeTrek/carrot"
	"github.com/stretchr/testify/assert"
)

var t1 = func() int { return 1 }
var r1 = func() int { return 2 }
var o1 = func() int { return 3 }

var t2 = func(p1 [2000]byte) int { fmt.Printf("%s%s%s%s", p1[0:1], p1[0:1], p1[0:1], p1[0:1]); return 1 }
var r2 = func(p1 [2000]byte) int { return 2 }
var o2 = func(p1 [2000]byte) int { return 3 }

var t3 = func() {}
var r3 = func() {}
var o3 = func() {}

var t4 = func() int { return 4 }
var r4 = func() int { fmt.Println(""); return 5 }
var o4 = func() int { return 6 }

func TestSample(t *testing.T) {
	carrot.Patch(t1, r1, o1)
	assert.Equal(t, 2, t1())
	assert.Equal(t, 1, o1())
	carrot.UnpatchAll()
	assert.Equal(t, 1, t1())
	assert.Equal(t, 3, o1())
}

func TestComplex(t *testing.T) {
	var b = [2000]byte{0}
	assert.True(t, carrot.Patch(t2, r2, o2))
	assert.Equal(t, 2, t2(b))
	carrot.UnpatchAll()
	assert.Equal(t, 1, t2(b))
}

func TestMulti(t *testing.T) {
	assert.True(t, carrot.Patch(t3, r3, o3))
	assert.False(t, carrot.Patch(t3, r3, o3))
	assert.Panics(t, func() { carrot.Patch(t1, r2, o1) })
	assert.True(t, carrot.Patch(t2, r2, o2))

	assert.True(t, carrot.Patch(t1, r1, o1))
	assert.Equal(t, 2, t1())
	assert.Equal(t, 2, r1())
	assert.Equal(t, 1, o1())

	assert.True(t, carrot.Patch(t4, r1, o4))

	assert.Equal(t, 2, t4())
	assert.Equal(t, 2, r1())
	assert.Equal(t, 4, o4())

	carrot.Unpatch(t1)

	assert.Equal(t, 1, t1())
	assert.Equal(t, 2, r1())
	assert.Equal(t, 3, o1())
}
