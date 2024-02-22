package goast /* 123 */

import (
	"go/format"
	"io"
	"log"
	"os"
	"strings"
)

func ParseFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	buf, err = format.Source(buf)
	if err != nil {
		return nil, err
	}

	file := &File{
		Path: path,
	}

	var cacheComment *unit
	tryCleanCacheComment := func() {
		if cacheComment == nil {
			return
		}
		file.Nodes = append(file.Nodes, &Node{
			Values: []*unit{cacheComment},
		})
		cacheComment = nil
	}

	rows := strings.Split(string(buf), "\n")
	for i := 0; i < len(rows); i++ {
		row := strings.TrimSpace(rows[i])

		if strings.HasPrefix(row, "package") {
			file.Package = extractPackage(row)
			continue
		}

		if strings.HasPrefix(row, "import") {
			file.Imports, i = extractImport(rows, i)
			continue
		}

		if strings.HasPrefix(row, "/") {
			tryCleanCacheComment()
			cacheComment, i = extractComment(rows, i)
			continue

		}

		if strings.HasPrefix(row, "const") {

		}

		if strings.HasPrefix(row, "var") {

		}

		if strings.HasPrefix(row, "type") {

		}

		if strings.HasPrefix(row, "func") {

		}

	}

	return file, nil
}

func tryCombineComment(spans []string, idx int) (string, int) {
	switch spans[idx] {
	case "//":
		result := make([]string, 0, len(spans))
		for i := idx; i < len(spans); i++ {
			result = append(result, spans[i])
			idx = i
		}
		return strings.Join(result, " "), idx
	case "/*":
		result := make([]string, 0, len(spans))
		for i := idx; i < len(spans); i++ {
			result = append(result, spans[i])
			idx = i
			if spans[i] == "*/" {
				break
			}
		}
		return strings.Join(result, " "), idx
	default:
		return spans[idx], idx
	}
}

// tryParsingComplexComponent parses string to Node when type is map/array/slice/multiline-string/struct/interface
func tryParsingComplexComponent(rows []string, idx int) (*Node, int) {
	spans := strings.Split(strings.TrimSpace(rows[idx]), " ")
	root := &Node{}
	complexComponent := false
	structIndex := -1
	p := -1
	openBucket := ""
	for i, span := range spans {
		span = strings.TrimSpace(span)
		t := parsingType(span, func(s string) (Type, bool) {
			if len(s) >= 1 && s[0] == '`' {
				return Special, true /* Special for multiline string */
			}
			return Raw, false
		})
		root.Values = append(root.Values, Unit(i, t, span))
		switch t {
		case Map, Slice, Array, Special, Structure, Interface:
			complexComponent = true
			structIndex = i
		case Keyword:
			switch span {
			case "(":
			case ")":
			case "{":
			case "}":
			}
		}
	}

	if !complexComponent {
		return root, idx
	}
	// start parsing complex component

}

func extractPackage(row string) []*unit {
	result := []*unit{}
	spans := strings.Split(row, " ")
	idx := 0
	span := ""
	for i := 0; i < len(spans); i, idx = i+1, idx+1 {
		span, i = tryCombineComment(spans, i)
		result = append(result, &unit{
			Index: idx,
			Type:  parsingType(span),
			Value: span,
		})
	}
	return result
}

func extractImport(rows []string, idx int) ([][]*unit, int) {
	fnRowParser := func(row string) []*unit {
		result := []*unit{}
		spans := strings.Split(row, " ")
		idx := 0
		span := ""
		for i := 0; i < len(spans); i, idx = i+1, idx+1 {
			span, i = tryCombineComment(spans, i)
			result = append(result, &unit{
				Index: idx,
				Type: parsingType(span, func(s string) (Type, bool) {
					switch s[0] {
					case '"':
						return Raw, true
					case '/':
						return Comment, true
					default:
						return Keyword, true
					}
				}),
				Value: span,
			})
		}
		return result
	}

	row := strings.TrimSpace(rows[idx])
	if !strings.Contains(row, "(") {
		return [][]*unit{fnRowParser(row)}, idx
	}

	result := [][]*unit{}
	for i := idx + 1; i < len(rows); i++ {
		row := strings.TrimSpace(rows[i])
		if len(row) == 0 {
			continue
		}

		if row[0] == ')' {
			break
		}

		result = append(result, fnRowParser(row))
		idx = i
	}
	return result, idx
}

func extractComment(rows []string, idx int) (*unit, int) {
	row := strings.TrimSpace(rows[idx])
	if len(row) <= 1 {
		log.Fatalf("extract comment, err: %s", row)
	}

	comments := []string{}
	switch string(row[:2]) {
	case "//":
		for i := idx; i < len(rows); i++ {
			row := strings.TrimSpace(rows[i])
			if len(row) <= 1 || string(row[:2]) != "//" {
				break
			}
			comments = append(comments, row)
			idx = i
		}
	case "/*":
		for i := idx; i < len(rows); i++ {
			row := strings.TrimSpace(rows[i])
			if len(row) <= 1 {
				break
			}
			comments = append(comments, row)
			idx = i
			if string(row[:2]) == "*/" {
				break
			}
		}
	}
	return &unit{
		Index: 0,
		Type:  Comment,
		Value: strings.Join(comments, "\n"),
	}, idx
}

func extractConst(rows []string, idx int) (*Node, int) {
	fnRowParser := func(rows []string, i int) (*Node, int) {
		result := &Node{}
		for i, span := range strings.Split(rows[i], " ") {
			result.Values = append(result.Values, &unit{
				Index: i,
				Type:  parsingType(span),
				Value: span,
			})
		}

		openComment := false
		openMultilineString := false
		last := result.Values[len(result.Values)-1]

		switch last.Type {
		case Comment:
			if strings.HasPrefix(last.Value, "/*") && !strings.HasSuffix(last.Value, "*/") {
				openComment = true
			}
		case String:
			if len(last.Value) != 0 && last.Value[0] == '`' && last.Value[len(last.Value)-1] != '`' {
				openMultilineString = true
			}
		}

		for i := i; i < len(rows); i++ {
			if openComment {

			}

			if openMultilineString {

			}
		}

		return result, i
	}

	var result *Node
	row := strings.TrimSpace(rows[idx])
	if len(row) <= 1 {
		log.Fatal("extract const error")
	}

	if !strings.Contains(row, "(") {
		// TODO: Implement me
		_ = fnRowParser
	}

	return result, idx
}
