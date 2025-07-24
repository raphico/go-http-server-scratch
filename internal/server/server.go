package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/raphico/go-http-server-scratch/internal/mux"
	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

type Server struct {
	addr string
	mux *mux.Mux
}

func New(addr string, mux *mux.Mux) *Server {
	return &Server {
		addr,
		mux,
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

	s.mux.Match(response, request)
}
