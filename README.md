# Build Your Own HTTP/1.1 Server in GO

This is a fully functional HTTP/1.1 server implemented from scratch using only low-level networking primitives provided by the Go standard library. It was created as a systems-level learning project to explore how web servers operate at the lowest level, including TCP socket management, manual HTTP parsing, and concurrent connection handling

## Motivation

This project was built to:

1. Deepen understanding of how HTTP and TCP networking work under the hood.
2. Strengthen systems programming skills, especially around concurrency, sockets, and I/O
3. Practice Go by focusing on idiomatic structuring, testing, and low-level networking construct

## Features

- Manual TCP server using `net.Listen` and `Accept`
- Manual HTTP/1.1 request parsing (method, path, headers)
- Supports basic HTTP routes:
  - `GET /` — Health check
  - `GET /echo/<msg>` — Echoes message in response body
  - `GET /user-agent` — Returns `User-Agent` header
  - `GET /files/<filename>` — Return a file content
  - `POST /files/<filename>` — Reads request body and Saves file to disk
- Gzip compression for clients that support it (`Accept-Encoding: gzip`)
- Concurrency using goroutines (one per connection)
- Persistent connections (`Connection: keep-alive`)
- Proper request body reading and `Content-Length` handling
- 404 Not Found for unknown routes
- Graceful connection closing (`Connection: close`, EOF handling)
- Prevents directory traversal (e.g., `../`)
- No third-party dependencies — only Go standard library

## How to run

1. clone the repository

```bash
git clone https://github.com/raphico/go-http-server-scratch.git
cd go-http-server-scratch
```

2. Run the server

```bash
go run cmd/server/main.go
```

3. Try it out

```bash
curl http://localhost:4221/
curl http://localhost:4221/echo/hello
curl -v --header "User-Agent: foobar/1.2.3" http://localhost:4221/user-agent
```

## Folder structure

| File                            | Purpose                                          |
| ------------------------------- | ------------------------------------------------ |
| `cmd/server/main.go`            | Start server                                     |
| `internal/server/server.go`     | Handles TCP socket and connection management     |
| `internal/protocol/request.go`  | HTTP request parsing                             |
| `internal/protocol/response.go` | HTTP response generation                         |
| `internal/protocol/header.go`   | Manages HTTP header parsing, storage, and access |
| `internal/server/handler`       | handler functions                                |
| `internal/server/mux`           | Defines multiplexer                              |

## Credits

Inspired by [@coder-crafters](https://codecrafters.io/) HTTP Server Challenge
All implementation written from scratch by me.
