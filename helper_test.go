package goast

import "testing"

func TestGolang(t *testing.T) {
	a := NewAssert(t)

	s := "0123456789"
	a.Require(string(s[2:5]) == "234", string(s[2:5]))
}
