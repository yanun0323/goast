package goast

type Type uint8

const (
	Unsupported Type = iota
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
	/* array */
	Array
	/* slice */
	Slice
	/* other */
	Function
	Method
	Other
)
