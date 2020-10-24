/* Copyright 2020 Kilobit Labs Inc. */

// Tests for the webbed package.
package webbed_test

import _ "fmt"
import _ "errors"
import "bytes"
import _ "strings"
import "io"
import "io/ioutil"
import _ "net/url"
import "net/http"
import "net/http/httptest"
import "testing"
import "kilobit.ca/go/webbed"
import "kilobit.ca/go/tested/assert"

func TestLimitsTest(t *testing.T) {

	assert.Expect(t, true, true)
}

var limitstestdata = []struct {
	desc  string
	size  int
	code  int
	isErr bool
}{
	{
		"too big",
		1<<20 + 1,
		413,
		true,
	},

	{
		"just right",
		500,
		200,
		false,
	},

	{
		"empty",
		0,
		200,
		false,
	},
}

func TestHTTPLimitsHandler(t *testing.T) {

	for _, data := range limitstestdata {

		t.Log(data.desc)

		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			_, err := ioutil.ReadAll(req.Body)

			switch {

			case data.isErr == false && err == nil || err == io.EOF:
				w.WriteHeader(http.StatusOK)

			case data.isErr == true && err.Error() == "http: request body too large":
				w.WriteHeader(http.StatusRequestEntityTooLarge)

			default:
				t.Logf("%#v", err)
				w.WriteHeader(http.StatusBadRequest)
			}
		})

		lh := webbed.NewHTTPLimitsHandler(handler)

		srv := httptest.NewServer(lh)
		defer srv.Close()

		client := srv.Client()

		bs := bytes.Repeat([]byte{'X'}, data.size)
		//t.Log(len(bs))

		resp, err := client.Post(srv.URL, "text/plain", bytes.NewReader(bs))
		assert.Ok(t, err)

		assert.Expect(t, data.code, resp.StatusCode)
	}
}
