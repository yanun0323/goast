package goast

import (
	"fmt"
	"slices"
	"strings"

	"github.com/yanun0323/goast/charset"
	"github.com/yanun0323/goast/helper"
	"github.com/yanun0323/goast/kind"
)

/*
Node stands for any element in go language.
*/

// NewNode creates a new Node.
func NewNode(line int, text string, kinds ...kind.Kind) *Node {
	if len(kinds) != 0 {
		return &Node{line: line, kind: kinds[0], text: text}
	}

	return &Node{line: line, kind: kind.New(text), text: text}
}

// NewNodes creates a chain of nodes.
func NewNodes(startLine int, texts ...string) *Node {
	if len(texts) == 0 {
		return nil
	}

	head := NewNode(startLine, texts[0])
	cur := head
	for _, text := range texts[1:] {
		cur.InsertNext(NewNode(startLine, text))
		cur = cur.Next()
		startLine += strings.Count(text, "\n")
	}

	return head
}

// Node represents a structure.
type Node struct {
	line int
	kind kind.Kind
	text string

	prev *Node
	next *Node
}

// Copy copies node.
//
// If 'copyRelativeNode' equals true, then copy relative nodes too (prev and next).
// Otherwise only copies absolute node itself without prev and next.
func (n *Node) Copy(copyRelativeNode ...bool) *Node {
	if n == nil {
		return nil
	}

	if len(copyRelativeNode) != 0 && copyRelativeNode[0] {
		return n.deepCopy(nil, nil)
	}

	return &Node{
		line: n.line,
		kind: n.kind,
		text: n.text,
	}
}

func (n *Node) deepCopy(prev, next *Node) *Node {
	if n == nil {
		return nil
	}

	nn := &Node{
		line: n.line,
		kind: n.kind,
		text: n.text,
	}

	if prev == nil {
		prev = n.Prev().deepCopy(nil, nn)
	}

	if next == nil {
		next = n.Next().deepCopy(nn, nil)
	}

	nn.prev = prev
	nn.next = next

	return nn
}

func (n *Node) loop(iter func(*Node) *Node, fn func(*Node) bool) *Node {
	for nn := n; nn != nil; nn = iter(nn) {
		if !fn(nn) {
			return nn
		}
	}

	return nil
}

// IterPrev iterates over the previous nodes of the node.
func (n *Node) IterPrev(fn func(*Node) bool) *Node {
	return n.loop(func(nn *Node) *Node { return nn.Prev() }, fn)
}

// IterNext iterates over the next nodes of the node.
func (n *Node) IterNext(fn func(*Node) bool) *Node {
	return n.loop(func(nn *Node) *Node { return nn.Next() }, fn)
}

// Prev returns the previous node of the node.
func (n *Node) Prev() *Node {
	if n != nil {
		return n.prev
	}

	return nil
}

// Next returns the next node of the node.
func (n *Node) Next() *Node {
	if n != nil {
		return n.next
	}

	return nil
}

// InsertPrev inserts incoming node into current node's before,
// then returns old previous/next node of incoming node.
func (n *Node) InsertPrev(nn *Node) (insertedOldPrev *Node, insertedOldNext *Node) {
	if n == nn || n.Prev() == nn {
		return
	}

	oldPrev, oldNext := nn.Prev(), nn.Next()
	nPrev := n.Prev()

	n.setPrev(nn)
	nn.setNext(n)

	nn.setPrev(nPrev)
	nPrev.setNext(nn)

	oldPrev.setNext(nil)
	oldNext.setPrev(nil)

	return oldPrev, oldNext
}

// InsertNext inserts incoming node into current node's after,
// then returns old previous/next node of incoming node.
func (n *Node) InsertNext(nn *Node) (insertedOldPrev *Node, insertedOldNext *Node) {
	if n == nn || n.Next() == nn {
		return
	}

	oldPrev, oldNext := nn.Prev(), nn.Next()
	nNext := n.Next()

	n.setNext(nn)
	nn.setPrev(n)

	nn.setNext(nNext)
	nNext.setPrev(nn)

	oldPrev.setNext(nil)
	oldNext.setPrev(nil)

	return oldPrev, oldNext
}

