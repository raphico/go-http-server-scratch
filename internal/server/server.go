package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

type Server struct {
	addr string
}

func New(addr string) *Server {
	return &Server {
		addr,
	}
}

func (s *Server) Start() error {
	fmt.Printf("Starting server on %s\n", s.addr)

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		port := strings.Split(s.addr, ":")[1]
		return fmt.Errorf("failed to start bind to port: %v", port);
	}

	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept connection: %v", err)
	}

	s.handleConnection(conn)

	return nil
}

func (s *Server) handleConnection(conn net.Conn)  {
	var response protocol.Response = protocol.NewResponse(conn)

	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read request: ", err.Error())
		response.Write(protocol.StatusInternalServerError, nil)
		response.Send()
		return;
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) < 3 {
		response.Write(protocol.StatusBadRequest, nil)
		response.Send()
		return;
	}

	path := parts[1]
	switch {
	case path == "/":
		response.Write(protocol.StatusOk, nil)
		response.Send()
	case strings.HasPrefix(path, "/echo/"):
		parts = strings.Split(path, "/")
		if len(parts) != 3 {
			response.Write(protocol.StatusBadRequest, nil)
			response.Send()
			return;
		}

		msg := parts[2];
		response.Write(protocol.StatusOk, []byte(msg))
		response.Header().Set("Content-Type", "text/plain")
		response.Header().Set("Content-Length", "3")
		response.Send()

	default:
		response.Write(protocol.StatusNotFound, nil)
		response.Send()
	}
}
