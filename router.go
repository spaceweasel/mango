package mango

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
)

type routes interface {
	AddHandlerFunc(pattern, method string, handlerFunc ContextHandlerFunc)
	GetResource(path string) (*Resource, bool)
	SetGlobalCORS(config CORSConfig)
	SetCORS(pattern string, config CORSConfig)
	AddCORS(pattern string, config CORSConfig)
}

// RequestLogFunc is the signature for implementing router RequestLogger
type RequestLogFunc func(*RequestLog)

// Router is the main mango object. Router implements the standard
// library http.Handler interface, so it can be used in the call to
// http.ListenAndServe method.
// New Router objects should be created using the NewRouter method
// rather than creating a new uninitialised instance.
// TODO: Add more info here.
type Router struct {
	ValidationHandler
	EncoderEngine
	routes                   routes
	preHooks                 []ContextHandlerFunc
	postHooks                []ContextHandlerFunc
	RequestLogger            RequestLogFunc
	ErrorLogger              func(error)
	AutoPopulateOptionsAllow bool
	modelValidator           ModelValidator
	CompMinLength            int
	staticHandler            http.Handler
}

// AddModelValidator adds a custom model validator to the collection.
func (r *Router) AddModelValidator(m interface{}, fn ValidateFunc) {
	r.modelValidator.AddCustomValidator(m, fn)
}

// NewRouter returns a pointer to a new Router instance.
// The Router will be initialised with a new EncoderEngine,
// Validation handlers for route parameters and models,
// and route handling functionality.
func NewRouter() *Router {
	r := Router{}
	r.ValidationHandler = newValidationHandler()
	r.routes = newTree(r.ValidationHandler)
	r.modelValidator = newModelValidator(r.ValidationHandler)
	r.EncoderEngine = newEncoderEngine()
	r.AutoPopulateOptionsAllow = true
	r.CompMinLength = 300
	return &r
}

// SetGlobalCORS sets the CORS configuration that will be used for
// a resource if it has no CORS configuration of its own. If the
// resource has no CORSConfig and tree.GlobalCORSConfig is nil
// then CORS request are treated like any other.
func (r *Router) SetGlobalCORS(config CORSConfig) {
	r.routes.SetGlobalCORS(config)
}

// SetCORS sets the CORS configuration that will be used for
// the resource matching the pattern.
// These settings override any global settings.
func (r *Router) SetCORS(pattern string, config CORSConfig) {
	r.routes.SetCORS(pattern, config)
}

// AddCORS sets the CORS configuration that will be used for
// the resource matching the pattern, by merging the supplied
// config with any globalCORSConfig.
// SetGlobalCORS MUST be called before this method!
// AddCORS will panic if GlobalCORS is nil.
func (r *Router) AddCORS(pattern string, config CORSConfig) {
	r.routes.AddCORS(pattern, config)
}

// Get registers a new handlerFunc that will be called when HTTP GET
// requests are made to URLs with paths that match pattern.
// If a GET handlerFunc already exists for pattern, Get panics.
func (r *Router) Get(pattern string, handlerFunc ContextHandlerFunc) {
	r.routes.AddHandlerFunc(pattern, "GET", handlerFunc)
}

// Post registers a new handlerFunc that will be called when HTTP POST
// requests are made to URLs with paths that match pattern.
// If a POST handlerFunc already exists for pattern, Post panics.
func (r *Router) Post(pattern string, handlerFunc ContextHandlerFunc) {
	r.routes.AddHandlerFunc(pattern, "POST", handlerFunc)
}

// Put registers a new handlerFunc that will be called when HTTP PUT
// requests are made to URLs with paths that match pattern.
// If a PUT handlerFunc already exists for pattern, Put panics.
func (r *Router) Put(pattern string, handlerFunc ContextHandlerFunc) {
	r.routes.AddHandlerFunc(pattern, "PUT", handlerFunc)
}

// Patch registers a new handlerFunc that will be called when HTTP PATCH
// requests are made to URLs with paths that match pattern.
// If a PATCH handlerFunc already exists for pattern, Patch panics.
func (r *Router) Patch(pattern string, handlerFunc ContextHandlerFunc) {
	r.routes.AddHandlerFunc(pattern, "PATCH", handlerFunc)
}

// Delete registers a new handlerFunc that will be called when HTTP DELETE
// requests are made to URLs with paths that match pattern.
// If a DELETE handlerFunc already exists for pattern, Delete panics.
func (r *Router) Delete(pattern string, handlerFunc ContextHandlerFunc) {
	r.routes.AddHandlerFunc(pattern, "DELETE", handlerFunc)
}

// Head registers a new handlerFunc that will be called when HTTP HEAD
// requests are made to URLs with paths that match pattern.
// If a HEAD handlerFunc already exists for pattern, Head panics.
func (r *Router) Head(pattern string, handlerFunc ContextHandlerFunc) {
	r.routes.AddHandlerFunc(pattern, "HEAD", handlerFunc)
}

// Options registers a new handlerFunc that will be called when HTTP OPTIONS
// requests are made to URLs with paths that match pattern.
// If a OPTIONS handlerFunc already exists for pattern, Options panics.
func (r *Router) Options(pattern string, handlerFunc ContextHandlerFunc) {
	r.routes.AddHandlerFunc(pattern, "OPTIONS", handlerFunc)
}

