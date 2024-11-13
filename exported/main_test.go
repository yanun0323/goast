package exported

import (
	"testing"

	"github.com/yanun0323/goast"
	"github.com/yanun0323/goast/kind"
	"github.com/yanun0323/goast/scope"
)

func TestExample(t *testing.T) {
	f := "./sample_test.go"
	ast, err := goast.ParseAst(f)
	if err != nil {
		t.Fatalf("parse ast from file: %s, err: %+v", f, err)
	}

	sc := make([]goast.Scope, 0, len(ast.Scope()))

	ast.IterScope(func(s goast.Scope) bool {
		switch s.Kind() {
		case scope.Const, scope.Variable, scope.Type, scope.Func, scope.Package, scope.Import:
			takePrev := false
			s.Node().IterNext(func(n *goast.Node) bool {
				if takePrev {
					n.TakePrev()
					takePrev = false
				}

				switch n.Kind() {
				case kind.Comment:
					takePrev = true
				}

				return true
			})
			sc = append(sc, s)
		}
		return true
	})

	ast = ast.SetScope(sc)
	// ast.IterScope(func(s goast.Scope) bool {
	// 	s.Node().IterNext(func(n *goast.Node) bool {
	// 		n.Print()
	// 		return true
	// 	})
	// 	return true
	// })

	s := "./output/save_test"
	if err := ast.Save(s, false); err != nil {
		t.Fatalf("save ast to file: %s, err: %+v", f, err)
	}
}

func TestUsecase(t *testing.T) {
	f := "usecase_test.go"
	ast, err := goast.ParseAst(f)
	if err != nil {
		t.Fatalf("parse ast from file: %s, err: %+v", f, err)
	}

	ast.IterScope(func(s goast.Scope) bool {
		if _, ok := s.GetInterfaceName(); !ok {
			return true
		}

		s.Print()

		s.Node().IterNext(func(n *goast.Node) bool {
			n.Print()
			return true
		})
		return true
	})
}
