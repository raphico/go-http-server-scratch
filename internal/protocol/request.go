package protocol

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	Headers Header
	Method  string
	URL     *url.URL
	Body    []byte
}

func ParseRequest(conn net.Conn) (*Request, error) {
	var request = Request{
		Headers: make(Header),
	}

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("failed to read request: invalid request line")
	}

	parsedUrl, err := url.Parse(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %w", err)
	}

	request.Method = parts[0]
	request.URL = parsedUrl

	// adds request's headers
	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read request: %w", err)
		}

		header = strings.TrimSpace(header)
		if header == "" {
			break
		}

		parts := strings.Split(header, ": ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("failed to read request: invalid header")
		}

		key, value := parts[0], parts[1]
		for v := range strings.SplitSeq(value, ", ") {
			request.Headers.Set(key, v)
		}
	}

	// adds request's body
	clStr := request.Headers.Get("Content-Length")
	if len(clStr) > 0 {
		contentLength, err := strconv.Atoi(clStr)
		if err != nil {
			return nil, fmt.Errorf("invalid Content-Length: %w", err)
		}

		buf := make([]byte, contentLength)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, fmt.Errorf("failed to read request's body: %w", err)
		}

		request.Body = buf
	} else {
		request.Body = nil
	}

	return &request, nil
}
