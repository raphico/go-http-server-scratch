package main

import (
	"fmt"
	"os"

	"github.com/raphico/go-http-server-scratch/internal/server"
)

const port = 4221

func main() {
	addr := fmt.Sprintf("0.0.0.0:%v", port)
	
	s := server.New(addr)
	if err := s.Start(); err != nil {
		fmt.Fprint(os.Stderr, "Failed to start server: ", err.Error())
		os.Exit(1)
	}
}
