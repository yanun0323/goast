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
}

func (n *node) Line() int {
	if n == nil {
		return 0
	}

	return n.line
}

func (n *node) Kind() Kind {
	if n == nil {
		return KindRaw
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
