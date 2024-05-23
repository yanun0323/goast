package goast

import "testing"

func TestSet(t *testing.T) {
	a := NewAssert(t)

	var s set[int]
	a.Require(s.Contain(5) == false, "set contain 5")
	s.Insert(6)
	a.Require(s.Contain(6), "set contain 6")
}
