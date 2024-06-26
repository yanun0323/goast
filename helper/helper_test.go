package helper

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/yanun0323/goast/assert"
)

func TestGolang(t *testing.T) {
	a := assert.New(t)

	s := "0123456789"
	a.Require(string(s[2:5]) == "234", string(s[2:5]))
}

func TestHasPrefix(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		s      []byte
		prefix string
		ok     bool
	}{
		{s: []byte("123456"), prefix: "", ok: true},
		{s: []byte("123456"), prefix: "123", ok: true},
		{s: []byte("123456"), prefix: "456", ok: false},
		{s: []byte("123456"), prefix: "123456", ok: true},
		{s: []byte("123456"), prefix: "1234567", ok: false},
		{s: []byte(""), prefix: "", ok: true},
		{s: []byte(""), prefix: "123", ok: false},
	}

	for ti, tc := range testCases {
		desc := fmt.Sprintf("test case: %d", ti)
		t.Run(desc, func(t *testing.T) {
			result := HasPrefix(tc.s, tc.prefix)
			a.Require(result == tc.ok, fmt.Sprintf("expected %v, but get %v", tc.ok, result))
		})
	}
}

func TestHasSuffix(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		s      []byte
		suffix string
		ok     bool
	}{
		{s: []byte("123456"), suffix: "", ok: true},
		{s: []byte("123456"), suffix: "456", ok: true},
		{s: []byte("123456"), suffix: "123", ok: false},
		{s: []byte("123456"), suffix: "123456", ok: true},
		{s: []byte("123456"), suffix: "0123456", ok: false},
		{s: []byte(""), suffix: "", ok: true},
		{s: []byte(""), suffix: "123", ok: false},
	}

	for ti, tc := range testCases {
		desc := fmt.Sprintf("test case: %d", ti)
		t.Run(desc, func(t *testing.T) {
			result := HasSuffix(tc.s, tc.suffix)
			a.Require(result == tc.ok, fmt.Sprintf("expected %v, but get %v", tc.ok, result))
		})
	}
}

func TestPrintTidy(t *testing.T) {
	a := assert.New(t)

	s := "\n \n \n\t\r"
	a.Equal(TidyText(s), "\\n·\\n·\\n -> ·")
}

func TestReadFile(t *testing.T) {
	a := assert.New(t)

	_, err := ReadFile("./empty/no.go")
	a.Require(errors.Is(err, os.ErrNotExist))
}
