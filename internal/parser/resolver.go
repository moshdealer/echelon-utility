package parser

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"echelon-utility/internal/document"
)

// Resolver решает какой парсер использовать

type Resolver struct {
	json Parser
	yaml Parser
}

func NewResolver() *Resolver {
	return &Resolver{
		json: &JSON{},
		yaml: &YAML{},
	}
}

func (r *Resolver) Parse(sourceName string, content []byte) (*document.Document, error) {
	ext := strings.ToLower(filepath.Ext(sourceName))
	switch ext {
	case ".yaml", ".yml":
		doc, err := r.yaml.Parse(content)
		if err != nil {
			return nil, err
		}
		if err := validate(doc); err != nil {
			return nil, err
		}
		return doc, nil
	case ".json":
		doc, err := r.json.Parse(content)
		if err != nil {
			return nil, err
		}
		if err := validate(doc); err != nil {
			return nil, err
		}
		return doc, nil
	case "":
		jsonDocument, jsonErr := r.json.Parse(content)
		if jsonErr == nil {
			if err := validate(jsonDocument); err != nil {
				return nil, err
			}
			return jsonDocument, nil
		}

		yamlDocument, yamlErr := r.yaml.Parse(content)
		if yamlErr == nil {
			if err := validate(yamlDocument); err != nil {
				return nil, err
			}
			return yamlDocument, nil
		}

		return nil, fmt.Errorf(
			"failed to parse source %q as JSON or YAML: %w", sourceName, errors.Join(jsonErr, yamlErr),
		)
	default:
		return nil, fmt.Errorf(
			"unsupported configuration format %q", ext,
		)
	}
}

func validate(doc *document.Document) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}

	if doc.Root == nil {
		return fmt.Errorf("configuration document is empty")
	}

	if _, ok := doc.Root.(map[string]any); !ok {
		return fmt.Errorf(
			"configuration root must be an object, got %T",
			doc.Root,
		)
	}

	return nil
}
