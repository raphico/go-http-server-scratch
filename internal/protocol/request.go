package protocol

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"strings"
)

type Request struct {
	Headers Header
	Method string 
	URL *url.URL
}

func ParseRequest(conn net.Conn) (*Request, error) {
	var request = Request {
		Headers: make(Header),
	}

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %v", err)
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("failed to read request: invalid request line")
	}

	parsedUrl, err := url.Parse(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %v", err)
	}

	request.Method = parts[0]
	request.URL = parsedUrl

	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read request: %v", err)
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
		request.Headers.Set(key, value)
	}

	return &request, nil
}