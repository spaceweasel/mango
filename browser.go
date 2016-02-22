package mango

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
)

// Browser is client used to simulate HTTP request to the service.
// This can be useful for verifying routings, testing request headers,
// as well as examining responses: headers, status code, content.
type Browser struct {
	handler http.Handler
}

// NewBrowser returns a *Browser which can be used to test server responses.
// This method takes a Router as a parameter to enable full and accurate
// testing/simulation to be performed.
func NewBrowser(h http.Handler) *Browser {
	return &Browser{handler: h}
}

// Get simulates an HTTP GET request to the server.
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) Get(url string, headers http.Header) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	w := httptest.NewRecorder()
	b.handler.ServeHTTP(w, req)
	return w, nil
}

// Post simulates an HTTP POST request to the server.
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) Post(url string, body io.Reader, headers http.Header) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	w := httptest.NewRecorder()
	b.handler.ServeHTTP(w, req)
	return w, nil
}

// PostS simulates an HTTP POST request to the server, with a string body.
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) PostS(url, body string, headers http.Header) (*httptest.ResponseRecorder, error) {
	bytes := bytes.NewBufferString(body)
	return b.Post(url, bytes, headers)
}

// Put simulates an HTTP PUT request to the server.
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) Put(url string, body io.Reader, headers http.Header) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	w := httptest.NewRecorder()
	b.handler.ServeHTTP(w, req)
	return w, nil
}

// PutS simulates an HTTP PUT request to the server, with a string body.
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) PutS(url, body string, headers http.Header) (*httptest.ResponseRecorder, error) {
	bytes := bytes.NewBufferString(body)
	return b.Put(url, bytes, headers)
}

// Del simulates an HTTP DELETE request to the server.
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) Del(url string, headers http.Header) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	w := httptest.NewRecorder()
	b.handler.ServeHTTP(w, req)
	return w, nil
}

// Patch simulates an HTTP PATCH request to the server, with a string body..
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) Patch(url string, body io.Reader, headers http.Header) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	w := httptest.NewRecorder()
	b.handler.ServeHTTP(w, req)
	return w, nil
}

// PatchS simulates an HTTP PATCH request to the server, with a string body..
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) PatchS(url, body string, headers http.Header) (*httptest.ResponseRecorder, error) {
	bytes := bytes.NewBufferString(body)
	return b.Patch(url, bytes, headers)
}

// Head simulates an HTTP HEAD request to the server.
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) Head(url string, headers http.Header) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	w := httptest.NewRecorder()
	b.handler.ServeHTTP(w, req)
	return w, nil
}

// Options simulates an HTTP OPTIONS request to the server.
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) Options(url string, headers http.Header) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	w := httptest.NewRecorder()
	b.handler.ServeHTTP(w, req)
	return w, nil
}
