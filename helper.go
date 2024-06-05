package goast

import (
	"go/format"
	"io"
	"os"
	"strings"
)

func readFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	formatted, err := format.Source(buf)
	if err != nil {
		return nil, err
	}

	return formatted, nil
}

func hasPrefix(s []byte, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}

	for i := range prefix {
		if s[i] != prefix[i] {
			return false
		}
	}

	return true
}

func hasSuffix(s []byte, suffix string) bool {
	if len(suffix) > len(s) {
		return false
	}
	ss := s[len(s)-len(suffix):]
	for i := range ss {
		if ss[i] != suffix[i] {
			return false
		}
	}

	return true
}

func printTidy(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, " ", "\\s")
	return s
}
