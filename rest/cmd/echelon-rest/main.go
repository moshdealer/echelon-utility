package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"echelon-utility/internal/config"
	"echelon-utility/internal/rule"
	"echelon-utility/internal/scanner"
	"echelon-utility/rest/server"
)

func main() {
	address := flag.String("address", ":8080", "HTTP server listen address")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("load config error: %w", err))
	}

	scanService := scanner.New(rule.Build(cfg.Rules))
	handler := server.New(scanService)

	log.Printf("REST server is listening on %s", *address)
	if err := http.ListenAndServe(*address, handler); err != nil {
		log.Fatal(fmt.Errorf("REST server error: %w", err))
	}
}
