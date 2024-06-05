package goast

import "testing"

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
/* comment */ Fn /* comment */ func( /* comment */
	/* comment */ int, /* comment */
	/* comment */ string, /* comment */
	/* comment */) /* comment */ ( /* comment */
	/* comment */ int, /* comment */
	/* comment */ error, /* comment */
	/* comment */) /* comment */
}`

	interfaceNode, err := extract([]byte(interfaceText))
	a.NoError(err, "extract interface text should be no error")

	interfaceResult := typeResetter{}.Run(interfaceNode)
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
	structNode.PrintAllNext()
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
