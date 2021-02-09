/* Copyright 2021 Kilobit Labs Inc. */

package methods_test

import _ "fmt"
import _ "errors"

import "net/http"

import "kilobit.ca/go/webbed/methods"

import "testing"
import "net/http/httptest"
import "kilobit.ca/go/tested/assert"

func TestRoutesTest(t *testing.T) {
	assert.Expect(t, true, true)
}

var okHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func TestHandleMethods(t *testing.T) {

	tests := []struct {
		methods []string
		req     *http.Request
		status  int
	}{
		{[]string{""}, httptest.NewRequest("GET", "/", nil), http.StatusMethodNotAllowed},
		{[]string{"GET"}, httptest.NewRequest("GET", "/", nil), http.StatusOK},
		{[]string{"GET"}, httptest.NewRequest("PUT", "/", nil), http.StatusMethodNotAllowed},
		{[]string{"PUT", "GET"}, httptest.NewRequest("GET", "/", nil), http.StatusOK},
		{[]string{"PUT", "GET"}, httptest.NewRequest("PUT", "/", nil), http.StatusOK},
		{[]string{"PUT", "GET"}, httptest.NewRequest("POST", "/", nil), http.StatusMethodNotAllowed},
	}

	mh := methods.New(nil)

	for _, test := range tests {

		w := httptest.NewRecorder()
		handler := mh.HandleMethods(okHandler, test.methods...)
		handler.ServeHTTP(w, test.req)
		resp := w.Result()
		assert.Expect(t, test.status, resp.StatusCode, test.methods, test.req, test.status)
	}
}
