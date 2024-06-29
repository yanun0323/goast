package helper

import (
	"fmt"
	"go/format"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

var (
	_debug bool
)

func SetDebug(debug bool) {
	_debug = debug
}

func DebugPrint(args ...any) {
	if _debug {
		ss := strings.Split(string(debug.Stack()), "\n")
		s := ""
		if len(ss) >= 7 {
			s = ss[6]
		}
		args = append(args, "\t", strings.Split(s, " ")[0])
		fmt.Println(args...)
	}
}

func ReadFile(file string) ([]byte, error) {
	if !HasSuffix([]byte(file), ".go") {
		file = file + ".go"
	}

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

func SaveFile(file string, data []byte) error {
	if !HasSuffix([]byte(file), ".go") {
		file = file + ".go"
	}

	formatted, err := format.Source(data)
	if err != nil {
		slog.Error(fmt.Sprintf("format ast data, err: %+v", err))
		formatted = data
	}

	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("mkdir %s, err: %w", dir, err)
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("create file %s, err: %w", file, err)
	}
	defer f.Close()

	if _, err := f.Write(formatted); err != nil {
		return fmt.Errorf("write file %s, err: %w", file, err)
	}

	return nil
}

func HasPrefix(s []byte, prefix string) bool {
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

func HasSuffix(s []byte, suffix string) bool {
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

func TidyText(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "·")
	s = strings.ReplaceAll(s, " ", "·")
	s = strings.ReplaceAll(s, "\t", " -> ")
	return s
}

func AppendUnrepeatable[Type comparable](slice []Type, elems ...Type) []Type {
	if len(elems) == 0 {
		return slice
	}

	if len(slice) == 0 {
		return elems
	}

	if slice[len(slice)-1] == elems[0] {
		elems = elems[1:]
	}

	return append(slice, elems...)
}
