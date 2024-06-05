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

func TestTypeResetter(t *testing.T) {
	a := NewAssert(t)

	interfaceText := `type Student interface {
		Meow(int32)
		SelfIntroduction(string) []byte
		Laugh(loud float64) (bool, error)
		Learn(fn func(string) map[string]string) (func(int)int64, error)
}`

	structText := `type Student struct {
Name 								string
Age /* negative means not born */ 	int8
Relationship						map[string]string
FuncRelationship					map[string]func(int, int8) error
FuncRelationship2					map[string]func(n int32, nn int64) error
FuncRelationship3					map[string]func(uint, uint8) (uint64, error)
}`

	interfaceNode, err := extract([]byte(interfaceText))
	a.NoError(err, "extract interface text should be no error")

	interfaceResult := typeResetter{}.Run(interfaceNode)
	// interfaceNode.PrintAllNext()
	a.Nil(interfaceResult, "reset interface node", interfaceResult.Text())

	containByteSlice := false
	interfaceNode.IterNext(func(n *Node) bool {
		switch n.Text() {
		case "Student":
			a.Equal(n.Kind(), KindTypeName, n.Text())
		case "Meow", "SelfIntroduction", "Laugh", "Learn":
			a.Equal(n.Kind(), KindFuncName, n.Text())
		case "loud", "fn":
			a.Equal(n.Kind(), KindParamName, n.Text())
		case "int32", "string", "float64", "bool", "error", "func(string) map[string]string", "func(int)int64":
			a.Equal(n.Kind(), KindParamType, n.Text())
		case "[]byte":
			containByteSlice = true
			a.Equal(n.Kind(), KindParamType, n.Text())
		}
		return true
	})
	a.Require(containByteSlice, "should contain []byte")

	structNode, err := extract([]byte(structText))
	a.NoError(err, "extract struct text should be no error")

	structResult := typeResetter{}.Run(structNode)
	a.Nil(structResult, "reset struct node", structResult.Text())
	structNode.IterNext(func(n *Node) bool {
		switch n.Text() {
		case "Student":
			a.Equal(n.Kind(), KindTypeName)
		}
		return true
	})
}
