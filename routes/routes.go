/* Copyright 2021 Kilobit Labs Inc. */

// Inspired by Ben Hoyt and Axel Wagner's ShiftPath.  This package
// encapsulates the ShiftPath technique, avoiding regexes but allowing
// route definition in one place.
//
package routes

import "fmt"
import _ "errors"

import "strings"
import _ "net/url"
import "net/http"

import "kilobit.ca/go/webbed"

type Route struct {
	nfh http.Handler
	rs map[string]http.Handler
}

func New(notFoundHandler http.Handler) *Route {
	return &Route{
		notFoundHandler,
		map[string]http.Handler{},
	}
}

func (route Route) Add(path string, handler http.Handler) error {

// 	u, err := url.Parse(path)
// 	if err != nil {
// 		return err
// 	}

	head, rest := webbed.ShiftPath(path)
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
			//fmt.Printf("Swap: '%s', '%s'\n%s", head, rest, subroute)
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
	head, r.URL.Path = webbed.ShiftPath(r.URL.Path)

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
			fmt.Fprint(sb, r.ToString(prefix + " -> " + seg))
		}
	}

	return sb.String()
}

func (route Route) String() string {

	return route.ToString("")
}
