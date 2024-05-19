package goast

import (
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type rawFile struct {
	dir   string
	name  string
	lines []string
}

func newRawFile(filePath string) (*rawFile, error) {
	f, err := os.Open(filePath)
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

	return &rawFile{
		dir:   filepath.Dir(filePath),
		name:  filepath.Base(filePath),
		lines: strings.Split(string(formatted), "\n"),
	}, nil
}
