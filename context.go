package mango

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

// Response is an object used to facilitate building a response.
type Response struct {
	context *Context
	model   interface{}
	status  int
}

// WithModel sets the Model that will be serialized for the response.
// The serialization mechanism will depend on the request Accept header,
// the encoding DefaultMediaType and whether the WithContentType method
// has been used.
// If the Accept request header is missing, empty or equal to "*/*", then
// the model will be encoded with the default media type. If required, the
// default media type can be overridden for the individual response using
// WithContentType.
// This method returns the Response object and can be chained.
func (r *Response) WithModel(m interface{}) *Response {
	r.context.model = m
	r.context.responseReady = true
	return r
}

// WithStatus sets the HTTP status code of the response.
// This method returns the Response object and can be chained.
func (r *Response) WithStatus(s int) *Response {
	r.context.status = s
	r.context.responseReady = true
	return r
}

// WithHeader adds a header to the response.
// This method returns the Response object and can be chained.
func (r *Response) WithHeader(key, value string) *Response {
	r.context.Writer.Header().Add(key, value)
	return r
}

// WithContentType sets the Content-Type header of the response.
// WithContentType overrides the default media type for this individual
// response. If the response contains a model and the Accept request header
// is missing, empty or equal to "*/*", then the model will be encoded with
// type ct.
// This method returns the Response object and can be chained.
func (r *Response) WithContentType(ct string) *Response {
	r.context.Writer.Header().Set("Content-Type", ct)
	return r
}

// Context is the request context.
// Context encapsulates the underlying Req and Writer, but exposes
// them if required. It provides many helper methods which are
// designed to keep your handler code clean and free from boiler
// code.
type Context struct {
	Request        *http.Request
	Writer         http.ResponseWriter
	status         int
	payload        []byte
	model          interface{}
	RouteParams    map[string]string
	encoderEngine  EncoderEngine
	Reader         io.ReadCloser
	Identity       Identity
	responseReady  bool
	modelValidator ModelValidator
	X              interface{}
}

// ContextHandlerFunc type is an adapter to allow the use of ordinary
// functions as HTTP handlers. It is similar to the standard library's
// http.HandlerFunc in that if f is a function with the appropriate
// signature, ContextHandlerFunc(f) is a Handler object that calls f.
type ContextHandlerFunc func(*Context)

// ServeHTTP calls f(c).
func (f ContextHandlerFunc) ServeHTTP(c *Context) {
	f(c)
}

// Respond returns a new context based Response object.
func (c *Context) Respond() *Response {
	return &Response{context: c}
}

// RespondWith is a generic method for producing a simple response.
// It takes a single parameter whose type will determine the action.
//
// Strings will be used for the response content.
// Integers will be used for the response status code.
// Any other type is deemed to be a model which will be serialized.
//
// The serialization mechanism for the model will depend on the
// request Accept header, the encoding DefaultMediaType and whether
// the WithContentType method has been used. See the Response struct
// for more details.
// This method returns the Response object and can be chained.
func (c *Context) RespondWith(d interface{}) *Response {
	response := &Response{context: c}

	switch t := d.(type) {
	case int:
		c.status = t
	case string:
		c.payload = []byte(t)
	default: //must be a model
		c.model = d
	}
	c.responseReady = true
	return response
}

// Authenticated returns true if a request user has been authenticated.
// Authentication should be performed in a pre-hook, assigning a valid
// Identity to the Context if authentication succeeds.
// This method simply examines whether the Context has a valid Identity.
func (c *Context) Authenticated() bool {
	return c.Identity != nil
}

func (c *Context) urlSchemeHost() string {
	if c.Request.TLS != nil {
		return "https://" + c.Request.Host
	}
	return "http://" + c.Request.Host
}

// Error sends the specified message and HTTP status code as a response.
// Request handlers should cease execution after calling this method.
func (c *Context) Error(msg string, code int) {
	http.Error(c.Writer, msg, code)
}

// Redirect sends a redirect response using the specified URL and HTTP
// status.
// Request handlers should cease execution after calling this method.
// TODO: Not yet implemented
func (c *Context) Redirect(urlStr string, code int) {
	http.Redirect(c.Writer, c.Request, urlStr, code)
}

// // Render executes a template using the supplied data.
// // Request handlers should cease execution after calling this method.
// // TODO: Not yet implemented
// func (c *Context) Render(tmpl string, data interface{}) {
// 	panic("not yet implemented")
// }

