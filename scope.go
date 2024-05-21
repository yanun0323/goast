package goast

import "strings"

/*
Scope means the first level declaration of a go file.

It could be a package/comment/import/var/const/type/func.
*/
type Scope interface {
	Kind() ScopeKind
	Line() int
	Print()
	Text() string
	Valuable() bool
}

func NewScope(line int, kind ScopeKind, values []rawLine) Scope {
	return &scope{
		line:   line,
		kind:   kind,
		values: values,
	}
}

type scope struct {
	line   int
	kind   ScopeKind
	values []rawLine
}

func (d *scope) Line() int {
	if d == nil {
		return 0
	}

	return d.line
}

func (d *scope) Kind() ScopeKind {
	if d == nil {
		return ScopeUnknown
	}

	return d.kind
}

func (d *scope) Valuable() bool {
	if d == nil {
		return false
	}

	return len(d.values) != 0
}

func (d *scope) Text() string {
	if d == nil {
		return ""
	}

	text := make([]string, 0, len(d.values))
	for _, l := range d.values {
		text = append(text, l.text)
	}

	return strings.Join(text, "\n")
}

func (d *scope) Print() {
	if d == nil {
		return
	}

	println(d.Line(), " ....", "Scope."+d.kind.String())
}

// ScopeKind
type ScopeKind string

const (
	ScopeUnknown      ScopeKind = ""
	ScopePackage      ScopeKind = "package" // package
	ScopeComment      ScopeKind = "//"      // comment
	ScopeInnerComment ScopeKind = "/*"      // inner comment
	ScopeImport       ScopeKind = "import"  // import
	ScopeVariable     ScopeKind = "var"     // var
	ScopeConst        ScopeKind = "const"   // const
	ScopeType         ScopeKind = "type"    // type
	ScopeFunc         ScopeKind = "func"    // func
)

func (k ScopeKind) String() string {
	switch k {
	case ScopeUnknown:
		return "Unknown"
	case ScopePackage:
		return "Package"
	case ScopeComment:
		return "Comment"
	case ScopeInnerComment:
		return "InnerComment"
	case ScopeImport:
		return "Import"
	case ScopeVariable:
		return "Variable"
	case ScopeConst:
		return "Const"
	case ScopeType:
		return "Type"
	case ScopeFunc:
		return "Func"
	}

	return ""
}

func NewScopeKind(s string) ScopeKind {
	kinds := []ScopeKind{
		ScopePackage,
		ScopeComment,
		ScopeInnerComment,
		ScopeImport,
		ScopeVariable,
		ScopeConst,
		ScopeType,
		ScopeFunc,
	}

	for _, kind := range kinds {
		if isScopeKind(s, kind) {
			return kind
		}
	}

	return ScopeUnknown
}

func isScopeKind(s string, k ScopeKind) bool {
	ks := string(k)
	if len(s) < len(ks) {
		return false
	}

	return s[:len(ks)] == ks
}
