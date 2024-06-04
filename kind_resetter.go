package goast

func kindReset(n *Node) *Node {
	return commonResetter().Run(n)
}

type kindResetter interface {
	Run(*Node) *Node
}

func commonResetter() kindResetter {
	return &resetter{
		DeeperResetterTable: map[Kind]kindResetter{
			KindFunc: funcResetter(),
		},
	}
}

func funcResetter() kindResetter {
	return &resetter{
		TriggerKind:     newSet(KindFunc),
		TriggerLimit:    -1,
		KindChangeTable: map[int]Kind{1: KindFuncName},
		ChangeableKind:  newSet(KindRaws),
		ReturnKind:      newSet(KindCurlyBracketLeft),
		DeeperResetterTable: map[Kind]kindResetter{
			KindParenthesisLeft: funcParamResetter(),
		},
	}
}

func inlineFuncResetter() kindResetter {
	return &resetter{
		// TODO: Implement me
	}
}

type resetter struct {
	TriggerKind         set[Kind]
	TriggerLimit        int
	KindChangeTable     map[int]Kind
	ChangeableKind      set[Kind]
	ReturnKind          set[Kind]
	DeeperResetterTable map[Kind]kindResetter
}

func (r *resetter) Run(head *Node) *Node {
	var jumpTo *Node

	triggered := false
	triggerIndex := 0
	triggerLimit := r.TriggerLimit

	return head.IterNext(func(n *Node) bool {
		if jumpTo != nil {
			if jumpTo == n {
				jumpTo = nil
			}
			return true
		}

		//  ReturnKind > deeperResetter > TriggerKind > ChangeableKind > UnchangeableKind
		kind := n.Kind()
		if r.ReturnKind.Contain(kind) {
			return false
		}

		if resetter, ok := r.DeeperResetterTable[kind]; ok && resetter != nil {
			jumpTo = resetter.Run(n)
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

func funcParamResetter() kindResetter {
	return &fnParamResetter{
		DeeperResetterTable: map[Kind]kindResetter{
			KindFunc: inlineFuncResetter(),
		},
	}
}

type fnParamResetter struct {
	DeeperResetterTable map[Kind]kindResetter
}

func (r *fnParamResetter) Run(head *Node) *Node {
	var (
		jumpTo   *Node
		nameNode *Node
	)

	head = head.Next() // ignore first '('
	nextIsFirstParam := true
	isFirstParamANameNode := false // determine by space
	buf := []*Node{}

	resetBuf := func() {
		buf = []*Node{}
		nameNode = nil
		isFirstParamANameNode = false
	}

	dealBuf := func() {
		defer resetBuf()

		n := nameNode
		if isFirstParamANameNode {
			if len(buf) == 0 {
				return
			}
			n = buf[0]
			buf = buf[1:]
		}

		if len(buf) == 0 {
			n.SetKind(KindParamType)
			return
		}

		next := buf[len(buf)-1].Next()
		n = n.CombineNext(KindParamType, buf...)
		n.ReplaceNext(next)
	}

	return head.IterNext(func(n *Node) bool {
		if jumpTo != nil {
			if jumpTo == n {
				jumpTo = nil
			}
			return true
		}
		kind := n.Kind()
		switch kind {
		case KindComments:
			return true
		case KindSpace:
			isFirstParamANameNode = nameNode != nil
			return true
		case KindParenthesisLeft:
			return true
		case KindParenthesisRight:
			dealBuf()
			return false
		case KindComma:
			dealBuf()
			nextIsFirstParam = true
			return true
		case KindRaws:
			if nextIsFirstParam {
				n.SetKind(KindParamName)
				nextIsFirstParam = false
				nameNode = n
				return true
			}
		case KindFunc:
			jumpTo = r.skipFuncParam(n)
			resetBuf()
			return true
		}

		if resetter, ok := r.DeeperResetterTable[kind]; ok && resetter != nil {
			jumpTo = resetter.Run(n)
			return true
		}

		buf = append(buf, n)

		return true
	})
}
func (r *fnParamResetter) skipFuncParam(fn *Node) *Node {
	buf := make([]*Node, 0, 10)
	parenthesisLeftCount := 0

	defer func() {
		if len(buf) == 0 {
			return
		}

		n := buf[0]
		next := buf[len(buf)-1].Next()
		n = n.CombineNext(KindParamType, buf[1:]...)
		n.ReplaceNext(next)
	}()

	tryAppendBuf := func(n *Node) bool {
		if parenthesisLeftCount == 0 {
			return false
		}
		buf = append(buf, n)
		return true
	}

	return fn.IterNext(func(n *Node) bool {
		kind := n.Kind()
		switch kind {
		case KindParenthesisLeft:
			parenthesisLeftCount++
			buf = append(buf, n)
			return true
		case KindParenthesisRight: // return symbol
			keep := tryAppendBuf(n)
			parenthesisLeftCount--
			return keep
		case KindComma: // return symbol
			return tryAppendBuf(n)
		default:
			buf = append(buf, n)
			return true
		}
	})
}
