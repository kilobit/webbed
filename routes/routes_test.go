/* Copyright 2021 Kilobit Labs Inc. */

package routes_test

import _ "fmt"
import _ "errors"
import "io/ioutil"
import "net/http"

import "kilobit.ca/go/webbed/routes"

import "testing"
import "net/http/httptest"
import "kilobit.ca/go/tested/assert"

func TestRoutesTest(t *testing.T) {
	assert.Expect(t, true, true)
}

func NewFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.EscapedPath()))
	})
}

func NewNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(r.URL.EscapedPath()))
	})
}

func TestNewRoute(t *testing.T) {

	routes.New(NewNotFoundHandler())
}

func TestRoutesHTTPHandler(t *testing.T) {

	tests := []struct {
		routes []string
		path   string
		status int
		body   string
	}{
		{[]string{"/"}, "/", 200, ""},
		{[]string{"/"}, "/bar", 404, ""},

		{[]string{"/foo"}, "/", 404, ""},
		{[]string{"/foo"}, "/foo", 200, ""},
		{[]string{"/foo"}, "/bar", 404, ""},

		{[]string{"/foo/bar"}, "/foo/bar", 200, ""},
		{[]string{"/foo/bar"}, "/foo", 404, ""},

		{[]string{"/foo", "/foo/bar"}, "/foo", 200, ""},
		{[]string{"/foo", "/foo/bar"}, "/foo/bar", 200, ""},

		{[]string{"/foo", "/foo/bar/bing"}, "/foo", 200, ""},
		{[]string{"/foo", "/foo/bar/bing"}, "/foo/bar", 404, ""},
		{[]string{"/foo", "/foo/bar/bing"}, "/foo/bar/bing", 200, ""},

		{[]string{"/foo/bar", "/foo/bar/bing"}, "/foo/bar", 200, ""},

		{[]string{"/foo/bar"}, "/foo/bar", 200, ""},
		{[]string{"/foo/bar", "/"}, "/foo/bar", 200, ""},
		{[]string{"/foo/bar", "/", "/foo/bar/bing"}, "/foo/bar", 200, ""},
		{[]string{"/foo/bar", "/", "/foo/bar/bing", "/foo/bar/bong"}, "/foo/bar", 200, ""},
	}

	for _, test := range tests {

		//t.Log(test)

		foundHandler := NewFoundHandler()

		notFoundHandler := NewNotFoundHandler()

		r := routes.New(notFoundHandler)
		for _, route := range test.routes {
			r.Add(route, foundHandler)
			//t.Logf("\n%#v\n", r)
			//t.Log("-----")
		}

		//t.Log(r.String())

		rw := httptest.NewRecorder()

		req := httptest.NewRequest("GET", test.path, nil)
		//t.Log(req)

		r.ServeHTTP(rw, req)

		res := rw.Result()

		assert.Expect(t, test.status, res.StatusCode, test, res)

		body, err := ioutil.ReadAll(res.Body)
		assert.Ok(t, err, test, req)
		assert.Expect(t, test.body, string(body), test, req)
	}
}

func TestRouteString(t *testing.T) {

	r := routes.New(NewNotFoundHandler())
	r.Add("/foo/bar", NewFoundHandler())
	r.Add("/", NewFoundHandler())
	r.Add("/foo/bar/bing", NewFoundHandler())
	r.Add("/foo/bar/bong", NewFoundHandler())

	//t.Log("\n", r)

	//fmt.Print(r)

	//assert.Expect(t, "", r.String())
}

func TestShiftPath(t *testing.T) {

	tests := []struct {
		path string
		exp  []string
	}{
		{"", []string{""}},
		{"/", []string{""}},
		{"/foo", []string{"foo", ""}},
		{"/foo/bar", []string{"foo", "bar", ""}},
		{"/foo/bar/bing", []string{"foo", "bar", "bing", ""}},
		{"/foo/bar/bing/", []string{"foo", "bar", "bing", ""}},
	}

	for _, test := range tests {

		result := []string{}
		rest := test.path

		for {
			var path string
			path, rest = routes.ShiftPath(rest)

			//t.Log(path, rest)

			result = append(result, path)
			if path == "" {
				break
			}
		}

		//t.Log(result)
		assert.ExpectDeep(t, test.exp, result)
	}
}
