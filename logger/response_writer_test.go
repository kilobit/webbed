/* Copyright 2021 Kilobit Labs Inc. */

package logger_test

import _ "fmt"
import _ "errors"
import "net/http"
import "net/http/httptest"

import "testing"
import "kilobit.ca/go/tested/assert"
import "kilobit.ca/go/tested/handlers"

import "kilobit.ca/go/webbed/logger"

func TestResponseWriterTest(t *testing.T) {
	assert.Expect(t, true, true)
}

func TestResponseWriter(t *testing.T) {

	handler := handlers.OkHandler

	lh := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		lrw := logger.NewResponseWriter(w)
		handler.ServeHTTP(lrw, req)

		assert.Expect(t, http.StatusOK, lrw.StatusCode())
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	lh.ServeHTTP(w, req)
	resp := w.Result()

	assert.Expect(t, http.StatusOK, resp.StatusCode)
}
