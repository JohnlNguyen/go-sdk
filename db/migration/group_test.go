package migration

import (
	"testing"

	"go-sdk/assert"
)

func TestNewGroup(t *testing.T) {
	assert := assert.New(t)

	g := Group(Step(Always(), NoOp))
	assert.Len(g, 1)
}
