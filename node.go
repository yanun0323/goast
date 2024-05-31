package goast

/*
Node stands for any element in go language.
*/

func NewNode(line int, text string, kind ...Kind) *Node {
	if len(kind) != 0 {
		return &Node{line: line, kind: kind[0], text: text}
	}

	return &Node{line: line, kind: NewKind(text), text: text}
}

type Node struct {
	line int
	kind Kind
	text string

	prev *Node
	next *Node
}

func (n *Node) loop(iter func(*Node) *Node, fn func(*Node) bool) *Node {
	for nn := n; nn != nil; nn = iter(nn) {
		if fn(nn) {
			continue
		}

		return nn
	}

	return nil
}

func (n *Node) IterPrev(fn func(*Node) bool) *Node {
	return n.loop(func(nn *Node) *Node { return nn.Prev() }, fn)
}

func (n *Node) IterNext(fn func(*Node) bool) *Node {
	return n.loop(func(nn *Node) *Node { return nn.Next() }, fn)
}

func (n *Node) Prev() *Node {
	if n != nil {
		return n.prev
	}

	return nil
}

func (n *Node) Next() *Node {
	if n != nil {
		return n.next
	}

	return nil
}

// InsertPrev inserts this node into current node's before,
// then returns old previous/next node of inserted node.
func (n *Node) InsertPrev(nn *Node) (insertedOldPrev *Node, insertedOldNext *Node) {
	oldPrev, oldNext := n.Prev(), nn.Next()
	nPrev := n.Prev()

	n.setPrev(nn)
	nn.setNext(n)

	nn.setPrev(nPrev)
	nPrev.setNext(nn)

	oldPrev.setNext(nil)
	oldNext.setPrev(nil)

	return oldPrev, oldNext
}

// InsertNext inserts this node into current node's after,
// then returns old previous/next node of inserted node.
func (n *Node) InsertNext(nn *Node) (insertedOldPrev *Node, insertedOldNext *Node) {
	oldPrev, oldNext := nn.Prev(), nn.Next()
	nNext := n.Next()

	n.setNext(nn)
	nn.setPrev(n)

	nn.setNext(nNext)
	nNext.setPrev(nn)

	oldPrev.setNext(nil)
	oldNext.setPrev(nil)

	return oldPrev, oldNext
}

// ReplacePrev replaces this node into current node's before,
// then returns the old previous node of current node and
// the old next node of replaced node.
func (n *Node) ReplacePrev(nn *Node) (currentOldPrev *Node, replacedOldNext *Node) {
	oldPrev, oldNext := n.Prev(), nn.Next()

	n.setPrev(nn)
	nn.setNext(n)

	oldPrev.setNext(nil)
	oldNext.setPrev(nil)

	return oldPrev, oldNext
}

// ReplaceNext replaces this node into current node's after,
// then returns the old previous node of replaced node and
// the old next node of current node.
func (n *Node) ReplaceNext(nn *Node) (replacedOldPrev *Node, currentOldNext *Node) {
	oldPrev, oldNext := nn.Prev(), n.Next()

	n.setNext(nn)
	nn.setPrev(n)

	oldPrev.setNext(nil)
	oldNext.setPrev(nil)

	return oldPrev, oldNext
}

// TakePrev takes the single node after current node, and connects leftover node.
func (n *Node) TakePrev() *Node {
	removed := n.Prev()

	n.setPrev(removed.Prev())
	n.Prev().setNext(n)

	removed.setNext(nil)
	removed.setPrev(nil)

	return removed
}

// TakeNext takes the single node before current node, and connects leftover node.
func (n *Node) TakeNext() *Node {
	removed := n.Next()

	n.setNext(removed.Next())
	n.Next().setPrev(n)

	removed.setNext(nil)
	removed.setPrev(nil)

	return removed
}

// RemovePrev removes node after current node, and returns removed node with it's all previous node.
func (n *Node) RemovePrev() *Node {
	removed := n.Prev()
	n.setPrev(nil)
	removed.setNext(nil)

	return removed
}

// RemoveNext removes node next current node, and returns removed node with it's all next node.
func (n *Node) RemoveNext() *Node {
	removed := n.Next()
	n.setNext(nil)
	removed.setPrev(nil)

	return removed
}

func (n *Node) Line() int {
	if n == nil {
		return 0
	}

	return n.line
}

func (n *Node) Kind() Kind {
	if n == nil {
		return KindRaws
	}

	return n.kind
}

func (n *Node) SetKind(k Kind) {
	if n != nil {
		n.kind = k
	}
}

func (n *Node) Valuable() bool {
	if n == nil {
		return false
	}

	return len(n.text) != 0
}

func (n *Node) Text() string {
	if n == nil {
		return ""
	}
	return n.text
}

func (n *Node) Print() {
	if n == nil {
		return
	}

	println("\t", n.Line()+1, " ....", "*Node."+n.kind.PointerString(), "....", printTidy(n.text))
}

func (n *Node) setPrev(nn *Node) {
	if n != nil {
		n.prev = nn
	}
}

func (n *Node) setNext(nn *Node) {
	if n != nil {
		n.next = nn
	}
}
