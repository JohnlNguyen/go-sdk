package async

import (
	"testing"

	"go-sdk/assert"
)

func TestExec(t *testing.T) {
	assert := assert.New(t)

	didRun := make(chan struct{})
	errors := Exec(func() error {
		defer close(didRun)
		return nil
	})
	<-didRun
	assert.Empty(errors)
}
