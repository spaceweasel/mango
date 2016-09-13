package mango

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestOriginAllowedReturnsFalseWhenNoOriginsInConfigOrigins(t *testing.T) {
	want := false
	cc := CORSConfig{}
	got := cc.originAllowed("http://elsewhere.com")
	if got != want {
		t.Errorf("Origin Allowed = %t, want %t", got, want)
	}
}

func TestOriginAllowedReturnsFalseWhenOriginNotInConfigOrigins(t *testing.T) {
	want := false
	cc := CORSConfig{}
	cc.Origins = []string{"http://somewhere.com", "http://elsewhere.com"}
	got := cc.originAllowed("http://mangowhere.com")
	if got != want {
		t.Errorf("Origin Allowed = %t, want %t", got, want)
	}
}

func TestOriginAllowedReturnsTrueWhenOriginInConfigOrigins(t *testing.T) {
	want := true
	cc := CORSConfig{}
	cc.Origins = []string{"http://somewhere.com", "http://elsewhere.com"}
	got := cc.originAllowed("http://elsewhere.com")
	if got != want {
		t.Errorf("Origin Allowed = %t, want %t", got, want)
	}
}

func TestOriginAllowedReturnsTrueWhenConfigOriginsContainsStar(t *testing.T) {
	want := true
	cc := CORSConfig{}
	cc.Origins = []string{"*"}
	got := cc.originAllowed("http://elsewhere.com")
	if got != want {
		t.Errorf("Origin Allowed = %t, want %t", got, want)
	}
}

func TestMethodAllowedReturnsFalseWhenNoMethodsInConfigMethods(t *testing.T) {
	want := false
	cc := CORSConfig{}
	got := cc.methodAllowed("POST")
	if got != want {
		t.Errorf("Method Allowed = %t, want %t", got, want)
	}
}

func TestMethodAllowedReturnsFalseWhenMethodNotInConfigMethods(t *testing.T) {
	want := false
	cc := CORSConfig{}
	cc.Methods = []string{"PATCH", "PUT"}
	got := cc.methodAllowed("POST")
	if got != want {
		t.Errorf("Method Allowed = %t, want %t", got, want)
	}
}

func TestMethodAllowedReturnsTrueWhenMethodInConfigMethods(t *testing.T) {
	want := true
	cc := CORSConfig{}
	cc.Methods = []string{"POST", "PUT"}
	got := cc.methodAllowed("POST")
	if got != want {
		t.Errorf("Method Allowed = %t, want %t", got, want)
	}
}

func TestHeaderAllowedReturnsFalseWhenNoHeadersInConfigHeaders(t *testing.T) {
	want := false
	cc := CORSConfig{}
	got := cc.headersAllowed("http://elsewhere.com")
	if got != want {
		t.Errorf("Header Allowed = %t, want %t", got, want)
	}
}

func TestHeaderAllowedReturnsFalseWhenHeaderNotInConfigHeaders(t *testing.T) {
	want := false
	cc := CORSConfig{}
	cc.Headers = []string{"X-Fruit", "X-Special"}
	got := cc.headersAllowed("X-Moondust")
	if got != want {
		t.Errorf("Header Allowed = %t, want %t", got, want)
	}
}

func TestHeaderAllowedReturnsTrueWhenHeaderInConfigHeaders(t *testing.T) {
	want := true
	cc := CORSConfig{}
	cc.Headers = []string{"X-Fruit", "X-Special"}
	got := cc.headersAllowed("X-Fruit")
	if got != want {
		t.Errorf("Header Allowed = %t, want %t", got, want)
	}
}

func TestHeaderAllowedReturnsTrueWhenAllHeadersInConfigHeaders(t *testing.T) {
	want := true
	cc := CORSConfig{}
	cc.Headers = []string{"X-Fruit", "X-Special", "X-Mango", "X-Biscuit"}
	got := cc.headersAllowed("X-Fruit, X-Biscuit, X-Mango")
	if got != want {
		t.Errorf("Headers Allowed = %t, want %t", got, want)
	}
}

