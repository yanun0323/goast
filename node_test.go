package goast

import (
	"testing"

	"github.com/yanun0323/goast/assert"
	"github.com/yanun0323/goast/kind"
)

func list(values ...string) (head, tail *Node) {
	if len(values) == 0 {
		return nil, nil
	}

	h := NewNode(0, values[0], kind.Raw)
	cur := h
	t := h.Next()
	for _, v := range values[1:] {
		n := NewNode(0, v, kind.Raw)

		cur.setNext(n)
		n.setPrev(cur)

		cur = n
		t = n
	}

	return h, t
}

type direction uint8

const (
	l direction = iota + 1
	r
)

func AssertNode(a assert.Assert, n *Node, d direction, expected ...string) {
	a.T().Helper()
	iter := func(n *Node) *Node { return n.Prev() }
	if d == r {
		iter = func(n *Node) *Node { return n.Next() }
	}

	i := 0
	for h := n; h != nil; h, i = iter(h), i+1 {
		if i >= len(expected) {
			a.T().Fatalf("%s: extra expected node (%s)", a.T().Name(), h.Text())
		}

		a.Equal(h.Text(), expected[i])
	}

	if i < len(expected) {
		a.T().Fatalf("%s: mismatch node length expected (%s), but got nil", a.T().Name(), expected[i])
	}
}

func TestNilNodeMethod(t *testing.T) {
	a := assert.New(t)

	var n *Node

	a.NoPanic(func() { n.loop(func(nn *Node) *Node { return nn.Next() }, func(n *Node) bool { return true }) })
	a.NoPanic(func() { n.IterPrev(func(n *Node) bool { return true }) })
	a.NoPanic(func() { n.IterNext(func(n *Node) bool { return true }) })
	a.NoPanic(func() { n.Prev() })
	a.NoPanic(func() { n.Next() })
	a.NoPanic(func() { n.InsertPrev(nil) })
	a.NoPanic(func() { n.InsertNext(nil) })
	a.NoPanic(func() { n.ReplacePrev(nil) })
	a.NoPanic(func() { n.ReplaceNext(nil) })
	a.NoPanic(func() { n.TakePrev() })
	a.NoPanic(func() { n.TakeNext() })
	a.NoPanic(func() { n.RemovePrev() })
	a.NoPanic(func() { n.RemoveNext() })
	a.NoPanic(func() { n.Line() })
	a.NoPanic(func() { n.Kind() })
	a.NoPanic(func() { n.SetKind(kind.Raw) })
	a.NoPanic(func() { n.Text() })
	a.NoPanic(func() { n.Print() })
	a.NoPanic(func() { n.setPrev(nil) })
	a.NoPanic(func() { n.setNext(nil) })
}

func TestNodeTake(t *testing.T) {
	a := assert.New(t)

	// TakePrev
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		AssertNode(a, head, r, "1", "2", "3", "4", "5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		took := cur.TakePrev()
		a.Equal(took.Text(), "2")
		a.Nil(took.Prev())
		a.Nil(took.Next())

		AssertNode(a, head, l, "1")
		AssertNode(a, head, r, "1", "3", "4", "5")
	})

	// TakeNext
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		AssertNode(a, head, r, "1", "2", "3", "4", "5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		took := cur.TakeNext()
		a.Equal(took.Text(), "4")
		a.Nil(took.Prev())
		a.Nil(took.Next())

		AssertNode(a, head, l, "1")
		AssertNode(a, head, r, "1", "2", "3", "5")
	})
}

func TestNodeRemove(t *testing.T) {
	a := assert.New(t)

	// RemovePrev
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		AssertNode(a, head, r, "1", "2", "3", "4", "5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		removed := cur.RemovePrev()

		AssertNode(a, removed, l, "2", "1")
		AssertNode(a, removed, r, "2")
		AssertNode(a, cur, l, "3")
		AssertNode(a, cur, r, "3", "4", "5")
		AssertNode(a, head, l, "1")
		AssertNode(a, head, r, "1", "2")
	})

	// RemoveNext
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		AssertNode(a, head, r, "1", "2", "3", "4", "5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		removed := cur.RemoveNext()

		AssertNode(a, removed, l, "4")
		AssertNode(a, removed, r, "4", "5")
		AssertNode(a, cur, l, "3", "2", "1")
		AssertNode(a, cur, r, "3")
		AssertNode(a, head, l, "1")
		AssertNode(a, head, r, "1", "2", "3")
	})
}

