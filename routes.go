/* Copyright 2020 Kilobit Labs Inc. */

// Accept HTTP Form Data and forward it as an email.
package informed

import "fmt"
import _ "errors"
import "strings"
import "path"
import "net/url"
import "net/http"

type Routes map[string]http.Handler

// Route inbound requests to the appropriate handler.
type HTTPRouteHandler struct {
	routes          Routes
	notFoundHandler http.Handler
}

func (rh HTTPRouteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	route, tail := ShiftPath(req.URL.EscapedPath())

	handler, ok := rh.routes[route]
	if !ok {
		rh.notFoundHandler.ServeHTTP(w, req)
		return
	}

	req.URL.Path = tail

	handler.ServeHTTP(w, req)
}

// Returns an empty RouteHandler with the defaultNotFoundHandler.
func NewHTTPRouteHandler() *HTTPRouteHandler {
	return &HTTPRouteHandler{Routes{}, defaultNotFoundHandler}
}

// Add or re-point a route to a handler.
//
// For nested paths, use this HTTPRouteHandler as the target.
//
func (rh *HTTPRouteHandler) SetRoute(route string, handler http.Handler) {
	rh.routes[route] = handler
}

// Returns an error if *this* RouteHandler is provided.
func (rh *HTTPRouteHandler) SetNotFoundHandler(handler http.Handler) error {

	if handler == rh {
		return fmt.Errorf("Refusing to set the RouteHandler as it's own not found handler.")
	}

	rh.notFoundHandler = handler

	return nil
}

var defaultNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	ServeError(w, http.StatusNotFound, fmt.Errorf("the route, %s was not found on this server", req.URL.RawPath))
})

// Shift a Path element from the URL.
//
func ShiftPath(p string) (head, tail string) {

	path := path.Clean(p)

	i := strings.Index(path[1:], "/")
	if i == -1 {
		i = len(path) - 1
	}

	i++

	head, err := url.QueryUnescape(path[1:i])
	if err != nil {
		head = path[1:i]
	}

	return head, path[i:]
}
