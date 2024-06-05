package goast

import (
	"fmt"
	"testing"
)

func TestTypeResetter(t *testing.T) {
	a := NewAssert(t)

	interfaceText := `type Student interface {
		Meow(int32)
		SelfIntroduction(string) []byte
		Laugh(loud float64) (bool, error)
		Learn(fn func(string) map[string]string) (func(int)int64, error)
}`

	structText := `type Student struct {
Name                                string
Age /* negative means not born */, Age2	int8
Relationship                        map[string]string
FuncRelationship                    map[string]func(int, int8) error
FuncRelationship2                   map[string]func(n int32, nn int64) error
FuncRelationship3                   map[string]func(uint, uint8) (uint64, error)
/* comment */ Fn /* comment */ func( /* comment */
	/* comment */ context.Context, /* comment */
	/* comment */ string, /* comment */
	/* comment */) /* comment */ ( /* comment */
	/* comment */ int, /* comment */
	/* comment */ error, /* comment */
	/* comment */) /* comment */
}`

	println()
	println("interface resetter test")
	interfaceNode, err := extract([]byte(interfaceText))
	a.NoError(err, "extract interface text should be no error")

	interfaceResult := typeResetter{}.Run(interfaceNode)
	a.Nil(interfaceResult, "reset interface node", interfaceResult.Text())

	interfaceAssertMap := map[string]Kind{
		"Student": KindTypeName,
		"Meow":    KindFuncName, "SelfIntroduction": KindFuncName, "Laugh": KindFuncName, "Learn": KindFuncName,
		"loud": KindParamName, "fn": KindParamName,
		"int32": KindParamType, "string": KindParamType, "float64": KindParamType, "bool": KindParamType,
		"error": KindParamType, "[]byte": KindParamType,
	}
	interfaceNode.IterNext(func(n *Node) bool {
		if expected, ok := interfaceAssertMap[n.Text()]; ok {
			a.Equal(n.Kind(), expected, n.Text())
			delete(interfaceAssertMap, n.Text())
		}
		return true
	})
	if len(interfaceAssertMap) != 0 {
		interfaceNode.PrintIter()
	}
	a.Equal(len(interfaceAssertMap), 0, "interfaceAssertMap", fmt.Sprintf("%+v", interfaceAssertMap))

	structNode, err := extract([]byte(structText))
	a.NoError(err, "extract struct text should be no error")

	structResult := typeResetter{}.Run(structNode)
	structNode.PrintIter()
	a.Nil(structResult, "reset struct node", structResult.Text())

	structAssertMap := map[string]Kind{
		"Student": KindTypeName,
		"Age":     KindParamName, "Age2": KindParamName, "Relationship": KindParamName, "FuncRelationship": KindParamName,
		"FuncRelationship2": KindParamName, "FuncRelationship3": KindParamName, "Fn": KindParamName,
		"n": KindParamName, "nn": KindParamName,
		"string": KindParamType, "int8": KindParamType, "map[string]string": KindParamType, "int": KindParamType,
		"int32": KindParamType, "int64": KindParamType, "uint": KindParamType, "uint64": KindParamType,
		"error": KindParamType, "context.Context": KindParamType,
	}
	structNode.IterNext(func(n *Node) bool {
		if expected, ok := structAssertMap[n.Text()]; ok {
			a.Equal(n.Kind(), expected, n.Text())
			delete(structAssertMap, n.Text())
		}
		return true
	})
	a.Equal(len(structAssertMap), 0, "structAssertMap", fmt.Sprintf("%+v", structAssertMap))
}

func TestStructResetterGetNameCount(t *testing.T) {
	a := NewAssert(t)

	s1 := "\t\tRepository\n"
	s1n, err := extract([]byte(s1))
	a.NoError(err, "s1n")
	a.Equal(structResetter{}.getRowNameCount(s1n), 0, "s1n count")

	s2 := "\t\trepo1, repo2 Repository\n"
	s2n, err := extract([]byte(s2))
	a.NoError(err, "s2n")
	a.Equal(structResetter{}.getRowNameCount(s2n), 2, "s2n count")

	s3 := "\t\trepo Repository\n"
	s3n, err := extract([]byte(s3))
	a.NoError(err, "s3n")
	a.Equal(structResetter{}.getRowNameCount(s3n), 1, "s3n count")

	s4 := "\t\trepo1, repo2, repo3, repo4 func(\n\tint,\n\tint,\n\t)(\n\tint,\n\terror,\n\t)\n"
	s4n, err := extract([]byte(s4))
	a.NoError(err, "s4n")
	a.Equal(structResetter{}.getRowNameCount(s4n), 4, "s4n count")

	s5 := "\t\trepository.Repository[Kind]\n"
	s5n, err := extract([]byte(s5))
	a.NoError(err, "s5n")
	a.Equal(structResetter{}.getRowNameCount(s5n), 0, "s5n count")
}
