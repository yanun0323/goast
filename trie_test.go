package goast

import "testing"

func TestTrie(t *testing.T) {
	a := NewAssert(t)

	root := &trie_[*ElementKind]{}
	basic := ElemBasicType
	_ = root.AddText("any", &basic)
	kind, ok := root.FindText("any")
	a.Require(ok, "find 'any'")
	a.Require(kind != nil && *kind == ElemBasicType, "kind 'any'", kind.String())

	kind, ok = root.FindText("an")
	a.Require(ok, "find 'an'")
	a.Require(kind == nil, "kind 'an'", kind.String())

	kind, ok = root.FindText("ann")
	a.Require(!ok, "find 'ann'")
	a.Require(kind == nil, "kind 'ann'", kind.String())

	comment := ElemComment
	_ = root.AddText("ayy", &comment)

	kind, ok = root.FindText("any")
	a.Require(ok, "find 'any' 2")
	a.Require(kind != nil && *kind == ElemBasicType, "kind 'any' 2", kind.String())

	kind, ok = root.FindText("ayy")
	a.Require(ok, "find 'ayy'")
	a.Require(kind != nil && *kind == ElemComment, "kind 'ayy'", kind.String())

	kind, ok = root.FindByte([]byte("ayy"))
	a.Require(ok, "find 'ayy' 2")
	a.Require(kind != nil && *kind == ElemComment, "kind 'ayy' 2", kind.String())
}
