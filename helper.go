package goast

import (
	"go/format"
	"io"
	"os"
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
