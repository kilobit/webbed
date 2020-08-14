/* Copyright 2020 Kilobit Labs Inc. */

// Accept HTTP Form Data and forward it as an email.
package informed

import _ "fmt"
import _ "errors"
import "net/http"

// Set particular limits on the inbound requests.
//
// Currently supports the MaxBytes for the req.Body.
// HTTP Timeouts are set on the Server and Listener structs.
//
// TODO: Implement per IP rate limiting.
//
type HTTPLimitsHandler struct {
	MaxBytes int64
	handler http.Handler
}

func (lh HTTPLimitsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	req.Body = http.MaxBytesReader(w, req.Body, lh.MaxBytes)
	lh.handler.ServeHTTP(w, req)
}

func NewHTTPLimitsHandler(handler http.Handler, maxBytes int64) *HTTPLimitsHandler {

	return &HTTPLimitsHandler{
		maxBytes,
		handler,
	}
}
