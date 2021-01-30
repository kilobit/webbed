/* Copyright 2020 Kilobit Labs Inc. */

// Tests for the webbed package.
package limits_test

import _ "fmt"
import _ "errors"
import "sync"
import "time"
import "bytes"
import _ "strings"
import "io"
import "io/ioutil"
import _ "net/url"
import "net/http"
import "net/http/httptest"

import "kilobit.ca/go/webbed/limits"

import "testing"
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

func TestHTTPLimitsHandlerMaxBytes(t *testing.T) {

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

		lh := limits.NewHTTPLimitsHandler(handler)

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

var exceededLimit limits.LimitFunc = func(w http.ResponseWriter, req *http.Request) bool {
	w.WriteHeader(http.StatusBadRequest)
	return false
}

var allowedLimit limits.LimitFunc = func(w http.ResponseWriter, req *http.Request) bool {
	return true
}

var okHandler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func sleepyOkHandler(t time.Duration) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(t)
		w.WriteHeader(http.StatusOK)
	})
}

func TestHTTPLimitsHandler(t *testing.T) {

	tests := []struct {
		limits []limits.Limit
		path   string
		body   io.Reader
		exp    int
	}{
		{[]limits.Limit{}, "/", nil, http.StatusOK},
		{[]limits.Limit{exceededLimit}, "/", nil, http.StatusBadRequest},
		{[]limits.Limit{allowedLimit}, "/", nil, http.StatusOK},
		{[]limits.Limit{allowedLimit, exceededLimit}, "/", nil, http.StatusBadRequest},
	}

	for _, test := range tests {

		//t.Logf("%#v\n", test.limits)
		lh := limits.NewHTTPLimitsHandler(okHandler, test.limits...)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", test.path, test.body)

		lh.ServeHTTP(w, req)

		resp := w.Result()

		assert.Expect(t, test.exp, resp.StatusCode)
	}
}

// TODO: Currently we don't check responses... this test will mostly
// look for panics etc.
//
func TestRateLimit(t *testing.T) {

	t.Skip("Skipped until concurrency issue with the test run is resolved.")

	tests := []struct {
		reqs   int
		limit  int
		period time.Duration
		wait   time.Duration
		exp    int
	}{
		{1, 10, 1 * time.Second, 0, http.StatusOK},
		{12, 10, 1 * time.Second, 1 * time.Second, http.StatusTooManyRequests},
	}

	for _, test := range tests {

		rl := limits.RateLimit(test.limit, test.period, test.wait)

		//t.Logf("%#v\n", test.limits)
		lh := limits.NewHTTPLimitsHandler(sleepyOkHandler(time.Second/10), rl)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		wg := sync.WaitGroup{}

		for i := 0; i < test.reqs; i++ {

			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				lh.ServeHTTP(w, req)

				//resp := w.Result()

				// if i < test.limit {
				// 	assert.Expect(t, http.StatusOK, resp.StatusCode)
				// } else {
				// 	assert.Expect(t, http.StatusTooManyRequests, resp.StatusCode)
				//}
			}(i)
		}

		wg.Wait()
	}
}
