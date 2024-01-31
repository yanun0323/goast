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

	// var commentCache *unit

	rows := strings.Split(string(buf), "\n")
	for i := 0; i < len(rows); i++ {
		row := strings.TrimSpace(rows[i])

		if strings.HasPrefix(row, "package") {
			file.Package = extractPackage(row)
			continue
		}

		if strings.HasPrefix(row, "import") {
			imports, idx := extractImport(rows, i)
			i = idx
			file.Imports = imports
			continue
		}

		if strings.HasPrefix(row, "/") {
			extractComment(rows, i)
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

func extractComment(rows []string, idx int) (unit, int) {
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
			if len(row) <= 1 || string(row[:2]) == "*/" {
				break
			}
			comments = append(comments, row)
			idx = i
		}
	}
	return unit{
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
