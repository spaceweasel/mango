package mango

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

// Browser is used for simulating HTTP request to the service.
// This can be useful for verifying routings, testing request headers,
// as well as examining responses: headers, status code, content.
type Browser struct {
	router *Router
}

// NewBrowser returns a *Browser which can be used to test server responses.
// This method takes a Router as a parameter to enable full and accurate
// testing/simulation to be performed.
func NewBrowser(r *Router) *Browser {
	return &Browser{router: r}
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
	b.router.ServeHTTP(w, req)
	return w, nil
}

// get info from result
// resp, err := br.Post(...)
// resp.HeaderMap.Get("Location") // important to use HeaderMap and not Header()!

// Post simulates an HTTP POST request to the server.
// The response can be examined afterwards to check status, headers
// and content.
func (b *Browser) Post(url, body string, headers http.Header) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header = headers
	w := httptest.NewRecorder()
	b.router.ServeHTTP(w, req)
	return w, nil
}
