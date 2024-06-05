package example

import (
	"testing"

	"github.com/yanun0323/goast"
	"github.com/yanun0323/goast/kind"
	"github.com/yanun0323/goast/scope"
)

func TestExample(t *testing.T) {
	f := "sample_test.go"
	ast, err := goast.ParseAst(f)
	if err != nil {
		t.Fatalf("parse ast from file: %s, err: %+v", f, err)
	}

	sc := make([]goast.Scope, 0, len(ast.Scope()))

	ast.IterScope(func(s goast.Scope) bool {
		switch s.Kind() {
		case scope.Type, scope.Func, scope.Package, scope.Import:
			sc = append(sc, s)

			var prev *goast.Node
			s.Node().IterNext(func(n *goast.Node) bool {
				if prev != nil {
					n.ReplacePrev(prev)
					prev = nil
				}

				switch n.Kind() {
				case kind.Comment:
					prev = n.Prev()
				}

				return true
			})
		}
		return true
	})

	ast = ast.SetScope(sc)

	s := "output/save_test"
	if err := ast.Save(s); err != nil {
		t.Fatalf("save ast to file: %s, err: %+v", f, err)
	}
}
