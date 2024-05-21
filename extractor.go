package goast

import (
	"errors"
	"strings"
)

var ErrOutOfRange = errors.New("out of range")

func extract(text []byte) ([]Node, error) {
	var i, line int
	return _commonExtractor.Run(text, &i, &line)
}

var (
	_commonExtractor = extractor{
		ReturnCharset:       nil,
		SeparatorCharset:    _separatorCharset,
		CommentKeyword:      _commentKeywordTrie,
		InnerCommentKeyword: _innerCommentKeywordTrie,
	}
	_parenthesisExtractor = extractor{
		ReturnCharset:       newCharset[byte](')'),
		SeparatorCharset:    _separatorCharset,
		CommentKeyword:      _commentKeywordTrie,
		InnerCommentKeyword: _innerCommentKeywordTrie,
	}

	_curlyBracketExtractor = extractor{
		ReturnCharset:       newCharset[byte]('}'),
		SeparatorCharset:    _separatorCharset,
		CommentKeyword:      _commentKeywordTrie,
		InnerCommentKeyword: _innerCommentKeywordTrie,
	}

	_commentKeywordTrie = newTrie(map[string]bool{"//": true})

	_innerCommentKeywordTrie = newTrie(map[string]trie[bool]{
		"/*": newTrie(map[string]bool{"*/": true}),
	})
	_deeperExtract = map[byte]extractor{
		'(': _parenthesisExtractor,
		'{': _curlyBracketExtractor,
	}
)

type extractor struct {
	ReturnCharset       charset[byte]
	SeparatorCharset    charset[byte]
	CommentKeyword      trie[bool]
	InnerCommentKeyword trie[trie[bool]]
}

func (e extractor) Run(text []byte, i *int, line *int) ([]Node, error) {
	var (
		char         byte
		result       []Node
		buf          strings.Builder
		comment      bool
		innerComment trie[bool]
	)

	tryBuf2Node := func(kind ...Kind) {
		if buf.Len() == 0 {
			return
		}
		result = append(result, NewNode(*line, buf.String(), kind...))
		buf.Reset()
	}

	for ; *i < len(text); *i++ {
		char = text[*i]

		if comment {
			// comment text
			if char == '\n' {
				// comment end
				tryBuf2Node(KindComment)
				result = append(result, NewNode(*line, string(char)))
				*line++
				comment = false
			} else {
				buf.WriteByte(char)
			}
			continue
		}

		if innerComment != nil {
			// inner comment text
			buf.WriteByte(char)
			if fit, ok := innerComment.FindByte(text[*i-1 : *i+1]); ok && fit {
				// inner comment end
				tryBuf2Node(KindComment)
				innerComment = nil
			}
			continue
		}

		if *i+2 <= len(text) {
			if fit, ok := e.CommentKeyword.FindByte(text[*i : *i+2]); ok && fit {
				// comment start
				tryBuf2Node()
				comment = true
				buf.WriteByte(char)
				continue
			}

			if set, ok := e.InnerCommentKeyword.FindByte(text[*i : *i+2]); ok && set != nil {
				// inner comment start
				tryBuf2Node()
				innerComment = set
				buf.WriteByte(char)
				continue
			}
		}

		if e, ok := _deeperExtract[char]; ok {
			// () or {}
			tryBuf2Node()
			result = append(result, NewNode(*line, string(char)))
			*i++
			ns, err := e.Run(text, i, line)
			if err != nil {
				return nil, err
			}
			result = append(result, ns...)
			continue
		}

		if e.SeparatorCharset.Contain(char) {
			// ' '
			tryBuf2Node()
			result = append(result, NewNode(*line, string(char)))
		} else {
			buf.WriteByte(char)
		}

		if e.ReturnCharset.Contain(char) {
			// inside ) or }
			tryBuf2Node()
			return result, nil
		}

		if char == '\n' {
			*line++
		}
	}
	tryBuf2Node()
	return result, nil
}
