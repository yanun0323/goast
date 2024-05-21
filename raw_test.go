package goast

import "testing"

func TestRaw(t *testing.T) {
	a := NewAssert(t)

	raw, err := newRawFile("sample_test.go")
	a.NoError(err)

	a.Require(raw != nil)
	a.Require(raw.name == "sample_test.go")
	a.Require(raw.dir == ".")
}
