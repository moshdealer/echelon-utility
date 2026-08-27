package main

import (
	"fmt"
	"os"

	"echelon-utility/internal/cli"
	"echelon-utility/internal/config"
	"echelon-utility/internal/output"
	"echelon-utility/internal/pathscanner"
	"echelon-utility/internal/rule"
	"echelon-utility/internal/scanner"
	"echelon-utility/internal/source"
)

type App struct {
	CLIOptions *cli.Options
	Cfg        *config.Config
	Scanner    scanner.Scanner
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
	Utility.Scanner = scanner.New(rule.Build(cfg.Rules))
	pathScanService := pathscanner.New(Utility.Scanner, rule.BuildFileRules(cfg.Rules))

	var result []rule.Finding
	var scanErrors []error
	if Utility.CLIOptions.FromStdin {
		input, err := source.NewStream("stdin", os.Stdin).Load()
		if err != nil {
			fmt.Println(fmt.Errorf("failed to read source: %w", err))
			os.Exit(1)
		}

		result, err = Utility.Scanner.Scan(input.Name, input.Content)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to parse document: %w", err))
			os.Exit(1)
		}
	} else {
		result, scanErrors = pathScanService.Scan(Utility.CLIOptions.InputPath)
	}

	// Вывод
	if err := output.WriteText(os.Stdout, result); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("write output error: %w", err))
		os.Exit(1)
	}

	for _, scanErr := range scanErrors {
		fmt.Println(fmt.Errorf("scan error: %w", scanErr))
	}
	if len(scanErrors) > 0 {
		os.Exit(1)
	}

	// Завершение программы
	if len(result) > 0 && !Utility.CLIOptions.Silent {
		os.Exit(1)
	}
}