func TestNodeInsert(t *testing.T) {
	a := assert.New(t)

	// InsertPrev
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		head2, _ := list("-1", "-2", "-3", "-4", "-5")
		AssertNode(a, head, r, "1", "2", "3", "4", "5")
		AssertNode(a, head2, r, "-1", "-2", "-3", "-4", "-5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		nn := head2.Next().Next()
		a.Equal(nn.Text(), "-3")

		prev, next := cur.InsertPrev(nn)

		AssertNode(a, head, l, "1")
		AssertNode(a, head, r, "1", "2", "-3", "3", "4", "5")

		AssertNode(a, head2, l, "-1")
		AssertNode(a, head2, r, "-1", "-2")

		AssertNode(a, prev, l, "-2", "-1")
		AssertNode(a, prev, r, "-2")

		AssertNode(a, next, l, "-4")
		AssertNode(a, next, r, "-4", "-5")

	})

	// InsertNext
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		head2, _ := list("-1", "-2", "-3", "-4", "-5")
		AssertNode(a, head, r, "1", "2", "3", "4", "5")
		AssertNode(a, head2, r, "-1", "-2", "-3", "-4", "-5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		nn := head2.Next().Next()
		a.Equal(nn.Text(), "-3")

		prev, next := cur.InsertNext(nn)

		AssertNode(a, head, l, "1")
		AssertNode(a, head, r, "1", "2", "3", "-3", "4", "5")

		AssertNode(a, head2, l, "-1")
		AssertNode(a, head2, r, "-1", "-2")

		AssertNode(a, prev, l, "-2", "-1")
		AssertNode(a, prev, r, "-2")

		AssertNode(a, next, l, "-4")
		AssertNode(a, next, r, "-4", "-5")

	})
}

func TestNodeReplace(t *testing.T) {
	a := assert.New(t)

	// ReplacePrev
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		head2, _ := list("-1", "-2", "-3", "-4", "-5")
		AssertNode(a, head, r, "1", "2", "3", "4", "5")
		AssertNode(a, head2, r, "-1", "-2", "-3", "-4", "-5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		nn := head2.Next().Next()
		a.Equal(nn.Text(), "-3")

		prev, next := cur.ReplacePrev(nn)

		AssertNode(a, head, l, "1")
		AssertNode(a, head, r, "1", "2")

		AssertNode(a, head2, l, "-1")
		AssertNode(a, head2, r, "-1", "-2", "-3", "3", "4", "5")

		AssertNode(a, prev, l, "2", "1")
		AssertNode(a, prev, r, "2")

		AssertNode(a, next, l, "-4")
		AssertNode(a, next, r, "-4", "-5")

	})

	// ReplaceNext
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		head2, _ := list("-1", "-2", "-3", "-4", "-5")
		AssertNode(a, head, r, "1", "2", "3", "4", "5")
		AssertNode(a, head2, r, "-1", "-2", "-3", "-4", "-5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		nn := head2.Next().Next()
		a.Equal(nn.Text(), "-3")

		prev, next := cur.ReplaceNext(nn)

		AssertNode(a, head, l, "1")
		AssertNode(a, head, r, "1", "2", "3", "-3", "-4", "-5")

		AssertNode(a, head2, l, "-1")
		AssertNode(a, head2, r, "-1", "-2")

		AssertNode(a, prev, l, "-2", "-1")
		AssertNode(a, prev, r, "-2")

		AssertNode(a, next, l, "4")
		AssertNode(a, next, r, "4", "5")

		n, _ := list("1", "2")
		oP, oN := n.ReplaceNext(n.Next())
		AssertNode(a, n, r, "1", "2")
		a.Nil(oP)
		a.Nil(oN)
	})
}

func TestIsolate(t *testing.T) {
	a := assert.New(t)
	head, tail := list("1", "2", "3")

	mid := head.Next()
	mid.Isolate()

	a.Nil(head.Prev())
	a.Nil(head.Next())

	a.Nil(tail.Prev())
	a.Nil(tail.Next())

	a.Nil(mid.Prev())
	a.Nil(mid.Next())
}

func TestCombine(t *testing.T) {
	a := assert.New(t)

	n1, _ := list("1")
	n2, _ := list("2")
	n3, _ := list("3")
	n4, _ := list("4", "5")
	n5 := n4.Next()

	n1.CombineNext(kind.Raw, n2, n3, n4)

	a.Equal(n1.Text(), "1234")
	AssertNode(a, n1, r, "1234")

	a.Nil(n5.Next())
	a.Nil(n5.Prev())

	var n *Node

	n = n.CombineNext(kind.Raw, n1, n5)
	a.Equal(n.Text(), "12345")
	AssertNode(a, n, r, "12345")
}
