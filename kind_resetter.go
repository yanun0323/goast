package goast

var (
	_commonResetter = &kindResetter{}
	_importResetter = &kindResetter{}
	_defineResetter = &kindResetter{}
	_typeResetter   = &kindResetter{}
	_funcResetter   = &kindResetter{}

	_commonDeeperResetter = deeperResetter{
		KindImport: _importResetter,
		KindVar:    _defineResetter,
		KindConst:  _defineResetter,
		KindType:   _typeResetter,
		KindFunc:   _funcResetter,
	}
)

type deeperResetter map[Kind]*kindResetter

// kindResetter re-set the node kind with node's position of scope
type kindResetter struct {
	KindChangeTable  map[int]Kind
	ChangeableKind   set[Kind]
	UnchangeableKind set[Kind]
}

func Run(ns []Node, i *int) {
	for _, n := range ns {
		_ = n
	}
}