// ReplacePrev replaces incoming node into current node's before,
// then returns the old previous node of current node and
// the old next node of incoming node.
func (n *Node) ReplacePrev(nn *Node) (currentOldPrev *Node, replacedOldNext *Node) {
	if n == nn || n.Prev() == nn {
		return
	}

	oldPrev, oldNext := n.Prev(), nn.Next()

	n.setPrev(nn)
	nn.setNext(n)

	oldPrev.setNext(nil)
	oldNext.setPrev(nil)

	return oldPrev, oldNext
}

// ReplaceNext replaces incoming node into current node's after,
// then returns the old previous node of incoming node and
// the old next node of current node.
func (n *Node) ReplaceNext(nn *Node) (replacedOldPrev *Node, currentOldNext *Node) {
	if n == nn || n.Next() == nn {
		return
	}

	oldPrev, oldNext := nn.Prev(), n.Next()

	n.setNext(nn)
	nn.setPrev(n)

	oldPrev.setNext(nil)
	oldNext.setPrev(nil)

	return oldPrev, oldNext
}

// TakePrev takes the single node before current node, and connects leftover node.
func (n *Node) TakePrev() *Node {
	removed := n.Prev()
	newPrev := removed.Prev()

	n.setPrev(newPrev)
	newPrev.setNext(n)

	removed.setNext(nil)
	removed.setPrev(nil)

	return removed
}

// TakeNext takes the single node after current node, and connects leftover node.
func (n *Node) TakeNext() *Node {
	removed := n.Next()
	newNext := removed.Next()

	n.setNext(newNext)
	newNext.setPrev(n)

	removed.setNext(nil)
	removed.setPrev(nil)

	return removed
}

// RemovePrev removes node previous current node, and returns removed node with it's all previous node.
func (n *Node) RemovePrev() *Node {
	removed := n.Prev()
	n.setPrev(nil)
	removed.setNext(nil)

	return removed
}

// RemoveNext removes node next current node, and returns removed node with it's all next node.
func (n *Node) RemoveNext() *Node {
	removed := n.Next()
	n.setNext(nil)
	removed.setPrev(nil)

	return removed
}

// CombinePrev combines all values of incoming nodes into current node,
// and returns new node.
//
// e.g. (nn3 - nn2 - nn1 - n)
func (n *Node) CombinePrev(k kind.Kind, nns ...*Node) *Node {
	n.SetKind(k)

	if len(nns) == 0 {
		return n
	}

	if n == nil {
		n = nns[0]
		nns = nns[1:]
	}

	buf := make([]string, 0, len(nns)+1)
	buf = append(buf, n.Text())
	for _, nn := range nns {
		if nn == nil || n == nn {
			continue
		}

		buf = append(buf, nn.Text())
		// make node recyclable to gc
		nn.Isolate()
	}

	slices.Reverse(buf)

	n.text = strings.Join(buf, "")

	return n
}

// CombineNext combines all values of incoming nodes into current node,
// and returns new node.
//
// e.g. (n - nn1 - nn2 - nn3)
func (n *Node) CombineNext(k kind.Kind, nns ...*Node) *Node {
	n.SetKind(k)

	if len(nns) == 0 {
		return n
	}

	if n == nil {
		n = nns[0]
		nns = nns[1:]
	}

	buf := make([]string, 0, len(nns)+1)
	buf = append(buf, n.Text())
	for _, nn := range nns {
		if nn == nil || n == nn {
			continue
		}

		buf = append(buf, nn.Text())
		// make node recyclable to gc
		nn.Isolate()
	}

	n.text = strings.Join(buf, "")

	return n
}

// Isolated disconnects current node from others.
func (n *Node) Isolate() {
	n.Prev().setNext(nil)
	n.Next().setPrev(nil)

	n.setPrev(nil)
	n.setNext(nil)
}

// Line returns the line number of the node.
func (n *Node) Line() int {
	if n == nil {
		return -2
	}

	return n.line
}

// Kind returns the kind of the node.
func (n *Node) Kind() kind.Kind {
	if n == nil {
		return kind.None
	}

	return n.kind
}

