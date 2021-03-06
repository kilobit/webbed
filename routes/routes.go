/* Copyright 2021 Kilobit Labs Inc. */

// Simple Web Routing in Golang.
//
// Inspired by Ben Hoyt and Axel Wagner's ShiftPath.  This package
// encapsulates the ShiftPath technique, avoiding regexes but allowing
// route definition in one place.  The route map is compiled into a
// Trie for efficient lookups.
//
// TODO: Benchmark this router using Ben Hoyt's benchmarks.
//
package routes

import "fmt"
import _ "errors"

import "strings"
import "path"
import "net/url"
import "net/http"

type Route struct {
	nfh http.Handler
	rs  map[string]http.Handler
}

func New(notFoundHandler http.Handler) *Route {
	return &Route{
		notFoundHandler,
		map[string]http.Handler{},
	}
}

func (route Route) Add(path string, handler http.Handler) error {

	head, rest := ShiftPath(path)
	isleaf := rest == ""
	current, hascurrent := route.rs[head]

	if hascurrent {

		subroute, isroute := current.(*Route)

		if isroute {
			subroute.Add(rest, handler)

			// not a route
		} else {
			// swap
			subroute := New(route.nfh)
			subroute.Add("", current)
			subroute.Add(rest, handler)
			route.rs[head] = subroute
		}

		// no current
	} else {
		if isleaf {
			route.rs[head] = handler

			// Not a leaf
		} else {
			subroute := New(route.nfh)
			subroute.Add(rest, handler)
			route.rs[head] = subroute
		}
	}

	return nil
}

func (route Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

	handler, found := route.rs[head]
	if !found {
		handler = route.nfh
	}

	handler.ServeHTTP(w, r)
}

func (route Route) ToString(prefix string) string {

	sb := &strings.Builder{}

	for seg, handler := range route.rs {

		hstr := "handler"
		r, ok := handler.(*Route)
		if ok {
			hstr = "route"
		}

		fmt.Fprintf(sb, "%s -> '%s' [%s]\n", prefix, seg, hstr)

		if ok {
			fmt.Fprint(sb, r.ToString(prefix+" -> "+seg))
		}
	}

	return sb.String()
}

func (route Route) String() string {

	return route.ToString("")
}

// Invoke the handler with the path returned to it's original state
// before being routed.
//
// As the router processes path segments it removes the handled
// prefixes.  Some http.Handlers such as the http.FileServer will
// behave poorly when the path has been changed.  Passing the
// FileServer to this Handler will clear up that problem.
//
func ResetPath(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL, _ = url.Parse(req.RequestURI)

		h.ServeHTTP(w, req)
	})
}

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
