package goast

import (
	"fmt"

	"github.com/yanun0323/goast/kind"
	"github.com/yanun0323/goast/scope"
)

/*
Scope means the first level declaration of a go file.

It could be a package/comment/import/var/const/type/func.
*/
type Scope interface {
	// Kind returns the kind of the scope.
	Kind() scope.Kind

	// Line returns the line number of the scope.
	Line() int

	// Description returns a description of the scope.
	Description() string

	// Print prints the description of the scope.
	Print()

	// Node returns the root node of the scope.
	Node() *Node

	// GetTypeName returns the type name of the scope if the scope is a type.
	GetTypeName() (string, bool)

	// GetStructName returns the struct name of the scope if the scope is a struct.
	GetStructName() (string, bool)

	// GetInterfaceName returns the interface name of the scope if the scope is a interface.
	GetInterfaceName() (string, bool)

	// GetFuncName returns the func name of the scope if the scope is a func.
	GetFuncName() (string, bool)

	// GetMethodName returns the method name of the scope if the scope is a method.
	GetMethodName() (string, bool)

	// GetMethodReceiver returns the method receiver of the scope if the scope is a method.
	GetMethodReceiver() (string, bool)

	// Copy returns a copy of the scope.
	Copy() Scope
}

// NewScope creates a new Scope.
func NewScope(line int, kind scope.Kind, node *Node) Scope {
	return &scopeStruct{
		line: line,
		kind: kind,
		node: node,
	}
}

// ParseScope parses the given text into scopes.
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

func (d *scopeStruct) Description() string {
	if d == nil {
		return ""
	}

	return fmt.Sprintf("%d .... *Scope.%s", d.Line()+1, d.kind.String())
}

func (d *scopeStruct) Print() {
	println(d.Description())
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
