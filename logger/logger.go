/* Copyright 2021 Kilobit Labs Inc. */

// A simple structured logging package.
//
// Fields are automatically added Log data.
// Parameters are Log data with pre-set names.
//
package logger

import "fmt"
import _ "errors"
import "io"
import "strings"
import "sort"
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

type Encoder func(m map[string]interface{}, keys []string) ([]byte, error)

var EncoderJSON Encoder = func(m map[string]interface{}, keys []string) ([]byte, error) {

	return json.Marshal(m)
}

var EncoderTXT Encoder = func(m map[string]interface{}, keys []string) ([]byte, error) {

	if keys == nil {
		keys = make([]string, len(m))
		i := 0
		for k, _ := range m {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
	}

	sb := &strings.Builder{}

	fmt.Fprintf(sb, "Entry: %d\n", time.Now().Unix())

	for _, key := range keys {
		fmt.Fprintf(sb, "\t%s: %v\n", key, m[key])
	}

	return []byte(sb.String()), nil
}

type ErrorHandler func(error)

func NOPErrorHandler(error) {}

func NewWriteErrorHandler(w io.Writer) ErrorHandler {

	return ErrorHandler(func(err error) {
		if err != nil {
			fmt.Fprintf(w, "Logger Error: %s", err.Error())
		}
	})
}

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

func (lg *Logger) EncodeAsText() {
	lg.enc = EncoderTXT
}

// Write errors encountered while logging to the writer.
//
func (lg *Logger) WriteLogErrors(w io.Writer) {
	lg.ehandler = NewWriteErrorHandler(w)
}

// Log the given message and parameters.
//
func (lg *Logger) Log(message interface{}, params ...interface{}) {

	var err error

	m := map[string]interface{}{}
	keys := []string{lg.messageName}

	for i, name := range lg.params {
		param, ok := readElement(params, i)
		if ok {
			m[name] = param
			keys = append(keys, name)
		}
	}

	fkeys := []string{}
	for name, fh := range lg.fields {
		m[name] = fh()
		fkeys = append(fkeys, name)
	}
	sort.Strings(fkeys)
	keys = append(keys, fkeys...)

	// Add extra params to the extras.
	if len(params) > len(lg.params) {
		m[lg.extrasName] = params[len(lg.params):]
		keys = append(keys, lg.extrasName)
	}

	m[lg.messageName] = message

	bs, err := lg.enc(m, keys)
	lg.ehandler(err)

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
