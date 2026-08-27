package parser

import (
	"encoding/json"
	"fmt"

	"echelon-utility/internal/document"
)

type JSON struct{}

func (p *JSON) Parse(content []byte) (*document.Document, error) {
	var root any

	if err := json.Unmarshal(content, &root); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return &document.Document{Root: root}, nil
}
