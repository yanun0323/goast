package goast

import (
	"fmt"
	"testing"

	"github.com/yanun0323/goast/assert"
	"github.com/yanun0323/goast/kind"
)

func TestInterfaceResetter(t *testing.T) {
	a := assert.New(t)

	interfaceText := `type Student interface {
		Meow(int32)
		SelfIntroduction(string) []byte
		Laugh(loud float64) (bool, error)
		Learn(fn func(string) map[string]string) (func(int)int64, error)
}`

	interfaceNode, err := extract([]byte(interfaceText))
	a.NoError(err, "extract interface text should be no error")

	interfaceResult := typeResetter{}.Run(interfaceNode)
	a.Nil(interfaceResult, "reset interface node", interfaceResult.Text())

	interfaceAssertMap := map[string]kind.Kind{
		"Student": kind.TypeName,
		"Meow":    kind.FuncName, "SelfIntroduction": kind.FuncName, "Laugh": kind.FuncName, "Learn": kind.FuncName,
		"loud": kind.ParamName, "fn": kind.ParamName,
		"[]byte": kind.ParamType,
	}
	interfaceNode.IterNext(func(n *Node) bool {
		if expected, ok := interfaceAssertMap[n.Text()]; ok {
			a.Equal(n.Kind(), expected, n.Text())
			delete(interfaceAssertMap, n.Text())
		}
		return true
	})
	a.Equal(len(interfaceAssertMap), 0, "interfaceAssertMap", fmt.Sprintf("%+v", interfaceAssertMap))
}

func TestStructResetter(t *testing.T) {
	a := assert.New(t)

	structText := `type Student[T any] struct {
		Name                                string
		Age /* negative means not born */, Age2	int8
		Relationship                        map[string]string
		FuncRelationship                    map[string]func(int, int8) error
		FuncRelationship2                   map[string]func(n int32, nn int64) error
		FuncRelationship3                   map[string]func(uint, uint8) (uint64, error)
		/* 123 */ Fn /* 234 */ func( /* 345 */
			/* 456 */ context.Context, /* 567 */
			/* 678 */ string, /* 789 */
			/* 321 */) /* 432 */ ( /* 543 */
			/* 654 */ int, /* 765 */
			/* 876 */ error, /* 987 */
			/* 000 */) /* comment */
		}`

	structNode, err := extract([]byte(structText))
	a.NoError(err, "extract struct text should be no error")

	structResult := typeResetter{}.Run(structNode)
	a.Nil(structResult, "reset struct node", structResult.Text())
	structAssertMap := map[string]kind.Kind{
		"Student": kind.TypeName,
		"Age":     kind.ParamName, "Age2": kind.ParamName, "Relationship": kind.ParamName, "FuncRelationship": kind.ParamName,
		"FuncRelationship2": kind.ParamName, "FuncRelationship3": kind.ParamName, "Fn": kind.ParamName,
		"n": kind.ParamName, "nn": kind.ParamName, "context.Context": kind.ParamType,
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

func TestOtherResetter(t *testing.T) {
	a := assert.New(t)

	otherText := `type setter  /* comment */ func( /* comment */
		/* comment */ context.Context, /* comment */
		/* comment */ string, /* comment */
		/* comment */) /* comment */ ( /* comment */
		/* comment */ int, /* comment */
		/* comment */ error, /* comment */
		/* comment */) /* comment */

type reMap map[uint8]float64

type reSlice []uint64`

	otherNode, err := extract([]byte(otherText))
	a.NoError(err, "extract other text should be no error")

	otherResult := typeResetter{}.Run(otherNode)
	a.Nil(otherResult, "reset other node", otherResult.Text())

	otherAssertMap := map[string]kind.Kind{}
	otherNode.IterNext(func(n *Node) bool {
		if expected, ok := otherAssertMap[n.Text()]; ok {
			a.Equal(n.Kind(), expected, n.Text())
			delete(otherAssertMap, n.Text())
		}
		return true
	})
	a.Equal(len(otherAssertMap), 0, "otherAssertMap", fmt.Sprintf("%+v", otherAssertMap))
}

func TestStructResetterGetNameCount(t *testing.T) {
	a := assert.New(t)

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
