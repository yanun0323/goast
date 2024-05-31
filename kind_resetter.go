package goast

func kindReset(n Node) {
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
		KindChangeTable:  map[int]Kind{1: KindTypeName},
		ChangeableKind:   newSet(KindRaws),
		UnchangeableKind: nil,
		ReturnKind:       newSet(KindInterface, KindStruct, KindKeywords),
	}

	_structCurlyBracketResetter = &kindResetter{}

	_interfaceCurlyBracketResetter = &kindResetter{}

	_funcResetter = &kindResetter{
		TriggerKind:      newSet(KindFunc),
		KindChangeTable:  map[int]Kind{1: KindFuncName},
		ChangeableKind:   newSet(KindRaws),
		UnchangeableKind: nil,
		ReturnKind:       newSet(KindParenthesisLeft, KindCurlyBracketLeft),
	}

	_funcParenthesisResetter = &kindResetter{
		TriggerKind:      newSet(KindParenthesisLeft),
		KindChangeTable:  map[int]Kind{1: KindParamName},
		ChangeableKind:   newSet(KindRaws),
		UnchangeableKind: nil,
		ReturnKind:       newSet(KindParenthesisRight),
	}

	_parenthesisResetter = &kindResetter{}
)

// kindResetter re-set the node kind with node's position of scope
type kindResetter struct {
	TriggerKind      set[Kind]
	KindChangeTable  map[int]Kind
	ChangeableKind   set[Kind]
	UnchangeableKind set[Kind]
	ReturnKind       set[Kind]
}

func (kr *kindResetter) Run(head Node) Node {
	var deeperResetterResultNode Node

	triggered := false
	triggerIndex := 0
	deeperResetter := _deeperResetterTable[kr]

	return head.IterNext(func(n Node) bool {
		if deeperResetterResultNode != nil {
			if deeperResetterResultNode == n {
				deeperResetterResultNode = nil
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
			deeperResetterResultNode = resetter.Run(n)
			return true
		}

		if kr.TriggerKind.Contain(kind) {
			triggered = true
			triggerIndex = 0
			return true
		}

		if !triggered {
			return true
		}

		if kr.UnchangeableKind.Contain(kind) || !kr.ChangeableKind.Contain(kind) {
			return true
		}

		triggerIndex++

		if k, ok := kr.KindChangeTable[triggerIndex]; ok && k != KindNone {
			n.SetKind(k)
		}

		return true
	})
}
