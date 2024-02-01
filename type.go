package goast

import "strings"

type Type uint8

const (
	Raw Type = iota
	Keyword
	Comment
	Structure
	Interface
	Bool
	Map
	/* string */
	String
	Byte
	Rune
	/* int */
	Int
	UInt
	Int8
	UInt8
	Int16
	UInt16
	Int32
	UInt32
	Int64
	UInt64
	UIntPtr
	/* float */
	Float32
	Float64
	/* complex */
	Complex64
	Complex128
	/* slice */
	Slice
	/* array */
	Array
	/* function */
	Function
	Method
)

func parsingType(span string, extraRule ...func(string) (Type, bool)) Type {
	for _, rule := range extraRule {
		t, ok := rule(span)
		if ok {
			return t
		}
	}

	if len(span) == 0 {
		return Raw
	}

	switch span {
	case "type", "var", "const", "import", "package", "=", "{", "}", "(", ")":
		return Keyword
	case "struct":
		return Structure
	case "interface":
		return Interface
	case "bool":
		return Bool
	case "string":
		return String
	case "byte":
		return Byte
	case "rune":
		return Rune
	case "int":
		return Int
	case "uint":
		return UInt
	case "int8":
		return Int8
	case "uint8":
		return UInt8
	case "int16":
		return Int16
	case "uint16":
		return UInt16
	case "int32":
		return Int32
	case "uint32":
		return UInt32
	case "int64":
		return Int64
	case "uint64":
		return UInt64
	case "uintptr":
		return UIntPtr
	case "float32":
		return Float32
	case "float64":
		return Float64
	case "complex64":
		return Complex64
	case "complex128":
		return Complex128
	}

	if strings.HasPrefix(span, "//") || strings.HasPrefix(span, "/*") {
		return Comment
	}

	if strings.HasPrefix(span, "map") {
		return Map
	}

	if strings.HasPrefix(span, "[]") {
		return Slice
	}

	if strings.HasPrefix(span, "[") {
		return Array
	}

	if strings.HasPrefix(span, "func (") {
		return Method
	}

	if strings.HasPrefix(span, "func") {
		return Function
	}

	return Raw
}
