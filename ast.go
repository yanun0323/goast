package goast

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/yanun0323/goast/helper"
	"github.com/yanun0323/goast/kind"
	"github.com/yanun0323/goast/scope"
)

var (
	ErrNoPackage = errors.New("no package name")
	ErrNotExist  = errors.New("file no exist")
)

type Ast interface {
	Package() (string, error)
	File() string
	Dir() string
	Name() string
	Ext() string
	Scope() []Scope
	IterScope(func(Scope) bool)
	SetScope([]Scope) Ast
	String() string
	Save(string) error
}

func ParseAst(file string) (Ast, error) {
	data, err := helper.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotExist
		}

		return nil, err
	}

	node, err := extract(data)
	if err != nil {
		return nil, err
	}

	scs := []Scope{}
	nodesToReset := []*Node{}

	head := node
	k := scope.Unknown
	line := -1

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

	if _, err := findPackageName(scs); err != nil {
		return nil, err
	}

	result := &ast{
		file:  file,
		dir:   filepath.Dir(file),
		name:  filepath.Base(file),
		ext:   filepath.Ext(file),
		scope: scs,
	}

	return result, nil
}

func findPackageName(scs []Scope) (string, error) {
	var packageScope Scope
	for _, sc := range scs {
		if sc.Kind() == scope.Package {
			packageScope = sc
			break
		}
	}

	if packageScope == nil {
		return "", ErrNoPackage
	}

	packageKeywordAppeared := false
	result := packageScope.Node().IterNext(func(n *Node) bool {
		if packageKeywordAppeared && n.Kind() == kind.Raw {
			return false
		}

		if n.Kind() == kind.Package {
			packageKeywordAppeared = true
		}

		return true
	})

	return result.Text(), nil
}

type ast struct {
	file  string
	dir   string
	name  string
	ext   string
	scope []Scope
}

func NewAst(scope []Scope) (Ast, error) {
	if _, err := findPackageName(scope); err != nil {
		return nil, err
	}

	return &ast{
		scope: scope,
	}, nil
}

func (f *ast) Package() (string, error) {
	return findPackageName(f.scope)
}

func (f *ast) File() string {
	if f == nil {
		return ""
	}

	return f.file
}

func (f *ast) Dir() string {
	if f == nil {
		return ""
	}

	return f.dir
}

func (f *ast) Name() string {
	if f == nil {
		return ""
	}

	return f.name
}

func (f *ast) Ext() string {
	if f == nil {
		return ""
	}

	return f.ext
}

func (f *ast) Scope() []Scope {
	if f == nil {
		return nil
	}

	return f.scope
}

func (f *ast) IterScope(fn func(Scope) bool) {
	for _, sc := range f.scope {
		if !fn(sc) {
			return
		}
	}
}

func (f *ast) SetScope(scope []Scope) Ast {
	return &ast{
		file:  f.file,
		dir:   f.dir,
		name:  f.name,
		ext:   f.ext,
		scope: scope,
	}
}

func (f *ast) String() string {
	buf := strings.Builder{}

	for _, sc := range f.Scope() {
		buf.WriteString(sc.Text())
	}

	return buf.String()
}

func (f *ast) Save(file string) error {
	buf := bytes.Buffer{}

	for _, sc := range f.Scope() {
		buf.WriteString(sc.Text())
	}

	return helper.SaveFile(file, buf.Bytes())
}
