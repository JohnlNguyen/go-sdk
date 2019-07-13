package sh

import (
	"bytes"
	"testing"

	"go-sdk/assert"
)

func TestPromptFrom(t *testing.T) {
	assert := assert.New(t)

	input := bytes.NewBufferString("test\n")
	output := bytes.NewBuffer(nil)
	assert.Equal("test", PromptFrom(output, input, "value: "))
}
