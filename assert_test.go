package goast

import (
	"strings"
	"testing"
)

type Assert struct {
	t *testing.T
}

func NewAssert(t *testing.T) Assert {
	return Assert{t: t}
}

func (a Assert) NoError(err error, msg ...string) {
	if err == nil {
		return
	}

	a.t.Helper()

	if len(msg) != 0 {
		a.t.Fatalf("error: %s, %s, err: %+v", a.t.Name(), strings.Join(msg, " "), err)
	}

	a.t.Fatalf("error: %s, err: %+v", a.t.Name(), err)
}

func (a Assert) Require(ok bool, msg ...string) {
	if ok {
		return
	}

	a.t.Helper()

	if len(msg) != 0 {
		a.t.Fatalf("require: %s, %s, err: ", a.t.Name(), strings.Join(msg, " "))
	}

	a.t.Fatalf("require: %s", a.t.Name())
}
