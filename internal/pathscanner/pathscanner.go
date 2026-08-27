package pathscanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"echelon-utility/internal/rule"
	"echelon-utility/internal/scanner"
	"echelon-utility/internal/source"
)

// Service проверяет один конфигурационный файл или все конфиги в директории
// (бонусное задание на рекурсивный обход директории)
type Service struct {
	contentScanner scanner.Scanner
	fileRules      []rule.FileRule
}

func New(contentScanner scanner.Scanner, fileRules []rule.FileRule) *Service {
	return &Service{
		contentScanner: contentScanner,
		fileRules:      fileRules,
	}
}

func (s *Service) Scan(inputPath string) ([]rule.Finding, []error) {
	info, err := os.Stat(inputPath)
	if err != nil {
		return nil, []error{fmt.Errorf("stat source %q: %w", inputPath, err)}
	}

	if !info.IsDir() {
		findings, err := s.scanFile(inputPath, inputPath, false)
		if err != nil {
			return findings, []error{err}
		}
		return findings, nil
	}

	files, scanErrors := configFiles(inputPath)
	if len(files) == 0 {
		scanErrors = append(scanErrors, fmt.Errorf("no JSON or YAML configuration files found in %q", inputPath))
		return nil, scanErrors
	}

	findings := make([]rule.Finding, 0)
	for _, filePath := range files {
		fileFindings, err := s.scanFile(filePath, filepath.ToSlash(filePath), true)
		findings = append(findings, fileFindings...)
		if err != nil {
			scanErrors = append(scanErrors, err)
		}
	}

	return findings, scanErrors
}

func (s *Service) scanFile(filePath string, displayName string, fromDirectory bool) ([]rule.Finding, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("stat file %q: %w", filePath, err)
	}

	fileFindings := make([]rule.Finding, 0)
	for _, currentRule := range s.fileRules {
		findingPath := displayName
		if fromDirectory {
			findingPath = ""
		}

		finding := currentRule.Check(findingPath, info.Mode())
		if finding == nil {
			continue
		}
		if fromDirectory {
			finding.Source = displayName
		}
		fileFindings = append(fileFindings, *finding)
	}

	input, err := source.NewFile(filePath).Load()
	if err != nil {
		return fileFindings, err
	}

	findings, err := s.contentScanner.Scan(input.Name, input.Content)
	if err != nil {
		return fileFindings, fmt.Errorf("scan %q: %w", filePath, err)
	}

	if fromDirectory {
		for index := range findings {
			findings[index].Source = displayName
		}
	}

	return append(findings, fileFindings...), nil
}

func configFiles(root string) ([]string, []error) {
	files := make([]string, 0)
	scanErrors := make([]error, 0)

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			scanErrors = append(scanErrors, fmt.Errorf("walk %q: %w", path, walkErr))
			if entry != nil && entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !isConfigFile(path) {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		scanErrors = append(scanErrors, fmt.Errorf("walk directory %q: %w", root, err))
	}

	sort.Strings(files)
	return files, scanErrors
}

func isConfigFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json", ".yaml", ".yml":
		return true
	default:
		return false
	}
}