func TestHeaderAllowedReturnsFalseWhenNotAllHeadersInConfigHeaders(t *testing.T) {
	want := false
	cc := CORSConfig{}
	cc.Headers = []string{"X-Fruit", "X-Special", "X-Biscuit"}
	got := cc.headersAllowed("X-Fruit, X-Biscuit, X-Mango")
	if got != want {
		t.Errorf("Headers Allowed = %t, want %t", got, want)
	}
}

func TestHeaderAllowedReturnsTrueWhenEmptyHeaderAndHeadersInConfigHeaders(t *testing.T) {
	want := true
	cc := CORSConfig{}
	cc.Headers = []string{"X-Fruit", "X-Special", "X-Biscuit"}
	got := cc.headersAllowed("")
	if got != want {
		t.Errorf("Headers Allowed = %t, want %t", got, want)
	}
}

func TestHeaderAllowedReturnsTrueWhenEmptyHeadersAndHeadersInConfigHeaders(t *testing.T) {
	want := true
	cc := CORSConfig{}
	cc.Headers = []string{"X-Fruit", "X-Special", "X-Biscuit"}
	got := cc.headersAllowed(" , ")
	if got != want {
		t.Errorf("Headers Allowed = %t, want %t", got, want)
	}
}

func TestHeaderAllowedReturnsTrueWhenEmptyHeaderAndNoHeadersInConfigHeaders(t *testing.T) {
	want := true
	cc := CORSConfig{}
	cc.Headers = []string{}
	got := cc.headersAllowed("")
	if got != want {
		t.Errorf("Headers Allowed = %t, want %t", got, want)
	}
}

func TestHandleCORSMakesNoChangesToResponseWhenResourceHasNoCORSConfig(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	w := httptest.NewRecorder()
	res := Resource{}
	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Origin")

	if got != want {
		t.Errorf("Origin = %q, want %q", got, want)
	}
}

func TestHandleCORSDoesNotAddAllowOriginHeaderWhenRequestHasNoOriginHeader(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	w := httptest.NewRecorder()
	res := Resource{}
	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Origin")

	if got != want {
		t.Errorf("Origin = %q, want %q", got, want)
	}
}

func TestHandleCORSAddsAllowOriginHeaderWhenRequestOriginAllowed(t *testing.T) {
	want := "http://greencheese.com"
	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
		},
	}

	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Origin")

	if got != want {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, want)
	}
}

func TestHandleCORSAddsVaryOriginHeaderWhenRequestOriginAllowed(t *testing.T) {
	want := "Origin"
	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
		},
	}
	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Vary")

	if got != want {
		t.Errorf("Vary = %q, want %q", got, want)
	}
}

func TestHandleCORSDoesNotAddAllowOriginHeaderWhenRequestOriginNotAllowed(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://bluecheese.com")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
		},
	}

	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Origin")

	if got != want {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, want)
	}
}

func TestHandleCORSDoesNotSetAllowCredentialsHeaderWhenResourceDoesNotAllowCredentials(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
		},
	}
	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Credentials")

	if got != want {
		t.Errorf("Access-Control-Allow-Credentials = %q, want %q", got, want)
	}
}

func TestHandleCORSSetsAllowCredentialsHeaderWhenResourceAllowsCredentials(t *testing.T) {
	want := "true"
	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins:     []string{"http://greencheese.com"},
			Credentials: true,
		},
	}
	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Credentials")

	if got != want {
		t.Errorf("Access-Control-Allow-Credentials = %q, want %q", got, want)
	}
}

func TestHandleCORSIncludesExposedHeadersWhenContainedInResourceConfig(t *testing.T) {
	want := "X-Cheese, X-Mangoes"
	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins:        []string{"http://greencheese.com"},
			ExposedHeaders: []string{"X-Cheese", "X-Mangoes"},
		},
	}
	handleCORS(req, w, &res)

	got := strings.Join(w.HeaderMap["Access-Control-Expose-Headers"], ", ")

	if got != want {
		t.Errorf("Access-Control-Expose-Headers = %q, want %q", got, want)
	}
}

func TestHandleCORSDoesNotAddAllowOriginHeaderWhenPreflightRequestHasNoACRequestMethod(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
		},
	}
	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Origin")

	if got != want {
		t.Errorf("Origin = %q, want %q", got, want)
	}
}

