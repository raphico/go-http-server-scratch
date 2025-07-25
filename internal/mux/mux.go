package mux

import (
	"strings"

	"github.com/raphico/go-http-server-scratch/internal/handler"
	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

type HandlerFunc func(w protocol.Response, r *protocol.Request)

type route struct {
	method  string
	pattern string
	handler HandlerFunc
}

type Mux struct {
	routes []route
}

func New() *Mux {
	return &Mux{
		routes: []route{},
	}
}

func (m *Mux) HandleFunc(pattern string, handler HandlerFunc) {
	method := "GET"

	parts := strings.SplitN(pattern, " ", 2)
	if len(parts) == 2 {
		method, pattern = parts[0], parts[1]
	}

	m.routes = append(m.routes, route{method, pattern, handler})
}

func (m *Mux) Match(w protocol.Response, r *protocol.Request) {
	for _, route := range m.routes {
		if !strings.EqualFold(route.method, r.Method) {
			continue
		}

		if matchPattern(r.URL.Path, route.pattern) {
			route.handler(w, r)
			return
		}
	}

	handler.NotFoundHandler(w, r)
}

func matchPattern(path, pattern string) bool {
	if pattern == path {
		return true
	}

	if pattern == "/" {
		return pattern == path
	}

	return strings.HasPrefix(path, pattern)
}
