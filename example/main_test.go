package example

import (
	"testing"

	"github.com/yanun0323/goast"
)

func TestExample(t *testing.T) {
	f := "sample_test.go"
	ast, err := goast.ParseAst(f)
	if err != nil {
		t.Fatalf("parse ast for file: %s, err: %+v", f, err)
	}

	for _, sc := range ast.Scope() {
		sc.Node().IterNext(func(n *goast.Node) bool {
			n.DebugPrint()
			return true
		})
	}

}
