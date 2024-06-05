package goast

func kindReset(n *Node) *Node {
	return generalResetter().Run(n)
}

type kindResetter interface {
	Run(*Node, ...func(*Node)) *Node
}

func handleHook(n *Node, hooks ...func(*Node)) {
	for _, hook := range hooks {
		hook(n)
	}
}

func generalResetter() kindResetter {
	return &genericResetter{
		DeeperResetterTable: map[Kind]kindResetter{
			KindFunc: funcResetter{isFuncKeywordLeading: true},
		},
	}
}

// genericResetter stands for COMMON genericResetter
type genericResetter struct {
	TriggerKind         set[Kind]
	TriggerLimit        int
	KindChangeTable     map[int]Kind
	ChangeableKind      set[Kind]
	ReturnKind          set[Kind]
	DeeperResetterTable map[Kind]kindResetter
}

func (r genericResetter) Run(head *Node, _ ...func(*Node)) *Node {
	var (
		skipAll      bool
		jumpTo       *Node
		triggered    bool
		triggerIndex int

		triggerLimit = r.TriggerLimit
	)

	return head.IterNext(func(n *Node) bool {
		if skipAll {
			return true
		}

		if jumpTo != nil {
			if jumpTo != n {
				return true
			}
			jumpTo = nil
		}

		//  ReturnKind > deeperResetter > TriggerKind > ChangeableKind > UnchangeableKind
		kind := n.Kind()
		if r.ReturnKind.Contain(kind) {
			return false
		}

		if resetter, ok := r.DeeperResetterTable[kind]; ok && resetter != nil {
			jumpTo = resetter.Run(n)
			skipAll = jumpTo == nil
			return true
		}

		if r.TriggerKind.Contain(kind) {
			triggered = true
			triggerIndex = 0
			return true
		}

		if triggerLimit == 0 || !triggered || !r.ChangeableKind.Contain(kind) {
			return true
		}

		triggerIndex++
		if triggerLimit > 0 {
			triggerLimit--
		}

		if k, ok := r.KindChangeTable[triggerIndex]; ok && k != KindNone {
			n.SetKind(k)
		}

		return true
	})
}

type scopeResetter struct{}

func (r scopeResetter) Reset(s Scope) {
	switch s.Kind() {
	case ScopePackage:
	case ScopeComment:
	case ScopeInnerComment:
	case ScopeImport:
	case ScopeVariable:
	case ScopeConst:
	case ScopeType:
	case ScopeFunc:
	default:

	}
}

type paramResetter struct {
	resetKind   Kind
	returnKinds []Kind
}

func (r paramResetter) Run(head *Node, _ ...func(*Node)) *Node {
	var (
		skipAll bool
		jumpTo  *Node
		buf     []*Node
	)

	if len(r.returnKinds) == 0 {
		return head.Next()
	}

	returnKindSet := newSet(r.returnKinds...)

	defer func() {
		if len(buf) != 0 {
			n := buf[0]
			next := buf[len(buf)-1].Next()
			n = n.CombineNext(r.resetKind, buf[1:]...)
			n.ReplaceNext(next)
		}
	}()

	return head.IterNext(func(n *Node) bool {
		if skipAll {
			return true
		}

		if jumpTo != nil {
			if n != jumpTo {
				return true
			}
			jumpTo = nil
		}

		if returnKindSet.Contain(n.Kind()) { // return kind
			return false
		}

		switch n.Kind() {
		case KindParenthesisLeft:
			jumpTo = parenthesisResetter{}.Run(n)
			skipAll = jumpTo == nil
			return true
		case KindCurlyBracketLeft:
			jumpTo = curlyBracketResetter{}.Run(n)
			skipAll = jumpTo == nil
			return true
		case KindSquareBracketLeft:
			jumpTo = squareBracketResetter{}.Run(n)
			skipAll = jumpTo == nil
			return true
		case KindFunc:
			jumpTo = funcResetter{isParameter: true}.Run(n)
			skipAll = jumpTo == nil
			return true
		case KindComment:
			return true
		default:
			return true
		}
	})
}
