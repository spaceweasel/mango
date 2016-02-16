package mango

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func NewMockRouter(status int) *mockRouter {
	r := mockRouter{respStatus: status}
	r.respHeader = http.Header{}
	return &r
}

type mockRouter struct {
	req        *http.Request
	respStatus int
	respBody   string
	respHeader http.Header
}

func (r *mockRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.req = req
	for k, v := range r.respHeader {
		for _, hv := range v {
			w.Header().Add(k, hv)
		}
	}
	w.WriteHeader(r.respStatus)

	if req.Body == nil {
		w.Write([]byte(r.respBody))
		return
	}
	// echo body
	rb, _ := ioutil.ReadAll(req.Body)
	w.Write(rb)
}

func TestNewBrowserReturnsPointerToInstanceOfBrowser(t *testing.T) {
	want := "Browser"
	r := Router{}
	b := NewBrowser(&r)
	got := reflect.TypeOf(*b).Name()
	if got != want {
		t.Errorf("Browser type = %q, want %q", got, want)
	}
}

// Get
func TestBrowserGetSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Get("/mangos", nil)
	got := r.req.URL.String()
	if got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
}

func TestBrowserGetSendsRequestWithHeadersToRouter(t *testing.T) {
	want := "application/json"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	b.Get("/mangos", hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserGetReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Get("/mangos", nil)
	got := resp.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestBrowserGetReturnsHeaders(t *testing.T) {
	want := "/the/moon"
	r := NewMockRouter(201)
	r.respHeader.Set("Location", "/the/moon")
	b := NewBrowser(r)
	resp, _ := b.Get("/mangos", nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserGetReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	r.respBody = "40 times around the deck is a mile"
	b := NewBrowser(r)
	resp, _ := b.Get("/mangos", nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

// Post
func TestBrowserPostSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Post("/mangos", "", nil)
	got := r.req.URL.String()
	if got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
}

func TestBrowserPostSendsRequestWithHeadersToRouter(t *testing.T) {
	want := "application/json"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	b.Post("/mangos", "", hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserPostReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Post("/mangos", "", nil)
	got := resp.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestBrowserPostReturnsHeaders(t *testing.T) {
	want := "/the/moon"
	r := NewMockRouter(201)
	r.respHeader.Set("Location", "/the/moon")
	b := NewBrowser(r)
	resp, _ := b.Post("/mangos", "", nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserPostReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	resp, _ := b.Post("/mangos", "40 times around the deck is a mile", nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

// Put
func TestBrowserPutSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Put("/mangos", "", nil)
	got := r.req.URL.String()
	if got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
}

func TestBrowserPutSendsRequestWithHeadersToRouter(t *testing.T) {
	want := "application/json"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	b.Put("/mangos", "", hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserPutReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Put("/mangos", "", nil)
	got := resp.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestBrowserPutReturnsHeaders(t *testing.T) {
	want := "/the/moon"
	r := NewMockRouter(201)
	r.respHeader.Set("Location", "/the/moon")
	b := NewBrowser(r)
	resp, _ := b.Put("/mangos", "", nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserPutReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	resp, _ := b.Put("/mangos", "40 times around the deck is a mile", nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

// Delete
func TestBrowserDelSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Del("/mangos", nil)
	got := r.req.URL.String()
	if got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
}

func TestBrowserDelSendsRequestWithHeadersToRouter(t *testing.T) {
	want := "application/json"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	b.Del("/mangos", hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserDelReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Del("/mangos", nil)
	got := resp.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestBrowserDelReturnsHeaders(t *testing.T) {
	want := "/the/moon"
	r := NewMockRouter(201)
	r.respHeader.Set("Location", "/the/moon")
	b := NewBrowser(r)
	resp, _ := b.Del("/mangos", nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserDelReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	r.respBody = "40 times around the deck is a mile"
	b := NewBrowser(r)
	resp, _ := b.Del("/mangos", nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

// Patch
func TestBrowserPatchSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Patch("/mangos", "", nil)
	got := r.req.URL.String()
	if got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
}

func TestBrowserPatchSendsRequestWithHeadersToRouter(t *testing.T) {
	want := "application/json"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	b.Patch("/mangos", "", hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserPatchReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Patch("/mangos", "", nil)
	got := resp.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestBrowserPatchReturnsHeaders(t *testing.T) {
	want := "/the/moon"
	r := NewMockRouter(201)
	r.respHeader.Set("Location", "/the/moon")
	b := NewBrowser(r)
	resp, _ := b.Patch("/mangos", "", nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserPatchReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	resp, _ := b.Patch("/mangos", "40 times around the deck is a mile", nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}
