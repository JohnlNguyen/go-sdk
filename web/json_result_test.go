package web

import (
	"bytes"
	"net/http"
	"testing"

	"go-sdk/assert"
	"go-sdk/webutil"
)

func TestJSONResultRender(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	w := webutil.NewMockResponse(buf)
	r := NewCtx(w, webutil.NewMockRequest("GET", "/"))

	jr := &JSONResult{
		StatusCode: http.StatusOK,
		Response: map[string]interface{}{
			"foo": "bar",
		},
	}

	assert.Nil(jr.Render(r))
	assert.Equal(http.StatusOK, w.StatusCode())
	assert.Equal("{\"foo\":\"bar\"}\n", buf.String())
}

func TestJSONResultRenderStatusCode(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	w := webutil.NewMockResponse(buf)
	r := NewCtx(w, webutil.NewMockRequest("GET", "/"))

	jr := &JSONResult{
		StatusCode: http.StatusBadRequest,
		Response: map[string]interface{}{
			"foo": "bar",
		},
	}

	assert.Nil(jr.Render(r))
	assert.Equal(http.StatusBadRequest, w.StatusCode())
	assert.Equal("{\"foo\":\"bar\"}\n", buf.String())
}
