package output

import (
	"fmt"
	"io"
	"strings"

	"echelon-utility/internal/rule"
)

// WriteText выводит найденные проблемы в читаемом формате
func WriteText(w io.Writer, findings []rule.Finding) error {
	for index, finding := range findings {
		message := joinSentences(finding.Message, finding.Recommendation)

		if _, err := fmt.Fprintf(w, "%s: [%s] %s\n", finding.Severity, finding.Path, message); err != nil {
			return fmt.Errorf("write finding %d: %w", index, err)
		}
	}

	return nil
}

func joinSentences(parts ...string) string {
	sentences := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if !strings.ContainsAny(part[len(part)-1:], ".!?") {
			part += "."
		}

		sentences = append(sentences, part)
	}

	return strings.Join(sentences, " ")
}
