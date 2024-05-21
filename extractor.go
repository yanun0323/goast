package goast

import (
	"strings"
)

func extract(text []byte) ([]Node, error) {
	var i, line int
	return _commonExtractor.Run(text, &i, &line)
}

var (
	_commonExtractor = extractor{
		ReturnCharset:    nil,
		SeparatorCharset: _separatorCharset,
		CommentKeyword:   _commentKeywordTrie,
	}
	_comboExtractor = extractor{
		ReturnCharset:    newCharset[byte](')'),
		SeparatorCharset: _separatorCharset,
		CommentKeyword:   _commentKeywordTrie,
	}

	_contextExtractor = extractor{
		ReturnCharset:    newCharset[byte]('}'),
		SeparatorCharset: _separatorCharset,
		CommentKeyword:   _commentKeywordTrie,
	}

	_commentKeywordTrie = newTrie(map[string]trie[bool]{
		"//": newTrie(map[string]bool{"\n": true}),
		"/*": newTrie(map[string]bool{"*/": true}),
	})
	_deeperExtract = map[byte]extractor{
		'(': _comboExtractor,
		'{': _contextExtractor,
	}
)

type extractor struct {
	ReturnCharset    charset[byte]
	SeparatorCharset charset[byte]
	CommentKeyword   trie[trie[bool]]
}

func (e extractor) Run(text []byte, i *int, line *int) ([]Node, error) {
	var (
		char    byte
		result  []Node
		buf     strings.Builder
		comment trie[bool]
	)

	tryBuf2Node := func(kind ...ElementKind) {
		if buf.Len() == 0 {
			return
		}
		result = append(result, NewElement(*line, buf.String(), kind...))
		buf.Reset()
	}

	for ; *i < len(text); *i++ {
		char = text[*i]

		if comment != nil {
			// comment text
			buf.WriteByte(char)
			if fit, ok := comment.FindByte(text[*i:]); ok && fit {
				tryBuf2Node(ElemComment)
				comment = nil
			}
			if char == '\n' {
				*line++
			}
			continue
		}

		if set, ok := e.CommentKeyword.FindByte(text[*i:]); ok && set != nil {
			// comment start
			tryBuf2Node()
			comment = set
			buf.WriteByte(char)
			continue
		}

		if e, ok := _deeperExtract[char]; ok {
			// () or {}
			tryBuf2Node()
			result = append(result, NewElement(*line, string(char)))
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
			result = append(result, NewElement(*line, string(char)))
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
