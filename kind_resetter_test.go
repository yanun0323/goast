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

func Hello(ctx context.Context) error {
	return nil
}`
	head, err := extract([]byte(text))
	a.NoError(err, fmt.Sprintf("extract text, err: %s", err))

	tail := _commonResetter.Run(head)

	_ = head.IterNext(func(n Node) bool {
		switch n.Text() {
		case "Hello":
			a.Require(n.Kind() == KindFuncName, "Hello should be KindFuncName")
		case "ctx":
			a.Require(n.Kind() == KindParamName, "ctx should be KindParamName")
		}
		return true
	})

	a.Require(tail == nil, "tail should be nil")
}
