/* Copyright 2020 Kilobit Labs Inc. */

// Tests for the informed package.
package informed_test

import _ "fmt"
import _ "errors"
import _ "strings"
import _ "io/ioutil"
import "net/url"
import "net/http"
import "net/http/httptest"
import "testing"
import "kilobit.ca/go/informed"
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

		rh := informed.NewHTTPRouteHandler()
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
