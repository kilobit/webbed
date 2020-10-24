/* Copyright 2020 Kilobit Labs Inc. a*/

package logger_test

import _ "fmt"
import _ "errors"
import "log"
import "net/url"

import "testing"
import "kilobit.ca/go/tested/assert"

import "kilobit.ca/go/webbed/logger"

func TestLoggerTest(t *testing.T) {

	assert.Expect(t, true, true)
}

func TestLogger(t *testing.T) {

	lg := logger.New(log.Writer())
	lg.SetParams("map", "fooer")
	lg.SetField("txid", logger.UUIDL4Field)
	lg.SetField("timestamp", logger.TimestampField)

	values := url.Values{}

	lg.Log(struct {
		Type   string
		Values url.Values
	}{
		"form-data",
		values,
	},
		map[string]interface{}{
			"Foo":      "bar",
			"severity": "high",
		},
		"Foo", "Bar", "Bing!",
	)
}
