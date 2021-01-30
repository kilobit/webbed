/* Copyright 2021 Kilobit Labs Inc. */

package webbed_test

import _ "fmt"
import _ "errors"
import "context"
import "os"

import "kilobit.ca/go/webbed"

import "kilobit.ca/go/tested/assert"
import "testing"

func TestWebbedTest(t *testing.T) {
	assert.Expect(t, true, true)
}

type cKey int

const (
	testKey1 cKey = iota
	testKey2
	testKey3
)

func TestStringFromCtx(t *testing.T) {

	tests := []struct {
		key interface{}
		val interface{}
		exp string
	}{
		{testKey1, "Hello World!", "Hello World!"},
		{testKey1, nil, ""},
		{testKey2, 2, ""},
	}

	for _, test := range tests {

		ctx := context.TODO()
		ctx = context.WithValue(ctx, test.key, test.val)

		result := webbed.StringFromCtx(ctx, test.key)
		assert.Expect(t, test.exp, result)
	}
}

func TestLoadCtxFromEnv(t *testing.T) {

	tests := []struct{
		env map[string]string
		keymap map[string]interface{}
	}{
		{map[string]string{"foo": "bar", "bing": "bang"},
		 map[string]interface{}{"foo": testKey1, "bing": testKey2}},
	}

	for _, test := range tests {

		for k, v := range test.env {
			_, found := os.LookupEnv(k)
			if found {
				t.Fatalf("Refusing to use existing env key in test, %s.", k)
			}

			err := os.Setenv(k, v)
			assert.Ok(t, err)
		}

		ctx := context.TODO()
		ctx = webbed.LoadCtxFromEnv(ctx, test.keymap)

		for ekey, ckey := range test.keymap {

			result := webbed.StringFromCtx(ctx, ckey)
			assert.Expect(t, test.env[ekey], result)
		}

		for ekey, _ := range test.env {
			err := os.Unsetenv(ekey)
			assert.Ok(t, err)
		}
	}
}
