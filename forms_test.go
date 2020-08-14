/* Copyright 2020 Kilobit Labs Inc. */

// Tests for the informed package.
package informed_test

import _ "fmt"
import _ "errors"
import "strings"
import "io/ioutil"
import "net/url"
import "net/http"
import "net/http/httptest"
import "testing"
import "kilobit.ca/go/informed"
import "kilobit.ca/go/tested/assert"

func TestFormsTest(t *testing.T) {

	assert.Expect(t, true, true)
}

var formtestdata = []struct {
	desc   string
	values url.Values
	code   int
}{
	{
		"simple",
		url.Values{"foo": []string{"bar"}},
		http.StatusOK,
	},

	{
		"empty",
		url.Values{},
		http.StatusOK,
	},
}

func TestHTTPFormHandlerPost(t *testing.T) {

	for _, data := range formtestdata {

		t.Log(data.desc)

		f := informed.FormHandlerFunc(func(values url.Values) (int, error) {
			assert.ExpectDeep(t, data.values, values)
			return 0, nil
		})

		h := informed.NewHTTPFormHandler(f)

		srv := httptest.NewServer(h)
		defer srv.Close()

		client := srv.Client()
		resp, err := client.PostForm(srv.URL, data.values)
		assert.Ok(t, err)

		bs, err := ioutil.ReadAll(resp.Body)

		assert.Expect(t, "Ok", strings.TrimSpace(string(bs)))
	}
}

func TestHTTPFormHandlerGet(t *testing.T) {

	for _, data := range formtestdata {

		t.Log(data.desc)

		f := informed.FormHandlerFunc(func(values url.Values) (int, error) {
			assert.ExpectDeep(t, data.values, values)
			return 0, nil
		})

		h := informed.NewHTTPFormHandler(f)

		srv := httptest.NewServer(h)
		defer srv.Close()

		client := srv.Client()

		_url, _ := url.Parse(srv.URL)
		_url.RawQuery = data.values.Encode()

		resp, err := client.Get(_url.String())
		assert.Ok(t, err)

		bs, err := ioutil.ReadAll(resp.Body)

		assert.Expect(t, "Ok", strings.TrimSpace(string(bs)))
	}
}