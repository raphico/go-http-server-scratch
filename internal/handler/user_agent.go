package handler

import (
	"fmt"

	"github.com/raphico/go-http-server-scratch/internal/compress"
	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

func UserAgentHandler(w protocol.Response, r *protocol.Request) {
	body := r.Headers.Get("User-Agent")
	compressedBody := compress.CompressIfSupported(w, r, []byte(body))

	w.Write(protocol.StatusOk, compressedBody)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(compressedBody)))
	w.Send()
}
