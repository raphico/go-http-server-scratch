package handler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/raphico/go-http-server-scratch/internal/compress"
	"github.com/raphico/go-http-server-scratch/internal/protocol"
)

func PostFileHandler(w protocol.Response, r *protocol.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.Write(protocol.StatusBadRequest, nil)
		w.Send()
		return
	}

	filename := parts[2]
	safePath, ok, err := isPathSafe("/tmp", filename)
	if err != nil || !ok {
		w.Write(protocol.StatusBadRequest, nil)
		w.Send()
		return
	}

	file, err := os.Create(safePath)
	if err != nil {
		w.Write(protocol.StatusInternalServerError, nil)
		w.Send()
		return
	}

	defer file.Close()

	_, err = file.Write(r.Body)
	if err != nil {
		w.Write(protocol.StatusInternalServerError, nil)
		w.Send()
		return
	}

	w.Write(protocol.StatusCreated, nil)
	w.Send()
}

func GetFileHandler(w protocol.Response, r *protocol.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.Write(protocol.StatusBadRequest, nil)
		w.Send()
		return
	}

	filename := parts[2]
	safePath, ok, err := isPathSafe("/tmp", filename)
	if err != nil || !ok {
		w.Write(protocol.StatusBadRequest, nil)
		w.Send()
		return
	}

	data, err := os.ReadFile(safePath)
	if err != nil {
		w.Write(protocol.StatusNotFound, nil)
		w.Send()
		return
	}

	compressedBody := compress.CompressIfSupported(w, r, data)
	w.Write(protocol.StatusOk, compressedBody)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprint(len(compressedBody)))
	w.Send()
}

func isPathSafe(basePath, userInput string) (string, bool, error) {
	fullPath := filepath.Join(basePath, userInput)

	// e.g. evaluates /tmp/.../etc/passwd -> /etc/passwd
	resolvedPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			resolvedPath = fullPath
		} else {
			return "", false, err
		}
	}

	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", false, err
	}

	absResolved, err := filepath.Abs(resolvedPath)
	if err != nil {
		return "", false, err
	}

	// Ensure absResolved is inside absBase (not just a prefix match)
	if !strings.HasPrefix(absResolved, absBase+string(os.PathSeparator)) && absBase != absResolved {
		return absResolved, false, nil
	}

	return absResolved, true, nil
}
