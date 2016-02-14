package mango

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

type routes interface {
	AddHandlerFunc(pattern, method string, handlerFunc ContextHandlerFunc)
	HandlerFuncs(path string) (map[string]ContextHandlerFunc, map[string]string, bool)
}

// HookFunc is the signature for implementing prehook and posthook
// functions.
type HookFunc func(*Context) error

// Router is the main mango object. Router implements the standard
// library http.Handler interface, so it can be used in the call to
// http.ListenAndServe method.
// New Router objects should be created using the NewRouter method
// rather than creating a new uninitialised instance.
// TODO: Add more info here.
type Router struct {
	routes        routes
	preHooks      []HookFunc
	postHooks     []HookFunc
	encoderEngine EncoderEngine
}

// NewRouter returns a pointer to a new Router instance.
// The Router will be initialised with a new EncoderEngine
// and route handling functionality.
func NewRouter() *Router {
	r := Router{}
	r.routes = newTree()
	r.encoderEngine = newEncoderEngine()
	return &r
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

// Del registers a new handlerFunc that will be called when HTTP DELETE
// requests are made to URLs with paths that match pattern.
// If a DELETE handlerFunc already exists for pattern, Del panics.
func (r *Router) Del(pattern string, handlerFunc ContextHandlerFunc) {
	r.routes.AddHandlerFunc(pattern, "DELETE", handlerFunc)
}

// ServeHTTP dispatches the request to the handler whose pattern
// matches the request URL.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//log.Printf("INCOMING REQUEST - %v: %v\n", req.Method, req.URL.String())

	handlerFuncs, params, ok := r.routes.HandlerFuncs(req.URL.Path)
	if !ok {
		http.NotFound(w, req)
		return
	}
	fn, ok := handlerFuncs[req.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	c := &Context{
		Request:       req,
		Writer:        w,
		RouteParams:   params,
		encoderEngine: r.encoderEngine}

	//call prehooks
	var err error
	for _, h := range r.preHooks {
		if err := h(c); err != nil {
			// just log and continue for the moment - might need to
			// revisit this in the future. Possibly return from here?
			name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
			log.Printf("error running prehook (%s): %v", name, err)
		}
	}

	fn.ServeHTTP(c)

	//perform content negotiation...
	var encoder Encoder
	var ct string

	if c.model != nil {
		encoder, ct, err = c.GetEncoder()
		if err != nil {
			//log.Printf("unable to get encoder: %v", err)
			msg := fmt.Sprintf("Unable to encode to requested acceptable formats: %q", req.Header.Get("Accept"))
			r.sendError(w, msg, http.StatusNotAcceptable)
			return
		}
		w.Header().Set("Content-Type", ct)
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

	// now call posthooks here
	for _, h := range r.postHooks {
		if err := h(c); err != nil {
			// just log and continue for the moment - might need to
			// revisit this in the future. Possibly return from here?
			name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
			log.Printf("error running posthook (%s): %v", name, err)
		}
	}

	if c.status != 0 && c.status != 200 {
		w.WriteHeader(c.status)
	}

	if encoder != nil {
		if err := encoder.Encode(c.model); err != nil {
			//log.Printf("Unable to encode model: %v", err)
			msg := "Sorry, something went wrong."
			r.sendError(w, msg, http.StatusInternalServerError)
			return
		}

	} else {
		w.Write(c.payload)
	}
}

func (r *Router) sendError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprint(w, error)
}

// AddPreHook adds a HookFunc that will be called before any handler
// function is called.
// They can be used to sanitize requests, authenticate users, adding
// CORS handling, request logging etc. and can respond directly,
// preventing any handler from executing if required.
// Note: PreHooks are executed in the order they are added.
func (r *Router) AddPreHook(hook HookFunc) {
	r.preHooks = append(r.preHooks, hook)
}

// AddPostHook adds a HookFunc that will be called after a handler
// function has been called.
// They can be used to perform cleanup tasks, logging etc.
// Note: PostHooks are executed in the order they are added.
func (r *Router) AddPostHook(hook HookFunc) {
	r.postHooks = append(r.postHooks, hook)
}
