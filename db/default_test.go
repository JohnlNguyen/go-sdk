package db

import (
	"testing"

	"go-sdk/assert"
)

func TestDefault(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(Default())
}
