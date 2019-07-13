package db

import (
	"testing"

	"go-sdk/assert"
	"go-sdk/exception"
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(Error(nil))

	var err error
	assert.Nil(Error(err))

	err = exception.New("this is only a test")
	assert.True(exception.Is(Error(err), exception.Class("this is only a test")))
}
