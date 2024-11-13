package goast

import (
	"github.com/yanun0323/goast/helper"
	"github.com/yanun0323/goast/kind"
)

// parenthesisResetter starts with '(', ends with ')'
type parenthesisResetter struct {
	skip       bool
	isReceiver bool
}

func (r parenthesisResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("parenthesisResetter.Run", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("parenthesisResetter.Run.Returned")

	if r.skip {
		return head.skipNestNext(kind.ParenthesisLeft, kind.ParenthesisRight, hooks...)
	}

	if head.Kind() != kind.ParenthesisLeft {
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
		case kind.ParenthesisRight: // return kind
			return false
		case kind.Comment, kind.Comma, kind.ParenthesisLeft, kind.Tab, kind.NewLine, kind.Space:
			return true
		default:
			jumpTo = r.handleParenthesisParam(n, r.isReceiver, hooks...)
			skipAll = jumpTo == nil
			return true
		}
	})
}

// handleParenthesisParam starts with next of '(' and ',', ends with ',' and ')'
func (r parenthesisResetter) handleParenthesisParam(head *Node, isReceiver bool, hooks ...func(*Node)) *Node {
	helper.DebugPrint("parenthesisResetter.handleParenthesisParam", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("parenthesisResetter.handleParenthesisParam.Returned")

	var (
		skipAll bool
		jumpTo  *Node

		firstRawHandled bool
	)

	nameKind := kind.ParamName
	typeKind := kind.ParamType
	if isReceiver {
		nameKind = kind.MethodReceiverName
		typeKind = kind.MethodReceiverType
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

		switch n.Kind() {
		case kind.Comma, kind.ParenthesisRight: // return kind
			return false
		case kind.SquareBracketLeft: // ignore including generic
			jumpTo = squareBracketResetter{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.Comment, kind.ParenthesisLeft, kind.Tab, kind.NewLine, kind.Space:
			return true
		case kind.Func:
			jumpTo = funcResetter{isParameter: true}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.Raw, kind.ParamType /* for struct resetter */, kind.ParamName /* for struct resetter */ :
			if !firstRawHandled {
				if n.Next().Kind() == kind.Space {
					n.SetKind(nameKind)
				} else {
					n.SetKind(typeKind)
				}
				firstRawHandled = true
				return true
			}
			n.SetKind(typeKind)
			return true
		default:
			return true
		}
	})
}

// curlyBracketResetter starts with '{', ends with '}'
type curlyBracketResetter struct {
	skip bool
}

func (r curlyBracketResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("curlyBracketResetter.Run", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("curlyBracketResetter.Run.Returned")

	if r.skip {
		return head.skipNestNext(kind.CurlyBracketLeft, kind.CurlyBracketRight, hooks...)
	}

	// TODO: handle content
	return head.skipNestNext(kind.CurlyBracketLeft, kind.CurlyBracketRight, hooks...)
}

// squareBracketResetter starts with '[', ends with ']'
type squareBracketResetter struct {
	skip bool
}

func (r squareBracketResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("squareBracketResetter.Run", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("squareBracketResetter.Run.Returned")

	if r.skip {
		return head.skipNestNext(kind.SquareBracketLeft, kind.SquareBracketRight, hooks...)
	}

	// TODO: handle generic
	return head.skipNestNext(kind.SquareBracketLeft, kind.SquareBracketRight, hooks...)
}

// TODO: Refactor parentheses resetters

// parenthesisResetter starts with '(', ends with ')'
type parenthesisResetter2 struct {
	skip           bool
	isFuncReceiver bool
}

func (r parenthesisResetter2) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("parenthesisResetter.Run", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("parenthesisResetter.Run.Returned")

	if r.skip {
		return head.skipNestNext(kind.ParenthesisLeft, kind.ParenthesisRight, hooks...)
	}

	if head.Kind() != kind.ParenthesisLeft {
		handleHook(head, hooks...)
		return head.Next()
	}

	handleHook(head, hooks...)
	head = head.Next() // skip first '('

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
		case kind.ParenthesisRight: // return kind
			return false
		case kind.ParenthesisLeft:
			jumpTo = parenthesisResetter2{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.Comment, kind.Comma, kind.Tab, kind.NewLine, kind.Space:
			return true
		default:
			jumpTo = r.handleParenthesisParam(n, r.isFuncReceiver, hooks...)
			skipAll = jumpTo == nil
			return true
		}

	})
}

// handleParenthesisParam starts with next of '(' and ',', ends with ',' and ')'
func (r parenthesisResetter2) handleParenthesisParam(head *Node, isReceiver bool, hooks ...func(*Node)) *Node {
	helper.DebugPrint("parenthesisResetter.handleParenthesisParam", "\t\t....", head.debugText(5))
	defer helper.DebugPrint("parenthesisResetter.handleParenthesisParam.Returned")

	collection := []*Node{}
	returned := head.IterNext(func(n *Node) bool {
		switch n.Kind() {
		case kind.Comma, kind.ParenthesisRight:
			return false
		case kind.Space:
			if len(collection) != 0 {
				collection = append(collection, n)
			}

			return true
		case kind.Tab, kind.Comment:
			return true
		default:
			collection = append(collection, n)
			return true
		}
	})

	return returned
}
