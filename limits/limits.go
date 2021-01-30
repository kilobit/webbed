/* Copyright 2020 Kilobit Labs Inc. */

package limits

import _ "fmt"
import _ "errors"
import "sync"
import "time"
import "net/http"

var defaultMaxBytes int64 = 1 << 20 // 1MB

// Examine a http.Request for conformity with a limit.
// Returns a boolean to indicate if processing should continue.
//
// A failed limit should send an appropriate response and return a
// false result.
//
// Important!! A passed limit must never invoke
// http.ResponseWriter.WriteHeader.
//
type Limit interface {
	Apply(w http.ResponseWriter, req *http.Request) bool
}

// A function type implementing the Limit interface.
//
type LimitFunc func(w http.ResponseWriter, req *http.Request) bool

func (f LimitFunc) Apply(w http.ResponseWriter, req *http.Request) bool {
	return f(w, req)
}

// Set particular limits on the inbound requests.
//
// Currently supports the MaxBytes for the req.Body.
// HTTP Timeouts are set on the Server and Listener structs.
//
// TODO: Implement per IP rate limiting.
//
type HTTPLimitsHandler struct {
	maxBytes int64
	limits   []Limit
	handler  http.Handler
}

func (lh HTTPLimitsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	for _, limit := range lh.limits {
		result := limit.Apply(w, req)
		if !result {
			return
		}
	}

	req.Body = http.MaxBytesReader(w, req.Body, lh.maxBytes)
	lh.handler.ServeHTTP(w, req)
}

func NewHTTPLimitsHandler(handler http.Handler, limits ...Limit) *HTTPLimitsHandler {

	return &HTTPLimitsHandler{
		defaultMaxBytes,
		limits,
		handler,
	}
}

// Limit the total number of requests allowed over a given period.
//
// Will wait before returning the error.
//
func RateLimit(limit int, period, wait time.Duration) Limit {

	mu := sync.Mutex{}
	start := time.Now()
	current := 0

	return LimitFunc(func(w http.ResponseWriter, req *http.Request) bool {

		exceeded := false
		now := time.Now()

		// Begin mutual exclusion
		mu.Lock()

		if now.Sub(start) > period {
			current = 0
		}

		current++

		exceeded = current > limit

		// End mutual exclusion
		mu.Unlock()

		if exceeded {
			http.Error(w, "Request threshold has been exceeded.", http.StatusTooManyRequests)
			return false
		}

		return true
	})

}
