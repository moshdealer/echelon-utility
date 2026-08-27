package analyzer

import (
	"fmt"
	"sort"

	"echelon-utility/internal/document"
	"echelon-utility/internal/rule"
)

type Analyzer struct {
	rules []rule.Rule
}

func New(rules []rule.Rule) *Analyzer {
	return &Analyzer{
		rules: rules,
	}
}

func (a *Analyzer) Analyze(doc *document.Document) []rule.Finding {
	result := make([]rule.Finding, 0)

	if doc == nil || doc.Root == nil {
		return result
	}

	a.walk(doc.Root, "", "", &result)

	return result
}

// walk - рекурсивный анализ файла
func (a *Analyzer) walk(value any, path string, key string, result *[]rule.Finding) {
	switch current := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(current))

		for childKey := range current {
			keys = append(keys, childKey)
		}

		sort.Strings(keys)

		for _, childKey := range keys {
			childValue := current[childKey]
			childPath := joinPath(path, childKey)

			a.walk(childValue, childPath, childKey, result)
		}

	case []any:
		for index, childValue := range current {
			childPath := fmt.Sprintf("%s[%d]", path, index)

			a.walk(childValue, childPath, key, result)
		}
	default:
		for _, currentRule := range a.rules {
			found := currentRule.Check(path, key, current)
			if found == nil {
				continue
			}

			*result = append(*result, *found)
		}
	}
}

func joinPath(parent string, key string) string {
	if parent == "" {
		return key
	}

	return parent + "." + key
}
