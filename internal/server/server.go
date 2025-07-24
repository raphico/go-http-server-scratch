package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/raphico/go-http-server-scratch/internal/mux"
	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

type Server struct {
	addr   string
	mux    *mux.Mux
	logger *slog.Logger
}

func New(addr string, mux *mux.Mux, logger *slog.Logger) *Server {
	return &Server{
		addr,
		mux,
		logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting server", "address", s.addr)

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
	defer func() {
		s.logger.Info("closing connection", "remote_addr", conn.RemoteAddr().String())
		conn.Close()
	}()

	s.logger.Info("new connection accepted", "remote_addr", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)

	for {
		response := protocol.NewResponse(conn)
		request, err := protocol.ParseRequest(reader)

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				s.logger.Info("connection closed by client", "remote_addr", conn.RemoteAddr().String())
				break
			}

			s.logger.Error("failed to parse request", "remote_addr", conn.RemoteAddr().String(), "error", err)
			response.Write(protocol.StatusBadRequest, nil)
			response.Send()
			continue
		}

		s.logger.Info("request received",
			"remote_addr", conn.RemoteAddr().String(),
			"method", request.Method,
			"path", request.URL.Path,
		)

		s.mux.Match(response, request)

		if request.Headers.Get("Connection") == "close" {
			s.logger.Info("client requested connection close", "remote_addr", conn.RemoteAddr().String())
			break
		}
	}
}
