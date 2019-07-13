package selector

import (
	"encoding/json"
	"testing"

	"go-sdk/assert"
)

func TestErrorJSON(t *testing.T) {
	// assert that the error can be serialized as json.
	assert := assert.New(t)

	testErr := Error("this is only a test")

	contents, err := json.Marshal(testErr)
	assert.Nil(err)
	assert.Equal("\"this is only a test\"", string(contents))
}
