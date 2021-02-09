/* Copyright 2021 Kilobit Labs Inc. */

// Simple utilities for building Web Apps in Golang.
//
package webbed

import _ "fmt"
import _ "errors"
import "context"
import "os"
import "net/http"

// Return an error on the calling HTTP connection.
//
// code is a http.Status code.
func ServeError(w http.ResponseWriter, code int, err error) {

	http.Error(w, err.Error(), code)
}

// Retrieve a string value from a context.
//
// Missing values or those of a non-string type will retrun "".
//
func StringFromCtx(ctx context.Context, key interface{}) string {
	str, ok := ctx.Value(key).(string)
	if !ok {
		str = ""
	}

	return str
}

// Convert the environment map to a context.
//
func LoadCtxFromEnv(ctx context.Context, keymap map[string]interface{}) context.Context {

	for ekey, ckey := range keymap {
		value, ok := os.LookupEnv(ekey)
		if ok {
			ctx = context.WithValue(ctx, ckey, value)
		}
	}

	return ctx
}
