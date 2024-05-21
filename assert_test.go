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

	if len(msg) != 0 {
		a.t.Fatalf("%s, err: %+v", strings.Join(msg, " "), err)
	}

	a.t.Fatalf("%+v", err)
}

func (a Assert) Require(ok bool, msg ...string) {
	if ok {
		return
	}

	if len(msg) != 0 {
		a.t.Fatalf("%s, err: expected true bug get false", strings.Join(msg, " "))
	}

	a.t.Fatal("err: expected true bug get false")
}
