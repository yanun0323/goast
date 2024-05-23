package goast

var (
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
	TriggerKind      map[Kind]bool
	KindChangeTable  map[int]Kind
	ChangeableKind   map[Kind]bool
	UnchangeableKind map[Kind]bool
}

func Run(ns []Node) {
	for _, n := range ns {
		_ = n
	}
}
