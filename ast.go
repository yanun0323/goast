package goast

import (
	"strings"
)

type Ast interface {
	Defines() []Define
}

func ParseAst(filePath string) (Ast, error) {
	raw, err := newRawFile(filePath)
	if err != nil {
		return nil, err
	}

	ff := &ast{}

	d := &define{
		kind: DefUnknown,
	}

	cache := []string{}

	for i, line := range raw.lines {
		span := strings.Split(line, "\t")[0]
		kind := NewDefineKind(strings.Split(span, " ")[0])
		if kind != DefUnknown {
			if d.Valuable() {
				ff.defines = append(ff.defines, d)
			}
			d = &define{
				line: i + 1,
				kind: kind,
			}
			cache = []string{}
		}
		cache = append(cache, line)
		if len(line) != 0 {
			d.values = append(d.values, cache...)
			cache = []string{}
		}
	}

	if d.Valuable() {
		ff.defines = append(ff.defines, d)
	}

	return ff, nil
}

type ast struct {
	defines []Define
}

func (f *ast) Defines() []Define {
	if f == nil {
		return nil
	}

	return f.defines
}
