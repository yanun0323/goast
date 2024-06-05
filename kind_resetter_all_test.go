package goast

import (
	"fmt"
	"testing"
)

func TestCommonResetter(t *testing.T) {
	a := NewAssert(t)
	text := `func Hello(ctx context.Context, m map[string]string, f1 func(int) error, f2 func(num int8) (int32, error)) error {
	return (int, nil)
}`

	println("Extract:")
	head, err := extract([]byte(text))
	a.NoError(err, fmt.Sprintf("extract text, err: %s", err))

	tail := kindReset(head)
	head.PrintAllNext()
	_ = head.IterNext(func(n *Node) bool {
		switch n.Text() {
		case "Hello":
			a.Equal(n.Kind(), KindFuncName, fmt.Sprintf("%s should be KindFuncName", n.Text()))
		case "ctx",
			"m",
			"f1",
			"f2":
			a.Equal(n.Kind(), KindParamName, fmt.Sprintf("%s should be KindParamName", n.Text()))
		case "context.Context",
			"map[string]string",
			"func(int) error",
			"func(num int8) (int32, error)",
			"string",
			"error":
			a.Equal(n.Kind(), KindParamType, fmt.Sprintf("%s should be KindParamType", n.Text()))
		}

		return true
	})

	a.Require(tail == nil, "tail should be nil")
}
