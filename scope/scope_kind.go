package scope

import "github.com/yanun0323/goast/kind"

// Kind
type Kind string

const (
	Unknown      Kind = ""
	Package      Kind = "package" // package
	Comment      Kind = "//"      // comment
	InnerComment Kind = "/*"      // inner comment
	Import       Kind = "import"  // import
	Variable     Kind = "var"     // var
	Const        Kind = "const"   // const
	Type         Kind = "type"    // type
	Func         Kind = "func"    // func
)

func (k Kind) String() string {
	switch k {
	case Unknown:
		return "Unknown"
	case Package:
		return "Package"
	case Comment:
		return "Comment"
	case InnerComment:
		return "InnerComment"
	case Import:
		return "Import"
	case Variable:
		return "Variable"
	case Const:
		return "Const"
	case Type:
		return "Type"
	case Func:
		return "Func"
	default:
		return ""
	}
}

func (k Kind) ToKind() kind.Kind {
	switch k {
	case Unknown:
		return kind.Raw
	case Package:
		return kind.Package
	case Comment:
		return kind.Comment
	case InnerComment:
		return kind.Comment
	case Import:
		return kind.Import
	case Variable:
		return kind.Var
	case Const:
		return kind.Const
	case Type:
		return kind.Type
	case Func:
		return kind.Func
	default:
		return kind.None
	}
}
