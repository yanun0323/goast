package goast

import (
	"testing"
)

func list(values ...string) (head, tail *Node) {
	if len(values) == 0 {
		return nil, nil
	}

	h := NewNode(0, values[0], KindRaws)
	cur := h
	t := h.Next()
	for _, v := range values[1:] {
		n := NewNode(0, v, KindRaws)

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

func (a Assert) assertNode(n *Node, d direction, expected ...string) {
	a.t.Helper()
	iter := func(n *Node) *Node { return n.Prev() }
	if d == r {
		iter = func(n *Node) *Node { return n.Next() }
	}

	i := 0
	for h := n; h != nil; h, i = iter(h), i+1 {
		if i >= len(expected) {
			a.t.Fatalf("%s: extra expected node (%s)", a.t.Name(), h.Text())
		}

		a.Equal(h.Text(), expected[i])
	}

	if i < len(expected) {
		a.t.Fatalf("%s: mismatch node length expected (%s), but got nil", a.t.Name(), expected[i])
	}
}

func TestNilNodeMethod(t *testing.T) {
	a := NewAssert(t)

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
	a.NoPanic(func() { n.SetKind(KindRaws) })
	a.NoPanic(func() { n.Valuable() })
	a.NoPanic(func() { n.Text() })
	a.NoPanic(func() { n.Print() })
	a.NoPanic(func() { n.setPrev(nil) })
	a.NoPanic(func() { n.setNext(nil) })
}

func TestNodeTake(t *testing.T) {
	a := NewAssert(t)

	// TakePrev
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		a.assertNode(head, r, "1", "2", "3", "4", "5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		took := cur.TakePrev()
		a.Equal(took.Text(), "2")
		a.Nil(took.Prev())
		a.Nil(took.Next())

		a.assertNode(head, l, "1")
		a.assertNode(head, r, "1", "3", "4", "5")
	})

	// TakeNext
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		a.assertNode(head, r, "1", "2", "3", "4", "5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		took := cur.TakeNext()
		a.Equal(took.Text(), "4")
		a.Nil(took.Prev())
		a.Nil(took.Next())

		a.assertNode(head, l, "1")
		a.assertNode(head, r, "1", "2", "3", "5")
	})
}

func TestNodeRemove(t *testing.T) {
	a := NewAssert(t)

	// RemovePrev
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		a.assertNode(head, r, "1", "2", "3", "4", "5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		removed := cur.RemovePrev()

		a.assertNode(removed, l, "2", "1")
		a.assertNode(removed, r, "2")
		a.assertNode(cur, l, "3")
		a.assertNode(cur, r, "3", "4", "5")
		a.assertNode(head, l, "1")
		a.assertNode(head, r, "1", "2")
	})

	// RemoveNext
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		a.assertNode(head, r, "1", "2", "3", "4", "5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		removed := cur.RemoveNext()

		a.assertNode(removed, l, "4")
		a.assertNode(removed, r, "4", "5")
		a.assertNode(cur, l, "3", "2", "1")
		a.assertNode(cur, r, "3")
		a.assertNode(head, l, "1")
		a.assertNode(head, r, "1", "2", "3")
	})
}

func TestNodeInsert(t *testing.T) {
	a := NewAssert(t)

	// InsertPrev
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		head2, _ := list("-1", "-2", "-3", "-4", "-5")
		a.assertNode(head, r, "1", "2", "3", "4", "5")
		a.assertNode(head2, r, "-1", "-2", "-3", "-4", "-5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		nn := head2.Next().Next()
		a.Equal(nn.Text(), "-3")

		prev, next := cur.InsertPrev(nn)

		a.assertNode(head, l, "1")
		a.assertNode(head, r, "1", "2", "-3", "3", "4", "5")

		a.assertNode(head2, l, "-1")
		a.assertNode(head2, r, "-1", "-2")

		a.assertNode(prev, l, "-2", "-1")
		a.assertNode(prev, r, "-2")

		a.assertNode(next, l, "-4")
		a.assertNode(next, r, "-4", "-5")

	})

	// InsertNext
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		head2, _ := list("-1", "-2", "-3", "-4", "-5")
		a.assertNode(head, r, "1", "2", "3", "4", "5")
		a.assertNode(head2, r, "-1", "-2", "-3", "-4", "-5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		nn := head2.Next().Next()
		a.Equal(nn.Text(), "-3")

		prev, next := cur.InsertNext(nn)

		a.assertNode(head, l, "1")
		a.assertNode(head, r, "1", "2", "3", "-3", "4", "5")

		a.assertNode(head2, l, "-1")
		a.assertNode(head2, r, "-1", "-2")

		a.assertNode(prev, l, "-2", "-1")
		a.assertNode(prev, r, "-2")

		a.assertNode(next, l, "-4")
		a.assertNode(next, r, "-4", "-5")

	})
}

func TestNodeReplace(t *testing.T) {
	a := NewAssert(t)

	// ReplacePrev
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		head2, _ := list("-1", "-2", "-3", "-4", "-5")
		a.assertNode(head, r, "1", "2", "3", "4", "5")
		a.assertNode(head2, r, "-1", "-2", "-3", "-4", "-5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		nn := head2.Next().Next()
		a.Equal(nn.Text(), "-3")

		prev, next := cur.ReplacePrev(nn)

		a.assertNode(head, l, "1")
		a.assertNode(head, r, "1", "2")

		a.assertNode(head2, l, "-1")
		a.assertNode(head2, r, "-1", "-2", "-3", "3", "4", "5")

		a.assertNode(prev, l, "2", "1")
		a.assertNode(prev, r, "2")

		a.assertNode(next, l, "-4")
		a.assertNode(next, r, "-4", "-5")

	})

	// ReplaceNext
	a.NoPanic(func() {
		head, _ := list("1", "2", "3", "4", "5")
		head2, _ := list("-1", "-2", "-3", "-4", "-5")
		a.assertNode(head, r, "1", "2", "3", "4", "5")
		a.assertNode(head2, r, "-1", "-2", "-3", "-4", "-5")

		cur := head.Next().Next()
		a.Equal(cur.Text(), "3")

		nn := head2.Next().Next()
		a.Equal(nn.Text(), "-3")

		prev, next := cur.ReplaceNext(nn)

		a.assertNode(head, l, "1")
		a.assertNode(head, r, "1", "2", "3", "-3", "-4", "-5")

		a.assertNode(head2, l, "-1")
		a.assertNode(head2, r, "-1", "-2")

		a.assertNode(prev, l, "-2", "-1")
		a.assertNode(prev, r, "-2")

		a.assertNode(next, l, "4")
		a.assertNode(next, r, "4", "5")
	})
}
