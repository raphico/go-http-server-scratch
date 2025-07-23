package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

const port = 4221

func main() {
	fmt.Printf("Starting server on port %d\n", port);

	address := fmt.Sprintf("0.0.0.0:%d", port)

	listener, err := net.Listen("tcp", address)
	if err != nil {

		fmt.Fprintln(os.Stderr, "Error accepting connection: ", err.Error())	
		os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to accept connection: ", err.Error())
		return;
	}

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
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
