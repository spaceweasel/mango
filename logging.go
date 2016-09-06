package mango

import (
	"fmt"
	"net/http"
	"time"
)

// extract Now() to own func to facilitate testing
var nowUTC = func() time.Time {
	return time.Now().UTC()
}

// NewRequestLog returns an initialized *RequestLog populated with information
// extracted from req.
func NewRequestLog(req *http.Request) *RequestLog {
	log := RequestLog{
		Start:         nowUTC(),
		RemoteAddr:    req.RemoteAddr,
		AccessRequest: req.Method + " " + req.URL.RequestURI() + " " + req.Proto,
		Host:          req.Host,
		Referer:       req.Referer(),
		UserAgent:     req.UserAgent(),
		UserID:        "-",
		Method:        req.Method,
		URI:           req.URL.RequestURI(),
		Protocol:      req.Proto,
		header:        req.Header,
	}
	return &log
}

// RequestLog is the structure used to record data about a request
// (and response). In addition to information extracted from the
// request object, this struct holds details about the duration,
// status and amount of data returned.
type RequestLog struct {
	identity Identity
	context  interface{}
	header   http.Header

	// Start is the time the request was received
	Start time.Time

	// Finish is the time the request was completed

	Finish time.Time
	// Host is the host on which the requested resource resides.
	// Format is "host:port" although port is omitted for standard
	// ports.
	// Example: www.somedomain.com
	Host string

	// AccessRequest is a concatenation of request information, in the
	// format: METHOD Path&Query protocol
	//
	// e.g. GET /somepath/more/thing.gif HTTP/1.1
	AccessRequest string

	// Status is the response status code
	Status int

	// BytesOut is the number of bytes returned by the response
	BytesOut int

	// Duration is the time taken to process the request.
	Duration time.Duration

	// UserAgent is the client's user agent string (if provided)
	UserAgent string

	// RemoteAddr is the network address that sent the request.
	// Format is "IP:port"
	RemoteAddr string

	// Referer is the referring URL (if provided).
	// Referer is misspelled deliberately to match HTTP specification.
	Referer string

	// UserID returns the UserID of the authenticated user making the
	// request. Returns "-" if the request user has not been authenticated.
	UserID string

	// Method is the HTTP request method.
	Method string

	// URI is the path and query of the request URI.
	URI string

	// Protocol is the HTTP Protocol used for the request.
	Protocol string
}

// CommonFormat returns request data as a string in W3C Common Log Format.
// (https://en.wikipedia.org/wiki/Common_Log_Format)
func (r *RequestLog) CommonFormat() string {
	timeStamp := r.Start.Format("02/Jan/2006:15:04:05 -0700")
	s := fmt.Sprintf("%s - %s [%s] \"%s\" %d %d",
		r.RemoteAddr, r.UserID, timeStamp, r.AccessRequest, r.Status, r.BytesOut)
	return s
}

// CombinedFormat returns request data as a string in Combined Log Format.
// Combined Log Format is similar to Common Log Format, with the addition
// of the Referer and UserAgent.
func (r *RequestLog) CombinedFormat() string {
	s := fmt.Sprintf("%s \"%s\" \"%s\"",
		r.CommonFormat(), r.Referer, r.UserAgent)
	return s
}

// Context returns the user defined X property of the mango context.
func (r *RequestLog) Context() interface{} {
	return r.context
}

// Identity returns the identity of the authenticated user.
// Identity is a property of the mango context and must be set in the
// consuming application (e.g. using a custom authenticaion hook).
func (r *RequestLog) Identity() Identity {
	return r.identity
}

// Header returns the request header value for the specified key.
// Returns an empty string if no matching header is found.
func (r *RequestLog) Header(key string) string {
	return r.header.Get(key)
}

func (r *RequestLog) stop() {
	r.Finish = nowUTC()
	r.Duration = r.Finish.Sub(r.Start)
}
