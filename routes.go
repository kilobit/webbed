/* Copyright 2020 Kilobit Labs Inc. */

// Accept HTTP Form Data and forward it as an email.
package informed

import "fmt"
import _ "errors"
import _ "net/url"
import "net/http"

// Route inbound requests to the appropriate handler.
type HTTPRouteHandler struct {
	handler http.Handler // Handler
}

func (fh HTTPRouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ok")
}

func NewHTTPRouteHandler(handler http.Handler) *HTTPRouteHandler {
	return &HTTPRouteHandler{handler}
}
