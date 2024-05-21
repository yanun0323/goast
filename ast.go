package goast

import "path/filepath"

type Ast interface {
	File() string
	Dir() string
	Name() string
	Ext() string
	Scope() []Scope
}

func ParseAst(file string) (Ast, error) {
	data, err := readFile(file)
	if err != nil {
		return nil, err
	}

	node, err := extract(data)
	if err != nil {
		return nil, err
	}

	sc := []Scope{}

	p := 0
	k := ScopeUnknown
	line := -1
	for i, n := range node {
		if n.Line() == line {
			continue
		}
		line = n.Line()
		nk := NewScopeKind(n.Text())
		if nk == ScopeUnknown {
			continue
		}

		if k != ScopeUnknown {
			sc = append(sc, NewScope(
				node[p].Line(),
				k,
				node[p:i],
			))
		}
		p = i
		k = nk
	}

	result := &ast{
		file:  file,
		dir:   filepath.Dir(file),
		name:  filepath.Base(file),
		ext:   filepath.Ext(file),
		scope: sc,
	}

	return result, nil
}

type ast struct {
	file  string
	dir   string
	name  string
	ext   string
	scope []Scope
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
