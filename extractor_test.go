package goast

import (
	"fmt"
	"testing"
)

func TestExtract(t *testing.T) {
	a := NewAssert(t)

	s := `
func Print(msg string, d decimal.Decimal) {
	// This is comment
	/* And This is (hello)
	Inner Comment */
	s := "Hello\nThere"
	println(msg)
} 
`
	n, err := extract([]byte(s))
	a.Require(err == nil, "extract no error", fmt.Sprintf("%+v", err))

	count := 0
	_ = n.IterNext(func(n *Node) bool {
		if a.Debug == 2 {
			n.Print()
		}
		count++
		return true
	})
	a.Require(count == 40, "nodes length", fmt.Sprintf("%d", count))
}
