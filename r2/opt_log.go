package r2

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"go-sdk/logger"
)

const (
	// maxLogBytes is the maximum number of bytes to log from a response.
	// it is currently set to 1mb.
	maxLogBytes = 1 << 20
)

// OptLogRequest adds an OnResponse listener to log the response of a call.
func OptLogRequest(log logger.Log) Option {
	return OptOnRequest(func(req *http.Request) error {
		event := NewEvent(Flag,
			OptEventRequest(req))
		log.Trigger(event)
		return nil
	})
}

// OptLogResponse adds an OnResponse listener to log the response of a call.
func OptLogResponse(log logger.Log) Option {
	return OptOnResponse(func(req *http.Request, res *http.Response, started time.Time, err error) error {
		if err != nil {
			return err
		}
		event := NewEvent(FlagResponse,
			OptEventStarted(started),
			OptEventRequest(req),
			OptEventResponse(res))

		log.Trigger(event)
		return nil
	})
}

// OptLogResponseWithBody adds an OnResponse listener to log the response of a call.
// It reads the contents of the response fully before emitting the event.
// Do not use this if the size of the responses can be large.
func OptLogResponseWithBody(log logger.Log) Option {
	return OptOnResponse(func(req *http.Request, res *http.Response, started time.Time, err error) error {
		if err != nil {
			return err
		}
		defer res.Body.Close()

		// read out the buffer in full
		buffer := new(bytes.Buffer)
		if _, err := io.Copy(buffer, res.Body); err != nil {
			return err
		}
		// set the body to the read contents
		res.Body = ioutil.NopCloser(bytes.NewReader(buffer.Bytes()))

		event := NewEvent(FlagResponse,
			OptEventStarted(started),
			OptEventRequest(req),
			OptEventResponse(res),
			OptEventBody(buffer.Bytes()))

		log.Trigger(event)
		return nil
	})
}
