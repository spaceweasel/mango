package mango

import (
	"bytes"
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

func TestBrowserGetReturnsErrorWhenInvalidURL(t *testing.T) {
	want := "parse ::/mangos: missing protocol scheme"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	_, err := b.Get("::/mangos", nil)
	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

// Post
func TestBrowserPostSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Post("/mangos", nil, nil)
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
	b.Post("/mangos", nil, hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserPostReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Post("/mangos", nil, nil)
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
	resp, _ := b.Post("/mangos", nil, nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserPostReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	bytes := bytes.NewBufferString("40 times around the deck is a mile")
	resp, _ := b.Post("/mangos", bytes, nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestBrowserPostReturnsErrorWhenInvalidURL(t *testing.T) {
	want := "parse ::/mangos: missing protocol scheme"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	bytes := bytes.NewBufferString("mangos")
	_, err := b.Post("::/mangos", bytes, nil)
	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

// PostS
func TestBrowserPostSReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	resp, _ := b.PostS("/mangos", "40 times around the deck is a mile", nil)
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
	b.Put("/mangos", nil, nil)
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
	b.Put("/mangos", nil, hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserPutReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Put("/mangos", nil, nil)
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
	resp, _ := b.Put("/mangos", nil, nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserPutReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	bytes := bytes.NewBufferString("40 times around the deck is a mile")
	resp, _ := b.Put("/mangos", bytes, nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestBrowserPutReturnsErrorWhenInvalidURL(t *testing.T) {
	want := "parse ::/mangos: missing protocol scheme"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	bytes := bytes.NewBufferString("mangos")
	_, err := b.Put("::/mangos", bytes, nil)
	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

// PutS
func TestBrowserPutSReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	resp, _ := b.PutS("/mangos", "40 times around the deck is a mile", nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

// Delete
func TestBrowserDeleteSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Delete("/mangos", nil)
	got := r.req.URL.String()
	if got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
}

func TestBrowserDeleteSendsRequestWithHeadersToRouter(t *testing.T) {
	want := "application/json"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	b.Delete("/mangos", hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserDeleteReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Delete("/mangos", nil)
	got := resp.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestBrowserDeleteReturnsHeaders(t *testing.T) {
	want := "/the/moon"
	r := NewMockRouter(201)
	r.respHeader.Set("Location", "/the/moon")
	b := NewBrowser(r)
	resp, _ := b.Delete("/mangos", nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserDeleteReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	r.respBody = "40 times around the deck is a mile"
	b := NewBrowser(r)
	resp, _ := b.Delete("/mangos", nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestBrowserDeleteReturnsErrorWhenInvalidURL(t *testing.T) {
	want := "parse ::/mangos: missing protocol scheme"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	_, err := b.Delete("::/mangos", nil)
	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

// Patch
func TestBrowserPatchSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Patch("/mangos", nil, nil)
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
	b.Patch("/mangos", nil, hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserPatchReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Patch("/mangos", nil, nil)
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
	resp, _ := b.Patch("/mangos", nil, nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserPatchReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	bytes := bytes.NewBufferString("40 times around the deck is a mile")
	resp, _ := b.Patch("/mangos", bytes, nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestBrowserPatchReturnsErrorWhenInvalidURL(t *testing.T) {
	want := "parse ::/mangos: missing protocol scheme"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	bytes := bytes.NewBufferString("mangos")
	_, err := b.Patch("::/mangos", bytes, nil)
	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

// PatchS
func TestBrowserPatchSReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	resp, _ := b.PatchS("/mangos", "40 times around the deck is a mile", nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

// Head
func TestBrowserHeadSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Head("/mangos", nil)
	got := r.req.URL.String()
	if got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
}

func TestBrowserHeadSendsRequestWithHeadersToRouter(t *testing.T) {
	want := "application/json"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	b.Head("/mangos", hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserHeadReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Head("/mangos", nil)
	got := resp.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestBrowserHeadReturnsHeaders(t *testing.T) {
	want := "/the/moon"
	r := NewMockRouter(201)
	r.respHeader.Set("Location", "/the/moon")
	b := NewBrowser(r)
	resp, _ := b.Head("/mangos", nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserHeadReturnsErrorWhenInvalidURL(t *testing.T) {
	want := "parse ::/mangos: missing protocol scheme"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	_, err := b.Head("::/mangos", nil)
	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

// Options
func TestBrowserOptionsSendsRequestWithCorrectURLToRouter(t *testing.T) {
	want := "/mangos"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	b.Options("/mangos", nil)
	got := r.req.URL.String()
	if got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
}

func TestBrowserOptionsSendsRequestWithHeadersToRouter(t *testing.T) {
	want := "application/json"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	b.Options("/mangos", hdrs)
	got := r.req.Header.Get("Accept")
	if got != want {
		t.Errorf("Accept header = %q, want %q", got, want)
	}
}

func TestBrowserOptionsReturnsStatusCode(t *testing.T) {
	want := 404
	r := NewMockRouter(404)
	b := NewBrowser(r)
	resp, _ := b.Options("/mangos", nil)
	got := resp.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestBrowserOptionsReturnsHeaders(t *testing.T) {
	want := "/the/moon"
	r := NewMockRouter(201)
	r.respHeader.Set("Location", "/the/moon")
	b := NewBrowser(r)
	resp, _ := b.Options("/mangos", nil)
	got := resp.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location header = %q, want %q", got, want)
	}
}

func TestBrowserOptionsReturnsBody(t *testing.T) {
	want := "40 times around the deck is a mile"
	r := NewMockRouter(200)
	r.respBody = "40 times around the deck is a mile"
	b := NewBrowser(r)
	resp, _ := b.Options("/mangos", nil)
	got := resp.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestBrowserOptionsReturnsErrorWhenInvalidURL(t *testing.T) {
	want := "parse ::/mangos: missing protocol scheme"
	r := NewMockRouter(200)
	b := NewBrowser(r)
	_, err := b.Options("::/mangos", nil)
	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}
