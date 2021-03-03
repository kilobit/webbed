/* Copyright 2021 Kilobit Labs Inc. a*/

package logger

import _ "fmt"
import _ "errors"
import "net/http"

type ResponseWriter struct {
	code int
	http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{0, w}
}

func (w *ResponseWriter) StatusCode() int {
	return w.code
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
