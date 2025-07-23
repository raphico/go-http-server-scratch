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
	var response Response = NewResponse(conn)

	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read request: ", err.Error())
		response.Write(StatusInternalServerError, nil)
		response.Send()
		return;
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) < 3 {
		response.Write(StatusBadRequest, nil)
		response.Send()
		return;
	}

	path := parts[1]
	switch {
	case path == "/":
		response.Write(StatusOk, nil)
		response.Send()
	case strings.HasPrefix(path, "/echo/"):
		parts = strings.Split(path, "/")
		if len(parts) != 3 {
			response.Write(StatusBadRequest, nil)
			response.Send()
			return;
		}

		msg := parts[2];
		response.Write(StatusOk, []byte(msg))
		response.headers.Set("Content-Type", "text/plain")
		response.headers.Set("Content-Length", "3")
		response.Send()

	default:
		response.Write(StatusNotFound, nil)
		response.Send()
	}
}

const (
	StatusOk = 200
	StatusBadRequest = 400
	StatusNotFound = 404
	StatusInternalServerError = 500
)

var StatusText = map[int]string {
	StatusOk: "OK",
	StatusBadRequest: "Bad Request",
	StatusNotFound: "Not Found",
	StatusInternalServerError: "Internal Server Error",
}

type Header map[string][]string

type Response struct {
	conn net.Conn
	statusCode int
	body []byte
	headers  Header
}

func NewResponse(conn net.Conn) Response {
	return Response {
		conn: conn,
		headers: make(Header),
	}
}

func (r *Response) Header() Header {
	return r.headers
}

func (h *Header) Set(key string, value string) {
	(*h)[key] = append((*h)[key], value)
}

func (r *Response) Write(statusCode int, body []byte) {
	r.body = body
	r.statusCode = statusCode
}

func (r *Response) Send () {
	var builder strings.Builder	
	
	// write status line
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.statusCode, StatusText[r.statusCode])
	builder.WriteString(statusLine)

	// write http headers
	for key, values := range r.headers {
		for _, v := range values {
			fmt.Fprintf(&builder, "%s: %s\r\n", key, v)
		}
	}

	// end headers
	builder.WriteString("\r\n")

	// write body
	builder.WriteString(string(r.body))

	r.conn.Write([]byte(builder.String()))
}