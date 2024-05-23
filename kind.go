package goast

type Kind [2]string

var (
	KindNone       Kind
	KindRaws       = Kind{"", "Raw"}
	KindComments   = Kind{"", "Comment"}
	KindSymbols    = Kind{"", "Symbol"}
	KindBasics     = Kind{"", "BasicType"}
	KindSeparators = Kind{"", "Separator"}
	KindKeywords   = Kind{"", "Keyword"}

	/* separator */
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

	/* keyword */
	KindImport    = Kind{"import", "Import"}
	KindVar       = Kind{"var", "Var"}
	KindConst     = Kind{"const", "Const"}
	KindMap       = Kind{"map", "Map"}
	KindChannel   = Kind{"chan", "Channel"}
	KindFunc      = Kind{"func", "Func"}
	KindType      = Kind{"type", "Type"}
	KindStruct    = Kind{"struct", "Struct"}
	KindInterface = Kind{"interface", "Interface"}
	KindPackage   = Kind{"package", "Package"}

	/* resetter set */
	KindFuncName  = Kind{"", "FuncName"}
	KindTypeName  = Kind{"", "TypeName"}
	KindParamName = Kind{"", "ParamName"}
	KindString    = Kind{"", "String"}
	KindMethod    = Kind{"", "Method"}
)

func (k *Kind) String() string {
	if k == nil {
		return ""
	}

	return string(k[1])
}

func NewKind(s string) Kind {
	if len(s) == 0 {
		return KindNone
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
	case KindImport[0]:
		return KindImport
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
		return KindSeparators
	}

	if buf := []byte(s); hasPrefix(buf, "\"") || hasPrefix(buf, "`") {
		return KindString
	}

	if _golangKeywords.Contain(s) {
		return KindKeywords
	}

	if _golangBasicType.Contain(s) {
		return KindBasics
	}

	if _golangSymbol.Contain(s) {
		return KindSymbols
	}

	return KindRaws
}
