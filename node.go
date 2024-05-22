package goast

/*
Node stands for any element in go language.
*/
type Node interface {
	Kind() Kind
	Line() int
	Print()
	Text() string
	Valuable() bool
}

func NewNode(line int, text string, kind ...Kind) Node {
	elem := &node{line: line, text: text}

	if len(kind) != 0 {
		elem.kind = kind[0]
		return elem
	}

	if len(text) == 1 && _separatorCharset.Contain(text[0]) {
		elem.kind = KindSeparator
		return elem
	}

	if buf := []byte(text); hasPrefix(buf, "\"") || hasPrefix(buf, "`") {
		elem.kind = KindString
		return elem
	}

	if _golangKeywords.Contain(text) {
		elem.kind = KindKeyword
		return elem
	}

	if _golangBasicType.Contain(text) {
		elem.kind = KindBasicType
		return elem
	}

	if _golangSymbol.Contain(text) {
		elem.kind = KindSymbol
		return elem
	}

	elem.kind = KindUnknown
	return elem
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
		return KindUnknown
	}

	return n.kind
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
