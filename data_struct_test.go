package goast

import (
	"testing"

	"github.com/yanun0323/goast/assert"
)

func TestSet(t *testing.T) {
	a := assert.New(t)

	var s set[int]
	a.Require(s.Contain(5) == false, "set contain 5")
	s.Insert(6)
	a.Require(s.Contain(6), "set contain 6")
}
