/* Copyright 2021 Kilobit Labs Inc. */

// This helper package filters requests with an unsupported method
// set.
//
package methods

import "fmt"
import _ "errors"

import "strings"
import "net/http"

type MethodHandler struct {
	errorHandler http.Handler
}

var defaultErrorHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	msg := fmt.Sprintf("The '%s' is not supported on this URL.", req.Method)

	http.Error(w, msg, http.StatusMethodNotAllowed)
})

// Create a new MethodHandler.
//
// Uses the defaultErrorHandler if errorHandler is nil.
//
func New(errorHandler http.Handler) *MethodHandler {

	if errorHandler == nil {
		errorHandler = defaultErrorHandler
	}

	return &MethodHandler{errorHandler}
}

// Return a http.Handler that runs the given handler if the
// Request.Method is in the given list.  Otherwise, the error handler
// is invoked.
//
// TODO: Consider pre-processing ToUpper on given methods.
//
func (mh *MethodHandler) HandleMethods(handler http.Handler, methods ...string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for _, method := range methods {
			if strings.ToUpper(req.Method) == strings.ToUpper(method) {
				handler.ServeHTTP(w, req)
				return
			}
		}

		mh.errorHandler.ServeHTTP(w, req)
	})
}

func (mh *MethodHandler) Get(handler http.Handler) http.Handler {
	return mh.HandleMethods(handler, http.MethodGet)
}

func (mh *MethodHandler) Head(handler http.Handler) http.Handler {
	return mh.HandleMethods(handler, http.MethodHead)
}

func (mh *MethodHandler) Post(handler http.Handler) http.Handler {
	return mh.HandleMethods(handler, http.MethodPost)
}

func (mh *MethodHandler) Put(handler http.Handler) http.Handler {
	return mh.HandleMethods(handler, http.MethodPut)
}

func (mh *MethodHandler) Patch(handler http.Handler) http.Handler {
	return mh.HandleMethods(handler, http.MethodPatch)
}

func (mh *MethodHandler) Delete(handler http.Handler) http.Handler {
	return mh.HandleMethods(handler, http.MethodDelete)
}

func (mh *MethodHandler) Connect(handler http.Handler) http.Handler {
	return mh.HandleMethods(handler, http.MethodConnect)
}

func (mh *MethodHandler) Options(handler http.Handler) http.Handler {
	return mh.HandleMethods(handler, http.MethodOptions)
}

func (mh *MethodHandler) Trace(handler http.Handler) http.Handler {
	return mh.HandleMethods(handler, http.MethodTrace)
}
