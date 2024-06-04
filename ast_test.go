package goast

import (
	"testing"
)

func TestParse(t *testing.T) {
	a := NewAssert(t)
	_ = a
	// return
	ff, err := ParseAst("sample_test.go")
	a.NoError(err, "parse ast error")
	a.Require(ff != nil, "nil ast check")
	a.Require(len(ff.Scope()) != 0, "scope length check")
	for _, sc := range ff.Scope() {
		_ = sc
		// sc.Print()
		// sc.Node().PrintAllNext()
	}
}
