package sh

import (
	"testing"

	"go-sdk/assert"
	"go-sdk/exception"
)

func TestParseFlagsTrailer(t *testing.T) {
	assert := assert.New(t)

	parsed, err := ArgsTrailer("foo", "bar")
	assert.True(exception.Is(err, ErrFlagsNoTrailer))
	assert.Empty(parsed)

	parsed, err = ArgsTrailer("foo", "bar", "--")
	assert.True(exception.Is(err, ErrFlagsNoTrailer))
	assert.Empty(parsed)

	parsed, err = ArgsTrailer("foo", "bar", "--", "echo", "'things'")
	assert.Nil(err)
	assert.Equal([]string{"echo", "'things'"}, parsed)
}
