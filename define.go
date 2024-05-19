package goast

import "strings"

/*
Define means the first level declaration of a go file.

It could be a package/comment/import/var/const/type/func.
*/
type Define interface {
	Kind() DefineKind
	Line() int
	Print()
	Text() string
	Valuable() bool
}

func NewDefine(line int, kind DefineKind, values []string) Define {
	return &define{
		line:   line,
		kind:   kind,
		values: values,
	}
}

type define struct {
	line   int
	kind   DefineKind
	values []string
}

func (d *define) Line() int {
	if d == nil {
		return 0
	}

	return d.line
}

func (d *define) Kind() DefineKind {
	if d == nil {
		return DefUnknown
	}

	return d.kind
}

func (d *define) Valuable() bool {
	if d == nil {
		return false
	}

	return len(d.values) != 0
}

func (d *define) Text() string {
	if d == nil {
		return ""
	}

	return strings.Join(d.values, "\n")
}

func (d *define) Print() {
	if d == nil {
		return
	}

	println(d.Line(), " ....", "Define."+d.kind.String())
}

// DefineKind
type DefineKind string

const (
	DefUnknown      DefineKind = ""
	DefPackage      DefineKind = "package" // package
	DefComment      DefineKind = "//"      // comment
	DefInnerComment DefineKind = "/*"      // inner comment
	DefImport       DefineKind = "import"  // import
	DefVariable     DefineKind = "var"     // var
	DefConst        DefineKind = "const"   // const
	DefType         DefineKind = "type"    // type
	DefFunc         DefineKind = "func"    // func
)

func (k DefineKind) String() string {
	switch k {
	case DefUnknown:
		return "Unknown"
	case DefPackage:
		return "Package"
	case DefComment:
		return "Comment"
	case DefInnerComment:
		return "InnerComment"
	case DefImport:
		return "Import"
	case DefVariable:
		return "Variable"
	case DefConst:
		return "Const"
	case DefType:
		return "Type"
	case DefFunc:
		return "Func"
	}

	return ""
}

func NewDefineKind(s string) DefineKind {
	switch DefineKind(s) {
	case DefPackage, DefComment, DefInnerComment,
		DefImport, DefVariable, DefConst,
		DefType, DefFunc:
		return DefineKind(s)
	default:
		return DefUnknown
	}
}
