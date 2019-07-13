package web

import (
	"context"
	"net/http"
	"testing"

	"go-sdk/assert"
	"go-sdk/stringutil"
)

func TestSessionAware(t *testing.T) {
	assert := assert.New(t)

	sessionID := NewSessionID()

	var didExecuteHandler bool
	var sessionWasSet bool

	app := New().WithAuth(NewLocalAuthManager())
	app.Auth().PersistHandler()(context.TODO(), &Session{SessionID: sessionID, UserID: "bailey"}, nil)

	app.GET("/", func(r *Ctx) Result {
		didExecuteHandler = true
		sessionWasSet = r.Session() != nil
		return r.Text().Result("COOL")
	}, SessionAware)

	meta, err := app.Mock().WithPathf("/").WithCookieValue(app.Auth().CookieName(), sessionID).ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(ContentTypeText, meta.Headers.Get(HeaderContentType))
	assert.True(didExecuteHandler, "we should have triggered the hander")
	assert.True(sessionWasSet, "the session should have been set by the middleware")

	unsetMeta, err := app.Mock().WithPathf("/").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, unsetMeta.StatusCode)
	assert.False(sessionWasSet)
}

func TestSessionRequired(t *testing.T) {
	assert := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	app := New().WithAuth(NewLocalAuthManager())
	app.Auth().PersistHandler()(context.TODO(), &Session{SessionID: sessionID, UserID: "bailey"}, nil)

	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session() != nil
		return r.Text().Result("COOL")
	}, SessionRequired)

	unsetMeta, err := app.Mock().WithPathf("/").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusForbidden, unsetMeta.StatusCode)
	assert.False(sessionWasSet)

	meta, err := app.Mock().WithPathf("/").WithCookieValue(app.Auth().CookieName(), sessionID).ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(sessionWasSet)
}

func TestSessionRequiredCustomParamName(t *testing.T) {
	assert := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	app := New().WithAuth(NewLocalAuthManager())
	app.Auth().PersistHandler()(context.TODO(), &Session{SessionID: sessionID, UserID: "bailey"}, nil)
	app.Auth().WithCookieName("web_auth")

	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session() != nil
		return r.Text().Result("COOL")
	}, SessionRequired)

	unsetMeta, err := app.Mock().WithPathf("/").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusForbidden, unsetMeta.StatusCode)
	assert.False(sessionWasSet)

	meta, err := app.Mock().WithPathf("/").WithCookieValue(app.Auth().CookieName(), sessionID).ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(sessionWasSet)

	meta, err = app.Mock().WithPathf("/").WithCookieValue(DefaultCookieName, sessionID).ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusForbidden, meta.StatusCode)
	assert.True(sessionWasSet)
}

func TestSessionMiddleware(t *testing.T) {
	assert := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	app := New().WithAuth(NewLocalAuthManager())
	app.Auth().PersistHandler()(context.TODO(), &Session{SessionID: sessionID, UserID: "bailey"}, nil)

	var calledCustom bool
	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session() != nil
		return r.Text().Result("COOL")
	}, SessionMiddleware(func(_ *Ctx) Result {
		calledCustom = true
		return NoContent
	}))

	unsetMeta, err := app.Mock().WithPathf("/").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, unsetMeta.StatusCode)
	assert.False(sessionWasSet)

	meta, err := app.Mock().WithPathf("/").WithCookieValue(app.Auth().CookieName(), sessionID).ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(sessionWasSet)

	assert.True(calledCustom)
}
