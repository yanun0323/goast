package goast

import (
	"path/filepath"

	"github.com/yanun0323/goast/helper"
	"github.com/yanun0323/goast/kind"
	"github.com/yanun0323/goast/scope"
)

type Ast interface {
	Package() string
	File() string
	Dir() string
	Name() string
	Ext() string
	Scope() []Scope
}

func ParseAst(file string) (Ast, error) {
	data, err := helper.ReadFile(file)
	if err != nil {
		return nil, err
	}

	node, err := extract(data)
	if err != nil {
		return nil, err
	}

	_ = generalResetter().Run(node)

	scs := []Scope{}

	head := node
	k := scope.Unknown
	line := -1

	tryAppendScope := func() {
		if k != scope.Unknown {
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

	result := &ast{
		pkg:   findPackageName(scs),
		file:  file,
		dir:   filepath.Dir(file),
		name:  filepath.Base(file),
		ext:   filepath.Ext(file),
		scope: scs,
	}

	return result, nil
}

func findPackageName(scs []Scope) string {
	var packageScope Scope
	for _, sc := range scs {
		if sc.Kind() == scope.Package {
			packageScope = sc
			break
		}
	}

	if packageScope == nil {
		return ""
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

	return result.Text()
}

type ast struct {
	pkg   string
	file  string
	dir   string
	name  string
	ext   string
	scope []Scope
}

func (f *ast) Package() string {
	if f == nil {
		return ""
	}

	return f.pkg
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
