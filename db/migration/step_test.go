package migration

import (
	"testing"

	"go-sdk/assert"
)

func TestStep(t *testing.T) {
	assert := assert.New(t)

	step := Step(Always(), NoOp)
	assert.NotNil(step.Guard)
	assert.NotNil(step.Body)
}
