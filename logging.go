package mango

import (
	"fmt"
	"net/http"
	"time"
)

// NewWatchedResponse returns an initialized instance of a WatchedResponse.
func NewWatchedResponse(w http.ResponseWriter) *WatchedResponse {
	wr := WatchedResponse{
		rw:     w,
		status: 200,
	}
	return &wr
}

// WatchedResponse implements the http.ResponseWriter interface and
// wraps the ResponseWriter provided to the ServeHTTP method. It's
// purpose is to collect data on the information written to provide
// more informative logging.
type WatchedResponse struct {
	rw        http.ResponseWriter
	byteCount int
	status    int
	readonly  bool
	responded bool
}

// Header returns the header map that will be sent by
// WriteHeader.
// See http.ResponseWriter interface for more information.
func (w *WatchedResponse) Header() http.Header {
	if w.readonly {
		// return a copy of the map
		h := http.Header{}
		origMap := map[string][]string(w.rw.Header())
		for k, s := range origMap {
			for _, v := range s {
				h.Add(k, v)
			}
		}
		return h
	}
	return w.rw.Header()
}

// WriteHeader sends an HTTP response header with status code
// to the underlying http.ResponseWriter. Status code is
// recorded to provide more informative logging.
// See http.ResponseWriter interface for more information.
func (w *WatchedResponse) WriteHeader(status int) {
	if w.readonly {
		return
	}
	w.responded = true
	w.status = status
	w.rw.WriteHeader(status)
}

// Write writes the data to the underlying http.ResponseWriter
// connection as part of an HTTP reply. The cumulative number
// of bytes is recorded to provide more informative logging.
// See http.ResponseWriter interface for more information.
func (w *WatchedResponse) Write(b []byte) (int, error) {
	if w.readonly {
		return 0, fmt.Errorf("write method has been called already")
	}
	w.responded = true
	i, err := w.rw.Write(b)
	w.byteCount += i
	return i, err
}

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
	}
	return &log
}

// RequestLog is the structure used to record data about a request
// (and response). In addition to information extracted from the
// request object, this struct holds details about the duration,
// status and amount of data returned.
type RequestLog struct {
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

func (r *RequestLog) stop() {
	r.Finish = nowUTC()
	r.Duration = r.Finish.Sub(r.Start)
}
