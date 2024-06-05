package goast

// typeResetter head should be 'type'
//
// return key kind: '}' '\n'
type typeResetter struct{}

func (r typeResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	if head.Kind() != KindType {
		handleHook(head, hooks...)
		return head.Next()
	}

	var (
		isTypeNameAssigned bool
		exception          bool
		skipAll            bool
		jumpTo             *Node
	)

	return head.IterNext(func(n *Node) bool {
		handleHook(n, hooks...)
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
			jumpTo = interfaceResetter{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case KindStruct:
			jumpTo = structResetter{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case KindRaw:
			jumpTo = r.otherResetter(n, hooks...)
			skipAll = jumpTo == nil
			return true
		default:
			return true
		}
	})
}

// interfaceResetter starts with 'interface'
//
// return key kind: '}'
type interfaceResetter struct{}

func (r interfaceResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	var (
		skipAll bool
		jumpTo  *Node
	)

	head = head.findNext([]Kind{KindCurlyBracketLeft}, findNodeOption{}, hooks...) // set head to '{'

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
		case KindCurlyBracketRight: // return kind
			return false
		case KindTab, KindComment, KindNewLine, KindSpace, KindCurlyBracketLeft:
			return true
		default:
			jumpTo = funcResetter{
				isInterfaceDefinition: true,
			}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		}
	})
}

// structResetter starts with 'struct'
//
// return key kind: '}'
type structResetter struct{}

func (r structResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	var (
		skipAll bool
		jumpTo  *Node
	)

	head = head.findNext([]Kind{KindCurlyBracketLeft}, findNodeOption{}, hooks...).Next() // skip first of head to '{'

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
		case KindCurlyBracketRight: // return kind
			return false
		case KindComment, KindCurlyBracketLeft, KindNewLine:
			return true
		default:
			jumpTo = r.handleStructRow(n, hooks...)
			skipAll = jumpTo == nil
			return true
		}
	})
}

// handleStructRow returns key kind: '\n'
//
//   - a, b, c int
//
//   - a, b, c func(int) (int, error)
//
//   - a, b, c struct{}
func (r structResetter) handleStructRow(head *Node, hooks ...func(*Node)) *Node {
	var (
		skipAll bool
		jumpTo  *Node

		paramNameCount = r.getRowNameCount(head)
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
		case KindNewLine: // return kind
			return false
		case KindComment, KindNewLine, KindTab, KindSpace, KindComma:
			return true
		case KindRaw:
			if paramNameCount != 0 {
				paramNameCount--
				n.SetKind(KindParamName)
				return true
			}
			jumpTo = paramResetter{resetKind: KindParamType, returnKinds: []Kind{KindNewLine}}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		default:
			jumpTo = paramResetter{resetKind: KindParamType, returnKinds: []Kind{KindNewLine}}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		}
	})
}

func (r structResetter) getRowNameCount(head *Node) int {
	var (
		rawCount           int
		hasComma           bool
		hasSpaceOrAfterRaw bool
		cacheNameCount     int
		nameCount          int
	)

	_ = head.IterNext(func(n *Node) bool {
		switch n.Kind() {
		case KindNewLine: // return kind
			if hasSpaceOrAfterRaw {
				nameCount = cacheNameCount
			}
			return false
		case KindComment, KindTab:
			return true
		case KindRaw:
			rawCount++
			return true
		case KindSpace:
			if rawCount != 0 {
				if !hasComma && cacheNameCount == 0 {
					cacheNameCount++
				}
				hasSpaceOrAfterRaw = true
			}
			return true
		case KindComma:
			if cacheNameCount == 0 {
				cacheNameCount++
			}
			hasComma = true
			cacheNameCount++
			nameCount = cacheNameCount
			return true
		default:
			if hasSpaceOrAfterRaw {
				nameCount = cacheNameCount
			}
			return false
		}
	})

	return nameCount
}

// otherResetter
//
// return key kind: '\n'
func (r typeResetter) otherResetter(head *Node, hooks ...func(*Node)) *Node {
	var (
		buf []*Node
	)

	defer func() {
		if len(buf) != 0 {
			n := buf[0]
			next := buf[len(buf)-1].Next()
			n = n.CombineNext(KindTypeAliasType, buf[1:]...)
			n.ReplaceNext(next)
		}
	}()

	return head.IterNext(func(n *Node) bool {
		handleHook(n, hooks...)
		switch n.Kind() {
		case KindNewLine: // return kind
			return false
		case KindComment:
			return true
		default:
			buf = appendUnrepeatable(buf, n)
			return true
		}
	})
}
