package goast

import (
	"github.com/yanun0323/goast/helper"
	"github.com/yanun0323/goast/kind"
)

// parenthesisResetter starts with '('
type parenthesisResetter struct {
	skip       bool
	isReceiver bool
}

func (r parenthesisResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("parenthesisResetter.Run", "\t\t....", head.DebugText(5))
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

// handleParenthesisParam starts with next of '(' and ','
func (r parenthesisResetter) handleParenthesisParam(head *Node, isReceiver bool, hooks ...func(*Node)) *Node {
	helper.DebugPrint("parenthesisResetter.handleParenthesisParam", "\t\t....", head.DebugText(5))
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
		case kind.Raw:
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

// handleParenthesisParam starts with next of '(' and ','
func (r parenthesisResetter) handleParenthesisParamCombined(head *Node, isReceiver bool, hooks ...func(*Node)) *Node {
	helper.DebugPrint("parenthesisResetter.handleParenthesisParamCombined", "\t\t....", head.DebugText(5))
	defer helper.DebugPrint("parenthesisResetter.handleParenthesisParamCombined.Returned")

	var (
		skipAll          bool
		jumpTo           *Node
		buf              []*Node
		hasSpaceAfterRaw bool
		hasName          bool
	)

	nameKind := kind.ParamName
	typeKind := kind.ParamType
	if isReceiver {
		nameKind = kind.MethodReceiverName
		typeKind = kind.MethodReceiverType
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
		case kind.Comma, kind.ParenthesisRight: // return kind
			return false
		case kind.SquareBracketLeft: // ignore including generic
			jumpTo = squareBracketResetter{}.Run(n, append(hooks, func(nn *Node) {
				buf = helper.AppendUnrepeatable(buf, nn)
			})...).Next()
			skipAll = jumpTo == nil
			return true
		case kind.Comment, kind.ParenthesisLeft, kind.Tab, kind.NewLine:
			return true
		case kind.Space:
			if len(buf) != 0 {
				hasSpaceAfterRaw = true
			}
			return true
		case kind.Func:
			jumpTo = funcResetter{isParameter: true}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		default:
			if hasSpaceAfterRaw {
				hasName = true
			}
			buf = helper.AppendUnrepeatable(buf, n)
			return true
		}
	})
}

// curlyBracketResetter starts with '{'
type curlyBracketResetter struct {
	skip bool
}

func (r curlyBracketResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("curlyBracketResetter.Run", "\t\t....", head.DebugText(5))
	defer helper.DebugPrint("curlyBracketResetter.Run.Returned")

	if r.skip {
		return head.skipNestNext(kind.CurlyBracketLeft, kind.CurlyBracketRight, hooks...)
	}

	// TODO: handle content
	return head.skipNestNext(kind.CurlyBracketLeft, kind.CurlyBracketRight, hooks...)
}

// squareBracketResetter
type squareBracketResetter struct {
	skip bool
}

func (r squareBracketResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("squareBracketResetter.Run", "\t\t....", head.DebugText(5))
	defer helper.DebugPrint("squareBracketResetter.Run.Returned")

	if r.skip {
		return head.skipNestNext(kind.SquareBracketLeft, kind.SquareBracketRight, hooks...)
	}

	// TODO: handle generic
	return head.skipNestNext(kind.SquareBracketLeft, kind.SquareBracketRight, hooks...)
}
