package goast

import (
	"fmt"
	"testing"
)

func TestCommonResetter(t *testing.T) {
	a := NewAssert(t)
	text := `package main 

import (
	"context"
)

func Hello(ctx context.Context, m map[string]string, f1 func(int) error, f2 func(num int8) (int32, error)) (string, error) {
	return (int, nil)
}`
	head, err := extract([]byte(text))
	a.NoError(err, fmt.Sprintf("extract text, err: %s", err))

	tail := kindReset(head)

	_ = head.IterNext(func(n *Node) bool {
		switch n.Text() {
		case "Hello":
			a.Require(n.Kind() == KindFuncName, "Hello should be KindFuncName")
		case "ctx":
			a.Require(n.Kind() == KindParamName, "ctx should be KindParamName")
		}

		n.Print()

		return true
	})

	a.Require(tail == nil, "tail should be nil")
}
