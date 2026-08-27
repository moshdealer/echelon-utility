package server

import (
	"context"

	echelonv1 "echelon-utility/grpc/api/echelon/v1"
	"echelon-utility/internal/rule"
	"echelon-utility/internal/scanner"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	echelonv1.UnimplementedScannerServiceServer
	scanner scanner.Scanner
}

func New(s scanner.Scanner) *Server {
	return &Server{scanner: s}
}

func (s *Server) Scan(_ context.Context, request *echelonv1.ScanRequest) (*echelonv1.ScanResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	findings, err := s.scanner.Scan(request.GetSourceName(), request.GetContent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "scan config: %v", err)
	}

	response := &echelonv1.ScanResponse{
		Findings: make([]*echelonv1.Finding, 0, len(findings)),
	}
	for _, finding := range findings {
		response.Findings = append(response.Findings, toProtoFinding(finding))
	}

	return response, nil
}

func toProtoFinding(finding rule.Finding) *echelonv1.Finding {
	return &echelonv1.Finding{
		Severity:       finding.Severity,
		Rule:           finding.Rule,
		Path:           finding.Path,
		Message:        finding.Message,
		Recommendation: finding.Recommendation,
	}
}
