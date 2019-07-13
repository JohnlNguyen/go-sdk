package exception

import (
	"testing"

	"go-sdk/assert"
)

func TestIs(t *testing.T) {
	assert := assert.New(t)

	errInvalidSomething := Class("invalid something")

	ex := New(errInvalidSomething)

	assert.True(Is(ex, errInvalidSomething))
	assert.True(Is(errInvalidSomething, errInvalidSomething))
}
