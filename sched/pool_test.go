package sched

import (
	"testing"

	"github.com/goesd/support/assert"
)

func TestPool(t *testing.T) {
	pool := newPool(10)

	pool.push(0, 2.5)
	pool.push(1, 1.0)
	pool.push(2, 2.0)

	assert.Equal(len(pool), 3, t)

	assert.Equal(pool.pop(), uint16(1), t)
	assert.Equal(pool.pop(), uint16(2), t)
	assert.Equal(pool.pop(), uint16(0), t)

	assert.Equal(len(pool), 0, t)
}