//
// func (c *Context) sendResponse() {
// 	fmt.Fprintf(c.Writer, "")
// }

func (c *Context) contentDecoder(r io.Reader) (Decoder, error) {
	ct := c.Request.Header.Get("Content-Type")
	ct = strings.Replace(ct, " ", "", -1)
	decoder, err := c.encoderEngine.GetDecoder(r, ct)
	if err != nil {
		// If the full Content-Type doesn't match try matching only up to the ;
		decoder, err = c.encoderEngine.GetDecoder(r, strings.Split(ct, ";")[0])
	}
	if err != nil {
		return nil,
			UnsupportedMediaTypeError{
				hdr: "Content-Type",
				val: ct,
			}
	}
	return decoder, nil
}

func (c *Context) acceptableMediaTypes() []string {
	hdr := c.Request.Header.Get("Accept")
	hdr = strings.Replace(hdr, " ", "", -1)
	types := strings.Split(hdr, ",")
	mt := make(mediaTypes, len(types))

	for i, t := range types {
		m, err := newMediaType(t)
		if err != nil {
			continue
		}
		mt[i] = *m
	}
	sort.Sort(mt)
	r := []string{}
	for _, t := range mt {
		if !t.Empty() {
			r = append(r, t.String())
		}
	}
	return r
}

// GetEncoder returns an Encoder suitable for serializing data in a response.
// The Encoder is selected based on the request Accept header (or default media
// type if no Accept header supplied).
// If successful, the an encoder and content-type are returned and a nil error.
// Success is determined by a nil error.
// The returned encoder will have been pre-injected with an io.Writer, so the
// Encode method can be called directly, passing the data to be encoded as the
// only parameter.
func (c *Context) GetEncoder() (Encoder, string, error) {
	mts := c.acceptableMediaTypes()
	var err error
	var mt string
	for _, mt = range mts {
		if mt == "*/*" {
			// use specified encoding if specified in response header
			mt = c.Writer.Header().Get("Content-Type")
			if mt == "" {
				mt = c.encoderEngine.DefaultMediaType()
			}
		}
		var encoder Encoder
		encoder, err = c.encoderEngine.GetEncoder(c.Writer, mt)
		if err == nil {
			return encoder, mt, nil
		}
	}
	return nil, mt, err
}

// Bind populates the supplied model with data from the request.
// This is performed in stages. initially, any requestbody content is
// deserialized.
//
// TODO: Following is not yet implemented:
//
// Route parameters are used next to populate any unset members.
// Finally, query parameters are used to populate any remaining unset members.
//
// This method is under review - currently Binding only uses deserialized
// request body content.
func (c *Context) Bind(m interface{}) error {
	r := c.Request.Body
	ce, err := c.contentEncoding()
	if err != nil {
		return err
	}
	if ce == "gzip" {
		r, err = gzip.NewReader(c.Request.Body)
		if err != nil {
			return err
		}
		defer r.Close()
	}

	decoder, err := c.contentDecoder(r)
	if err != nil {
		return err
	}
	err = decoder.Decode(m)
	if err != nil {
		return err
	}

	// TODO: now update any missing empty properties from url path/query params

	return nil
}

func (c *Context) contentEncoding() (string, error) {
	ce := c.Request.Header.Get("Content-Encoding")
	if ce == "" {
		// Some proxies (e.g. Zuul) strip the 'content-encoding'
		// header, so check for custom style variety, 'X-...'
		ce = c.Request.Header.Get("X-Content-Encoding")
	}
	ce = strings.ToLower(ce)
	acceptable := []string{"", "gzip"}
	if !stringInSlice(ce, acceptable) {
		return "", UnsupportedMediaTypeError{
			hdr: "Content-Encoding",
			val: ce,
		}
	}
	return ce, nil
}

// Validate validates the properties of the model m.
func (c *Context) Validate(m interface{}) (map[string][]ValidationFailure, bool) {
	return c.modelValidator.Validate(m)
}

// UnsupportedMediaTypeError means that the request payload was not in an
// acceptable format. This error occurs if the header value for either
// 'content-type' or 'content-encoding' is not recognised.
type UnsupportedMediaTypeError struct {
	hdr string
	val string
}

func (e UnsupportedMediaTypeError) Error() string {
	return fmt.Sprintf("unsupported media type (%s: %s)", e.hdr, e.val)
}
