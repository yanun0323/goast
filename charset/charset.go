package charset

type charset[T comparable] map[T]struct{}

func (cs charset[T]) Contain(key T) bool {
	if cs == nil {
		return false
	}
	_, ok := cs[key]
	return ok
}

func NewCharset[T comparable](chars ...T) charset[T] {
	set := make(charset[T], len(chars))
	for _, char := range chars {
		set[char] = struct{}{}
	}
	return set
}

var (
	SeparatorCharset = NewCharset[byte](
		' ', '[', ']', '(', ')', '{', '}', ',', ':', ';', '\n', '\t', '\r',
	)

	GolangKeywords = NewCharset(
		"break",
		"default",
		"func",
		"interface",
		"select",
		"case",
		"defer",
		"go",
		"map",
		"struct",
		"chan",
		"else",
		"goto",
		"package",
		"switch",
		"const",
		"fallthrough",
		"if",
		"range",
		"type",
		"continue",
		"for",
		"import",
		"return",
		"var",
	)

	GolangSymbol = NewCharset(
		"+",
		"&",
		"+=",
		"&=",
		"&&",
		"==",
		"!=",
		"(",
		")",
		"-",
		"|",
		"-=",
		"|=",
		"||",
		"<",
		"<=",
		"[",
		"]",
		"*",
		"^",
		"*=",
		"^=",
		"<-",
		">",
		">=",
		"{",
		"}",
		"/",
		"<<",
		"/=",
		"<<=",
		"++",
		"=",
		":=",
		",",
		";",
		"%",
		">>",
		"%=",
		">>=",
		"--",
		"!",
		"...",
		".",
		":",
		"&^",
		"&^=",
		"~",
	)

	GolangBasicType = NewCharset(
		"bool",
		"string",
		"byte",
		"rune",
		"int",
		"uint",
		"int8",
		"uint8",
		"int16",
		"uint16",
		"int32",
		"uint32",
		"int64",
		"uint64",
		"uintptr",
		"float32",
		"float64",
		"complex64",
		"complex128",
	)
)
