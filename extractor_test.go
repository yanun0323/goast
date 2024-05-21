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
	a.Require(len(n) != 0, "extracted nodes not empty")
	for i := range n {
		n[i].Print()
	}
}
