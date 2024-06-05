package goast

// import (
// 	"github.com/yanun0323/goast/charset"
// 	"github.com/yanun0323/goast/helper"
// )

// type Kind [2]string

// var (
// 	KindNone      Kind
// 	KindRaw       = Kind{"", "Raw"}
// 	KindComment   = Kind{"", "Comment"}
// 	KindSymbol    = Kind{"", "Symbol"}
// 	KindBasic     = Kind{"", "BasicType"}
// 	KindSeparator = Kind{"", "Separator"}
// 	KindKeyword   = Kind{"", "Keyword"}

// 	/* separator */
// 	KindTab                = Kind{"\t", "Tab"}
// 	KindSpace              = Kind{" ", "Space"}
// 	KindComma              = Kind{",", "Comma"}
// 	KindColon              = Kind{":", "Colon"}
// 	KindNewLine            = Kind{"\n", "NewLine"}
// 	KindSemicolon          = Kind{";", "Semicolon"}
// 	KindParenthesisLeft    = Kind{"(", "ParenthesisLeft"}
// 	KindParenthesisRight   = Kind{")", "ParenthesisRight"}
// 	KindCurlyBracketLeft   = Kind{"{", "CurlyBracketLeft"}
// 	KindCurlyBracketRight  = Kind{"}", "CurlyBracketRight"}
// 	KindSquareBracketLeft  = Kind{"[", "SquareBracketLeft"}
// 	KindSquareBracketRight = Kind{"]", "SquareBracketRight"}
// 	/* keyword */
// 	KindImport    = Kind{"import", "Import"}
// 	KindVar       = Kind{"var", "Var"}
// 	KindConst     = Kind{"const", "Const"}
// 	KindMap       = Kind{"map", "Map"}
// 	KindChannel   = Kind{"chan", "Channel"}
// 	KindFunc      = Kind{"func", "Func"}
// 	KindType      = Kind{"type", "Type"}
// 	KindStruct    = Kind{"struct", "Struct"}
// 	KindInterface = Kind{"interface", "Interface"}
// 	KindPackage   = Kind{"package", "Package"}

// 	/* resetter set */
// 	KindFuncName           = Kind{"", "FuncName"}
// 	KindTypeName           = Kind{"", "TypeName"}
// 	KindTypeAliasType      = Kind{"", "TypeAliasType"}
// 	KindParamName          = Kind{"", "ParamName"}
// 	KindParamType          = Kind{"", "ParamType"}
// 	KindMethodReceiverName = Kind{"", "MethodReceiverName"}
// 	KindMethodReceiverType = Kind{"", "MethodReceiverType"}
// 	KindString             = Kind{"", "String"}
// 	KindMethod             = Kind{"", "Method"}
// )

// func (k *Kind) PointerString() string {
// 	if k == nil {
// 		return ""
// 	}

// 	return string(k[1])
// }

// func (k Kind) String() string {
// 	return string(k[1])
// }

// func (k Kind) GetClose() (Kind, bool) {
// 	switch k {
// 	case KindParenthesisLeft:
// 		return KindParenthesisRight, true
// 	case KindCurlyBracketLeft:
// 		return KindCurlyBracketRight, true
// 	case KindSquareBracketLeft:
// 		return KindSquareBracketRight, true
// 	case KindParenthesisRight:
// 		return KindParenthesisLeft, true
// 	case KindCurlyBracketRight:
// 		return KindCurlyBracketLeft, true
// 	case KindSquareBracketRight:
// 		return KindSquareBracketLeft, true
// 	default:
// 		return KindNone, false
// 	}
// }

// func NewKind(s string) Kind {
// 	if len(s) == 0 {
// 		return KindNone
// 	}

// 	switch s {
// 	case KindTab[0]:
// 		return KindTab
// 	case KindSpace[0]:
// 		return KindSpace
// 	case KindComma[0]:
// 		return KindComma
// 	case KindColon[0]:
// 		return KindColon
// 	case KindNewLine[0]:
// 		return KindNewLine
// 	case KindSemicolon[0]:
// 		return KindSemicolon
// 	case KindParenthesisLeft[0]:
// 		return KindParenthesisLeft
// 	case KindParenthesisRight[0]:
// 		return KindParenthesisRight
// 	case KindCurlyBracketLeft[0]:
// 		return KindCurlyBracketLeft
// 	case KindCurlyBracketRight[0]:
// 		return KindCurlyBracketRight
// 	case KindSquareBracketLeft[0]:
// 		return KindSquareBracketLeft
// 	case KindSquareBracketRight[0]:
// 		return KindSquareBracketRight
// 	case KindImport[0]:
// 		return KindImport
// 	case KindVar[0]:
// 		return KindVar
// 	case KindConst[0]:
// 		return KindConst
// 	case KindMap[0]:
// 		return KindMap
// 	case KindChannel[0]:
// 		return KindChannel
// 	case KindFunc[0]:
// 		return KindFunc
// 	case KindType[0]:
// 		return KindType
// 	case KindStruct[0]:
// 		return KindStruct
// 	case KindInterface[0]:
// 		return KindInterface
// 	case KindPackage[0]:
// 		return KindPackage
// 	}

// 	if len(s) == 1 && charset.SeparatorCharset.Contain(s[0]) {
// 		return KindSeparator
// 	}

// 	if buf := []byte(s); helper.HasPrefix(buf, "\"") || helper.HasPrefix(buf, "`") {
// 		return KindString
// 	}

// 	if charset.GolangKeywords.Contain(s) {
// 		return KindKeyword
// 	}

// 	if charset.GolangBasicType.Contain(s) {
// 		return KindBasic
// 	}

// 	if charset.GolangSymbol.Contain(s) {
// 		return KindSymbol
// 	}

// 	return KindRaw
// }
