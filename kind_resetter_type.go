package goast

// typeResetter head should be 'type'
//
// return key kind: '}' '\n'
type typeResetter struct{}

func (r typeResetter) Run(head *Node) *Node {
	if head.Kind() != KindType {
		return head.Next()
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
			if jumpTo == n {
				jumpTo = nil
			}
			return true
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
	})
}

// interfaceResetter
//
// return key kind: '}'
func (r typeResetter) interfaceResetter(head *Node) *Node {
	var (
		skipAll bool
		jumpTo  *Node
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
		switch n.Kind() {
		case KindCurlyBracketRight: // return kind
			return false
		case KindTab, KindComment, KindNewLine, KindSpace, KindCurlyBracketLeft:
			return true
		default:
			println("typeResetter.funcResetter.Start:", n.TidiedText(), n.Kind().String())
			jumpTo = funcResetter{
				isInterfaceDefinition: true,
			}.Run(n)
			skipAll = jumpTo == nil
			println("typeResetter.funcResetter.jumpTo:", jumpTo.TidiedText(), jumpTo.kind.String())
			println()
			return true
		}
	})
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
	})
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
	})
}
