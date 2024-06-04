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

// funcResetter head should be 'func'
func funcResetter() kindResetter {
	return &resetter{
		TriggerKind:     newSet(KindFunc),
		TriggerLimit:    -1,
		KindChangeTable: map[int]Kind{1: KindFuncName},
		ChangeableKind:  newSet(KindRaw),
		ReturnKind:      newSet(KindCurlyBracketLeft),
		DeeperResetterTable: map[Kind]kindResetter{
			KindParenthesisLeft: &funcParamResetter{},
		},
	}
}

// resetter stands for COMMON resetter
type resetter struct {
	TriggerKind         set[Kind]
	TriggerLimit        int
	KindChangeTable     map[int]Kind
	ChangeableKind      set[Kind]
	ReturnKind          set[Kind]
	DeeperResetterTable map[Kind]kindResetter
}

func (r resetter) Run(head *Node) *Node {
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
	}).Next()
}

// funcParamResetter head should be '('
//
// return key kind: '{'
type funcParamResetter struct {
	returnKeyKind Kind
}

func (r funcParamResetter) Run(head *Node) *Node {
	if head.Kind() != KindParenthesisLeft {
		return head
	}

	var (
		skipAll               bool
		jumpTo                *Node
		nameNode              *Node
		buf                   []*Node
		nextIsFirstParam      bool
		isFirstParamANameNode bool // determine by space
	)

	head = head.Next() // ignore first '('
	nextIsFirstParam = true

	resetBuf := func() {
		buf = []*Node{}
		nameNode = nil
		nextIsFirstParam = true
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
		if skipAll {
			return true
		}

		if jumpTo != nil {
			if jumpTo != n {
				return true
			}
			jumpTo = nil
		}
		kind := n.Kind()
		switch kind {
		case KindParenthesisRight, r.returnKeyKind: // return kind
			dealBuf()
			return false
		case KindComment:
			return true
		case KindSpace:
			isFirstParamANameNode = nameNode != nil
			return true
		case KindParenthesisLeft:
			return true
		case KindComma:
			dealBuf()
			nextIsFirstParam = true
			return true
		case KindRaw:
			if nextIsFirstParam {
				n.SetKind(KindParamName)
				nextIsFirstParam = false
				nameNode = n
				return true
			}
		case KindFunc:
			jumpTo = paramFuncResetter{}.Run(n)
			resetBuf()
			skipAll = jumpTo == nil
			return true
		}

		buf = append(buf, n)
		return true
	}).Next()
}

// paramFuncResetter makes all input turns into KindParamType
//
// return key kind: ')' ','
type paramFuncResetter struct {
	returnKind Kind
}

func (r paramFuncResetter) Run(head *Node) *Node {
	var (
		buf                  []*Node
		parenthesisLeftCount int
	)

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

	return head.IterNext(func(n *Node) bool {
		kind := n.Kind()
		switch kind {
		case KindParenthesisRight, r.returnKind: // return kind
			keep := tryAppendBuf(n)
			parenthesisLeftCount--
			return keep
		case KindComma: // return kind
			return tryAppendBuf(n)
		case KindParenthesisLeft:
			parenthesisLeftCount++
			buf = append(buf, n)
			return true
		default:
			buf = append(buf, n)
			return true
		}
	}).Next()
}

// funcSingleReturnTypeResetter head should be ' ', and next should not be '('.
//
// return key kind: '{'
type funcSingleReturnTypeResetter struct {
	returnKind Kind
}

func (r funcSingleReturnTypeResetter) Run(head *Node) *Node {
	if head.Kind() != KindSpace {
		return head
	}

	if head.Next().Kind() == KindParenthesisLeft {
		return head.Next()
	}

	var (
		skipAll bool
		jumpTo  *Node
		buf     []*Node
	)

	defer func() {
		if len(buf) == 0 {
			return
		}

		n := buf[0]
		next := buf[len(buf)-1].Next()
		n.CombineNext(KindParamType, buf[1:]...)
		n.ReplaceNext(next)
	}()

	head = head.Next() // skip first space

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

		switch n.Kind() {
		case KindCurlyBracketLeft, r.returnKind: // return kind
			return false
		case KindParenthesisRight:
			jumpTo = funcParamResetter{
				returnKeyKind: KindNewLine,
			}.Run(n)
			skipAll = jumpTo == nil
			return true
		case KindComment:
			return true
		default:
			buf = append(buf, n)
			return true
		}
	}).Next()
}

