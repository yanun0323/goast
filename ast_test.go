package goast

import (
	"testing"
)

func TestParse(t *testing.T) {
	a := NewAssert(t)

	ff, err := ParseAst("sample_test.go")
	a.NoError(err)
	for _, sc := range ff.Scope() {
		sc.Print()
		for _, n := range sc.Node() {
			n.Print()
		}
	}
}
