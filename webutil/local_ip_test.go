package webutil

import (
	"testing"

	"go-sdk/assert"
)

func TestLocalIP(t *testing.T) {
	assert := assert.New(t)

	assert.NotEmpty(LocalIP())
}
