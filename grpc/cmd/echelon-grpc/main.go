package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	echelonv1 "echelon-utility/grpc/api/echelon/v1"
	grpcserver "echelon-utility/grpc/server"
	"echelon-utility/internal/config"
	"echelon-utility/internal/rule"
	"echelon-utility/internal/scanner"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	address := flag.String("address", ":9090", "gRPC server listen address")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("load config error: %w", err))
	}

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatal(fmt.Errorf("listen gRPC address %s: %w", *address, err))
	}

	scanService := scanner.New(rule.Build(cfg.Rules))
	server := grpc.NewServer()
	echelonv1.RegisterScannerServiceServer(server, grpcserver.New(scanService))
	reflection.Register(server)

	log.Printf("gRPC server is listening on %s", *address)
	if err := server.Serve(listener); err != nil {
		log.Fatal(fmt.Errorf("gRPC server error: %w", err))
	}
}
