package parser

import "echelon-utility/internal/document"

// Parser преобразует сырое содержимое конфигурации в единое внутреннее представление
// Интерфейс необходим для универсальной работы с разными форматами json/yaml
type Parser interface {
	Parse(content []byte) (*document.Document, error)
}
