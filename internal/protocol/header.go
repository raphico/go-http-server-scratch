package protocol

import "strings"

type Header map[string][]string

func canonical(key string) string {
	return strings.ToLower(key)
}

func (h Header) Set(key string, value string) {
	k := canonical(key)
	h[k] = append(h[k], value)
}

func (h Header) Get(key string) string {
	values := h.Values(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (h Header) Values(key string) []string {
	return h[canonical(key)]
}
