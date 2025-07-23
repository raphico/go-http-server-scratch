package protocol

import (
	"fmt"
	"net"
	"strings"
)

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

func (h Header) Set(key string, value string) {
	h[key] = append(h[key], value)
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