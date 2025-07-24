package handler

import (
	"fmt"

	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

func UserAgentHandler(w protocol.Response, r *protocol.Request) {
	body := r.Headers.Get("User-Agent")
	w.Write(protocol.StatusOk, []byte(body))
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	w.Send()
}
