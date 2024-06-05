package goast

// parenthesisResetter starts with '('
type parenthesisResetter struct {
	isReceiver bool
}

func (r parenthesisResetter) Run(head *Node) *Node {
	if head.Kind() != KindParenthesisLeft {
		return head.Next()
	}

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
		case KindParenthesisRight: // return kind
			return false
		case KindComment, KindComma, KindParenthesisLeft:
			return true
		default:
			jumpTo = r.handleParenthesisParam(n, r.isReceiver)
			skipAll = jumpTo == nil
			return true
		}

	})
}

// handleParenthesisParam starts with next of '(' and ','
func (r parenthesisResetter) handleParenthesisParam(head *Node, isReceiver bool) *Node {
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
				buf[0].Print()
				buf = buf[1:]
			}

			h := buf[0]
			next := buf[len(buf)-1].Next()
			h = h.CombineNext(typeKind, buf[1:]...)
			h.ReplaceNext(next)
			h.Print()
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

		switch n.Kind() {
		case KindComma, KindParenthesisRight: // return kind
			return false
		case KindSquareBracketLeft: // ignore including generic
			jumpTo = squareBracketResetter{}.Run(n, func(nn *Node) {
				buf = append(buf, nn)
			}).Next()
			skipAll = jumpTo == nil
			return true
		case KindComment, KindParenthesisLeft:
			return true
		case KindSpace:
			if len(buf) != 0 {
				hasSpaceAfterRaw = true
			}
			return true
		case KindFunc:
			jumpTo = funcResetter{isParameter: true}.Run(n)
			skipAll = jumpTo == nil
			return true
		default:
			if hasSpaceAfterRaw {
				hasName = true
			}
			buf = append(buf, n)
			return true
		}
	})
}

// curlyBracketResetter starts with '{'
type curlyBracketResetter struct{}

func (r curlyBracketResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	// TODO: handle content
	return head.skipNestNext(KindCurlyBracketLeft, KindCurlyBracketRight, hooks...)
}

// squareBracketResetter
type squareBracketResetter struct{}

func (r squareBracketResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	// TODO: handle generic
	return head.skipNestNext(KindSquareBracketLeft, KindSquareBracketRight, hooks...)
}
