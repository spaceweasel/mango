package mango

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// NewResponseWriter returns an initialized instance of a ResponseWriter.
// If compMinLen is set above zero, then responses may be compressed if:
//
// a) the body is longer than compMinLen bytes,
// b) a suitable format has been requested in the accept-encoding header
//    (currently this must be gzip or deflate),
// c) headers have not already been sent using the WriteHeader method
func NewResponseWriter(w http.ResponseWriter, acceptedEncoding string, compMinLen int) *ResponseWriter {
	wr := ResponseWriter{
		rw:               w,
		status:           200,
		acceptedEncoding: acceptedEncoding,
		compMinLength:    compMinLen,
	}
	return &wr
}

// ResponseWriter implements the http.ResponseWriter interface and
// wraps the ResponseWriter provided to the ServeHTTP method. It's
// primary purpose is to collect data on the information written to provide
// more informative logging, but is used also for response compression.
type ResponseWriter struct {
	rw               http.ResponseWriter
	byteCount        int
	status           int
	readonly         bool
	responded        bool
	headersSent      bool
	compMinLength    int
	acceptedEncoding string
}

// Header returns the header map that will be sent by
// WriteHeader.
// See http.ResponseWriter interface for more information.
func (r *ResponseWriter) Header() http.Header {
	if r.headersSent || r.readonly {
		// return a copy of the map
		h := http.Header{}
		origMap := map[string][]string(r.rw.Header())
		for k, s := range origMap {
			for _, v := range s {
				h.Add(k, v)
			}
		}
		return h
	}
	return r.rw.Header()
}

// WriteHeader sends an HTTP response header with status code
// to the underlying http.ResponseWriter. Status code is
// recorded to provide more informative logging.
// See http.ResponseWriter interface for more information.
func (r *ResponseWriter) WriteHeader(status int) {
	if r.headersSent || r.readonly {
		return
	}
	r.headersSent = true
	r.responded = true
	r.status = status
	r.rw.WriteHeader(status)
}

func (r *ResponseWriter) compressor(w io.Writer, l int) io.WriteCloser {
	// If headers sent then we're too late for compression.
	if r.headersSent {
		return nil
	}
	if r.compMinLength == 0 || l < r.compMinLength {
		return nil
	}
	e := strings.Split(r.acceptedEncoding, ",")
	for _, ae := range e {
		switch strings.TrimSpace(ae) {
		case "gzip":
			c := gzip.NewWriter(w)
			r.rw.Header().Set("Content-Encoding", "gzip")
			return c
		case "deflate":
			c, err := flate.NewWriter(w, flate.DefaultCompression)
			if err != nil {
				return nil
			}
			r.rw.Header().Set("Content-Encoding", "deflate")
			return c
		}
	}
	return nil
}

// Write writes the data to the underlying http.ResponseWriter
// connection as part of an HTTP reply. The cumulative number
// of bytes is recorded to provide more informative logging.
// See http.ResponseWriter interface for more information.
func (r *ResponseWriter) Write(b []byte) (int, error) {
	if r.readonly {
		return 0, fmt.Errorf("write method has been called already")
	}

	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		c := r.compressor(writer, len(b))
		if c != nil {
			defer c.Close()
			c.Write(b)
		} else {
			writer.Write(b)
		}
	}()

	// TODO: check the error before updating anything
	bc, err := io.Copy(r.rw, reader)
	i := int(bc)

	r.byteCount += i
	r.headersSent = true
	r.responded = true
	return i, err
}
