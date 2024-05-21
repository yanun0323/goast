package goast

/*
Scope means the first level declaration of a go file.

It could be a package/comment/import/var/const/type/func.
*/
type Scope interface {
	Kind() ScopeKind
	Line() int
	Print()
	Node() []Node
}

func NewScope(line int, kind ScopeKind, node []Node) Scope {
	return &scope{
		line: line,
		kind: kind,
		node: node,
	}
}

type scope struct {
	line int
	kind ScopeKind
	node []Node
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

func (d *scope) Node() []Node {
	return d.node
}

func (d *scope) Print() {
	if d == nil {
		return
	}

	println(d.Line()+1, " ....", "Scope."+d.kind.String())
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