// typeResetter head should be 'type'
//
// return key kind: '}' '\n'
type typeResetter struct{}

func (r typeResetter) Run(head *Node) *Node {
	if head.Kind() != KindType {
		return head
	}

	var (
		isTypeNameAssigned bool
		isInterface        bool
		isStruct           bool
		exception          bool
		skipAll            bool
		jumpTo             *Node
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

		if exception {
			switch n.Kind() {
			case KindCurlyBracketRight, KindNewLine: // return kind
				return false
			default:
				return true
			}
		}

		if !isTypeNameAssigned {
			switch n.Kind() {
			case KindCurlyBracketRight, KindNewLine: // return kind
				return false
			case KindRaw:
				n.SetKind(KindTypeName)
				isTypeNameAssigned = true
			default:
				return true
			}
		}

		switch n.Kind() {
		case KindComment:
			return true
		case KindInterface:
			isInterface = true
			return true
		case KindStruct:
			isStruct = true
			return true
		case KindCurlyBracketLeft:
			if isInterface {
				jumpTo = r.interfaceResetter(n)
				skipAll = jumpTo == nil
				return true
			}

			if isStruct {
				jumpTo = r.structResetter(n)
				skipAll = jumpTo == nil
				return true
			}

			exception = true
			return true
		case KindRaw:
			jumpTo = r.otherResetter(n)
			skipAll = jumpTo == nil
			return true
		default:
			return true
		}
	}).Next()
}

// interfaceResetter
//
// return key kind: '}'
func (r typeResetter) interfaceResetter(head *Node) *Node {
	var (
		skipAll                 bool
		jumpTo                  *Node
		canAssignFuncReturnType bool

		canAssignFuncName = true
	)

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

		if canAssignFuncReturnType { // ' ' after ')'
			println("canAssignFuncReturnType:")
			n.Print()
			n.Next().Print()
			switch n.Kind() {
			case KindComment:
				return true
			case KindSpace:
				canAssignFuncReturnType = false
				jumpTo = funcSingleReturnTypeResetter{
					returnKind: KindNewLine,
				}.Run(n)
				canAssignFuncName = jumpTo.Kind() != KindNewLine
				skipAll = jumpTo == nil
				return true
			default:
				canAssignFuncReturnType = false
				// keep going
			}
			println("canAssignFuncReturnType done")
		}

		switch n.Kind() {
		case KindCurlyBracketRight: // return kind
			return false
		case KindComment, KindCurlyBracketLeft, KindTab:
			return true
		case KindNewLine:
			canAssignFuncName = true
			return true
		case KindParenthesisLeft:
			println("(:", n.Text(), n.Next().Text())
			jumpTo = funcParamResetter{
				returnKeyKind: KindNewLine,
			}.Run(n)
			canAssignFuncName = jumpTo.Kind() != KindNewLine
			canAssignFuncReturnType = jumpTo.Prev().Kind() == KindParenthesisRight
			skipAll = jumpTo == nil
			println("jumpTO:", jumpTo.Text(), jumpTo.Kind().String())
			return true
		case KindRaw:
			if canAssignFuncName {
				n.SetKind(KindFuncName)
				canAssignFuncName = false
			}
			return true
		default:
			return true
		}
	}).Next()
}

// structResetter
//
// return key kind: '}'
func (r typeResetter) structResetter(head *Node) *Node {
	return head.IterNext(func(n *Node) bool {
		switch n.Kind() {
		case KindComment, KindCurlyBracketLeft, KindTab:
			return true
		case KindCurlyBracketRight:
			return false
		default:
			return true
		}
	}).Next()
}

// otherResetter
//
// return key kind: '\n'
func (r typeResetter) otherResetter(head *Node) *Node {
	var (
		buf []*Node
	)

	defer func() {
		if len(buf) == 0 {
			return
		}

		n := buf[0]
		next := buf[len(buf)-1].Next()
		n = n.CombineNext(KindTypeAliasType, buf[1:]...)
		n.ReplaceNext(next)
	}()

	return head.IterNext(func(n *Node) bool {
		switch n.Kind() {
		case KindNewLine: // return kind
			return false
		case KindComment:
			return true
		default:
			buf = append(buf)
			return true
		}
	}).Next()
}
