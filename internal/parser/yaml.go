package parser

import (
	"fmt"

	"echelon-utility/internal/document"
	yaml "go.yaml.in/yaml/v3"
)

type YAML struct{}

func (y *YAML) Parse(content []byte) (*document.Document, error) {
	var root any

	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	return &document.Document{Root: root}, nil
}
