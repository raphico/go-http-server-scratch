package handler

import (
	"fmt"
	"strings"

	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

func EchoHandler(w protocol.Response, r *protocol.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.Write(protocol.StatusBadRequest, nil)
		w.Send()
		return
	}

	body := parts[2]
	w.Write(protocol.StatusOk, []byte(body))
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	w.Send()
}
