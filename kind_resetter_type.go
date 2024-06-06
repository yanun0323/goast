package goast

import (
	"github.com/yanun0323/goast/helper"
	"github.com/yanun0323/goast/kind"
)

// typeResetter head should be 'type'
//
// return key kind: '}' '\n'
type typeResetter struct{}

func (r typeResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("typeResetter.Run", "\t\t....", head.DebugText(5))
	defer helper.DebugPrint("typeResetter.Run.Returned")

	if head.Kind() != kind.Type {
		handleHook(head, hooks...)
		return head.Next()
	}

	var (
		isTypeNameAssigned   bool
		isSpaceAfterTypeName bool
		skipAll              bool
		jumpTo               *Node
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

		if !isTypeNameAssigned {
			switch n.Kind() {
			case kind.CurlyBracketRight, kind.NewLine: // return kind
				return false
			case kind.Raw:
				n.SetKind(kind.TypeName)
				isTypeNameAssigned = true
				return true
			default:
				return true
			}
		}

		if !isSpaceAfterTypeName {
			switch n.Kind() {
			case kind.Space:
				isSpaceAfterTypeName = true
			case kind.SquareBracketLeft:
				jumpTo = squareBracketResetter{skip: true}.Run(n, hooks...)
				skipAll = jumpTo == nil
				return true
			}
			return true
		}

		switch n.Kind() {
		case kind.Comment:
			return true
		case kind.Interface:
			jumpTo = interfaceResetter{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.Struct:
			jumpTo = structResetter{}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.Func:
			jumpTo = funcResetter{isNotMethod: true}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.Raw:
			jumpTo = paramResetter{resetKind: kind.ParamType, returnKinds: []kind.Kind{kind.NewLine}}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.SquareBracketLeft:
			jumpTo = squareBracketResetter{skip: true}.Run(n, hooks...)
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
	helper.DebugPrint("interfaceResetter.Run", "\t\t....", head.DebugText(5))
	defer helper.DebugPrint("interfaceResetter.Run.Returned")

	var (
		skipAll bool
		jumpTo  *Node
	)

	head = head.findNext([]kind.Kind{kind.CurlyBracketLeft}, findNodeOption{}, hooks...) // set head to '{'

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
		case kind.CurlyBracketRight: // return kind
			return false
		case kind.Tab, kind.Comment, kind.NewLine, kind.Space, kind.CurlyBracketLeft:
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
	helper.DebugPrint("structResetter.Run", "\t\t....", head.DebugText(5))
	defer helper.DebugPrint("structResetter.Run.Returned")

	var (
		skipAll bool
		jumpTo  *Node
	)

	head = head.findNext([]kind.Kind{kind.CurlyBracketLeft}, findNodeOption{}, hooks...).Next() // skip first of head to '{'

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
		case kind.CurlyBracketRight: // return kind
			return false
		case kind.Comment, kind.CurlyBracketLeft, kind.NewLine:
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
	helper.DebugPrint("funcResetter.handleStructRow", "\t\t....", head.DebugText(5))
	defer helper.DebugPrint("funcResetter.handleStructRow.Returned")
	
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
		case kind.NewLine: // return kind
			return false
		case kind.Comment, kind.NewLine, kind.Tab, kind.Space, kind.Comma:
			return true
		case kind.Raw:
			if paramNameCount != 0 {
				paramNameCount--
				n.SetKind(kind.ParamName)
				return true
			}
			jumpTo = paramResetter{resetKind: kind.ParamType, returnKinds: []kind.Kind{kind.NewLine}}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		default:
			jumpTo = paramResetter{resetKind: kind.ParamType, returnKinds: []kind.Kind{kind.NewLine}}.Run(n, hooks...)
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
		case kind.NewLine: // return kind
			if hasSpaceOrAfterRaw {
				nameCount = cacheNameCount
			}
			return false
		case kind.Comment, kind.Tab:
			return true
		case kind.Raw:
			rawCount++
			return true
		case kind.Space:
			if rawCount != 0 {
				if !hasComma && cacheNameCount == 0 {
					cacheNameCount++
				}
				hasSpaceOrAfterRaw = true
			}
			return true
		case kind.Comma:
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
