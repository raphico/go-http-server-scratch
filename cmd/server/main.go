package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/raphico/go-http-server-scratch/internal/handler"
	"github.com/raphico/go-http-server-scratch/internal/mux"
	"github.com/raphico/go-http-server-scratch/internal/server"
)

const port = 4221

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	addr := fmt.Sprintf("0.0.0.0:%v", port)

	mux := mux.New()

	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/echo/", handler.EchoHandler)
	mux.HandleFunc("/user-agent", handler.UserAgentHandler)
	mux.HandleFunc("GET /files", handler.GetFileHandler)
	mux.HandleFunc("POST /files", handler.PostFileHandler)

	s := server.New(addr, mux, logger)
	if err := s.Start(); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
