package goast

// funcResetter includes:
//
//   - function with name and '{}' : func Hello() (string, error) {}
//
//   - function with no name but with '{}': func () string {}
//
//   - function with no name and no '{}': func () string
//
//   - function (method) with naming receiver: func (r *Receiver) Hello() string {}
//
//   - function (method) with no naming receiver: func (*Receiver) Hello() string {}
//
//   - [x] function in interface definition: Hello(string) string
//
//   - [x] function as parameter: (fn func(string) error)
//
//   - [x] temporary defined function: func(int, string) { ... }(5, "hello")
type funcResetter struct {
	isParameter           bool
	isFuncKeywordLeading  bool
	isInterfaceDefinition bool
}

func (r funcResetter) Run(head *Node) *Node {
	if r.isParameter {
		return r.handleParameterFunc(head)
	}

	if r.isInterfaceDefinition {
		return r.handleGeneralFunc(head)
	}

	if r.isTemporaryFunc(head) {
		return r.handleGeneralFunc(head)
	}

	// handle func keyword leading
	if r.isFuncKeywordLeading {
		if head.Kind() != KindFunc {
			return head.Next()
		}

		head = head.Next() // skip first 'func' keyword
	}

	var (
		foundFuncOrMethod bool
		isMethod          bool
		skipAll           bool
		jumpTo            *Node
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

		// determine func or method
		if !foundFuncOrMethod {
			switch n.Kind() {
			case KindRaw:
				foundFuncOrMethod = true
			case KindParenthesisLeft:
				foundFuncOrMethod = true
				isMethod = true
			default:
				return true
			}
			// keep going when find func/method
		}

		if isMethod {
			println("\t", "funcResetter.handleParenthesisParam.Start:", n.TidiedText(), n.Kind().String())
			isMethod = false
			jumpTo = r.handleParenthesisParam(head, true)
			skipAll = jumpTo == nil
			println("\t", "funcResetter.handleParenthesisParam.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
			println()
			return true
		}

		println("\t", "funcResetter.handleGeneralFunc.Start:", n.TidiedText(), n.Kind().String())
		jumpTo = r.handleGeneralFunc(head)
		skipAll = jumpTo == nil
		println("\t", "funcResetter.handleGeneralFunc.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
		println()
		return true
	})
}

// handleTemporaryFunc starts with 'func' or next of 'func'
func (r funcResetter) handleGeneralFunc(head *Node, returnKinds ...Kind) *Node {
	var (
		skipAll bool
		jumpTo  *Node

		isFuncNameAssigned bool
		isFuncParamHandled bool
		isReturnKind       set[Kind]
	)

	if len(returnKinds) != 0 {
		isReturnKind = newSet(returnKinds...)
	} else {
		isReturnKind = newSet(KindNewLine)
	}

	if head.Kind() == KindFunc {
		head = head.Next() // skip first 'func'
	}

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

		// return kind
		if isReturnKind.Contain(n.Kind()) {
			return !isFuncParamHandled
		}

		switch n.Kind() {
		case KindComment:
			return true
		case KindRaw:
			if !isFuncNameAssigned {
				n.SetKind(KindFuncName)
				isFuncNameAssigned = true
			}
			return true
		case KindSpace:
			if isFuncParamHandled {
				println("\t", "handleGeneralFunc.handleSingleReturnType.Start:", n.TidiedText(), n.Kind().String())
				jumpTo = r.handleSingleReturnType(n)
				skipAll = jumpTo == nil
				println("\t", "handleGeneralFunc.handleSingleReturnType.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
				println()
				return true
			}
			return true
		case KindParenthesisLeft:
			println("\t", "handleGeneralFunc.handleParenthesis.Start:", n.TidiedText(), n.Kind().String())
			isFuncParamHandled = true
			jumpTo = r.handleParenthesis(n)
			skipAll = jumpTo == nil
			println("\t", "handleGeneralFunc.handleParenthesis.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
			println()
			return true
		case KindCurlyBracketLeft:
			println("\t", "handleGeneralFunc.handleCurlyBracket.Start:", n.TidiedText(), n.Kind().String())
			jumpTo = r.handleCurlyBracket(n)
			skipAll = jumpTo == nil
			println("\t", "handleGeneralFunc.handleCurlyBracket.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
			println()
			return true
		case KindSquareBracketLeft:
			println("\t", "handleGeneralFunc.handleSquareBracket.Start:", n.TidiedText(), n.Kind().String())
			isFuncParamHandled = true
			jumpTo = r.handleSquareBracket(n)
			skipAll = jumpTo == nil
			println("\t", "handleGeneralFunc.handleSquareBracket.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
			println()
			return true
		default:
			return true
		}
	})
}

// isTemporaryFunc starts with 'func' or next of 'func'
func (r funcResetter) isTemporaryFunc(head *Node) bool {
	found := head.findNext(
		newSet(KindNewLine, KindCurlyBracketRight),
		findNodeOption{
			IsOutsideParenthesis:   true,
			IsOutsideCurlyBracket:  true,
			IsOutsideSquareBracket: true,
		},
	)

	if found.Kind() != KindCurlyBracketRight {
		return false
	}

	return found.Next().Kind() == KindParenthesisLeft
}

// handleParameterFunc starts with 'func'
func (r funcResetter) handleParameterFunc(head *Node) *Node {
	n := r.handleGeneralFunc(head, KindComma, KindParenthesisRight)
	if n.Kind() == KindParenthesisRight {
		return n.Next()
	}
	return n
}

// handleParenthesis starts with '('
func (r funcResetter) handleParenthesis(head *Node) *Node {
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
			println("\t", "handleParenthesis.handleParenthesisParam.Start:", n.TidiedText(), n.Kind().String())
			jumpTo = r.handleParenthesisParam(n, false)
			skipAll = jumpTo == nil
			println("\t", "handleParenthesis.handleParenthesisParam.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
			println()
			return true
		}

	})
}

// handleParenthesisParam starts with next of '(' and ','
func (r funcResetter) handleParenthesisParam(head *Node, isReceiver bool) *Node {
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
			if len(buf) == 0 {
				return false
			}

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
			return false
		case KindSquareBracketLeft: // ignore including generic
			println("\t", "handleParenthesisParam.handleSquareBracket.Start:", n.TidiedText(), n.Kind().String())
			jumpTo = r.handleSquareBracket(n, func(nn *Node) {
				buf = append(buf, nn)
			}).Next()
			skipAll = jumpTo == nil
			println("\t", "handleParenthesisParam.handleSquareBracket.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
			println()
			return true
		case KindComment, KindParenthesisLeft:
			return true
		case KindSpace:
			if len(buf) != 0 {
				hasSpaceAfterRaw = true
			}
			return true
		case KindFunc:
			println("\t", "handleParenthesisParam.funcResetter.Start:", n.TidiedText(), n.Kind().String())
			jumpTo = funcResetter{isParameter: true}.Run(n)
			skipAll = jumpTo == nil
			println("\t", "handleParenthesisParam.funcResetter.JumpTo:", jumpTo.TidiedText(), jumpTo.Kind().String())
			println()
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

// handleSingleReturnType starts with ' '
func (r funcResetter) handleSingleReturnType(head *Node) *Node {
	if head.Kind() != KindSpace {
		return head.Next()
	}

	head = head.Next() // skip first space

	found := head.findNext(newSet(KindComment), findNodeOption{TargetReverse: true})
	switch found.Kind() {
	case KindParenthesisLeft:
		return r.handleParenthesis(head)
	case KindFunc:
		return r.handleParameterFunc(head)
	}

	var (
		buf []*Node
	)

	return head.IterNext(func(n *Node) bool {
		switch n.Kind() {
		case KindNewLine, KindSpace: // return kind
			if len(buf) != 0 {
				n := buf[0]
				next := buf[len(buf)-1].Next()
				n = n.CombineNext(KindParamType, buf[1:]...)
				n.ReplaceNext(next)
			}
			return false
		case KindComment:
			return true
		default:
			buf = append(buf, n)
			return true
		}
	})
}

// handleCurlyBracket starts with '{'
func (r funcResetter) handleCurlyBracket(head *Node, hook ...func(*Node)) *Node {
	// TODO: handle content
	return head.skipNestNext(KindCurlyBracketLeft, KindCurlyBracketRight, hook...)
}

func (r funcResetter) handleSquareBracket(head *Node, hook ...func(*Node)) *Node {
	// TODO: handle generic
	return head.skipNestNext(KindSquareBracketLeft, KindSquareBracketRight, hook...)
}
