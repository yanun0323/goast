package goast

type Kind uint8

const (
	KindUnknown Kind = iota
	KindComment
	KindKeyword
	KindSymbol
	KindBasicType
	KindSeparator
	KindName /* manual define */
)

func (k *Kind) String() string {
	if k == nil {
		return ""
	}
	switch *k {
	case KindUnknown:
		return "Unknown"
	case KindComment:
		return "Comment"
	case KindKeyword:
		return "Keyword"
	case KindSymbol:
		return "Symbol"
	case KindBasicType:
		return "BasicType"
	case KindSeparator:
		return "Separator"
	case KindName:
		return "Name"
	}

	return ""
}
