package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read request: ", err.Error())
		return;
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) < 3 {
		writeResponse(conn, "400", "Bad Request")
		return;
	}

	path := parts[1]
	switch path {
	case "/":
		writeResponse(conn, "200", "OK")
	default:
		writeResponse(conn, "404", "Not Found")
	}
}


func writeResponse(conn net.Conn, statusCode string, reason string) {
	response := fmt.Sprintf(`HTTP/1.1 %s %s\r\n\r\n`, statusCode, reason)
	conn.Write([]byte(response))
}