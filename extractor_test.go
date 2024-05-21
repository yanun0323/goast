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
	println(msg)
} 
`
	n, err := extract([]byte(s))
	a.Require(err == nil, "extract no error", fmt.Sprintf("%+v", err))
	for i := range n {
		n[i].Print()
	}
	a.Require(len(n) == 19, "nodes length", fmt.Sprintf("%d", len(n)))
}
