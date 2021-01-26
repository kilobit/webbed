/* Copyright 2021 Kilobit Labs Inc. */

package webbed_test

import _ "fmt"
import _ "errors"
import "context"

import "kilobit.ca/go/webbed"

import "kilobit.ca/go/tested/assert"
import "testing"

func TestServerTest(t *testing.T) {

	assert.Expect(t, true, true)
}

type cKey int

const(
	testKey1 cKey = iota
	testKey2
	testKey3
)

func TestStringFromCtx(t *testing.T) {

	tests := []struct{
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