func TestHandleCORSDoesNotAddAllowOriginHeaderWhenPreflightRequestHasACRequestMethodNotInCORSConfig(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"GET"},
		},
	}
	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Origin")

	if got != want {
		t.Errorf("Origin = %q, want %q", got, want)
	}
}

func TestHandleCORSDoesNotAddAllowOriginHeaderWhenPreflightRequestHasACRequestHeadersNotInCORSConfig(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese", "X-Biscuits"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST": nil,
		},
	}
	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Origin")

	if got != want {
		t.Errorf("Origin = %q, want %q", got, want)
	}
}

func TestHandleCORSSetsAllowOriginWhenPreflightSucceeds(t *testing.T) {
	want := "http://greencheese.com"
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese", "X-Mangoes"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST": nil,
		},
	}
	handleCORS(req, w, &res)

	got := w.HeaderMap.Get("Access-Control-Allow-Origin")

	if got != want {
		t.Errorf("Origin = %q, want %q", got, want)
	}
}

func TestHandleCORSSetsAllowMethodsWhenPreflightSucceeds(t *testing.T) {
	want := "POST, PATCH"
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese", "X-Mangoes"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST":  nil,
			"PATCH": nil,
			"PUT":   nil,
		},
	}
	handleCORS(req, w, &res)
	got := strings.Join(w.HeaderMap["Access-Control-Allow-Methods"], ", ")

	if got != want {
		t.Errorf("Methods = %q, want %q", got, want)
	}
}

func TestHandleCORSSetsAllowHeadersWhenPreflightSucceeds(t *testing.T) {
	want := "X-Cheese, X-Mangoes"
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese", "X-Mangoes"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST": nil,
		},
	}
	handleCORS(req, w, &res)
	got := strings.Join(w.HeaderMap["Access-Control-Allow-Headers"], ", ")

	if got != want {
		t.Errorf("Headers = %q, want %q", got, want)
	}
}

func TestHandleCORSDoesNotSetsExposeHeadersWhenPreflightSucceeds(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins:        []string{"http://greencheese.com"},
			Methods:        []string{"POST", "PATCH"},
			Headers:        []string{"X-Cheese", "X-Mangoes"},
			ExposedHeaders: []string{"X-Cheese", "X-Mangoes"},
		},
	}
	handleCORS(req, w, &res)
	got := w.HeaderMap.Get("Access-Control-Expose-Headers")

	if got != want {
		t.Errorf("Expose Headers = %q, want %q", got, want)
	}
}

func TestHandleCORSSetsMaxAgeWhenPreflightSucceedsAndMaxAgeGreaterthanZero(t *testing.T) {
	want := "172800"
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese", "X-Mangoes"},
			MaxAge:  172800,
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST": nil,
		},
	}
	handleCORS(req, w, &res)
	got := w.HeaderMap.Get("Access-Control-Max-Age")

	if got != want {
		t.Errorf("MaxAge = %q, want %q", got, want)
	}
}

func TestHandleCORSDoesNotSetMaxAgeWhenPreflightSucceedsWhenMaxAgeEqualToZero(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins:        []string{"http://greencheese.com"},
			Methods:        []string{"POST", "PATCH"},
			Headers:        []string{"X-Cheese", "X-Mangoes"},
			ExposedHeaders: []string{"X-Cheese", "X-Mangoes"},
		},
	}
	handleCORS(req, w, &res)
	got := w.HeaderMap.Get("Access-Control-Max-Age")

	if got != want {
		t.Errorf("MaxAge = %q, want %q", got, want)
	}
}

func TestHandleCORSOnlyIncludesResourceMethodsInAllowMethodsWhenPreflightSucceeds(t *testing.T) {
	want := "PATCH, PUT"
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "PUT")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH", "PUT"},
			Headers: []string{"X-Cheese", "X-Mangoes"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"PATCH": nil,
			"PUT":   nil,
		},
	}
	handleCORS(req, w, &res)
	got := strings.Join(w.HeaderMap["Access-Control-Allow-Methods"], ", ")

	if got != want {
		t.Errorf("Methods = %q, want %q", got, want)
	}
}

