package goast

import (
	"sort"
	"strings"

	"github.com/yanun0323/goast/kind"
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
	Text() string

	GetTypeName() (string, bool)
	GetStructName() (string, bool)
	GetInterfaceName() (string, bool)
	GetFuncName() (string, bool)
	GetMethodName() (string, bool)
	GetMethodReceiver() (string, bool)
	Copy() Scope
}

func NewScope(line int, kind scope.Kind, node *Node) Scope {
	return &scopeStruct{
		line: line,
		kind: kind,
		node: node,
	}
}

func ParseScope(startLine int, text []byte) ([]Scope, error) {
	node, err := extract(text)
	if err != nil {
		return nil, err
	}

	scs := []Scope{}
	nodesToReset := []*Node{}

	head := node
	k := scope.Unknown
	line := startLine - 1

	tryAppendScope := func() {
		if k != scope.Unknown {
			nodesToReset = append(nodesToReset, head)
			scs = append(scs, NewScope(
				head.Line(),
				k,
				head,
			))
		}

		if head.Prev() != nil {
			_ = head.RemovePrev()
		}
	}

	_ = node.IterNext(func(n *Node) bool {
		if n.Line() == line {
			return true
		}

		line = n.Line()
		nk := newScopeKind(n.Text())
		if nk == scope.Unknown {
			return true
		}

		tryAppendScope()

		head = n
		k = nk
		return true
	})

	tryAppendScope()

	for _, n := range nodesToReset {
		resetKind(n)
	}

	return scs, nil
}

type scopeStruct struct {
	line int
	kind scope.Kind
	node *Node
}

func (d *scopeStruct) Line() int {
	if d == nil {
		return -1
	}

	return d.line + 1
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

func (d *scopeStruct) Text() string {
	buf := strings.Builder{}

	d.Node().IterNext(func(n *Node) bool {
		buf.WriteString(n.Text())
		return true
	})

	return buf.String()
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

func (d *scopeStruct) GetFuncName() (string, bool) {
	if d.Kind() != scope.Func {
		return "", false
	}

	var (
		name  string
		found bool
	)

	d.Node().IterNext(func(n *Node) bool {
		if n.Kind() == kind.FuncName {
			name = n.Text()
			found = true
			return false
		}

		return true
	})

	return name, found
}

func (d *scopeStruct) GetTypeName() (string, bool) {
	if d.Kind() != scope.Type {
		return "", false
	}

	var (
		name  string
		found bool
	)

	d.Node().IterNext(func(n *Node) bool {
		found = n.Kind() == kind.TypeName
		name = n.Text()
		return !found
	})

	return name, found
}

func (d *scopeStruct) GetStructName() (string, bool) {
	if d.Kind() != scope.Type {
		return "", false
	}

	isStruct := false

	d.Node().IterNext(func(n *Node) bool {
		isStruct = n.Kind() == kind.Struct
		return !isStruct
	})

	return d.GetTypeName()
}

func (d *scopeStruct) GetInterfaceName() (string, bool) {
	if d.Kind() != scope.Type {
		return "", false
	}

	isInterface := false

	d.Node().IterNext(func(n *Node) bool {
		isInterface = n.Kind() == kind.Interface
		return !isInterface
	})

	return d.GetTypeName()
}

func (d *scopeStruct) GetMethodName() (string, bool) {
	if d.Kind() != scope.Func {
		return "", false
	}

	var (
		name  string
		found bool
	)

	if _, isMethod := d.GetMethodReceiver(); !isMethod {
		return "", false
	}

	d.Node().IterNext(func(n *Node) bool {
		found = n.Kind() == kind.FuncName
		name = n.Text()
		return !found
	})

	return name, found
}

func (d *scopeStruct) GetMethodReceiver() (string, bool) {
	if d.Kind() != scope.Func {
		return "", false
	}

	var (
		name  string
		found bool
	)

	d.Node().IterNext(func(n *Node) bool {
		found = n.Kind() == kind.MethodReceiverType
		name = n.Text()
		return !found
	})

	return name, found
}

func (d *scopeStruct) Copy() Scope {
	return &scopeStruct{
		line: d.line,
		kind: d.kind,
		node: d.node.Copy(true),
	}
}
