package goast

import (
	"errors"
	"strings"
)

var ErrOutOfRange = errors.New("out of range")

// extract parses file text content into nodes.
func extract(text []byte) (*Node, error) {
	var (
		i, line int
		buf     strings.Builder
	)
	node, err := _commonExtractor.Run(text, &buf, &i, &line)
	if err != nil {
		return nil, err
	}

	return node, nil
}

var (
	_commonExtractor = &extractor{
		SeparatorCharset:  _separatorCharset,
		ReturnKeyword:     "",
		SkipReturnKeyword: "",
	}

	_parenthesisExtractor = &extractor{
		kind:              KindKeywords,
		SeparatorCharset:  _separatorCharset,
		ReturnKeyword:     ")",
		SkipReturnKeyword: "",
	}

	_curlyBracketExtractor = &extractor{
		kind:              KindKeywords,
		SeparatorCharset:  _separatorCharset,
		ReturnKeyword:     "}",
		SkipReturnKeyword: "",
	}

	_commentExtractor = &extractor{
		kind:              KindComments,
		IncludeOpen:       true,
		SeparatorCharset:  nil,
		ReturnKeyword:     "\n",
		SkipReturnKeyword: "",
	}

	_innerCommentExtractor = &extractor{
		kind:              KindComments,
		IncludeOpen:       true,
		IncludeClose:      true,
		SeparatorCharset:  nil,
		ReturnKeyword:     "*/",
		SkipReturnKeyword: "",
	}

	_stringExtractor = &extractor{
		kind:              KindString,
		IncludeOpen:       true,
		IncludeClose:      true,
		SeparatorCharset:  nil,
		ReturnKeyword:     "\"",
		SkipReturnKeyword: "\\\"",
	}

	_multilineStringExtractor = &extractor{
		kind:              KindString,
		IncludeOpen:       true,
		IncludeClose:      true,
		SeparatorCharset:  nil,
		ReturnKeyword:     "`",
		SkipReturnKeyword: "",
	}

	_deeperExtractTable = map[*extractor]deeperExtract{
		_commonExtractor:       _commonDeeperExtract,
		_parenthesisExtractor:  _commonDeeperExtract,
		_curlyBracketExtractor: _commonDeeperExtract,
	}
	_commonDeeperExtract = deeperExtract{
		"(":  _parenthesisExtractor,
		"{":  _curlyBracketExtractor,
		"//": _commentExtractor,
		"/*": _innerCommentExtractor,
		"\"": _stringExtractor,
		"`":  _multilineStringExtractor,
	}
)

type deeperExtract map[string]*extractor

func (de deeperExtract) PrefixFit(s []byte) (*extractor, bool) {
	if de == nil {
		return nil, false
	}

	for k, v := range de {
		if hasPrefix(s, k) {
			return v, true
		}
	}
	return nil, false
}

type extractor struct {
	kind              Kind
	IncludeOpen       bool
	IncludeClose      bool
	SeparatorCharset  charset[byte]
	ReturnKeyword     string
	SkipReturnKeyword string
}

func (e *extractor) Run(text []byte, buf *strings.Builder, i *int, line *int) (*Node, error) {
	if e == nil {
		return nil, nil
	}

	var (
		char      byte
		head, cur *Node
	)

	bufLine := *line

	insertNode := func(n *Node) {
		if n == nil {
			return
		}

		if head == nil {
			head = n
		}

		if cur != nil {
			cur.ReplaceNext(n)
		} else {
			cur = n
		}

		cur = n.IterNext(func(n *Node) bool {
			return n.Next() != nil
		})
	}

	pushNode := func(useLine bool, kind ...Kind) {
		if buf.Len() == 0 {
			return
		}
		if useLine {
			insertNode(NewNode(*line, buf.String(), kind...))
		} else {
			insertNode(NewNode(bufLine, buf.String(), kind...))
		}
		buf.Reset()
		bufLine = *line
	}

	lineStep := func() {
		if text[*i] == '\n' {
			*line++
		}
	}

	for ; *i < len(text); *i++ {
		char = text[*i]

		if ee, ok := e.DeeperExtract().PrefixFit(text[*i:]); ok {
			// () {} /**/ "" `` //\n
			pushNode(true)
			if ee != nil && ee.IncludeOpen { // /* // " `
				buf.WriteByte(char)
			} else { // ( {
				insertNode(NewNode(bufLine, string(char)))
			}
			*i++
			ns, err := ee.Run(text, buf, i, line)
			if err != nil {
				return nil, err
			}
			insertNode(ns)
			continue
		}

		trailing := text[:*i+1]
		if len(e.SkipReturnKeyword) != 0 && hasSuffix(trailing, e.SkipReturnKeyword) {
			// skip return
			buf.WriteByte(char)
			lineStep()
			continue
		}

		if len(e.ReturnKeyword) != 0 && hasSuffix(trailing, e.ReturnKeyword) {
			// inside ) } */ " ` \n
			if e.IncludeClose { // */ " `
				buf.WriteByte(char)
				pushNode(false, e.Kind()...)
			} else {
				pushNode(true, e.Kind()...)
				insertNode(NewNode(*line, string(char)))
			}
			lineStep()
			return head, nil
		}

		if e.SeparatorCharset.Contain(char) {
			// ' '
			pushNode(true)
			insertNode(NewNode(*line, string(char)))
		} else {
			buf.WriteByte(char)
		}

		lineStep()
	}
	pushNode(true)
	return head, nil
}

func (e *extractor) Kind() []Kind {
	if e == nil || e.kind == KindRaws {
		return nil
	}

	return []Kind{e.kind}
}

func (e *extractor) DeeperExtract() deeperExtract {
	return _deeperExtractTable[e]
}
