package goast

/*
Element stands for any element in go language.
*/
type Element interface {
	Kind() ElementKind
	Line() int
	Print()
	Text() string
	Valuable() bool
}

func NewElement(line int, value string, kind ...ElementKind) Element {
	elem := &element{line: line, value: value}

	if len(kind) != 0 {
		elem.kind = kind[0]
		return elem
	}

	if len(value) == 1 && _separatorCharset.Contain(value[0]) {
		elem.kind = ElemSeparator
		return elem
	}

	if _golangKeywords.Contain(value) {
		elem.kind = ElemKeyword
		return elem
	}

	if _golangBasicType.Contain(value) {
		elem.kind = ElemBasicType
		return elem
	}

	if _golangSymbol.Contain(value) {
		elem.kind = ElemSymbol
		return elem
	}

	elem.kind = ElemUnknown
	return elem
}

type element struct {
	line  int
	kind  ElementKind
	value string
}

func (e *element) Line() int {
	if e == nil {
		return 0
	}

	return e.line
}

func (e *element) Kind() ElementKind {
	if e == nil {
		return ElemUnknown
	}

	return e.kind
}

func (e *element) Valuable() bool {
	if e == nil {
		return false
	}

	return len(e.value) != 0
}

func (e *element) Text() string {
	if e == nil {
		return ""
	}

	return e.value
}

func (e *element) Print() {
	if e == nil {
		return
	}

	println("\t", e.Line(), " ....", "Element."+e.kind.String(), "....", e.value)
}

type ElementKind uint8

const (
	ElemUnknown ElementKind = iota
	ElemComment
	ElemKeyword
	ElemSymbol
	ElemBasicType
	ElemSeparator
	ElemName /* manual define */
)

func (k *ElementKind) String() string {
	if k == nil {
		return ""
	}
	switch *k {
	case ElemUnknown:
		return "Unknown"
	case ElemComment:
		return "Comment"
	case ElemKeyword:
		return "Keyword"
	case ElemSymbol:
		return "Symbol"
	case ElemBasicType:
		return "BasicType"
	case ElemSeparator:
		return "Separator"
	case ElemName:
		return "Name"
	}

	return ""
}
