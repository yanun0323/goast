package goast

func kindReset(n *Node) *Node {
	return generalResetter().Run(n)
}

type kindResetter interface {
	Run(*Node) *Node
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

func (r genericResetter) Run(head *Node) *Node {
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
			println("genericResetter.funcResetter.Start:", n.TidiedText(), n.Kind().String())
			jumpTo = resetter.Run(n)
			skipAll = jumpTo == nil
			println("genericResetter.funcResetter.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
			println()
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

// // paramFuncResetter makes all input turns into KindParamType
// //
// // return key kind: ')' ','
// type paramFuncResetter struct {
// 	returnKind Kind
// }

// func (r paramFuncResetter) Run(head *Node) *Node {
// 	var (
// 		buf                  []*Node
// 		parenthesisLeftCount int
// 	)

// 	defer func() {
// 		if len(buf) == 0 {
// 			return
// 		}

// 		n := buf[0]
// 		next := buf[len(buf)-1].Next()
// 		n = n.CombineNext(KindParamType, buf[1:]...)
// 		n.ReplaceNext(next)
// 	}()

// 	tryAppendBuf := func(n *Node) bool {
// 		if parenthesisLeftCount == 0 {
// 			return false
// 		}
// 		buf = append(buf, n)
// 		return true
// 	}

// 	return head.IterNext(func(n *Node) bool {
// 		kind := n.Kind()
// 		switch kind {
// 		case KindParenthesisRight, r.returnKind: // return kind
// 			keep := tryAppendBuf(n)
// 			parenthesisLeftCount--
// 			return keep
// 		case KindComma: // return kind
// 			return tryAppendBuf(n)
// 		case KindParenthesisLeft:
// 			parenthesisLeftCount++
// 			buf = append(buf, n)
// 			return true
// 		default:
// 			buf = append(buf, n)
// 			return true
// 		}
// 	})
// }

// // funcSingleReturnTypeResetter head should be ' ', and next should not be '('.
// //
// // return key kind: '{'
// type funcSingleReturnTypeResetter struct {
// 	returnKind Kind
// }

// func (r funcSingleReturnTypeResetter) Run(head *Node) *Node {
// 	if head.Kind() != KindSpace {
// 		return head.Next()
// 	}

// 	if head.Next().Kind() == KindParenthesisLeft {
// 		return head.Next()
// 	}

// 	var (
// 		skipAll bool
// 		jumpTo  *Node
// 		buf     []*Node
// 	)

// 	defer func() {
// 		if len(buf) == 0 {
// 			return
// 		}

// 		n := buf[0]
// 		next := buf[len(buf)-1].Next()
// 		n.CombineNext(KindParamType, buf[1:]...)
// 		n.ReplaceNext(next)
// 	}()

// 	head = head.Next() // skip first space

// 	return head.IterNext(func(n *Node) bool {
// 		if skipAll {
// 			return true
// 		}

// 		if jumpTo != nil {
// 			if n != jumpTo {
// 				return true
// 			}
// 			jumpTo = nil
// 		}

// 		switch n.Kind() {
// 		case KindCurlyBracketLeft, r.returnKind: // return kind
// 			return false
// 		case KindParenthesisRight:
// 			jumpTo = funcParamResetter{
// 				returnKeyKind: KindNewLine,
// 			}.Run(n)
// 			skipAll = jumpTo == nil
// 			return true
// 		case KindComment:
// 			return true
// 		default:
// 			buf = append(buf, n)
// 			return true
// 		}
// 	})
// }

// funcParamResetter head should be '('
//
// return key kind: '{'
// type funcParamResetter struct {
// 	returnKeyKind Kind
// }

// func (r funcParamResetter) Run(head *Node) *Node {
// 	if head.Kind() != KindParenthesisLeft {
// 		return head.Next()
// 	}

// 	var (
// 		skipAll               bool
// 		jumpTo                *Node
// 		nameNode              *Node
// 		buf                   []*Node
// 		nextIsFirstParam      bool
// 		isFirstParamANameNode bool // determine by space
// 	)

// 	head = head.Next() // ignore first '('
// 	nextIsFirstParam = true

// 	resetBuf := func() {
// 		buf = []*Node{}
// 		nameNode = nil
// 		nextIsFirstParam = true
// 		isFirstParamANameNode = false
// 	}

// 	dealBuf := func() {
// 		defer resetBuf()

// 		n := nameNode
// 		if isFirstParamANameNode {
// 			if len(buf) == 0 {
// 				return
// 			}
// 			n = buf[0]
// 			buf = buf[1:]
// 		}

// 		if len(buf) == 0 {
// 			n.SetKind(KindParamType)
// 			return
// 		}

// 		next := buf[len(buf)-1].Next()
// 		n = n.CombineNext(KindParamType, buf...)
// 		n.ReplaceNext(next)
// 	}

// 	return head.IterNext(func(n *Node) bool {
// 		if skipAll {
// 			return true
// 		}

// 		if jumpTo != nil {
// 			if jumpTo != n {
// 				return true
// 			}
// 			jumpTo = nil
// 		}
// 		kind := n.Kind()
// 		switch kind {
// 		case KindParenthesisRight, r.returnKeyKind: // return kind
// 			dealBuf()
// 			return false
// 		case KindComment:
// 			return true
// 		case KindSpace:
// 			isFirstParamANameNode = nameNode != nil
// 			return true
// 		case KindParenthesisLeft:
// 			return true
// 		case KindComma:
// 			dealBuf()
// 			nextIsFirstParam = true
// 			return true
// 		case KindRaw:
// 			if nextIsFirstParam {
// 				n.SetKind(KindParamName)
// 				nextIsFirstParam = false
// 				nameNode = n
// 				return true
// 			}
// 		case KindFunc:
// 			jumpTo = paramFuncResetter{}.Run(n)
// 			resetBuf()
// 			skipAll = jumpTo == nil
// 			return true
// 		}

// 		buf = append(buf, n)
// 		return true
// 	})
// }
