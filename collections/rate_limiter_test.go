package collections

import (
	"testing"
	"time"

	"go-sdk/assert"
)

func TestRateLimiter(t *testing.T) {
	it := assert.New(t)

	rl := NewRateLimiter(5, 1*time.Millisecond)

	it.False(rl.Check("a"))
	it.False(rl.Check("b"))
	it.False(rl.Check("b"))
	it.False(rl.Check("b"))
	it.False(rl.Check("b"))
	it.False(rl.Check("a"))
	it.False(rl.Check("a"))
	it.False(rl.Check("a"))
	it.True(rl.Check("a"))

	time.Sleep(1 * time.Millisecond)

	it.False(rl.Check("a"))
	it.False(rl.Check("b"))
	it.False(rl.Check("b"))
	it.False(rl.Check("b"))
	it.False(rl.Check("b"))
	it.False(rl.Check("a"))
	it.False(rl.Check("a"))
	it.False(rl.Check("a"))
	it.True(rl.Check("a"))
}
