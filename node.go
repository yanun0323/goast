package goast

/*
Node stands for any element in go language.
*/
type Node interface {
	Kind() Kind
	SetKind(Kind)
	Line() int
	Print()
	Text() string
	Valuable() bool

	Prev() Node
	Next() Node
	InsertPrev(Node)
	InsertNext(Node)
	RemovePrev() Node
	RemoveNext() Node
	IterPrev(func(Node) bool) Node
	IterNext(func(Node) bool) Node

	setPrev(Node)
	setNext(Node)
}

func NewNode(line int, text string, kind ...Kind) Node {
	if len(kind) != 0 {
		return &node{line: line, kind: kind[0], text: text}
	}

	return &node{line: line, kind: NewKind(text), text: text}
}

type node struct {
	line int
	kind Kind
	text string

	prev Node
	next Node
}

func (n *node) loop(iter func(Node) Node, fn func(Node) bool) Node {
	for nn := Node(n); nn != nil; nn = iter(nn) {
		if fn(nn) {
			continue
		}

		return nn
	}

	return nil
}

func (n *node) IterPrev(fn func(Node) bool) Node {
	return n.loop(func(n Node) Node { return n.Prev() }, fn)
}

func (n *node) IterNext(fn func(Node) bool) Node {
	return n.loop(func(n Node) Node { return n.Next() }, fn)
}

func (n *node) Prev() Node {
	if n != nil {
		return n.prev
	}

	return nil
}

func (n *node) Next() Node {
	if n != nil {
		return n.next
	}

	return nil
}

func (n *node) setPrev(nn Node) {
	if n != nil {
		n.prev = nn
	}
}

func (n *node) setNext(nn Node) {
	if n != nil {
		n.next = nn
	}
}

func (n *node) InsertPrev(nn Node) {
	n.setPrev(nn)
	nn.setNext(n)
}

func (n *node) InsertNext(nn Node) {
	n.setNext(nn)
	nn.setPrev(n)
}

func (n *node) RemovePrev() Node {
	nn := n.Prev()
	n.setPrev(nil)
	nn.setNext(nil)
	return nn
}

func (n *node) RemoveNext() Node {
	nn := n.Next()
	n.setNext(nil)
	nn.setPrev(nil)
	return nn
}

func (n *node) Line() int {
	if n == nil {
		return 0
	}

	return n.line
}

func (n *node) Kind() Kind {
	if n == nil {
		return KindRaws
	}

	return n.kind
}

func (n *node) SetKind(k Kind) {
	if n != nil {
		n.kind = k
	}
}

func (n *node) Valuable() bool {
	if n == nil {
		return false
	}

	return len(n.text) != 0
}

func (n *node) Text() string {
	if n == nil {
		return ""
	}

	return n.text
}

func (n *node) Print() {
	if n == nil {
		return
	}

	println("\t", n.Line()+1, " ....", "Node."+n.kind.String(), "....", printTidy(n.text))
}
