package goast

import (
	"github.com/yanun0323/goast/charset"
	"github.com/yanun0323/goast/helper"
	"github.com/yanun0323/goast/kind"
)

func resetKind(n *Node) *Node {
	switch n.Kind() {
	case kind.Func:
		return funcResetter{isFuncKeywordLeading: true}.Run(n)
	case kind.Type:
		return typeResetter{}.Run(n)
	}
	return nil
}

type kindResetter interface {
	Run(*Node, ...func(*Node)) *Node
}

func handleHook(n *Node, hooks ...func(*Node)) {
	for _, hook := range hooks {
		hook(n)
	}
}

// func generalResetter() kindResetter {
// 	return &genericResetter{
// 		DeeperResetterTable: map[kind.Kind]kindResetter{
// 			kind.Func: funcResetter{isFuncKeywordLeading: true},
// 			kind.Type: typeResetter{},
// 			kind.
// 		},
// 	}
// }

// genericResetter stands for COMMON genericResetter
// type genericResetter struct {
// 	TriggerKind         charset.Set[kind.Kind]
// 	TriggerLimit        int
// 	KindChangeTable     map[int]kind.Kind
// 	ChangeableKind      charset.Set[kind.Kind]
// 	ReturnKind          charset.Set[kind.Kind]
// 	DeeperResetterTable map[kind.Kind]kindResetter
// }

// func (r genericResetter) Run(head *Node, _ ...func(*Node)) *Node {
// 	helper.DebugPrint("genericResetter.Run", "\t\t....", head.DebugText(5))
// 	defer helper.DebugPrint("genericResetter.Run.Returned")

// 	var (
// 		skipAll      bool
// 		jumpTo       *Node
// 		triggered    bool
// 		triggerIndex int

// 		triggerLimit = r.TriggerLimit
// 	)

// 	return head.IterNext(func(n *Node) bool {
// 		if skipAll {
// 			return true
// 		}

// 		if jumpTo != nil {
// 			if jumpTo != n {
// 				return true
// 			}
// 			jumpTo = nil
// 		}

// 		//  ReturnKind > deeperResetter > TriggerKind > ChangeableKind > UnchangeableKind
// 		if r.ReturnKind.Contain(n.Kind()) {
// 			return false
// 		}

// 		if resetter, ok := r.DeeperResetterTable[n.Kind()]; ok && resetter != nil {
// 			jumpTo = resetter.Run(n)
// 			skipAll = jumpTo == nil
// 			return true
// 		}

// 		if r.TriggerKind.Contain(n.Kind()) {
// 			triggered = true
// 			triggerIndex = 0
// 			return true
// 		}

// 		if triggerLimit == 0 || !triggered || !r.ChangeableKind.Contain(n.Kind()) {
// 			return true
// 		}

// 		triggerIndex++
// 		if triggerLimit > 0 {
// 			triggerLimit--
// 		}

// 		if k, ok := r.KindChangeTable[triggerIndex]; ok && k != kind.None {
// 			n.SetKind(k)
// 		}

// 		return true
// 	})
// }

type paramResetter struct {
	skip        bool
	resetKind   kind.Kind
	returnKinds []kind.Kind
}

func (r paramResetter) Run(head *Node, hooks ...func(*Node)) *Node {
	helper.DebugPrint("paramResetter.Run", "\t\t....", head.DebugText(5))
	defer helper.DebugPrint("paramResetter.Run.Returned")

	var (
		skipAll bool
		jumpTo  *Node
		buf     []*Node
	)

	if len(r.returnKinds) == 0 {
		handleHook(head, hooks...)
		return head.Next()
	}

	returnKindSet := charset.New(r.returnKinds...)

	defer func() {
		if len(buf) != 0 {
			n := buf[0]
			next := buf[len(buf)-1].Next()
			n = n.CombineNext(r.resetKind, buf[1:]...)
			n.ReplaceNext(next)
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

		if returnKindSet.Contain(n.Kind()) { // return kind
			return false
		}

		hooksCopy := make([]func(*Node), len(hooks))
		copy(hooksCopy, hooks)
		appendedHooks := append(hooksCopy, func(nn *Node) {
			buf = helper.AppendUnrepeatable(buf, nn)
		})

		switch n.Kind() {
		case kind.ParenthesisLeft:
			jumpTo = parenthesisResetter{skip: r.skip}.Run(n, appendedHooks...)
			skipAll = jumpTo == nil
			return true
		case kind.CurlyBracketLeft:
			jumpTo = curlyBracketResetter{skip: r.skip}.Run(n, appendedHooks...)
			skipAll = jumpTo == nil
			return true
		case kind.SquareBracketLeft:
			jumpTo = squareBracketResetter{skip: r.skip}.Run(n, appendedHooks...)
			skipAll = jumpTo == nil
			return true
		case kind.Func:
			jumpTo = funcResetter{isNotMethod: true}.Run(n, hooks...)
			skipAll = jumpTo == nil
			return true
		case kind.Comment:
			return true
		default:
			if !r.skip {
				buf = helper.AppendUnrepeatable(buf, n)
			}

			return true
		}
	})
}
