package server

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"echelon-utility/internal/rule"
	"echelon-utility/internal/scanner"
)

type Server struct {
	scanner scanner.Scanner
}

type scanResponse struct {
	Findings []rule.Finding `json:"findings"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func New(s scanner.Scanner) http.Handler {
	server := &Server{scanner: s}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /scan", server.scan)

	return mux
}

func (s *Server) scan(w http.ResponseWriter, r *http.Request) {
	sourceName, err := sourceName(r.Header.Get("Content-Type"))
	if err != nil {
		writeJSON(w, http.StatusUnsupportedMediaType, errorResponse{Error: err.Error()})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "failed to read request body"})
		return
	}

	findings, err := s.scanner.Scan(sourceName, body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, scanResponse{Findings: findings})
}

func sourceName(contentType string) (string, error) {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", errors.New("invalid Content-Type")
	}

	switch mediaType {
	case "application/json":
		return "request.json", nil
	case "application/yaml", "application/x-yaml", "text/yaml":
		return "request.yaml", nil
	default:
		return "", errors.New("Content-Type must be application/json or application/yaml")
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}
