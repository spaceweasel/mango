package mango

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type mockRoutes struct {
	routes           map[string]map[string]ContextHandlerFunc
	validators       map[string]Validator
	corsConfigs      map[string]CORSConfig
	globalCORSConfig CORSConfig
}

func (m *mockRoutes) TestValidators(s, constraint string) bool {
	v, ok := m.validators[constraint]
	if !ok {
		panic("Unknown constraint" + constraint)
	}
	return v.Validate(s, []string{})
}

func (m *mockRoutes) AddHandlerFunc(pattern, method string, handlerFunc ContextHandlerFunc) {
	_, ok := m.routes[pattern]
	if !ok {
		m.routes[pattern] = make(map[string]ContextHandlerFunc)
	}
	_, dup := m.routes[pattern][method]
	if dup {
		panic(fmt.Sprintf("duplicate route handler method: \"%s %s\"", method, pattern))
	}
	m.routes[pattern][method] = handlerFunc
}

func (m *mockRoutes) GetResource(path string) (*Resource, bool) {
	hm, ok := m.routes[path]
	res := Resource{
		Handlers: hm,
	}
	return &res, ok
}

func (m *mockRoutes) SetGlobalCORS(config CORSConfig) {
	m.globalCORSConfig = config
}

func (m *mockRoutes) SetCORS(pattern string, config CORSConfig) {
	m.corsConfigs[pattern] = config
}

func (m *mockRoutes) AddCORS(pattern string, config CORSConfig) {
	m.corsConfigs[pattern] = config
}

func newMockRoutes() *mockRoutes {
	mr := mockRoutes{}
	mr.routes = make(map[string]map[string]ContextHandlerFunc)
	mr.validators = make(map[string]Validator)
	mr.corsConfigs = make(map[string]CORSConfig)
	return &mr
}

func TestGetAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Get("/test", testFunc)
	resource, ok := rtr.routes.GetResource("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := resource.Handlers["GET"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function got %q, want %q", got, want)
	}
}

func TestPostAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Post("/test", testFunc)
	resource, ok := rtr.routes.GetResource("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := resource.Handlers["POST"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function got %q, want %q", got, want)
	}
}

func TestPutAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Put("/test", testFunc)
	resource, ok := rtr.routes.GetResource("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := resource.Handlers["PUT"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function got %q, want %q", got, want)
	}
}

func TestPatchAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Patch("/test", testFunc)
	resource, ok := rtr.routes.GetResource("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := resource.Handlers["PATCH"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function got %q, want %q", got, want)
	}
}

func TestDeleteAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Delete("/test", testFunc)
	resource, ok := rtr.routes.GetResource("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := resource.Handlers["DELETE"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function got %q, want %q", got, want)
	}
}

func TestHeadAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Head("/test", testFunc)
	resource, ok := rtr.routes.GetResource("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := resource.Handlers["HEAD"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function got %q, want %q", got, want)
	}
}

func TestOptionsAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Options("/test", testFunc)
	resource, ok := rtr.routes.GetResource("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := resource.Handlers["OPTIONS"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function got %q, want %q", got, want)
	}
}

func TestWhenNoMatchingRouteServeHTTPReturns404NotFound(t *testing.T) {
	want := 404
	rtr := Router{}
	rtr.routes = newMockRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)

	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}
}

func TestWhenNoMatchingHandlerServeHTTPReturns405MethodNotAllowed(t *testing.T) {
	want := 405
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Delete("/test", testFunc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)

	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}
}

func TestWhenNoMatchingHandlerForOPTIONSRequestAndAutoPopulateOptionsAllow(t *testing.T) {
	want := "DELETE, GET, POST"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AutoPopulateOptionsAllow = true
	rtr.Get("/test", testFunc)
	rtr.Post("/test", testFunc)
	rtr.Delete("/test", testFunc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	rtr.ServeHTTP(w, req)
	sort.Strings(w.HeaderMap["Allow"])
	got := strings.Join(w.HeaderMap["Allow"], ", ")
	if got != want {
		t.Errorf("Allow got %q, want %q", got, want)
	}
}

func TestWhenNoMatchingHandlerRequestAndAutoPopulateOptionsAllow(t *testing.T) {
	want := 405
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AutoPopulateOptionsAllow = true
	rtr.Get("/test", testFunc)
	rtr.Post("/test", testFunc)
	rtr.Delete("/test", testFunc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/test", nil)
	rtr.ServeHTTP(w, req)

	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}
}

func TestWhenNoErrorAndNoStatusSetServeHTTPReturns200OK(t *testing.T) {
	want := 200
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Get("/test", testFunc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)

	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}
}

