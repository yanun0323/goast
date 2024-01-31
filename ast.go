package goast

import (
	"io"
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

	rows := strings.Split(string(buf), "\n")
	for i := 0; i < len(rows); i++ {
		// TODO: parse to nodes
	}

	return &File{}, nil
}
