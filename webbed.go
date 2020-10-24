/* Copyright 2020 Kilobit Labs Inc. */

package webbed

import _ "fmt"
import _ "errors"
import "net/http"

// Return an error on the calling HTTP connection.
//
// code is a http.Status code.
func ServeError(w http.ResponseWriter, code int, err error) {

	http.Error(w, err.Error(), code)
}
