package compress

import (
	"bytes"
	"compress/gzip"
	"slices"

	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

func CompressIfSupported(w protocol.Response, r *protocol.Request, body []byte) []byte {
	accept := r.Headers.Values("Accept-Encoding")

	if slices.Contains(accept, "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		data, err := gzipBytes(body)
		if err != nil {
			return body
		}

		return data
	}

	return body
}

func gzipBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
