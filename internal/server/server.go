package server

import (
	"fmt"
	"net"
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
	response := protocol.NewResponse(conn)

	request, err := protocol.ParseRequest(conn)
	if err != nil {
		fmt.Printf("%s", err.Error())
		response.Write(protocol.StatusBadRequest, nil)
		response.Send()
		return
	}

	r := *request

	switch {
	case r.URL.Path == "/":
		response.Write(protocol.StatusOk, nil)
		response.Send()

	case strings.HasPrefix(r.URL.Path, "/echo/"):
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			response.Write(protocol.StatusBadRequest, nil)
			response.Send()
			return;
		}

		body := parts[2];
		response.Write(protocol.StatusOk, []byte(body))
		response.Header().Set("Content-Type", "text/plain")
		response.Header().Set("Content-Length", "3")
		response.Send()

	case r.URL.Path == "/user-agent":
		body := r.Headers.Get("User-Agent")
		response.Write(protocol.StatusOk, []byte(body))
		response.Header().Set("Content-Type", "text/plain")
		response.Header().Set("Content-Length", fmt.Sprint(len(body)))
		response.Send()

	default:
		response.Write(protocol.StatusNotFound, nil)
		response.Send()
	}
}
