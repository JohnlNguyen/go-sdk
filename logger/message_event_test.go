package logger

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMessageEventListener(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)

	textBuffer := bytes.NewBuffer(nil)
	jsonBuffer := bytes.NewBuffer(nil)
	all := New().WithFlags(AllFlags()).WithRecoverPanics(false).
		WithWriter(NewTextWriter(textBuffer)).
		WithWriter(NewJSONWriter(jsonBuffer))

	defer all.Close()
	all.Listen(Flag("test-flag"), "default", NewMessageEventListener(func(e *MessageEvent) {
		defer wg.Done()
		assert.Equal("test-flag", e.Flag())
		assert.Equal("foo bar", e.Message())
	}))

	go func() { all.Trigger(Messagef(Flag("test-flag"), "foo %s", "bar")) }()
	go func() { all.Trigger(Messagef(Flag("test-flag"), "foo %s", "bar")) }()
	wg.Wait()
	all.Drain()

	assert.NotEmpty(textBuffer.String())
	assert.NotEmpty(jsonBuffer.String())
}

func TestMessageEventInterfaces(t *testing.T) {
	assert := assert.New(t)

	ee := Messagef(Info, "this is a test").
		WithHeading("heading").
		WithLabel("foo", "bar")

	eventProvider, isEvent := MarshalEvent(ee)
	assert.True(isEvent)
	assert.Equal(Info, eventProvider.Flag())
	assert.False(eventProvider.Timestamp().IsZero())

	headingProvider, isHeadingProvider := MarshalEventHeading(ee)
	assert.True(isHeadingProvider)
	assert.Equal("heading", headingProvider.Heading())

	metaProvider, isMetaProvider := MarshalEventMeta(ee)
	assert.True(isMetaProvider)
	assert.Equal("bar", metaProvider.Labels()["foo"])
}

func TestMessageEventProperties(t *testing.T) {
	assert := assert.New(t)

	e := Messagef(Info, "")

	assert.False(e.Timestamp().IsZero())
	assert.True(e.WithTimestamp(time.Time{}).Timestamp().IsZero())

	assert.Empty(e.Labels())
	assert.Equal("bar", e.WithLabel("foo", "bar").Labels()["foo"])

	assert.Empty(e.Annotations())
	assert.Equal("zar", e.WithAnnotation("moo", "zar").Annotations()["moo"])

	assert.Equal(Info, e.Flag())
	assert.Equal(Error, e.WithFlag(Error).Flag())

	assert.Empty(e.Heading())
	assert.Equal("Heading", e.WithHeading("Heading").Heading())

	assert.Empty(e.Message())
	assert.Equal("Message", e.WithMessage("Message").Message())

	assert.Empty(e.FlagTextColor())
	assert.Equal(ColorWhite, e.WithFlagTextColor(ColorWhite).FlagTextColor())
}
