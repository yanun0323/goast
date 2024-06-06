package goast

import (
	"fmt"
	"testing"

	"github.com/yanun0323/goast/assert"
	"github.com/yanun0323/goast/kind"
)

func TestCommonResetter(t *testing.T) {
	a := assert.New(t)
	text := `func Hello(ctx context.Context, m map[string]string, f1 func(int) error, f2 func(num int8) (int32, error)) error {
	return (int, nil)
}`

	head, err := extract([]byte(text))
	a.NoError(err, fmt.Sprintf("extract text, err: %s", err))

	tail := resetKind(head)
	_ = head.IterNext(func(n *Node) bool {
		switch n.Text() {
		case "Hello":
			a.Equal(n.Kind(), kind.FuncName, fmt.Sprintf("%s should be kind.FuncName", n.Text()))
		case "ctx",
			"m",
			"f1",
			"f2":
			a.Equal(n.Kind(), kind.ParamName, fmt.Sprintf("%s should be kind.ParamName", n.Text()))
		case "context.Context",
			"map[string]string",
			"func(int) error",
			"func(num int8) (int32, error)":
			a.Equal(n.Kind(), kind.ParamType, fmt.Sprintf("%s should be kind.ParamType", n.Text()))
		}

		return true
	})

	a.Equal(tail.Next(), nil, "tail should be nil")
}
