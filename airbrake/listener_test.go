package airbrake

import (
	"fmt"
	"testing"

	"go-sdk/assert"
	"go-sdk/exception"
)

func TestFrames(t *testing.T) {
	assert := assert.New(t)

	ex := exception.New("this is a test")
	fr := frames(ex.Stack())
	assert.NotEmpty(fr, fmt.Sprintf("%T", ex.Stack()))

	fr = frames(exception.New("this is a test").WithStack(exception.StackStrings([]string{"foo", "bar"})).Stack())
	assert.Empty(fr)
}
