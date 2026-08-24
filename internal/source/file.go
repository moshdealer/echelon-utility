package source

import (
	"fmt"
	"os"
)

type File struct {
	path string
}

func NewFile(path string) *File {
	return &File{
		path: path,
	}
}

func (f *File) Load() (*Input, error) {
	content, err := os.ReadFile(f.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", f.path, err)
	}

	return &Input{
		Name:    f.path,
		Content: content,
	}, nil
}
