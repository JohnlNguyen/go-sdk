package web

import (
	"bytes"
	"testing"

	"go-sdk/assert"
	"go-sdk/webutil"
)

func TestDefaultProviderMiddlewares(t *testing.T) {
	assert := assert.New(t)

	r := applyMiddleware(JSONProviderAsDefault)
	_, ok := r.DefaultResultProvider().(JSONResultProvider)
	assert.True(ok)

	r = applyMiddleware(ViewProviderAsDefault)
	_, ok = r.DefaultResultProvider().(*ViewCache)
	assert.True(ok)

	r = applyMiddleware(XMLProviderAsDefault)
	_, ok = r.DefaultResultProvider().(XMLResultProvider)
	assert.True(ok)

	r = applyMiddleware(TextProviderAsDefault)
	_, ok = r.DefaultResultProvider().(TextResultProvider)
	assert.True(ok)
}

func applyMiddleware(middleware Middleware) (output *Ctx) {
	middleware(func(ctx *Ctx) Result {
		output = ctx
		return NoContent
	})(NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/")))
	return
}
