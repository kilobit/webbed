/* Copyright 2020 Kilobit Labs Inc. */

// Tests for the webbed package.
package webbed_test

import _ "fmt"
import _ "errors"
import _ "strings"
import _ "io/ioutil"
import "net/url"
import "net/http"
import "net/http/httptest"
import "testing"
import "kilobit.ca/go/webbed"
import "kilobit.ca/go/tested/assert"

func TestRoutesTest(t *testing.T) {

	assert.Expect(t, true, true)
}

var routestestdata = []struct {
	desc  string
	route string
	path  string
	tail  string
	code  int
}{
	{
		"simple",
		"one",
		"one/",
		"",
		200,
	},

	{
		"2-element",
		"one",
		"one/two",
		"/two",
		200,
	},

	{
		"no path",
		"",
		"",
		"",
		200,
	},

	{
		"not found",
		"one",
		"/",
		"",
		404,
	},
}

func TestHTTPRouteHandler(t *testing.T) {

	for _, data := range routestestdata {

		t.Log(data.desc)

		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			assert.Expect(t, data.tail, req.URL.EscapedPath())
			w.WriteHeader(http.StatusOK)
		})

		rh := webbed.NewHTTPRouteHandler()
		rh.SetRoute(data.route, handler)

		srv := httptest.NewServer(rh)
		defer srv.Close()

		client := srv.Client()

		_url, _ := url.Parse(srv.URL)
		_url.Path = data.path

		_, err := client.Get(_url.String())
		assert.Ok(t, err)
	}
}

func TestNestedHTTPRouteHandler(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	nested := webbed.NewHTTPRouteHandler()
	nested.SetRoute("nested", handler)

	root := webbed.NewHTTPRouteHandler()
	root.SetRoute("root", nested)

	srv := httptest.NewServer(root)
	defer srv.Close()

	client := srv.Client()

	u, _ := url.Parse(srv.URL)
	u.Path = "root/nested"

	resp, err := client.Get(u.String())
	assert.Ok(t, err)
	assert.Expect(t, http.StatusOK, resp.StatusCode)
}
