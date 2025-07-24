package mux

import (
	"strings"

	"github.com/raphico/go-http-server-scratch/internal/handler"
	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

type HandlerFunc func(w protocol.Response, r *protocol.Request)

type Routes map[string]HandlerFunc

type Mux struct {
	Routes Routes
}

func New() *Mux {
	return &Mux {
		Routes: make(Routes),
	}
}

func (m *Mux) HandleFunc(pattern string, handler HandlerFunc) {
	m.Routes[pattern] = handler
}

func (m *Mux) Match(w protocol.Response, r *protocol.Request) {
	for pattern, handler := range m.Routes {
		// an exact match
		if pattern == r.URL.Path {
			handler(w, r)
			return
		}

		// for dynamic routes
		if 
			pattern != "/" &&
			strings.HasPrefix(pattern, "/") && 
			strings.HasPrefix(r.URL.Path, pattern) {
				handler(w, r)
				return
			}
	}

	handler.NotFoundHandler(w, r)
}