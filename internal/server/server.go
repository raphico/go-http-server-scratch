package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/raphico/go-http-server-scratch/internal/mux"
	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

type Server struct {
	addr string
	mux  *mux.Mux
}

func New(addr string, mux *mux.Mux) *Server {
	return &Server{
		addr,
		mux,
	}
}

func (s *Server) Start() error {
	fmt.Printf("Starting server on %s\n", s.addr)

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to start bind to %s: %w", s.addr, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// close the tcp connection once done
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		response := protocol.NewResponse(conn)
		request, err := protocol.ParseRequest(reader)

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				break
			}

			fmt.Printf("%s", err.Error())
			response.Write(protocol.StatusBadRequest, nil)
			response.Send()
			continue
		}

		s.mux.Match(response, request)

		if request.Headers.Get("Connection") == "close" {
			break
		}
	}
}
