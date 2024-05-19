package goast

import "testing"

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
		a.t.Fatalf("%s, err: %+v", msg[0], err)
	}

	a.t.Fatalf("%+v", err)
}

func (a Assert) Require(ok bool, msg ...string) {
	if ok {
		return
	}

	if len(msg) != 0 {
		a.t.Fatal("require")
	}

	a.t.Fatalf("require, %s", msg[0])
}
