/* Copyright 2020 Kilobit Labs Inc. */

package informed

import "fmt"
import _ "errors"
import "net/url"
import "net/http"

// Do something with the form data.
//
// Return a http response code and error.
type FormHandler interface {
	Handle(url.Values) (int, error)
}

// Do something with the form data. Implements the FormHandler
// interface.
//
// Return a http response code and error.
type FormHandlerFunc func(url.Values) (int, error)

func (f FormHandlerFunc) Handle(values url.Values) (int, error) {
	return f(values)
}

// A http.Handler that reads in form data from an http.Request and
// passes it to a FormHandler.
type HTTPFormHandler struct {
	handler FormHandler // Handler
}

func (fh HTTPFormHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		ServeError(w, http.StatusBadRequest, fmt.Errorf("Error while parsing the form data, %w", err))
		return
	}

	code, err := fh.handler.Handle(r.Form)
	if err != nil {
		ServeError(w, code, fmt.Errorf("Error while processing the form data, %w", err))
		return
	}

	if code != 0 {
		w.WriteHeader(code)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	fmt.Fprintf(w, "Ok")
}

func NewHTTPFormHandler(handler FormHandler) *HTTPFormHandler {
	return &HTTPFormHandler{handler}
}
