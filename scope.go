package goast

import (
	"sort"
	"strings"

	"github.com/yanun0323/goast/scope"
)

/*
Scope means the first level declaration of a go file.

It could be a package/comment/import/var/const/type/func.
*/
type Scope interface {
	Kind() scope.Kind
	Line() int
	Print()
	Node() *Node
}

func NewScope(line int, kind scope.Kind, node *Node) Scope {
	return &scopeStruct{
		line: line,
		kind: kind,
		node: node,
	}
}

type scopeStruct struct {
	line int
	kind scope.Kind
	node *Node
}

func (d *scopeStruct) Line() int {
	if d == nil {
		return 0
	}

	return d.line
}

func (d *scopeStruct) Kind() scope.Kind {
	if d == nil {
		return scope.Unknown
	}

	return d.kind
}

func (d *scopeStruct) Node() *Node {
	return d.node
}

func (d *scopeStruct) Print() {
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

func newScopeKind(s string) scope.Kind {
	kinds := []scope.Kind{
		scope.Package,
		scope.Comment,
		scope.InnerComment,
		scope.Import,
		scope.Variable,
		scope.Const,
		scope.Type,
		scope.Func,
	}

	for _, kind := range kinds {
		if isScopeKind(s, kind) {
			return kind
		}
	}

	return scope.Unknown
}

func isScopeKind(s string, k scope.Kind) bool {
	ks := string(k)
	if len(s) < len(ks) {
		return false
	}

	return s[:len(ks)] == ks
}