// SetKind sets the kind of the node.
func (n *Node) SetKind(k kind.Kind) {
	if n != nil {
		n.kind = k
	}
}

// Text returns the raw text of the node.
func (n *Node) Text() string {
	if n == nil {
		return ""
	}
	return n.text
}

func (n *Node) debugText(nextCount ...uint) string {
	if n == nil {
		return ""
	}

	var (
		buf   strings.Builder
		count uint = 1
	)

	if len(nextCount) != 0 {
		count += nextCount[0]
	}

	n.IterNext(func(n *Node) bool {
		buf.WriteString(helper.TidyText(n.Text()))
		count--
		return count != 0
	})
	return buf.String()
}

// SetText sets the raw text of the node.
func (n *Node) SetText(text string) {
	if n == nil {
		return
	}
	n.text = text
}

// Description returns node description.
//
// It returns "<nil>" if node is nil.
func (n *Node) Description() string {
	if n == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%d .... %s .... *Node.%s", n.Line()+1, helper.TidyText(n.Text()), n.Kind().String())
}

// Print prints the description of the node.
func (n *Node) Print() {
	println("\t", n.Description())
}

func (n *Node) setPrev(nn *Node) {
	if n == nn || n == nil {
		return
	}

	// TODO: how to prevent leak (no gc)
	if n.prev != nil {
		n.prev.next = nil
	}

	n.prev = nn
}

func (n *Node) setNext(nn *Node) {
	if n == nn || n == nil {
		return
	}

	// TODO: how to prevent leak (no gc)
	if n.next != nil {
		n.next.prev = nil
	}

	n.next = nn
}

// skipNestNext helper
func (n *Node) skipNestNext(nestLeft, nestRight kind.Kind, hooks ...func(*Node)) *Node {
	count := 1
	if n.Kind() == nestLeft {
		handleHook(n, hooks...)
		n = n.findNext([]kind.Kind{nestLeft}, findNodeOption{}, hooks...).Next() // skip first nestLeft
	}

	return n.IterNext(func(n *Node) bool {
		handleHook(n, hooks...)

		switch n.Kind() {
		case nestLeft:
			count++
			return true
		case nestRight:
			count--
			return count != 0
		default:
			return true
		}
	})
}

// findNext helper
func (n *Node) findNext(
	target []kind.Kind,
	opt findNodeOption,
	hooks ...func(*Node),
) *Node {
	var (
		parenthesisLeftCount   int
		curlyBracketLeftCount  int
		squareBracketLeftCount int

		targetKindSet = charset.New(target...)
	)

	return n.IterNext(func(n *Node) bool {
		handleHook(n, hooks...)
		switch n.Kind() {
		case kind.ParenthesisLeft:
			parenthesisLeftCount++
		case kind.CurlyBracketLeft:
			curlyBracketLeftCount++
		case kind.SquareBracketLeft:
			squareBracketLeftCount++
		case kind.ParenthesisRight:
			parenthesisLeftCount--
		case kind.CurlyBracketRight:
			curlyBracketLeftCount--
		case kind.SquareBracketRight:
			squareBracketLeftCount--
		}

		if opt.IsInsideParenthesis && parenthesisLeftCount == 0 {
			return true
		}

		if opt.IsInsideCurlyBracket && curlyBracketLeftCount == 0 {
			return true
		}

		if opt.IsInsideSquareBracket && squareBracketLeftCount == 0 {
			return true
		}

		if opt.IsOutsideParenthesis && parenthesisLeftCount != 0 {
			return true
		}

		if opt.IsOutsideCurlyBracket && curlyBracketLeftCount != 0 {
			return true
		}

		if opt.IsOutsideSquareBracket && squareBracketLeftCount != 0 {
			return true
		}

		if targetKindSet.Contain(n.Kind()) {
			return opt.TargetReverse
		}

		return !opt.TargetReverse
	})
}

type findNodeOption struct {
	TargetReverse          bool
	IsInsideParenthesis    bool
	IsInsideCurlyBracket   bool
	IsInsideSquareBracket  bool
	IsOutsideParenthesis   bool
	IsOutsideCurlyBracket  bool
	IsOutsideSquareBracket bool
}
