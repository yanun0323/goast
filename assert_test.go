package goast

import (
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
)

type Assert struct {
	t     *testing.T
	Debug int
}

func NewAssert(t *testing.T) Assert {
	d, err := strconv.Atoi(os.Getenv("DEBUG"))
	if err != nil {
		d = 0
	}

	return Assert{
		t:     t,
		Debug: d,
	}
}

func (a Assert) Nil(value any, msg ...string) {
	if a.isNil(value) {
		return
	}

	a.t.Helper()

	if len(msg) != 0 {
		a.t.Fatalf("%s: should be nil, but got (%+v), %s", a.t.Name(), value, strings.Join(msg, " "))
	}

	a.t.Fatalf("%s: should be nil, but got (%+v)", a.t.Name(), value)
}

func (a Assert) NotNil(value any, msg ...string) {
	if !a.isNil(value) {
		return
	}

	a.t.Helper()

	if len(msg) != 0 {
		a.t.Fatalf("%s: should be not nil, but got (%+v), %s", a.t.Name(), value, strings.Join(msg, " "))
	}

	a.t.Fatalf("%s: should be not nil, but got (%+v)", a.t.Name(), value)
}

func (a Assert) Equal(value, expected any, msg ...string) {
	rv, re := reflect.ValueOf(value), reflect.ValueOf(expected)
	if rv.Kind() == re.Kind() && rv.Equal(re) {
		return
	}

	a.t.Helper()

	if len(msg) != 0 {
		a.t.Fatalf("%s: should be equal, but got value: (%+v), expected: (%+v), %s", a.t.Name(), value, expected, strings.Join(msg, " "))
	}

	a.t.Fatalf("%s: should be equal, but got value: (%+v), expected: (%+v)", a.t.Name(), value, expected)
}

func (a Assert) NotEqual(value, expected any, msg ...string) {
	rv, re := reflect.ValueOf(value), reflect.ValueOf(expected)
	if rv.Kind() != re.Kind() || !rv.Equal(re) {
		return
	}

	a.t.Helper()

	if len(msg) != 0 {
		a.t.Fatalf("%s: should be no equal, but got value: (%+v), expected: (%+v), %s", a.t.Name(), value, expected, strings.Join(msg, " "))
	}

	a.t.Fatalf("%s: should be no equal, but got value: (%+v), expected: (%+v)", a.t.Name(), value, expected)
}

func (a Assert) Error(err error, msg ...string) {
	if err != nil {
		return
	}

	a.t.Helper()

	if len(msg) != 0 {
		a.t.Fatalf("%s: should be error, %s", a.t.Name(), strings.Join(msg, " "))
	}

	a.t.Fatalf("%s: should be error", a.t.Name())
}

func (a Assert) NoError(err error, msg ...string) {
	if err == nil {
		return
	}

	a.t.Helper()

	if len(msg) != 0 {
		a.t.Fatalf("%s: should be no error, %s, err: %+v", a.t.Name(), strings.Join(msg, " "), err)
	}

	a.t.Fatalf("%s:should be no error, err: %+v", a.t.Name(), err)
}

func (a Assert) Require(ok bool, msg ...string) {
	if ok {
		return
	}

	a.t.Helper()

	if len(msg) != 0 {
		a.t.Fatalf("require: %s, %s", a.t.Name(), strings.Join(msg, " "))
	}

	a.t.Fatalf("require: %s", a.t.Name())
}

func (a Assert) NoPanic(fn func()) {
	a.t.Helper()
	if ok, msg, _ := a.Call(fn); ok {
		a.t.Fatalf("%s: should be no panic, err: %s", a.t.Name(), msg)
	}
}

func (a Assert) Panic(fn func()) {
	a.t.Helper()
	if ok, msg, _ := a.Call(fn); !ok {
		a.t.Fatalf("%s: should be panic, err: %s", a.t.Name(), msg)
	}
}

func (Assert) Call(fn func()) (didPanic bool, message interface{}, stack string) {
	didPanic = true
	defer func() {
		message = recover()
		if didPanic {
			stack = string(debug.Stack())
		}
	}()

	// call the target function
	fn()
	didPanic = false

	return
}

func (Assert) isNil(i any) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
