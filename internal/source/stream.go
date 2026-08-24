package source

import (
	"fmt"
	"io"
)

type Stream struct {
	name   string
	reader io.Reader
}

func NewStream(name string, r io.Reader) *Stream {
	return &Stream{
		name:   name,
		reader: r,
	}
}

func (s *Stream) Load() (*Input, error) {
	content, err := io.ReadAll(s.reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read stream %q: %w", s.name, err)
	}

	return &Input{
		Name:    s.name,
		Content: content,
	}, nil
}
