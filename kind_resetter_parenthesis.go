package goast

// parenthesisResetter starts with '('
type parenthesisResetter struct {
	skip       bool
	isReceiver bool
}

func (r parenthesisResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	if r.skip {
		return head.skipNestNext(KindParenthesisLeft, KindParenthesisRight, hooks...)
	}

	if head.Kind() != KindParenthesisLeft {
		handleHook(head, hooks...)
		return head.Next()
	}

	var (
		skipAll bool
		jumpTo  *Node
	)

	return head.IterNext(func(n *Node) bool {
		handleHook(n, hooks...)
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
		case KindParenthesisRight: // return kind
			return false
		case KindComment, KindComma, KindParenthesisLeft, KindTab, KindNewLine:
			return true
		default:
			jumpTo = r.handleParenthesisParam(n, r.isReceiver, hooks...)
			skipAll = jumpTo == nil
			return true
		}

	})
}

// handleParenthesisParam starts with next of '(' and ','
func (r parenthesisResetter) handleParenthesisParam(head *Node, isReceiver bool, hooks ...func(*Node)) *Node {
	println("handleParenthesisParam:", head.debugText(3))
	var (
		skipAll          bool
		jumpTo           *Node
		buf              []*Node
		hasSpaceAfterRaw bool
		hasName          bool
	)

	nameKind := KindParamName
	typeKind := KindParamType
	if isReceiver {
		nameKind = KindMethodReceiverName
		typeKind = KindMethodReceiverType
	}

	defer func() {
		if len(buf) != 0 {
			if hasName {
				buf[0].SetKind(nameKind)
				buf = buf[1:]
			}

			h := buf[0]
			next := buf[len(buf)-1].Next()
			h = h.CombineNext(typeKind, buf[1:]...)
			h.ReplaceNext(next)
		}
	}()

	return head.IterNext(func(n *Node) bool {
		handleHook(n, hooks...)
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
		case KindComma, KindParenthesisRight: // return kind
			return false
		case KindSquareBracketLeft: // ignore including generic
			jumpTo = squareBracketResetter{}.Run(n, append(hooks, func(nn *Node) {
				buf = appendUnrepeatable(buf, nn)
			})...).Next()
			skipAll = jumpTo == nil
			return true
		case KindComment, KindParenthesisLeft, KindTab, KindNewLine:
			return true
		case KindSpace:
			if len(buf) != 0 {
				hasSpaceAfterRaw = true
			}
			return true
		case KindFunc:
			jumpTo = funcResetter{isParameter: true}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		default:
			if hasSpaceAfterRaw {
				hasName = true
			}
			buf = appendUnrepeatable(buf, n)
			return true
		}
	})
}

// curlyBracketResetter starts with '{'
type curlyBracketResetter struct {
	skip bool
}

func (r curlyBracketResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	if r.skip {
		return head.skipNestNext(KindCurlyBracketLeft, KindCurlyBracketRight, hooks...)
	}

	// TODO: handle content
	return head.skipNestNext(KindCurlyBracketLeft, KindCurlyBracketRight, hooks...)
}

// squareBracketResetter
type squareBracketResetter struct {
	skip bool
}

func (r squareBracketResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	if r.skip {
		return head.skipNestNext(KindSquareBracketLeft, KindSquareBracketRight, hooks...)
	}

	// TODO: handle generic
	return head.skipNestNext(KindSquareBracketLeft, KindSquareBracketRight, hooks...)
}
