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
	return &Mux{
		Routes: make(Routes),
	}
}

func (m *Mux) HandleFunc(pattern string, handler HandlerFunc) {
	m.Routes[pattern] = handler
}

func (m *Mux) Match(w protocol.Response, r *protocol.Request) {
	for pattern, handler := range m.Routes {
		p := pattern
		method := "GET"

		parts := strings.Split(pattern, " ")
		if len(parts) > 2 {
			w.Write(protocol.StatusInternalServerError, nil)
			w.Send()
			return
		}

		if len(parts) == 2 {
			p = parts[1]
			method = parts[0]
		}

		// an exact match
		if p == r.URL.Path && method == r.Method {
			handler(w, r)
			return
		}

		// for dynamic routes
		if p != "/" &&
			strings.HasPrefix(p, "/") &&
			strings.HasPrefix(r.URL.Path, p) &&
			method == r.Method {
			handler(w, r)
			return
		}
	}

	handler.NotFoundHandler(w, r)
}
