package proxy

import (
	"bufio"
	"fmt"
	"golang.org/x/net/http/httpguts"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go-sdk/assert"
	"go-sdk/webutil"
)

func urlMustParse(urlToParse string) *url.URL {
	url, err := url.Parse(urlToParse)
	if err != nil {
		panic(err)
	}
	return url
}

func TestProxy(t *testing.T) {
	assert := assert.New(t)

	mockedEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if protoHeader := r.Header.Get(webutil.HeaderXForwardedProto); protoHeader == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No `X-Forwarded-Proto` header!"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok!"))
		return
	}))
	defer mockedEndpoint.Close()

	target, err := url.Parse(mockedEndpoint.URL)
	assert.Nil(err)

	proxy := New().WithUpstream(NewUpstream(target))
	proxy.WithUpstreamHeader(webutil.HeaderXForwardedProto, webutil.SchemeHTTP)

	mockedProxy := httptest.NewServer(proxy)

	res, err := http.Get(mockedProxy.URL)
	assert.Nil(err)
	defer res.Body.Close()

	fullBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)

	mockedContents := string(fullBody)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Equal("Ok!", mockedContents)
}

func TestReverseProxyWebSocket(t *testing.T) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if upgradeType(r.Header) != "websocket" {
			t.Error("unexpected backend request")
			http.Error(w, "unexpected request", 400)
			return
		}
		c, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error(err)
			return
		}
		defer c.Close()
		io.WriteString(c, "HTTP/1.1 101 Switching Protocols\r\nConnection: upgrade\r\nUpgrade: WebSocket\r\n\r\n")
		bs := bufio.NewScanner(c)
		if !bs.Scan() {
			t.Errorf("backend failed to read line from client: %v", bs.Err())
			return
		}
		fmt.Fprintf(c, "backend got %q\n", bs.Text())
	}))
	defer backendServer.Close()

	backURL, _ := url.Parse(backendServer.URL)
	proxy := New().WithUpstream(NewUpstream(backURL))
	proxy.WithUpstreamHeader(webutil.HeaderXForwardedProto, webutil.SchemeHTTP)

	frontendProxy := httptest.NewServer(proxy)
	defer frontendProxy.Close()

	req, _ := http.NewRequest("GET", frontendProxy.URL, nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")

	c := frontendProxy.Client()
	res, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 101 {
		t.Fatalf("status = %v; want 101", res.Status)
	}
	if upgradeType(res.Header) != "websocket" {
		t.Fatalf("not websocket upgrade; got %#v", res.Header)
	}
	rwc, ok := res.Body.(io.ReadWriteCloser)
	if !ok {
		t.Fatalf("response body is of type %T; does not implement ReadWriteCloser", res.Body)
	}
	defer rwc.Close()

	io.WriteString(rwc, "Hello\n")
	bs := bufio.NewScanner(rwc)
	if !bs.Scan() {
		t.Fatalf("Scan: %v", bs.Err())
	}
	got := bs.Text()
	want := `backend got "Hello"`
	if got != want {
		t.Errorf("got %#q, want %#q", got, want)
	}
}

func upgradeType(h http.Header) string {
	if !httpguts.HeaderValuesContainsToken(h["Connection"], "Upgrade") {
		return ""
	}
	return strings.ToLower(h.Get("Upgrade"))
}