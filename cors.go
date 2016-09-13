package mango

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig holds CORS configuration. It can be used as
// the configuration for an individual resource or as a
// global configuration for the entire router tree.
type CORSConfig struct {
	// Origins list all permitted origins. A CORS request
	// origin MUST be in the Origins list for the response
	// headers to be populated with the correct response.
	// Values must contain the scheme (e.g. http://here.com).
	// A wildcard * can be used which will match ALL origins,
	// however the Access-Control-Allow-Origin response header
	// will always echo the request Origin if the remaining CORS
	// criteria is met.
	Origins []string
	// Methods available for the resource. If there are methods
	// listed here for which there is no handler, then that
	// method will not be included in the Access-Control-Allow-Methods
	// response header.
	Methods []string
	// Headers lists the custom headers in a request that the
	// server will accept
	Headers []string
	// ExposedHeaders are custom headers that the client browser
	// is allowed to access
	ExposedHeaders []string
	// Credentials controls the Access-Control-Allow-Credentials header.
	// The header is only included in the response if Credentials is true,
	// in which case the header has a value of "true".
	// A value of true allows the client browser to access response cookies.
	Credentials bool
	// MaxAge is the cache duration (in seconds) that is returned
	// in a Preflight Access-Control-Max-Age response header.
	// A value of zero means the header won't be sent.
	MaxAge int
}

func (c *CORSConfig) originAllowed(origin string) bool {
	return stringInSlice("*", c.Origins) ||
		stringInSlice(origin, c.Origins)
}

func (c *CORSConfig) methodAllowed(method string) bool {
	return stringInSlice(method, c.Methods)
}

func (c *CORSConfig) headersAllowed(headers string) bool {
	ha := strings.Split(headers, ",")
	for i := 0; i < len(ha); i++ {
		h := strings.TrimSpace(ha[i])
		if h == "" {
			continue
		}
		if !stringInSlice(h, c.Headers) {
			return false
		}
	}
	return true
}

func (c *CORSConfig) clone() *CORSConfig {
	cl := CORSConfig{
		Origins:        c.Origins,
		Methods:        c.Methods,
		Headers:        c.Headers,
		ExposedHeaders: c.ExposedHeaders,
		Credentials:    c.Credentials,
		MaxAge:         c.MaxAge,
	}
	return &cl
}

func (c *CORSConfig) merge(m CORSConfig) {
	c.Origins = appendIfNotExists(c.Origins, m.Origins)
	c.Methods = appendIfNotExists(c.Methods, m.Methods)
	c.Headers = appendIfNotExists(c.Headers, m.Headers)
	c.ExposedHeaders = appendIfNotExists(c.ExposedHeaders, m.ExposedHeaders)
	c.Credentials = m.Credentials
	c.MaxAge = m.MaxAge
}

func appendIfNotExists(dest []string, src []string) []string {
	for _, v := range src {
		if !stringInSlice(v, dest) {
			dest = append(dest, v)
		}
	}
	return dest
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// type CORSType int
//
// const (
// 	NoCORS CORSType = 0 + iota
// 	SimpleCORS
// 	PreflightCORS
// )

func handleCORS(req *http.Request, w http.ResponseWriter, resource *Resource) (preflight bool) {
	origin := req.Header.Get("Origin")
	corsConf := (*resource).CORSConfig
	if corsConf == nil {
		return
	}
	if !(*corsConf).originAllowed(origin) {
		return
	}

	if req.Method == "OPTIONS" {
		// check for preflight
		method := req.Header.Get("Access-Control-Request-Method")
		if method == "" {
			return
		}
		if !corsConf.methodAllowed(method) {
			return
		}
		if _, ok := resource.Handlers[method]; !ok {
			return
		}
		reqHeaders := req.Header.Get("Access-Control-Request-Headers")
		if !corsConf.headersAllowed(reqHeaders) {
			return
		}

		for _, m := range (*corsConf).Methods {
			if _, ok := resource.Handlers[m]; !ok {
				continue
			}
			w.Header().Add("Access-Control-Allow-Methods", m)
		}

		for _, h := range (*corsConf).Headers {
			w.Header().Add("Access-Control-Allow-Headers", h)
		}
		if (*corsConf).MaxAge > 0 {
			maStr := strconv.Itoa((*corsConf).MaxAge)
			w.Header().Set("Access-Control-Max-Age", maStr)
		}
		preflight = true
	} else {
		// normal request
		for _, h := range (*corsConf).ExposedHeaders {
			w.Header().Add("Access-Control-Expose-Headers", h)
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Add("Vary", "Origin")
	if (*corsConf).Credentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	return
}
