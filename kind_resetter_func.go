package goast

import "github.com/yanun0323/goast/helper"

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

func (r funcResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	if r.isParameter {
		return r.handleParameterFunc(head, hooks...)
	}

	if r.isInterfaceDefinition {
		return r.handleGeneralFunc(head, []Kind{KindNewLine}, hooks...)
	}

	if r.isTemporaryFunc(head) {
		return r.handleGeneralFunc(head, []Kind{KindNewLine}, hooks...)
	}

	// handle func keyword leading
	if r.isFuncKeywordLeading {
		if head.Kind() != KindFunc {
			handleHook(head, hooks...)
			return head.Next()
		}

		handleHook(head, hooks...)
		head = head.Next() // skip first 'func' keyword
	}

	var (
		foundFuncOrMethod bool
		isMethod          bool
		skipAll           bool
		jumpTo            *Node
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

		if isMethod { // '('
			isMethod = false
			jumpTo = parenthesisResetter{isReceiver: true}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		}

		switch n.Kind() {
		case KindNewLine: // return kind
			return false
		default:
			jumpTo = r.handleGeneralFunc(head, []Kind{KindNewLine}, hooks...)
			skipAll = jumpTo == nil
			return true
		}
	})
}

// handleTemporaryFunc starts with 'func' or next of 'func'
func (r funcResetter) handleGeneralFunc(head *Node, returnKinds []Kind, hooks ...func(*Node)) *Node {
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
		handleHook(head, hooks...)
		head = head.Next() // skip first 'func'
	}

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
				jumpTo = r.handleSingleReturnType(n, hooks...)
				skipAll = jumpTo == nil
				return true
			}
			return true
		case KindParenthesisLeft:
			isFuncParamHandled = true
			jumpTo = parenthesisResetter{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case KindCurlyBracketLeft:
			jumpTo = curlyBracketResetter{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case KindSquareBracketLeft:
			isFuncParamHandled = true
			jumpTo = squareBracketResetter{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		default:
			return true
		}
	})
}

// isTemporaryFunc starts with 'func' or next of 'func'
func (r funcResetter) isTemporaryFunc(head *Node) bool {
	found := head.findNext(
		[]Kind{KindNewLine, KindCurlyBracketRight},
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
func (r funcResetter) handleParameterFunc(head *Node, hooks ...func(*Node)) *Node {
	n := r.handleGeneralFunc(head, []Kind{KindComma, KindParenthesisRight}, hooks...)
	if n.Kind() == KindParenthesisRight {
		handleHook(n, hooks...)
		return n.Next()
	}
	return n
}

// handleSingleReturnType starts with ' '
func (r funcResetter) handleSingleReturnType(head *Node, hooks ...func(*Node)) *Node {
	if head.Kind() != KindSpace {
		handleHook(head, hooks...)
		return head.Next()
	}

	head = head.Next() // skip first space

	found := head.findNext([]Kind{KindComment}, findNodeOption{TargetReverse: true}, hooks...)
	switch found.Kind() {
	case KindParenthesisLeft:
		return parenthesisResetter{}.Run(head, hooks...)
	case KindFunc:
		return r.handleParameterFunc(head, hooks...)
	}

	var (
		buf []*Node
	)

	defer func() {
		if len(buf) != 0 {
			n := buf[0]
			next := buf[len(buf)-1].Next()
			n = n.CombineNext(KindParamType, buf[1:]...)
			n.ReplaceNext(next)
		}
	}()

	return head.IterNext(func(n *Node) bool {
		handleHook(n, hooks...)
		switch n.Kind() {
		case KindNewLine, KindSpace: // return kind
			return false
		case KindComment:
			return true
		default:
			buf = helper.AppendUnrepeatable(buf, n)
			return true
		}
	})
}
