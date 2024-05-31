package goast

func kindReset(n *Node) {
	_ = _commonResetter.Run(n)
}

type deeperResetter map[Kind]*kindResetter

var (
	_deeperResetterTable = map[*kindResetter]deeperResetter{
		_commonResetter: _commonDeeperResetter,
		_funcResetter:   _funcDeeperResetter,
	}

	_commonDeeperResetter = deeperResetter{
		KindImport: _importResetter,
		KindVar:    _variableResetter,
		KindConst:  _variableResetter,
		KindType:   _typeResetter,
		KindFunc:   _funcResetter,
	}

	_funcDeeperResetter = deeperResetter{
		KindParenthesisLeft: _funcParenthesisResetter,
	}
)

var (
	_commonResetter   = &kindResetter{}
	_importResetter   = &kindResetter{}
	_variableResetter = &kindResetter{}

	_typeResetter = &kindResetter{
		TriggerKind:     newSet(KindType),
		TriggerLimit:    -1,
		KindChangeTable: map[int]Kind{1: KindTypeName},
		ChangeableKind:  newSet(KindRaws),
		ReturnKind:      newSet(KindInterface, KindStruct, KindKeywords),
	}

	_structCurlyBracketResetter = &kindResetter{}

	_interfaceCurlyBracketResetter = &kindResetter{}

	_funcResetter = &kindResetter{
		TriggerKind:     newSet(KindFunc),
		TriggerLimit:    -1,
		KindChangeTable: map[int]Kind{1: KindFuncName},
		ChangeableKind:  newSet(KindRaws),
		ReturnKind:      newSet(KindParenthesisLeft, KindCurlyBracketLeft),
	}

	_funcParenthesisResetter = &kindResetter{
		TriggerKind:     newSet(KindParenthesisLeft),
		TriggerLimit:    -1,
		KindChangeTable: map[int]Kind{1: KindParamName},
		ChangeableKind:  newSet(KindRaws),
		ReturnKind:      newSet(KindParenthesisRight),
	}

	_parenthesisResetter = &kindResetter{}
)

// kindResetter re-set the node kind with node's position of scope
type kindResetter struct {
	TriggerKind     set[Kind]
	TriggerLimit    int
	KindChangeTable map[int]Kind
	ChangeableKind  set[Kind]
	ReturnKind      set[Kind]
}

func (kr *kindResetter) Run(head *Node) *Node {
	var deeperResetterResult *Node

	triggered := false
	triggerIndex := 0
	deeperResetter := _deeperResetterTable[kr]
	triggerLimit := kr.TriggerLimit

	return head.IterNext(func(n *Node) bool {
		if deeperResetterResult != nil {
			if deeperResetterResult == n {
				deeperResetterResult = nil
			}
			return true
		}

		// TODO: Implement me
		//  ReturnKind > deeperResetter > TriggerKind > ChangeableKind > UnchangeableKind
		kind := n.Kind()
		if kr.ReturnKind.Contain(kind) {
			return false
		}

		if resetter, ok := deeperResetter[kind]; ok && resetter != nil {
			deeperResetterResult = resetter.Run(n)
			return true
		}

		if kr.TriggerKind.Contain(kind) {
			triggered = true
			triggerIndex = 0
			return true
		}

		if triggerLimit == 0 || !triggered || !kr.ChangeableKind.Contain(kind) {
			return true
		}

		triggerIndex++
		if triggerLimit > 0 {
			triggerLimit--
		}

		if k, ok := kr.KindChangeTable[triggerIndex]; ok && k != KindNone {
			n.SetKind(k)
		}

		return true
	})
}
