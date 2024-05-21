package goast

type Ast interface {
	Defines() []Scope
}

func ParseAst(filePath string) (Ast, error) {
	rf, err := newRawFile(filePath)
	if err != nil {
		return nil, err
	}

	ff := &ast{}

	_ = rf

	return ff, nil
}

type ast struct {
	scopes []Scope
}

func (f *ast) Defines() []Scope {
	if f == nil {
		return nil
	}

	return f.scopes
}
