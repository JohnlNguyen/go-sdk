package r2

import (
	"testing"

	"go-sdk/assert"
)

func TestOptBasicAuth(t *testing.T) {
	assert := assert.New(t)

	opt := OptBasicAuth("foo", "bar")

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Header)
	assert.NotEmpty(req.Header.Get("Authorization"))
}
