package goast

import (
	"github.com/yanun0323/goast/charset"
	"github.com/yanun0323/goast/helper"
	"github.com/yanun0323/goast/kind"
)

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
	isNotMethod           bool
}

func (r funcResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("funcResetter.Run", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("funcResetter.Run.Returned")

	if r.isParameter {
		return r.handleParameterFunc(head, hooks...)
	}

	if r.isInterfaceDefinition {
		return r.handleGeneralFunc(head, []kind.Kind{kind.NewLine}, hooks...)
	}

	if r.isTemporaryFunc(head) {
		return r.handleGeneralFunc(head, []kind.Kind{kind.NewLine}, hooks...)
	}

	// handle func keyword leading
	if r.isFuncKeywordLeading {
		if head.Kind() != kind.Func {
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

	if r.isNotMethod {
		foundFuncOrMethod = true
		isMethod = false
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

		// determine func or method
		if !foundFuncOrMethod {
			switch n.Kind() {
			case kind.Raw:
				foundFuncOrMethod = true
			case kind.ParenthesisLeft:
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
		case kind.NewLine: // return kind
			return false
		default:
			jumpTo = r.handleGeneralFunc(head, []kind.Kind{kind.NewLine}, hooks...)
			skipAll = jumpTo == nil
			return true
		}
	})
}

// handleTemporaryFunc starts with 'func' or next of 'func'
func (r funcResetter) handleGeneralFunc(head *Node, returnKinds []kind.Kind, hooks ...func(*Node)) *Node {
	helper.DebugPrint("funcResetter.handleGeneralFunc", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("funcResetter.handleGeneralFunc.Returned")

	var (
		skipAll bool
		jumpTo  *Node

		isFuncNameAssigned bool
		isFuncParamHandled bool
		isReturnKind       charset.Set[kind.Kind]
	)

	if len(returnKinds) != 0 {
		isReturnKind = charset.New(returnKinds...)
	} else {
		isReturnKind = charset.New(kind.NewLine)
	}

	if head.Kind() == kind.Func {
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
		case kind.Comment:
			return true
		case kind.Raw:
			if !isFuncNameAssigned {
				n.SetKind(kind.FuncName)
				isFuncNameAssigned = true
			}
			return true
		case kind.Space:
			if isFuncNameAssigned && isFuncParamHandled {
				jumpTo = r.handleSingleReturnType(n, hooks...)
				skipAll = jumpTo == nil
				return true
			}
			return true
		case kind.ParenthesisLeft:
			isFuncParamHandled = true
			jumpTo = parenthesisResetter{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.CurlyBracketLeft:
			jumpTo = curlyBracketResetter{skip: true}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.SquareBracketLeft:
			isFuncParamHandled = true
			jumpTo = squareBracketResetter{skip: true}.Run(n, hooks...)
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
		[]kind.Kind{kind.NewLine, kind.CurlyBracketRight},
		findNodeOption{
			IsOutsideParenthesis:   true,
			IsOutsideCurlyBracket:  true,
			IsOutsideSquareBracket: true,
		},
	)

	if found.Kind() != kind.CurlyBracketRight {
		return false
	}

	return found.Next().Kind() == kind.ParenthesisLeft
}

// handleParameterFunc starts with 'func'
func (r funcResetter) handleParameterFunc(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("funcResetter.handleParameterFunc", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("funcResetter.handleParameterFunc.Returned")

	n := r.handleGeneralFunc(head, []kind.Kind{kind.Comma, kind.ParenthesisRight, kind.NewLine}, hooks...)
	if n.Kind() == kind.ParenthesisRight {
		handleHook(n, hooks...)
		return n.Next()
	}
	return n
}

// handleSingleReturnType starts with '\s', ends with '\n' and '\s' and '{'
func (r funcResetter) handleSingleReturnType(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("funcResetter.handleSingleReturnType", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("funcResetter.handleSingleReturnType.Returned")

	if head.Kind() != kind.Space {
		handleHook(head, hooks...)
		return head.Next()
	}

	head = head.Next() // skip first space

	found := head.findNext([]kind.Kind{kind.Comment}, findNodeOption{TargetReverse: true}, hooks...)
	switch found.Kind() {
	case kind.ParenthesisLeft:
		return parenthesisResetter{}.Run(head, hooks...)
	case kind.Func:
		return r.handleParameterFunc(head, hooks...)
	}

	var (
		buf []*Node
	)

	defer func() {
		if len(buf) != 0 {
			n := buf[0]
			next := buf[len(buf)-1].Next()
			n = n.CombineNext(kind.ParamType, buf[1:]...)
			n.ReplaceNext(next)
		}
	}()

	return head.IterNext(func(n *Node) bool {
		handleHook(n, hooks...)
		switch n.Kind() {
		case kind.NewLine, kind.Space, kind.CurlyBracketLeft: // return kind
			return false
		case kind.Comment:
			return true
		default:
			buf = helper.AppendUnrepeatable(buf, n)
			return true
		}
	})
}
