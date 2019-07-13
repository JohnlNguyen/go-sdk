package webutil

import (
	"testing"

	"go-sdk/assert"
)

func TestGetUseragent(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("go-sdk test", GetUserAgent(NewMockRequest("GET", "/")))
}
