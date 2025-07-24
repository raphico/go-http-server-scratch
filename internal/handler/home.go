package handler

import (
	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

func HomeHandler(w protocol.Response, r *protocol.Request) {
	w.Write(protocol.StatusOk, nil)
	w.Send()
}