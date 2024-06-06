package kind

import (
	"github.com/yanun0323/goast/charset"
	"github.com/yanun0323/goast/helper"
)

type Kind [2]string

var (
	None      Kind
	Raw       = Kind{"", "Raw"}
	Comment   = Kind{"", "Comment"}
	Number    = Kind{"", "Number"}
	Symbol    = Kind{"", "Symbol"}
	Basic     = Kind{"", "BasicType"}
	Separator = Kind{"", "Separator"}
	Keyword   = Kind{"", "Keyword"}

	/* separator */
	Tab                = Kind{"\t", "Tab"}
	Space              = Kind{" ", "Space"}
	Comma              = Kind{",", "Comma"}
	Colon              = Kind{":", "Colon"}
	NewLine            = Kind{"\n", "NewLine"}
	Semicolon          = Kind{";", "Semicolon"}
	ParenthesisLeft    = Kind{"(", "ParenthesisLeft"}
	ParenthesisRight   = Kind{")", "ParenthesisRight"}
	CurlyBracketLeft   = Kind{"{", "CurlyBracketLeft"}
	CurlyBracketRight  = Kind{"}", "CurlyBracketRight"}
	SquareBracketLeft  = Kind{"[", "SquareBracketLeft"}
	SquareBracketRight = Kind{"]", "SquareBracketRight"}
	/* keyword */
	Import    = Kind{"import", "Import"}
	Var       = Kind{"var", "Var"}
	Const     = Kind{"const", "Const"}
	Map       = Kind{"map", "Map"}
	Channel   = Kind{"chan", "Channel"}
	Func      = Kind{"func", "Func"}
	Type      = Kind{"type", "Type"}
	Struct    = Kind{"struct", "Struct"}
	Interface = Kind{"interface", "Interface"}
	Package   = Kind{"package", "Package"}

	/* resetter set */
	FuncName           = Kind{"", "FuncName"}
	TypeName           = Kind{"", "TypeName"}
	TypeAliasType      = Kind{"", "TypeAliasType"}
	ParamName          = Kind{"", "ParamName"}
	ParamType          = Kind{"", "ParamType"}
	MethodReceiverName = Kind{"", "MethodReceiverName"}
	MethodReceiverType = Kind{"", "MethodReceiverType"}
	String             = Kind{"", "String"}
	Method             = Kind{"", "Method"}
)

func (k *Kind) PointerString() string {
	if k == nil {
		return ""
	}

	return string(k[1])
}

func (k Kind) String() string {
	return string(k[1])
}

func (k Kind) GetClose() (Kind, bool) {
	switch k {
	case ParenthesisLeft:
		return ParenthesisRight, true
	case CurlyBracketLeft:
		return CurlyBracketRight, true
	case SquareBracketLeft:
		return SquareBracketRight, true
	case ParenthesisRight:
		return ParenthesisLeft, true
	case CurlyBracketRight:
		return CurlyBracketLeft, true
	case SquareBracketRight:
		return SquareBracketLeft, true
	default:
		return None, false
	}
}

func New(s string) Kind {
	if len(s) == 0 {
		return None
	}

	switch s {
	case Tab[0]:
		return Tab
	case Space[0]:
		return Space
	case Comma[0]:
		return Comma
	case Colon[0]:
		return Colon
	case NewLine[0]:
		return NewLine
	case Semicolon[0]:
		return Semicolon
	case ParenthesisLeft[0]:
		return ParenthesisLeft
	case ParenthesisRight[0]:
		return ParenthesisRight
	case CurlyBracketLeft[0]:
		return CurlyBracketLeft
	case CurlyBracketRight[0]:
		return CurlyBracketRight
	case SquareBracketLeft[0]:
		return SquareBracketLeft
	case SquareBracketRight[0]:
		return SquareBracketRight
	case Import[0]:
		return Import
	case Var[0]:
		return Var
	case Const[0]:
		return Const
	case Map[0]:
		return Map
	case Channel[0]:
		return Channel
	case Func[0]:
		return Func
	case Type[0]:
		return Type
	case Struct[0]:
		return Struct
	case Interface[0]:
		return Interface
	case Package[0]:
		return Package
	}

	if len(s) != 0 {
		isNumber := true
		for i := range s {
			isNumber = isNumber && charset.NumberCharset.Contain(s[i])
			if !isNumber {
				break
			}
		}
		if isNumber {
			return Number
		}
	}

	if len(s) == 1 && charset.SeparatorCharset.Contain(s[0]) {
		return Separator
	}

	if buf := []byte(s); helper.HasPrefix(buf, "\"") || helper.HasPrefix(buf, "`") {
		return String
	}

	if charset.GolangKeywords.Contain(s) {
		return Keyword
	}

	if charset.GolangBasicType.Contain(s) {
		return Basic
	}

	if len(s) != 0 {
		isSymbol := true
		for i := range s {
			isSymbol = isSymbol && charset.GolangSymbol.Contain(s[i])
			if !isSymbol {
				break
			}
		}
		if isSymbol {
			return Symbol
		}
	}

	return Raw
}
