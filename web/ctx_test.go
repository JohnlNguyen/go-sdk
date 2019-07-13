package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"go-sdk/assert"
)

func TestCtxGetState(t *testing.T) {
	assert := assert.New(t)

	context := NewCtx(nil, nil)
	context.WithStateValue("foo", "bar")
	assert.Equal("bar", context.StateValue("foo"))
}

func TestCtxParamQuery(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).WithQueryString("foo", "bar").CreateCtx(nil)
	assert.Nil(err)
	param, err := context.Param("foo")
	assert.Nil(err)
	assert.Equal("bar", param)

	param, err = context.QueryValue("foo")
	assert.Nil(err)
	assert.Equal("bar", param)
}

func TestCtxParamHeader(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).WithHeader("foo", "bar").CreateCtx(nil)
	assert.Nil(err)
	param, err := context.Param("foo")
	assert.Nil(err)
	assert.Equal("bar", param)

	param, err = context.HeaderValue("foo")
	assert.Nil(err)
	assert.Equal("bar", param)
}

func TestCtxParamForm(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).WithFormValue("foo", "bar").CreateCtx(nil)
	assert.Nil(err)
	param, err := context.Param("foo")
	assert.Nil(err)
	assert.Equal("bar", param)

	param, err = context.FormValue("foo")
	assert.Nil(err)
	assert.Equal("bar", param)
}

func TestCtxParamCookie(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).WithCookie(&http.Cookie{Name: "foo", Value: "bar"}).CreateCtx(nil)
	assert.Nil(err)
	param, err := context.Param("foo")
	assert.Nil(err)
	assert.Equal("bar", param)
}

func TestCtxPostBodyAsString(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).WithPostBody([]byte("test payload")).CreateCtx(nil)
	assert.Nil(err)
	body, err := context.PostBodyAsString()
	assert.Nil(err)
	assert.Equal("test payload", body)

	context, err = NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)
	body, err = context.PostBodyAsString()
	assert.Nil(err)
	assert.Empty(body)
}

func TestCtxPostBodyAsJSON(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).WithPostBody([]byte(`{"test":"test payload"}`)).CreateCtx(nil)
	assert.Nil(err)

	var contents map[string]interface{}
	err = context.PostBodyAsJSON(&contents)
	assert.Nil(err)
	assert.Equal("test payload", contents["test"])

	context, err = NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)
	contents = make(map[string]interface{})
	err = context.PostBodyAsJSON(&contents)
	assert.NotNil(err)
}

type postXMLTest string

func TestCtxPostBodyAsXML(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).WithPostBody([]byte(`<postXMLTest>test payload</postXMLTest>`)).CreateCtx(nil)
	assert.Nil(err)

	var contents postXMLTest
	err = context.PostBodyAsXML(&contents)
	assert.Nil(err)
	assert.Equal("test payload", string(contents))
}

func TestCtxPostBody(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)
	body, err := context.PostBody()
	assert.Nil(err)
	assert.Empty(body)

	context, err = NewMockRequestBuilder(New()).WithPostBody([]byte(`testbytes`)).CreateCtx(nil)
	assert.Nil(err)
	body, err = context.PostBody()
	assert.Equal([]byte(`testbytes`), body)
}

func TestCtxPostedFiles(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)
	postedFiles, err := context.PostedFiles()
	assert.Nil(err)
	assert.Empty(postedFiles)

	context, err = NewMockRequestBuilder(New()).WithPostedFile(PostedFile{
		Key:      "file",
		FileName: "test.txt",
		Contents: []byte("this is only a test")}).CreateCtx(nil)
	assert.Nil(err)

	postedFiles, err = context.PostedFiles()
	assert.Nil(err)
	assert.NotEmpty(postedFiles)
	assert.Equal("file", postedFiles[0].Key)
	assert.Equal("test.txt", postedFiles[0].FileName)
	assert.Equal("this is only a test", string(postedFiles[0].Contents))
}

func TestCtxRouteParam(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).CreateCtx(RouteParameters{"foo": "bar"})
	assert.Nil(err)
	value, err := context.RouteParam("foo")
	assert.Nil(err)
	assert.Equal("bar", value)
}

func TestCtxGetCookie(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).WithCookie(&http.Cookie{Name: "foo", Value: "bar"}).CreateCtx(nil)
	assert.Nil(err)
	assert.Equal("bar", context.GetCookie("foo").Value)
}

func TestCtxHeaderParam(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)
	value, err := context.HeaderValue("test")
	assert.NotNil(err)
	assert.Empty(value)

	context, err = NewMockRequestBuilder(New()).WithHeader("test", "foo").CreateCtx(nil)
	assert.Nil(err)
	value, err = context.HeaderValue("test")
	assert.Nil(err)
	assert.Equal("foo", value)
}

func TestCtxWriteNewCookie(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)

	context.WriteNewCookie("foo", "bar", time.Time{}, "/foo/bar", true)
	assert.Equal("foo=bar; Path=/foo/bar; HttpOnly; Secure", context.Response().Header().Get("Set-Cookie"))
}

func TestCtxExtendCookie(t *testing.T) {
	assert := assert.New(t)

	ctx, err := NewMockRequestBuilder(New()).WithCookieValue("foo", "bar").CreateCtx(nil)
	assert.Nil(err)

	ctx.ExtendCookie("foo", "/", 0, 0, 1)

	cookies := ReadSetCookies(ctx.Response().Header())
	assert.NotEmpty(cookies)
	cookie := cookies[0]
	assert.False(cookie.Expires.IsZero())
}

func TestCtxExtendCookieByDuration(t *testing.T) {
	assert := assert.New(t)

	ctx, err := NewMockRequestBuilder(New()).WithCookieValue("foo", "bar").CreateCtx(nil)
	assert.Nil(err)

	ctx.ExtendCookieByDuration("foo", "/", time.Hour)

	cookies := ReadSetCookies(ctx.Response().Header())
	assert.NotEmpty(cookies)
	cookie := cookies[0]
	assert.False(cookie.Expires.IsZero())
}

func TestCtxRedirect(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)

	result := context.Redirect("foo%sbar")
	assert.Empty(result.Method)
	assert.Equal("foo%sbar", result.RedirectURI)
}

func TestCtxRedirectf(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)

	result := context.Redirectf("foo%sbar", "buzz")
	assert.Empty(result.Method)
	assert.Equal("foobuzzbar", result.RedirectURI)
}

func TestCtxRedirectWithMethod(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)

	result := context.RedirectWithMethod("POST", "foo%sbar")
	assert.Equal("POST", result.Method)
	assert.Equal("foo%sbar", result.RedirectURI)
}

func TestCtxRedirectWithMethodf(t *testing.T) {
	assert := assert.New(t)

	context, err := NewMockRequestBuilder(New()).CreateCtx(nil)
	assert.Nil(err)

	result := context.RedirectWithMethodf("POST", "foo%sbar", "buzz")
	assert.Equal("POST", result.Method)
	assert.Equal("foobuzzbar", result.RedirectURI)
}

func TestCtxPostBodyReparse(t *testing.T) {
	assert := assert.New(t)

	values := url.Values{}
	values.Add("foo", "bar")
	contents := []byte(values.Encode())

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := NewCtx(nil, &http.Request{Method: "POST", Header: headers, Body: ioutil.NopCloser(bytes.NewBuffer(contents))})

	body, err := ctx.PostBody()
	assert.Nil(err)
	assert.Equal(contents, body)

	value, err := ctx.FormValue("foo")
	assert.Nil(err)
	assert.Equal("bar", value)
}
