package goast

import (
	"errors"
	"testing"

	"github.com/yanun0323/goast/assert"
)

func TestParse(t *testing.T) {
	a := assert.New(t)
	_ = a
	{
		ff, err := ParseAst("./ast.go")
		a.NoError(err, "parse ast error")
		a.Require(ff != nil, "nil ast check")
		a.Require(len(ff.Scope()) != 0, "scope length check")
		for _, sc := range ff.Scope() {
			_ = sc
			// sc.Print()
			println(sc.Kind().String())
			sc.Node().Print()
			sc.Node().IterNext(func(n *Node) bool {
				if n != nil {
					n.Print()
				}
				return n != nil
			})
		}
	}

	{
		_, err := ParseAst("./ast_not_exist.go")
		a.Error(err, "parse ast error")
		a.Equal(errors.Is(err, ErrNotExist), true)
	}
}
