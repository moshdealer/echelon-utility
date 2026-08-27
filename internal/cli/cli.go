package cli

import (
	"flag"
	"fmt"
)

// Пакет cli нужен для обработки флагов и аргументов, переданных при запуске утилиты из командной строки

type Options struct {
	InputPath string
	Silent    bool
	FromStdin bool
}

func Parse(args []string) (*Options, error) {
	var opts Options

	fs := flag.NewFlagSet(
		"echelon-utility",
		flag.ContinueOnError,
	)

	fs.BoolVar(
		&opts.Silent,
		"s",
		false,
		"do not return an error when problems are found",
	)

	fs.BoolVar(
		&opts.Silent,
		"silent",
		false,
		"do not return an error when problems are found",
	)

	fs.BoolVar(
		&opts.FromStdin,
		"stdin",
		false,
		"read configuration from standard input",
	)

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("parse flags: %w", err)
	}

	positionalArgs := fs.Args()

	if opts.FromStdin {
		if len(positionalArgs) != 0 {
			return nil, fmt.Errorf("positional argument and input from a stdin cannot be used together")
		}
	} else {
		if len(positionalArgs) != 1 {
			return nil, fmt.Errorf("expected one positional argument")
		}
		opts.InputPath = positionalArgs[0]
	}

	return &opts, nil
}