func TestAddPreHookAppendsToHookList(t *testing.T) {
	want := "prehookhandler"
	callStack := ""
	ph := func(ctx *Context) {
		callStack += "prehook"
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AddPreHook(ph)
	rtr.Get("/test", func(ctx *Context) {
		callStack += "handler"
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)

	got := callStack
	if got != want {
		t.Errorf("Status got %q, want %q", got, want)
	}
}

func TestPreHooksCalledInOrder(t *testing.T) {
	want := "prehook1prehook2prehook3handler"
	callStack := ""
	ph1 := func(ctx *Context) {
		callStack += "prehook1"
	}
	ph2 := func(ctx *Context) {
		callStack += "prehook2"
	}
	ph3 := func(ctx *Context) {
		callStack += "prehook3"
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AddPreHook(ph1)
	rtr.AddPreHook(ph2)
	rtr.AddPreHook(ph3)
	rtr.Get("/test", func(ctx *Context) {
		callStack += "handler"
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)

	got := callStack
	if got != want {
		t.Errorf("Status got %q, want %q", got, want)
	}
}

func TestAddPostHookAppendsToHookList(t *testing.T) {
	want := "handlerposthook"
	callStack := ""
	ph := func(ctx *Context) {
		callStack += "posthook"
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AddPostHook(ph)
	rtr.Get("/test", func(ctx *Context) {
		callStack += "handler"
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)

	got := callStack
	if got != want {
		t.Errorf("Status got %q, want %q", got, want)
	}
}

func TestPostHooksCalledInOrder(t *testing.T) {
	want := "handlerposthook1posthook2posthook3"
	callStack := ""
	ph1 := func(ctx *Context) {
		callStack += "posthook1"
	}
	ph2 := func(ctx *Context) {
		callStack += "posthook2"
	}
	ph3 := func(ctx *Context) {
		callStack += "posthook3"
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AddPostHook(ph1)
	rtr.AddPostHook(ph2)
	rtr.AddPostHook(ph3)
	rtr.Get("/test", func(ctx *Context) {
		callStack += "handler"
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)

	got := callStack
	if got != want {
		t.Errorf("Status got %q, want %q", got, want)
	}
}

func TestPostHookResponsesAreIgnored(t *testing.T) {
	ph := func(ctx *Context) {
		ctx.RespondWith("with biscuits").WithStatus(204)
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AddPostHook(ph)
	rtr.Get("/test", func(ctx *Context) {
		ctx.RespondWith("Mango trees").WithStatus(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)
	want := 200
	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}

	wantBody := "Mango trees"
	gotBody := w.Body.String()
	if gotBody != wantBody {
		t.Errorf("Body got %q, want %q", gotBody, wantBody)
	}
}

func TestPreHookWriterResponsesPreventMainHandlerRunning(t *testing.T) {
	ph := func(ctx *Context) {
		ctx.Writer.WriteHeader(204)
		ctx.Writer.Write([]byte("with biscuits"))
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AddPreHook(ph)
	rtr.Get("/test", func(ctx *Context) {
		ctx.RespondWith("Mango trees").WithStatus(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)
	want := 204
	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}

	wantBody := "with biscuits"
	gotBody := w.Body.String()
	if gotBody != wantBody {
		t.Errorf("Body got %q, want %q", gotBody, wantBody)
	}
}

func TestPreHookContextResponsesPreventMainHandlerRunning(t *testing.T) {
	ph := func(ctx *Context) {
		ctx.RespondWith("with biscuits").WithStatus(204)
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AddPreHook(ph)
	rtr.Get("/test", func(ctx *Context) {
		ctx.RespondWith("Mango trees").WithStatus(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)
	want := 204
	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}

	wantBody := "with biscuits"
	gotBody := w.Body.String()
	if gotBody != wantBody {
		t.Errorf("Body got %q, want %q", gotBody, wantBody)
	}
}

func TestPreHookWriterResponsesPreventSubsequentPreHooksRunning(t *testing.T) {
	ph1 := func(ctx *Context) {
		ctx.Writer.WriteHeader(204)
		ctx.Writer.Write([]byte("with biscuits"))
	}
	ph2 := func(ctx *Context) {
		t.Errorf("Subsequent PreHooks not ignored")
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AddPreHook(ph1)
	rtr.AddPreHook(ph2)
	rtr.Get("/test", func(ctx *Context) {
		ctx.RespondWith("Mango trees").WithStatus(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)
}

func TestPreHookContextResponsesPreventSubsequentPreHooksRunning(t *testing.T) {
	ph1 := func(ctx *Context) {
		ctx.RespondWith("with biscuits").WithStatus(204)
	}
	ph2 := func(ctx *Context) {
		t.Errorf("Subsequent PreHooks not ignored")
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.AddPreHook(ph1)
	rtr.AddPreHook(ph2)
	rtr.Get("/test", func(ctx *Context) {
		ctx.RespondWith("Mango trees").WithStatus(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)
}

func TestPreHookContextResponsesCanSerializeModel(t *testing.T) {
	want := "{\"Name\":\"Mango\",\"Size\":34,\"Edible\":true}\n"
	ph := func(ctx *Context) {
		m := struct {
			Name   string
			Size   int
			Edible bool
		}{
			"Mango", 34, true,
		}
		ctx.RespondWith(m)
	}

	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.EncoderEngine = newEncoderEngine()
	rtr.AddPreHook(ph)
	rtr.Get("/test", func(ctx *Context) {
		ctx.RespondWith("Mango trees").WithStatus(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept", "application/json")
	rtr.ServeHTTP(w, req)

	got := w.Body.String()
	if got != want {
		t.Errorf("Body got %s, want %s", got, want)
	}
}

func TestGetSimpleTextResponse(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()

	rtr.Get("/test", func(ctx *Context) {
		ctx.RespondWith("We're just two lost souls swimming in a fish bowl")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	rtr.ServeHTTP(w, req)
	want := "We're just two lost souls swimming in a fish bowl"
	got := w.Body.String()
	if got != want {
		t.Errorf("Response got %q, want %q", got, want)
	}
}

func TestGetResponseStatus(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()

	rtr.Get("/test", func(ctx *Context) {
		ctx.Respond().WithStatus(404)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	rtr.ServeHTTP(w, req)
	want := 404
	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}
}

func TestGetEncodedResponse(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.EncoderEngine = &mockEncoderEngine{}

	var model = struct {
		a string
		b string
		c int
	}{
		"mango", "biscuits", 34,
	}

	rtr.Get("/test", func(ctx *Context) {
		ctx.Respond().WithModel(fmt.Sprint(model))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept", "test/test")
	rtr.ServeHTTP(w, req)

	want := fmt.Sprint(model)
	got := w.Body.String()
	if got != want {
		t.Errorf("Response got %q, want %q", got, want)
	}
}

func TestResponseCodeWhenRequestAcceptHeaderIsUnsupported(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.EncoderEngine = &mockEncoderEngine{}

	var model = struct {
		a string
		b string
		c int
	}{
		"mango", "biscuits", 34,
	}

	rtr.Get("/test", func(ctx *Context) {
		ctx.Respond().WithModel(fmt.Sprint(model))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept", "test/mango")
	rtr.ServeHTTP(w, req)

	want := 406
	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}
}

func TestResponseMessageWhenRequestAcceptHeaderIsUnsupported(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.EncoderEngine = &mockEncoderEngine{}

	var model = struct {
		a string
		b string
		c int
	}{
		"mango", "biscuits", 34,
	}

	rtr.Get("/test", func(ctx *Context) {
		ctx.Respond().WithModel(fmt.Sprint(model))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept", "test/mango")
	rtr.ServeHTTP(w, req)

	want := "Unable to encode to requested acceptable formats: \"test/mango\"\n"
	got := w.Body.String()
	if got != want {
		t.Errorf("Error message got %q, want %q", got, want)
	}
}

func TestResponseCodeWhenErrorEncodingPayload(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.EncoderEngine = &mockEncoderEngine{}

	var model = struct {
		a string
		b string
		c int
	}{
		"mango", "biscuits", 34,
	}

	rtr.Get("/test", func(ctx *Context) {
		ctx.Respond().WithModel(model)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept", "test/test")
	rtr.ServeHTTP(w, req)

	want := 500
	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}
}

func TestResponseMessageWhenErrorEncodingPayload(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.EncoderEngine = &mockEncoderEngine{}

	var model = struct {
		a string
		b string
		c int
	}{
		"mango", "biscuits", 34,
	}

	rtr.Get("/test", func(ctx *Context) {
		ctx.Respond().WithModel(model)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept", "test/test")
	rtr.ServeHTTP(w, req)

	want := "Internal Server Error\n"
	got := w.Body.String()
	if got != want {
		t.Errorf("Error message got %q, want %q", got, want)
	}
}

func TestGetStaticResponse(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()

	rtr.StaticDir("static")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test.html", nil)

	rtr.ServeHTTP(w, req)
	want := "<html>Test Static</html>\n"
	got := w.Body.String()
	if got != want {
		t.Errorf("Response got %q, want %q", got, want)
	}
}

func TestGetStaticResponseStatus(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()

	rtr.StaticDir("static")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test.html", nil)

	rtr.ServeHTTP(w, req)
	want := 200
	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}
}

func TestGetStaticRedirectsWhenIndexHtml(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()

	rtr.StaticDir("static")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/index.html", nil)

	rtr.ServeHTTP(w, req)
	want := 301
	got := w.Code
	if got != want {
		t.Errorf("Status got %d, want %d", got, want)
	}
}

func TestGetStaticRedirectsToRootWhenIndexHtml(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()

	rtr.StaticDir("static")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/index.html", nil)

	rtr.ServeHTTP(w, req)
	want := "./"
	fmt.Println(w.HeaderMap)
	got := w.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location got %q, want %q", got, want)
	}

}

func TestGetStaticRootReturnsIndexHtml(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()

	rtr.StaticDir("static")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "", nil)

	rtr.ServeHTTP(w, req)
	want := "<html>Index Static</html>\n"
	got := w.Body.String()
	if got != want {
		t.Errorf("Response got %q, want %q", got, want)
	}
}

func TestNewRouterSetsValidationHandler(t *testing.T) {
	want := reflect.TypeOf(&elementValidationHandler{}).String()
	r := NewRouter()
	if r.ValidationHandler == nil {
		t.Errorf("ValidationHandler type = \"<nil>\", want %q", want)
		return
	}
	got := reflect.TypeOf(r.ValidationHandler).String()
	if got != want {
		t.Errorf("ValidationHandler got %q, want %q", got, want)
	}
}

func TestNewRouterSetsRoutes(t *testing.T) {
	want := reflect.TypeOf(&tree{}).String()
	r := NewRouter()
	if r.routes == nil {
		t.Errorf("Routes type = \"<nil>\", want %q", want)
		return
	}
	got := reflect.TypeOf(r.routes).String()
	if got != want {
		t.Errorf("Routes type got %q, want %q", got, want)
	}
}

func TestNewRouterSetsRoutesValidationHandler(t *testing.T) {
	want := reflect.TypeOf(&elementValidationHandler{}).String()
	r := NewRouter()
	if r.routes == nil {
		t.Errorf("Routes type = \"<nil>\", want %q", want)
		return
	}
	tr := r.routes.(*tree)
	got := reflect.TypeOf(tr.validators).String()
	if got != want {
		t.Errorf("Routes type got %q, want %q", got, want)
	}
}

func TestNewRouterSetsEncoderEngine(t *testing.T) {
	want := reflect.TypeOf(&encoderEngine{}).String()
	r := NewRouter()
	if r.EncoderEngine == nil {
		t.Errorf("EncoderEngine type = \"<nil>\", want %q", want)
		return
	}
	got := reflect.TypeOf(r.EncoderEngine).String()
	if got != want {
		t.Errorf("EncoderEngine got %q, want %q", got, want)
	}
}

func TestNewRouterInitialisesEncoderEngineWithDefaultMediaType(t *testing.T) {
	want := DefaultMediaType
	r := NewRouter()
	got := r.EncoderEngine.DefaultMediaType()
	if got != want {
		t.Errorf("EncoderEngine.DefaultMediaType got %q, want %q", got, want)
	}
}

func TestNewRouterSetsAutoPopulateOptionsAllowToTrue(t *testing.T) {
	want := true
	r := NewRouter()
	got := r.AutoPopulateOptionsAllow
	if got != want {
		t.Errorf("AutoPopulateOptionsAllow = %t, want %t", got, want)
	}
}

func TestRegisterModulesWithEmptyModuleRegistersNoNewRoutes(t *testing.T) {
	want := 0
	r := Router{}
	r.routes = newMockRoutes()

	r.RegisterModules([]Registerer{
		emptyTestModule{},
	})
	got := len(r.routes.(*mockRoutes).routes)
	if got != want {
		t.Errorf("Route count got %d, want %d", got, want)
	}
}

func TestRegisterModulesWithSingleModuleRegistersRoutes(t *testing.T) {
	want := 1
	r := Router{}
	r.routes = newMockRoutes()

	r.RegisterModules([]Registerer{
		singleRouteTestModule{},
	})
	s, _ := r.routes.GetResource("/single")
	got := len(s.Handlers)
	if got != want {
		t.Errorf("Route count got %d, want %d", got, want)
	}
}

func TestRegisterModulesWithMultipleModulesRegistersRoutes(t *testing.T) {
	want := 3
	r := Router{}
	r.routes = newMockRoutes()

	r.RegisterModules([]Registerer{
		singleRouteTestModule{},
		multiRouteTestModule{},
	})
	s, _ := r.routes.GetResource("/single")
	m, _ := r.routes.GetResource("/multi")
	got := len(s.Handlers) + len(m.Handlers)

	if got != want {
		t.Errorf("Route count got %d, want %d", got, want)
	}
}

func TestRegisterModulesDoesNotAffectExisingRegistrations(t *testing.T) {
	want := 3
	r := Router{}
	r.routes = newMockRoutes()
	r.Get("/single", testFunc)

	r.RegisterModules([]Registerer{
		multiRouteTestModule{},
	})

	s, _ := r.routes.GetResource("/single")
	m, _ := r.routes.GetResource("/multi")
	got := len(s.Handlers) + len(m.Handlers)

	if got != want {
		t.Errorf("Route count got %d, want %d", got, want)
	}
}

func TestRegisterModulesPanicsWhenAttemptingDuplicateRoute(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	r := Router{}
	r.routes = newMockRoutes()

	r.RegisterModules([]Registerer{
		singleRouteTestModule{},
		duplicateRouteTestModule{},
	})
}

func encFunc(io.Writer) Encoder {
	return nil
}

func decFunc(io.ReadCloser) Decoder {
	return nil
}

func TestRouterAddValidator(t *testing.T) {
	want := true
	r := Router{}
	evh := elementValidationHandler{}
	evh.validators = make(map[string]Validator)
	r.ValidationHandler = &evh
	r.AddValidator(Int32Validator{})
	_, valid := r.IsValid("123", "int32")
	got := valid
	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestRouterAddValidators(t *testing.T) {
	want := true
	r := Router{}
	evh := elementValidationHandler{}
	evh.validators = make(map[string]Validator)
	r.ValidationHandler = &evh
	r.AddValidators([]Validator{Int32Validator{}})
	_, valid := r.IsValid("123", "int32")
	got := valid
	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestRouterAddValidatorPanicsIfConstraintConflicts(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "conflicting constraint type: int32"
			got := r
			if got != want {
				t.Errorf("Error message got %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()

	r := Router{}
	evh := elementValidationHandler{}
	evh.validators = make(map[string]Validator)
	r.ValidationHandler = &evh
	r.AddValidator(Int32Validator{})
	r.AddValidator(dupValidator{})
}

func TestRouterAddValidatorsPanicsIfConstraintConflicts(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "conflicting constraint type: int32"
			got := r
			if got != want {
				t.Errorf("Error message got %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()

	r := Router{}
	evh := elementValidationHandler{}
	evh.validators = make(map[string]Validator)
	r.ValidationHandler = &evh
	r.AddValidators([]Validator{
		Int32Validator{},
		dupValidator{},
	})
}

func TestRouterRequestLoggingWhenNotFound(t *testing.T) {
	// no routes, so code: 404 and "404 page not found\n" - 19 bytes
	want := "127.0.0.1 - - [11/Oct/2000:13:55:36 -0700] \"GET /spaceweasel/mango/stone.png HTTP/1.1\" 404 19"

	location := time.FixedZone("test", -25200)
	start := time.Date(2000, 10, 11, 13, 55, 36, 0, location)
	nowUTC = func() time.Time {
		return start
	}

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	req.RemoteAddr = "127.0.0.1"
	w := httptest.NewRecorder()
	got := ""
	r := Router{}
	r.RequestLogger = func(l *RequestLog) {
		got = l.CommonFormat()
		if got != want {
			t.Errorf("Log got %q, want %q", got, want)
		}
	}
	r.routes = newMockRoutes()
	r.ServeHTTP(w, req)
}

func TestRouterRequestLoggingWithUserHandler(t *testing.T) {
	want := "127.0.0.1 - - [11/Oct/2000:13:55:36 -0700] \"GET /mango HTTP/1.1\" 200 19"

	location := time.FixedZone("test", -25200)
	start := time.Date(2000, 10, 11, 13, 55, 36, 0, location)
	nowUTC = func() time.Time {
		return start
	}

	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.RemoteAddr = "127.0.0.1"
	w := httptest.NewRecorder()
	got := ""
	r := Router{}
	r.RequestLogger = func(l *RequestLog) {
		got = l.CommonFormat()
		if got != want {
			t.Errorf("Log got %q, want %q", got, want)
		}
	}
	r.routes = newMockRoutes()
	r.routes.AddHandlerFunc("/mango", "GET", func(c *Context) {
		c.RespondWith("A mango in the hand")
	})
	r.ServeHTTP(w, req)
}

func TestRouterRequestLoggingWithNoContentResponse(t *testing.T) {
	want := "127.0.0.1 - - [11/Oct/2000:13:55:36 -0700] \"GET /mango HTTP/1.1\" 204 0"

	location := time.FixedZone("test", -25200)
	start := time.Date(2000, 10, 11, 13, 55, 36, 0, location)
	nowUTC = func() time.Time {
		return start
	}

	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.RemoteAddr = "127.0.0.1"
	w := httptest.NewRecorder()
	got := ""
	r := Router{}
	r.RequestLogger = func(l *RequestLog) {
		got = l.CommonFormat()
		if got != want {
			t.Errorf("Log got %q, want %q", got, want)
		}
	}
	r.routes = newMockRoutes()
	r.routes.AddHandlerFunc("/mango", "GET", func(c *Context) {
		c.RespondWith(204)
	})
	r.ServeHTTP(w, req)
}

func TestRouterRequestLoggingWhenUnRecoveredPanic(t *testing.T) {
	msg := "Internal Server Error\n"
	bCount := strconv.Itoa(len(msg))
	want := "127.0.0.1 - - [11/Oct/2000:13:55:36 -0700] \"GET /mango HTTP/1.1\" 500 " + bCount

	location := time.FixedZone("test", -25200)
	start := time.Date(2000, 10, 11, 13, 55, 36, 0, location)
	nowUTC = func() time.Time {
		return start
	}

	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.RemoteAddr = "127.0.0.1"
	w := httptest.NewRecorder()
	r := Router{}
	ch := make(chan string)
	r.RequestLogger = func(l *RequestLog) {
		ch <- l.CommonFormat()
	}
	r.routes = newMockRoutes()
	r.routes.AddHandlerFunc("/mango", "GET", func(c *Context) {
		panic("what no mangoes!")
	})

	r.ServeHTTP(w, req)

	select {
	case got := <-ch:
		if got != want {
			t.Errorf("Log got %q, want %q", got, want)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out")
	}
}

func TestRouterRequestLoggerIsUpdatedWhenAuthenticated(t *testing.T) {
	want := "Mungo"

	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.RemoteAddr = "127.0.0.1"
	w := httptest.NewRecorder()
	got := ""
	r := Router{}
	r.RequestLogger = func(l *RequestLog) {
		got = l.UserID
		if got != want {
			t.Errorf("UserID got %q, want %q", got, want)
		}
	}
	r.routes = newMockRoutes()
	r.routes.AddHandlerFunc("/mango", "GET", func(c *Context) {
		//c.RespondWith("A mango in the hand")
	})
	r.AddPreHook(func(c *Context) {
		c.Identity = BasicIdentity{Username: "Mungo"}
	})
	r.ServeHTTP(w, req)
}

func TestRouterErrorLoggingMsgHasSummaryAsFirstLineWhenUnRecoveredPanic(t *testing.T) {
	want := "what no mangoes!"
	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	w := httptest.NewRecorder()
	r := Router{}
	ch := make(chan error)
	r.ErrorLogger = func(err error) {
		ch <- err
	}
	r.routes = newMockRoutes()
	r.routes.AddHandlerFunc("/mango", "GET", func(c *Context) {
		panic("what no mangoes!")
	})

	r.ServeHTTP(w, req)

	select {
	case err := <-ch:
		msg := err.Error()
		lines := strings.Split(msg, "\n")
		got := lines[0]
		if got != want {
			t.Errorf("Error Summary got %q, want %q", got, want)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out")
	}
}

func TestRouterErrorLoggingMsgHasReqDetailAsSecondLineWhenUnRecoveredPanic(t *testing.T) {
	want := "127.0.0.1 - - [11/Oct/2000:13:55:36 -0700] \"GET /mango HTTP/1.1\" 0 0"

	location := time.FixedZone("test", -25200)
	start := time.Date(2000, 10, 11, 13, 55, 36, 0, location)
	nowUTC = func() time.Time {
		return start
	}

	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
	req.RemoteAddr = "127.0.0.1"
	w := httptest.NewRecorder()
	r := Router{}
	ch := make(chan error)
	r.ErrorLogger = func(err error) {
		ch <- err
	}
	r.routes = newMockRoutes()
	r.routes.AddHandlerFunc("/mango", "GET", func(c *Context) {
		panic("what no mangoes!")
	})

	r.ServeHTTP(w, req)

	select {
	case err := <-ch:
		msg := err.Error()
		lines := strings.Split(msg, "\n")
		got := lines[1]
		if got != want {
			t.Errorf("Error request got %q, want %q", got, want)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out")
	}
}

func TestSetCORSForwardsToTree(t *testing.T) {
	want := "http://greencheese.com"

	mr := newMockRoutes()
	r := Router{
		routes: mr,
	}
	config := CORSConfig{
		Origins: []string{"http://greencheese.com"},
		Methods: []string{"POST", "PATCH"},
		Headers: []string{"X-Cheese", "X-Mangoes"},
	}
	r.SetCORS("/mango", config)

	c, ok := mr.corsConfigs["/mango"]
	if !ok {
		t.Errorf("CORS Config not sent")
		return
	}

	got := strings.Join(c.Origins, ", ")

	if got != want {
		t.Errorf("Origins got %q, want %q", got, want)
	}
}

func TestSetGlobalCORSForwardsToTree(t *testing.T) {
	want := "http://greencheese.com"

	mr := newMockRoutes()
	r := Router{
		routes: mr,
	}
	config := CORSConfig{
		Origins: []string{"http://greencheese.com"},
		Methods: []string{"POST", "PATCH"},
		Headers: []string{"X-Cheese", "X-Mangoes"},
	}
	r.SetGlobalCORS(config)

	gc := mr.globalCORSConfig
	got := strings.Join(gc.Origins, ", ")

	if got != want {
		t.Errorf("Origins got %q, want %q", got, want)
	}
}

func TestAddCORSForwardsToTree(t *testing.T) {
	want := "http://greencheese.com"

	mr := newMockRoutes()
	r := Router{
		routes: mr,
	}
	config := CORSConfig{
		Origins: []string{"http://greencheese.com"},
		Methods: []string{"POST", "PATCH"},
		Headers: []string{"X-Cheese", "X-Mangoes"},
	}
	r.AddCORS("/mango", config)

	c, ok := mr.corsConfigs["/mango"]
	if !ok {
		t.Errorf("CORS Config not sent")
		return
	}

	got := strings.Join(c.Origins, ", ")

	if got != want {
		t.Errorf("Origins got %q, want %q", got, want)
	}
}

func TestRouterSetsContextModelValidator(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	evh := elementValidationHandler{}
	evh.validators = make(map[string]Validator)
	rtr.ValidationHandler = &evh
	rtr.modelValidator = newModelValidator(rtr.ValidationHandler)

	rtr.Get("/test", func(ctx *Context) {
		if ctx.modelValidator == nil {
			t.Errorf("ModelValidator = <nil>, want %T", contextModelValidator{})
		}
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	rtr.ServeHTTP(w, req)
}

func TestRouterAddModelValidatorAddsToCollection(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	evh := elementValidationHandler{}
	evh.validators = make(map[string]Validator)
	rtr.ValidationHandler = &evh
	rtr.modelValidator = newModelValidator(rtr.ValidationHandler)

	type model struct {
		Name string
		Age  int
	}
	custValidator := func(m interface{}) (map[string][]ValidationFailure, bool) {
		results := make(map[string][]ValidationFailure)
		cm := m.(model)
		if cm.Name != "Mango" {
			results["Name"] = []ValidationFailure{
				{Code: "wrongname",
					Message: "Name must be Mango"},
			}
		}
		return results, len(results) == 0
	}

	rtr.AddModelValidator(model{}, custValidator)

	rtr.Get("/test", func(ctx *Context) {
		testModel := model{"Mingo", 45}
		details, ok := ctx.Validate(testModel)
		if ok {
			t.Errorf("Validate result = true, want false")
			return
		}
		want := "Name must be Mango"
		got := details["Name"][0].Message
		if got != want {
			t.Errorf("Validate message got %q, want %q", got, want)
			return
		}
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	rtr.ServeHTTP(w, req)
}

func TestRouterNewValidatorsAddedAreAvailableToModelValidator(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	evh := elementValidationHandler{}
	evh.validators = make(map[string]Validator)
	rtr.ValidationHandler = &evh
	rtr.modelValidator = newModelValidator(rtr.ValidationHandler)

	rtr.AddValidator(CheeseValidator{})

	type model struct {
		Name string `validate:"cheese"`
		Age  int
	}

	rtr.Get("/test", func(ctx *Context) {
		testModel := model{"Mingo", 45}
		details, ok := ctx.Validate(testModel)
		if ok {
			t.Errorf("Validate result = true, want false")
			return
		}
		want := "Must contain the word cheese."
		got := details["Name"][0].Message
		if got != want {
			t.Errorf("Validate message got %q, want %q", got, want)
			return
		}
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	rtr.ServeHTTP(w, req)
}

type CheeseValidator struct{}

func (CheeseValidator) Validate(val interface{}, params []string) bool {
	s := val.(string)
	return strings.Contains(s, "cheese")
}

func (CheeseValidator) Type() string {
	return "cheese"
}

func (CheeseValidator) FailureMsg() string {
	return "Must contain the word cheese."
}

// // TODO:
// func TestRouterCallsHandleCORS(t *testing.T) {
//
// 	req, _ := http.NewRequest("GET", "https://somewhere.com/mango", nil)
// 	req.RemoteAddr = "127.0.0.1"
// 	w := httptest.NewRecorder()
// 	got := ""
// 	r := Router{}
// 	r.routes = newMockRoutes()
// 	r.routes.AddHandlerFunc("/mango", "GET", func(c *Context) {
// 		c.RespondWith("A mango in the hand")
// 	})
// 	r.ServeHTTP(w, req)
// }

type emptyTestModule struct{}

func (t emptyTestModule) Register(r *Router) {}

type singleRouteTestModule struct{}

func (t singleRouteTestModule) Register(r *Router) {
	r.Get("/single", testFunc)
}

type multiRouteTestModule struct{}

func (t multiRouteTestModule) Register(r *Router) {
	r.Put("/multi", testFunc2)
	r.Post("/multi", testFunc3)
}

type duplicateRouteTestModule struct{}

func (t duplicateRouteTestModule) Register(r *Router) {
	r.Get("/single", testFunc)
}
