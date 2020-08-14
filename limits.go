/* Copyright 2020 Kilobit Labs Inc. */

package informed

import _ "fmt"
import _ "errors"
import "net/http"

var defaultMaxBytes int64 = 1 << 20 // 1MB

// Set particular limits on the inbound requests.
//
// Currently supports the MaxBytes for the req.Body.
// HTTP Timeouts are set on the Server and Listener structs.
//
// TODO: Implement per IP rate limiting.
//
type HTTPLimitsHandler struct {
	maxBytes int64
	handler  http.Handler
}

func (lh HTTPLimitsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	req.Body = http.MaxBytesReader(w, req.Body, lh.maxBytes)
	lh.handler.ServeHTTP(w, req)
}

func NewHTTPLimitsHandler(handler http.Handler) *HTTPLimitsHandler {

	return &HTTPLimitsHandler{
		defaultMaxBytes,
		handler,
	}
}
