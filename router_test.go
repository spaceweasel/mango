package mango

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type mockRoutes struct {
	routes     map[string]map[string]ContextHandlerFunc
	validators map[string]ParamValidator
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

func (m *mockRoutes) HandlerFuncs(path string) (map[string]ContextHandlerFunc, map[string]string, bool) {
	hm, ok := m.routes[path]
	return hm, nil, ok
}

func (m *mockRoutes) AddRouteParamValidator(v ParamValidator) {
	if _, ok := m.validators[v.Type()]; ok {
		panic("conflicting constraint type: " + v.Type())
	}
	m.validators[v.Type()] = v
}

func (m *mockRoutes) AddRouteParamValidators(validators []ParamValidator) {
	for _, v := range validators {
		m.AddRouteParamValidator(v)
	}
}

func newMockRoutes() *mockRoutes {
	mr := mockRoutes{}
	mr.routes = make(map[string]map[string]ContextHandlerFunc)
	mr.validators = make(map[string]ParamValidator)
	return &mr
}

func TestGetAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Get("/test", testFunc)
	handlers, _, ok := rtr.routes.HandlerFuncs("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := handlers["GET"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function = %q, want %q", got, want)
	}
}

func TestPostAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Post("/test", testFunc)
	handlers, _, ok := rtr.routes.HandlerFuncs("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := handlers["POST"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function = %q, want %q", got, want)
	}
}

func TestPutAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Put("/test", testFunc)
	handlers, _, ok := rtr.routes.HandlerFuncs("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := handlers["PUT"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function = %q, want %q", got, want)
	}
}

func TestPatchAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Patch("/test", testFunc)
	handlers, _, ok := rtr.routes.HandlerFuncs("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := handlers["PATCH"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function = %q, want %q", got, want)
	}
}

