package goast

import (
	"fmt"
	"testing"
)

func TestGolang(t *testing.T) {
	a := NewAssert(t)

	s := "0123456789"
	a.Require(string(s[2:5]) == "234", string(s[2:5]))
}

func TestHasPrefix(t *testing.T) {
	a := NewAssert(t)

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
			t.Log(desc)
			result := hasPrefix(tc.s, tc.prefix)
			a.Require(result == tc.ok, fmt.Sprintf("expected %v, but get %v", tc.ok, result))
		})
	}
}

func TestHasSuffix(t *testing.T) {
	a := NewAssert(t)

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
			t.Log(desc)
			result := hasSuffix(tc.s, tc.suffix)
			a.Require(result == tc.ok, fmt.Sprintf("expected %v, but get %v", tc.ok, result))
		})
	}
}

func TestPrintTidy(t *testing.T) {
	a := NewAssert(t)

	s := "\n \n \n\t\r"
	a.Equal(printTidy(s), "\\n\\s\\n\\s\\n\\t\\r")
}
