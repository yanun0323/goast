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
		cur = n
		t = n
	}

	return h, t
}

func TestNodeNilMethod(t *testing.T) {
	a := NewAssert(t)

	var n *Node

	a.NoPanic(func() {
		n.loop(func(nn *Node) *Node { return nn.Next() }, func(n *Node) bool { return true })
		n.IterPrev(func(n *Node) bool { return true })
		n.IterNext(func(n *Node) bool { return true })
		n.Prev()
		n.Next()
		n.InsertPrev(nil)
		n.InsertNext(nil)
		n.ReplacePrev(nil)
		n.ReplaceNext(nil)
		n.TakePrev()
		n.TakeNext()
		n.RemovePrev()
		n.RemoveNext()
		n.Line()
		n.Kind()
		n.SetKind(KindRaws)
		n.Valuable()
		n.Text()
		n.Print()
		n.setPrev(nil)
		n.setNext(nil)
	})
}

func TestNodeTakePrev(t *testing.T) {
	a := NewAssert(t)

	cur, _ := list("1", "2", "3", "4", "5")
	cur = cur.Next().Next()

	a.Require(cur.Text() == "3", "center node should be 3")
}
