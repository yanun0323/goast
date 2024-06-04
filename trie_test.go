package goast

import "testing"

func TestTrie(t *testing.T) {
	a := NewAssert(t)

	root := newTrie(map[string]*Kind{})
	basic := KindBasic
	_ = root.AddText("any", &basic)
	kind, ok := root.FindText("any")
	a.Require(ok, "find 'any'")
	a.Require(kind != nil && *kind == KindBasic, "kind 'any'", kind.PointerString())

	kind, ok = root.FindText("an")
	a.Require(ok, "find 'an'")
	a.Require(kind == nil, "kind 'an'", kind.PointerString())

	kind, ok = root.FindText("ann")
	a.Require(!ok, "find 'ann'")
	a.Require(kind == nil, "kind 'ann'", kind.PointerString())

	comment := KindComment
	_ = root.AddText("ayy", &comment)

	kind, ok = root.FindText("any")
	a.Require(ok, "find 'any' 2")
	a.Require(kind != nil && *kind == KindBasic, "kind 'any' 2", kind.PointerString())

	kind, ok = root.FindText("ayy")
	a.Require(ok, "find 'ayy'")
	a.Require(kind != nil && *kind == KindComment, "kind 'ayy'", kind.PointerString())

	kind, ok = root.FindByte([]byte("ayy"))
	a.Require(ok, "find 'ayy' 2")
	a.Require(kind != nil && *kind == KindComment, "kind 'ayy' 2", kind.PointerString())
}

func TestTrieComment(t *testing.T) {
	a := NewAssert(t)

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
	a := NewAssert(t)

	tr := newTrie(newCharset("123"))
	_, ok := tr.FindText("123")
	a.Require(ok, "find 123")

	_, ok = tr.FindTextReversely("321")
	a.Require(ok, "find reversely 321")

	_, ok = tr.FindByte([]byte("123"))
	a.Require(ok, "find byte 123")

	_, ok = tr.FindByteReversely([]byte("321"))
	a.Require(ok, "find byte reversely 321")
}
