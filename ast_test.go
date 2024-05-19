package goast

import (
	"testing"
)

func TestParse(t *testing.T) {
	a := NewAssert(t)

	ff, err := ParseAst("sample_test.go")
	a.NoError(err)
	for _, scope := range ff.Defines() {
		scope.Print()
	}
}
