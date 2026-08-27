package scanner

import (
	"echelon-utility/internal/analyzer"
	"echelon-utility/internal/parser"
	"echelon-utility/internal/rule"
)

// Scanner описывает общий для всех адаптеров способ проверки конфигурации
type Scanner interface {
	Scan(sourceName string, content []byte) ([]rule.Finding, error)
}

type Service struct {
	resolver *parser.Resolver
	analyzer *analyzer.Analyzer
}

func New(rules []rule.Rule) *Service {
	return &Service{
		resolver: parser.NewResolver(),
		analyzer: analyzer.New(rules),
	}
}

func (s *Service) Scan(sourceName string, content []byte) ([]rule.Finding, error) {
	doc, err := s.resolver.Parse(sourceName, content)
	if err != nil {
		return nil, err
	}

	return s.analyzer.Analyze(doc), nil
}
