package handler

import (
	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

func NotFoundHandler(w protocol.Response, r *protocol.Request) {
	w.Write(protocol.StatusNotFound, nil)
	w.Send()
}