// StaticDir sets a root directory for serving static files.
// Root can be an absolute or relative directory path.
// As a special case, any request ending in "/index.html" is redirected to the
// same path, without the final "index.html".
func (r *Router) StaticDir(root string) {
	r.staticHandler = http.FileServer(http.Dir(root))
}

// ServeHTTP dispatches the request to the handler whose pattern
// matches the request URL.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ae := req.Header.Get("Accept-Encoding")
	resp := NewResponseWriter(w, ae, r.CompMinLength)
	reqLog := NewRequestLog(req)
	defer func() {
		if r.RequestLogger == nil {
			return
		}
		reqLog.stop()
		reqLog.BytesOut = resp.byteCount
		reqLog.Status = resp.status
		// don't let logging hinder sending response
		go r.RequestLogger(reqLog)
	}()
	defer func() {
		// although the calling code handles panics, we'll do it
		// here so the RequestLogger can capture it too.
		if rec := recover(); rec != nil {
			http.Error(resp, "Internal Server Error", 500)
			if r.ErrorLogger != nil {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				go func() {
					buf = bytes.Trim(buf, "\x00")
					err := fmt.Errorf("%v\n%s\n%s\n", rec, reqLog.CommonFormat(), buf)
					r.ErrorLogger(err)
				}()
			}
		}
	}()

	resource, ok := r.routes.GetResource(req.URL.Path)
	if !ok {
		// try static files
		if r.staticHandler != nil {
			r.staticHandler.ServeHTTP(resp, req)
			return
		}

		http.NotFound(resp, req)
		return
	}

	if handleCORS(req, resp, resource) {
		return
	}

	fn, ok := resource.Handlers[req.Method]
	if !ok {
		if r.AutoPopulateOptionsAllow {
			// if a dedicated OPTIONS handler hasn't been added to the resource
			// then just return with ALLOW header.
			for k := range resource.Handlers {
				resp.Header().Add("Allow", k)
			}
		}
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	c := &Context{
		Request:        req,
		Writer:         resp,
		RouteParams:    resource.RouteParams,
		encoderEngine:  r.EncoderEngine,
		modelValidator: r.modelValidator,
	}

	//call prehooks
	for _, h := range r.preHooks {
		h(c)
		if resp.responded || c.responseReady {
			break
		}
	}

	// update request log with any post prehook changes
	if c.Identity != nil {
		reqLog.Identity = c.Identity
		reqLog.UserID = c.Identity.UserID()
	}
	reqLog.Context = c.X

	// TODO: record name of handler function in reqLog
	// handlerName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	// only run handler if a prehook hasn't responded already
	if !resp.responded && !c.responseReady {
		fn.ServeHTTP(c)
	}

	//perform content negotiation...
	var encoder Encoder
	var ct string
	var err error
	if c.model != nil {
		encoder, ct, err = c.GetEncoder()
		if err != nil {
			msg := fmt.Sprintf("Unable to encode to requested acceptable formats: %q", req.Header.Get("Accept"))
			http.Error(resp, msg, http.StatusNotAcceptable)
			return
		}
		resp.Header().Set("Content-Type", ct)
	}
	// if c.Reader != nil {
	// 	defer c.Reader.Close()
	// 	_, err = io.Copy(c.Writer, c.Reader)
	// 	if err != nil {
	// 		log.Printf("unable to copy stream to context writer", err)
	// 		msg := fmt.Sprintf("Unable to stream data")
	// 		r.sendError(w, msg, http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	if c.status != 0 && c.status != 200 {
		resp.WriteHeader(c.status)
	}

	if encoder != nil {
		if err := encoder.Encode(c.model); err != nil {
			panic(fmt.Sprintf("unable to encode model: %v", err))
		}
	} else {
		resp.Write(c.payload)
	}

	resp.readonly = true // prevent PostHooks from altering the response
	for _, h := range r.postHooks {
		h(c)
	}
}

// AddPreHook adds a ContextHandlerFunc that will be called before any
// handler function is called.
// They can be used to sanitize requests, authenticate users, adding
// CORS handling etc. and can respond directly, preventing any handler
// from executing if required.
// Note: PreHooks are executed in the order they are added.
func (r *Router) AddPreHook(hook ContextHandlerFunc) {
	r.preHooks = append(r.preHooks, hook)
}

// AddPostHook adds a ContextHandlerFunc that will be called after a
// handler function has been called.
// PostHooks can be used to perform cleanup tasks etc., but unlike
// PreHooks, they cannot alter a response.
// Note: PostHooks are executed in the order they are added.
func (r *Router) AddPostHook(hook ContextHandlerFunc) {
	r.postHooks = append(r.postHooks, hook)
}

// Registerer is the interface that handler function modules need to
// implement.
type Registerer interface {
	Register(r *Router)
}

// RegisterModules registers the route handler functions in each of
// the modules.
// If a specific pattern-method handlerFunc already exists, RegisterModules panics.
func (r *Router) RegisterModules(modules []Registerer) {
	for _, m := range modules {
		m.Register(r)
	}
}
