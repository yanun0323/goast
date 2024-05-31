package goast

import "testing"

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

func TestNodeTakePrev(t *testing.T) {
	a := NewAssert(t)

	a.NoPanic(func() {
		cur, _ := list("1", "2", "3", "4", "5")
		cur = cur.Next().Next()

		a.Require(cur.Text() == "3", "center node should be 3")

		took := cur.TakePrev()
		a.Equal(took.Text(), "2")
		a.Nil(took.Prev())
		a.Nil(took.Next())

		a.Equal(cur.Prev().Text(), "1")
		a.Equal(cur.Prev().Next().Text(), "3")
	})
}

func TestNodeTakeNext(t *testing.T) {
	a := NewAssert(t)

	a.NoPanic(func() {
		cur, _ := list("1", "2", "3", "4", "5")
		cur = cur.Next().Next()

		a.Require(cur.Text() == "3", "center node should be 3")

		took := cur.TakeNext()
		a.Equal(took.Text(), "4")
		a.Nil(took.Prev())
		a.Nil(took.Next())

		a.Equal(cur.Next().Text(), "5")
		a.Equal(cur.Next().Prev().Text(), "3")
	})
}
