/* Copyright 2020 Kilobit Labs Inc. a*/

// A simple structured logging package.
//
package logger

import "fmt"
import _ "errors"
import "io"
import "time"
import "crypto/rand"
import "encoding/json"

type FieldHandler func() interface{}

func UUIDL4Field() interface{} {
	uuid, _ := newUUID()
	return uuid
}

func TimestampField() interface{} {
	return time.Now()
}

type Encoder func(map[string]interface{}) ([]byte, error)

var EncoderJSON Encoder = func(m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type ErrorHandler func(error)

func NOPErrorHandler(error) {}

type Logger struct {
	fields      map[string]FieldHandler
	params      []string // nth parameter name
	messageName string
	extrasName  string
	enc         Encoder
	w           io.Writer
	ehandler    ErrorHandler
}

func New(w io.Writer) *Logger {
	return &Logger{
		map[string]FieldHandler{},
		[]string{},
		"message",
		"extras",
		EncoderJSON,
		w,
		NOPErrorHandler,
	}
}

// Log the given message and parameters.
//
func (lg Logger) Log(message interface{}, params ...interface{}) {

	var err error

	m := map[string]interface{}{}

	for i, name := range lg.params {
		param, ok := readElement(params, i)
		if ok {
			m[name] = param
		}
	}

	for name, fh := range lg.fields {
		m[name] = fh()
	}

	// Add extra params to the extras.
	if len(params) > len(lg.params) {
		m[lg.extrasName] = params[len(lg.params):]
	}

	m[lg.messageName] = message

	bs, err := lg.enc(m)
	_, err = fmt.Fprintln(lg.w, string(bs))
	lg.ehandler(err)
}

func (lg *Logger) SetParams(names ...string) {
	lg.params = names
}

func (lg *Logger) SetField(name string, fh FieldHandler) {
	lg.fields[name] = fh
}

// newUUID generates a random UUID according to RFC 4122
//
// Retrieved from https://play.golang.org/p/4FkNSiUDMg
//
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func readElement(els []interface{}, i int) (interface{}, bool) {

	if i < 0 || i > len(els) {
		return "", false
	}

	return els[i], true
}