func TestHandleCORSFailsPreflightIfRequestMethodNotInResourceMethodsEvenWhenInCORSConfig(t *testing.T) {
	want := ""
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH", "PUT"},
			Headers: []string{"X-Cheese", "X-Mangoes"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"PATCH": nil,
			"PUT":   nil,
		},
	}
	handleCORS(req, w, &res)
	got := w.HeaderMap.Get("Access-Control-Allow-Origin")

	if got != want {
		t.Errorf("Origin = %q, want %q", got, want)
	}
}

func TestHandleCORSReturnsFalseWhenNoOrigin(t *testing.T) {
	want := false
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST": nil,
		},
	}
	got := handleCORS(req, w, &res)

	if got != want {
		t.Errorf("Preflight = %t, want %t", got, want)
	}
}

func TestHandleCORSReturnsFalseWhenNoAccessControlRequestMethod(t *testing.T) {
	want := false
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST": nil,
		},
	}
	got := handleCORS(req, w, &res)

	if got != want {
		t.Errorf("Preflight = %t, want %t", got, want)
	}
}

func TestHandleCORSReturnsFalseWhenNotOPTIONSRequest(t *testing.T) {
	want := false
	req, _ := http.NewRequest("POST", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST": nil,
		},
	}
	got := handleCORS(req, w, &res)

	if got != want {
		t.Errorf("Preflight = %t, want %t", got, want)
	}
}

func TestHandleCORSReturnsTrueWhenPreflightFails(t *testing.T) {
	want := true
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST": nil,
		},
	}
	got := handleCORS(req, w, &res)

	if got != want {
		t.Errorf("Preflight = %t, want %t", got, want)
	}
}

func TestHandleCORSReturnsTrueWhenPreflightSucceeds(t *testing.T) {
	want := true
	req, _ := http.NewRequest("OPTIONS", "https://somewhere.com/mango", nil)
	req.Header.Set("Origin", "http://greencheese.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Mangoes")
	w := httptest.NewRecorder()
	res := Resource{
		CORSConfig: &CORSConfig{
			Origins: []string{"http://greencheese.com"},
			Methods: []string{"POST", "PATCH"},
			Headers: []string{"X-Cheese", "X-Mangoes"},
		},
		Handlers: map[string]ContextHandlerFunc{
			"POST": nil,
		},
	}
	got := handleCORS(req, w, &res)

	if got != want {
		t.Errorf("Preflight = %t, want %t", got, want)
	}
}

