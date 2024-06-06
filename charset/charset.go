package charset

type Set[T comparable] map[T]struct{}

func (cs Set[T]) Contain(key T) bool {
	if cs == nil {
		return false
	}
	_, ok := cs[key]
	return ok
}

func New[T comparable](chars ...T) Set[T] {
	set := make(Set[T], len(chars))
	for _, char := range chars {
		set[char] = struct{}{}
	}
	return set
}

func NewSlice[T comparable](chars ...[]T) Set[T] {
	set := make(Set[T], len(chars))
	for i := range chars {
		for _, char := range chars[i] {
			set[char] = struct{}{}
		}
	}
	return set
}

var (
	SeparatorCharset = New[byte](
		' ', '[', ']', '(', ')', '{', '}', ',', ':', ';', '\n', '\t', '\r',
	)

	NumberCharset = New[byte](
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
	)

	GolangKeywords = New(
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

	GolangSymbol = NewSlice[byte](
		[]byte("+"),
		[]byte("&"),
		[]byte("+="),
		[]byte("&="),
		[]byte("&&"),
		[]byte("=="),
		[]byte("!="),
		[]byte("("),
		[]byte(")"),
		[]byte("-"),
		[]byte("|"),
		[]byte("-="),
		[]byte("|="),
		[]byte("||"),
		[]byte("<"),
		[]byte("<="),
		[]byte("["),
		[]byte("]"),
		[]byte("*"),
		[]byte("^"),
		[]byte("*="),
		[]byte("^="),
		[]byte("<-"),
		[]byte(">"),
		[]byte(">="),
		[]byte("{"),
		[]byte("}"),
		[]byte("/"),
		[]byte("<<"),
		[]byte("/="),
		[]byte("<<="),
		[]byte("++"),
		[]byte("="),
		[]byte(":="),
		[]byte(","),
		[]byte(";"),
		[]byte("%"),
		[]byte(">>"),
		[]byte("%="),
		[]byte(">>="),
		[]byte("--"),
		[]byte("!"),
		[]byte("..."),
		[]byte("."),
		[]byte(":"),
		[]byte("&^"),
		[]byte("&^="),
		[]byte("~"),
		[]byte("_"),
	)

	GolangBasicType = New(
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
