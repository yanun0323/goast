package goast

import (
	"os"
	"testing"
)

var (
	_testMultilineMessage = `hello there,
This is a multiline message.
It's for parsing file test.`
	_testSlice = []string{"foo", "bar"}
	_testArray = [2]string{"foo", "bar"}
)

const (
	_testConst = "Hello World!"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
