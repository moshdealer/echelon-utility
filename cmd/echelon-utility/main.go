package main

import (
	"fmt"
	"os"

	"echelon-utility/internal/analyzer"
	"echelon-utility/internal/cli"
	"echelon-utility/internal/config"
	"echelon-utility/internal/document"
	"echelon-utility/internal/output"
	"echelon-utility/internal/parser"
	"echelon-utility/internal/rule"
	"echelon-utility/internal/source"
)

type App struct {
	CLIOptions *cli.Options
	Cfg        *config.Config
	Input      *source.Input
	Resolver   *parser.Resolver
	Document   *document.Document
	Rules      []rule.Rule
	Analyzer   *analyzer.Analyzer
}

func main() {
	Utility := App{}

	// Анализ CLI аргументов и флагов
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(fmt.Errorf("cli error: %w", err))
		os.Exit(1)
	}
	Utility.CLIOptions = opts

	// Чтение конфига утилиты
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(fmt.Errorf("load config error: %w", err))
		os.Exit(1)
	}
	Utility.Cfg = cfg
	// Извлечение данных в виде байтов из файла/stdin
	var s source.Loader
	if Utility.CLIOptions.FromStdin {
		s = source.NewStream("stdin", os.Stdin)
	} else {
		s = source.NewFile(Utility.CLIOptions.InputPath)
	}

	Utility.Input, err = s.Load()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read source: %w", err))
		os.Exit(1)
	}

	// Подключаем Resolver для парсинга и валидации json/yaml
	Utility.Resolver = parser.NewResolver()
	Utility.Document, err = Utility.Resolver.Parse(Utility.Input.Name, Utility.Input.Content)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to parse document: %w", err))
		os.Exit(1)
	}

	// Создаем список правил и анализатор
	Utility.Rules = rule.Build(cfg.Rules)
	Utility.Analyzer = analyzer.New(Utility.Rules)

	result := Utility.Analyzer.Analyze(Utility.Document)

	// Вывод
	if err := output.WriteText(os.Stdout, result); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("write output error: %w", err))
		os.Exit(1)
	}

	// Завершение программы
	if len(result) > 0 && !Utility.CLIOptions.Silent {
		os.Exit(1)
	}
}
