package main

import (
	"fmt"
	"os"

	"github.com/raphico/go-http-server-scratch/internal/handler"
	"github.com/raphico/go-http-server-scratch/internal/mux"
	"github.com/raphico/go-http-server-scratch/internal/server"
)

const port = 4221

func main() {
	addr := fmt.Sprintf("0.0.0.0:%v", port)

	mux := mux.New()

	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/echo/", handler.EchoHandler)
	mux.HandleFunc("/user-agent", handler.UserAgentHandler)
	
	s := server.New(addr, mux)
	if err := s.Start(); err != nil {
		fmt.Fprint(os.Stderr, "Failed to start server: ", err.Error())
		os.Exit(1)
	}
}
