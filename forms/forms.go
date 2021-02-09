/* Copyright 2020 Kilobit Labs Inc. */

// Simple HTTP form handling for Golang.
//
package forms

import "fmt"
import _ "errors"
import "context"
import "net/url"
import "net/http"

import "kilobit.ca/go/webbed"

// Do something with the form data.
//
// Return a http response code and error.
type FormHandler interface {
	Handle(context.Context, url.Values) (int, error)
}

// Do something with the form data. Implements the FormHandler
// interface.
//
// Return a http response code and error.
type FormHandlerFunc func(context.Context, url.Values) (int, error)

func (f FormHandlerFunc) Handle(ctx context.Context, values url.Values) (int, error) {
	return f(ctx, values)
}

// A http.Handler that reads in form data from an http.Request and
// passes it to a FormHandler.
type HTTPFormHandler struct {
	handler FormHandler // Handler
}

func (fh HTTPFormHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		webbed.ServeError(w, http.StatusBadRequest, fmt.Errorf("Error while parsing the form data, %w", err))
		return
	}

	code, err := fh.handler.Handle(req.Context(), req.Form)
	if err != nil {
		webbed.ServeError(w, code, fmt.Errorf("Error while processing the form data, %w", err))
		return
	}

	if code != 0 {
		w.WriteHeader(code)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	fmt.Fprintf(w, "Ok")
}

func New(handler FormHandler) *HTTPFormHandler {
	return &HTTPFormHandler{handler}
}
