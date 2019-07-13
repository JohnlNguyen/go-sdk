package async

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"go-sdk/assert"
	"go-sdk/graceful"
)

// Assert a latch is graceful
var (
	_ graceful.Graceful = (*AutoflushBuffer)(nil)
)

func TestAutoflushBuffer(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)
	buffer := NewAutoflushBuffer(10, time.Hour).WithFlushHandler(func(objects []interface{}) {
		defer wg.Done()
		assert.Len(objects, 10)
	})

	buffer.Start()
	defer buffer.Stop()

	for x := 0; x < 20; x++ {
		buffer.Add(fmt.Sprintf("foo%d", x))
	}

	wg.Wait()
}

func TestAutoflushBufferTicker(t *testing.T) {
	assert := assert.New(t)
	assert.StartTimeout(500 * time.Millisecond)
	defer assert.EndTimeout()

	wg := sync.WaitGroup{}
	wg.Add(20)
	buffer := NewAutoflushBuffer(100, time.Millisecond).WithFlushHandler(func(objects []interface{}) {
		for range objects {
			wg.Done()
		}
	})

	buffer.Start()
	defer buffer.Stop()

	for x := 0; x < 20; x++ {
		buffer.Add(fmt.Sprintf("foo%d", x))
	}
	wg.Wait()
}

func BenchmarkAutoflushBuffer(b *testing.B) {
	buffer := NewAutoflushBuffer(128, 500*time.Millisecond).WithFlushHandler(func(objects []interface{}) {
		if len(objects) > 128 {
			b.Fail()
		}
	})

	buffer.Start()
	defer buffer.Stop()

	for x := 0; x < b.N; x++ {
		for y := 0; y < 1000; y++ {
			buffer.Add(fmt.Sprintf("asdf%d%d", x, y))
		}
	}
}
