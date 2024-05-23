package goast

type Kind [2]string

var (
	KindNone               = Kind{"", ""}
	KindRaw                = Kind{"", "Raw"}
	KindComment            = Kind{"", "Comment"}
	KindSymbol             = Kind{"", "Symbol"}
	KindBasicType          = Kind{"", "BasicType"}
	KindSeparator          = Kind{"", "Separator"} /* separator */
	KindSpace              = Kind{" ", "Space"}
	KindComma              = Kind{",", "Comma"}
	KindColon              = Kind{":", "Colon"}
	KindSemicolon          = Kind{";", "Semicolon"}
	KindParenthesisLeft    = Kind{"(", "ParenthesisLeft"}
	KindParenthesisRight   = Kind{")", "ParenthesisRight"}
	KindCurlyBracketLeft   = Kind{"{", "CurlyBracketLeft"}
	KindCurlyBracketRight  = Kind{"}", "CurlyBracketRight"}
	KindSquareBracketLeft  = Kind{"[", "SquareBracketLeft"}
	KindSquareBracketRight = Kind{"]", "SquareBracketRight"}
	KindKeyword            = Kind{"", "Keyword"} /* keyword */
	KindVar                = Kind{"var", "Var"}
	KindConst              = Kind{"const", "Const"}
	KindMap                = Kind{"map", "Map"}
	KindChannel            = Kind{"chan", "Channel"}
	KindFunc               = Kind{"func", "Func"}
	KindType               = Kind{"type", "Type"}
	KindStruct             = Kind{"struct", "Struct"}
	KindInterface          = Kind{"interface", "Interface"}
	KindPackage            = Kind{"package", "Package"}
	KindName               = Kind{"", "Name"}   /* manual define */
	KindString             = Kind{"", "String"} /* manual define */
	KindMethod             = Kind{"", "Method"} /* manual define */
)

func (k *Kind) String() string {
	if k == nil {
		return ""
	}

	return string(k[1])
}

func NewKind(s string) Kind {
	if len(s) == 0 {
		return KindRaw
	}

	switch s {
	case KindSpace[0]:
		return KindSpace
	case KindComma[0]:
		return KindComma
	case KindColon[0]:
		return KindColon
	case KindSemicolon[0]:
		return KindSemicolon
	case KindParenthesisLeft[0]:
		return KindParenthesisLeft
	case KindParenthesisRight[0]:
		return KindParenthesisRight
	case KindCurlyBracketLeft[0]:
		return KindCurlyBracketLeft
	case KindCurlyBracketRight[0]:
		return KindCurlyBracketRight
	case KindSquareBracketLeft[0]:
		return KindSquareBracketLeft
	case KindSquareBracketRight[0]:
		return KindSquareBracketRight
	case KindVar[0]:
		return KindVar
	case KindConst[0]:
		return KindConst
	case KindMap[0]:
		return KindMap
	case KindChannel[0]:
		return KindChannel
	case KindFunc[0]:
		return KindFunc
	case KindType[0]:
		return KindType
	case KindStruct[0]:
		return KindStruct
	case KindInterface[0]:
		return KindInterface
	case KindPackage[0]:
		return KindPackage
	}

	if len(s) == 1 && _separatorCharset.Contain(s[0]) {
		return KindSeparator
	}

	if buf := []byte(s); hasPrefix(buf, "\"") || hasPrefix(buf, "`") {
		return KindString
	}

	if _golangKeywords.Contain(s) {
		return KindKeyword
	}

	if _golangBasicType.Contain(s) {
		return KindBasicType
	}

	if _golangSymbol.Contain(s) {
		return KindSymbol
	}

	return KindRaw
}