func TestDeleteAddsHandlerToRoutes(t *testing.T) {
	want := "testFunc"
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Del("/test", testFunc)
	handlers, _, ok := rtr.routes.HandlerFuncs("/test")
	if !ok {
		t.Errorf("Handler not added")
	}
	h := handlers["DELETE"]
	got := extractFnName(h)
	if got != want {
		t.Errorf("Handler function = %q, want %q", got, want)
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
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestWhenNoMatchingHandlerServeHTTPReturns405MethodNotAllowed(t *testing.T) {
	want := 405
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.Del("/test", testFunc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	rtr.ServeHTTP(w, req)

	got := w.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
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
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestAddPreHookAppendsToHookList(t *testing.T) {
	want := "prehookhandler"
	callStack := ""
	ph := func(ctx *Context) error {
		callStack += "prehook"
		return nil
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
		t.Errorf("Status = %q, want %q", got, want)
	}
}

func TestPreHooksCalledInOrder(t *testing.T) {
	want := "prehook1prehook2prehook3handler"
	callStack := ""
	ph1 := func(ctx *Context) error {
		callStack += "prehook1"
		return nil
	}
	ph2 := func(ctx *Context) error {
		callStack += "prehook2"
		return nil
	}
	ph3 := func(ctx *Context) error {
		callStack += "prehook3"
		return nil
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
		t.Errorf("Status = %q, want %q", got, want)
	}
}

func TestAddPostHookAppendsToHookList(t *testing.T) {
	want := "handlerposthook"
	callStack := ""
	ph := func(ctx *Context) error {
		callStack += "posthook"
		return nil
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
		t.Errorf("Status = %q, want %q", got, want)
	}
}

func TestPostHooksCalledInOrder(t *testing.T) {
	want := "handlerposthook1posthook2posthook3"
	callStack := ""
	ph1 := func(ctx *Context) error {
		callStack += "posthook1"
		return nil
	}
	ph2 := func(ctx *Context) error {
		callStack += "posthook2"
		return nil
	}
	ph3 := func(ctx *Context) error {
		callStack += "posthook3"
		return nil
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
		t.Errorf("Status = %q, want %q", got, want)
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
		t.Errorf("Response = %q, want %q", got, want)
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
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestGetEncodedResponse(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.encoderEngine = &mockEncoderEngine{}

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
		t.Errorf("Response = %q, want %q", got, want)
	}
}

func TestResponseCodeWhenRequestAcceptHeaderIsUnsupported(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.encoderEngine = &mockEncoderEngine{}

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
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestResponseMessageWhenRequestAcceptHeaderIsUnsupported(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.encoderEngine = &mockEncoderEngine{}

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
		t.Errorf("Error message = %q, want %q", got, want)
	}
}

func TestResponseCodeWhenErrorEncodingPayload(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.encoderEngine = &mockEncoderEngine{}

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
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestResponseMessageWhenErrorEncodingPayload(t *testing.T) {
	rtr := Router{}
	rtr.routes = newMockRoutes()
	rtr.encoderEngine = &mockEncoderEngine{}

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
		t.Errorf("Error message = %q, want %q", got, want)
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
		t.Errorf("Routes type = %q, want %q", got, want)
	}
}

func TestNewRouterSetsEncoderEngine(t *testing.T) {
	want := reflect.TypeOf(&encoderEngine{}).String()
	r := NewRouter()
	if r.encoderEngine == nil {
		t.Errorf("EncoderEngine type = \"<nil>\", want %q", want)
		return
	}
	got := reflect.TypeOf(r.encoderEngine).String()
	if got != want {
		t.Errorf("EncoderEngine = %q, want %q", got, want)
	}
}

func TestNewRouterInitialisesEncoderEngineWithDefaultMediaType(t *testing.T) {
	want := DefaultMediaType
	r := NewRouter()
	got := r.encoderEngine.DefaultMediaType()
	if got != want {
		t.Errorf("EncoderEngine.DefaultMediaType = %q, want %q", got, want)
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
		t.Errorf("Route count = %d, want %d", got, want)
	}
}

func TestRegisterModulesWithSingleModuleRegistersRoutes(t *testing.T) {
	want := 1
	r := Router{}
	r.routes = newMockRoutes()

	r.RegisterModules([]Registerer{
		singleRouteTestModule{},
	})
	s, _, _ := r.routes.HandlerFuncs("/single")
	got := len(s)
	if got != want {
		t.Errorf("Route count = %d, want %d", got, want)
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
	s, _, _ := r.routes.HandlerFuncs("/single")
	m, _, _ := r.routes.HandlerFuncs("/multi")
	got := len(s) + len(m)

	if got != want {
		t.Errorf("Route count = %d, want %d", got, want)
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

	s, _, _ := r.routes.HandlerFuncs("/single")
	m, _, _ := r.routes.HandlerFuncs("/multi")
	got := len(s) + len(m)

	if got != want {
		t.Errorf("Route count = %d, want %d", got, want)
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

func TestAddEncoderFuncForwardsRequestToEncoderEngine(t *testing.T) {
	want := "mango/test-encFunc"
	r := Router{}
	ee := mockEncoderEngine{}
	r.encoderEngine = &ee

	r.AddEncoderFunc("mango/test", encFunc)
	if len(ee.EncoderRequests) == 0 {
		t.Errorf("Requests = 0, want 1")
		return
	}
	got := ee.EncoderRequests[0]
	if got != want {
		t.Errorf("Recorded request = %q, want %q", got, want)
	}
}

func TestAddEncoderFuncCapturesEncoderEngineError(t *testing.T) {
	want := "error/error"
	r := Router{}
	r.encoderEngine = &mockEncoderEngine{}

	err := r.AddEncoderFunc("error/error", encFunc)
	if err == nil {
		t.Errorf("Error = <nil>, want %q", want)
		return
	}
	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

func decFunc(io.ReadCloser) Decoder {
	return nil
}

func TestAddDecoderFuncForwardsRequestToEncoderEngine(t *testing.T) {
	want := "mango/test-decFunc"
	r := Router{}
	ee := mockEncoderEngine{}
	r.encoderEngine = &ee

	r.AddDecoderFunc("mango/test", decFunc)
	if len(ee.DecoderRequests) == 0 {
		t.Errorf("Requests = 0, want 1")
		return
	}
	got := ee.DecoderRequests[0]
	if got != want {
		t.Errorf("Recorded request = %q, want %q", got, want)
	}
}

func TestAddDecoderFuncCapturesEncoderEngineError(t *testing.T) {
	want := "error/error"
	r := Router{}
	r.encoderEngine = &mockEncoderEngine{}

	err := r.AddDecoderFunc("error/error", decFunc)
	if err == nil {
		t.Errorf("Error = <nil>, want %q", want)
		return
	}
	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

func TestRouterAddRouteParamValidator(t *testing.T) {
	want := true
	r := Router{}
	routes := newMockRoutes()
	r.routes = routes
	r.AddRouteParamValidator(Int32Validator{})
	valid := routes.TestValidators("123", "int32")
	got := valid
	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestRouterAddRouteParamValidators(t *testing.T) {
	want := true
	r := Router{}
	routes := newMockRoutes()
	r.routes = routes
	r.AddRouteParamValidators([]ParamValidator{Int32Validator{}})
	valid := routes.TestValidators("123", "int32")
	got := valid
	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestRouterAddRouteParamValidatorPanicsIfConstraintConflicts(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "conflicting constraint type: int32"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()

	r := Router{}
	routes := newMockRoutes()
	r.routes = routes

	r.AddRouteParamValidator(Int32Validator{})
	r.AddRouteParamValidator(dupValidator{})
}

func TestRouterAddRouteParamValidatorsPanicsIfConstraintConflicts(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "conflicting constraint type: int32"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()

	r := Router{}
	routes := newMockRoutes()
	r.routes = routes

	r.AddRouteParamValidators([]ParamValidator{
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
			t.Errorf("Log = %q, want %q", got, want)
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
			t.Errorf("Log = %q, want %q", got, want)
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
			t.Errorf("Log = %q, want %q", got, want)
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
			t.Errorf("Log = %q, want %q", got, want)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out")
	}
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
			t.Errorf("Error Summary = %q, want %q", got, want)
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
			t.Errorf("Error request = %q, want %q", got, want)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out")
	}
}

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
