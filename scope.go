package goast

import (
	"sort"
	"strings"

	"github.com/yanun0323/goast/kind"
)

/*
Scope means the first level declaration of a go file.

It could be a package/comment/import/var/const/type/func.
*/
type Scope interface {
	Kind() ScopeKind
	Line() int
	Print()
	Node() *Node
}

func NewScope(line int, kind ScopeKind, node *Node) Scope {
	return &scope{
		line: line,
		kind: kind,
		node: node,
	}
}

type scope struct {
	line int
	kind ScopeKind
	node *Node
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

func (d *scope) Node() *Node {
	return d.node
}

func (d *scope) Print() {
	if d == nil {
		return
	}

	println(d.Line()+1, "....", "Scope."+d.kind.String())

	buf := map[int][]string{}
	lines := []int{}
	_ = d.Node().IterNext(func(n *Node) bool {
		if _, ok := buf[n.Line()]; !ok {
			lines = append(lines, n.Line())
		}
		buf[n.Line()] = append(buf[n.Line()], n.Text())

		return true
	})

	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})

	for _, l := range lines {
		println("\t", l+1, "....", strings.Join(buf[l], ""))
	}
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
	default:
		return ""
	}
}

func (k ScopeKind) ToKind() kind.Kind {
	switch k {
	case ScopeUnknown:
		return kind.Raw
	case ScopePackage:
		return kind.Package
	case ScopeComment:
		return kind.Comment
	case ScopeInnerComment:
		return kind.Comment
	case ScopeImport:
		return kind.Import
	case ScopeVariable:
		return kind.Var
	case ScopeConst:
		return kind.Const
	case ScopeType:
		return kind.Type
	case ScopeFunc:
		return kind.Func
	default:
		return kind.None
	}
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
