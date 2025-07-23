# Build Your Own HTTP/1.1 Server in GO

A low-level HTTP/1.1 server implementation from scratch using only TCP primitives provided by the GO standard library. Inspired by the CoderCrafters challenge.

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
  - `GET /files/<filename>` — Serves static file from `./public/`
  - `POST /files/<filename>` — Saves file to disk in `./public/`
- Gzip compression for clients that support it (`Accept-Encoding: gzip`)
- Concurrency using goroutines (one per connection)
- Persistent connections (`Connection: keep-alive`)
- Proper request body reading and `Content-Length` handling
- 404 Not Found for unknown routes
- Graceful connection closing (`Connection: close`, EOF handling)
- Prevents directory traversal (e.g., `../`)
- No third-party dependencies — only Go standard library

## Credits

Inspired by [@coder-crafters](https://codecrafters.io/) HTTP Server Challenge
All implementation written from scratch by me.
