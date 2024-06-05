package goast

import (
	"testing"

	"github.com/yanun0323/goast/assert"
	"github.com/yanun0323/goast/charset"
	"github.com/yanun0323/goast/kind"
)

func TestTrie(t *testing.T) {
	a := assert.New(t)

	root := newTrie(map[string]*kind.Kind{})
	basic := kind.Basic
	_ = root.AddText("any", &basic)
	k, ok := root.FindText("any")
	a.Require(ok, "find 'any'")
	a.Require(k != nil && *k == kind.Basic, "k 'any'", k.PointerString())

	k, ok = root.FindText("an")
	a.Require(ok, "find 'an'")
	a.Require(k == nil, "k 'an'", k.PointerString())

	k, ok = root.FindText("ann")
	a.Require(!ok, "find 'ann'")
	a.Require(k == nil, "k 'ann'", k.PointerString())

	comment := kind.Comment
	_ = root.AddText("ayy", &comment)

	k, ok = root.FindText("any")
	a.Require(ok, "find 'any' 2")
	a.Require(k != nil && *k == kind.Basic, "k 'any' 2", k.PointerString())

	k, ok = root.FindText("ayy")
	a.Require(ok, "find 'ayy'")
	a.Require(k != nil && *k == kind.Comment, "k 'ayy'", k.PointerString())

	k, ok = root.FindByte([]byte("ayy"))
	a.Require(ok, "find 'ayy' 2")
	a.Require(k != nil && *k == kind.Comment, "k 'ayy' 2", k.PointerString())
}

func TestTrieComment(t *testing.T) {
	a := assert.New(t)

	root := newTrie(map[string]trie[bool]{
		"//": newTrie(map[string]bool{"\n": true}),
		"/*": newTrie(map[string]bool{"*/": true}),
	})

	comment := []byte("//")
	innerComment := []byte("/*")

	tr, ok := root.FindByte(comment)
	a.Require(ok, "found comment ok")
	a.Require(tr != nil, "found comment trie")

	tr, ok = root.FindByte(innerComment)
	a.Require(ok, "found innerComment ok")
	a.Require(tr != nil, "found innerComment trie")
}

func TestTrieFindReversely(t *testing.T) {
	a := assert.New(t)

	tr := newTrie(charset.New("123"))
	_, ok := tr.FindText("123")
	a.Require(ok, "find 123")

	_, ok = tr.FindTextReversely("321")
	a.Require(ok, "find reversely 321")

	_, ok = tr.FindByte([]byte("123"))
	a.Require(ok, "find byte 123")

	_, ok = tr.FindByteReversely([]byte("321"))
	a.Require(ok, "find byte reversely 321")
}