func TestCORSConfigClone(t *testing.T) {
	config := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}

	tests := []struct {
		want string
		name string
		fn   func(c *CORSConfig) string
	}{
		{
			"http://bluecheese.com",
			"Origins",
			func(c *CORSConfig) string {
				return strings.Join(c.Origins, ", ")
			},
		},
		{
			"POST, PATCH",
			"Methods",
			func(c *CORSConfig) string {
				return strings.Join(c.Methods, ", ")
			},
		},
		{
			"X-Cheese, X-Mangoes",
			"Headers",
			func(c *CORSConfig) string {
				return strings.Join(c.Headers, ", ")
			},
		},
		{
			"X-Biscuits, X-Mangoes",
			"ExposedHeaders",
			func(c *CORSConfig) string {
				return strings.Join(c.ExposedHeaders, ", ")
			},
		},
		{
			"true",
			"Credentials",
			func(c *CORSConfig) string {
				return strconv.FormatBool(c.Credentials)
			},
		},
		{
			"45",
			"MaxAge",
			func(c *CORSConfig) string {
				return strconv.Itoa(c.MaxAge)
			},
		},
	}

	clone := config.clone()

	for _, test := range tests {
		if got := test.fn(clone); got != test.want {
			t.Errorf("CORSConfig.%s = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestCORSConfigMerge(t *testing.T) {

	destConfig := &CORSConfig{
		Origins:        []string{"http://greencheese.com"},
		Methods:        []string{"PUT"},
		Headers:        []string{"X-Custard", "X-Fish"},
		ExposedHeaders: []string{"X-Onions"},
	}

	srcConfig := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}

	tests := []struct {
		want string
		name string
		fn   func(c *CORSConfig) string
	}{
		{
			"http://greencheese.com, http://bluecheese.com",
			"Origins",
			func(c *CORSConfig) string {
				return strings.Join(c.Origins, ", ")
			},
		},
		{
			"PUT, POST, PATCH",
			"Methods",
			func(c *CORSConfig) string {
				return strings.Join(c.Methods, ", ")
			},
		},
		{
			"X-Custard, X-Fish, X-Cheese, X-Mangoes",
			"Headers",
			func(c *CORSConfig) string {
				return strings.Join(c.Headers, ", ")
			},
		},
		{
			"X-Onions, X-Biscuits, X-Mangoes",
			"ExposedHeaders",
			func(c *CORSConfig) string {
				return strings.Join(c.ExposedHeaders, ", ")
			},
		},
		{
			"true",
			"Credentials",
			func(c *CORSConfig) string {
				return strconv.FormatBool(c.Credentials)
			},
		},
		{
			"45",
			"MaxAge",
			func(c *CORSConfig) string {
				return strconv.Itoa(c.MaxAge)
			},
		},
	}

	destConfig.merge(srcConfig)

	for _, test := range tests {
		if got := test.fn(destConfig); got != test.want {
			t.Errorf("CORSConfig.%s = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestCORSConfigMergeRemovesDuplicates(t *testing.T) {
	destConfig := &CORSConfig{
		Origins:        []string{"http://greencheese.com", "http://bluecheese.com"},
		Methods:        []string{"PUT", "PATCH"},
		Headers:        []string{"X-Custard", "X-Fish", "X-Mangoes"},
		ExposedHeaders: []string{"X-Onions", "X-Biscuits"},
	}

	srcConfig := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}

	tests := []struct {
		want string
		name string
		fn   func(c *CORSConfig) string
	}{
		{
			"http://greencheese.com, http://bluecheese.com",
			"Origins",
			func(c *CORSConfig) string {
				return strings.Join(c.Origins, ", ")
			},
		},
		{
			"PUT, PATCH, POST",
			"Methods",
			func(c *CORSConfig) string {
				return strings.Join(c.Methods, ", ")
			},
		},
		{
			"X-Custard, X-Fish, X-Mangoes, X-Cheese",
			"Headers",
			func(c *CORSConfig) string {
				return strings.Join(c.Headers, ", ")
			},
		},
		{
			"X-Onions, X-Biscuits, X-Mangoes",
			"ExposedHeaders",
			func(c *CORSConfig) string {
				return strings.Join(c.ExposedHeaders, ", ")
			},
		},
		{
			"true",
			"Credentials",
			func(c *CORSConfig) string {
				return strconv.FormatBool(c.Credentials)
			},
		},
		{
			"45",
			"MaxAge",
			func(c *CORSConfig) string {
				return strconv.Itoa(c.MaxAge)
			},
		},
	}

	destConfig.merge(srcConfig)

	for _, test := range tests {
		if got := test.fn(destConfig); got != test.want {
			t.Errorf("CORSConfig.%s = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestCORSConfigMergeSrcCredentialsAndMaxAgeOverrideDest(t *testing.T) {
	destConfig := &CORSConfig{
		Credentials: false,
		MaxAge:      32,
	}

	srcConfig := CORSConfig{
		Credentials: true,
		MaxAge:      45,
	}

	tests := []struct {
		want string
		name string
		fn   func(c *CORSConfig) string
	}{
		{
			"true",
			"Credentials",
			func(c *CORSConfig) string {
				return strconv.FormatBool(c.Credentials)
			},
		},
		{
			"45",
			"MaxAge",
			func(c *CORSConfig) string {
				return strconv.Itoa(c.MaxAge)
			},
		},
	}

	destConfig.merge(srcConfig)

	for _, test := range tests {
		if got := test.fn(destConfig); got != test.want {
			t.Errorf("CORSConfig.%s = %q, want %q", test.name, got, test.want)
		}
	}
}

// Examples

func ExampleCORSConfig() {
	// Simple Request using Test Browser
	r := NewRouter()
	r.SetGlobalCORS(CORSConfig{
		Origins:        []string{"*"},
		Methods:        []string{"POST", "PUT"},
		Headers:        []string{"X-Mangoes"},
		ExposedHeaders: []string{"X-Mangoes"},
	})

	r.Get("/fruits", func(c *Context) {
		c.RespondWith("GET fruits")
	})

	br := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Origin", "http://bluecheese.com")
	resp, err := br.Get("http://greencheese.com/fruits", hdrs)

	if err != nil {
		fmt.Println(err)
		return
	}
	allowOrigin := resp.HeaderMap.Get("Access-Control-Allow-Origin")
	exposedHeaders := resp.HeaderMap.Get("Access-Control-Expose-Headers")
	fmt.Println(allowOrigin)
	fmt.Println(exposedHeaders)
	// Output:
	// http://bluecheese.com
	// X-Mangoes
}

func ExampleCORSConfig_second() {
	// Preflight Request using Test Browser
	r := NewRouter()
	r.SetGlobalCORS(CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PUT"},
		Headers:        []string{"X-Mangoes"},
		ExposedHeaders: []string{"X-Mangoes"},
	})
	r.Post("/fruits", func(c *Context) {
		c.RespondWith("POST fruits")
	})
	r.Get("/fruits", func(c *Context) {
		c.RespondWith("GET fruits")
	})

	br := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Origin", "http://bluecheese.com")
	hdrs.Set("Access-Control-Request-Method", "POST")
	hdrs.Set("Access-Control-Request-Headers", "X-Mangoes")
	resp, err := br.Options("http://greencheese.com/fruits", hdrs)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Examine the response headers...
	allowOrigin := resp.HeaderMap.Get("Access-Control-Allow-Origin")
	allowMethods := resp.HeaderMap.Get("Access-Control-Allow-Methods")
	allowHeaders := resp.HeaderMap.Get("Access-Control-Allow-Headers")
	vary := resp.HeaderMap.Get("Vary")

	fmt.Println(allowOrigin)  // http://bluecheese.com
	fmt.Println(allowMethods) // POST (PUT has no handler, so is removed from list)
	fmt.Println(allowHeaders) // X-Mangoes
	fmt.Println(vary)         // Origin
	fmt.Println(resp.Code)    // 200
	fmt.Println(resp.Body)    // Body should be empty

	// Output:
	// http://bluecheese.com
	// POST
	// X-Mangoes
	// Origin
	// 200
	//
}

func ExampleCORSConfig_third() {
	// Preflight Request using Test Browser
	// Still returns status 200 even though preflight check fails
	r := NewRouter()
	r.SetGlobalCORS(CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PUT"},
		Headers:        []string{"X-Mangoes"},
		ExposedHeaders: []string{"X-Mangoes"},
	})
	r.Post("/fruits", func(c *Context) {
		c.RespondWith("POST fruits")
	})
	r.Get("/fruits", func(c *Context) {
		c.RespondWith("GET fruits")
	})

	br := NewBrowser(r)
	hdrs := http.Header{}
	hdrs.Set("Origin", "http://bluecheese.com")
	hdrs.Set("Access-Control-Request-Method", "PATCH")
	hdrs.Set("Access-Control-Request-Headers", "X-Mangoes")
	resp, err := br.Options("http://greencheese.com/fruits", hdrs)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Examine the response headers...
	allowOrigin := resp.HeaderMap.Get("Access-Control-Allow-Origin")
	allowMethods := resp.HeaderMap.Get("Access-Control-Allow-Methods")
	allowHeaders := resp.HeaderMap.Get("Access-Control-Allow-Headers")
	vary := resp.HeaderMap.Get("Vary")

	// preflight check fails so all but the Resp.Code will be empty

	fmt.Println(allowOrigin)
	fmt.Println(allowMethods)
	fmt.Println(allowHeaders)
	fmt.Println(vary)
	fmt.Println(resp.Code) // 200
	fmt.Println(resp.Body)

	// Output:
	//
	//
	//
	//
	// 200
	//
}
