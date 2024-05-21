package goast

import (
	"strings"
	"testing"
)

func TestTrimSpace(t *testing.T) {
	a := NewAssert(t)

	s := "\r\t\n123\n\t\r"
	s = strings.TrimSpace(s)
	a.Require(s == "123")
}
