package goast

import "strings"

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
	v := strings.TrimSpace(value)
	if len(v) == 0 {
		return nil
	}

	k := ElemUnknown
	if len(kind) != 0 {
		k = kind[0]
	} else {
		k = newElementKind(v)
	}

	return &element{
		line:  line,
		kind:  k,
		value: v,
	}
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

	println("\t", e.Line(), " ....", "Element."+e.kind.String())
}

type ElementKind uint8

const (
	ElemUnknown ElementKind = iota
	ElemComment
	ElemKeyword
	ElemSymbol
	ElemBasic
	ElemName /* manual define */
)

func (k ElementKind) String() string {
	switch k {
	case ElemUnknown:
		return "Unknown"
	case ElemComment:
		return "Comment"
	case ElemKeyword:
		return "Keyword"
	case ElemSymbol:
		return "Symbol"
	case ElemBasic:
		return "Basic"
	case ElemName:
		return "Name"
	}

	return ""
}

var _keywordTable = map[string]struct{}{
	"break":       {},
	"default":     {},
	"func":        {},
	"interface":   {},
	"interface{}": {},
	"select":      {},
	"case":        {},
	"defer":       {},
	"go":          {},
	"map":         {},
	"struct":      {},
	"struct{}":    {},
	"chan":        {},
	"else":        {},
	"goto":        {},
	"package":     {},
	"switch":      {},
	"const":       {},
	"fallthrough": {},
	"if":          {},
	"range":       {},
	"type":        {},
	"continue":    {},
	"for":         {},
	"import":      {},
	"return":      {},
	"var":         {},
}

var _symbolTable = map[string]struct{}{
	"+":   {},
	"&":   {},
	"+=":  {},
	"&=":  {},
	"&&":  {},
	"==":  {},
	"!=":  {},
	"(":   {},
	")":   {},
	"-":   {},
	"|":   {},
	"-=":  {},
	"|=":  {},
	"||":  {},
	"<":   {},
	"<=":  {},
	"[":   {},
	"]":   {},
	"*":   {},
	"^":   {},
	"*=":  {},
	"^=":  {},
	"<-":  {},
	">":   {},
	">=":  {},
	"{":   {},
	"}":   {},
	"/":   {},
	"<<":  {},
	"/=":  {},
	"<<=": {},
	"++":  {},
	"=":   {},
	":=":  {},
	",":   {},
	";":   {},
	"%":   {},
	">>":  {},
	"%=":  {},
	">>=": {},
	"--":  {},
	"!":   {},
	"...": {},
	".":   {},
	":":   {},
	"&^":  {},
	"&^=": {},
	"~":   {},
}

var _basicTable = map[string]struct{}{
	"bool":       {},
	"string":     {},
	"byte":       {},
	"rune":       {},
	"int":        {},
	"uint":       {},
	"int8":       {},
	"uint8":      {},
	"int16":      {},
	"uint16":     {},
	"int32":      {},
	"uint32":     {},
	"int64":      {},
	"uint64":     {},
	"uintptr":    {},
	"float32":    {},
	"float64":    {},
	"complex64":  {},
	"complex128": {},
}

func newElementKind(v string) ElementKind {
	if _, isBasic := _basicTable[v]; isBasic {
		return ElemBasic
	}

	if _, isKeyword := _keywordTable[v]; isKeyword {
		return ElemKeyword
	}

	if _, isSymbol := _symbolTable[v]; isSymbol {
		return ElemSymbol
	}

	if len(v) >= 2 && (v[:2] == "//" || v[:2] == "/*") {
		return ElemComment
	}

	if len(v) >= 4 && v[:4] == "map[" {
		return ElemKeyword
	}

	if len(v) >= 2 && v[:2] == "[]" {
		return ElemKeyword
	}

	if v[0] == '[' {
		return ElemKeyword
	}

	if len(v) >= 6 && v[:6] == "func (" {
		return ElemKeyword
	}

	if len(v) >= 4 && v[:4] == "func" {
		return ElemKeyword
	}

	return ElemUnknown
}
