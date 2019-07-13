package async

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"

	"go-sdk/assert"
)

func TestBatch(t *testing.T) {
	assert := assert.New(t)

	var items []interface{}
	for x := 0; x < 32; x++ {
		items = append(items, "hello"+strconv.Itoa(x))
	}

	var processed int32
	errors := make(chan error, 32)
	b := NewBatch(func(_ context.Context, v interface{}) error {
		atomic.AddInt32(&processed, 1)
		return fmt.Errorf("this is only a test")
	}, items...).WithErrors(errors)

	b.ProcessContext(context.Background())

	assert.Equal(32, processed)
	assert.Equal(32, len(errors))
}
