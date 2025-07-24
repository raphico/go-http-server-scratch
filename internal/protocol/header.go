package protocol

type Header map[string][]string

func (h Header) Set(key string, value string) {
	h[key] = append(h[key], value)
}

func (h Header) Get(key string) string {
	if values, ok := h[key]; ok && len(values) > 0 {
		return values[0]
	}

	return ""
}

func (h Header) Values(key string) []string {
	if values, ok := h[key]; ok && len(values) > 0 {
		return values
	}

	return nil
}
