package compress

import "github.com/raphico/go-http-server-scratch/internal/protocol"

func CompressIfSupported(w protocol.Response, r *protocol.Request, body []byte) []byte {
	encoding := r.Headers.Get("Accept-Encoding")

	switch encoding {
	case "gzip":
		w.Header().Set("Content-Encoding", "gzip")
		return body
	default:
		return body
	}
}
